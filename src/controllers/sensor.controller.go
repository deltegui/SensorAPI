package controllers

import (
	"encoding/json"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/mux"
)

func registerSensorRoutes(app phoenix.App) {
	app.MapGroup("/sensor", func(m phoenix.Mapper) {
		m.Post("", SaveSensor)
		m.Get("/{name}", GetSensorByName)
		m.Delete("/{name}", DeleteSensorByName)
		m.Post("/update", UpdateSensor)
		m.Get("/{name}/now", SensorNow)
	})
}

func SaveSensor(saveSensor domain.SaveSensorCase, builder domain.SensorBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		var viewModel domain.SensorViewModel
		if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
			renderer.RenderError(domain.MalformedRequestErr)
			return
		}
		viewModel.SensorBuilder = builder
		response, err := saveSensor(viewModel)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(response)
	}
}

func GetSensorByName(getSensorBy domain.GetSensorCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		name := mux.Vars(req)["name"]
		sensor, err := getSensorBy(name)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(sensor)
	}
}

func DeleteSensorByName(deleteSensorBy domain.DeleteSensorCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		name := mux.Vars(req)["name"]
		sensor, err := deleteSensorBy(name)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(sensor)
	}
}

func UpdateSensor(updateSensorCase domain.UpdateSensorCase, builder domain.SensorBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		var viewModel domain.SensorViewModel
		if err := json.NewDecoder(req.Body).Decode(&viewModel); err != nil {
			renderer.RenderError(domain.MalformedRequestErr)
			return
		}
		viewModel.SensorBuilder = builder
		sensor, err := updateSensorCase(viewModel)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(sensor)
	}
}

func SensorNow(sensorNowCase domain.SensorNowCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		sensorName := mux.Vars(req)["name"]
		response, err := sensorNowCase(sensorName)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(response)
	}
}
