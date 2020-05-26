package controllers

import (
	"encoding/json"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/mux"
)

func registerSensorRoutes() {
	phoenix.MapGroup("/sensor", func(mapper phoenix.Mapper) {
		mapper.MapAll([]phoenix.Mapping{
			{Method: phoenix.Post, Builder: SaveSensor, Endpoint: ""},
			{Method: phoenix.Get, Builder: GetSensorByName, Endpoint: "/{name}"},
			{Method: phoenix.Delete, Builder: DeleteSensorByName, Endpoint: "/{name}"},
			{Method: phoenix.Post, Builder: UpdateSensor, Endpoint: "/update"},
			{Method: phoenix.Get, Builder: SensorNow, Endpoint: "/{name}/now"},
		})
	})
}

func SaveSensor(saveSensorCase domain.SaveSensorCase, builder domain.SensorBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		var viewModel domain.SensorViewModel
		if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
			presenter.PresentError(domain.MalformedRequestErr)
			return
		}
		viewModel.SensorBuilder = builder
		saveSensorCase.Exec(presenter, viewModel)
	}
}

func GetSensorByName(getSensorCase domain.GetSensorCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		sensorName := mux.Vars(req)["name"]
		getSensorCase.Exec(presenter, sensorName)
	}
}

func DeleteSensorByName(deleteSensorCase domain.DeleteSensorCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		sensorName := mux.Vars(req)["name"]
		deleteSensorCase.Exec(presenter, sensorName)
	}
}

func UpdateSensor(updateSensorCase domain.UpdateSensorCase, builder domain.SensorBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		var viewModel domain.SensorViewModel
		if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
			presenter.PresentError(domain.MalformedRequestErr)
			return
		}
		viewModel.SensorBuilder = builder
		updateSensorCase.Exec(presenter, viewModel)
	}
}

func SensorNow(sensorNowCase domain.SensorNowCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		sensorName := mux.Vars(req)["name"]
		sensorNowCase.Exec(presenter, sensorName)
	}
}
