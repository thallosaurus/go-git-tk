package config

import (
	"encoding/json"
	"go-git-tk/pkg/layouts"
	"log"
	"os"
)

const DEFAULT_CONFIG_PATH = "/etc/gittk/config.json"

type Config struct {
	EnableSSH        bool
	SSHBaseDomain    string
	SSHUser          string
	ShowBorders      bool
	ShowFullRepoPath bool
	ColorTint        string
	RepoWorkdir      string
}

var Conf Config

func InitConfig(path ...string) {
	p := DEFAULT_CONFIG_PATH
	if len(path) == 1 {
		p = path[0]
	}
	Conf = ReadConfig(p)
	layouts.SetAppColor(Conf.ColorTint)
}

func ReadConfig(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var data Config
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
