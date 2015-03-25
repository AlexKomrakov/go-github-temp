package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	config_file = ".config.yml"
)

type ServerConfig struct {
	Adress string
}

//type DeployConfig struct {
//	Host         string
//	Pull_request struct{ Commands []map[string]string }
//	Push         struct {
//		Branch   string
//		Commands []map[string]string
//	}
//}

func GetServerConfig() (config ServerConfig, err error) {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &config)
	return
}

//func GetYamlConfig(file []byte) (config DeployConfig, err error) {
//	err = yaml.Unmarshal(file, &config)
//	return
//}
