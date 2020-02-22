package persistence

import (
	"log"
	"sensorapi/src/domain"

	"github.com/jmoiron/sqlx"
)

type SqlxReportTypeRepo struct {
	db *sqlx.DB
}

func NewSqlxReportTypeRepo(conn SqlxConnection) domain.ReportTypeRepo {
	return SqlxReportTypeRepo{conn.GetConnection()}
}

func (repo SqlxReportTypeRepo) Save(rType domain.ReportType) error {
	insert := "insert into REPORT_TYPES values($1)"
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec(insert, string(rType))
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
