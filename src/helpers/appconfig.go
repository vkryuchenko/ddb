/*
author Vyacheslav Kryuchenko
*/
package helpers

import (
	"encoding/json"
	"log"
	"os"
)

type AppConfig struct {
	Listen     string         `json:"listen"`
	Log        string         `json:"log"`
	Appname    string         `json:"appname"`
	Secret     string         `json:"secret"`
	Target     string         `json:"target"`
	Apiversion string         `json:"apiversion"`
	Ports      string         `json:"ports"`
	DB         PostgresConfig `json:"db"`
}

func (ac *AppConfig) Read(configPath string) {
	log.Println("Check config.")
	cf, err := os.Open(configPath)
	defer cf.Close()
	if err != nil {
		log.Fatalf("can't read config file!\n%s", err)
	}

	decoder := json.NewDecoder(cf)
	decodeErr := decoder.Decode(ac)
	if decodeErr != nil {
		log.Fatalf("can't decode %s\n%s", configPath, decodeErr)
	}
	if ac.Log != "" {
		log.Printf("logging into %s", ac.Log)
		logFile, err := os.OpenFile(ac.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("can't open log file -- %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	} else {
		log.Print("logging into STDOUT")
	}
}
