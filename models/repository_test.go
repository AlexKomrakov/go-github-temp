package models
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/google/go-github/github"
)

func TestRepository(t *testing.T) {
	repo := Repository{Login: "owner_login", Name: "repo_name"}
	number, err := repo.Store()
	assert.Equal(t, 1, number)
	assert.Nil(t, err)
	assert.NotEmpty(t, repo.Id)
	assert.False(t, repo.Enabled)

	repo.Enabled = true
	number, err = repo.Update()
	assert.Equal(t, 1, number)
	assert.Nil(t, err)
	assert.True(t, repo.Enabled)

	second_repo := Repository{Login: "owner_login", Name: "repo_name"}
	second_repo.FindOne()
	assert.Equal(t, second_repo.Id, repo.Id)
	assert.True(t, second_repo.Enabled)

	third_repo := Repository{Login: "owner_login", Name: "not_existing_repo_name"}
	success, err := third_repo.FindOne()
	assert.False(t, success)
	assert.Nil(t, err)

	success, err = third_repo.FindOrCreate()
	assert.True(t, success)
	assert.Nil(t, err)
	assert.NotEmpty(t, third_repo.Id)

	defer Repository{Login: "owner_login"}.Delete()
}

func TestGetGithubRepositoriesIntersection(t *testing.T) {
	name  := "name"
	login := "AlexKomrakov"

	repo := Repository{Login: login, Name: name}
	repo.Store()

	user  := github.User{Login: &login}
	github_repos := []github.Repository{{Name: &name, Owner: &user}}
	repos, _ := GetGithubRepositoriesIntersection(github_repos)
	assert.NotEmpty(t, repos)

	defer Repository{Login: login}.Delete()
}