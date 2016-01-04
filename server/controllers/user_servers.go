package controllers
import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
)

func UserServers(res http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	servers := models.GetServers(user)

	Render(res, req, "servers", map[string]interface{}{"Servers": servers})
}
