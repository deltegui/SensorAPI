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

type ConnectionType string

const (
	HTTP ConnectionType = "http"
)

type Sensor struct {
	Name             string         `db:"Name"`
	ConnType         ConnectionType `db:"ConnType"`
	ConnValue        string         `db:"ConnValue"`
	UpdateInterval   int64          `db:"UpdateInterval"`
	Deleted          bool           `db:"Deleted"`
	SupportedReports []ReportType
}

type SensorRepo interface {
	Save(Sensor) error
	GetAll() []Sensor
	GetByName(name string) *Sensor
	Update(sensor Sensor) bool
}
