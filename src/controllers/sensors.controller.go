package controllers

import (
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/phoenix"
)

func GetAllSensors(getAllSensorsCase domain.GetAllSensorsCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		wantDeleted := req.URL.Query()["deleted"]
		var reqCase domain.GetAllRequest
		reqCase.WantDeleted = !(len(wantDeleted) < 1 || len(wantDeleted[0]) == 0 || wantDeleted[0] == "false")
		getAllSensorsCase.Exec(presenter, reqCase)
	}
}

func AllSensorNow(allSensorNowCase domain.AllSensorNowCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := NewJSONPresenter(w)
		allSensorNowCase.Exec(presenter, nil)
	}
}

func registerSensorsRoutes() {
	phoenix.MapGroup("/sensors", func(m phoenix.Mapper) {
		m.MapAll([]phoenix.Mapping{
			{Method: phoenix.Get, Builder: GetAllSensors, Endpoint: ""},
			{Method: phoenix.Get, Builder: AllSensorNow, Endpoint: "/all/now"},
		})
	})
}
