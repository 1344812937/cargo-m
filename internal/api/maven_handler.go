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

func (handler *MavenRepoHandler) Register(router *gin.RouterGroup) {
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"gav": "hh"})
	})
}
