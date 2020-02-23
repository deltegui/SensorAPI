package controllers

import (
	"encoding/json"
	"log"
	"sensorapi/src/connectors"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/locomotive"
	"github.com/gorilla/mux"
)

type ConnetionViewModel struct {
	ConnType  domain.ConnectionType `validate:"required" json:"type"`
	ConnValue string                `validate:"required" json:"value"`
}

type SensorViewModel struct {
	Name             string              `validate:"required" json:"name"`
	Connection       ConnetionViewModel  `validate:"required" json:"connection"`
	UpdateInterval   int64               `validate:"required" json:"updateInterval"`
	Deleted          bool                `json:"deleted"`
	SupportedReports []domain.ReportType `validate:"required" json:"supportedReports"`
}

func (model SensorViewModel) toSensor() domain.Sensor {
	return domain.Sensor{
		Name:             model.Name,
		ConnType:         model.Connection.ConnType,
		ConnValue:        model.Connection.ConnValue,
		UpdateInterval:   model.UpdateInterval,
		SupportedReports: model.SupportedReports,
		Connector:        connectors.HTTPConnector{IP: model.Connection.ConnValue},
	}
}

func createViewModelFromSensor(sensor domain.Sensor) SensorViewModel {
	return SensorViewModel{
		Name: sensor.Name,
		Connection: ConnetionViewModel{
			ConnType:  sensor.ConnType,
			ConnValue: sensor.ConnValue,
		},
		UpdateInterval:   sensor.UpdateInterval,
		SupportedReports: sensor.SupportedReports,
	}
}

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
	var sensors []domain.Sensor
	if len(wantDeleted) < 1 || len(wantDeleted[0]) == 0 || wantDeleted[0] == "false" {
		sensors = controller.sensorRepo.GetAll(domain.WithoutDeletedSensors)
	} else {
		sensors = controller.sensorRepo.GetAll(domain.WithDeletedSensors)
	}
	var viewModels []SensorViewModel
	for _, sensor := range sensors {
		viewModels = append(viewModels, createViewModelFromSensor(sensor))
	}
	presenter.Present(viewModels)
}

func (controller SensorController) SaveSensor(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	var viewModel SensorViewModel
	if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.validator.Validate(viewModel); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.sensorRepo.Save(viewModel.toSensor()); err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	presenter.Present(viewModel)
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
	presenter.Present(createViewModelFromSensor(sensor))
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
	var viewModel SensorViewModel
	if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if err := controller.validator.Validate(viewModel); err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	if !controller.sensorRepo.Update(viewModel.toSensor()) {
		presenter.PresentError(domain.UpdateErr)
		return
	}
	presenter.Present(viewModel)
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
		{Method: locomotive.Get, Handler: controller.GetAllSensors, Endpoint: ""},
		{Method: locomotive.Post, Handler: controller.SaveSensor, Endpoint: ""},
		{Method: locomotive.Get, Handler: controller.GetSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Delete, Handler: controller.DeleteSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Post, Handler: controller.UpdateSensor, Endpoint: "/update"},
		{Method: locomotive.Get, Handler: controller.SensorNow, Endpoint: "/{name}/now"},
	}
}
