package model

import (
	"cargo-m/internal/until"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id         *uint64    `gorm:"primary_key;AUTO_INCREMENT;comment:主键ID"`
	Valid      int        `gorm:"default:1"`
	Sort       *int       `gorm:"default:0"`
	CreateTime *time.Time `gorm:"type:datetime;not null"`
	ModifyTime *time.Time `gorm:"type:datetime;not null"`
}

// BeforeCreate 数据插入数据之前处理
func (e *BaseModel) BeforeCreate(tx *gorm.DB) error {
	// 单独设置默认值逻辑（如果有需要）
	if e.Id == nil {
		id, err := until.IdGenerate.NextId()
		if err != nil {
			return err
		}
		e.Id = &id
	}
	if e.Valid == 0 {
		e.Valid = 1
	}
	now := time.Now()
	if e.CreateTime == nil {
		e.CreateTime = &now
	}
	if e.ModifyTime == nil {
		e.ModifyTime = &now
	}
	return nil
}
