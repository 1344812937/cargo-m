package api

import (
	"cargo-m/internal/core"
	"cargo-m/internal/service"
	"cargo-m/internal/until"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

type BlueCatApi struct {
	mavenService *service.MavenService
}

func NewBlueCat(mavenService *service.MavenService) *BlueCatApi {
	return &BlueCatApi{mavenService: mavenService}
}

func (bc *BlueCatApi) GetRunnerInfo(c *gin.Context) {
	app := core.GetCurrentApp()
	webConfig := app.ApplicationConfig.WebConfig
	var serverHost []string

	if len(webConfig.Host) == 0 {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			until.Log.Error("Can't get interface addresses", err)
		}
		for _, addr := range addrs {
			// 检查ip地址类型，跳过回环地址
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					serverHost = append(serverHost, fmt.Sprintf("http://%s:%s/maven-repo/getRepo/", ipnet.IP, webConfig.Port))
				}
			}
		}
	} else {
		serverHost = append(serverHost, fmt.Sprintf("http://%s:%s/maven-repo/getRepo/", webConfig.Host, webConfig.Port))
	}
	res := gin.H{
		"mavenRepoUrl": serverHost,
		"mavenRepo":    app.ApplicationConfig.LocalRepoConfig,
		"startTime":    app.StartTime.Format("2006-01-02 15:04:05.000"),
	}
	c.JSON(200, res)
}

func (bc *BlueCatApi) Register(group *gin.RouterGroup) {
	group.GET("/", bc.GetRunnerInfo)
}
