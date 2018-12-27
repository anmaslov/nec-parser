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

type Redis struct{
	Addr 		string
	Password	string
	Db			int
}

type Configuration struct {
	Env 	  	string
	//Phones    	[]Stantion
	Database	Db
	Redis		Redis
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
