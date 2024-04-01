package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"tiktok/common/config"
	"tiktok/common/log"
	_ "tiktok/common/log"
)

// In OSS , object = file , bucket = folder , but you can create a folder in the bucket

var AliyunClient *oss.Client

func AliyunInit() {
	client, err := oss.New(
		config.Get().Aliyun.Endpoint,
		config.Get().Aliyun.AccessKeyID,
		config.Get().Aliyun.AccessKeySecret)
	if err != nil {
		log.Logger.Error("aliyun oss init failed")
		fmt.Println("aliyun oss init failed:", err)
		return
	}
	AliyunClient = client
	fmt.Println("aliyun oss init success")
}

// truly , i don't know where to set the CreateBucket function
func CreateBucket(name string) {
	err := AliyunClient.CreateBucket(name)
	if err != nil {
		exist, err := AliyunClient.IsBucketExist(name)
		if err == nil && exist {
			log.Logger.Info(fmt.Sprintf("We already own %s\n", name))
		} else {
			log.Logger.Error("create bucket error")
			return
		}
	}
	log.Logger.Info(fmt.Sprintf("Successfully created %s\n", name))
}

func UploadVideoToOSS(bucketName string, objectName string, reader multipart.File) (bool, error) {
	bucket, err := AliyunClient.Bucket(bucketName) //"simple-tiktok-app"
	if err != nil {
		log.Logger.Error("get bucket error")
		return false, err
	}
	err = bucket.PutObject(objectName, reader)
	// bucket.UploadFile(objectKey, filePath string, partSize int64, options ...Option)
	if err != nil {
		log.Logger.Error("put object error")
		return false, err
	}
	log.Logger.Info("upload video to oss success")
	return true, nil
}

// if you want to get the object under the folder
// the objectName is such like "folderName/objectName"

func GetOssVideoUrlAndImgUrl(bucketName string, objectName string) (string, string, error) {
	url := "https://" + bucketName + "." + "oss-cn-hangzhou.aliyuncs.com" + "/" + objectName
	return url, url + "?x-oss-process=video/snapshot,t_0,f_jpg,w_0,h_0,m_fast,ar_auto", nil
}
