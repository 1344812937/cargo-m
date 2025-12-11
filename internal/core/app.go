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

var applicationInstance *Application

func GetCurrentApp() *Application {
	return applicationInstance
}

type Application struct {
	ApplicationConfig *config.ApplicationConfig
	CornTask          *tasks.CronTask
	WebEngine         *gin.Engine
	ProxyServer       *proxy.SocksProxy
	StartTime         *time.Time
}

func NewApplication(applicationConfig *config.ApplicationConfig, cornTask *tasks.CronTask, webEngine *gin.Engine, proxyServer *proxy.SocksProxy) *Application {
	if applicationInstance == nil {
		applicationInstance = &Application{ApplicationConfig: applicationConfig, CornTask: cornTask, WebEngine: webEngine, ProxyServer: proxyServer}
	}
	return applicationInstance
}

func (app *Application) Start() {
	now := time.Now()
	app.StartTime = &now
	if app.CornTask != nil {
		app.CornTask.Start()
	}
	proxyConfig := app.ApplicationConfig.ProxyConfig
	if proxyConfig.Enabled {
		go app.ProxyServer.Run(proxyConfig.Port, proxyConfig.AuthUser, proxyConfig.AuthPass)
	}
	if app.WebEngine != nil {
		go func() {
			err := app.WebEngine.Run(fmt.Sprintf("%s:%s", app.ApplicationConfig.WebConfig.Host, app.ApplicationConfig.WebConfig.Port))
			if err != nil {
				panic(err)
			}
		}()
	}
	until.Log.Info("Application started")
}

func (app *Application) Close() {
	if app.CornTask != nil {
		app.CornTask.Stop()
	}
}
