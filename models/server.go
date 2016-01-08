package models

import (
	"golang.org/x/crypto/ssh"
	"strings"
)

// TODO Добавить проверку на униальность хостов
type Server struct {
	Id        int64  `json:"id"`
	User      string `json:"user"`
	User_host string `json:"user_host"`
	Password  string `json:"password"`
	Checked   bool   `json:"checked"`
}

func GetServers(user string) (servers []Server) {
	err := Orm.Find(&servers, &Server{User: user})
	if err != nil {
		panic(err)
	}

	return
}

func (r Server) Store() {
	_, err := Orm.Insert(&r)
	if err != nil {
		panic(err)
	}
}

func (r Server) Delete() (int64, error) {
	return Orm.Delete(&r)
}

func (r Server) Find() Server {
	_, err := Orm.Get(&r)
	if err != nil {
		panic(err)
	}

	return r
}

func (s Server) Check() bool {
	_, err := GetSshClient(s.User_host, s.Password)

	return err == nil
}

func (s Server) Client() (client *ssh.Client, err error) {
	return GetSshClient(s.User_host, s.Password)
}

func GetSshClient(user_host, password string) (client *ssh.Client, err error) {
	params := strings.Split(user_host, "@")
	if len(params) != 2 {
		panic("Wrong ssh user@host: " + user_host)
	}
	user := params[0]
	host := params[1]

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err = ssh.Dial("tcp", host, config)

	return client, err
}