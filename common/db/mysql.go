package db

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"tiktok/common/config"
	"tiktok/common/model"
)

type Mysql struct {
	*gorm.DB
}

type User struct {
	model.User
}

var mysqlDB = &Mysql{}

func MysqlInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Get().Mysql.User,
		config.Get().Mysql.Password,
		config.Get().Mysql.Address,
		config.Get().Mysql.DBName)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 191,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tiktok_",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Println("failed to connect database", err)
		return
	}
	mysqlDB = &Mysql{db}
	err = mysqlDB.AutoMigrateTable()
	if err != nil {
		log.Println("failed to migrate database", err)
		return
	}

}

func (m *Mysql) AutoMigrateTable() error {
	err := m.AutoMigrate(
		&model.User{},
		&model.UnFollow{},
		&model.Video{},
		&model.Comment{},
		&model.Follow{},
		&model.VideoFavorite{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate table: %v", err)
	}
	return nil
}

func GetMysql() *Mysql {
	return mysqlDB
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = int64(uuid.New().ID())
	return nil
}
