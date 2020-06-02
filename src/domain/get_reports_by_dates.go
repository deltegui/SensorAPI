package domain

import (
	"time"
)

type ReportsByDatesRequest struct {
	From time.Time `validate:"required" json:"from"`
	To   time.Time `validate:"required" json:"to"`
}

type GetReportsByDates UseCase

func NewGetReportsByDates(reportRepo ReportRepo, validator Validator) GetReportsByDates {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		datesReq := req.(ReportsByDatesRequest)
		if err := validator.Validate(datesReq); err != nil {
			return nil, MalformedRequestErr
		}
		reports := reportRepo.GetBetweenDates(datesReq.From, datesReq.To)
		return reports, nil
	}
}
