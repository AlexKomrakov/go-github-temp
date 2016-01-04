package controllers
import (
	"net/http"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

func GithubLoginCallback(w http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != oauthStateString {
		l.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	code := req.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		l.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		l.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	process_login(w, req, user, token.AccessToken)
	http.Redirect(w, req, "/login", http.StatusFound)
}
