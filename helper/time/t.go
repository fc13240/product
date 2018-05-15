package time

import (
	"fmt"
	"time"
)

func ElapseTime(t time.Time) string {

	since := time.Since(t)
	hours := since.Hours()

	if month := int(hours / (24 * 30)); month > 0 {
		return fmt.Sprintf("%d%s", month, "个月前")
	}

	if day := int(hours / (24)); day > 0 {
		return fmt.Sprintf("%d%s", day, "天前")
	}

	if h := int(hours); h > 0 {
		return fmt.Sprintf("%d%s", h, "小时前")
	}

	if min := int(since.Minutes()); min > 1 {
		return fmt.Sprintf("%d%s", min, "分钟前")
	}

	return fmt.Sprintf("%d%s", int(since.Seconds()), "秒前")
}

func String() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Unix() int64 {
	return time.Now().Unix()
}
func Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}
