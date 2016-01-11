package models
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	repo := Repository{Login: "owner_login", Name: "repo_name"}
	number, err := repo.Store()
	assert.Equal(t, 1, number)
	assert.Nil(t, err)
	assert.NotEmpty(t, repo.Id)

	second_repo := Repository{Login: "owner_login", Name: "repo_name"}
	second_repo.FindOne()
	assert.Equal(t, second_repo.Id, repo.Id)

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