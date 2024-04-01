package service

import (
	"gorm.io/gorm"
	"log"
	"tiktok/common/db"
	"tiktok/common/model"
)

type Comment struct{}

type CommentActionReq struct {
	CommentId  int64  `json:"comment_id"`
	VideoId    int64  `json:"video_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	ParentId   int64  `json:"parent_id"`
	ActionType byte   `json:"action_type"`
}

type CommentResp struct {
	CommentId int64  `json:"comment_id"`
	Content   string `json:"content"`
	//...
}

func (comm *Comment) CommentPublish(commentId int64, userId int64, videoId int64, content string) (interface{}, error) {
	comment := model.Comment{
		ID:       commentId,
		UserId:   userId,
		VideoId:  videoId,
		Content:  content,
		ParentId: 0,
	}
	tx := db.GetMysql().Begin()
	err := tx.Debug().Model(&model.Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("create comment error")
		tx.Rollback()
		return nil, err
	}
	tx.Debug().Model(&model.Video{}).Where("id = ?", videoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	err = tx.Commit().Error

	var commentResp CommentResp
	commentResp.CommentId = commentId
	commentResp.Content = content
	return commentResp, err
}

func (comm *Comment) CommentList(videoId int64) {

}

func (comm *Comment) CommentDelete(commentId int64, userId int64) {

}

func (comm *Comment) CommentReply(userId int64, videoId int64, content string, commentId int64, parentId int64) {

}
