package config

import (
	"byung/log"
	"encoding/json"
	"os"
)

type Config struct {
	DataDirectory             string
	ListenAddress             string
	CertFile                  string
	KeyFile                   string
	Https                     bool
	Statics                   string
	DefaultAvatar             string
	DefaultArticleAttachImage string
	MaxUsers                  int
}

var Conf Config
var confFile = "./conf/config.json"

func init() {
	configFile, err := os.Open(confFile)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&Conf)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
