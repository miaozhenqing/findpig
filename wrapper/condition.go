package wrapper

import (
	"findpig/util"
	"fmt"
	"time"
)

type ConditionWrapper interface {
	Check(userId int64) bool
}

// 注册天数小于
type RegisterDaysWithinLimit struct {
	Target int
}

func (cond *RegisterDaysWithinLimit) Check(userId int64) bool {
	registerTime := getRegisterTime(userId)
	days := util.DaysBetween(registerTime, time.Now())
	return days <= cond.Target
}

func getRegisterTime(userId int64) time.Time {
	//fake data
	return time.Now()
}

// 获取条件包装
func GetActConditionWrapper(condition int, param string, target int) ConditionWrapper {
	if condition == 1 {
		registerDayLess := RegisterDaysWithinLimit{
			Target: target,
		}
		return &registerDayLess
	}
	panic(fmt.Sprintf("unexpected condition: %d", condition))
}
