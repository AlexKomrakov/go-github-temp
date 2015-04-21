package mongo

import (
//	"code.google.com/p/goauth2/oauth"
//	"github.com/google/go-github/github"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	url               = "localhost"
	repos_collection  = "repositories"
	builds_collection = "builds"
	database          = "gohub"
	tokens_collection = "tokens"
)

type Repository struct {
	User       string `json:"user"`
	Repository string `json:"repository"`
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

func (r Repository) Store() {
	err := getDb().C(repos_collection).Insert(&r)
	if err != nil {
		panic(err)
	}
}

func (r Repository) Delete() {
	err := getDb().C(repos_collection).Remove(r)
	if err != nil {
		panic(err)
	}
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

//func (r *Repository) GetGithubClient() *github.Client {
//	transport := &oauth.Transport{
//		Token: &oauth.Token{AccessToken: r.Token},
//	}
//	return github.NewClient(transport.Client())
//}
//
//type Build struct {
//	Branch   Branch        `json:"branch,omitempty"`
//	Event    interface{}   `json:"event,omitempty"`
//	Commands []Command     `json:"commands,omitempty"`
//	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
//	Success  bool          `json:"success,omitempty"`
//	Time     int64         `json:"time,omitempty"`
//}
//
//func (b *Build) GetId() string {
//	return b.Id.Hex()
//}
//
//type Command struct {
//	Type   string `json:"type,omitempty"`
//	Action string `json:"action,omitempty"`
//	Out    string `json:"out,omitempty"`
//	Err    string `json:"err,omitempty"`
//}
//
//type Branch struct {
//	Owner string
//	Repo  string
//	Sha   string
//}
//
//func (b *Build) Save() {
//	c := getDb().C(builds_collection)
//	err := c.Insert(&b)
//	if err != nil {
//		panic(err)
//	}
//}
//
//func GetBuilds(user, repo string) (builds []Build) {
//	c := getDb().C(builds_collection)
//	err := c.Find(bson.M{"branch.owner": user, "branch.repo": repo}).All(&builds)
//	if err != nil {
//		panic(err)
//	}
//	return
//}
//
//func GetBuild(id string) (build Build) {
//	c := getDb().C(builds_collection)
//	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&build)
//	if err != nil {
//		panic(err)
//	}
//	return
//}
//
//func (b *Branch) GetRepository() (r *Repository) {
//	c := getDb().C(repos_collection)
//	c.Find(bson.M{"user": b.Owner, "repository": b.Repo}).One(&r)
//	return
//}
//

//
//func RemoveRepository(repo *Repository) {
//	c := getDb().C(repos_collection)
//	err := c.Remove(repo)
//	if err != nil {
//		panic(err)
//	}
//}
