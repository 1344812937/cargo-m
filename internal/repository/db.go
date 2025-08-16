package repository

import (
	"cargo-m/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open(".\\data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.MavenArtifactModel{})
	Db = *db
}
