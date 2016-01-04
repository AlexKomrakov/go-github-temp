package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	url                = "localhost"
	builds_collection  = "builds"
	database           = "gohub"
)

type DeployScenario struct {
	Branch string
	Host   string
	Commands []map[string]string
    Error    []map[string]string
}

type Build struct {
	Id               bson.ObjectId 			   `bson:"_id,omitempty" json:"id"`
	Login 			 string                    `json:"login"`
	Name 			 string                    `json:"name"`
	SHA  			 string                    `json:"sha"`
	Event            string 				   `json:"event"`
	Created_at       time.Time 				   `json:"created_at"`
	DeployFile       map[string]DeployScenario `json:"deployScenario"`
	CommandResponses []CommandResponse		   `bson:"commandresponses,omitempty" json:"commandResponses"`
}
func (r Build) GetBuilds() (builds []Build, err error) {
	err = getDb().C(builds_collection).Find(bson.M{"login": r.Login, "name": r.Name}).Sort("-_id").All(&builds)
	return
}
func (b Build) HasError() bool {
	for _, val := range b.CommandResponses {
		if val.Error != "" {
			return true
		}
	}

	return false
}
func (b *Build) Store() (err error) {
	b.Created_at = time.Now()
	b.Id = bson.NewObjectId()

	return getDb().C(builds_collection).Insert(&b)
}
func (b *Build) AddCommand(c CommandResponse) (err error) {
	b.CommandResponses = append(b.CommandResponses, c)

	return getDb().C(builds_collection).Update(bson.M{"_id": b.Id}, bson.M{"$set": bson.M{"commandresponses": b.CommandResponses}})
}
func FindBuildById(id interface {}) (b Build, err error) {
	err = getDb().C(builds_collection).Find(bson.M{"_id": bson.ObjectIdHex(id.(string))}).One(&b)

	return
}

type CommandResponse struct {
	Type    string  `bson:"type,omitempty"`
	Command string  `bson:"command,omitempty"`
	Error   string  `bson:"error,omitempty"`
	Success string  `bson:"success,omitempty"`
}

//TODO defer session.Close()
func getDb() (db *mgo.Database) {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	db = session.DB(database)

	return
}