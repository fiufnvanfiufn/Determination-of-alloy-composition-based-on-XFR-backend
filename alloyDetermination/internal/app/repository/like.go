package repository

import (
	"alloyDetermination/internal/app/ds"

	"gorm.io/gorm/clause"
)

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

func (r *Repository) SetLike(userID, alloyID uint, like bool) error {
	if like {
		l := ds.Like{UserID: userID, AlloyID: alloyID}
		return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&l).Error
	}
	return r.db.Where("user_id = ? AND alloy_id = ?", userID, alloyID).Delete(&ds.Like{}).Error
}
