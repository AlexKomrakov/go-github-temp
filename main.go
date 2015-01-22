package main

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

const (
	server   = "localhost:8080"
	log_file = "/gohub.log"

	github_user  = "alexkomrakov"
	github_repo  = "gohub"
	github_ref   = "fe66d4fd869c82487b64b25ac85c6b278de4356br"
	github_token = "" // https://gist.github.com/AlexKomrakov/38323424693224a3c03d

	ssh_host = "komrakov-stage.smart-crowd.ru:22"
	ssh_user = "komrakov"

	deploy_config = ".deploy.yml"
)

var (
	Error *log.Logger
)

func main() {
	file := loggerInit(log_file)
	defer file.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/github", githubHandler)

	http.Handle("/", router)
	err := http.ListenAndServe(server, nil)
	if err != nil {
		Error.Println("Error on starting server: %v", err)
	}
}

func githubHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello github")

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: github_token},
	}
	client := github.NewClient(transport.Client())

	listOptions := new(github.ListOptions)
	json.Unmarshal([]byte(`{"page": 0, "perPage": 5}`), &listOptions)
	z, c, _ := client.Repositories.ListStatuses(github_user, github_repo, github_ref, listOptions)

	io.WriteString(w, fmt.Sprintf("%v", z))
	io.WriteString(w, fmt.Sprintf("%v", c))

	setGitStatus(client, "success")
}

func setGitStatus(client *github.Client, state string) {
	context := "continuous-integration/gorgon-ci"
	status := &github.RepoStatus{State: &state, Context: &context}
	_, resp, err := client.Repositories.CreateStatus(github_user, github_repo, github_ref, status)
	fmt.Print(resp)
	fmt.Print(err)
}

func getGithubFileContent(client *github.Client, owner, repo, filename, branch string) ([]byte, error) {
	repoOptions := &github.RepositoryContentGetOptions{branch}
	a, _, _, err1 := client.Repositories.GetContents(owner, repo, filename, repoOptions)
	if err1 != nil {
		Error.Println("Error on getting file from github: %v", err1)
		return nil, err1
	}

	fileContent, err2 := a.Decode()
	if err2 != nil {
		Error.Println("Error on decoding file from github: %v", err2)
		return nil, err2
	}

	return fileContent, nil
}

type ymlConfig struct {
	Host     []interface{}
	Commands []interface{}
}

func runCommands(client *github.Client, config ymlConfig) {
	sshClient := getSshClient()
	defer sshClient.Close()

	for _, command := range config.Commands {

		switch v := command.(type) {
		case map[interface{}]interface{}:
			ma := command.(map[interface{}]interface{})
			for commandType, action := range ma {
				actionStr := action.(string)
				if commandType == "status" {
					setGitStatus(client, actionStr)
				}
				if commandType == "ssh" {
					out, err := execSshCommand(sshClient, actionStr)
					fmt.Println(out.String())
					fmt.Println(err.String())
				}
			}
		default:
			Error.Printf("Error on run yaml config commands. %v", v)
		}

		//
		//		fmt.Print(command)
		//
		//		str := command.(string)
		//		out, _ := execSshCommand(client, str)
	}
}

func readYamlConfig(file []byte) (ymlConfig, error) {
	config := ymlConfig{}
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		Error.Println("Error on reading yaml config")
		return config, err
	}

	return config, nil
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello world")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Error.Println("error reading body")
	}
	Error.Println(string(body))
}

func loggerInit(filename string) (file *os.File) {
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		fmt.Print("Error on parsing abs path for logger: %v", err1)
	}

	filename = filepath.Join(dir, filename)

	f, err2 := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err2 != nil {
		fmt.Print("Error opening file: %v", err2)
	}
	err3 := f.Truncate(0)
	if err3 != nil {
		fmt.Print("Error on clearing log file: %v", err3)
	}

	Error = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	Error.Println("Logger start")

	return f
}

//TODO Defer close
func getSshClient() *ssh.Client {
	key, err := getKeyFile()
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: ssh_user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(ssh_pass),
			ssh.PublicKeys(key),
		},
	}
	client, err := ssh.Dial("tcp", ssh_host, config)
	if err != nil {
		fmt.Printf("unable to connect: %s", err)
	}

	return client
}

func getKeyFile() (key ssh.Signer, err error) {
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func execSshCommand(client *ssh.Client, command string) (bytes.Buffer, bytes.Buffer) {
	session, err1 := client.NewSession()
	if err1 != nil {
		panic("Failed to create session: " + err1.Error())
	}
	defer session.Close()

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	session.Run(command) //TODO Check error
	//	if err2 != nil {
	//		panic("Failed to run: " + err2.Error())
	//	}

	return outBuf, errBuf
}
