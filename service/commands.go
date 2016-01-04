package service

import (
	"errors"
	"bytes"
	"strings"
	"github.com/google/go-github/github"
	"github.com/alexkomrakov/gohub/mongo"
	"github.com/alexkomrakov/gohub/models"
)

func ProcessHook(event, body string) {
	params := make(map[string]string)

	switch event {
	case "pull_request":
		pullRequestEvent, _ := ParsePullRequestEvent(body)
		params["user"]   = *pullRequestEvent.Repo.Owner.Login
		params["repo"]   = *pullRequestEvent.Repo.Name
		params["sha"]    = *pullRequestEvent.PullRequest.Head.SHA
		params["branch"] = *pullRequestEvent.PullRequest.Head.Ref
	case "push":
		pushEvent, _ := ParsePushEvent(body)
		params["user"]   = *pushEvent.Repo.Owner.Name
		params["repo"]   = *pushEvent.Repo.Name
		params["sha"]    = *pushEvent.After
		params["branch"] = *pushEvent.Ref
	default:
		panic("Not supported event: " + event)
	}

	token, _    := models.GetToken(params["user"])
	client      := GetGithubClient(token)
	file, _     := GetFileContent(client, params["user"], params["repo"], params["sha"], GetServerConfig().DeployFile)
	string_file := ReplaceVariables(params, string(file))
	deploy, _   := GetYamlConfig([]byte(string_file))

	if deploy[event].Branch == "" || deploy[event].Branch == params["branch"] {
		RunCommands(deploy, client, event, mongo.Build{Login: params["user"], Name: params["repo"], SHA: params["sha"]})
	}
}

func ReplaceVariables(params map[string]string, text string) string {
	for variable, value := range params {
		r := strings.NewReplacer("{{" + variable + "}}", value)
		text = r.Replace(text)
	}

	return text
}

func RunCommands(deploy map[string]mongo.DeployScenario, client *github.Client, event string, current_build mongo.Build) (build mongo.Build) {
	current_build.DeployFile = deploy
	current_build.Event = event
	current_build.Store()

	config := deploy[event]
	server := models.Server{User: current_build.Login, User_host: config.Host}.Find()
    has_error := false
    error := ""
    for _, command := range config.Commands {
		for commandType, actionStr := range command {
            if has_error == true && error != "" {
                continue
            }
            error = ""
			if commandType == "status" {
				out, err := SetGitStatus(client, current_build.Login, current_build.Name, current_build.SHA, actionStr)
                if err != nil {
                    error = err.Error()
                    has_error = true
                }
				current_build.AddCommand(mongo.CommandResponse{Type: commandType, Command: actionStr, Success: out, Error: error})
			}
			if commandType == "ssh" {
				out, err := ExecSshCommand(server, actionStr)
                if err != nil {
                    error = err.Error()
                    has_error = true
                }
				current_build.AddCommand(mongo.CommandResponse{Type: commandType, Command: actionStr, Success: out, Error: error})
			}
		}
	}
    // TODO Refactor this shit
    if has_error == true {
        for _, command := range config.Error {
            for commandType, actionStr := range command {
                error = ""
                if commandType == "status" {
                    out, err := SetGitStatus(client, current_build.Login, current_build.Name, current_build.SHA, actionStr)
                    if err != nil {
                        error = err.Error()
                    }
					current_build.AddCommand(mongo.CommandResponse{Type: commandType, Command: actionStr, Success: out, Error: error})
                }
                if commandType == "ssh" {
                    out, err := ExecSshCommand(server, actionStr)
                    if err != nil {
                        error = err.Error()
                    }
					current_build.AddCommand(mongo.CommandResponse{Type: commandType, Command: actionStr, Success: out, Error: error})
                }
            }
        }
    }

	return
}

func ExecSshCommand(server models.Server, command string) (out string, err error) {
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
