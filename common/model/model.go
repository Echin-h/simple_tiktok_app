package model

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Base
	Name     string   `gorm:"not null;unique;index"`
	UserName string   `gorm:"not null;unique;index"`
	Password string   `gorm:"not null"`
	Videos   []*Video `gorm:"foreignKey:user_id"`
	// GPT说这些计数不应该在数据库层面存储，应该在应用层进行处理
}

type Video struct {
	Base                     // 这里的ID可以表示某一条视频
	Title         string     `gorm:"varchar(20)"`
	Content       string     `gorm:"varchar(100)"`
	UserId        int64      `gorm:"primaryKey"`
	Comments      []*Comment `gorm:"foreignKey:video_id"`
	PicUrl        string     `gorm:"comment:封面地址"`
	PlayUrl       string     `gorm:"comment:播放地址"`
	FavoriteCount int64
	CommentCount  int64
	// 具体的视频信息.......
}

type Comment struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignerKey:user_id"`
	UserId    int64
	Content   string `gorm:"varchar(100);not null"`
	VideoId   int64
	ParentId  int64      `gorm:"default:null;index"` // 评论底下还有评论
	Children  []*Comment `gorm:"-"`
}

// 很巧妙，这个是逻辑上的多表
type Follow struct {
	Id           int64 `gorm:"primaryKey"`           // 关注的人的ID
	User         User  `gorm:"foreignKey:follow_id"` // 关注人
	FollowId     int64
	FollowedUser User `gorm:"foreignKey:user_id"` // 被关注人
	UserId       int64
	CreateTime   int64 `gorm:"autoCreateTime"`
}

type VideoFavorite struct {
	Id         int64 `gorm:"primaryKey"`
	User       User  `gorm:"foreignKey:user_id"`
	UserId     int64
	Video      Video `gorm:"foreignKey:video_id"`
	VideoId    int64
	CreateTime int64 `gorm:"autoCreateTime"`
}

type UnFollow struct {
	Id           int64 `gorm:"primaryKey"`
	User         User  `gorm:"foreignKey:FollowId"` // 关注人
	FollowId     int64
	FollowedUser User `gorm:"foreignKey:UserId"` // 被关注人
	UserId       int64
	CreateTime   int64 `gorm:"autoCreateTime"`
}
