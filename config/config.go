package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/decosblockchain/audittrail-client/logging"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ServerUrl  string
	ListenPort int
}

var configuration Configuration

func Init() {
	err := gonfig.GetConf(path.Join(BaseDir(), "config.json"), &configuration)
	if err != nil {
		logging.Error.Println(err)
		os.Exit(6)
	}
}

func ListenPort() int {
	return configuration.ListenPort
}

func ServerUrl() string {
	return configuration.ServerUrl
}

func SendUrl() string {
	return ServerUrl() + "send"
}

func BaseDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func EnsurePathsExist() {
	if _, err := os.Stat(path.Join(BaseDir(), "log")); os.IsNotExist(err) {
		os.Mkdir(path.Join(BaseDir(), "log"), 0775)
	}

	if _, err := os.Stat(path.Join(BaseDir(), "log", "archive")); os.IsNotExist(err) {
		os.Mkdir(path.Join(BaseDir(), "log", "archive"), 0775)
	}

	if _, err := os.Stat(path.Join(BaseDir(), "data")); os.IsNotExist(err) {
		os.Mkdir(path.Join(BaseDir(), "data"), 0775)
	}

}
