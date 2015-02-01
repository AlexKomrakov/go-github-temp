package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/google/go-github/github"
	"code.google.com/p/goauth2/oauth"
)

const (
	url              = "localhost"
	repos_collection = "repositories"
	database         = "gohub"
)

type Repository struct {
	User       string `json:"user"`
	Repository string `json:"repository"`
	Token      string `json:"token"`
}

func (r *Repository) GetGithubClient() *github.Client {
	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: r.Token},
	}
	return github.NewClient(transport.Client())
}

type Build struct {
	Branch *Branch                  `json:"branch,omitempty"`
	Event  *github.PullRequestEvent `json:"event,omitempty"`
}

type Branch struct {
	Owner string
	Repo  string
	Sha   string
}

func (b *Branch) GetRepository() *Repository {
	c := getDb().C(repos_collection)
	var r *Repository
	c.Find(bson.M{"user": b.Owner, "repository": b.Repo}).One(&r)
	return r
}

//TODO defer session.Close()
func getDb() (*mgo.Database) {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	db := session.DB(database)
	return db
}

func AddRepository(repo *Repository) {
	c := getDb().C(repos_collection)
	err := c.Insert(&repo)
	if err != nil {
		panic(err)
	}
}

func GetRepositories() []Repository {
	c := getDb().C(repos_collection)
	var repositories []Repository
	err := c.Find(bson.M{}).All(&repositories)
	if err != nil {
		panic(err)
	}
	return repositories
}

func RemoveRepository(repo *Repository) {
	c := getDb().C(repos_collection)
	err := c.Remove(repo)
	if err != nil {
		panic(err)
	}
}
