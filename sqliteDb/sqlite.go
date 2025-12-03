package sqliteDb

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"path/filepath"

	// 根据安装选择驱动
	"github.com/glebarez/sqlite"
)

var sqliteDb *gorm.DB

func init() {
	sqliteDb = conn()
}

func ensureFileExists(src, dst string) {
	// 1. 检查目标文件是否存在
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		// 2. 确保目标目录存在（关键！）
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		// 3. 拷贝文件（安全！）
		srcFile, err := os.Open(src)
		if err != nil {
			log.Fatalf("打开源文件失败: %v", err)
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			log.Fatalf("创建目标文件失败: %v", err)
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			log.Fatalf("拷贝文件失败: %v", err)
		}

		log.Printf("文件已创建: %s", dst)
	}
}

func conn() *gorm.DB {
	homeDir, _ := os.UserHomeDir()
	dbPath := filepath.Join(homeDir, ".wailsPostman")
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	dbPath = filepath.Join(dbPath, "database.db")
	ensureFileExists("./test.db", dbPath)
	// 基础连接（文件数据库）
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("数据库连接失败")
	}
	return db
}
