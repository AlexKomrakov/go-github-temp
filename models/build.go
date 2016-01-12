package models

import (
	"time"
)

type Build struct {
	Id               int64 			           `json:"id"`
	RepositoryId     int64					   `json:"repository_id" xorm:"index"`
	SHA  			 string                    `json:"sha"`
	Event            string 				   `json:"event"`
	// TODO add xorm created_at
	Created_at       time.Time 				   `json:"created_at"`
	// TODO save deploy file assigned to build
	//DeployFile       map[string]DeployScenario `json:"deployScenario"`
}

type CommandResponse struct {
	Id      int64   `json:"id"`
	BuildId int64   `json:"build_id" xorm:"index"`
	Type    string
	Command string
	Error   string
	Success string
}
func (b Build) CommandResponses() (command_responses []CommandResponse, err error) {
	// TODO check sorting order
	err = Orm.Find(&command_responses, &CommandResponse{BuildId: b.Id})
	return
}
func (b Build) HasError() bool {
	command_responses, err := b.CommandResponses()
	if err != nil {
		return false
	}
	for _, val := range command_responses {
		if val.Error != "" {
			return true
		}
	}

	return false
}
func (b *Build) Store() (err error) {
	// TODO add xorm created_at
	b.Created_at = time.Now()
	_, err = Orm.Insert(&b)

	return
}
func (b *Build) AddCommand(c CommandResponse) (err error) {
	c.BuildId = b.Id
	_, err = Orm.Insert(&c)

	return
}
func FindBuildById(id interface {}) (b Build, err error) {
	_, err = Orm.Id(id).Get(&b)

	return
}