package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sensorapi/src/domain"
	"time"

	"github.com/deltegui/locomotive"
)

type ReportController struct {
	reportRepo       domain.ReportRepo
	getReportsByDate domain.UseCase
}

func NewReportController(reportRepo domain.ReportRepo, getReportsByDate domain.GetReportsByDates) ReportController {
	return ReportController{reportRepo, getReportsByDate}
}

func (controller ReportController) GetAllReports(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	presenter.Present(controller.reportRepo.GetAll())
}

func getQueryFrom(req *http.Request, query string) (string, error) {
	key, ok := req.URL.Query()[query]
	if !ok || len(key[0]) < 1 {
		return "", fmt.Errorf("Key does not exist")
	}
	return key[0], nil
}

func getDateFrom(req *http.Request, query string) (time.Time, error) {
	dateStr, err := getQueryFrom(req, query)
	if err != nil {
		return time.Now(), err
	}
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}

func (controller ReportController) GetReportsBetweenDates(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	from, err := getDateFrom(req, "from")
	if err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	to, err := getDateFrom(req, "to")
	if err != nil {
		log.Println(err)
		presenter.PresentError(domain.MalformedRequestErr)
		return
	}
	log.Println(from, to)
	controller.getReportsByDate.Exec(presenter, domain.ReportsByDatesRequest{
		From: from,
		To:   to,
	})
}

func (controller ReportController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: controller.GetAllReports, Endpoint: "/"},
		{Method: locomotive.Get, Handler: controller.GetReportsBetweenDates, Endpoint: "/dates"},
	}
}
