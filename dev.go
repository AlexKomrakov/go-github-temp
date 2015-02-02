package main

import (
	"bytes"
	"strings"
	//"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"code.google.com/p/goauth2/oauth"
)

const (
	config_file = ".config.yml"
	deploy_file = ".deploy.yml"
)

var (
	Error        *log.Logger
	github_token  string
)

func readConfig() (map[string]string) {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		panic(err)
	}

	var config map[string]string
	err2 := yaml.Unmarshal(b, &config)
	if err2 != nil {
		fmt.Println("Error on reading yaml config")
	}

	return config
}

func start() {
	config := readConfig()
	github_token = config["github_token"]

	file := loggerInit(config["log_file"])
	defer file.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/github", githubHandler)

	http.Handle("/", router)
	err := http.ListenAndServe(config["server"], nil)
	if err != nil {
		Error.Println("Error on starting server: %v", err)
	}
}

func githubHandler(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")

	var data github.PullRequestEvent
	json.Unmarshal([]byte(body), &data)

	owner_name := data.Repo.Owner.Login
	repo_name  := data.Repo.Name
	sha        := data.PullRequest.Head.SHA

	currentBranch := branch{owner_name, repo_name, sha}

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: github_token},
	}
	client := github.NewClient(transport.Client())

	content, _ := getGithubFileContent(client, currentBranch, deploy_file)
	conf, _ := readYamlConfig(content)

	runCommands(client, currentBranch, conf)
}

func getGithubFileContent(client *github.Client, br mongo.Branch, filename string) ([]byte, error) {
	repoOptions := &github.RepositoryContentGetOptions{br.Sha}
	a, _, _, err1 := client.Repositories.GetContents(br.Owner, br.Repo, filename, repoOptions)
	if err1 != nil {
		fmt.Println("Error on getting file from github: %v", err1)
		return nil, err1
	}

	fileContent, err2 := a.Decode()
	if err2 != nil {
		fmt.Println("Error on decoding file from github: %v", err2)
		return nil, err2
	}

	return fileContent, nil
}

func runCommands(client *github.Client, br mongo.Branch, config ymlConfig) {
	sshClient := getSshClient(config.Host[0].(string))
	defer sshClient.Close()

	for _, command := range config.Commands {

		switch v := command.(type) {
		case map[interface{}]interface{}:
			ma := command.(map[interface{}]interface{})
			setGitStatus(client, br, "pending")
			for commandType, action := range ma {
				actionStr := action.(string)
				if commandType == "status" {
					setGitStatus(client, br, actionStr)
				}
				if commandType == "ssh" {
					out, err := execSshCommand(sshClient, actionStr)
					fmt.Println(out.String())
					fmt.Println(err.String())
				}
			}
			setGitStatus(client, br, "success")
		default:
			Error.Printf("Error on run yaml config commands. %v", v)
		}
	}
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


func getKeyFile() (key ssh.Signer, err error) {
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}


