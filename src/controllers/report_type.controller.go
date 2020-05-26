package controllers

import (
	"log"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/mux"
)

type ReportTypeController struct {
	reportTypeRepo domain.ReportTypeRepo
}

func NewReportTypeController(reportTypeRepo domain.ReportTypeRepo) ReportTypeController {
	return ReportTypeController{reportTypeRepo}
}

func (controller ReportTypeController) GetReportTypes(w http.ResponseWriter, req *http.Request) {
	renderer := phoenix.JSONRenderer{w}
	renderer.Render(controller.reportTypeRepo.GetAll())
}

func (controller ReportTypeController) SaveReportType(w http.ResponseWriter, req *http.Request) {
	renderer := phoenix.JSONRenderer{w}
	reportType := domain.ReportType(mux.Vars(req)["name"])
	if err := controller.reportTypeRepo.Save(reportType); err != nil {
		log.Println(err)
		renderer.RenderError(domain.MalformedRequestErr)
		return
	}
	renderer.Render(struct{ ReportType domain.ReportType }{reportType})
}

func (controller ReportTypeController) GetMappings() []phoenix.CMapping {
	return []phoenix.CMapping{
		{Method: phoenix.Get, Handler: controller.GetReportTypes, Endpoint: "/all"},
		{Method: phoenix.Post, Handler: controller.SaveReportType, Endpoint: "/create/{name}"},
	}
}
