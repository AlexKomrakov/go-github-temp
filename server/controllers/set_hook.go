package controllers
import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/google/go-github/github"
	"github.com/alexkomrakov/gohub/models"
	"github.com/alexkomrakov/gohub/service"
	"github.com/gorilla/mux"
)

func SetHook(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	token, _ := models.GetToken(user)
	client := service.GetGithubClient(token)

	url  := config.Url + "/hooks"
	hook := &github.Hook{Name: github.String("web"), Active: github.Bool(true), Events: config.Events, Config: map[string]interface {}{"url": url}}

	_, _, err := client.Repositories.CreateHook(params["user"], params["repo"], hook)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, req, "/repos/" + params["user"] + "/" + params["repo"], http.StatusTemporaryRedirect)
}
