package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	user  := "test_user"
	token := Token{user, "test_token"}

	n, err := token.Store()
	assert.Nil(t, err)
	assert.NotEmpty(t, n)

	new_token := Token{User: user}
	result, err := new_token.FindOne()
	assert.True(t, result)
	assert.Nil(t, err)
	assert.NotEmpty(t, new_token.Token)

	n, err = token.Delete()
	assert.Nil(t, err)
	assert.NotEmpty(t, n)
}