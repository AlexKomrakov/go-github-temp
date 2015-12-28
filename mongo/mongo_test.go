package mongo

import (
	"testing"
	"fmt"
	"golang.org/x/crypto/ssh"
	"github.com/mitchellh/go-linereader"
)

//func TestMongo(t *testing.T) {
//	repo := Repository{"alexkomrakov", "gohub", "testKey"}
//	AddRepository(&repo)
//	repos := GetRepositories()
//	assert.Len(t, repos, 1)
//	assert.Equal(t, repos[0], repo)
//	for _, v := range repos {
//		RemoveRepository(&v)
//	}
//	repos = GetRepositories()
//	assert.Len(t, repos, 0)
//}
//
//func TestGetRepository(t *testing.T) {
//	branch := Branch{"AlexKomrakov", "gohub", "asdsad"}
//	repo := branch.GetRepository()
//	fmt.Print(repo)
//}

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