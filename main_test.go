package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
	//	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
)

func TestConnect(t *testing.T) {
		hook := new(github.Hook)
		json.Unmarshal([]byte(`{"name": "web", "active": false, "events": ["pull_request"],	"config": {	"url": "http://requestb.in/13ymqwv1"}}`), &hook)
		fmt.Print(hook)

		token := "" // https://gist.github.com/AlexKomrakov/a55a5867b17eed3057ac
		transport := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token},
		}
		client := github.NewClient(transport.Client())
		client.Repositories.CreateHook("alexkomrakov", "gohub", hook)

//		status := new(github.RepoStatus)
//		json.Unmarshal([]byte(`{"state": "success"}`), &status)
//		_, b, _ := client.Repositories.CreateStatus("alexkomrakov", "gohub", "", status)
//		fmt.Println(b)
//
//		listOptions := new(github.ListOptions)
//		json.Unmarshal([]byte(`{"page": 0, "perPage": 5}`), &listOptions)
//		z, c, _ := client.Repositories.ListStatuses("alexkomrakov", "gohub", "", listOptions)
//		fmt.Println(z)
//		fmt.Println(c)
}

func TestReadFile(t *testing.T) {
	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: github_token},
	}
	client := github.NewClient(transport.Client())

	content, _ := getGithubFileContent(client, github_user, github_repo, ".deploy.yml", github_ref)
	fmt.Print(string(content))

	conf, _ := readYamlConfig(content)
	fmt.Println(conf)
}

func TestExecSshCommand(t *testing.T) {
//	client := getSshClient()
//	defer client.Close()
//
//	out, _ := execSshCommand(client, "git")
//	fmt.Print(out.String())
}

func TestRunCommands(t *testing.T) {
	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: github_token},
	}
	client := github.NewClient(transport.Client())

	content, _ := getGithubFileContent(client, github_user, github_repo, ".deploy.yml", github_ref)
	fmt.Print(string(content))

	conf, _ := readYamlConfig(content)

	runCommands(client, conf)
}
