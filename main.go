package main

import (
	"reflect"
	"sensorapi/src/configuration"
	"sensorapi/src/controllers"
	"sensorapi/src/domain"

	"github.com/deltegui/locomotive"
	"github.com/deltegui/locomotive/injector"
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
	injector.ShowAvailableBuilders()
	reporterType := reflect.TypeOf((*domain.Reporter)(nil)).Elem()
	reporter := injector.GetByType(reporterType).(domain.Reporter)
	go reporter.Start()
	locomotive.Run(config.ListenURL)
}
