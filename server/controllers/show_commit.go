package controllers
import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
	"github.com/alexkomrakov/gohub/service"
)

func ShowCommit(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	token, _ := models.GetToken(user)
	client := service.GetGithubClient(token)
	repo, _, _ := client.Repositories.Get(params["user"], params["repo"])
	commit, _, _ := client.Repositories.GetCommit(params["user"], params["repo"], params["sha"])
	file, _ := service.GetFileContent(client, params["user"], params["repo"], params["sha"], config.DeployFile)
	deploy, _ := service.GetYamlConfig(file)

	Render(res, req, "commit", map[string]interface{}{"Repo": repo, "Commit": commit, "File": string(file), "Deploy": deploy})
}
