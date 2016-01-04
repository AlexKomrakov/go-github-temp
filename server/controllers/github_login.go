package controllers
import (
	"net/http"
	"golang.org/x/oauth2"
)

func GithubLogin(w http.ResponseWriter, req *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}
