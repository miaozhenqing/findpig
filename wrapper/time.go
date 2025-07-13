package wrapper

import (
	"findpig/util"
	"fmt"
	"time"
)

type TimeWrapper interface {
	Check(userId int64) bool
}

type BetweenTimeWrapper struct {
	Start time.Time
	End   time.Time
}

func (bw *BetweenTimeWrapper) Check(userId int64) bool {
	return util.IsTimeBetween(bw.Start, bw.End)
}

// 获取时间包装
func GetTimeWrapper(timeType int, timeStartStr string, timeEndStr string) TimeWrapper {
	if timeType == 1 {
		s := util.ParseTime(timeStartStr)
		e := util.ParseTime(timeEndStr)
		timeWrapper := BetweenTimeWrapper{Start: s, End: e}
		return &timeWrapper
	}
	panic(fmt.Sprintf("unexpected time type: %d", timeType))
}
