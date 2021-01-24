package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Hosts []struct {
		HostName   string `yaml:"name"`
		Connection string `yaml:"connection"`
		UserName   string `yaml:"username"`
		PassWord   string `yaml:"password"`
		Commands   []struct {
			Name       string `yaml:"name"`
			String     string `yaml:"string"`
			UserInput  bool   `yaml:"userinput"`
			WhiteSpace bool   `yaml:"whitespace"`
		} `yaml:"commands"`
	} `yaml:"hosts"`
}

func GetConfig() {
	if _, err := os.Stat("./config/config.yml"); err == nil { // check if config file exists
		yamlFile, err := ioutil.ReadFile("./config/config.yml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) { // config file not included, use embedded config
		yamlFile, err := Asset("config/config.yml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Schrodinger: file may or may not exist. See err for details.")
		// panic(err)
	}
}
