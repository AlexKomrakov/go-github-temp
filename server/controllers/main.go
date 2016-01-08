package controllers

import (
	"github.com/alexkomrakov/gohub/models"
    "github.com/google/go-github/github"
	"github.com/unrolled/render"
	"net/http"
    "log"
    "github.com/alexkomrakov/gohub/service"
    "github.com/goincremental/negroni-sessions"
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

func process_login(res http.ResponseWriter, req *http.Request, user *github.User, token string) {
    if user != nil {
        session := sessions.GetSession(req)
        session.Set("user", user.Login)
        token := models.Token{User: *user.Login, Token: token}
        token.Store()

        http.Redirect(res, req, "/", http.StatusFound)
    }
}
