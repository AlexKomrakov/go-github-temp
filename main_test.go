package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
	//	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"testing"
)

func TestConnect(t *testing.T) {
	//	hook := new(github.Hook)
	//	json.Unmarshal([]byte(`{"name": "web", "active": false, "events": ["pull_request"],	"config": {	"url": "http://requestb.in/1c8ldr11"}}`), &new(github.Hook))
	//	//	fmt.Print(hook)
	//
	//	token := "" // https://gist.github.com/AlexKomrakov/a55a5867b17eed3057ac
	//	transport := &oauth.Transport{
	//		Token: &oauth.Token{AccessToken: token},
	//	}
	//	client := github.NewClient(transport.Client())
	//	fmt.Print(client.Repositories.CreateHook("alexkomrakov", "gohub", hook))
	//
	//	status := new(github.RepoStatus)
	//	json.Unmarshal([]byte(`{"state": "success"}`), &status)
	//	_, b, _ := client.Repositories.CreateStatus("alexkomrakov", "gohub", "", status)
	//	fmt.Println(b)
	//
	//	listOptions := new(github.ListOptions)
	//	json.Unmarshal([]byte(`{"page": 0, "perPage": 5}`), &listOptions)
	//	z, c, _ := client.Repositories.ListStatuses("alexkomrakov", "gohub", "", listOptions)
	//	fmt.Println(z)
	//	fmt.Println(c)
}

func TestReadFile(t *testing.T) {
	token := github_token
	owner := github_user
	repo := github_repo
	branch := github_ref
	filename := deploy_config

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	client := github.NewClient(transport.Client())

	content, _ := getGithubFileContent(client, owner, repo, filename, branch)
	fmt.Print(string(content))

	conf, _ := readYamlConfig(content)
	fmt.Println(conf)
}

func TestExecSshCommand(t *testing.T) {
	config := &ssh.ClientConfig{
		User: ssh_user,
		Auth: []ssh.AuthMethod{
			ssh.Password(ssh_pass),
		},
	}
	client, err := ssh.Dial("tcp", "komrakov-stage.smart-crowd.ru:22", config)
	if err != nil {
		fmt.Printf("unable to connect: %s", err)
	}
	defer client.Close()

	out, _ := execSshCommand(client, "git")
	fmt.Print(out.String())
}
