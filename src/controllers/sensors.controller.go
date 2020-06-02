package controllers

import (
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/phoenix"
)

func GetAllSensors(getAllSensors domain.GetAllSensorsCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		wantDeleted := req.URL.Query()["deleted"]
		var reqCase domain.GetAllRequest
		reqCase.WantDeleted = !(len(wantDeleted) < 1 || len(wantDeleted[0]) == 0 || wantDeleted[0] == "false")
		sensors, err := getAllSensors(reqCase)
		renderer := phoenix.NewJSONRenderer(w)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(sensors)
	}
}

func AllSensorNow(getSensorState domain.AllSensorNowCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reports, err := getSensorState(nil)
		renderer := phoenix.NewJSONRenderer(w)
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(reports)
	}
}

func registerSensorsRoutes(app phoenix.App) {
	app.MapGroup("/sensors", func(m phoenix.Mapper) {
		m.MapRoot(GetAllSensors)
		m.Get("/all/now", AllSensorNow)
	})
}
