package service

import (
	"errors"
	"bytes"
	"github.com/google/go-github/github"
	"github.com/alexkomrakov/gohub/mongo"
)

type CommandResponse struct {
	Type    string
	Command string
	Error   error
	Output  string
}

func RunCommands(config DeployScenario, client *github.Client, user, repo, sha string) (result []CommandResponse) {
	server := mongo.Server{User: user, User_host: config.Host}.Find()
	for _, command := range config.Commands {
		for commandType, actionStr := range command {
			if commandType == "status" {
				out, err := SetGitStatus(client, user, repo, sha, actionStr)
				result = append(result, CommandResponse{Type: commandType, Command: actionStr, Output: out, Error: err})
			}
			if commandType == "ssh" {
				out, err := ExecSshCommand(server, actionStr)
				result = append(result, CommandResponse{Type: commandType, Command: actionStr, Output: out, Error: err})
			}
		}
	}
	return
}

func ExecSshCommand(server mongo.Server, command string) (out string, err error) {
	client, err := server.Client()
	if err != nil {
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	err = session.Run(command)
	if err != nil {
		return
	}

	return outBuf.String(), errors.New(errBuf.String())
}

func SetGitStatus(client *github.Client, user, repo, sha, state string) (out string, err error) {
	_, _, err = client.Repositories.CreateStatus(user, repo, sha, &github.RepoStatus{State: &state})
	out = "Success. Current github branch status: " + state
	return
}
