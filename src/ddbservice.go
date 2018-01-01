/*
author Vyacheslav Kryuchenko
*/
package main

import (
	"flag"
	"helpers"
	"log"
	"web/app"
)

// go-bindata -pkg web -prefix "src/web/resources/" -o src/web/resources.go src/web/resources/...

var (
	configPath = *flag.String("config", "appconfig.json", "configuration file path")
)

func main() {
	appConfig := helpers.AppConfig{}
	appConfig.Read(configPath)
	err := appConfig.DB.InitDatabase()
	if err != nil {
		log.Panic(err)
	}
	ldapClient, err := appConfig.LDAP.Init()
	if err != nil {
		log.Panic(err)
	}
	provider := app.Provider{
		Develop:         appConfig.Develop,
		Listen:          appConfig.Listen,
		ApplicationName: appConfig.Appname,
		Secret:          appConfig.Secret,
		LDAPClient:      &ldapClient,
	}
	app.StartServer(&provider)
}
