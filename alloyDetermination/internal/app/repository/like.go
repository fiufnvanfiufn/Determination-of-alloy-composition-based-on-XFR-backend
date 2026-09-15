package repository

import "alloyDetermination/internal/app/ds"

// Количество лайков для одной услуги.
func (r *Repository) GetLikesCount(alloyID uint) (int64, error) {
	var count int64
	err := r.db.
		Model(&ds.Like{}).
		Where("alloy_id = ?", alloyID).
		Count(&count).Error
	return count, err
}

// Количество лайков для набора услуг — одним запросом, чтобы не гонять БД в цикле.
func (r *Repository) GetLikesCountsForAlloys(ids []uint) (map[uint]int64, error) {
	if len(ids) == 0 {
		return map[uint]int64{}, nil
	}

	type row struct {
		AlloyID uint
		Cnt     int64
	}
	var rows []row

	err := r.db.
		Model(&ds.Like{}).
		Select("alloy_id, COUNT(*) as cnt").
		Where("alloy_id IN ?", ids).
		Group("alloy_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[uint]int64, len(rows))
	for _, r := range rows {
		res[r.AlloyID] = r.Cnt
	}
	return res, nil
}

// (Опционально) список лайков конкретной услуги, если нужен на странице "подробнее".
func (r *Repository) GetLikesForAlloy(alloyID uint) ([]ds.Like, error) {
	var likes []ds.Like
	err := r.db.Where("alloy_id = ?", alloyID).Find(&likes).Error
	return likes, err
}
