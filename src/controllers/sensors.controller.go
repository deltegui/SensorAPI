package controllers

import (
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/locomotive"
)

type SensorsController struct {
	getAllSensorsCase domain.UseCase
	allSensorNowCase  domain.UseCase
}

func NewSensorsController(getAllSensorsCase domain.GetAllSensorsCase, allSensorNowCase domain.AllSensorNowCase) SensorsController {
	return SensorsController{
		getAllSensorsCase,
		allSensorNowCase,
	}
}

func (controller SensorsController) GetAllSensors(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	wantDeleted := req.URL.Query()["deleted"]
	var reqCase domain.GetAllRequest
	reqCase.WantDeleted = !(len(wantDeleted) < 1 || len(wantDeleted[0]) == 0 || wantDeleted[0] == "false")
	controller.getAllSensorsCase.Exec(presenter, reqCase)
}

func (controller SensorsController) AllSensorNow(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	controller.allSensorNowCase.Exec(presenter, nil)
}

func (controller SensorsController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetAllSensors, Endpoint: ""},
		{Method: locomotive.Get, Handler: controller.AllSensorNow, Endpoint: "/all/now"},
	}
}
