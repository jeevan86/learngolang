package tm

import "time"

func TruncToMinuteTs(sec int64) int64 {
	return sec / 60 * 60
}

func CurrentMinuteTs() int64 {
	return time.Now().Unix() / 60 * 60
}
