package repository

import (
	"alloyDetermination/internal/app/ds"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrAlloyNotFound = errors.New("alloy not found or not published")

func (r *Repository) GetLikesCount(alloyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("alloy_id = ?", alloyID).Count(&count).Error
	return count, err
}

func (r *Repository) GetLikesCountsForAlloys(ids []uint) (map[uint]int64, error) {
	result := make(map[uint]int64, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	type row struct {
		AlloyID uint
		Cnt     int64
	}
	var rows []row
	err := r.db.Model(&ds.Like{}).
		Select("alloy_id, COUNT(*) AS cnt").
		Where("alloy_id IN ?", ids).
		Group("alloy_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, x := range rows {
		result[x.AlloyID] = x.Cnt
	}
	return result, nil
}

func (r *Repository) SetLike(alloyID, userID uint, likeValue int) (int, error) {
	// 1. Валидация likeValue: только 0 или 1
	if likeValue != 0 && likeValue != 1 {
		return 0, errors.New("like value must be 0 or 1")
	}

	// 2. Проверка, что сплав существует и опубликован
	var alloy ds.Alloy
	err := r.db.
		Where("alloy_id = ? AND alloy_status = ?", alloyID, "опубликован").
		First(&alloy).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrAlloyNotFound
	}
	if err != nil {
		return 0, err
	}

	// 3. Поставить или снять лайк
	if likeValue == 1 {
		l := ds.Like{UserID: userID, AlloyID: alloyID}
		if err := r.db.
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&l).Error; err != nil {
			return 0, err
		}
	} else {
		if err := r.db.
			Where("user_id = ? AND alloy_id = ?", userID, alloyID).
			Delete(&ds.Like{}).Error; err != nil {
			return 0, err
		}
	}

	var count int64
	if err := r.db.
		Model(&ds.Like{}).
		Where("alloy_id = ?", alloyID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) IsLikedByUser(userID, alloyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).
		Where("user_id = ? AND alloy_id = ?", userID, alloyID).
		Count(&count).Error
	return count > 0, err
}
