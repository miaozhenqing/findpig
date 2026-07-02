package model

import (
	"encoding/json"
	"findpig/util"
)

type ActEntity struct {
	//数据库记录字段
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	MainID     int    `json:"main_id"`
	SubID      int    `json:"sub_id"`
	Progress   string `json:"progress"`
	State      int    `json:"state"`
	UpdateTime int64  `json:"update_time"`
	CreateTime int64  `json:"create_time"`
	ModifyTime int64  `json:"modify_time"`
	//不记录到数据库
	ActProgresses []*ActProgress `gorm:"-"`
	//不记录到数据库
	Dirty bool `gorm:"-"`
}

type ActProgress struct {
	R int `json:"r"`
	P any `json:"p"`
	T int `json:"t"`
}

func (act *ActEntity) Init() {
	if act.Progress != "" {
		err := json.Unmarshal([]byte(act.Progress), &act.ActProgresses)
		if err != nil {
			panic(err)
		}
	}
}
func (act *ActEntity) Post() {
	if act.ActProgresses != nil {
		marshal, _ := json.Marshal(act.ActProgresses)
		act.Progress = string(marshal)
	}
}

func NewActEntity(userId int64, mainId int, subId int, ruleIds []int) *ActEntity {
	pStr, p := newActProgress(ruleIds)
	entity := &ActEntity{
		UserID:        userId,
		SubID:         subId,
		MainID:        mainId,
		Progress:      pStr,
		CreateTime:    util.CurrentTime().UnixMilli(),
		ActProgresses: p,
	}
	return entity
}

func newActProgress(ruleIds []int) (string, []*ActProgress) {
	p := make([]*ActProgress, 0, len(ruleIds))
	for _, id := range ruleIds {
		progress := &ActProgress{R: id}
		p = append(p, progress)
	}
	marshal, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(marshal), p
}

func (act *ActEntity) IsInProgress() bool {
	return act.State == 0
}
func (act *ActEntity) IsCompleted() bool {
	return act.State == 1
}
func (act *ActEntity) MarkCompleted() {
	act.State = 1
}
func (act *ActEntity) IsReward() bool {
	return act.State == 2
}
func (act *ActEntity) MarkReward() {
	act.State = 2
}
func (act *ActEntity) Flush(flushType int) {
	//todo
}
