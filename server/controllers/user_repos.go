package controllers

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/service"
	"github.com/alexkomrakov/gohub/models"
	"fmt"
)

func UserRepos(res http.ResponseWriter, req *http.Request) {
	session  := sessions.GetSession(req)
	user     := session.Get("user").(string)

	token_model := models.Token{User: user}
	result, err := token_model.FindOne()
	if result == false || err != nil {
		panic("Cant find user token")
	}
	token := token_model.Token
	client   := service.GetGithubClient(token)
	github_repos, _, _ := client.Repositories.List("", nil)

	Render(res, req, "repos", map[string]interface{}{"Github": github_repos})
}