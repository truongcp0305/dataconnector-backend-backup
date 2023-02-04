package main

import (
	"data-connector/library"
	"data-connector/log"
	"data-connector/model"
	"data-connector/router"
	"data-connector/service"
)

func main() {
	log.Info("Start server...", map[string]interface{}{
		"scope": log.Trace(),
	})
	err := library.InitEnvironment()
	if err != nil {
		log.Error("Missing env file config", map[string]interface{}{
			"scope": log.Trace(),
		})
	}
	model.ResetStatus()
	service.InitJob()
	router.InitRoutes()

}
