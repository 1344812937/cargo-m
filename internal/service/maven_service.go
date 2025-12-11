package service

import (
	"cargo-m/internal/config"
	"cargo-m/internal/model"
	"cargo-m/internal/repository"
	"cargo-m/internal/until"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type MavenService struct {
	*BaseService[model.MavenArtifactModel]
	repo              *repository.MavenRepo
	applicationConfig *config.ApplicationConfig
	httpClient        *http.Client
}

func NewMavenService(mavenRepo *repository.MavenRepo, applicationConfig *config.ApplicationConfig) *MavenService {
	return &MavenService{
		BaseService:       NewBaseService(mavenRepo),
		repo:              mavenRepo,
		applicationConfig: applicationConfig,
		httpClient:        http.DefaultClient,
	}
}

func (ms *MavenService) GetRepo(c *gin.Context) {
	mavenGavInfo := parseMavenPath(c.Param("path"))
	if mavenGavInfo != nil {
		key, err := ms.repo.GetByKey(mavenGavInfo.Key)
		if err != nil {
			c.JSON(501, gin.H{})
			return
		}
		if key == nil {
			repoConfig := ms.applicationConfig.LocalRepoConfig
			if len(repoConfig.RemoteRepo) > 0 {
				ms.ProxyHttp(c, repoConfig.RemoteRepo)
				return
			}
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

func (ms *MavenService) CheckRepo(c *gin.Context) {
	mavenGavInfo := parseMavenPath(c.Param("path"))
	if mavenGavInfo != nil {
		key, err := ms.repo.GetByKey(mavenGavInfo.Key)
		if err != nil {
			c.JSON(501, gin.H{})
			return
		}
		if key == nil {
			repoConfig := ms.applicationConfig.LocalRepoConfig
			if len(repoConfig.RemoteRepo) > 0 {
				err := ms.ProxyHttp(c, repoConfig.RemoteRepo)
				if err != nil {
					c.JSON(502, gin.H{"error": err.Error()})
					return
				}
			}
			c.JSON(404, gin.H{})
			return
		}
		c.Status(200)
		return
	}
	c.JSON(500, gin.H{})
}

func (ms *MavenService) ProxyHttp(c *gin.Context, proxyUrl string) error {
	if len(proxyUrl) == 0 {
		return fmt.Errorf("proxy url is empty")
	}
	if strings.LastIndex(proxyUrl, "/") == len(proxyUrl)-1 {
		proxyUrl = proxyUrl[:len(proxyUrl)-1]
	}
	requestUrl := proxyUrl + c.Param("path")
	until.Log.Info("无本地资源转发至：", requestUrl)
	// 创建HTTP请求（复用原始请求的Method、Header等）
	req, err := http.NewRequest(c.Request.Method, requestUrl, c.Request.Body)
	if err != nil {
		return err
	}
	// 复制原始请求的Header（特别是Range头支持断点续传）
	for k, v := range c.Request.Header {
		req.Header[k] = v
	}
	// 发送请求到远程仓库
	resp, err := ms.httpClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 设置响应头（复制远程仓库的Header）
	for k, v := range resp.Header {
		c.Header(k, v[0]) // 取第一个值
	}
	c.Status(resp.StatusCode)

	// 7. 流式传输响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		// 捕获客户端断开连接的错误，避免日志污染
		if strings.Contains(err.Error(), "client disconnected") {
			return nil
		}
		return fmt.Errorf("streaming response failed: %w", err)
	}
	return nil
}

// GetLocalMavenRepo 本地maven仓库扫描，识别当前环境中的所有资源
func (ms *MavenService) GetLocalMavenRepo() {
	var localMavenGav []*model.MavenArtifactModel
	localRepoPath := ms.applicationConfig.LocalRepoConfig.LocalPath
	if localRepoPath == "" || localRepoPath == "." || localRepoPath == "/" {
		until.Log.Error(`无效路径`)
		return
	}
	filePaths := scan(localRepoPath, 10)
	for _, path := range filePaths {
		// 排除无效文件
		if strings.Contains(path, "_remote.repositories") {
			continue
		}
		gavPath := strings.Replace(path, localRepoPath, "", 1)
		mavenArtifact := parseMavenPath(gavPath)
		if mavenArtifact != nil {
			mavenArtifact.FilePath = path
			localMavenGav = append(localMavenGav, mavenArtifact)
		} else {
			until.Log.Error(`解析失败：` + path)
		}
	}
	all, err := ms.repo.FindAll()
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
	e := ms.repo.Save(needInsertData)
	if e != nil {
		panic(e)
	}
	until.Log.Info(fmt.Sprintf("当前扫描结果%d/%d", len(needInsertData), len(localMavenGav)))
}

// 将路径信息转换成maven坐标信息
func parseMavenPath(path string) *model.MavenArtifactModel {
	systemSeparator := string(filepath.Separator)
	res := &model.MavenArtifactModel{BaseModel: &model.BaseModel{}}
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
