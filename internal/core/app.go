package core

import (
	"cargo-m/internal/config"
	"cargo-m/internal/tasks"
	"cargo-m/internal/until"
	"fmt"

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
		err := a.WebEngine.Run(fmt.Sprintf("%s:%s", "", a.ApplicationConfig.WebConfig.Port))
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
