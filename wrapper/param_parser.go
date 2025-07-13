package wrapper

import "fmt"

// 参数解析器
type ParamParser interface {
	Parse(any) (bool, int)
}

type ParserResult struct {
	RuleWrapper       *RuleWrapper
	SuccessfulProcess int
}

// 总是为true，并且返回原值
type AlwaysTrueReturn1 struct {
}

func (p *AlwaysTrueReturn1) Parse(param any) (bool, int) {
	return true, 1
}

// 通过规则id获取参数解析器
func GetParamParserByRule(ruleId int) ParamParser {
	if ruleId == 1 {
		return &AlwaysTrueReturn1{}
	} else if ruleId == 2 {
		return &AlwaysTrueReturn1{}
	}
	panic(fmt.Sprintf("no param parser found for ruleId %d", ruleId))
}
