package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTest(t *testing.T) {
	assert.Equal(t, test(), "hello world")
}

func TestConnect(t *testing.T) {
	hook := new(github.Hook)
	json.Unmarshal([]byte(`{"name": "web", "active": false, "events": ["pull_request"],	"config": {	"url": "http://requestb.in/1c8ldr11"}}`), &hook)
	//	fmt.Print(hook)

	token := "389924dc1c4981bdd9ffce7bb6de96f7ce18faef" // https://gist.github.com/AlexKomrakov/a55a5867b17eed3057ac
	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	client := github.NewClient(transport.Client())
	fmt.Print(client.Repositories.CreateHook("alexkomrakov", "gohub", hook))

	//	options :=& github.ListOptions{1,5}
	//	res, _, _ := client.Repositories.ListHooks("alexkomrakov", "gohub", options)
	//	fmt.Println(res)
}
