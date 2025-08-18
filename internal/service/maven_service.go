package service

import (
	"cargo-m/internal/config"
	"cargo-m/internal/model"
	"cargo-m/internal/repository"
	"cargo-m/internal/until"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type MavenService struct {
	applicationConfig *config.ApplicationConfig
	mavenRepo         *repository.MavenRepo
}

func NewMavenService(mavenRepo *repository.MavenRepo, applicationConfig *config.ApplicationConfig) *MavenService {
	return &MavenService{mavenRepo: mavenRepo, applicationConfig: applicationConfig}
}

func (t *MavenService) GetRepo(c *gin.Context) {
	mavenGavInfo := parseMavenPath(c.Param("path"))
	if mavenGavInfo != nil {
		key, err := t.mavenRepo.GetByKey(mavenGavInfo.Key)
		if err != nil {
			c.JSON(501, gin.H{})
			return
		}
		if key == nil {
			c.JSON(404, gin.H{})
			return
		}
		file, err := os.Open(key.FilePath)
		if err != nil {
			until.Log.Error("文件读取失败: %v", err)
			c.JSON(502, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			until.Log.Error("Failed to stat file: %v", err)
			c.JSON(500, gin.H{"error": "Failed to get file info"})
			return
		}

		// 设置Content-Length
		c.DataFromReader(http.StatusOK, stat.Size(), "application/java-archive", file, nil)
		return
	}
	c.JSON(500, gin.H{})
}

// GetLocalMavenRepo 本地maven仓库扫描，识别当前环境中的所有资源
func (t *MavenService) GetLocalMavenRepo() {
	var localMavenGav []*model.MavenArtifactModel
	localRepoPath := t.applicationConfig.LocalRepoConfig.LocalPath
	if localRepoPath == "" || localRepoPath == "." || localRepoPath == "/" {
		until.Log.Error(`无效路径`)
		return
	}
	filePaths := scan(localRepoPath, 10)
	for _, path := range filePaths {
		gavPath := strings.Replace(path, localRepoPath, "", 1)
		mavenArtifact := parseMavenPath(gavPath)
		if mavenArtifact != nil {
			mavenArtifact.FilePath = path
			localMavenGav = append(localMavenGav, mavenArtifact)
		} else {
			until.Log.Error(`解析失败：` + path)
		}
	}
	all, err := t.mavenRepo.FindAll()
	if err != nil {
		until.Log.Error(`数据库查询失败` + err.Error())
		return
	}
	allDataLen := len(all)
	until.Log.Info("当前查询: ", allDataLen)
	// 将获取到的数据一条条遍历并判断那些需要插入数据库 那些需要废弃
	mavenGavMap := make(map[string]model.MavenArtifactModel, allDataLen)
	for _, artifactModel := range all {
		mavenGavMap[artifactModel.Key] = artifactModel
	}
	var needInsertData []*model.MavenArtifactModel
	for _, mavenGAV := range localMavenGav {
		_, exists := mavenGavMap[mavenGAV.Key]
		if !exists {
			needInsertData = append(needInsertData, mavenGAV)
		}
	}
	e := t.mavenRepo.Save(needInsertData)
	if e != nil {
		panic(e)
	}
	until.Log.Info(fmt.Sprintf("当前扫描结果%d/%d", len(needInsertData), len(localMavenGav)))
}

// 将路径信息转换成maven坐标信息
func parseMavenPath(path string) *model.MavenArtifactModel {
	systemSeparator := string(filepath.Separator)
	res := &model.MavenArtifactModel{}
	// 标准化路径并分割
	normalized := filepath.Clean(path)
	index := strings.Index(normalized, systemSeparator)
	if index == 0 {
		normalized = normalized[1:]
	}
	res.Key = normalized
	parts := strings.Split(normalized, systemSeparator)

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

func scan(path string, maxDepth int) []string {
	var res []string
	rootDepth := strings.Count(path, string(os.PathSeparator))
	if path == "" || path == "." || path == "."+string(os.PathSeparator) {
		rootDepth = 0 // 当前目录的特殊处理
	}
	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			println(err)
		}
		depth := strings.Count(path, string(os.PathSeparator)) - rootDepth
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			res = append(res, path)
		}
		return nil
	})

	if err != nil {
		return nil
	}
	return res
}
