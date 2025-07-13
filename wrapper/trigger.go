package wrapper

import (
	"findpig/cache"
	"findpig/db/model"
	"findpig/db/repository/inter"
	"findpig/util"
	"fmt"
	"time"
)

// 入参
type ContextWrapper[T any] struct {
	UserId       int64
	SubWrapper   *SubWrapper
	TriggerRule  int
	TriggerParam T
	DataManager  *DataManager
	ParserResult []*ParserResult
	SuccessFlag  bool
}

func (cw ContextWrapper[T]) GetActivity() *model.ActEntity {
	return cw.DataManager.GetActivity(cw.UserId, true, cw.SubWrapper)
}

type DataManager struct {
	TableCache      *cache.Cache
	ActRepo         inter.ActRepo
	ActEntityCache  []*model.ActEntity
	InitSelectFlag  bool
	ContextWrappers []*ContextWrapper[any]
}

func (dm *DataManager) GetActivity(userId int64, createIfNotExists bool, wrapper *SubWrapper) *model.ActEntity {
	dm.TryToInitSelect(userId)
	var res *model.ActEntity
	if dm.ActEntityCache != nil {
		for _, entity := range dm.ActEntityCache {
			if entity.MainID == wrapper.MainId && entity.SubID == wrapper.Id && entity.UserID == userId {
				res = entity
				break
			}
		}

		if res == nil && createIfNotExists {
			ruleIds := util.Map(wrapper.Rules, func(t *RuleWrapper) int {
				return t.Id
			})
			entity := model.NewActEntity(userId, wrapper.Id, wrapper.MainId, ruleIds)
			err := dm.ActRepo.Insert(entity.MainID, entity)
			if err != nil {
				panic(err)
			}
			dm.ActEntityCache = append(dm.ActEntityCache, entity)
			res = entity
		}
	}
	if res != nil {
		res.Flush(wrapper.Flush.FlushType)
	}
	return res
}

func NewDataManager(tc *cache.Cache, actRepo inter.ActRepo) *DataManager {
	return &DataManager{
		TableCache: tc,
		ActRepo:    actRepo,
	}
}

func (dm *DataManager) UpdateDirtyActivity() {
	//for i, entity := range dm.ActEntityCache {
	//
	//}
}
func (dm *DataManager) TryToInitSelect(userId int64) {
	if dm.InitSelectFlag {
		return
	}
	dm.InitSelectFlag = true
	notExistSubIds := make(map[int][]SubWrapper)
	groupByMain := dm.groupByMain()
	for mainId, subWrappers := range groupByMain {
		dm.createTableIfNotExist(mainId)
		subIds := util.Map(subWrappers, func(t SubWrapper) int {
			return t.Id
		})
		dbEntities, err := dm.ActRepo.SelectInReturnList(userId, mainId, subIds)
		if err != nil {
			panic(err)
		}
		dm.ActEntityCache = dbEntities
		//找到数据库里面没有的
		existSubIds := make(map[int]bool)
		for _, entity := range dbEntities {
			existSubIds[entity.SubID] = true
		}
		for _, subWrapper := range subWrappers {
			if !existSubIds[subWrapper.Id] {
				notExistSubIds[subWrapper.MainId] = append(notExistSubIds[subWrapper.MainId], subWrapper)
			}
		}
	}
	for mainId, subWrappers := range notExistSubIds {
		entities := make([]*model.ActEntity, 0, len(subWrappers))
		for _, subWrapper := range subWrappers {
			ruleIds := util.Map(subWrapper.Rules, func(t *RuleWrapper) int {
				return t.Id
			})
			entity := model.NewActEntity(userId, mainId, subWrapper.Id, ruleIds)
			entities = append(entities, entity)
		}
		_, err := dm.ActRepo.InsertByBatch(mainId, entities)
		if err != nil {
			panic(err)
		}
		for _, entity := range entities {
			dm.ActEntityCache = append(dm.ActEntityCache, entity)
		}
	}
}
func (dm *DataManager) groupByMain() map[int][]SubWrapper {
	var mainToSubsMap = make(map[int][]SubWrapper)
	for _, c := range dm.ContextWrappers {
		mainId := c.SubWrapper.MainId
		mainToSubsMap[mainId] = append(mainToSubsMap[mainId], *c.SubWrapper)
	}
	return mainToSubsMap
}

func (dm *DataManager) createTableIfNotExist(mainId int) {
	_, b := dm.TableCache.Get(string(mainId))
	if b {
		return
	}
	tableName := fmt.Sprintf("%s_%d", "act", mainId)
	err := dm.ActRepo.CreateTableIfNotExist(tableName)
	if err != nil {
		panic(err)
	}
	dm.TableCache.Set(string(mainId), true, time.Hour.Milliseconds()*24*7)
}
