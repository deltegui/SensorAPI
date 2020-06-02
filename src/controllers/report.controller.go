package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sensorapi/src/domain"
	"time"

	"github.com/deltegui/phoenix"
)

func GetAllReports(reportRepo domain.ReportRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.JSONRenderer{w}
		renderer.Render(reportRepo.GetAll())
	}
}

func GetReportsBetweenDates(getReportsByDate domain.GetReportsByDates) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		renderer := phoenix.NewJSONRenderer(w)
		from, err := getDateFrom(req, "from")
		if err != nil {
			log.Println(err)
			renderer.RenderError(domain.MalformedRequestErr)
			return
		}
		to, err := getDateFrom(req, "to")
		if err != nil {
			log.Println(err)
			renderer.RenderError(domain.MalformedRequestErr)
			return
		}
		log.Println(from, to)
		reports, err := getReportsByDate(domain.ReportsByDatesRequest{
			From: from,
			To:   to,
		})
		if err != nil {
			renderer.RenderError(err)
			return
		}
		renderer.Render(reports)
	}
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

func registerReportRoutes(app phoenix.App) {
	app.MapGroup("/routes", func(m phoenix.Mapper) {
		m.Get("/", GetAllReports)
		m.Get("/dates", GetReportsBetweenDates)
	})
}
