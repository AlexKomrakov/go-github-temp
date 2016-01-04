package controllers
import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
)

func Index(res http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	user := session.Get("user")
	if user == nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	http.Redirect(res, req, "/repos/" + user.(string), http.StatusFound)
	return
}