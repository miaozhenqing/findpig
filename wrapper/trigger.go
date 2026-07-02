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
	TableCache        *cache.Cache
	ActRepo           inter.ActRepo
	MainToEntityCache map[int][]*model.ActEntity
	InitSelectFlag    bool
	ContextWrappers   []*ContextWrapper[any]
}

func (dm *DataManager) GetActivity(userId int64, createIfNotExists bool, wrapper *SubWrapper) *model.ActEntity {
	dm.TryToInitSelect(userId)
	var res *model.ActEntity
	if dm.MainToEntityCache != nil {
		entities := dm.MainToEntityCache[wrapper.MainId]
		if entities != nil {
			for _, entity := range entities {
				if entity.SubID == wrapper.Id && entity.UserID == userId {
					res = entity
					break
				}
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
			dm.MainToEntityCache[wrapper.MainId] = append(dm.MainToEntityCache[wrapper.MainId], entity)
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
	dirtyMap := make(map[int][]*model.ActEntity)
	for mainId, entities := range dm.MainToEntityCache {
		var dirtyEntities []*model.ActEntity
		for _, entity := range entities {
			if entity.Dirty {
				entity.Post()
				entity.Dirty = false
				dirtyEntities = append(dirtyEntities, entity)
			}
		}
		// 只有存在 dirty 数据时才加入最终 map
		if len(dirtyEntities) > 0 {
			dirtyMap[mainId] = dirtyEntities
		}
	}
	if len(dirtyMap) == 0 {
		return
	}
	_, err := dm.ActRepo.UpdateByBatch(dirtyMap)
	if err != nil {
		panic(err)
	}
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
		dbEntities, err := dm.ActRepo.SelectByMainIdMap(userId, mainId, subIds)
		if err != nil {
			panic(err)
		}
		existSubIds := make(map[int]bool)
		for _, entity := range dbEntities {
			for _, actEntity := range entity {
				actEntity.Init()
				existSubIds[actEntity.SubID] = true
			}
		}
		dm.MainToEntityCache = dbEntities
		for _, subWrapper := range subWrappers {
			if !existSubIds[subWrapper.Id] {
				notExistSubIds[subWrapper.MainId] = append(notExistSubIds[subWrapper.MainId], subWrapper)
			}
		}
	}
	var mainToEntityMap = make(map[int][]*model.ActEntity)
	for mainId, subWrappers := range notExistSubIds {
		entities := make([]*model.ActEntity, 0, len(subWrappers))
		for _, subWrapper := range subWrappers {
			ruleIds := util.Map(subWrapper.Rules, func(t *RuleWrapper) int {
				return t.Id
			})
			entity := model.NewActEntity(userId, mainId, subWrapper.Id, ruleIds)
			entities = append(entities, entity)
		}
		mainToEntityMap[mainId] = entities
	}
	if len(mainToEntityMap) > 0 {
		_, err := dm.ActRepo.AddByBatch(mainToEntityMap)
		if err != nil {
			panic(err)
		}
		// 合并新数据到缓存
		for mainId, newEntities := range mainToEntityMap {
			if existing, ok := dm.MainToEntityCache[mainId]; ok {
				dm.MainToEntityCache[mainId] = append(existing, newEntities...)
			} else {
				dm.MainToEntityCache[mainId] = newEntities
			}
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
