package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	DB        string
	SiteName  string
	SecretKey string
	SitePort  string
	SiteUrl   string
	JwtSecret string
}

func LoadConf() *Config {
	source, errRead := ioutil.ReadFile("./conf/conf.yml")
	if errRead != nil {
		log.Println(errRead)
	}
	var config *Config
	err := yaml.Unmarshal(source, &config)
	if err != nil {
		log.Println(err)
	}
	return config
}
