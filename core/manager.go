package core

import (
	"findpig/db/repository/inter"
	"findpig/template"
	"findpig/util"
	"findpig/wrapper"
)

var (
	idToSubActContainer   = make(map[int]*wrapper.SubWrapper)
	ruleToSubActContainer = make(map[int][]*wrapper.SubWrapper)
	idToMainContainer     = make(map[int]*wrapper.MainWrapper)
)
var Repo inter.ActRepo

// 根据规则获取子活动
func GetSubActsByRule(rule int) []*wrapper.SubWrapper {
	wrappers := ruleToSubActContainer[rule]
	return wrappers
}

// 注册活动
func RegisterAct(main template.MainTpl, subTpls []template.SubTpl) {

	for i := range subTpls {
		subTpl := subTpls[i]
		//包装子活动
		subConfig := wrapper.SubWrapper{
			Id:     subTpl.Id,
			MainId: main.Id,
			Name:   subTpl.Name,
			Desc:   subTpl.Desc,
			Rules:  make([]*wrapper.RuleWrapper, 0, len(subTpl.RuleIds)),
		}

		//包装规则
		ruleCount := len(subTpl.RuleIds)
		for i := 0; i < ruleCount; i++ {
			ruleId := util.SafeGet(subTpl.RuleIds, i)
			ruleWrapper := wrapper.NewRuleWrapper(util.SafeGet(subTpl.RuleIds, i),
				util.SafeGet(subTpl.RuleParams, i),
				util.SafeGet(subTpl.RuleTargets, i),
				i)
			subConfig.Rules = append(subConfig.Rules, ruleWrapper)
			//添加到规则容器中
			ruleToSubActContainer[ruleId] = append(ruleToSubActContainer[ruleId], &subConfig)
		}

		//包装条件
		conditionCount := len(subTpl.DoConditionTypes)
		for i := 0; i < conditionCount; i++ {
			conditionWrapper := wrapper.GetActConditionWrapper(util.SafeGet(subTpl.DoConditionTypes, i),
				util.SafeGet(subTpl.DoConditionParams, i),
				util.SafeGet(subTpl.DoConditionTargets, i))
			subConfig.Conditions = append(subConfig.Conditions, conditionWrapper)
		}

		//包装时间
		subConfig.DoTime = wrapper.GetTimeWrapper(subTpl.DoTimeType, subTpl.DoTimeStart, subTpl.DoTimeEnd)
		subConfig.ShowTime = wrapper.GetTimeWrapper(subTpl.ShowTimeType, subTpl.ShowTimeStart, subTpl.ShowTimeEnd)
		subConfig.RewardTime = wrapper.GetTimeWrapper(subTpl.RewardTimeType, subTpl.RewardTimeStart, subTpl.RewardTimeEnd)

		//包装刷新规则
		subConfig.Flush = &wrapper.FlushWrapper{FlushType: subTpl.FlushType}

		//添加到id容器中
		idToSubActContainer[subTpl.Id] = &subConfig
	}

	//包装主活动
	mainConfig := wrapper.MainWrapper{
		Id:   main.Id,
		Name: main.Name,
		Desc: main.Desc,
	}
	//添加子活动
	subIds := make([]int, len(subTpls))
	for i := range subTpls {
		subIds[i] = subTpls[i].Id
	}
	mainConfig.Subs = subIds
	//包装时间
	mainConfig.ShowTime = wrapper.GetTimeWrapper(main.ShowTimeType, main.ShowTimeStart, main.ShowTimeEnd)

	//添加到容器中
	idToMainContainer[main.Id] = &mainConfig
}

func RegisterActRepo(repoImpl inter.ActRepo) {
	Repo = repoImpl
}
