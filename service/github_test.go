package service
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/alexkomrakov/gohub/services"
)

func TestGetGithubClient(t *testing.T) {
    client := services.GetGithubClient("")
    assert.NotNil(t, client)
}


