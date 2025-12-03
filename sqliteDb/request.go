package sqliteDb

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	ID            int64      `gorm:"primaryKey"`
	Active        int        //是否显示当前请求详情
	CollectionsId int64      //项目空间主键ID
	Name          string     //请求名
	Addr          string     //请求url
	Content       string     //请求内容
	CreateTime    *time.Time `gorm:"type:datetime"`
	UpdateTime    *time.Time `gorm:"type:datetime"`
}

type RequestModel struct {
	Conn *gorm.DB
}

func NewRequestModel() *RequestModel {
	conn := conn()
	return &RequestModel{Conn: conn}
}
func (u *RequestModel) Create(requestInfo Request) (int64, error) {
	// 迁移 schema
	u.Conn.AutoMigrate(&Request{})
	result := u.Conn.Create(&requestInfo) // 通过数据的指针来创建
	return requestInfo.ID, result.Error
}

func (u *RequestModel) Update(requestId int64, update Request) (int64, error) {
	var info Request
	result := u.Conn.Model(&info).Where("id = ?", requestId).Updates(&update) // 通过数据的指针来创建
	return result.RowsAffected, result.Error
}

func (u *RequestModel) UpdateByNotMap(requestId int64, update map[string]interface{}) (int64, error) {
	var info Request
	result := u.Conn.Model(&info).Where("id != ?", requestId).Updates(&update) // 通过数据的指针来创建
	return result.RowsAffected, result.Error
}

func (u *RequestModel) Delete(id int64) (int64, error) {
	var info Request
	result := u.Conn.Delete(&info, id)
	return result.RowsAffected, result.Error
}

func (u *RequestModel) FindByCollectionsId(collectionsId int64) ([]Request, error) {
	var list []Request
	result := u.Conn.Where("collections_id = ?", collectionsId).Find(&list)
	return list, result.Error
}
