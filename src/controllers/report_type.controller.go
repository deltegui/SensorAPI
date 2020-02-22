package controllers

import (
	"log"
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/locomotive"
	"github.com/gorilla/mux"
)

type ReportTypeController struct {
	reportTypeRepo domain.ReportTypeRepo
}

func NewReportTypeController(reportTypeRepo domain.ReportTypeRepo) ReportTypeController {
	return ReportTypeController{reportTypeRepo}
}

func (controller ReportTypeController) GetReportTypes(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	presenter.Present(controller.reportTypeRepo.GetAll())
}

func (controller ReportTypeController) SaveReportType(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	reportType := domain.ReportType(mux.Vars(req)["name"])
	if err := controller.reportTypeRepo.Save(reportType); err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	presenter.Present(struct{ ReportType domain.ReportType }{reportType})
}

func (controller ReportTypeController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetReportTypes, Endpoint: "/all"},
		{Method: locomotive.Post, Handler: controller.SaveReportType, Endpoint: "/create/{name}"},
	}
}
