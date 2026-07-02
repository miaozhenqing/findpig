package wrapper

import (
	"findpig/db/model"
	"fmt"
)

type Accumulator interface {
	Accum(parameter int, target int, actProg *model.ActProgress) bool
	CheckCompleted(actProg *model.ActProgress, target int) bool
}

type AccumInt struct {
}

func (a AccumInt) Accum(parameter int, target int, actProg *model.ActProgress) bool {
	if a.CheckCompleted(actProg, target) {
		return false
	}
	if actProg.P == nil {
		actProg.P = parameter
	} else {
		var old int
		if v, ok := actProg.P.(float64); ok {
			old = int(v)
		} else {
			old = actProg.P.(int)
		}
		actProg.P = old + parameter
	}
	return true
}
func (a AccumInt) CheckCompleted(actProg *model.ActProgress, target int) bool {
	if actProg.P == nil {
		return false
	}
	var v int
	if vf, ok := actProg.P.(float64); ok {
		v = int(vf)
	} else {
		v = actProg.P.(int)
	}
	return v >= target
}

//todo 完成了int类型的进度累计器，下一步尝试将actProgress从actEntity剥离开来

func GetAccumulatorByRule(rule int) Accumulator {
	if rule == 1 {
		return &AccumInt{}
	} else if rule == 2 {
		return &AccumInt{}
	}
	panic(fmt.Errorf("GetAccumulatorByRule called with rule %d", rule))
}
