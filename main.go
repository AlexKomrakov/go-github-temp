package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/alexkomrakov/gohub/mongo"
	"net/http"
	"encoding/json"
)

func GetReposApi (r render.Render) {
	repositories := mongo.GetRepositories()
	r.JSON(200, repositories)
}

func PostReposApi (res http.ResponseWriter, req *http.Request, r render.Render) {
	decoder := json.NewDecoder(req.Body)
	var repo mongo.Repository
	err := decoder.Decode(&repo)
	if err != nil {
		panic(err)
	}
	mongo.AddRepository(&repo)
	r.JSON(200, "{status: 'ok'}")
}

func Index (r render.Render) {
	r.HTML(200, "index", "golang")
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Layout: "base"}))

	m.Get("/", Index)
	m.Get("/repos", GetReposApi)
	m.Post("/repos", PostReposApi)
	m.Run()
}
