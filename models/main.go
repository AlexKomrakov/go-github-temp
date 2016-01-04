package models

import (
	"github.com/google/go-github/github"
	"time"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Db_file = "sqllite.db"
	Orm *xorm.Engine
)

type DeployScenario struct {
	Branch string
	Host   string
	Commands []map[string]string
	Error    []map[string]string
}

type Repository struct {
	github.Repository
	Id interface{} `bson:"_id,omitempty" json:"id,omitempty"`
}

type RepositoryCredentials struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type CommitCredentials struct {
	RepositoryCredentials
	SHA string `json:"sha"`
}

type Build struct {
	CommitCredentials
	Id               interface{} 			   `bson:"_id,omitempty" json:"id"`
	DeployFile       map[string]DeployScenario `json:"deployScenario"`
	Event            string 				   `json:"event"`
	Created_at       time.Time 				   `json:"created_at"`
	CommandResponses []CommandResponse		   `bson:"commandresponses,omitempty" json:"commandResponses"`
}

type CommandResponse struct {
	Type    string  `bson:"type,omitempty"`
	Command string  `bson:"command,omitempty"`
	Error   string  `bson:"error,omitempty"`
	Success string  `bson:"success,omitempty"`
}

func init() {
	// TODO Вынести в конфиг-файл
	orm, err := xorm.NewEngine("sqlite3", Db_file)
	if err != nil {
		panic(err)
	}
	Orm = orm
}