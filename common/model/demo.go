package model

// this is a demo, unused in the project
import "time"

type VideoDemo struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	PlayUrl       string     `json:"play_url"`
	PicUrl        string     `json:"pic_url"`
	FavoriteCount int64      `json:"favorite_count"`
	CommentCount  int64      `json:"comment_count"`
	CreatedAt     time.Time  `json:"created_at"`
	Content       string     `json:"content"`
	UserId        int64      `json:"user_id"`
	Comments      []*Comment `json:"comments"`
}

type CommentDemo struct {
	ID      int64    `json:"id"`
	User    UserDemo `json:"user"`
	UserId  int64    `json:"user_id"`
	Content string   `json:"content"`
	VideoId int64    `json:"video_id"`
}

type UserDemo struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
