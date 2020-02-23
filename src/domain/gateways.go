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
	Name             string         `db:"NAME"`
	ConnType         ConnectionType `db:"CONNTYPE"`
	ConnValue        string         `db:"CONNVALUE"`
	UpdateInterval   int64          `db:"UPDATE_INTERVAL"`
	Deleted          bool           `db:"DELETED"`
	SupportedReports []ReportType
	Connector        SensorConnector `json:"-"`
}

func (sensor Sensor) GetCurrentState() ([]Report, error) {
	reports, err := sensor.Connector.ReadDataFor(sensor)
	if err != nil {
		return reports, SensorNotRespondErr
	}
	return reports, nil
}

type ShowDeleted bool

const (
	WithDeletedSensors    ShowDeleted = true
	WithoutDeletedSensors ShowDeleted = false
)

type SensorRepo interface {
	Save(Sensor) error
	GetAll(showDeleted ShowDeleted) []Sensor
	GetByName(name string) (Sensor, error)
	Update(sensor Sensor) bool
}

// Report represents a Sensor report.
type Report struct {
	ReportType string    `db:"TYPE"`
	SensorName string    `db:"SENSOR"`
	Date       time.Time `db:"REPORT_DATE"`
	Value      float32   `db:"VALUE"`
}

// SensorConnector is an abstraction over
// the way to communicate to real sensor
type SensorConnector interface {
	ReadDataFor(Sensor) ([]Report, error)
}

type ReportRepo interface {
	Save(Report)
	GetAll() []Report
}

type ScheluderJob func()

type ReportScheluder interface {
	AddJobEvery(ScheluderJob, int64)
	Start() // DEBE SER NO BLOQUEANTE (QUE ARRANQUE UNA CORRUTINA)
	Stop()
}
