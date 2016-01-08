package controllers
import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/goincremental/negroni-sessions"
	"github.com/alexkomrakov/gohub/models"
	"github.com/alexkomrakov/gohub/service"
	"strings"
	"github.com/google/go-github/github"
)

func ShowRepo(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	session := sessions.GetSession(req)
	user := session.Get("user").(string)
	token, _ := models.GetToken(user)
	client := service.GetGithubClient(token)
	repo, _, _ := client.Repositories.Get(params["user"], params["repo"])
	builds, _ := models.Build{Login: params["user"], Name: params["repo"]}.GetBuilds()
	hooks, _, _ := client.Repositories.ListHooks(params["user"], params["repo"], &github.ListOptions{})
	branches, _, _ := client.Git.ListRefs(params["user"], params["repo"], &github.ReferenceListOptions{})
	var filtered_branches []github.Reference
	for _, branch := range branches {
		if strings.HasPrefix(*branch.Ref, "refs/heads/") {
			filtered_branches = append(filtered_branches, branch)
		}
	}

	Render(res, req, "repo", map[string]interface{}{"Repo": repo, "Hooks": hooks, "Branches": filtered_branches, "Builds": builds})
}