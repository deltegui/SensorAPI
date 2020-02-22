package controllers

import (
	"encoding/json"
	"log"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/locomotive"
	"github.com/gorilla/mux"
)

type SensorController struct {
	sensorRepo domain.SensorRepo
}

func NewSensorController(sensorRepo domain.SensorRepo) SensorController {
	return SensorController{sensorRepo}
}

func (controller SensorController) GetAllSensors(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	presenter.Present(controller.sensorRepo.GetAll())
}

func (controller SensorController) SaveSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var sensor domain.Sensor
	if err := json.NewDecoder(req.Body).Decode(&sensor); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.sensorRepo.Save(sensor); err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	presenter.Present(sensor)
}

func (controller SensorController) GetSensorByName(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	sensor := controller.sensorRepo.GetByName(sensorName)
	if sensor == nil {
		presenter.PresentError(domain.SensorNotFoundErr)
		return
	}
	presenter.Present(sensor)
}

func (controller SensorController) DeleteSensorByName(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	sensor := controller.sensorRepo.GetByName(sensorName)
	if sensor == nil {
		presenter.PresentError(domain.SensorNotFoundErr)
		return
	}
	sensor.Deleted = true
	if controller.sensorRepo.Update(*sensor) {
		presenter.Present(struct{ Deleted bool }{true})
		return
	}
	presenter.PresentError(domain.InternalErr)
}

func (controller SensorController) UpdateSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var sensor domain.Sensor
	if err := json.NewDecoder(req.Body).Decode(&sensor); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if !controller.sensorRepo.Update(sensor) {
		presenter.PresentError(domain.UpdateErr)
		return
	}
	presenter.Present(sensor)
}

func (controller SensorController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetAllSensors, Endpoint: "/all"},
		{Method: locomotive.Post, Handler: controller.SaveSensor, Endpoint: "/"},
		{Method: locomotive.Get, Handler: controller.GetSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Delete, Handler: controller.DeleteSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Post, Handler: controller.UpdateSensor, Endpoint: "/update"},
	}
}
