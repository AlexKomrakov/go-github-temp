package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexkomrakov/gohub/mongo"
	"github.com/go-martini/martini"
	"github.com/google/go-github/github"
	"github.com/martini-contrib/render"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os/user"
	"strings"
)

const (
	deploy_file = ".deploy.yml"
	config_file = ".config.yml"
)

var config ServerConfig

type ymlConfig struct {
	Host         string
	Pull_request struct{ Commands []map[string]string }
	Push         struct {
		Branch   string
		Commands []map[string]string
	}
}

func getGithubFileContent(client *github.Client, br mongo.Branch, filename string) ([]byte, error) {
	repoOptions := &github.RepositoryContentGetOptions{br.Sha}
	a, _, _, err1 := client.Repositories.GetContents(br.Owner, br.Repo, filename, repoOptions)
	if err1 != nil {
		panic(err1)
	}

	fileContent, err2 := a.Decode()
	if err2 != nil {
		panic(err2)
	}

	return fileContent, nil
}

// Statuses: pending, success, error, or failure
func setGitStatus(client *github.Client, build *mongo.Build, state string) (out string, err error) {
	context := "continuous-integration/gorgon-ci"
	url := config.Adress + "/repos/" + build.Branch.Owner + "/" + build.Branch.Repo + "/" + build.Id.Hex()
	status := &github.RepoStatus{State: &state, Context: &context, TargetURL: &url}
	repoStatus, _, err := client.Repositories.CreateStatus(build.Branch.Owner, build.Branch.Repo, build.Branch.Sha, status)
	out = "Success. Current github branch status: " + *repoStatus.State
	return
}

func readYamlConfig(file []byte) (ymlConfig, error) {
	config := ymlConfig{}
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Print(err.Error())
		panic("Error on reading yaml config")
	}

	return config, nil
}

func getKeyFile() (key ssh.Signer, err error) {
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	return
}

func getSshClient(user_host string) (client *ssh.Client, err error) {
	key, err := getKeyFile()
	if err != nil {
		return
	}

	params := strings.Split(user_host, "@")
	if len(params) != 2 {
		panic("Wrong ssh user@host in config: " + user_host)
	}
	user := params[0]
	host := params[1]

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(ssh_pass),
			ssh.PublicKeys(key),
		},
	}
	client, err = ssh.Dial("tcp", host, config)
	if err != nil {
		return
	}

	return client, err
}

func runCommands(build *mongo.Build, client *github.Client, event string, config ymlConfig) {
	var commands []mongo.Command
	var out string
	var err error
	var actions []map[string]string

	if event == "push" {
		actions = config.Push.Commands
	} else if event == "pull_request" {
		actions = config.Pull_request.Commands
	}

	for _, command := range actions {
		for commandType, actionStr := range command {
			if commandType == "status" {
				out, err = setGitStatus(client, build, actionStr)
			}
			if commandType == "ssh" {
				out, err = execSshCommand(config.Host, actionStr)
			}
			if err != nil {
				commands = append(commands, mongo.Command{commandType, actionStr, out, err.Error()})
			} else {
				commands = append(commands, mongo.Command{Type: commandType, Action: actionStr, Out: out})
			}
		}
		if err != nil {
			out, err = setGitStatus(client, build, "error")
			if err != nil {
				commands = append(commands, mongo.Command{"status", "error", out, err.Error()})
			} else {
				commands = append(commands, mongo.Command{Type: "status", Action: "error", Out: out})
			}
			break
		}
	}
	build.Commands = commands
	build.Save()
}

func execSshCommand(host string, command string) (out string, err error) {
	client, err := getSshClient(host)
	if err != nil {
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	err = session.Run(command) //TODO Check error
	if err != nil {
		return
	}

	return outBuf.String(), errors.New(errBuf.String())
}

type PushEvent struct {
	HeadCommit *PushEventCommit   `json:"head_commit,omitempty"`
	Forced     *bool              `json:"forced,omitempty"`
	Created    *bool              `json:"created,omitempty"`
	Deleted    *bool              `json:"deleted,omitempty"`
	Ref        *string            `json:"ref,omitempty"`
	Before     *string            `json:"before,omitempty"`
	After      *string            `json:"after,omitempty"`
	Compare    *string            `json:"compare,omitempty"`
	Size       *int               `json:"size,omitempty"`
	Commits    []PushEventCommit  `json:"commits,omitempty"`
	Repo       *github.Repository `json:"repository,omitempty"`
}

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	ID       *string              `json:"id,omitempty"`
	Message  *string              `json:"message,omitempty"`
	Author   *github.CommitAuthor `json:"author,omitempty"`
	URL      *string              `json:"url,omitempty"`
	Distinct *bool                `json:"distinct,omitempty"`
	Added    []string             `json:"added,omitempty"`
	Removed  []string             `json:"removed,omitempty"`
	Modified []string             `json:"modified,omitempty"`
}

func GithubHookApi(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")
	event := req.Header["X-Github-Event"][0]
	switch event {
	case "pull_request":
		var pullRequestEvent github.PullRequestEvent
		json.Unmarshal([]byte(body), &pullRequestEvent)
		actions := map[string]bool{"opened": true, "reopened": true, "synchronize": true}
		if actions[*pullRequestEvent.Action] {
			branch := mongo.Branch{*pullRequestEvent.Repo.Owner.Login, *pullRequestEvent.Repo.Name, *pullRequestEvent.PullRequest.Head.SHA}
			build := &mongo.Build{branch, pullRequestEvent, nil, bson.NewObjectId()}
			client := build.Branch.GetRepository().GetGithubClient()
			content, _ := getGithubFileContent(client, build.Branch, deploy_file)
			content = []byte(strings.Replace(string(content), "{{sha}}", build.Branch.Sha, -1))
			config, _ := readYamlConfig(content)

			runCommands(build, client, event, config)
		} else {
			fmt.Print("Skipping pull request event type: " + *pullRequestEvent.Action)
		}
	case "push":
		var pushEvent PushEvent
		json.Unmarshal([]byte(body), &pushEvent)
		branch := mongo.Branch{*pushEvent.Repo.Owner.Name, *pushEvent.Repo.Name, *pushEvent.After}
		build := &mongo.Build{branch, pushEvent, nil, bson.NewObjectId()}
		client := build.Branch.GetRepository().GetGithubClient()
		content, _ := getGithubFileContent(client, build.Branch, deploy_file)
		content = []byte(strings.Replace(string(content), "{{sha}}", build.Branch.Sha, -1))
		config, _ := readYamlConfig(content)

		if config.Push.Branch == *pushEvent.Ref {
			runCommands(build, client, event, config)
		}
	default:
		fmt.Println("Not supported event: " + req.Header["X-Github-Event"][0])
		fmt.Println(body)
	} 
}

func GetReposApi(r render.Render) {
	repositories := mongo.GetRepositories()
	r.JSON(200, repositories)
}

func PostReposApi(res http.ResponseWriter, req *http.Request, r render.Render) {
	decoder := json.NewDecoder(req.Body)
	var repo mongo.Repository
	err := decoder.Decode(&repo)
	if err != nil {
		panic(err)
	}
	mongo.AddRepository(&repo)
	r.JSON(200, map[string]string{"status": "ok"})
}

func RepoPage(params martini.Params, r render.Render) {
	data := make(map[string]interface{})
	data["params"] = params
	data["builds"] = mongo.GetBuilds(params["user"], params["repo"])
	r.HTML(200, "repo", data)
}

func BuildPage(params martini.Params, r render.Render) {
	data := make(map[string]interface{})
	data["params"] = params
	data["builds"] = mongo.GetBuilds(params["user"], params["repo"])
	data["build"] = mongo.GetBuild(params["build"])
	r.HTML(200, "repo", data)
}

func Index(r render.Render) {
	r.HTML(200, "index", nil)
}

type ServerConfig struct {
	Server string
	Adress string
}

func readConfig() ServerConfig {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		panic(err)
	}

	var config ServerConfig
	err2 := yaml.Unmarshal(b, &config)
	if err2 != nil {
		fmt.Println("Error on reading yaml config")
	}

	return config
}

func main() {
	config = readConfig()
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Layout: "base"}))
	m.Get("/", Index)
	m.Get("/repos", GetReposApi)
	m.Post("/repos", PostReposApi)
	m.Post("/hooks", GithubHookApi)
	m.Get("/repos/:user/:repo", RepoPage)
	m.Get("/repos/:user/:repo/:build", BuildPage)
	m.RunOnAddr(config.Server)
}
