package server

import (
	"github.com/gorilla/mux"
)

func Router() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")

	router.HandleFunc("/login/github", GithubLogin).Methods("GET")
	router.HandleFunc("/login/github/callback", GithubLogin).Methods("GET")

	router.HandleFunc("/login", Login).Methods("GET", "POST")
	router.HandleFunc("/logout", Logout).Methods("GET")

	router.HandleFunc("/logs/{name}", Logs).Methods("GET")

	router.HandleFunc("/repos/{user}", UserRepos).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}", ShowRepo).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/add", AddRepo).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/delete", DeleteRepo).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/commit/{sha}", ShowCommit).Methods("GET")
	router.HandleFunc("/repos/{user}/{repo}/commit/{sha}/run/{scenario}", RunScenario).Methods("GET")


	router.HandleFunc("/servers/{user}", UserServers).Methods("GET")
	router.HandleFunc("/servers/{user}", AddServer).Methods("POST")
	router.HandleFunc("/servers/{user}/delete", DeleteServer).Methods("POST")

	router.HandleFunc("/hooks", GithubHookApi).Methods("POST")
	router.HandleFunc("/hooks", GithubHookApi).Methods("POST")


	//	router.HandleFunc("/repos", GetReposApi).Methods("GET")
//	router.HandleFunc("/repos/{user}/{repo}/hook", SetHook).Methods("GET")
//	router.HandleFunc("/repos/{user}/{repo}/{build}", BuildPage).Methods("GET")
//	router.HandleFunc("/repos", PostReposApi).Methods("POST")

	return
}
