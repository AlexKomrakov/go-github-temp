package service
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "fmt"
)

func TestGetGithubClient(t *testing.T) {
    client := GetGithubClient("")
    assert.NotNil(t, client)

//    user, _, _ := client.Users.Get("")
//    fmt.Print(user)

    repos, _, _ := client.Repositories.List("", nil)
    fmt.Print(repos)

}


