package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type Repository struct {
	User       string `json:"user"`
	Repository string `json:"repository"`
	Token      string `json:"token"`
}

func GetReposApi (r render.Render) {
	result := Repository{"user1", "passwordhash", "/user1"}
	r.JSON(200, result)
}

func Index (r render.Render) {
	r.HTML(200, "index", "golang")
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Layout: "base"}))

	m.Get("/", Index)
	m.Get("/repos", GetReposApi)
	m.Run()
}
