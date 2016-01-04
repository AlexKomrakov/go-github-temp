package models

import (
	"testing"
	"fmt"
)

func TestServerStore(t *testing.T) {
	server := Server{
		User: "test",
		User_host: "user:host",
		Password: "password",
		Checked: true,
	}
	server.Store()
	server.Store()

	result := server.Find()
	fmt.Println(result)

	count, _ := server.Delete()
	fmt.Println(count)
}