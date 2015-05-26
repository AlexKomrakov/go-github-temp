package server

import (
	"github.com/alexkomrakov/gohub/mongo"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
    "github.com/ActiveState/tail"

	"net/http"
    "log"
    "github.com/alexkomrakov/gohub/service"
    "os"

    "github.com/goincremental/negroni-sessions"

    "github.com/gorilla/schema"
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

    r.HTML(res, http.StatusOK, view, ViewData{User: user, Data: data})
}

func init() {
    config := service.GetServerConfig()
	r = render.New(render.Options{Layout: "base"})
    l = service.GetFileLogger(config.Logs.Gohub)
}

func Index(res http.ResponseWriter, req *http.Request) {
    Render(res, req, "index", nil)
}

func UserRepos(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    github_repos, _, _ := client.Repositories.List("", nil)
    gohub_repos := mongo.GetRepositories(user)

    Render(res, req, "repos", map[string]interface{}{"Github": github_repos, "Gohub": gohub_repos})
}

func AddRepo(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)

    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    repo, _, _ := client.Repositories.Get(user, params["repo"])
    if repo != nil {
        mongo.Repository{User: *repo.Owner.Login, Repository: *repo.Name}.Store()
        http.Redirect(res, req, "/repos/"+*repo.Owner.Login, http.StatusFound)
    } else {
        panic("Can't find repo")
    }
}

func DeleteRepo(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    mongo.Repository{User: user, Repository: params["repo"]}.Delete()

    http.Redirect(res, req, "/repos/"+user, http.StatusFound)
}

func UserServers(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    servers := mongo.GetServers(user)

    Render(res, req, "servers", map[string]interface{}{"Servers": servers})
}

func AddServer(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    req.ParseForm()

    server := mongo.Server{User: user, User_host: req.PostFormValue("user_host"), Password: req.PostFormValue("password") }
    server.Checked = server.Check()
    server.Store()

    http.Redirect(res, req, "/servers/"+user, http.StatusFound)
}

func DeleteServer(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    req.ParseForm()

    server := new(mongo.Server)
    schema.NewDecoder().Decode(server, req.PostForm)
    log.Println(server)
    server.Delete()

    http.Redirect(res, req, "/servers/"+user, http.StatusFound)
}

func Logout(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    session.Delete("user")

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
        mongo.Token{User: *user.Login, Token: token}.Store()

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
