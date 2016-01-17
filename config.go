package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type jsonconfig struct {
	Nick     string
	Server   string
	Port     string
	Channels []string
	// For later use
	//Networks NetworksType
}

type NetworksType struct {
	Server   string
	Port     string
	Channels []string
}

func get_config(conf string) jsonconfig {
	file, e := ioutil.ReadFile("./config")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var config jsonconfig
	json.Unmarshal(file, &config)
	return config
}
