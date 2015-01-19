package main

import (
	//	"bytes"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"testing"
)

func TestExecuteSSh(t *testing.T) {
	config := &ssh.ClientConfig{
		User: "komrakov",
		Auth: []ssh.AuthMethod{
			ssh.Password("31IXdDDu"),
		},
	}
	// Dial your ssh server.
	client, err := ssh.Dial("tcp", "komrakov-stage.smart-crowd.ru:22", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer client.Close()

	session, err1 := client.NewSession()
	if err1 != nil {
		panic("Failed to create session: " + err1.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	var z bytes.Buffer
	session.Stdout = &b
	session.Stderr = &z
	err2 := session.Run("git") //TODO Check session.StdinPipe
	fmt.Println(err2)
	//	if err2 != nil {
	//		panic("Failed to run: " + err2.Error())
	//	}
	fmt.Println(b.String())
	fmt.Println(z.String())
	//	fmt.Println(b.String())
}
