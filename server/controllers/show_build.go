package controllers
import (
	"net/http"
	"github.com/alexkomrakov/gohub/mongo"
	"github.com/gorilla/mux"
)

func ShowBuild(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	build, _ := mongo.FindBuildById(params["id"])

	Render(res, req, "build", map[string]interface{}{"Params": params, "Build": build})
}
