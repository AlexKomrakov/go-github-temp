package server
//
//import (
//	"encoding/json"
//	"github.com/google/go-github/github"
//)
//
//type PushEvent struct {
//	HeadCommit *PushEventCommit   `json:"head_commit,omitempty"`
//	Forced     *bool              `json:"forced,omitempty"`
//	Created    *bool              `json:"created,omitempty"`
//	Deleted    *bool              `json:"deleted,omitempty"`
//	Ref        *string            `json:"ref,omitempty"`
//	Before     *string            `json:"before,omitempty"`
//	After      *string            `json:"after,omitempty"`
//	Compare    *string            `json:"compare,omitempty"`
//	Size       *int               `json:"size,omitempty"`
//	Commits    []PushEventCommit  `json:"commits,omitempty"`
//	Repo       *github.Repository `json:"repository,omitempty"`
//}
//
//type PushEventCommit struct {
//	ID       *string              `json:"id,omitempty"`
//	Message  *string              `json:"message,omitempty"`
//	Author   *github.CommitAuthor `json:"author,omitempty"`
//	URL      *string              `json:"url,omitempty"`
//	Distinct *bool                `json:"distinct,omitempty"`
//	Added    []string             `json:"added,omitempty"`
//	Removed  []string             `json:"removed,omitempty"`
//	Modified []string             `json:"modified,omitempty"`
//}
//
//func ParsePullRequestEvent(s string) (event github.PullRequestEvent, err error) {
//	err = json.Unmarshal([]byte(s), &event)
//	return
//}
//
//func ParsePushEvent(s string) (event PushEvent, err error) {
//	err = json.Unmarshal([]byte(s), &event)
//	return
//}
