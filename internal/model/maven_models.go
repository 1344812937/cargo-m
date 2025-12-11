package model

type MavenArtifactModel struct {
	*BaseModel
	Key        string `gorm:"key;"`
	GroupId    string `gorm:"group_id;comment:组ID"`
	ArtifactId string `gorm:"artifact_id"`
	Version    string `gorm:"version;comment:版本号"`
	FileName   string `gorm:"file_name;comment:文件名称"`
	Classifier string `gorm:"classifier;"`
	Extension  string `gorm:"extension;comment:扩展名"`
	FilePath   string `gorm:"file_path;comment:文件夹路径"`
}

func (t *MavenArtifactModel) TableName() string {
	return "maven_artifact"
}
