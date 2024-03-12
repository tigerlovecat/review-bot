package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"unicode"
)

var (
	BaseFormatTime      = "2006-01-02 15:04:05"
	MaxBatchInsertLimit = 2000
)

func TimeInt64ToString(baseTime int64) string {
	return time.Unix(baseTime, 0).Format(BaseFormatTime)
}

// DiffTimeForShow 将时间戳转为 HH:MM:SS
func DiffTimeForShow(diff int64) string {

	timeDiff := time.Duration(diff * int64(time.Second))

	// 将时间戳转换为HH:MM:SS格式
	hours := int(timeDiff.Hours())
	minutes := int(timeDiff.Minutes()) % 60
	seconds := int(timeDiff.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func GamerNameForShow(name string) string {
	for key, i2 := range name {
		if key == 0 && unicode.Is(unicode.Scripts["Han"], i2) {
			// 是中文
			nameForShow := name[0:3]
			return string(nameForShow) + "*"
		} else {
			nameForShow := name[0]
			return string(nameForShow) + "*"
		}
	}
	return "**"
}

func Md5String(str string) string {

	// 创建md5哈希对象
	hash := md5.New()

	// 将密码添加到哈希对象中
	hash.Write([]byte(str))

	// 计算md5哈希值
	hashedPassword := hash.Sum(nil)

	// 将哈希值转换为16进制字符串
	hexHashedPassword := hex.EncodeToString(hashedPassword)

	return hexHashedPassword
}
