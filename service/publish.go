package service

import (
	"bytes"
	"errors"
	"github.com/h2non/filetype"
	"io"
	"log"
	"mime/multipart"
	"tiktok/common/db"
	l "tiktok/common/log"
	"tiktok/common/model"
	"tiktok/common/oss"
)

var BucketName = "simple-tiktok-app"

type Video struct{}

// PublishAction 已登陆的用户上传视频
func (v Video) PublishAction(data *multipart.FileHeader, title string, userId int64, content string) error {
	file, err := data.Open()
	if err != nil {
		return err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	// 判断是否为视频 (important)
	check, err := data.Open()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, check); err != nil {
		log.Println("copy file error", err)
		return err
	}
	if filetype.IsVideo(buf.Bytes()) == false {
		log.Println("file is not video")
		return errors.New("not a video")
	}
	check.Close()

	right, err := oss.UploadVideoToOSS(BucketName, data.Filename, file)
	if err != nil || !right {
		return err
	}

	urlString, imgString, _ := oss.GetOssVideoUrlAndImgUrl(BucketName, data.Filename)
	if urlString == "" || imgString == "" {
		return errors.New("upload video error, can't get url")
	}

	video := &model.Video{
		Title:         title,
		UserId:        userId,
		PicUrl:        imgString,
		PlayUrl:       urlString,
		Content:       content,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	l.Logger.Info("start to save video to db")
	err = db.GetMysql().Model(&model.Video{}).Create(video).Error
	if err != nil {
		return err
	}
	l.Logger.Info("save video to db success")
	return nil
}

// 注意： 这里的VideoDemo是可以自定义自己想要返回的值
func (v Video) PublishList(userId int64) (interface{}, error) {
	var videosDemo []model.VideoDemo
	var videos []model.Video
	err := db.GetMysql().Model(&model.Video{}).Where("user_id = ?", userId).Find(&videos).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, v := range videos {
		videosDemo = append(videosDemo, model.VideoDemo{
			ID:            v.ID,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			PicUrl:        v.PicUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			CreatedAt:     v.CreatedAt,
			Content:       v.Content,
			UserId:        v.UserId,
			Comments:      v.Comments,
		})
	}
	return videosDemo, nil
}
