package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
)

func EnableRepo(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	params := mux.Vars(req)
	repo := models.Repository{Login: params["user"], Name: params["repo"]}
	// TODO Check repo permissions
	success, err  := repo.FindOne()
	if success != true {
		panic(err)
	}
	repo.Enabled = !repo.Enabled
	_, err = repo.Update()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, req, "/repos/" + user, http.StatusFound)
}
