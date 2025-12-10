package core

import (
	"cargo-m/internal/config"
	"cargo-m/internal/proxy"
	"cargo-m/internal/tasks"
	"cargo-m/internal/until"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var ApplicationInstance *Application

type Application struct {
	ApplicationConfig *config.ApplicationConfig
	CornTask          *tasks.CronTask
	WebEngine         *gin.Engine
	ProxyServer       *proxy.SocksProxy
	StartTime         *time.Time
}

func NewApplication(applicationConfig *config.ApplicationConfig, cornTask *tasks.CronTask, webEngine *gin.Engine, proxyServer *proxy.SocksProxy) *Application {
	if ApplicationInstance == nil {
		ApplicationInstance = &Application{ApplicationConfig: applicationConfig, CornTask: cornTask, WebEngine: webEngine, ProxyServer: proxyServer}
	}
	return ApplicationInstance
}

func (app *Application) Start() {
	app.StartTime = &time.Time{}
	if app.CornTask != nil {
		app.CornTask.Start()
	}
	proxyConfig := app.ApplicationConfig.ProxyConfig
	if proxyConfig.Enabled {
		go app.ProxyServer.Run(proxyConfig.Port, proxyConfig.AuthUser, proxyConfig.AuthPass)
	}
	if app.WebEngine != nil {
		err := app.WebEngine.Run(fmt.Sprintf("%s:%s", app.ApplicationConfig.WebConfig.Host, app.ApplicationConfig.WebConfig.Port))
		if err != nil {
			panic(err)
		}
	}
	until.Log.Info("Application started")
}

func (app *Application) Close() {
	if app.CornTask != nil {
		app.CornTask.Stop()
	}
}
