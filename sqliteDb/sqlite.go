package sqliteDb

import (
	"fmt"
	"gorm.io/gorm"
	// 根据安装选择驱动
	"github.com/glebarez/sqlite"
)

var sqliteDb *gorm.DB

func init() {
	sqliteDb = conn()
}

func conn() *gorm.DB {
	// 基础连接（文件数据库）
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("数据库连接失败")
	}
	return db
}
