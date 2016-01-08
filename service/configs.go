package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
    "golang.org/x/oauth2"
	"github.com/alexkomrakov/gohub/models"
)

const (
	config_file = "config.yml"
)

type ServerConfig struct {
	Adress           string
    DeployFile       string
    SessionSecretKey string
    Logs   struct {
        Error string
        Gohub string
    }
    Oauth            oauth2.Config
    OauthStateString string
    Url              string
	Events           []string
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

func GetYamlConfig(file []byte) (config map[string]models.DeployScenario, err error) {
	err = yaml.Unmarshal(file, &config)
	return
}
