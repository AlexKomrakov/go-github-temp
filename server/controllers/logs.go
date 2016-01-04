package controllers
import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ActiveState/tail"
	"os"
)

func Logs(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	t, err := tail.TailFile(params["name"] + ".log", tail.Config{Follow: false, Location: &tail.SeekInfo{0, os.SEEK_SET}})
	if err != nil {
		panic(err)
	}

	Render(res, req, "logs", t)
}
