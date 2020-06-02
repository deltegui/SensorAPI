package controllers

import (
	"sensorapi/src/builders"
	"sensorapi/src/configuration"
	"sensorapi/src/cronscheluder"
	"sensorapi/src/domain"
	"sensorapi/src/persistence"
	"sensorapi/src/validator"

	"github.com/deltegui/phoenix"
)

func registerUseCases(injector *phoenix.Injector) {
	injector.Add(domain.NewGetAllSensorsCase)
	injector.Add(domain.NewAllSensorNowCase)
	injector.Add(domain.NewDeleteSensorCase)
	injector.Add(domain.NewGetSensorCase)
	injector.Add(domain.NewSaveSensorCase)
	injector.Add(domain.NewSensorNowCase)
	injector.Add(domain.NewUpdateSensorCase)
	injector.Add(domain.NewGetReportsByDates)
	injector.Add(domain.NewLoginUserCase)
}

func registerDependencies(injector *phoenix.Injector) {
	injector.Add(builders.NewHttpSensorBuilder)
	injector.Add(validator.NewPlaygroundValidator)
	injector.Add(persistence.NewSqlxReportTypeRepo)
	injector.Add(persistence.NewSqlxSensorRepo)
	injector.Add(persistence.NewSqlxReportRepo)
	injector.Add(persistence.NewSqlxUserRepo)
	injector.Add(domain.NewReporter)
	injector.Add(cronscheluder.NewCronScheluder)
}

func Register(app phoenix.App, config configuration.Configuration) {
	registerUseCases(app.Injector)
	registerDependencies(app.Injector)
	conn := persistence.NewSqlxConnection(config)
	app.Injector.Add(func() *persistence.SqlxConnection { return &conn })
	app.Injector.Add(func() configuration.Configuration { return config })

	app.Map(phoenix.Mapping{Method: phoenix.Get, Builder: NotFound, Endpoint: "404"})
	registerReportTypesRoutes(app)
	registerSensorsRoutes(app)
	registerSensorRoutes(app)
	registerReportRoutes(app)
	registerUserRoutes(app)
}
