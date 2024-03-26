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
	v := viper.New()

	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("../")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		log.Println("Read config error: ", err)
		return
	}
	err1 := v.Unmarshal(&c)
	if err1 != nil {
		log.Println("Unmarshal config error: ", err)
		return
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		err := v.ReadInConfig()
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
