package template

type SubTpl struct {
	//基本信息
	Id   int
	Name string
	Desc string
	//规则配置
	RuleIds     []int
	RuleParams  []string
	RuleTargets []int
	//显示时间
	ShowTimeType  int
	ShowTimeStart string
	ShowTimeEnd   string
	//累计时间
	DoTimeType  int
	DoTimeStart string
	DoTimeEnd   string
	//领奖时间
	RewardTimeType  int
	RewardTimeStart string
	RewardTimeEnd   string
	//累计条件
	DoConditionTypes   []int
	DoConditionParams  []string
	DoConditionTargets []int
	//刷新规则
	FlushType int
}

type MainTpl struct {
	//基本信息
	Id   int
	Name string
	Desc string
	//显示时间
	ShowTimeType  int
	ShowTimeStart string
	ShowTimeEnd   string
	//子活动列表
	Subs []int
}
