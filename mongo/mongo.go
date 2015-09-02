package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/ssh"

	"strings"
	"time"
	"github.com/google/go-github/github"
)

const (
	url                = "localhost"
	repos_collection   = "repositories"
	servers_collection = "servers"
	builds_collection  = "builds"
	database           = "gohub"
	tokens_collection  = "tokens"
)

type Repository struct {
	github.Repository
	Id bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
}

type RepositoryCredentials struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

func (r Repository) Store() {
	err := getDb().C(repos_collection).Insert(&r)
	if err != nil {
		panic(err)
	}
}

//func (r Repository) Delete() {
//	err := getDb().C(repos_collection).Remove(r)
//	if err != nil {
//		panic(err)
//	}
//}

func FindRepository(q interface{}) (repo Repository, err error) {
	err = getDb().C(repos_collection).Find(q).One(&repo)

	return
}

func (r RepositoryCredentials) Find() (repo Repository, err error) {
	return FindRepository(bson.M{"repository.name": r.Name, "repository.owner.login": r.Login})
}


type Commit struct {
	Repository Repository `json:"repository"`
	SHA        string     `json:"sha"`
}

type Build struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Commit       Commit		   `json:"commit"`
	Start_time   time.Time     `json:"start_time"`
	End_time     time.Time     `json:"end_time"`
	Success      bool          `json:"success"`
}

type Server struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User      string `json:"user"`
	User_host string `json:"user_host"`
	Password  string `json:"password"`
	Checked   bool   `json:"checked"`
}

type Token struct {
	User  string `json:"user"`
	Token string `json:"token"`
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

func GetRepositories(user string) (repositories []Repository) {
	c := getDb().C(repos_collection)
	err := c.Find(bson.M{"user": user}).All(&repositories)
	if err != nil {
		panic(err)
	}

	return
}

func GetServers(user string) (servers []Server) {
	c := getDb().C(servers_collection)
	err := c.Find(bson.M{"user": user}).All(&servers)
	if err != nil {
		panic(err)
	}

	return
}

func (r Server) Store() {
	err := getDb().C(servers_collection).Insert(&r)
	if err != nil {
		panic(err)
	}
}

func (r Server) Delete() {
	err := getDb().C(servers_collection).Remove(r)
	if err != nil {
		panic(err)
	}
}

func (r Server) Find() (s Server) {
	err := getDb().C(servers_collection).Find(bson.M{"user": r.User, "user_host": r.User_host}).One(&s)
	if err != nil {
		panic(r)
	}

	return
}

func (s Server) Check() bool {
	_, err := GetSshClient(s.User_host, s.Password)

	return err == nil
}

func (s Server) Client() (client *ssh.Client, err error) {
	return GetSshClient(s.User_host, s.Password)
}

func (t Token) Store() {
	err := getDb().C(tokens_collection).Insert(&t)
	if err != nil {
		panic(err)
	}
}

func GetToken(user string) string {
	var t Token
	err := getDb().C(tokens_collection).Find(bson.M{"user": user}).One(&t)
	if err != nil {
		panic(err)
	}
	return t.Token
}

func GetSshClient(user_host, password string) (client *ssh.Client, err error) {
	params := strings.Split(user_host, "@")
	if len(params) != 2 {
		panic("Wrong ssh user@host: " + user_host)
	}
	user := params[0]
	host := params[1]

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err = ssh.Dial("tcp", host, config)

	return client, err
}
