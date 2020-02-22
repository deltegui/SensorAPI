package domain

type Presenter interface {
	Present(data interface{})
	PresentError(data error)
}

type UseCaseRequest interface{}

var EmptyRequest UseCaseRequest = struct{}{}

type UseCase interface {
	Exec(Presenter, UseCaseRequest)
}

type ReportType string

type ReportTypeRepo interface {
	Save(ReportType) error
	GetAll() []ReportType
}
