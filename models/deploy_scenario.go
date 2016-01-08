package models

type DeployScenario struct {
	Branch string
	Host   string
	Commands []map[string]string
	Error    []map[string]string
}

