package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	url = "localhost"
	collection = "repositories"
	database = "gohub"
)

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
	c := getDb().C(collection)
	err := c.Insert(&repo)
	if err != nil {
		panic(err)
	}
}

func GetRepositories() []Repository {
	c := getDb().C(collection)
	var repositories []Repository
	err := c.Find(bson.M{}).All(&repositories)
	if err != nil {
		panic(err)
	}
	return repositories
}

func RemoveRepository(repo *Repository) {
	c := getDb().C(collection)
	err := c.Remove(repo)
	if err != nil {
		panic(err)
	}
}
