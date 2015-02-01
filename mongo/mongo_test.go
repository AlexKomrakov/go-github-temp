package mongo

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestMongo(t *testing.T) {
	repo := Repository{"alexkomrakov", "gohub", "testKey"}
	AddRepository(&repo)
	repos := GetRepositories()
	assert.Len(t, repos, 1)
	assert.Equal(t, repos[0], repo)
	for _, v := range repos {
		RemoveRepository(&v)
	}
	repos = GetRepositories()
	assert.Len(t, repos, 0)
}

func TestGetRepository(t *testing.T) {
	branch := Branch{"AlexKomrakov", "gohub", "asdsad"}
	repo := branch.GetRepository()
	fmt.Print(repo)
}

