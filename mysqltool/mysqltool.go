package mysqltool

import (
	"database/sql"
	"fmt"
)

type Mysqltool struct {
	name     string
	password string
	host     string
	port     int
	dBname   string
	db       *sql.DB
}

func NewMysqlTool(name, password, host string, port int, dbName string) *Mysqltool {
	return &Mysqltool{
		name:     name,
		password: password,
		host:     host,
		port:     port,
		dBname:   dbName,
	}
}

func (mt *Mysqltool) Connect() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		mt.name,
		mt.password,
		mt.host,
		mt.port,
		mt.dBname,
	)
	mt.db, err = sql.Open("tcp", dsn)
	if err != nil {
		return err
	}
	return mt.db.Ping()
}

func (mt *Mysqltool) Close() error {
	return mt.db.Close()
}

// func (mt *Mysqltool) GetCreatTable(table string) (string, error) {
// 	r, err := mt.db.Query("show creat table ?;", table)
// 	if err != nil {
// 		return "", err
// 	}
// 	s, err := r.Columns()
// 	fmt.Println(r.Scan())
// }

// func (mt *Mysqltool) GetQueryRows(rows *sql.Rows) (map[string]interface{}, error) {
// 	s, err := rows.Columns()
// 	if err != nil {
// 		return nil, err
// 	}

// }
