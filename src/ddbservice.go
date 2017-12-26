/*
author Vyacheslav Kryuchenko
*/
package main

import (
	"flag"
	"helpers"
	"web/app"
)

// go-bindata -pkg web -prefix "src/web/resources/" -o src/web/resources.go src/web/resources/...

var (
	configPath = *flag.String("config", "appconfig.json", "configuration file path")
)

func main() {
	appConfig := helpers.AppConfig{}
	appConfig.Read(configPath)
	provider := app.Provider{
		Listen:          appConfig.Listen,
		ApplicationName: appConfig.Appname,
		Secret:          appConfig.Secret,
	}
	app.StartServer(&provider)
}
