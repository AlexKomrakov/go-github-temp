package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
// https://gist.github.com/AlexKomrakov/3c3a7bee69da1fb2a328
)

var (
	Error *log.Logger
)

func main() {
	file := loggerInit(log_file)
	defer file.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/github", githubHandler)

	http.Handle("/", router)
	err := http.ListenAndServe(server, nil)
	if err != nil {
		Error.Println("Error on starting server: %v", err)
	}
}

func githubHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello github")

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: github_token},
	}
	client := github.NewClient(transport.Client())

	listOptions := new(github.ListOptions)
	json.Unmarshal([]byte(`{"page": 0, "perPage": 5}`), &listOptions)
	z, c, _ := client.Repositories.ListStatuses(github_user, github_repo, github_token, listOptions)

	io.WriteString(w, fmt.Sprintf("%v", z))
	io.WriteString(w, fmt.Sprintf("%v", c))

	setGitStatus(client)
}

func setGitStatus(client *github.Client) {
	status := new(github.RepoStatus)
	json.Unmarshal([]byte(`{"state": "success", "context": "continuous-integration/travis-ci"}`), &status)
	client.Repositories.CreateStatus(github_user, github_repo, github_token, status)
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello world")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Error.Println("error reading body")
	}
	Error.Println(string(body))
}

func loggerInit(filename string) (file *os.File) {
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		fmt.Print("Error on parsing abs path for logger: %v", err1)
	}

	filename = filepath.Join(dir, filename)

	f, err2 := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err2 != nil {
		fmt.Print("Error opening file: %v", err2)
	}
	err3 := f.Truncate(0)
	if err3 != nil {
		fmt.Print("Error on clearing log file: %v", err3)
	}

	Error = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	Error.Println("Logger start")

	return f
}
