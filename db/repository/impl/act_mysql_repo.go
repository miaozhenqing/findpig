package impl

import (
	"errors"
	"findpig/db/model"
	"findpig/db/repository/inter"
	"fmt"
	"gorm.io/gorm"
)

type ActGormRepo struct {
	db *gorm.DB
}

func NewActGormRepo(db *gorm.DB) inter.ActRepo {
	return &ActGormRepo{db: db}
}

func buildTableName(mainId int) string {
	return fmt.Sprintf("act_%d", mainId)
}

func (a ActGormRepo) SelectInReturnList(userId int64, mainId int, subIds []int) ([]*model.ActEntity, error) {
	if len(subIds) == 0 {
		return nil, errors.New("subIds cannot be empty")
	}
	var acts []*model.ActEntity
	err := a.db.Table(buildTableName(mainId)).Where("user_id = ? and main_id = ? and sub_id in ?", userId, mainId, subIds).Find(&acts).Error
	if err != nil {
		return nil, err
	}
	return acts, nil
}

func (a *ActGormRepo) SelectByMainIdMap(userId int64, mainId int, subIds []int) (map[int][]*model.ActEntity, error) {
	if len(subIds) == 0 {
		return nil, errors.New("subIds cannot be empty")
	}

	// 1. 查询原始数据
	var acts []*model.ActEntity
	err := a.db.Table(buildTableName(mainId)).
		Where("user_id = ? AND main_id = ? AND sub_id IN ?", userId, mainId, subIds).
		Find(&acts).Error
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// 2. 转换为 map[mainId][]*model.ActEntity
	result := make(map[int][]*model.ActEntity)
	for _, act := range acts {
		result[act.MainID] = append(result[act.MainID], act)
	}

	return result, nil
}

func (a ActGormRepo) InsertByBatch(mainId int, acts []*model.ActEntity) (int64, error) {
	if len(acts) == 0 {
		return 0, nil
	}
	tx := a.db.Table(buildTableName(mainId)).CreateInBatches(acts, len(acts))
	return tx.RowsAffected, tx.Error
}

func (a ActGormRepo) Insert(mainId int, act *model.ActEntity) error {
	tx := a.db.Table(buildTableName(mainId)).Create(act)
	return tx.Error
}

func (a ActGormRepo) CreateTableIfNotExist(tableName string) error {
	if a.db.Migrator().HasTable(tableName) {
		return nil
	}
	var err = a.db.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
            user_id     BIGINT NOT NULL COMMENT '用户ID',
            main_id     INT NOT NULL COMMENT '主活动ID',
            sub_id      INT NOT NULL COMMENT '子活动ID',
            progress    JSON COMMENT '进度数据',
            state       TINYINT NOT NULL DEFAULT 1 COMMENT '状态',
            update_time BIGINT NOT NULL COMMENT '进度更新时间戳',
            create_time BIGINT NOT NULL COMMENT '数据创建时间戳',
            modify_time BIGINT NOT NULL COMMENT '数据修改时间戳',
            INDEX idx_user_id (user_id),
            INDEX idx_main_id (main_id),
            INDEX idx_sub_id (sub_id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
    `, tableName)).Error

	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}
	return nil
}
