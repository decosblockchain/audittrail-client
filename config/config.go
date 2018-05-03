package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ServerUrl string
}

var configuration Configuration

func Init() {
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}
}

func ServerUrl() string {
	return configuration.ServerUrl
}

func SendUrl() string {
	return ServerUrl() + "send"
}
