package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	config_file = ".config.yml"
)

type ServerConfig struct {
	Adress string
	Deploy string
    Logs   struct {
        Error string
        Gohub string
    }
}

type DeployScenario struct {
	Branch string
	Host   string
	Commands []map[string]string
}

type DeployConfig struct {
	Scenarios []map[string]DeployScenario
}

func GetServerConfig() (config ServerConfig) {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, &config)
    if err != nil {
        panic(err)
    }

	return
}

func GetYamlConfig(file []byte) (config DeployConfig, err error) {
	err = yaml.Unmarshal(file, &config)
	return
}
