package core

import (
	"cargo-m/internal/config"
	"cargo-m/internal/proxy"
	"cargo-m/internal/tasks"
	"cargo-m/internal/until"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Application struct {
	ApplicationConfig *config.ApplicationConfig
	CornTask          *tasks.CronTask
	WebEngine         *gin.Engine
	ProxyServer       *proxy.SocksProxy
}

func NewApplication(applicationConfig *config.ApplicationConfig, cornTask *tasks.CronTask, webEngine *gin.Engine, proxyServer *proxy.SocksProxy) *Application {
	return &Application{ApplicationConfig: applicationConfig, CornTask: cornTask, WebEngine: webEngine, ProxyServer: proxyServer}
}

func (app *Application) Start() {
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
