package service
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "fmt"
)

func TestGetGithubClient(t *testing.T) {
    client := GetGithubClient("7bd2ae71bf78ab0489052ef560ab53771f372980")
    assert.NotNil(t, client)

    user, _, _ := client.Users.Get("")
    fmt.Print(user)
}


