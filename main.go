package main

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// https://gist.github.com/AlexKomrakov/412a549b693c5f0a03d6
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
	z, c, _ := client.Repositories.ListStatuses(github_user, github_repo, github_ref, listOptions)

	io.WriteString(w, fmt.Sprintf("%v", z))
	io.WriteString(w, fmt.Sprintf("%v", c))

	setGitStatus(client, "success")
}

func setGitStatus(client *github.Client, state string) {
	context := "continuous-integration/gorgon-ci"
	status := &github.RepoStatus{State: &state, Context: &context}
	client.Repositories.CreateStatus(github_user, github_repo, github_ref, status)
}

func getGithubFileContent(client *github.Client, owner, repo, filename, branch string) ([]byte, error) {
	repoOptions := &github.RepositoryContentGetOptions{branch}
	a, _, _, err1 := client.Repositories.GetContents(owner, repo, filename, repoOptions)
	if err1 != nil {
		Error.Println("Error on getting file from github: %v", err1)
		return nil, err1
	}

	fileContent, err2 := a.Decode()
	if err2 != nil {
		Error.Println("Error on decoding file from github: %v", err2)
		return nil, err2
	}

	return fileContent, nil
}

type ymlConfig struct {
	Host     []interface{}
	Commands []interface{}
}

func readYamlConfig(file []byte) (ymlConfig, error) {
	config := ymlConfig{}
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		Error.Println("Error on reading yaml config")
		return config, err
	}

	return config, nil
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

func execSshCommand(client *ssh.Client, command string) (bytes.Buffer, bytes.Buffer) {
	session, err1 := client.NewSession()
	if err1 != nil {
		panic("Failed to create session: " + err1.Error())
	}
	defer session.Close()

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	session.Run(command) //TODO Check error
	//	if err2 != nil {
	//		panic("Failed to run: " + err2.Error())
	//	}

	return outBuf, errBuf
}
