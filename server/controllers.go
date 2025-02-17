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
}

func UserRepos(res http.ResponseWriter, req *http.Request) {
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    github_repos, _, _ := client.Repositories.List("", nil)

    Render(res, req, "repos", map[string]interface{}{"Github": github_repos})
}

func ShowRepo(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    repo, _, _ := client.Repositories.Get(params["user"], params["repo"])
    repo_cred := mongo.RepositoryCredentials{params["user"], params["repo"]}
    _, err := repo_cred.GetRepository()
    if err != nil {
        mongo.Repository{Repository:*repo}.Store()
    }
    builds, _ := repo_cred.GetBuilds()
    hooks, _, _ := client.Repositories.ListHooks(params["user"], params["repo"], &github.ListOptions{})
    branches, _, _ := client.Git.ListRefs(params["user"], params["repo"], &github.ReferenceListOptions{})
    var filtered_branches []github.Reference
    for _, branch := range branches {
        if strings.HasPrefix(*branch.Ref, "refs/heads/") {
            filtered_branches = append(filtered_branches, branch)
        }
    }

    Render(res, req, "repo", map[string]interface{}{"Repo": repo, "Hooks": hooks, "Branches": filtered_branches, "Builds": builds})
}

func ShowBuild(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    build, _ := mongo.FindBuildById(params["id"])

    Render(res, req, "build", map[string]interface{}{"Params": params, "Build": build})
}

func ShowCommit(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    repo, _, _ := client.Repositories.Get(params["user"], params["repo"])
    commit, _, _ := client.Repositories.GetCommit(params["user"], params["repo"], params["sha"])
    file, _ := service.GetFileContent(client, params["user"], params["repo"], params["sha"], config.DeployFile)
    deploy, _ := service.GetYamlConfig(file)

    Render(res, req, "commit", map[string]interface{}{"Repo": repo, "Commit": commit, "File": string(file), "Deploy": deploy})
}

func RunScenario(res http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)
    file, _ := service.GetFileContent(client, params["user"], params["repo"], params["sha"], config.DeployFile)
    string_file := service.ReplaceVariables(params, string(file))
    deploy, _ := service.GetYamlConfig([]byte(string_file))
    build := service.RunCommands(deploy, client, params["scenario"], mongo.CommitCredentials{mongo.RepositoryCredentials{params["user"], params["repo"]}, params["sha"]})

    http.Redirect(res, req, "/repos/"+params["user"]+"/"+params["repo"]+"/build/"+build.Id.Hex(), http.StatusFound)
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
    l.Println(body)
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

    _, _, err := client.Repositories.CreateHook(params["user"], params["repo"], hook)
    if err != nil {
        panic(err)
    }

    http.Redirect(w, req, "/repos/" + params["user"] + "/" + params["repo"], http.StatusTemporaryRedirect)
}

func DeleteHook(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    session := sessions.GetSession(req)
    user := session.Get("user").(string)
    token := mongo.GetToken(user)
    client := service.GetGithubClient(token)

    id, _ := strconv.Atoi(params["id"])
    client.Repositories.DeleteHook(params["user"], params["repo"], id)

    http.Redirect(w, req, "/repos/" + params["user"] + "/" + params["repo"], http.StatusTemporaryRedirect)
}
