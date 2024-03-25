package config

import "time"

type Config struct {
	App   App
	Mysql Mysql
}

type App struct {
	Host      string
	Port      string
	JwtSecret string
	Release   string
	RunMode   string
}

type Mysql struct {
	Address     string
	User        string
	Password    string
	DBName      string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}
