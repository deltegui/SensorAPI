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
		phoenix.NewJSONRenderer(w).Render(reportTypeRepo.GetAll())
	}
}

func SaveReportType(reportTypeRepo domain.ReportTypeRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		reportType := domain.ReportType(mux.Vars(req)["name"])
		if err := reportTypeRepo.Save(reportType); err != nil {
			log.Println(err)
			renderer.RenderError(domain.MalformedRequestErr)
			return
		}
		renderer.Render(struct{ ReportType domain.ReportType }{reportType})
	}
}

func registerReportTypesRoutes(app phoenix.App) {
	app.MapGroup("/reporttypes", func(m phoenix.Mapper) {
		m.Get("/all", GetReportTypes)
		m.Post("/create/{name}", SaveReportType)
	})
}
