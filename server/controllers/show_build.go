package controllers
import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/alexkomrakov/gohub/models"
)

func ShowBuild(res http.ResponseWriter, req *http.Request) {
	params   := mux.Vars(req)
	build, _ := models.FindBuildById(params["id"])
	command_responses, _ := build.CommandResponses()

	Render(res, req, "build", map[string]interface{}{"Params": params, "Build": build, "CommandResponses": command_responses})
}
