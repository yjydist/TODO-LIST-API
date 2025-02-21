package model

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	username = viper.GetString("mysql.username")
	password = viper.GetString("mysql.password")
	host     = viper.GetString("mysql.host")
	port     = viper.GetInt("mysql.port")
	dbname   = viper.GetString("mysql.dbname")
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect MySQL:", err)
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbname)
	err = DB.Exec(sql).Error
	if err != nil {
		log.Fatal("Unable to create database:", err)
	}
	fmt.Println("Database created successfully or already exists")

	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	DB, err = gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}

	// **Step 4: 使用 AutoMigrate 创建表（如果不存在）**
	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}
	fmt.Println("数据表已创建或已存在")
}
