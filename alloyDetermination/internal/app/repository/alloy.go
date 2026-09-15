package repository

import (
	"database/sql"
	"errors"
	"time"

	"alloyDetermination/internal/app/ds"
)

// Список опубликованных услуг (для плитки/каталога), с поиском по названию.
// Чистый ORM.
func (r *Repository) GetPublishedAlloys(search string) ([]ds.Alloy, error) {
	var alloys []ds.Alloy
	q := r.db.Where("status = ?", "опубликован")
	if search != "" {
		q = q.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := q.Order("date_create DESC").Find(&alloys).Error; err != nil {
		return nil, err
	}
	return alloys, nil
}

// Один образец по ID — только не удалённый. Используем "курсор" (raw SQL + Row()).
func (r *Repository) GetAlloyByID(id int) (*ds.Alloy, error) {
	query := `
		SELECT id, name, description, status, img_url, video_url,
		       energy_kev, intensity_cps, date_create, date_formed, creator_id
		FROM alloys
		WHERE id = $1 AND status <> 'удален'`

	row := r.db.Raw(query, id).Row()

	a := &ds.Alloy{}
	err := row.Scan(
		&a.ID, &a.Name, &a.Description, &a.Status,
		&a.ImgURL, &a.VideoURL,
		&a.EnergyKev, &a.IntensityCps,
		&a.DateCreate, &a.DateFormed, &a.CreatorID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

// Черновик конкретного пользователя (или nil, если нет).
func (r *Repository) GetDraftByCreator(creatorID uint) (*ds.Alloy, error) {
	var a ds.Alloy
	err := r.db.
		Where("creator_id = ? AND status = ?", creatorID, "черновик").
		First(&a).Error
	if err != nil {
		if errors.Is(err, gormErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// Создание черновика. Чистый ORM.
func (r *Repository) CreateDraft(creatorID uint, name, imgURL, videoURL string) (*ds.Alloy, error) {
	a := &ds.Alloy{
		Name:       name,
		ImgURL:     imgURL,
		VideoURL:   videoURL,
		Status:     "черновик",
		DateCreate: time.Now(),
		CreatorID:  creatorID,
	}
	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

// Публикация: заполняем описание и два "предметных" поля, меняем статус.
// Чистый ORM.
func (r *Repository) PublishDraft(id uint, description string, energyKev, intensityCps float64) error {
	return r.db.Model(&ds.Alloy{}).
		Where("id = ? AND status = ?", id, "черновик").
		Updates(map[string]interface{}{
			"description":   description,
			"energy_kev":    energyKev,
			"intensity_cps": intensityCps,
			"status":        "опубликован",
			"date_formed":   time.Now(),
		}).Error
}

// Логическое удаление — БЕЗ ORM, через raw SQL ("курсор").
func (r *Repository) SoftDeleteAlloy(id uint) error {
	row := r.db.Raw(`UPDATE alloys SET status = 'удален' WHERE id = $1`, id).Row()
	defer row.Close()
	return nil
}
