package controllers

import (
	"log"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/mux"
)

func GetReportTypes(reportTypeRepo domain.ReportTypeRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.JSONRenderer{w}
		renderer.Render(reportTypeRepo.GetAll())
	}
}

func SaveReportType(reportTypeRepo domain.ReportTypeRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.JSONRenderer{w}
		reportType := domain.ReportType(mux.Vars(req)["name"])
		if err := reportTypeRepo.Save(reportType); err != nil {
			log.Println(err)
			renderer.RenderError(domain.MalformedRequestErr)
			return
		}
		renderer.Render(struct{ ReportType domain.ReportType }{reportType})
	}
}

func registerReportTypesRoutes() {
	phoenix.MapGroup("/reporttypes", func(m phoenix.Mapper) {
		m.MapAll([]phoenix.Mapping{
			{Method: phoenix.Get, Builder: GetReportTypes, Endpoint: "/all"},
			{Method: phoenix.Post, Builder: SaveReportType, Endpoint: "/create/{name}"},
		})
	})
}
