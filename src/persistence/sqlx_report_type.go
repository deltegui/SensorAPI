package persistence

import (
	"log"
	"sensorapi/src/domain"

	"github.com/jmoiron/sqlx"
)

type SqlxReportTypeRepo struct {
	db *sqlx.DB
}

func NewSqlxReportTypeRepo(conn *SqlxConnection) domain.ReportTypeRepo {
	return SqlxReportTypeRepo{conn.GetConnection()}
}

func (repo SqlxReportTypeRepo) Save(rType domain.ReportType) error {
	insert := "insert into REPORT_TYPES values(?)"
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec(insert, rType)
	return err
}

func (repo SqlxReportTypeRepo) GetAll() []domain.ReportType {
	tx, err := repo.db.Beginx()
	if err != nil {
		return []domain.ReportType{}
	}
	defer tx.Commit()
	var types []domain.ReportType
	err = tx.Select(&types, "SELECT TYPE_NAME FROM REPORT_TYPES")
	if err != nil {
		log.Println(err)
		return []domain.ReportType{}
	}
	return types
}

type SqlxSensorRepo struct {
	db *sqlx.DB
}

func NewSqlxSensorRepo(conn *SqlxConnection) domain.SensorRepo {
	return SqlxSensorRepo{conn.GetConnection()}
}

func (repo SqlxSensorRepo) Save(sensor domain.Sensor) error {
	insertSensor := "insert into SENSORS (NAME, CONNTYPE, CONNVALUE, UPDATE_INTERVAL)values(?, ?, ? ,?)"
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	if _, err = tx.Exec(insertSensor, sensor.Name, sensor.ConnType, sensor.ConnValue, sensor.UpdateInterval); err != nil {
		return err
	}
	return repo.saveSupportedReportsForSensor(tx, sensor)
}

func (repo SqlxSensorRepo) saveSupportedReportsForSensor(tx *sqlx.Tx, sensor domain.Sensor) error {
	insertReport := "insert into USED_REPORT_TYPES (SENSOR, REPORT_TYPE, ADD_DATE)values(?, ?, NOW())"
	if _, err := tx.Exec("DELETE FROM USED_REPORT_TYPES WHERE SENSOR LIKE ?", sensor.Name); err != nil {
		return err
	}
	for _, reportType := range sensor.SupportedReports {
		if _, err := tx.Exec(insertReport, sensor.Name, reportType); err != nil {
			return err
		}
	}
	return nil
}

func (repo SqlxSensorRepo) GetAll() []domain.Sensor {
	tx, err := repo.db.Beginx()
	if err != nil {
		return []domain.Sensor{}
	}
	defer tx.Commit()
	var sensors []domain.Sensor
	err = tx.Select(&sensors, "SELECT NAME as Name, CONNTYPE as ConnType, CONNVALUE as ConnValue, UPDATE_INTERVAL as UpdateInterval, DELETED as Deleted FROM SENSORS")
	if err != nil {
		log.Println(err)
		return []domain.Sensor{}
	}
	for i := 0; i < len(sensors); i++ {
		repo.FillSupportedReportsForSensor(tx, &sensors[i])
	}
	return sensors
}

func (repo SqlxSensorRepo) FillSupportedReportsForSensor(tx *sqlx.Tx, sensor *domain.Sensor) {
	var reports []domain.ReportType
	err := tx.Select(&reports, "SELECT REPORT_TYPE FROM USED_REPORT_TYPES WHERE SENSOR LIKE ?", sensor.Name)
	if err != nil {
		log.Println(err)
		return
	}
	sensor.SupportedReports = reports
}

func (repo SqlxSensorRepo) GetByName(name string) *domain.Sensor {
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil
	}
	defer tx.Commit()
	var sensor []domain.Sensor
	err = tx.Select(&sensor, "SELECT NAME as Name, CONNTYPE as ConnType, CONNVALUE as ConnValue, UPDATE_INTERVAL as UpdateInterval, DELETED as Deleted FROM SENSORS WHERE NAME LIKE ?", name)
	if err != nil || len(sensor) < 1 {
		log.Println(err)
		return nil
	}
	repo.FillSupportedReportsForSensor(tx, &sensor[0])
	return &sensor[0]
}

func (repo SqlxSensorRepo) Update(sensor domain.Sensor) bool {
	update := "UPDATE SENSORS SET CONNTYPE = ?, CONNVALUE = ?, UPDATE_INTERVAL = ?, DELETED = ? WHERE NAME LIKE ?"
	tx, err := repo.db.Beginx()
	if err != nil {
		log.Println(err)
		return false
	}
	defer tx.Commit()
	if _, err = tx.Exec(update, sensor.ConnType, sensor.ConnValue, sensor.UpdateInterval, sensor.Deleted, sensor.Name); err != nil {
		log.Println(err)
		return false
	}
	if err = repo.saveSupportedReportsForSensor(tx, sensor); err != nil {
		log.Println(err)
		return false
	}
	return true
}
