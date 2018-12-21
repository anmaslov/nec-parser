package main

import (
	"encoding/json"
	"log"
	"os"
)

type Db struct {
	Host		string
	Dbname		string
	Username	string
	Password	string
}

type Stantion struct {
	Name 	string
	Ip		string
	Port	string
}

type Configuration struct {
	Env 	  	string
	//Phones    	[]Stantion
	Database	Db
	Log			bool
}

func (conf *Configuration)loadConfig() {
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal("error load config: ", err)
	}
}
