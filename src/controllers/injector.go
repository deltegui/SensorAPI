package controllers

import (
	"sensorapi/src/configuration"
	"sensorapi/src/persistence"

	"github.com/deltegui/locomotive"
	"github.com/deltegui/locomotive/injector"
)

func Register(config configuration.Configuration) {
	injector.Add(func() persistence.SqlxConnection {
		return persistence.NewSqlxConnection(config)
	})
	injector.Add(persistence.NewSqlxReportTypeRepo)

	locomotive.MapRoot(NewErrorController)
	locomotive.Map("/reporttypes", NewReportTypeController)
}
