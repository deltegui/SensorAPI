package controllers

import (
	"sensorapi/src/builders"
	"sensorapi/src/configuration"
	"sensorapi/src/cronscheluder"
	"sensorapi/src/domain"
	"sensorapi/src/persistence"
	"sensorapi/src/validator"

	"github.com/deltegui/locomotive"
	"github.com/deltegui/locomotive/injector"
)

func registerUseCases() {
	injector.Add(domain.NewGetAllSensorsCase)
	injector.Add(domain.NewAllSensorNowCase)
	injector.Add(domain.NewDeleteSensorCase)
	injector.Add(domain.NewGetSensorCase)
	injector.Add(domain.NewSaveSensorCase)
	injector.Add(domain.NewSensorNowCase)
	injector.Add(domain.NewUpdateSensorCase)
}

func registerDependencies() {
	injector.Add(builders.NewHttpSensorBuilder)
	injector.Add(validator.NewPlaygroundValidator)
	injector.Add(persistence.NewSqlxReportTypeRepo)
	injector.Add(persistence.NewSqlxSensorRepo)
	injector.Add(persistence.NewSqlxReportRepo)
	injector.Add(domain.NewReporter)
	injector.Add(cronscheluder.NewCronScheluder)
}

func Register(config configuration.Configuration) {
	registerUseCases()
	registerDependencies()
	conn := persistence.NewSqlxConnection(config)
	injector.Add(func() *persistence.SqlxConnection { return &conn })

	locomotive.MapRoot(NewErrorController)
	locomotive.Map("/reporttypes", NewReportTypeController)
	locomotive.Map("/sensors", NewSensorsController)
	locomotive.Map("/sensor", NewSensorController)
	locomotive.Map("/reports", NewReportController)
}
