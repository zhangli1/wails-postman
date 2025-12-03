package sqliteDb

import (
	"gorm.io/gorm"
	"time"
)

type Collections struct {
	ID         int64      `gorm:"primaryKey"`
	ObjectId   int64      //项目主键ID
	Name       string     `gorm:"uniqueIndex"` //目录名
	Desc       string     //描述
	CreateTime *time.Time `gorm:"type:datetime"`
	UpdateTime *time.Time `gorm:"type:datetime"`
}

type CollectionsModel struct {
	Conn *gorm.DB
}

func NewCollectionsModel() *CollectionsModel {
	conn := conn()
	return &CollectionsModel{Conn: conn}
}
func (u *CollectionsModel) Create(collectionsInfo Collections) (int64, error) {
	// 迁移 schema
	u.Conn.AutoMigrate(&Collections{})
	result := u.Conn.Create(&collectionsInfo) // 通过数据的指针来创建
	return collectionsInfo.ID, result.Error
}

func (u *CollectionsModel) Update(requestId int64, update Collections) (int64, error) {
	var info Collections
	result := u.Conn.Model(&info).Where("id = ?", requestId).Updates(&update) // 通过数据的指针来创建
	return result.RowsAffected, result.Error
}

func (u *CollectionsModel) Delete(id int64) (int64, error) {
	var info Collections
	result := u.Conn.Delete(&info, id)
	return result.RowsAffected, result.Error
}

func (u *CollectionsModel) GetList() ([]Collections, error) {
	var list []Collections
	result := u.Conn.Find(&list)
	return list, result.Error
}
