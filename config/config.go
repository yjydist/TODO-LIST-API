package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Init 初始化配置
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "todo_list")
	viper.SetDefault("jwt.secret", "your_jwt_secret_key")
	viper.SetDefault("jwt.expire", 24) // 小时
	viper.SetDefault("ratelimit.requests_per_second", 10)
	viper.SetDefault("ratelimit.burst", 20)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认值
			fmt.Println("No config file found, using defaults")
		} else {
			// 配置文件存在但读取出错
			panic(fmt.Errorf("fatal error reading config file: %s", err))
		}
	}
}

func PrintConfig() {
	println("Config: ", viper.AllSettings())
}
