package controllers
import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
	"github.com/alexkomrakov/gohub/service"
	"strconv"
)

func DeleteHook(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	token, _ := models.GetToken(user)
	client := service.GetGithubClient(token)

	id, _ := strconv.Atoi(params["id"])
	client.Repositories.DeleteHook(params["user"], params["repo"], id)

	http.Redirect(w, req, "/repos/" + params["user"] + "/" + params["repo"], http.StatusTemporaryRedirect)
}
