package ts

import (
	"database/sql"
	"fmt"

	"github.com/taosdata/driver-go/v3/af"
	_ "github.com/taosdata/driver-go/v3/taosSql" // need cgo enabled
)

type Ts struct {
	Host string
	Port int
	User string
	Pass string
}
type TDengineDSN struct {
	Host string
	Port int

	UserName string
	Password string
	Protocol string
	Address  string
	DbName   string
	Params   map[string]string
}

func (d TDengineDSN) String() string {
	// [username[:password]@][protocol[(address)]]/[dbname][?param1=value1&...&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?", d.UserName, d.Password, d.Protocol, d.Address, d.DbName)
	for s, s2 := range d.Params {
		dsn = fmt.Sprintf("%s&%s=%s", dsn, s, s2)
	}
	return dsn
}

type TDengine struct {
	db *sql.DB
	*af.Connector
}

func NewTDengine(dsn *TDengineDSN) (*TDengine, error) {
	var retVal = &TDengine{}
	connector, err := af.Open(dsn.Host, dsn.UserName, dsn.Password, dsn.DbName, dsn.Port)
	if err != nil {
		return nil, err
	}
	retVal.Connector = connector
	return retVal, nil
}
