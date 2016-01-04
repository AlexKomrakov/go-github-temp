package mongo

import (
	"testing"
	"fmt"
	"golang.org/x/crypto/ssh"
	"github.com/mitchellh/go-linereader"
)

func TestSSH(t *testing.T) {
	fmt.Println("Hello Starting SSH connect")
	user     := "root"
	password := "Zxcfrt6"
	host     := "188.166.34.149:22"

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", host, config)
	fmt.Println(err)

	session, err := client.NewSession()
	fmt.Println(err)
	defer session.Close()

	reader, _ := session.StdoutPipe()
	go func() {
		err := session.Run("traceroute 8.8.8.8")
		fmt.Println(err)
	}()

	lr := linereader.New(reader)
	for line := range lr.Ch {
		fmt.Println(line)
	}

	fmt.Println(22)
}