package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var (
	c = Config{}
)

var Path = "app.yaml"

func Init() {
	viper.SetConfigFile(Path)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Read config error: ", err)
		return
	}
	err1 := viper.Unmarshal(&c)
	if err1 != nil {
		log.Println("Unmarshal config error: ", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.ReadInConfig()
		if err != nil {
			log.Println("Config file update; change: ", e.Name)
			return
		}
	})

	fmt.Println(Get().App.Host, Get().App.Port)
}

func Get() Config {
	return c
}
