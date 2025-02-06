package config

import "github.com/spf13/viper"

func init(){
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func PrintConfig() {
	println("Config: ", viper.AllSettings())
}