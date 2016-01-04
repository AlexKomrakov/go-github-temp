package controllers
import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
)

func Logout(res http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	session.Delete("user")

	http.Redirect(res, req, "/", http.StatusFound)
}
