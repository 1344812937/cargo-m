package service

import (
	"cargo-m/internal/repository"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

type MavenService struct {
	mavenRepo *repository.MavenRepo
}

func NewMavenService(mavenRepo *repository.MavenRepo) *MavenService {
	return &MavenService{mavenRepo: mavenRepo}
}

type MavenGAV struct {
	GroupId    string `json:"groupId" binding:"required"`
	ArtifactId string `json:"artifactId" binding:"required"`
	Version    string `json:"version" binding:"required"`
	FileName   string `json:"filename" binding:"required"`
	Classifier string `json:"classifier,omitempty"` // 可选
	Extension  string `json:"extension,omitempty"`  // 可选
}

func (t *MavenService) GetRepo(c *gin.Context) {
	mavenGavInfo := parseMavenPath(c.Param("path"))
	c.JSON(200, gin.H{"data": mavenGavInfo})
}

func parseMavenPath(path string) *MavenGAV {
	res := &MavenGAV{}
	// 标准化路径并分割
	normalized := filepath.Clean(path)
	index := strings.Index(normalized, "\\")
	if index == 0 {
		normalized = normalized[1:]
	}
	parts := strings.Split(normalized, "\\")

	// 1. 检查路径深度（至少需要4部分：g/a/v/filename）
	if len(parts) < 4 {
		return nil
	}
	// 2. 提取文件名并移除
	res.FileName = parts[len(parts)-1]
	parts = parts[:len(parts)-1]
	// 3. 提取版本号（路径最后一段）
	res.Version = parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	// 4. 提取artifactId（路径最后一段）
	res.ArtifactId = parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	// 5. 剩余部分为groupId（转换为点分隔）
	res.GroupId = strings.Join(parts, ".")

	// 6. 解析文件名
	baseName := strings.TrimSuffix(res.FileName, filepath.Ext(res.FileName))
	expectedPrefix := res.ArtifactId + "-" + res.Version

	// 处理带分类器的文件名
	if strings.HasPrefix(baseName, expectedPrefix) {
		// 示例：myapp-1.0-sources → classifier = "sources"
		suffix := strings.TrimPrefix(baseName, expectedPrefix)
		if suffix != "" {
			// 分类器以 '-' 开头
			if strings.HasPrefix(suffix, "-") {
				res.Classifier = strings.TrimPrefix(suffix, "-")
			}
		}
	}

	// 7. 获取扩展名
	res.Extension = strings.TrimPrefix(filepath.Ext(res.FileName), ".")

	return res
}
