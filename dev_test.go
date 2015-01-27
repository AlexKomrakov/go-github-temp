package main

import (
//	"code.google.com/p/goauth2/oauth"
//	"fmt"
//	"github.com/google/go-github/github"
//	//	"github.com/stretchr/testify/assert"
	"testing"
//	"encoding/json"
)

func TestConnect(t *testing.T) {
//		hook := new(github.Hook)
//		json.Unmarshal([]byte(`{"name": "web", "active": false, "events": ["pull_request"],	"config": {	"url": "komrakov-stage.smart-crowd.ru:8080/github"}}`), &hook)
//		fmt.Println(hook)
//
//		token := "389924dc1c4981bdd9ffce7bb6de96f7ce18faef" // https://gist.github.com/AlexKomrakov/a55a5867b17eed3057ac
//		transport := &oauth.Transport{
//			Token: &oauth.Token{AccessToken: token},
//		}
//		client := github.NewClient(transport.Client())
//		_, _, err := client.Repositories.CreateHook("alexkomrakov", "gohub", hook)
//		fmt.Println(err)


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
