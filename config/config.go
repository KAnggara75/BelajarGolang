package config

import "github.com/spf13/viper"

func GetPort() string {
	return ":" + viper.GetString("PORT")
}

func GetDatabaseURL() string {
	return viper.GetString("DATABASE_URL")
}

func GetDatabaseDriver() string {
	return viper.GetString("DATABASE_DRIVER")
}
