package oss

import (
	"fmt"
	"testing"
	"tiktok/common/log"
)

func TestNew(t *testing.T) {
	log.Init()
	AliyunInit()
	CreateBucket("simple-tiktok-app")
	log.Logger.Info("TestNew is success")
	fmt.Println("TestNew is success")
}
