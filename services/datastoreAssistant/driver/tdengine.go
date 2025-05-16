package driver

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"sync"

	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"github.com/bytedance/sonic"
	"github.com/taosdata/driver-go/v3/af"
	_ "github.com/taosdata/driver-go/v3/taosSql"
)

type TDEngine struct {
	td               *af.Connector
	schemaless       bool
	maxWriteLineSize int
	table            string
	//insterPool       *ants.Pool
	userId string

	mutex  sync.Mutex
	buffer []map[string]any
}

func (TDEngine) GetDriverName() string {
	return "TDEngine"
}

func (t *TDEngine) SetTableName(tableName string) {
	t.table = tableName
}

func (t TDEngine) GetNewDriver(tableName string) IDriver {
	return &TDEngine{
		schemaless:       t.schemaless,
		td:               t.td,
		maxWriteLineSize: 100,
		table:            tableName,
		buffer:           make([]map[string]any, 0),
		//insterPool:       t.insterPool,
		userId: t.userId,
	}
}

func (t TDEngine) CreateColumnsWithType(data map[string]string) error {
	var (
		columns []string
		types   []string
	)
	for k, v := range data {
		columns = append(columns, k)
		types = append(types, v)
	}
	var columnDefs []string
	for i := 0; i < len(columns); i++ {
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s", columns[i], types[i]))
	}

	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);",
		t.table,
		strings.Join(columnDefs, ", "),
	)

	if _, err := t.td.Exec(stmt); err != nil {
		logg.L.
			WithAction("DataStore-CreateTable").
			Errorf("failed to create table: %v", err)
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (t *TDEngine) SetBatchWriteTxnSize(size int) (int, error) {
	if size == 0 {
		return t.maxWriteLineSize, nil
	}
	if size < 0 || size > 100 {
		return 0, ErrWriteTxnSize
	}
	t.maxWriteLineSize = size
	return t.maxWriteLineSize, nil
}

func (t TDEngine) Migrate(m map[string]string) error {
	//TODO implement me
	panic("implement me")
}
func (TDEngine) anys2DriverValues(m []any) []driver.Value {
	var retVal []driver.Value
	for _, v := range m {
		retVal = append(retVal, v)
	}
	return retVal
}
func (t *TDEngine) Insert(m map[string]any) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.maxWriteLineSize != 1 {
		t.buffer = append(t.buffer, m)

		if len(t.buffer) < t.maxWriteLineSize {
			return nil
		} else {
			defer func() {
				t.buffer = make([]map[string]any, 0)
			}()
		}
	}

	if t.schemaless {
		if t.maxWriteLineSize == 1 {
			line, _ := sonic.MarshalString(m)
			err := t.td.OpenTSDBInsertJsonPayload(line)
			if err != nil {
				logg.L.Errorf("failed to insert schemaless data:%v", err)
			}
		} else {
			for _, m := range t.buffer {
				line, _ := sonic.MarshalString(m)
				err := t.td.OpenTSDBInsertJsonPayload(line)
				if err != nil {
					logg.L.Errorf("failed to insert schemaless data:%v", err)
				}
			}
		}
		return nil
	}

	var (
		placeholders []string
		columns      []string
		values       []any
	)

	for k := range t.buffer[0] {
		columns = append(columns, k)
	}

	for _, record := range t.buffer {
		for _, col := range columns {
			values = append(values, record[col])
		}
		placeholders = append(placeholders, "("+strings.Repeat("?,", len(columns)-1)+"?)")
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		t.table,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","),
	)
	if _, err := t.td.Exec(stmt, t.anys2DriverValues(values)...); err != nil {
		logg.L.Errorf("failed to insert data:%v", err)
	}
	return nil
}

func (t *TDEngine) Write(p []byte) (n int, err error) {
	var data map[string]any
	err = sonic.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}
	return len(p), t.Insert(data)
}

func (t TDEngine) GetData(conditions ...ConditionOption) map[string][]any {
	qc := &QueryCondition{}
	for _, opt := range conditions {
		opt(qc)
	}

	var whereClauses []string
	var args []interface{}

	if qc.ConTimeRange != nil {
		if qc.StartAt != nil {
			whereClauses = append(whereClauses, "ts >= ?")
			args = append(args, *qc.StartAt)
		}
		if qc.EndAt != nil {
			whereClauses = append(whereClauses, "ts <= ?")
			args = append(args, *qc.EndAt)
		}
	}

	if qc.ConNode != nil {
		if qc.NodeName != nil {
			id, err := nodeNameGetNodeId(*qc.NodeName, t.userId)
			if err != nil {
				logg.L.
					WithAction("DataStore-NodeNameGetNodeId").
					Errorf("failed to get node id by nodeName:%v", err)
				return nil
			}
			qc.NodeID = &id
		}
		if qc.NodeID != nil {
			whereClauses = append(whereClauses, "node_id = ?")
			args = append(args, *qc.NodeID)
		}
	}

	for _, vr := range qc.ValRange {
		clause := vr.Key
		if vr.StartVal != nil && vr.EndVal != nil {
			clause += " BETWEEN ? AND ?"
			args = append(args, *vr.StartVal, *vr.EndVal)
		} else if vr.StartVal != nil {
			clause += " >= ?"
			args = append(args, *vr.StartVal)
		} else if vr.EndVal != nil {
			clause += " <= ?"
			args = append(args, *vr.EndVal)
		}
		whereClauses = append(whereClauses, clause)
	}

	query := fmt.Sprintf("SELECT * FROM %s", t.table)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := t.td.Query(query, t.anys2DriverValues(args)...)
	if err != nil {
		return make(map[string][]any)
	}
	defer rows.Close()

	result := make(map[string][]any)
	columns := rows.Columns()

	for _, col := range columns {
		result[col] = []any{}
	}

	var dest []driver.Value
	for rows.Next(dest) != io.EOF {
		for i, column := range columns {
			result[column] = append(result[column], dest[i])
		}
	}

	return result
}

func NewTDEngineDriver(userId, host, username, password, db string,
	port int, table *string, useSchemaless *bool) (IDriver, func() error, error) {
	conn, err := af.Open(host, username, password, db, port)
	if err != nil {
		logg.L.Errorf("failed to open tdengine link:%v", err)
		return nil, nil, err
	}
	if table == nil {
		t := "meters"
		table = &t
	}
	//ants.NewPool(2048, ants.WithPreAlloc(true))
	retval := TDEngine{
		td:               conn,
		maxWriteLineSize: 1,
		buffer:           make([]map[string]any, 0),
		userId:           userId,
		table:            *table,
	}
	if useSchemaless != nil {
		retval.schemaless = *useSchemaless
	}
	return &retval, conn.Close, nil
}
