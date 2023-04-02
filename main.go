package main

import (
	"plugo/plugo"
)

func main() {
	pluginManager := plugo.NewPluginManager("localhost", "8080")

	pluginManager.LoadPlugins("./plugins")

	pluginManager.Serve()
}
