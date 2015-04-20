package server

import (
//	"github.com/alexkomrakov/gohub/mongo"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
    "github.com/ActiveState/tail"

	"net/http"
    "log"
    "github.com/alexkomrakov/gohub/service"
    "os"

    "github.com/goincremental/negroni-sessions"
)

var r *render.Render
var l *log.Logger

type ViewData struct {
    User  interface{} `json:"user,omitempty"`
    Token interface{} `json:"token,omitempty"`
    Data  interface{} `json:"data,omitempty"`
}

func Render(res http.ResponseWriter, req *http.Request, view string, data interface {}) {
    session := sessions.GetSession(req)
    user  := session.Get("user")
    token := session.Get("token")

    r.HTML(res, http.StatusOK, view, ViewData{User: user, Token: token, Data: data})
}

func init() {
    config := service.GetServerConfig()
	r = render.New(render.Options{Layout: "base"})
    l = service.GetFileLogger(config.Logs.Gohub)
}

func Index(res http.ResponseWriter, req *http.Request) {
    Render(res, req, "index", nil)
}

func Logout(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    session.Delete("user")
    session.Delete("token")

    http.Redirect(res, req, "/", http.StatusFound)
}

func Login(res http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    token := req.FormValue("token")

    client := service.GetGithubClient(token)
    user, _, _ := client.Users.Get("")
    if user != nil {
        session := sessions.GetSession(req)
        session.Set("user", user.Login)
        session.Set("token", token)
        http.Redirect(res, req, "/", http.StatusFound)
    }

    Render(res, req, "login", nil)
}

func Logs(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    t, err := tail.TailFile(params["name"] + ".log", tail.Config{Follow: false, Location: &tail.SeekInfo{0, os.SEEK_SET}})
    if err != nil {
        panic(err)
    }

    Render(res, req, "logs", t)
}

//func SetHook(res http.ResponseWriter, req *http.Request) {
//	params := mux.Vars(req)
//	result := setGithubHook(params["user"], params["repo"])
//	r.JSON(res, http.StatusOK, result)
//}
//
//func GetReposApi(res http.ResponseWriter, req *http.Request) {
//	repositories := mongo.GetRepositories()
//	r.JSON(res, http.StatusOK, repositories)
//}
//
//func RepoPage(res http.ResponseWriter, req *http.Request) {
//	data := make(map[string]interface{})
//	params := mux.Vars(req)
//	data["params"] = params
//	data["builds"] = mongo.GetBuilds(params["user"], params["repo"])
//	r.HTML(res, http.StatusOK, "repo", data)
//}
//
//func GithubHookApi(w http.ResponseWriter, req *http.Request) {
//	body := req.FormValue("payload")
//	event := req.Header["X-Github-Event"][0]
//	ProcessHook(event, body)
//}
//
//func PostReposApi(res http.ResponseWriter, req *http.Request) {
//	var repo mongo.Repository
//	decoder := json.NewDecoder(req.Body)
//	err := decoder.Decode(&repo)
//	if err != nil {
//		panic(err)
//	}
//	mongo.AddRepository(&repo)
//	r.JSON(res, http.StatusOK, nil)
//}
//
//func BuildPage(res http.ResponseWriter, req *http.Request) {
//	data := make(map[string]interface{})
//	params := mux.Vars(req)
//	data["params"] = params
//	data["builds"] = mongo.GetBuilds(params["user"], params["repo"])
//	data["build"]  = mongo.GetBuild(params["build"])
//	r.HTML(res, http.StatusOK, "repo", data)
//}
