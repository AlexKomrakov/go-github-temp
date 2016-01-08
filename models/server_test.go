package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestServerStore(t *testing.T) {
	server := Server{
		User: "test",
		User_host: "user:host",
		Password: "password",
		Checked: true,
	}
	server.Store()

	new_server := Server{User: "test", User_host: "user:host"}
	result, err := new_server.FindOne()
	assert.True(t, result)
	assert.Nil(t, err)
	assert.NotEmpty(t, new_server.Password)

	n, err := server.Delete()
	assert.NotEmpty(t, n)
	assert.Nil(t, err)
}