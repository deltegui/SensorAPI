package controllers

import (
	"net/http"
	"sensorapi/src/domain"

	"github.com/deltegui/locomotive"
)

type ReportController struct {
	reportRepo domain.ReportRepo
}

func NewReportController(reportRepo domain.ReportRepo) ReportController {
	return ReportController{reportRepo}
}

func (controller ReportController) GetAllReports(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	presenter.Present(controller.reportRepo.GetAll())
}

func (controller ReportController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetAllReports, Endpoint: "/"},
	}
}
