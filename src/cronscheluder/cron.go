package cronscheluder

import (
	"fmt"
	"log"
	"sensorapi/src/domain"

	"github.com/robfig/cron"
)

type CronScheluder struct {
	cron *cron.Cron
}

func NewCronScheluder() domain.ReportScheluder {
	return CronScheluder{cron.New()}
}

func parseInterval(interval int64) string {
	if interval == 0 {
		return "0"
	}
	return fmt.Sprintf("*/%d", interval)
}

func (scheluder CronScheluder) AddJobEvery(job domain.ScheluderJob, interval int64) {
	minutes := parseInterval(interval % 60)
	cronExpr := fmt.Sprintf("%s * * * *", minutes)
	log.Printf("One job have %d inerval, traduced to cron expression: %s\n", interval, cronExpr)
	scheluder.cron.AddFunc(cronExpr, job)
}

func (scheluder CronScheluder) Run() {
	scheluder.cron.Run()
}
