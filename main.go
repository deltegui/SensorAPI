package main

import (
	"sensorapi/src/configuration"
	"sensorapi/src/controllers"

	"github.com/deltegui/locomotive"
	"github.com/deltegui/locomotive/vars"
)

func setVariables() {
	vars.Name = "sensorapi"
	vars.Version = "0.1.0"
	vars.EnableStaticServer = false
	vars.EnableTemplates = false
}

func main() {
	setVariables()
	config := configuration.Load()
	controllers.Register(config)
	locomotive.Run(config.ListenURL)
}
