package controllers
import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
	"github.com/alexkomrakov/gohub/mongo"
	"github.com/alexkomrakov/gohub/service"
)

func RunScenario(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	token, _ := models.GetToken(user)
	client := service.GetGithubClient(token)
	file, _ := service.GetFileContent(client, params["user"], params["repo"], params["sha"], config.DeployFile)
	string_file := service.ReplaceVariables(params, string(file))
	deploy, _ := service.GetYamlConfig([]byte(string_file))
	build := service.RunCommands(deploy, client, params["scenario"], mongo.CommitCredentials{mongo.RepositoryCredentials{params["user"], params["repo"]}, params["sha"]})

	http.Redirect(res, req, "/repos/"+params["user"]+"/"+params["repo"]+"/build/"+build.Id.Hex(), http.StatusFound)
}