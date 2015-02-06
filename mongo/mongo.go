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
	builds_collection = "builds"
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
	Branch   Branch                  `json:"branch,omitempty"`
	Event    github.PullRequestEvent `json:"event,omitempty"`
	Commands []Command 			  `json:"commands,omitempty"`
}

type Command struct {
	Type   string
	Action string
	Out    string
	Err    string
}

type Branch struct {
	Owner string
	Repo  string
	Sha   string
}

func (b *Build) Save() {
	c := getDb().C(builds_collection)
	err := c.Insert(&b)
	if err != nil {
		panic(err)
	}
}

func GetBuilds(user, repo string) (builds []Build) {
	c := getDb().C(builds_collection)
	err := c.Find(bson.M{"branch.owner": user, "branch.repo": repo}).All(&builds)
	if err != nil {
		panic(err)
	}
	return
}

func (b *Branch) GetRepository() (r *Repository) {
	c := getDb().C(repos_collection)
	c.Find(bson.M{"user": b.Owner, "repository": b.Repo}).One(&r)
	return
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
