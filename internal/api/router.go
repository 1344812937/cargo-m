package api

import "github.com/gin-gonic/gin"

// NewRouter 路由初始化
func NewRouter(mavenRepoHandler *MavenRepoHandler) *gin.Engine {
	webEngine := gin.Default()
	mavenRepoHandler.Register(webEngine.Group("maven-repo"))
	return webEngine
}
