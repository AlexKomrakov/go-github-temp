package controllers
import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
	"github.com/gorilla/schema"
)

func DeleteServer(res http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	req.ParseForm()

	server := new(models.Server)
	schema.NewDecoder().Decode(server, req.PostForm)
	server.Delete()

	http.Redirect(res, req, "/servers/"+user, http.StatusFound)
}