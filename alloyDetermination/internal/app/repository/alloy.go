package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"alloyDetermination/internal/app/ds"

	"gorm.io/gorm"

	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (r *Repository) GetAlloysCatalog(maxEnergy float64) ([]ds.Alloy, error) {
	var alloys []ds.Alloy
	q := r.db.Where("alloy_status = ?", "опубликован")

	q = q.Where("alloy_energy_kev <= ?", maxEnergy) // ← без if

	if err := q.Order("alloy_id").Find(&alloys).Error; err != nil {
		return nil, err
	}
	return alloys, nil
}

func (r *Repository) GetAlloyByID(id uint) (*ds.Alloy, error) {
	var a ds.Alloy
	err := r.db.Where("alloy_id = ? AND alloy_status <> ?", id, "удален").First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) GetDraftByCreator(creatorID uint) (*ds.Alloy, error) {
	var a ds.Alloy
	err := r.db.Where("creator_id = ? AND alloy_status = ?", creatorID, "черновик").First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) CreateDraft(creatorID uint, name, imgURL, videoURL string) (*ds.Alloy, error) {
	existing, err := r.GetDraftByCreator(creatorID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	a := &ds.Alloy{
		Name:       name,
		Status:     "черновик",
		ImgURL:     imgURL,
		VideoURL:   videoURL,
		DateCreate: time.Now(),
		CreatorID:  creatorID,
	}
	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (r *Repository) PublishAlloy(id uint, description string, energyKev, intensityCps float64) error {
	updates := map[string]interface{}{
		"alloy_description":   description,
		"alloy_energy_kev":    energyKev,
		"alloy_intensity_cps": intensityCps,
		"alloy_status":        "опубликован",
		"alloy_formed_at":     sql.NullTime{Time: time.Now(), Valid: true},
	}
	res := r.db.Model(&ds.Alloy{}).Where("alloy_id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("alloy %d not found", id)
	}
	return nil
}

func (r *Repository) DeleteAlloySQL(id uint) error {
	query := "UPDATE alloys SET alloy_status = $1 WHERE alloy_id = $2"
	row := r.db.Raw(query, "удален", id).Row()

	var dummy string
	if err := row.Scan(&dummy); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("delete alloy error: %w", err)
	}
	return nil
}

func (r *Repository) GetAlloyFeed(id uint, next bool) (*ds.Alloy, error) {
	var a ds.Alloy
	q := r.db.Where("alloy_status = ?", "опубликован")

	if id == 0 {
		if err := q.Order("alloy_id").First(&a).Error; err != nil {
			return nil, err
		}
		return &a, nil
	}

	op := "<"
	order := "alloy_id DESC"
	if next {
		op = ">"
		order = "alloy_id ASC"
	}
	err := q.Where("alloy_id "+op+" ?", id).Order(order).First(&a).Error
	if err != nil {
		var fallback ds.Alloy
		ord := "alloy_id ASC"
		if !next {
			ord = "alloy_id DESC"
		}
		if err := r.db.Where("alloy_status = ?", "опубликован").
			Order(ord).First(&fallback).Error; err != nil {
			return nil, err
		}
		return &fallback, nil
	}
	return &a, nil
}

func (r *Repository) UploadFile(file *multipart.FileHeader, name string) error {
	logrus.Infof("UploadFile START: %s (%d bytes)", name, file.Size)
	start := time.Now()

	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = r.minio.PutObject(context.Background(), r.bucket, name, f, file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})

	logrus.Infof("UploadFile DONE: %s in %v, err=%v", name, time.Since(start), err)
	return err
}

// UpdateAlloyFiles — сохраняет имена файлов в БД.
func (r *Repository) UpdateAlloyFiles(id uint, imgName, videoName string) error {
	return r.db.Model(&ds.Alloy{}).
		Where("alloy_id = ?", id).
		Updates(map[string]interface{}{
			"alloy_image_url": imgName,
			"alloy_video_url": videoName,
		}).Error
}

func (r *Repository) SoftDeleteAlloy(id, userID uint) error {
	res := r.db.Model(&ds.Alloy{}).
		Where("alloy_id = ? AND creator_id = ?", id, userID).
		Update("alloy_status", "удален")
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("alloy %d not found or not yours", id)
	}
	return nil
}
