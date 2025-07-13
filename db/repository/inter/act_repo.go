package inter

import (
	"findpig/db/model"
)

type ActRepo interface {
	SelectInReturnList(userId int64, mainId int, subIds []int) ([]*model.ActEntity, error)

	SelectByMainIdMap(userId int64, mainId int, subIds []int) (map[int][]*model.ActEntity, error)

	InsertByBatch(mainId int, acts []*model.ActEntity) (int64, error)

	Insert(mainId int, act *model.ActEntity) error

	CreateTableIfNotExist(tableName string) error
}
