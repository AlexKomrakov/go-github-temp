package main

import (
	"bytes"
	"strings"
	//"code.google.com/p/goauth2/oauth"
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
	"code.google.com/p/goauth2/oauth"
)

const (
	config_file = ".config.yml"
	deploy_file = ".deploy.yml"
)

var (
	Error        *log.Logger
	github_token  string
)

func readConfig() (map[string]string) {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		panic(err)
	}

	var config map[string]string
	err2 := yaml.Unmarshal(b, &config)
	if err2 != nil {
		fmt.Println("Error on reading yaml config")
	}

	return config
}

func start() {
	config := readConfig()
	github_token = config["github_token"]

	file := loggerInit(config["log_file"])
	defer file.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/github", githubHandler)

	http.Handle("/", router)
	err := http.ListenAndServe(config["server"], nil)
	if err != nil {
		Error.Println("Error on starting server: %v", err)
	}
}

func githubHandler(w http.ResponseWriter, req *http.Request) {
	body := req.FormValue("payload")
//	body := `{"action":"opened","number":5,"pull_request":{"url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/5","id":27982242,"html_url":"https://github.com/AlexKomrakov/gohub/pull/5","diff_url":"https://github.com/AlexKomrakov/gohub/pull/5.diff","patch_url":"https://github.com/AlexKomrakov/gohub/pull/5.patch","issue_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/5","number":5,"state":"open","locked":false,"title":"+ save","user":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false},"body":"","created_at":"2015-01-24T21:01:22Z","updated_at":"2015-01-24T21:01:22Z","closed_at":null,"merged_at":null,"merge_commit_sha":null,"assignee":null,"milestone":null,"commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/5/commits","review_comments_url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/5/comments","review_comment_url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/comments/{number}","comments_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/5/comments","statuses_url":"https://api.github.com/repos/AlexKomrakov/gohub/statuses/35c18d42f4ff038e7503125b1263dddccb6e3204","head":{"label":"AlexKomrakov:config","ref":"config","sha":"35c18d42f4ff038e7503125b1263dddccb6e3204","user":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false},"repo":{"id":29361502,"name":"gohub","full_name":"AlexKomrakov/gohub","owner":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false},"private":false,"html_url":"https://github.com/AlexKomrakov/gohub","description":"","fork":false,"url":"https://api.github.com/repos/AlexKomrakov/gohub","forks_url":"https://api.github.com/repos/AlexKomrakov/gohub/forks","keys_url":"https://api.github.com/repos/AlexKomrakov/gohub/keys{/key_id}","collaborators_url":"https://api.github.com/repos/AlexKomrakov/gohub/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/AlexKomrakov/gohub/teams","hooks_url":"https://api.github.com/repos/AlexKomrakov/gohub/hooks","issue_events_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/events{/number}","events_url":"https://api.github.com/repos/AlexKomrakov/gohub/events","assignees_url":"https://api.github.com/repos/AlexKomrakov/gohub/assignees{/user}","branches_url":"https://api.github.com/repos/AlexKomrakov/gohub/branches{/branch}","tags_url":"https://api.github.com/repos/AlexKomrakov/gohub/tags","blobs_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/refs{/sha}","trees_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/trees{/sha}","statuses_url":"https://api.github.com/repos/AlexKomrakov/gohub/statuses/{sha}","languages_url":"https://api.github.com/repos/AlexKomrakov/gohub/languages","stargazers_url":"https://api.github.com/repos/AlexKomrakov/gohub/stargazers","contributors_url":"https://api.github.com/repos/AlexKomrakov/gohub/contributors","subscribers_url":"https://api.github.com/repos/AlexKomrakov/gohub/subscribers","subscription_url":"https://api.github.com/repos/AlexKomrakov/gohub/subscription","commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/commits{/sha}","git_commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/commits{/sha}","comments_url":"https://api.github.com/repos/AlexKomrakov/gohub/comments{/number}","issue_comment_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/comments/{number}","contents_url":"https://api.github.com/repos/AlexKomrakov/gohub/contents/{+path}","compare_url":"https://api.github.com/repos/AlexKomrakov/gohub/compare/{base}...{head}","merges_url":"https://api.github.com/repos/AlexKomrakov/gohub/merges","archive_url":"https://api.github.com/repos/AlexKomrakov/gohub/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/AlexKomrakov/gohub/downloads","issues_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues{/number}","pulls_url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls{/number}","milestones_url":"https://api.github.com/repos/AlexKomrakov/gohub/milestones{/number}","notifications_url":"https://api.github.com/repos/AlexKomrakov/gohub/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/AlexKomrakov/gohub/labels{/name}","releases_url":"https://api.github.com/repos/AlexKomrakov/gohub/releases{/id}","created_at":"2015-01-16T18:21:40Z","updated_at":"2015-01-23T19:29:59Z","pushed_at":"2015-01-24T21:01:11Z","git_url":"git://github.com/AlexKomrakov/gohub.git","ssh_url":"git@github.com:AlexKomrakov/gohub.git","clone_url":"https://github.com/AlexKomrakov/gohub.git","svn_url":"https://github.com/AlexKomrakov/gohub","homepage":null,"size":236,"stargazers_count":0,"watchers_count":0,"language":"Go","has_issues":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"open_issues_count":2,"forks":0,"open_issues":2,"watchers":0,"default_branch":"master"}},"base":{"label":"AlexKomrakov:master","ref":"master","sha":"87dafdec25a7e38f5b69f4268efac3ab869b076f","user":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false},"repo":{"id":29361502,"name":"gohub","full_name":"AlexKomrakov/gohub","owner":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false},"private":false,"html_url":"https://github.com/AlexKomrakov/gohub","description":"","fork":false,"url":"https://api.github.com/repos/AlexKomrakov/gohub","forks_url":"https://api.github.com/repos/AlexKomrakov/gohub/forks","keys_url":"https://api.github.com/repos/AlexKomrakov/gohub/keys{/key_id}","collaborators_url":"https://api.github.com/repos/AlexKomrakov/gohub/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/AlexKomrakov/gohub/teams","hooks_url":"https://api.github.com/repos/AlexKomrakov/gohub/hooks","issue_events_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/events{/number}","events_url":"https://api.github.com/repos/AlexKomrakov/gohub/events","assignees_url":"https://api.github.com/repos/AlexKomrakov/gohub/assignees{/user}","branches_url":"https://api.github.com/repos/AlexKomrakov/gohub/branches{/branch}","tags_url":"https://api.github.com/repos/AlexKomrakov/gohub/tags","blobs_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/refs{/sha}","trees_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/trees{/sha}","statuses_url":"https://api.github.com/repos/AlexKomrakov/gohub/statuses/{sha}","languages_url":"https://api.github.com/repos/AlexKomrakov/gohub/languages","stargazers_url":"https://api.github.com/repos/AlexKomrakov/gohub/stargazers","contributors_url":"https://api.github.com/repos/AlexKomrakov/gohub/contributors","subscribers_url":"https://api.github.com/repos/AlexKomrakov/gohub/subscribers","subscription_url":"https://api.github.com/repos/AlexKomrakov/gohub/subscription","commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/commits{/sha}","git_commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/commits{/sha}","comments_url":"https://api.github.com/repos/AlexKomrakov/gohub/comments{/number}","issue_comment_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/comments/{number}","contents_url":"https://api.github.com/repos/AlexKomrakov/gohub/contents/{+path}","compare_url":"https://api.github.com/repos/AlexKomrakov/gohub/compare/{base}...{head}","merges_url":"https://api.github.com/repos/AlexKomrakov/gohub/merges","archive_url":"https://api.github.com/repos/AlexKomrakov/gohub/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/AlexKomrakov/gohub/downloads","issues_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues{/number}","pulls_url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls{/number}","milestones_url":"https://api.github.com/repos/AlexKomrakov/gohub/milestones{/number}","notifications_url":"https://api.github.com/repos/AlexKomrakov/gohub/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/AlexKomrakov/gohub/labels{/name}","releases_url":"https://api.github.com/repos/AlexKomrakov/gohub/releases{/id}","created_at":"2015-01-16T18:21:40Z","updated_at":"2015-01-23T19:29:59Z","pushed_at":"2015-01-24T21:01:11Z","git_url":"git://github.com/AlexKomrakov/gohub.git","ssh_url":"git@github.com:AlexKomrakov/gohub.git","clone_url":"https://github.com/AlexKomrakov/gohub.git","svn_url":"https://github.com/AlexKomrakov/gohub","homepage":null,"size":236,"stargazers_count":0,"watchers_count":0,"language":"Go","has_issues":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"open_issues_count":2,"forks":0,"open_issues":2,"watchers":0,"default_branch":"master"}},"_links":{"self":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/5"},"html":{"href":"https://github.com/AlexKomrakov/gohub/pull/5"},"issue":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/issues/5"},"comments":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/issues/5/comments"},"review_comments":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/5/comments"},"review_comment":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/comments/{number}"},"commits":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/pulls/5/commits"},"statuses":{"href":"https://api.github.com/repos/AlexKomrakov/gohub/statuses/35c18d42f4ff038e7503125b1263dddccb6e3204"}},"merged":false,"mergeable":null,"mergeable_state":"unknown","merged_by":null,"comments":0,"review_comments":0,"commits":1,"additions":78,"deletions":58,"changed_files":5},"repository":{"id":29361502,"name":"gohub","full_name":"AlexKomrakov/gohub","owner":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false},"private":false,"html_url":"https://github.com/AlexKomrakov/gohub","description":"","fork":false,"url":"https://api.github.com/repos/AlexKomrakov/gohub","forks_url":"https://api.github.com/repos/AlexKomrakov/gohub/forks","keys_url":"https://api.github.com/repos/AlexKomrakov/gohub/keys{/key_id}","collaborators_url":"https://api.github.com/repos/AlexKomrakov/gohub/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/AlexKomrakov/gohub/teams","hooks_url":"https://api.github.com/repos/AlexKomrakov/gohub/hooks","issue_events_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/events{/number}","events_url":"https://api.github.com/repos/AlexKomrakov/gohub/events","assignees_url":"https://api.github.com/repos/AlexKomrakov/gohub/assignees{/user}","branches_url":"https://api.github.com/repos/AlexKomrakov/gohub/branches{/branch}","tags_url":"https://api.github.com/repos/AlexKomrakov/gohub/tags","blobs_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/refs{/sha}","trees_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/trees{/sha}","statuses_url":"https://api.github.com/repos/AlexKomrakov/gohub/statuses/{sha}","languages_url":"https://api.github.com/repos/AlexKomrakov/gohub/languages","stargazers_url":"https://api.github.com/repos/AlexKomrakov/gohub/stargazers","contributors_url":"https://api.github.com/repos/AlexKomrakov/gohub/contributors","subscribers_url":"https://api.github.com/repos/AlexKomrakov/gohub/subscribers","subscription_url":"https://api.github.com/repos/AlexKomrakov/gohub/subscription","commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/commits{/sha}","git_commits_url":"https://api.github.com/repos/AlexKomrakov/gohub/git/commits{/sha}","comments_url":"https://api.github.com/repos/AlexKomrakov/gohub/comments{/number}","issue_comment_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues/comments/{number}","contents_url":"https://api.github.com/repos/AlexKomrakov/gohub/contents/{+path}","compare_url":"https://api.github.com/repos/AlexKomrakov/gohub/compare/{base}...{head}","merges_url":"https://api.github.com/repos/AlexKomrakov/gohub/merges","archive_url":"https://api.github.com/repos/AlexKomrakov/gohub/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/AlexKomrakov/gohub/downloads","issues_url":"https://api.github.com/repos/AlexKomrakov/gohub/issues{/number}","pulls_url":"https://api.github.com/repos/AlexKomrakov/gohub/pulls{/number}","milestones_url":"https://api.github.com/repos/AlexKomrakov/gohub/milestones{/number}","notifications_url":"https://api.github.com/repos/AlexKomrakov/gohub/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/AlexKomrakov/gohub/labels{/name}","releases_url":"https://api.github.com/repos/AlexKomrakov/gohub/releases{/id}","created_at":"2015-01-16T18:21:40Z","updated_at":"2015-01-23T19:29:59Z","pushed_at":"2015-01-24T21:01:11Z","git_url":"git://github.com/AlexKomrakov/gohub.git","ssh_url":"git@github.com:AlexKomrakov/gohub.git","clone_url":"https://github.com/AlexKomrakov/gohub.git","svn_url":"https://github.com/AlexKomrakov/gohub","homepage":null,"size":236,"stargazers_count":0,"watchers_count":0,"language":"Go","has_issues":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"open_issues_count":2,"forks":0,"open_issues":2,"watchers":0,"default_branch":"master"},"sender":{"login":"AlexKomrakov","id":7386252,"avatar_url":"https://avatars.githubusercontent.com/u/7386252?v=3","gravatar_id":"","url":"https://api.github.com/users/AlexKomrakov","html_url":"https://github.com/AlexKomrakov","followers_url":"https://api.github.com/users/AlexKomrakov/followers","following_url":"https://api.github.com/users/AlexKomrakov/following{/other_user}","gists_url":"https://api.github.com/users/AlexKomrakov/gists{/gist_id}","starred_url":"https://api.github.com/users/AlexKomrakov/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AlexKomrakov/subscriptions","organizations_url":"https://api.github.com/users/AlexKomrakov/orgs","repos_url":"https://api.github.com/users/AlexKomrakov/repos","events_url":"https://api.github.com/users/AlexKomrakov/events{/privacy}","received_events_url":"https://api.github.com/users/AlexKomrakov/received_events","type":"User","site_admin":false}}`
	var data map[string]interface {}
	json.Unmarshal([]byte(body), &data)

	pull_request := data["pull_request"].(map[string]interface {})
	head := pull_request["head"].(map[string]interface {})
	sha := head["sha"].(string)
	fmt.Println(sha)

	repository := data["repository"].(map[string]interface {})
	repo_name := repository["name"].(string)
	fmt.Println(repo_name)

	owner := repository["owner"].(map[string]interface {})
	owner_name := owner["login"].(string)
	fmt.Println(owner_name)

	currentBranch := branch{owner_name, repo_name, sha}

	transport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: github_token},
	}
	client := github.NewClient(transport.Client())

	content, _ := getGithubFileContent(client, currentBranch, deploy_file)
	conf, _ := readYamlConfig(content)

	runCommands(client, currentBranch, conf)
}

func setGitStatus(client *github.Client, br branch, state string) {
	context := "continuous-integration/gorgon-ci"
	status := &github.RepoStatus{State: &state, Context: &context}
	_, resp, err := client.Repositories.CreateStatus(br.Owner, br.Repo, br.Sha, status)
	fmt.Print(resp)
	fmt.Print(err)
}

func getGithubFileContent(client *github.Client, br branch, filename string) ([]byte, error) {
	repoOptions := &github.RepositoryContentGetOptions{br.Sha}
	a, _, _, err1 := client.Repositories.GetContents(br.Owner, br.Repo, filename, repoOptions)
	if err1 != nil {
		fmt.Println("Error on getting file from github: %v", err1)
		return nil, err1
	}

	fileContent, err2 := a.Decode()
	if err2 != nil {
		fmt.Println("Error on decoding file from github: %v", err2)
		return nil, err2
	}

	return fileContent, nil
}

type ymlConfig struct {
	Host     []interface{}
	Commands []interface{}
}

type branch struct {
	Owner string
	Repo  string
	Sha   string
}

func runCommands(client *github.Client, br branch, config ymlConfig) {
	sshClient := getSshClient(config.Host[0].(string))
	defer sshClient.Close()

	for _, command := range config.Commands {

		switch v := command.(type) {
		case map[interface{}]interface{}:
			ma := command.(map[interface{}]interface{})
			setGitStatus(client, br, "pending")
			for commandType, action := range ma {
				actionStr := action.(string)
				if commandType == "status" {
					setGitStatus(client, br, actionStr)
				}
				if commandType == "ssh" {
					out, err := execSshCommand(sshClient, actionStr)
					fmt.Println(out.String())
					fmt.Println(err.String())
				}
			}
			setGitStatus(client, br, "success")
		default:
			Error.Printf("Error on run yaml config commands. %v", v)
		}
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
func getSshClient(user_host string) *ssh.Client {
	key, err := getKeyFile()
	if err != nil {
		panic(err)
	}

	params := strings.Split(user_host, "@")
	if len(params) != 2 {
		panic("Wrong ssh user@host in config: " + user_host)
	}
	user := params[0]
	host := params[1]

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(ssh_pass),
			ssh.PublicKeys(key),
		},
	}
	client, err := ssh.Dial("tcp", host, config)
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
	err2 := session.Run(command) //TODO Check error
	if err2 != nil {
		panic("Error: " + err2.Error())
	}

	return outBuf, errBuf
}
