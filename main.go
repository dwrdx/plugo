package main

import (
	"plugo/plugo"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	pluginManager := plugo.NewPluginManager("localhost", "8080")

	pluginManager.LoadPlugins("./plugins")

	pluginManager.Serve()
}
