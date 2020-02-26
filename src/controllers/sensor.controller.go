package controllers

import (
	"encoding/json"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/locomotive"
	"github.com/gorilla/mux"
)

type SensorController struct {
	saveSensorCase   domain.UseCase
	deleteSensorCase domain.UseCase
	updateSensorCase domain.UseCase
	sensorNowCase    domain.UseCase
	getSensorCase    domain.UseCase
	SensorBuilder    domain.SensorBuilder
}

func NewSensorController(
	saveSensorCase domain.SaveSensorCase,
	deleteSensorCase domain.DeleteSensorCase,
	updateSensorCase domain.UpdateSensorCase,
	sensorNowCase domain.SensorNowCase,
	getSensorCase domain.GetSensorCase,
	builder domain.SensorBuilder) SensorController {
	return SensorController{
		saveSensorCase,
		deleteSensorCase,
		updateSensorCase,
		sensorNowCase,
		getSensorCase,
		builder,
	}
}

func (controller SensorController) SaveSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var viewModel domain.SensorViewModel
	if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	viewModel.SensorBuilder = controller.SensorBuilder
	controller.saveSensorCase.Exec(presenter, viewModel)
}

func (controller SensorController) GetSensorByName(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	controller.getSensorCase.Exec(presenter, sensorName)
}

func (controller SensorController) DeleteSensorByName(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	controller.deleteSensorCase.Exec(presenter, sensorName)
}

func (controller SensorController) UpdateSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var viewModel domain.SensorViewModel
	if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	viewModel.SensorBuilder = controller.SensorBuilder
	controller.updateSensorCase.Exec(presenter, viewModel)
}

func (controller SensorController) SensorNow(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	controller.sensorNowCase.Exec(presenter, sensorName)
}

func (controller SensorController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Post, Handler: controller.SaveSensor, Endpoint: ""},
		{Method: locomotive.Get, Handler: controller.GetSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Delete, Handler: controller.DeleteSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Post, Handler: controller.UpdateSensor, Endpoint: "/update"},
		{Method: locomotive.Get, Handler: controller.SensorNow, Endpoint: "/{name}/now"},
	}
}
