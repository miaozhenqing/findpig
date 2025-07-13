package wrapper

// 规则
type RuleWrapper struct {
	Id          int
	Param       string
	Target      int
	Index       int
	ParamParser ParamParser
	Accumulator Accumulator
}

func NewRuleWrapper(ruleId int, param string, target int, index int) *RuleWrapper {
	return &RuleWrapper{
		Id:          ruleId,
		Param:       param,
		Target:      target,
		Index:       index,
		ParamParser: GetParamParserByRule(ruleId),
		Accumulator: GetAccumulatorByRule(ruleId),
	}
}
