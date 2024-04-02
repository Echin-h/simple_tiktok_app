package util

import (
	"strconv"
	"strings"
)

// videoId::userId -- string
func Connect(videoId int64, userId int64) string {
	var str strings.Builder
	str.WriteString(strconv.FormatInt(videoId, 10))
	str.WriteString("::")
	str.WriteString(strconv.FormatInt(userId, 10))
	return str.String()
}

func Separate(Id string) (int64, int64) {
	StringId := strings.Split(Id, "::")
	videoId := StringId[0]
	userId := StringId[1]
	video, _ := strconv.ParseInt(videoId, 10, 64)
	user, _ := strconv.ParseInt(userId, 10, 64)
	return video, user

}
