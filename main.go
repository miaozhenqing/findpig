package main

import (
	"findpig/core"
	"findpig/db/config"
	"findpig/db/repository/impl"
	"findpig/template"
)

func main() {

	testMainTpl := template.MainTpl{
		Id:            1,
		Name:          "主活动1",
		Desc:          "主活动1",
		Subs:          []int{1, 2},
		ShowTimeType:  1,
		ShowTimeStart: "2025-07-03 21:35:17",
		ShowTimeEnd:   "2025-11-03 21:35:17",
	}

	testSubTpl1 := template.SubTpl{
		Id:                 1,
		Name:               "子活动1",
		Desc:               "子活动1",
		RuleIds:            []int{1, 2},
		RuleParams:         nil,
		RuleTargets:        []int{1, 7},
		ShowTimeType:       1,
		ShowTimeStart:      "2025-07-03 21:35:17",
		ShowTimeEnd:        "2025-11-03 21:35:17",
		DoTimeType:         1,
		DoTimeStart:        "2025-07-03 21:35:17",
		DoTimeEnd:          "2025-11-03 21:35:17",
		RewardTimeType:     1,
		RewardTimeStart:    "2025-07-03 21:35:17",
		RewardTimeEnd:      "2025-11-03 21:35:17",
		DoConditionTypes:   []int{1},
		DoConditionParams:  nil,
		DoConditionTargets: []int{30},
		FlushType:          1,
	}
	testSubTpl2 := template.SubTpl{
		Id:                 2,
		Name:               "子活动2",
		Desc:               "子活动2",
		RuleIds:            []int{1, 2},
		RuleParams:         nil,
		RuleTargets:        []int{1, 7},
		ShowTimeType:       1,
		ShowTimeStart:      "2025-07-03 21:35:17",
		ShowTimeEnd:        "2025-11-03 21:35:17",
		DoTimeType:         1,
		DoTimeStart:        "2025-07-03 21:35:17",
		DoTimeEnd:          "2025-11-03 21:35:17",
		RewardTimeType:     1,
		RewardTimeStart:    "2025-07-03 21:35:17",
		RewardTimeEnd:      "2025-11-03 21:35:17",
		DoConditionTypes:   []int{1},
		DoConditionParams:  nil,
		DoConditionTargets: []int{30},
		FlushType:          1,
	}

	subs := []template.SubTpl{testSubTpl1, testSubTpl2}
	core.RegisterAct(testMainTpl, subs)

	//注册db-repo
	db := config.InitDB()
	repo := impl.NewActGormRepo(db)
	core.RegisterActRepo(repo)

	//测试数据
	userId := int64(123)
	triggerRule := 1
	triggerParam := 1

	core.ExecByRule(userId, triggerRule, triggerParam)

}
