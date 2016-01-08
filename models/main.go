package models

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Db_file = "sqllite.db"
	Orm *xorm.Engine
)

func init() {
	// TODO Вынести в конфиг-файл
	orm, err := xorm.NewEngine("sqlite3", Db_file)
	if err != nil {
		panic(err)
	}
	Orm = orm

	err = Orm.CreateTables(&Build{}, &CommandResponse{}, &Server{}, &Token{})
	if err != nil {
		panic(err)
	}
}