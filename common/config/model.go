package config

import "time"

type Config struct {
	App    App
	Mysql  Mysql
	Redis  Redis
	Aliyun Aliyun
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

type Redis struct {
	Address  string
	Password string
	Db       int
}

type Aliyun struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}
