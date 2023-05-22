package utils

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("conf/")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Viper load error: %v", err)
	}
}