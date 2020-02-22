package controllers

import (
	"sensorapi/src/configuration"
	"sensorapi/src/persistence"

	"github.com/deltegui/locomotive"
	"github.com/deltegui/locomotive/injector"
)

func Register(config configuration.Configuration) {
	conn := persistence.NewSqlxConnection(config)
	injector.Add(func() *persistence.SqlxConnection { return &conn })
	injector.Add(persistence.NewSqlxReportTypeRepo)
	injector.Add(persistence.NewSqlxSensorRepo)

	locomotive.MapRoot(NewErrorController)
	locomotive.Map("/reporttypes", NewReportTypeController)
	locomotive.Map("/sensors", NewSensorController)
}
