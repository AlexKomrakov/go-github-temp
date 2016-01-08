package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
	"github.com/alexkomrakov/gohub/service"
	"strconv"
)

func RunScenario(res http.ResponseWriter, req *http.Request) {
	params      := mux.Vars(req)
	session     := sessions.GetSession(req)
	user        := session.Get("user").(string)

	token_model := models.Token{User: user}
	token_model.FindOne()
	token := token_model.Token

	client      := service.GetGithubClient(token)
	file, _     := service.GetFileContent(client, params["user"], params["repo"], params["sha"], config.DeployFile)
	string_file := service.ReplaceVariables(params, string(file))
	deploy, _   := service.GetYamlConfig([]byte(string_file))
	build       := service.RunCommands(deploy, client, params["scenario"], models.Build{Login: params["user"], Name: params["repo"], SHA: params["sha"]})

	http.Redirect(res, req, "/repos/" + params["user"] + "/" + params["repo"] + "/build/" + strconv.FormatInt(build.Id, 10), http.StatusFound)
}