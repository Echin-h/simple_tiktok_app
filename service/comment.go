package service

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"tiktok/common/db"
	"tiktok/common/model"
)

type Comment struct{}

type CommentActionReq struct {
	CommentId    int64  `json:"comment_id"`
	VideoId      int64  `json:"video_id" binding:"required"`
	Content      string `json:"content" binding:"required"`
	ParentId     int64  `json:"parent_id"`
	ActionType   byte   `json:"action_type"`
	DelCommentId int64  `json:"del_comment_id"`
	DelVideoId   int64  `json:"del_video_id"`
}

type CommentActionResp struct {
	CommentId int64  `json:"comment_id"`
	Content   string `json:"content"`
	//...
}

type CommentListReq struct {
	CommentId int64 `json:"comment_id"`
	VideoId   int64 `json:"video_id"`
}

type CommentListResp struct {
	Id         int64    `json:"id"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
	User       UserResp `json:"user"`
}

type UserResp struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
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

	var commentResp CommentActionResp
	commentResp.CommentId = commentId
	commentResp.Content = content
	return commentResp, err
}

func (comm *Comment) CommentList(videoId int64, commentId int64) (interface{}, error) {
	var comments []model.Comment
	if commentId == 0 {
		err := db.GetMysql().Debug().Model(&model.Comment{}).Where("video_id = ?", videoId).Find(&comments).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println("the video have no comments")
				return nil, err
			} else {
				log.Println("list comments is wrong")
				return nil, err
			}
		}
	} else {
		err := db.GetMysql().Debug().Model(&model.Comment{}).
			Where("parent_id = ? AND video_id = ?", commentId, videoId).
			Find(&comments).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println("comments are none")
				return nil, err
			} else {
				log.Println("list son comments is wrong")
				return nil, err
			}
		}
	}
	resp, err := genListResp(comments)
	return resp, err
}

func (comm *Comment) CommentDelete(commentId int64, userId int64, videoId int64) error {
	//var delComm model.Comment
	tx := db.GetMysql().Begin()
	//err := tx.Debug().Model(&model.Comment{}).Where("id = ? and user_id = ? and video_id = ?", commentId, userId, videoId).
	//	First(&delComm).Error
	//if err != nil {
	//	log.Println("find the record is wrong")
	//	tx.Rollback()
	//	return err
	//}
	err := tx.Debug().Model(&model.Comment{}).
		Where("id = ? and user_id = ? and video_id = ?", commentId, userId, videoId).
		Delete(&model.Comment{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			return errors.New("the committed record is not exist")
		} else {
			tx.Rollback()
			return err
		}
	}
	// 同时删除该评论下的所有子评论
	err = tx.Debug().Model(&model.Comment{}).Where("parent_id = ?", commentId).Delete(&model.Comment{}).Error
	if err != nil {
		log.Println("删除子评论错误")
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func (comm *Comment) CommentReply(userId int64, videoId int64, content string, commentId int64, parentId int64) (interface{}, error) {
	newComm := model.Comment{
		ID:       commentId,
		UserId:   userId,
		Content:  content,
		VideoId:  videoId,
		ParentId: parentId,
	}

	tx := db.GetMysql().Begin()
	err := tx.Debug().Model(&model.Comment{}).Create(&newComm).Error
	if err != nil {
		log.Println("create comment error")
		tx.Rollback()
		return nil, err
	}
	tx.Debug().Model(&model.Video{}).Where("id = ?", videoId).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	err = tx.Commit().Error

	var commentResp CommentActionResp
	commentResp.CommentId = commentId
	commentResp.Content = content
	return commentResp, err
}

func genListResp(comments []model.Comment) ([]CommentListResp, error) {
	lenSize := len(comments)
	comms := make([]CommentListResp, lenSize)
	for i, v := range comments {
		comms[i].Id = v.ID
		comms[i].Content = v.Content
		comms[i].CreateDate = v.CreatedAt.Format("2006-01-02 03:04:05 PM")
		comms[i].User = UserResp{
			Id:            v.User.ID,
			Name:          v.User.UserName,
			FollowCount:   v.User.FollowerCount,
			FollowerCount: v.User.FollowerCount,
		}
	}
	return comms, nil
}
