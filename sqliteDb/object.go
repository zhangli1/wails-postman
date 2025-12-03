package sqliteDb

import (
	"gorm.io/gorm"
	"time"
)

type Object struct {
	ID         int64      `gorm:"primaryKey"`
	Name       string     `gorm:"uniqueIndex"` //项目名
	Desc       string     //描述
	CreateTime *time.Time `gorm:"type:datetime"`
	UpdateTime *time.Time `gorm:"type:datetime"`
}

type ObjectModel struct {
	Conn *gorm.DB
}

func NewObjectModel() *ObjectModel {
	conn := conn()
	return &ObjectModel{Conn: conn}
}
func (u *ObjectModel) Create(objectInfo Object) (int64, error) {
	// 迁移 schema
	u.Conn.AutoMigrate(&Object{})
	result := u.Conn.Create(&objectInfo) // 通过数据的指针来创建
	return objectInfo.ID, result.Error
}
