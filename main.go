package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/alexkomrakov/gohub/mongo"
	"github.com/google/go-github/github"
	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"os/user"
	"io/ioutil"
	"bytes"
)

const (
	deploy_file = ".deploy.yml"
)

type ymlConfig struct {
	Host     []interface{}
	Commands []interface{}
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

func setGitStatus(client *github.Client, br *mongo.Branch, state string) {
	context := "continuous-integration/gorgon-ci"
	status := &github.RepoStatus{State: &state, Context: &context}
	_, resp, err := client.Repositories.CreateStatus(br.Owner, br.Repo, br.Sha, status)
	fmt.Print(resp)
	fmt.Print(err)
}

func readYamlConfig(file []byte) (ymlConfig, error) {
	config := ymlConfig{}
	err := yaml.Unmarshal(file, &config)
	if err != nil {
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
	if err != nil {
		return
	}
	return
}

//TODO Defer close
func getSshClient(user_host string) *ssh.Client {
	key, err := getKeyFile()
	if err != nil {
		panic(err)
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
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
	fmt.Printf("unable to connect: %s", err)
	}

	return client
}

func runCommands(client *github.Client, build *mongo.Build, config ymlConfig) {
	sshClient := getSshClient(config.Host[0].(string))
	defer sshClient.Close()

	for _, command := range config.Commands {
		switch v := command.(type) {
		case map[interface{}]interface{}:
			ma := command.(map[interface{}]interface{})
			setGitStatus(client, build.Branch, "pending")
			for commandType, action := range ma {
				actionStr := action.(string)
				if commandType == "status" {
					setGitStatus(client, build.Branch, actionStr)
				}
				if commandType == "ssh" {
					out, err := execSshCommand(sshClient, actionStr)
					fmt.Println(out.String())
					fmt.Println(err.String())
					fmt.Println(v)
				}
			}
			setGitStatus(client, build.Branch, "success")
		default:
			panic("Error on run yaml config commands")
		}
	}
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
	err2 := session.Run(command) //TODO Check error
	if err2 != nil {
		panic("Error: " + err2.Error())
	}

	return outBuf, errBuf
}

func GithubHookApi(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")

	var data github.PullRequestEvent
	json.Unmarshal([]byte(body), &data)

	owner_name := *data.Repo.Owner.Login
	repo_name  := *data.Repo.Name
	sha        := *data.PullRequest.Head.SHA

	branch     := mongo.Branch{owner_name, repo_name, sha}
	build      := &mongo.Build{&branch, &data}
	repository := branch.GetRepository()
	client     := repository.GetGithubClient()

	content, _ := getGithubFileContent(client, branch, deploy_file)
	conf, _    := readYamlConfig(content)

	runCommands(client, build, conf)
}

func GetReposApi (r render.Render) {
	repositories := mongo.GetRepositories()
	r.JSON(200, repositories)
}

func PostReposApi (res http.ResponseWriter, req *http.Request, r render.Render) {
	decoder := json.NewDecoder(req.Body)
	var repo mongo.Repository
	err := decoder.Decode(&repo)
	if err != nil {
		panic(err)
	}
	mongo.AddRepository(&repo)
	r.JSON(200, map[string]string{"status": "ok"})
}

func RepoPage (params martini.Params, r render.Render) {
	r.HTML(200, "repo", params)
}

func Index (r render.Render) {
	r.HTML(200, "index", nil)
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Layout: "base"}))

	m.Get("/", Index)
	m.Get("/repos", GetReposApi)
	m.Post("/repos", PostReposApi)
	m.Post("/hooks", GithubHookApi)
	m.Get("/repos/:user/:repo", RepoPage)
	m.Run()
}
