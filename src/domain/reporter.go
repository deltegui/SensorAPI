package domain

import "log"

func newScheluderJob(sensor Sensor, reportRepo ReportRepo) ScheluderJob {
	return func() {
		log.Println("Running job for %s", sensor.Name)
		currentReports := sensor.GetCurrentState()
		for _, report := range currentReports {
			reportRepo.Save(report)
		}
	}
}

type Reporter struct {
	sensorRepo SensorRepo
	reportRepo ReportRepo
	scheluder  ReportScheluder
}

func NewReporter(sensorRepo SensorRepo, reportRepo ReportRepo, scheluder ReportScheluder) Reporter {
	return Reporter{
		sensorRepo,
		reportRepo,
		scheluder,
	}
}

func (reporter Reporter) Start() {
	sensors := reporter.sensorRepo.GetAll()
	for _, sensor := range sensors {
		job := newScheluderJob(sensor, reporter.reportRepo)
		reporter.scheluder.AddJobEvery(job, sensor.UpdateInterval)
	}
	reporter.scheluder.Run()
}
