package controllers
import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/service"
	"github.com/alexkomrakov/gohub/models"
)

func UserRepos(res http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	token, _ := models.GetToken(user)
	client := service.GetGithubClient(token)
	github_repos, _, _ := client.Repositories.List("", nil)

	Render(res, req, "repos", map[string]interface{}{"Github": github_repos})
}