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
	Host     string
	Commands []map[string]string
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

func runCommands(build *mongo.Build) {
	var commands []mongo.Command
	var out string
	var err error

	client := build.Branch.GetRepository().GetGithubClient()

	content, _ := getGithubFileContent(client, build.Branch, deploy_file)
	content = []byte(strings.Replace(string(content), "{{sha}}", build.Branch.Sha, -1))
	config, _ := readYamlConfig(content)

	for _, command := range config.Commands {
		out, err = setGitStatus(client, build, "pending")
		if err != nil {
			commands = append(commands, mongo.Command{"status", "pending", out, err.Error()})
		} else {
			commands = append(commands, mongo.Command{Type: "status", Action: "pending", Out: out})
		}
		for commandType, action := range command {
			actionStr := action
			if commandType == "status" {
				out, err = setGitStatus(client, build, actionStr)
			}
			if commandType == "ssh" {
				out, err = execSshCommand(config.Host, actionStr)
			}
			if err != nil {
				commands = append(commands, mongo.Command{commandType, actionStr, out, err.Error()})
				break
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
		} else {
			out, err = setGitStatus(client, build, "success")
			if err != nil {
				commands = append(commands, mongo.Command{"status", "success", out, err.Error()})
			} else {
				commands = append(commands, mongo.Command{Type: "status", Action: "success", Out: out})
			}
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

func GithubHookApi(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")

	switch req.Header["X-Github-Event"][0] {
	case "pull_request":
		var pullRequestEvent github.PullRequestEvent
		json.Unmarshal([]byte(body), &pullRequestEvent)
		actions := map[string]bool{"opened": true, "reopened": true, "synchronize": true}
		if actions[*pullRequestEvent.Action] {
			branch := mongo.Branch{*pullRequestEvent.Repo.Owner.Login, *pullRequestEvent.Repo.Name, *pullRequestEvent.PullRequest.Head.SHA}
			build := &mongo.Build{branch, pullRequestEvent, nil, bson.NewObjectId()}
			runCommands(build)
		} else {
			fmt.Print("Skipping pull request event type: " + *pullRequestEvent.Action)
		}
	case "push":
		var pushEvent github.PushEvent
		json.Unmarshal([]byte(body), &pushEvent)
		fmt.Println(*pushEvent.Ref)
		fmt.Println("Recieved push event")
		fmt.Println(body)
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
