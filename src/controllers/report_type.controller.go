package controllers

import (
	"sensorapi/src/domain"

	"net/http"

	"github.com/deltegui/locomotive"
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

func (controller ReportTypeController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetReportTypes, Endpoint: "/all"},
	}
}
