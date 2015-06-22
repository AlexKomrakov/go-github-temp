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

func GetFileContent(client *github.Client, owner, repo, sha, filename string) (fileContent []byte, err error) {
    repoOptions := &github.RepositoryContentGetOptions{sha}
    a, _, _, err := client.Repositories.GetContents(owner, repo, filename, repoOptions)
    if err != nil {
        return
    }

    fileContent, err = a.Decode()
    return
}
