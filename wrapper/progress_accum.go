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
		actProg.P = actProg.P.(int) + parameter
	}
	return true
}
func (a AccumInt) CheckCompleted(actProg *model.ActProgress, target int) bool {
	if actProg.P == nil {
		return false
	}
	return actProg.P.(int) >= target
}

func GetAccumulatorByRule(rule int) Accumulator {
	if rule == 1 {
		return &AccumInt{}
	} else if rule == 2 {
		return &AccumInt{}
	}
	panic(fmt.Errorf("GetAccumulatorByRule called with rule %d", rule))
}
