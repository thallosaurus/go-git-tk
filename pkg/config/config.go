package config

import (
	"encoding/json"
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

var Conf = ReadConfig("./scripts/gittk/shell.json")

func ReadConfig(path ...string) Config {
	p := DEFAULT_CONFIG_PATH
	if len(path) == 1 {
		p = path[0]
	}

	b, err := os.ReadFile(p)
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
