package controllers

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
)

func AddServer(res http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	req.ParseForm()

	server := models.Server{User: user, User_host: req.PostFormValue("user_host"), Password: req.PostFormValue("password") }
	server.Checked = server.Check()
	server.Store()

	http.Redirect(res, req, "/servers/"+user, http.StatusFound)
}
