// Package driver realize data read and write support for different databases(datastore provider)
//
//	All driver object must implement IDriver interface to support read/write data for IoT devices
//
//	As an excellent IoT data engineer, the IoT data storage structure you design should have the following parts:
//
//	Data collection time, data collection location, [] (data object name, data value)
//
//	At the same time, in order to comply with the table architecture specifications of the project, we will set the following fields for the stored data:
//
//	Table ID,
//
//	Node ID (this is the AgentID set by the user (usually you)) It is also your `data collection location`
//
//	Data name
//
//	Time
//
//	Its format is:
//
//	(id, node_id, [your dataColumns], ts)
//
//	You don't need to consider id, nodeId and ts, we have already designed it for you!
//	You just need to make sure that the data object you pass in conforms to the deserializable json object.
package driver

import (
	"errors"
	"io"
)

var ErrWriteTxnSize = errors.New("write txn size should be 1<= size <=100")

type IDriver interface {
	GetDriverName() string
	// SetTableName Set write/reader table name
	SetTableName(tableName string)
	// GetNewDriver Get a new writer/reader driver
	GetNewDriver(tableName string) IDriver
	// CreateColumnsWithType Create columns with type: key is column name, value is column type
	CreateColumnsWithType(map[string]string) error
	// SetBatchWriteTxnSize Set batch write transaction size: it should be 1<= size <=100.(default is 1)
	//  if size = 0,it will return currentSize,nil
	//  if size > 100,or size < 1,it will return 0,and error
	SetBatchWriteTxnSize(size int) (int, error)
	// Migrate auto migrate table
	Migrate(map[string]string) error
	IWriter
	IReader
}

type IWriter interface {
	// Insert insert data: key is column name, value is column value
	Insert(map[string]any) error
	// Writer will auto conversion data []byte to json-key-value, and write to table.
	io.Writer
}

type IReader interface {
	GetData(conditions ...ConditionOption) map[string][]any
}
