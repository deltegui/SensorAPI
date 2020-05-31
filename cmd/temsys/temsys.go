package main

import (
	"reflect"
	"sensorapi/src/configuration"
	"sensorapi/src/controllers"
	"sensorapi/src/domain"

	"github.com/deltegui/phoenix"
	"github.com/deltegui/phoenix/injector"
)

func main() {
	phoenix.Configure().
		SetProjectInfo("sensorapi", "0.1.0").
		EnableStaticServer().
		EnableTemplates().
		EnableLogoFile()
	config := configuration.Load()
	controllers.Register(config)
	injector.ShowAvailableBuilders()
	reporterType := reflect.TypeOf((*domain.Reporter)(nil)).Elem()
	reporter := injector.GetByType(reporterType).(domain.Reporter)
	go reporter.Start()
	phoenix.Run(config.ListenURL)
}
