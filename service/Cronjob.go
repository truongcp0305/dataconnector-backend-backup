package service

import (
	"data-connector/library"
	"data-connector/log"
	"data-connector/model"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Config string `json:"config"`
}

type JobInterface interface {
	InitJob() *cron.Cron
}

var c *cron.Cron

func InitJob() {
	if c != nil {
		c.Stop()
		c = nil
	}
	job := new(Job)
	job.StartJob()
}

func (job Job) StartJob() *cron.Cron {
	c = cron.New()
	cronjob := new(model.Cronjob)
	values, err := cronjob.GetListCronjob("1")
	if err == nil && len(values) > 0 {
		for i := 0; i < len(values); i++ {
			config := values[i]["config"]
			apiUuid := values[i]["api_query_uuid"]
			if err == nil {
				c.AddFunc(config, func() {
					req := new(library.Request)
					log.Info("job "+config+" with api uuid "+apiUuid+" callback", map[string]interface{}{
						"scope": log.Trace(),
					})
					req.Url = "https://data-connector.symper.vn/apiQueries/loadData"
					body := map[string]string{
						"uuid": apiUuid,
					}
					h := library.HEADER
					h["Content-Type"] = "application/x-www-form-urlencoded"
					req.Header = h
					req.Body = body
					req.Method = "POST"
					go req.Send()
				})
			}
		}
	}
	c.Start()
	return c
}
