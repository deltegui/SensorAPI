package domain

import (
	"log"
)

type ConnetionViewModel struct {
	ConnType  ConnectionType `validate:"required" json:"type"`
	ConnValue string         `validate:"required" json:"value"`
}

type SensorViewModel struct {
	Name             string             `validate:"required" json:"name"`
	Connection       ConnetionViewModel `validate:"required" json:"connection"`
	UpdateInterval   int64              `validate:"required" json:"updateInterval"`
	Deleted          bool               `json:"deleted"`
	SupportedReports []ReportType       `validate:"required" json:"supportedReports"`
	SensorBuilder    SensorBuilder      `json:"-"`
}

func (model SensorViewModel) ToSensor() Sensor {
	return model.SensorBuilder.
		WithName(model.Name).
		WithConnection(model.Connection.ConnType, model.Connection.ConnValue).
		WithSupportedReports(model.SupportedReports).
		WithUpdateInterval(model.UpdateInterval).
		IsDeleted(model.Deleted).
		Build()
}

func CreateViewModelFromSensor(sensor Sensor) SensorViewModel {
	return SensorViewModel{
		Name: sensor.Name,
		Connection: ConnetionViewModel{
			ConnType:  sensor.ConnType,
			ConnValue: sensor.ConnValue,
		},
		Deleted:          sensor.Deleted,
		UpdateInterval:   sensor.UpdateInterval,
		SupportedReports: sensor.SupportedReports,
	}
}

type GetAllRequest struct {
	WantDeleted bool
}

type GetAllSensorsCase UseCase

func NewGetAllSensorsCase(sensorRepo SensorRepo) GetAllSensorsCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		getAllReq := req.(GetAllRequest)
		var sensors []Sensor
		var err error
		if getAllReq.WantDeleted {
			sensors, err = sensorRepo.GetAll(WithDeletedSensors)
		} else {
			sensors, err = sensorRepo.GetAll(WithoutDeletedSensors)
		}
		if err != nil {
			return nil, SensorNotFoundErr
		}
		viewModels := []SensorViewModel{}
		for _, sensor := range sensors {
			viewModels = append(viewModels, CreateViewModelFromSensor(sensor))
		}
		return viewModels, nil
	}
}

type SaveSensorCase UseCase

func NewSaveSensorCase(sensorRepo SensorRepo, validator Validator, reporter Reporter, reportTypeRepo ReportTypeRepo) SaveSensorCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		viewModel := req.(SensorViewModel)
		if err := validator.Validate(viewModel); err != nil {
			return nil, MalformedRequestErr
		}
		if !haveReportTypes(reportTypeRepo, viewModel.SupportedReports) {
			return nil, ReportTypeDoesNotExists
		}
		if err := sensorRepo.Save(viewModel.ToSensor()); err != nil {
			log.Println(err)
			return nil, SensorAlreadyExist
		}
		reporter.Restart()
		return viewModel, nil
	}
}

func haveReportTypes(reportTypeRepo ReportTypeRepo, reportTypes []ReportType) bool {
	foundRTypes := reportTypeRepo.GetAll()
	for _, rType := range reportTypes {
		found := false
		for _, haveRType := range foundRTypes {
			if rType == haveRType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type DeleteSensorCase UseCase

func NewDeleteSensorCase(sensorRepo SensorRepo, validator Validator, reporter Reporter) DeleteSensorCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		sensorName := req.(string)
		sensor, err := sensorRepo.GetByName(sensorName)
		if err != nil {
			return nil, SensorNotFoundErr
		}
		sensor.Deleted = true
		if sensorRepo.Update(sensor) {
			return struct{ Deleted bool }{true}, nil
		}
		reporter.Restart()
		return nil, InternalErr
	}
}

type UpdateSensorCase UseCase

func NewUpdateSensorCase(sensorRepo SensorRepo, validator Validator, reporter Reporter, reportTypeRepo ReportTypeRepo) UpdateSensorCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		viewModel := req.(SensorViewModel)
		if err := validator.Validate(viewModel); err != nil {
			log.Println(err)
			return nil, MalformedRequestErr
		}
		if !haveReportTypes(reportTypeRepo, viewModel.SupportedReports) {
			return nil, ReportTypeDoesNotExists
		}
		if !sensorRepo.Update(viewModel.ToSensor()) {
			return nil, UpdateErr
		}
		reporter.Restart()
		return viewModel, nil
	}
}

type SensorNowCase UseCase

func NewSensorNowCase(sensorRepo SensorRepo, validator Validator, reporter Reporter) SensorNowCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		sensorName := req.(string)
		sensor, err := sensorRepo.GetByName(sensorName)
		if err != nil {
			return nil, SensorNotFoundErr
		}
		reports, err := sensor.GetCurrentState()
		if err != nil {
			return nil, err
		}
		return reports, nil
	}
}

type AllSensorNowCase UseCase

func NewAllSensorNowCase(sensorRepo SensorRepo) AllSensorNowCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		sensors, err := sensorRepo.GetAll(WithoutDeletedSensors)
		if err != nil {
			return nil, SensorNotFoundErr
		}
		var reports []Report
		for _, sensor := range sensors {
			r, err := sensor.GetCurrentState()
			if err == nil {
				reports = append(reports, r...)
			}
		}
		return reports, nil
	}
}

type GetSensorCase UseCase

func NewGetSensorCase(sensorRepo SensorRepo) GetSensorCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		sensorName := req.(string)
		sensor, err := sensorRepo.GetByName(sensorName)
		if err != nil {
			return nil, SensorNotFoundErr
		}
		return CreateViewModelFromSensor(sensor), nil
	}
}
