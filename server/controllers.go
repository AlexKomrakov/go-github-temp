package server

import (
	"github.com/alexkomrakov/gohub/mongo"
    "github.com/google/go-github/github"

    "github.com/gorilla/mux"
	"github.com/unrolled/render"
    "github.com/ActiveState/tail"

	"net/http"
    "log"
    "github.com/alexkomrakov/gohub/service"
    "os"

    "github.com/goincremental/negroni-sessions"

    "github.com/gorilla/schema"
    "strings"
    "strconv"

    "golang.org/x/oauth2"
    githuboauth "golang.org/x/oauth2/github"
)

var (
    r *render.Render
    l *log.Logger
    config service.ServerConfig

    // You must register the app at https://github.com/settings/applications
    // Set callback to http://127.0.0.1:7000/github_oauth_cb
    // Set ClientId and ClientSecret to
    oauthConf oauth2.Config
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString string
)

type ViewData struct {
    User  interface{} `json:"user,omitempty"`
    Token interface{} `json:"token,omitempty"`
    Data  interface{} `json:"data,omitempty"`
}

func Render(res http.ResponseWriter, req *http.Request, view string, data interface {}) {
    session := sessions.GetSession(req)
    user := session.Get("user")

    r.HTML(res, http.StatusOK, view, ViewData{User: user, Data: data})
}

func init() {
    config = service.GetServerConfig()
	r = render.New(render.Options{
        Layout: "base",
        IsDevelopment: true,
    })
    l = service.GetFileLogger(config.Logs.Gohub)

    oauthConf = config.Oauth
    oauthConf.Endpoint = githuboauth.Endpoint
    oauthStateString = config.OauthStateString
}

func Index(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user")
    if user == nil {
        http.Redirect(res, req, "/login", http.StatusFound)
    }
    http.Redirect(res, req, "/repos/" + user.(string), http.StatusFound)

    //    Render(res, req, "index", nil)
}

func UserRepos(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    github_repos, _, _ := client.Repositories.List("", nil)
    gohub_repos := mongo.GetRepositories(user)
    github_repos = filterRepos(github_repos, gohub_repos)

    Render(res, req, "repos", map[string]interface{}{"Github": github_repos, "Gohub": gohub_repos})
}

func filterRepos(github []github.Repository, database []mongo.Repository) (output []github.Repository) {
    for _, value := range github {
        found := false
        for _, db_value := range database {
            if (*value.Owner.Login == db_value.User && *value.Name == db_value.Repository) {
                found = true
            }
        }
        if (found == false) {
            output = append(output, value)
        }
    }
    return
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

func ShowRepo(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    repo, _, _ := client.Repositories.Get(user, params["repo"])
    hooks, _, _ := client.Repositories.ListHooks(user, params["repo"], &github.ListOptions{})
    branches, _, _ := client.Git.ListRefs(user, params["repo"], &github.ReferenceListOptions{})
    var filtered_branches []github.Reference
    for _, branch := range branches {
        if strings.HasPrefix(*branch.Ref, "refs/heads/") {
            filtered_branches = append(filtered_branches, branch)
        }
    }

    Render(res, req, "repo", map[string]interface{}{"Repo": repo, "Hooks": hooks, "Branches": filtered_branches})
}

func ShowCommit(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    repo, _, _ := client.Repositories.Get(user, params["repo"])
    commit, _, _ := client.Repositories.GetCommit(user, params["repo"], params["sha"])
    file, _ := service.GetFileContent(client, user, params["repo"], params["sha"], config.DeployFile)
    deploy, _ := service.GetYamlConfig(file)

    Render(res, req, "commit", map[string]interface{}{"Repo": repo, "Commit": commit, "File": string(file), "Deploy": deploy})
}

func RunScenario(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    file, _ := service.GetFileContent(client, user, params["repo"], params["sha"], config.DeployFile)
    deploy, _ := service.GetYamlConfig(file)
    commands := service.RunCommands(deploy[params["scenario"]], client, user, params["repo"], params["sha"])

    Render(res, req, "run", map[string]interface{}{"Params": params, "Commands": commands})
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
    process_login(res, req, user, token)

    Render(res, req, "login", nil)
}

func process_login(res http.ResponseWriter, req *http.Request, user *github.User, token string) {
    if user != nil {
        session := sessions.GetSession(req)
        session.Set("user", user.Login)
        mongo.Token{User: *user.Login, Token: token}.Store()

        http.Redirect(res, req, "/", http.StatusFound)
    }
}

func Logs(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    t, err := tail.TailFile(params["name"] + ".log", tail.Config{Follow: false, Location: &tail.SeekInfo{0, os.SEEK_SET}})
    if err != nil {
        panic(err)
    }

    Render(res, req, "logs", t)
}

func GithubHookApi(w http.ResponseWriter, req *http.Request) {
    body := req.FormValue("payload")
    event := req.Header["X-Github-Event"][0]
    l.Println(event)
    service.ProcessHook(event, body)
}

func GithubLogin(w http.ResponseWriter, req *http.Request) {
    url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
    http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func GithubLoginCallback(w http.ResponseWriter, req *http.Request) {
    state := req.FormValue("state")
    if state != oauthStateString {
        l.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
        http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
        return
    }

    code := req.FormValue("code")
    token, err := oauthConf.Exchange(oauth2.NoContext, code)
    if err != nil {
        l.Printf("oauthConf.Exchange() failed with '%s'\n", err)
        http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
        return
    }

    oauthClient := oauthConf.Client(oauth2.NoContext, token)
    client := github.NewClient(oauthClient)
    user, _, err := client.Users.Get("")
    if err != nil {
        l.Printf("client.Users.Get() faled with '%s'\n", err)
        http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
        return
    }

    process_login(w, req, user, token.AccessToken)
    http.Redirect(w, req, "/login", http.StatusFound)
}

func SetHook(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)

    url  := config.Url + "/hooks"
    hook := &github.Hook{Name: github.String("web"), Active: github.Bool(true), Events: config.Events, Config: map[string]interface {}{"url": url}}

    client.Repositories.CreateHook(user, params["repo"], hook)

    http.Redirect(w, req, "/repos/" + user + "/" + params["repo"], http.StatusTemporaryRedirect)
}

func DeleteHook(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)

    id, _ := strconv.Atoi(params["id"])
    client.Repositories.DeleteHook(user, params["repo"], id)

    http.Redirect(w, req, "/repos/" + user + "/" + params["repo"], http.StatusTemporaryRedirect)
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
