package driver

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/bytedance/sonic"
)

type MySql struct {
	db        *sql.DB
	tableName string
	userID    string
	txnSize   int

	mutex  sync.Mutex
	buffer []map[string]any
}

func (MySql) GetDriverName() string {
	return "mysql"
}

func (m *MySql) SetTableName(tableName string) {
	m.tableName = tableName
}

func (m MySql) GetNewDriver(tableName string) IDriver {
	return &MySql{tableName: tableName, txnSize: m.txnSize, db: m.db, userID: m.userID}
}

func (m MySql) CreateColumnsWithType(data map[string]string) error {
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
		m.tableName,
		strings.Join(columnDefs, ", "),
	)

	if _, err := m.db.Exec(stmt); err != nil {
		logg.L.
			WithAction("DataStore-CreateTable").
			Errorf("failed to create table: %v", err)
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (m *MySql) SetBatchWriteTxnSize(size int) (int, error) {
	if size == 0 {
		return m.txnSize, nil
	}
	if size < 0 || size > 100 {
		return 0, ErrWriteTxnSize
	}
	m.txnSize = size
	return m.txnSize, nil
}

func (m MySql) Migrate(columns map[string]string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MySql) Insert(data map[string]any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.buffer = append(m.buffer, data)

	if len(m.buffer) < m.txnSize {
		return nil
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	var (
		placeholders []string
		columns      []string
		values       []any
	)

	for k := range m.buffer[0] {
		columns = append(columns, k)
	}

	for _, record := range m.buffer {
		for _, col := range columns {
			values = append(values, record[col])
		}
		placeholders = append(placeholders, "("+strings.Repeat("?,", len(columns)-1)+"?)")
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		m.tableName,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","),
	)

	if _, err := tx.Exec(stmt, values...); err != nil {
		tx.Rollback()
		return fmt.Errorf("batch insert failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	m.buffer = nil

	return nil
}

func (m *MySql) Write(p []byte) (n int, err error) {
	var data map[string]any
	err = sonic.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}
	return len(p), m.Insert(data)
}

func (m MySql) GetData(conditions ...ConditionOption) map[string][]any {
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
			id, err := nodeNameGetNodeId(*qc.NodeName, m.userID)
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

	query := fmt.Sprintf("SELECT * FROM %s", m.tableName)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return make(map[string][]any)
	}
	defer rows.Close()

	result := make(map[string][]any)
	columns, _ := rows.Columns()

	for _, col := range columns {
		result[col] = []any{}
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return make(map[string][]any)
		}

		for i, val := range values {
			result[columns[i]] = append(result[columns[i]], val)
		}
	}

	return result
}

func NewMySQLDriver(dsn, tableName string, userID string) (IDriver, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logg.L.
			WithAction("DataStore-LinkToMySQL").
			Errorf("failed to open mysql link:%v", err)
		return nil, err
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	return &MySql{db: db, txnSize: 1, tableName: tableName, userID: userID}, nil
}
