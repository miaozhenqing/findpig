package wrapper

// 主活动
type MainWrapper struct {
	Id       int
	Name     string
	Desc     string
	Subs     []int
	ShowTime TimeWrapper
}

// 子活动
type SubWrapper struct {
	Id         int
	MainId     int
	Name       string
	Desc       string
	Rules      []*RuleWrapper
	Conditions []ConditionWrapper
	ShowTime   TimeWrapper
	DoTime     TimeWrapper
	RewardTime TimeWrapper
	Flush      *FlushWrapper
}
