package server

import (
	"github.com/alexkomrakov/gohub/mongo"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"encoding/json"
	"net/http"
)

const deploy_file = ".deploy.yml"

var r *render.Render

func init() {
	r = render.New(render.Options{Layout: "base"})
}

func Index(res http.ResponseWriter, req *http.Request) {
	r.HTML(res, http.StatusOK, "index", nil)
}

func GetReposApi(res http.ResponseWriter, req *http.Request) {
	repositories := mongo.GetRepositories()
	r.JSON(res, http.StatusOK, repositories)
}

func RepoPage(res http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	params := mux.Vars(req)
	data["params"] = params
	data["builds"] = mongo.GetBuilds(params["user"], params["repo"])
	r.HTML(res, http.StatusOK, "repo", data)
}

func GithubHookApi(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")
	event := req.Header["X-Github-Event"][0]
	ProcessHook(event, body)
}

func PostReposApi(res http.ResponseWriter, req *http.Request) {
	var repo mongo.Repository
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&repo)
	if err != nil {
		panic(err)
	}
	mongo.AddRepository(&repo)
	r.JSON(res, http.StatusOK, nil)
}

func BuildPage(res http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	params := mux.Vars(req)
	data["params"] = params
	data["builds"] = mongo.GetBuilds(params["user"], params["repo"])
	data["build"] = mongo.GetBuild(params["build"])
	r.HTML(res, http.StatusOK, "repo", data)
}
