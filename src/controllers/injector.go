package controllers

import (
	"sensorapi/src/configuration"
	"sensorapi/src/cronscheluder"
	"sensorapi/src/domain"
	"sensorapi/src/persistence"
	"sensorapi/src/validator"

	"github.com/deltegui/locomotive"
	"github.com/deltegui/locomotive/injector"
)

func Register(config configuration.Configuration) {
	conn := persistence.NewSqlxConnection(config)
	injector.Add(func() *persistence.SqlxConnection { return &conn })
	injector.Add(validator.NewPlaygroundValidator)
	injector.Add(persistence.NewSqlxReportTypeRepo)
	injector.Add(persistence.NewSqlxSensorRepo)
	injector.Add(persistence.NewSqlxReportRepo)
	injector.Add(domain.NewReporter)
	injector.Add(cronscheluder.NewCronScheluder)

	locomotive.MapRoot(NewErrorController)
	locomotive.Map("/reporttypes", NewReportTypeController)
	locomotive.Map("/sensors", NewSensorController)
	locomotive.Map("/reports", NewReportController)
}
