package server

//import (
//	"github.com/alexkomrakov/gohub/mongo"
//
//	"github.com/google/go-github/github"
//	"golang.org/x/crypto/ssh"
//	"gopkg.in/mgo.v2/bson"
//
//	"bytes"
//	"errors"
//	"fmt"
//	"io/ioutil"
//	"os/exec"
//	"os/user"
//	"strings"
//	"time"
//)
//
//const deploy_file = ".deploy.yml"
//
//func ProcessHook(event, body string) {
//	switch event {
//	case "pull_request":
//		pullRequestEvent, _ := ParsePullRequestEvent(body)
//		actions := map[string]bool{"opened": true, "reopened": true, "synchronize": true}
//		if actions[*pullRequestEvent.Action] {
//			branch := mongo.Branch{*pullRequestEvent.Repo.Owner.Login, *pullRequestEvent.Repo.Name, *pullRequestEvent.PullRequest.Head.SHA}
//			build := &mongo.Build{branch, pullRequestEvent, nil, bson.NewObjectId(), true, time.Now().Unix()}
//			client := build.Branch.GetRepository().GetGithubClient()
//			content, _ := getGithubFileContent(client, build.Branch, deploy_file)
//			content = []byte(strings.Replace(string(content), "{{sha}}", build.Branch.Sha, -1))
//			config, _ := GetYamlConfig(content)
//
//			runCommands(build, client, event, config)
//		} else {
//			fmt.Print("Skipping pull request event type: " + *pullRequestEvent.Action)
//		}
//	case "push":
//		pushEvent, err := ParsePushEvent(body)
//		if err != nil {
//			panic(err)
//		}
//		branch := mongo.Branch{*pushEvent.Repo.Owner.Name, *pushEvent.Repo.Name, *pushEvent.After}
//		build := &mongo.Build{branch, pushEvent, nil, bson.NewObjectId(), true, time.Now().Unix()}
//		client := build.Branch.GetRepository().GetGithubClient()
//		content, _ := getGithubFileContent(client, build.Branch, deploy_file)
//		content = []byte(strings.Replace(string(content), "{{sha}}", build.Branch.Sha, -1))
//		config, _ := GetYamlConfig(content)
//
//		if config.Push.Branch == *pushEvent.Ref {
//			runCommands(build, client, event, config)
//		}
//	default:
//		fmt.Println("Not supported event: " + event)
//		fmt.Println(body)
//	}
//}
//
//func getGithubFileContent(client *github.Client, br mongo.Branch, filename string) (fileContent []byte, err error) {
//	repoOptions := &github.RepositoryContentGetOptions{br.Sha}
//	a, _, _, err := client.Repositories.GetContents(br.Owner, br.Repo, filename, repoOptions)
//	if err != nil {
//		return
//	}
//
//	fileContent, err = a.Decode()
//	return
//}
//
//// Statuses: pending, success, error, or failure
//func setGitStatus(client *github.Client, build *mongo.Build, state string) (out string, err error) {
//	context := "continuous-integration/gorgon-ci"
//	config, _ := GetServerConfig()
//	url := config.Url + "/repos/" + build.Branch.Owner + "/" + build.Branch.Repo + "/" + build.Id.Hex()
//	status := &github.RepoStatus{State: &state, Context: &context, TargetURL: &url}
//	repoStatus, _, err := client.Repositories.CreateStatus(build.Branch.Owner, build.Branch.Repo, build.Branch.Sha, status)
//	out = "Success. Current github branch status: " + *repoStatus.State
//	return
//}
//
//func getKeyFile() (key ssh.Signer, err error) {
//	usr, _ := user.Current()
//	file := usr.HomeDir + "/.ssh/id_rsa"
//	buf, err := ioutil.ReadFile(file)
//	if err != nil {
//		return
//	}
//	key, err = ssh.ParsePrivateKey(buf)
//	return
//}


//
//func runCommands(build *mongo.Build, client *github.Client, event string, config DeployConfig) {
//	var commands []mongo.Command
//	var out string
//	var err error
//	var actions []map[string]string
//
//	if event == "push" {
//		actions = config.Push.Commands
//	} else if event == "pull_request" {
//		actions = config.Pull_request.Commands
//	}
//
//	for _, command := range actions {
//		for commandType, actionStr := range command {
//			if commandType == "status" {
//				out, err = setGitStatus(client, build, actionStr)
//			}
//			if commandType == "exec" {
//				out, err = execCommand(actionStr)
//			}
//			if commandType == "ssh" {
//				out, err = execSshCommand(config.Host, actionStr)
//			}
//			if err != nil {
//				commands = append(commands, mongo.Command{commandType, actionStr, out, err.Error()})
//			} else {
//				commands = append(commands, mongo.Command{Type: commandType, Action: actionStr, Out: out})
//			}
//		}
//		if err != nil {
//			out, err = setGitStatus(client, build, "error")
//			build.Success = false
//			if err != nil {
//				commands = append(commands, mongo.Command{"status", "error", out, err.Error()})
//			} else {
//				commands = append(commands, mongo.Command{Type: "status", Action: "error", Out: out})
//			}
//			break
//		}
//	}
//	build.Commands = commands
//	build.Save()
//}
//
//func execSshCommand(host string, command string) (out string, err error) {
//	client, err := getSshClient(host)
//	if err != nil {
//		return
//	}
//	defer client.Close()
//
//	session, err := client.NewSession()
//	if err != nil {
//		return
//	}
//	defer session.Close()
//
//	var outBuf bytes.Buffer
//	var errBuf bytes.Buffer
//	session.Stdout = &outBuf
//	session.Stderr = &errBuf
//	err = session.Run(command)
//	if err != nil {
//		return
//	}
//
//	return outBuf.String(), errors.New(errBuf.String())
//}
//
//func execCommand(cmd string) (string, error) {
//	parts := strings.Fields(cmd)
//	head := parts[0]
//	parts = parts[1:len(parts)]
//
//	out, err := exec.Command(head, parts...).Output()
//
//	return string(out), err
//}
//
//func setGithubHook(user, repo string) map[string]interface{} {
//	branch := mongo.Branch{user, repo, ""}
//	client := branch.GetRepository().GetGithubClient()
//
//	config, _ := GetServerConfig()
//	url  := config.Url + "/hooks"
//	hook := &github.Hook{Name: github.String("web"), Active: github.Bool(true), Events: []string{"pull_request", "push"}, Config: map[string]interface {}{"url": url}}
//
//	hook, response, error := client.Repositories.CreateHook(user, repo, hook)
//
//	result := make(map[string]interface{})
//	result["hook"]     = hook
//	result["response"] = response
//	result["error"]    = error
//
//	return result
//}
