package ts

type ITimeSeries interface {
	UsingDatabase(dbName string) error
	CreateTable(tableName string, columns ...interface{}) error
	InsertInto(tableName string, datas ...interface{}) error
	Query(tableName string, query string) ([]map[string]interface{}, error)
	GetLatestDataByColumns(tableName string, columns ...interface{}) ([]map[string]interface{}, error)
}
