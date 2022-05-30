package config

import (
	"fmt"
	"github.com/jeevan86/lf4go/factory"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type config struct {
	Logging factory.Logging `yaml:"logging"`
}

var configYml = "./config/config.yml"
var Config = loadConfigYml(configYml)

func loadConfigYml(path string) *config {
	if len(path) == 0 {
		path = configYml
	}
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("解析配置错误：%s", err.Error()))
	}
	var c = new(config)
	err = yaml.Unmarshal(yml, c)
	if err != nil {
		fmt.Println(fmt.Sprintf("解析配置错误：%s", err.Error()))
		return nil
	}
	return c
}
