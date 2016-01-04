package server

import (
	"github.com/gorilla/mux"
	"github.com/alexkomrakov/gohub/server/controllers"
)

func Router() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/", controllers.Index).Methods("GET")

	router.HandleFunc("/login/github", controllers.GithubLogin).Methods("GET")
	router.HandleFunc("/login/github/callback", controllers.GithubLoginCallback).Methods("GET")

	router.HandleFunc("/login", controllers.Login).Methods("GET", "POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("GET")

	router.HandleFunc("/logs/{name}", controllers.Logs).Methods("GET")

	router.HandleFunc("/repos/{user}", controllers.UserRepos).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}", controllers.ShowRepo).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/hook", controllers.SetHook).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/hook/{id}/delete", controllers.DeleteHook).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/build/{id}", controllers.ShowBuild).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/commit/{sha}", controllers.ShowCommit).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/commit/{sha}/run/{scenario}", controllers.RunScenario).Methods("GET")

	router.HandleFunc("/servers/{user}", controllers.UserServers).Methods("GET")
	router.HandleFunc("/servers/{user}", controllers.AddServer).Methods("POST")
	router.HandleFunc("/servers/{user}/delete", controllers.DeleteServer).Methods("POST")

	router.HandleFunc("/hooks", controllers.GithubHookApi).Methods("POST")

	return
}
