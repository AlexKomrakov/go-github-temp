package controllers
import (
	"net/http"
	"github.com/alexkomrakov/gohub/service"
)

func GithubHookApi(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")
	event := req.Header["X-Github-Event"][0]
	l.Println(event)
	l.Println(body)
	service.ProcessHook(event, body)
}
