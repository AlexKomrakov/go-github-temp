package controllers
import (
	"net/http"
	"github.com/alexkomrakov/gohub/service"
)

func Login(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	token := req.FormValue("token")

	client := service.GetGithubClient(token)
	user, _, _ := client.Users.Get("")
	process_login(res, req, user, token)

	Render(res, req, "login", nil)
}
