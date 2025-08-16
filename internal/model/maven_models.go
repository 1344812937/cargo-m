package model

import (
	"cargo-m/internal/until"

	"gorm.io/gorm"
)

type MavenArtifactModel struct {
	Id         *uint64 `gorm:"primary_key;auto_increment;comment:主键ID"`
	Key        string  `gorm:"key;"`
	GroupId    string  `gorm:"group_id;comment:组ID"`
	ArtifactId string  `gorm:"artifact_id"`
	Version    string  `gorm:"version;comment:版本号"`
	FileName   string  `gorm:"file_name;comment:文件名称"`
	Classifier string  `gorm:"classifier;"`
	Extension  string  `gorm:"extension;comment:扩展名"`
	FilePath   string  `gorm:"file_path;comment:文件夹路径"`
	Valid      int     `gorm:"valid" default:"1"`
}

func (t *MavenArtifactModel) BeforeCreate(tx *gorm.DB) error {
	if t.Id == nil {
		id, err := until.IdGenerate.NextId()
		if err != nil {
			return err
		}
		t.Id = &id
	}
	if t.Valid == 0 {
		t.Valid = 1
	}
	return nil
}

func (t *MavenArtifactModel) TableName() string {
	return "maven_artifact"
}
