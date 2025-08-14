package core

import (
	"cargo-m/internal/config"
	"cargo-m/internal/tasks"
	"cargo-m/internal/until"

	"github.com/gin-gonic/gin"
)

type Application struct {
	ApplicationConfig *config.ApplicationConfig
	CornTask          *tasks.CronTask
	WebEngine         *gin.Engine
}

func NewApplication(applicationConfig *config.ApplicationConfig, cornTask *tasks.CronTask, webEngine *gin.Engine) *Application {
	return &Application{ApplicationConfig: applicationConfig, CornTask: cornTask, WebEngine: webEngine}
}

func (a *Application) Start() {
	if a.CornTask != nil {
		a.CornTask.Start()
	}
	if a.WebEngine != nil {
		err := a.WebEngine.Run(":9080")
		if err != nil {
			panic(err)
		}
	}
	until.Log.Info("Application started")
}

func (a *Application) Close() {
	if a.CornTask != nil {
		a.CornTask.Stop()
	}
}
