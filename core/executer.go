package core

import (
	"findpig/cache"
	"findpig/db/repository/inter"
	"findpig/wrapper"
	"time"
)

func ExecByRule(userId int64, ruleId int, param any) {
	subWrappers := GetSubActsByRule(ruleId)
	if len(subWrappers) == 0 {
		return
	}
	dataManager := wrapper.NewDataManager(cache.NewCache(time.Hour*24), GetActRepo())
	validContexts := make([]*wrapper.ContextWrapper[any], 0, len(subWrappers))
	for _, subWrapper := range subWrappers {
		timeCheck := subWrapper.DoTime.Check(userId)
		if timeCheck {
			validContexts = append(validContexts, &wrapper.ContextWrapper[any]{UserId: userId,
				SubWrapper: subWrapper, TriggerRule: ruleId, TriggerParam: param, DataManager: dataManager})
		}
	}
	dataManager.ContextWrappers = append(dataManager.ContextWrappers, validContexts...)
	if len(validContexts) == 0 {
		return
	}
	for _, contextWrapper := range validContexts {

		existSuccess := false
		for _, ruleWrapper := range contextWrapper.SubWrapper.Rules {
			if ruleWrapper.Id == contextWrapper.TriggerRule {
				//获取规则对应的参数解析器
				succ, progress := ruleWrapper.ParamParser.Parse(contextWrapper.TriggerParam)
				if succ {
					existSuccess = true
					if contextWrapper.ParserResult == nil {
						contextWrapper.ParserResult = make([]*wrapper.ParserResult, 0, 1)
					}
					contextWrapper.ParserResult = append(contextWrapper.ParserResult, &wrapper.ParserResult{
						RuleWrapper:       ruleWrapper,
						SuccessfulProcess: progress,
					})
				}
			}
		}
		if !existSuccess {
			continue
		}
		//校验前置条件
		condCheck := true
		for _, condition := range contextWrapper.SubWrapper.Conditions {
			if !condition.Check(userId) {
				condCheck = false
				break
			}
		}
		if !condCheck {
			continue
		}
		//累计进度
		activity := contextWrapper.GetActivity()
		if !activity.IsInProgress() {
			continue
		}
		for _, result := range contextWrapper.ParserResult {
			index := result.RuleWrapper.Index
			progress := activity.ActProgresses[index]
			success := result.RuleWrapper.Accumulator.Accum(result.SuccessfulProcess, result.RuleWrapper.Target, progress)
			if success {
				activity.UpdateTime = time.Now().UnixMilli()
				activity.ModifyTime = time.Now().UnixMilli()
				contextWrapper.SuccessFlag = true
				activity.Dirty = true
			}
		}
		//校验是否全部完成
		allCompleted := true
		for _, ruleWrapper := range contextWrapper.SubWrapper.Rules {
			progress := activity.ActProgresses[ruleWrapper.Index]
			completed := ruleWrapper.Accumulator.CheckCompleted(progress, ruleWrapper.Target)
			if !completed {
				allCompleted = false
			}
		}
		if allCompleted {
			activity.MarkCompleted()
		}
	}
	dataManager.UpdateDirtyActivity()
}

func GetActRepo() inter.ActRepo {
	return Repo
}
