package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"tiktok/common/db"
	l "tiktok/common/log"
	"tiktok/common/model"
	"tiktok/util"
)

type VideoFavorite struct{}

type VideoResp struct {
	Id            int64        `json:"id"`
	Author        UserResponse `json:"author"`
	PlayUrl       string       `json:"play_url"`
	PicUrl        string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"`
	CommentCount  int64        `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

type UserResponse struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// HSet("myhash", "key1", "value1", "key2", "value2")
func (favorite *VideoFavorite) FavoriteAction(videoId int64, userId int64) error {
	rdb := db.GetRedis()
	log.Println("videoId: ", videoId, "userId: ", userId)
	latestFlag, err := rdb.HGet(context.Background(), "favoriteHash", util.Connect(videoId, userId)).Result()
	if err != nil {
		log.Println("err: ", err)
		latestFlag = "0"
	}
	l.Logger.Info("latestFlag is set")
	if latestFlag != "1" {
		_, err := rdb.HSet(context.Background(), "favoriteHash", util.Connect(videoId, userId), 1).Result()
		if err != nil {
			log.Println("set hash error: ", err)
			return err
		}
		log.Println("videoId: ", videoId, "userId: ", userId)
		// 点赞计数
		err = rdb.HIncrBy(context.Background(), "favoriteCount", strconv.FormatInt(videoId, 10), 1).Err()
		if err != nil {
			log.Println("redis set incr error: ", err)
			return err
		}
		log.Println("increase like num success")
	}
	return nil
}

func (favorite *VideoFavorite) RemoveFavor(videoId int64, userId int64) error {
	rdb := db.GetRedis()
	latestFlag, err := rdb.HGet(context.Background(), "favoriteHash", util.Connect(videoId, userId)).Result()
	if err != nil {
		log.Println("err: ", err)
		latestFlag = "1"
	}
	log.Println("latestFlag: ", latestFlag)
	if latestFlag != "0" {
		_, err := rdb.HSet(context.Background(), "favoriteHash", util.Connect(videoId, userId), 0).Result()
		if err != nil {
			log.Println("set hash error: ", err)
			return err
		}
		log.Println("video: ", videoId, "userId: ", userId)
		// 视频点赞减一
		err = rdb.HIncrBy(context.Background(), "favoriteCount", strconv.FormatInt(videoId, 10), -1).Err()
		if err != nil {
			log.Println("redis set incr error: ", err)
			return err
		}
		log.Println("decrease like num success")
	}
	return nil
}

func (favorite *VideoFavorite) FavoriteList(videoId int64, userId int64) (*VideoResp, error) {
	var mu sync.Mutex
	func() {
		defer mu.Unlock()
		mu.Lock()
		// TODO: 更新业务--同步
		RegularUpdate()
	}()
	rdb := db.GetRedis()
	// 获取视频点赞数
	favoriteCount, err := rdb.HGet(context.Background(), "favoriteCount", strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		log.Println("get favorite count error: ", err)
		return nil, err
	}
	log.Println("favoriteCount: ", favoriteCount)
	// 获取用户点赞状态
	isFavorite, err := rdb.HGet(context.Background(), "favoriteHash", util.Connect(videoId, userId)).Result()
	if err != nil {
		log.Println("get favorite status error: ", err)
		return nil, err
	}
	log.Println("isFavorite: ", isFavorite)
	// 获取视频信息
	var video model.Video
	err = db.GetMysql().Model(&model.Video{}).Where("id = ?", videoId).First(&video).Error
	if err != nil {
		log.Println("get video info error: ", err)
		return nil, err
	}
	// 获取作者信息
	var author model.User
	err = db.GetMysql().Model(&model.User{}).Where("id = ?", video.UserId).First(&author).Error
	if err != nil {
		log.Println("get author info error: ", err)
		return nil, err
	}
	// 返回数据
	videoResp := VideoResp{
		Id: video.ID,
		Author: UserResponse{
			Id:            author.ID,
			Name:          author.Name,
			FollowCount:   author.FollowCount,
			FollowerCount: author.FollowerCount,
		},
		PlayUrl:       video.PlayUrl,
		PicUrl:        video.PicUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    isFavorite == "1",
		Title:         video.Title,
	}

	return &videoResp, nil
}

func RegularUpdate() {
	UpdateMysql()
	DeleteRedis()
	l.Logger.Info("regular updating!")
}

func UpdateMysql() error {
	logrus.Info("Update starting1...")
	// 更新
	rdb := db.GetRedis()
	// update mysql
	pairs, err := rdb.HGetAll(context.Background(), "FavoriteHash").Result()
	if err != nil {
		logrus.Error("get pairs failed", err)
		return err
	}
	logrus.Info("pairs", pairs)

	for pair, flag := range pairs {
		logrus.Info("Update starting3...")
		videoId, userId := util.Separate(pair)
		logrus.Info("Update starting4...")
		var favors model.VideoFavorite
		favors.UserId = userId
		favors.VideoId = videoId
		logrus.Info(userId, videoId, flag)
		if flag == "1" {
			// 更新点赞表
			// 先删除，再添加
			err = db.GetMysql().Debug().Model(&model.VideoFavorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&model.VideoFavorite{}).Error
			if err != nil {
				logrus.Error("update video favorite_count failed", err)
				return err
			}
			if err := db.GetMysql().Debug().Model(&model.VideoFavorite{}).Create(&favors).Error; err != nil {
				logrus.Error("mysql error in creating video favorite")
			}
			logrus.Info("update video_favorite success")
		} else if flag == "0" {
			if err = db.GetMysql().Debug().Model(&model.VideoFavorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&model.VideoFavorite{}).Error; err != nil {
				logrus.Error("mysql error in deleting video favorite")
			}
		}
		// 更新视频点赞数
		delta, err := rdb.HGet(context.Background(), "FavoriteCount", strconv.FormatInt(videoId, 10)).Result()
		if err != nil {
			logrus.Error("get delta failed", err)
		}
		logrus.Info("delta: ", delta)
		if err := db.GetMysql().Debug().Model(&model.Video{}).
			Where("id = ?", videoId).
			Update("favorite_count", gorm.Expr("favorite_count + ?", delta)).Error; err != nil {
			logrus.Error("mysql error in updating favorite_count")
			return err
		}

	}
	return nil
}
func DeleteRedis() error {
	//视频点赞计数可以直接删除
	err := db.GetRedis().Del(context.Background(), "FavoriteCount").Err()
	if err != nil {
		l.Logger.Error("delete redis error")
		return err
	}
	l.Logger.Info("delete redis count success")
	err = db.GetRedis().Del(context.Background(), "FavoriteHash").Err()
	if err != nil {
		l.Logger.Error("delete redis error")
		return err
	}
	l.Logger.Info("delete redis hash success")
	return nil
}
