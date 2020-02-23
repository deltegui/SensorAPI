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
	reporter   domain.Reporter
	validator  domain.Validator
}

func NewSensorController(validator domain.Validator, sensorRepo domain.SensorRepo, reporter domain.Reporter) SensorController {
	return SensorController{sensorRepo, reporter, validator}
}

func (controller SensorController) GetAllSensors(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	wantDeleted := req.URL.Query()["deleted"]
	if len(wantDeleted) < 1 || len(wantDeleted[0]) == 0 || wantDeleted[0] == "false" {
		presenter.Present(controller.sensorRepo.GetAll(domain.WithoutDeletedSensors))
		return
	}
	presenter.Present(controller.sensorRepo.GetAll(domain.WithDeletedSensors))
}

func (controller SensorController) SaveSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var sensor domain.Sensor
	if err := json.NewDecoder(req.Body).Decode(&sensor); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.validator.Validate(sensor); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.sensorRepo.Save(sensor); err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	presenter.Present(sensor)
	controller.reporter.Restart()
}

func (controller SensorController) GetSensorByName(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	sensor, err := controller.sensorRepo.GetByName(sensorName)
	if err != nil {
		presenter.PresentError(domain.SensorNotFoundErr)
		return
	}
	presenter.Present(sensor)
}

func (controller SensorController) DeleteSensorByName(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	sensor, err := controller.sensorRepo.GetByName(sensorName)
	if err != nil {
		presenter.PresentError(domain.SensorNotFoundErr)
		return
	}
	sensor.Deleted = true
	if controller.sensorRepo.Update(sensor) {
		presenter.Present(struct{ Deleted bool }{true})
		return
	}
	presenter.PresentError(domain.InternalErr)
	controller.reporter.Restart()
}

func (controller SensorController) UpdateSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var sensor domain.Sensor
	if err := json.NewDecoder(req.Body).Decode(&sensor); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.validator.Validate(sensor); err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if !controller.sensorRepo.Update(sensor) {
		presenter.PresentError(domain.UpdateErr)
		return
	}
	presenter.Present(sensor)
	controller.reporter.Restart()
}

func (controller SensorController) SensorNow(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	sensorName := mux.Vars(req)["name"]
	sensor, err := controller.sensorRepo.GetByName(sensorName)
	if err != nil {
		presenter.PresentError(domain.SensorNotFoundErr)
		return
	}
	reports, err := sensor.GetCurrentState()
	if err != nil {
		presenter.PresentError(err)
		return
	}
	presenter.Present(reports)
}

func (controller SensorController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetAllSensors, Endpoint: "/"},
		{Method: locomotive.Post, Handler: controller.SaveSensor, Endpoint: "/"},
		{Method: locomotive.Get, Handler: controller.GetSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Delete, Handler: controller.DeleteSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Post, Handler: controller.UpdateSensor, Endpoint: "/update"},
		{Method: locomotive.Get, Handler: controller.SensorNow, Endpoint: "/{name}/now"},
	}
}
