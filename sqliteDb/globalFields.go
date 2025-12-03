package sqliteDb

import (
	"gorm.io/gorm"
)

type GlobalField struct {
	Field string `gorm:"uniqueIndex"` //字段名
	Value string
}

type GlobalFieldModel struct {
	Conn *gorm.DB
}

func NewGlobalFieldModel() *GlobalFieldModel {
	conn := conn()
	return &GlobalFieldModel{Conn: conn}
}
func (u *GlobalFieldModel) Create(info []GlobalField) (int64, error) {
	// 迁移 schema
	u.Conn.AutoMigrate(&GlobalField{})
	result := u.Conn.Create(&info) // 通过数据的指针来创建
	return result.RowsAffected, result.Error
}

func (u *GlobalFieldModel) GetList() ([]GlobalField, error) {
	var list []GlobalField
	result := u.Conn.Find(&list)
	return list, result.Error
}

func (u *GlobalFieldModel) Update(requestId int64, update GlobalField) (int64, error) {
	var info GlobalField
	result := u.Conn.Model(&info).Where("id = ?", requestId).Updates(&update) // 通过数据的指针来创建
	return result.RowsAffected, result.Error
}

func (u *GlobalFieldModel) UpdateByNotMap(requestId int64, update map[string]interface{}) (int64, error) {
	var info GlobalField
	result := u.Conn.Model(&info).Where("id != ?", requestId).Updates(&update) // 通过数据的指针来创建
	return result.RowsAffected, result.Error
}

func (u *GlobalFieldModel) Delete(id int64) (int64, error) {
	var info GlobalField
	result := u.Conn.Delete(&info, id)
	return result.RowsAffected, result.Error
}

func (u *GlobalFieldModel) BatchDelete() (int64, error) {
	var info GlobalField
	result := u.Conn.Where("1 = 1").Delete(&info)
	return result.RowsAffected, result.Error
}

func (u *GlobalFieldModel) FindByCollectionsId(collectionsId int64) ([]GlobalField, error) {
	var list []GlobalField
	result := u.Conn.Where("collections_id = ?", collectionsId).Find(&list)
	return list, result.Error
}
