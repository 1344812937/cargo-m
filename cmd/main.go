package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type MavenGAV struct {
	GroupId    string `json:"groupId" binding:"required"`
	ArtifactId string `json:"artifactId" binding:"required"`
	Version    string `json:"version" binding:"required"`
	FileName   string `json:"filename" binding:"required"`
	Classifier string `json:"classifier,omitempty"` // 可选
	Extension  string `json:"extension,omitempty"`  // 可选
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		path := c.FullPath()
		println(path)
		c.JSONP(200, "a")
	})
	r.GET("/:group/:artifact/:version/*file", func(c *gin.Context) {
		// 提取路径参数
		groupId := strings.Replace(c.Param("group"), "/", ".", -1)
		artifactId := c.Param("artifact")
		version := c.Param("version")
		fileName := c.Param("file")[1:] // 移除开头的 "/"

		gav := MavenGAV{
			GroupId:    groupId,
			ArtifactId: artifactId,
			Version:    version,
			FileName:   fileName,
		}

		c.JSON(200, gin.H{"gav": gav})
	})
	err := r.Run(":9090")
	if err != nil {
		println("运行发生异常", err.Error())
		return
	}
}
