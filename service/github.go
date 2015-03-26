package service
import (
    "github.com/google/go-github/github"
    "code.google.com/p/goauth2/oauth"
)


func GetGithubClient(token string) *github.Client {
    transport := &oauth.Transport{
        Token: &oauth.Token{AccessToken: token},
    }
    return github.NewClient(transport.Client())
}