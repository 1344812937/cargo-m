package api

import (
	"cargo-m/internal/service"

	"github.com/gin-gonic/gin"
)

type MavenRepoHandler struct {
	mavenService *service.MavenService
}

// NewMavenRepoHandler 初始化方法
func NewMavenRepoHandler(mavenService *service.MavenService) *MavenRepoHandler {
	return &MavenRepoHandler{mavenService}
}

// Register 路由注册
// 将需要暴露的路由方法注册
func (handler *MavenRepoHandler) Register(router *gin.RouterGroup) {
	router.GET("/getRepo/*path", handler.mavenService.GetRepo)
	router.HEAD("/getRepo/*path", handler.mavenService.CheckRepo)
}
