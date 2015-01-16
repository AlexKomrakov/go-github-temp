package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/google/go-github/github"
	"encoding/json"
	"fmt"
	"code.google.com/p/goauth2/oauth"
)

func TestTest(t *testing.T){
	assert.Equal(t, test() , "hello world")
}

func TestConnect(t *testing.T){
	hook := new(github.Hook)
	json.Unmarshal([]byte(`{"name": "web", "active": false, "events": ["pull_request"],	"config": {	"url": "https://godoc.org/golang.org/x/oauth2"}}`), &hook)
//	fmt.Print(hook)

	token := "" // https://gist.github.com/AlexKomrakov/a55a5867b17eed3057ac
	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	client := github.NewClient(transport.Client())
	fmt.Print(client.Repositories.CreateHook("alexkomrakov", "gohub", hook))

//	options :=& github.ListOptions{1,5}
//	res, _, _ := client.Repositories.ListHooks("alexkomrakov", "gohub", options)
//	fmt.Println(res)
}


