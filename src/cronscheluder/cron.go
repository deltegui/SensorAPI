package cronscheluder

import (
	"fmt"
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
		return "*"
	}
	return fmt.Sprintf("*/%d", interval)
}

func (scheluder CronScheluder) AddJobEvery(job domain.ScheluderJob, interval int64) {
	minutes := parseInterval(interval % 60)
	hours := parseInterval((interval / 60) % 24)
	cronExpr := fmt.Sprintf("%s %s * * *", minutes, hours)
	scheluder.cron.AddFunc(cronExpr, job)
}

func (scheluder CronScheluder) Run() {
	scheluder.cron.Run()
}
