package domain

import "time"

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
	Connector        SensorConnector `json:"-"`
}

func (sensor Sensor) GetCurrentState() []Report {
	return sensor.Connector.ReadDataFor(sensor)
}

type SensorRepo interface {
	Save(Sensor) error
	GetAll() []Sensor
	GetByName(name string) (Sensor, error)
	Update(sensor Sensor) bool
}

// Report represents a Sensor report.
type Report struct {
	ReportType string
	SensorName string
	Date       time.Time
	Value      float32
}

// SensorConnector is an abstraction over
// the way to communicate to real sensor
type SensorConnector interface {
	ReadDataFor(Sensor) []Report
}

type ReportRepo interface {
	Save(Report)
	GetAll() []Report
}

type ScheluderJob func()

type ReportScheluder interface {
	AddJobEvery(ScheluderJob, int64)
	Run()
}
