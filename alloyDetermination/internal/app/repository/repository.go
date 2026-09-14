package repository

import (
	"fmt"
	"strconv"
)

type Alloy struct {
	ID           int
	Title        string
	Description  string
	EnergyKev    float64
	SecondEnergyKev float64
	IntensityCps int
	Status       string
	ImageURL     string
	VideoURL     string
	Likes        []int
}

type Repository struct {
	data []Alloy
}

func NewRepository() (*Repository, error) {
	data := []Alloy{
		{
			ID:           1,
			Title:        "Оловянистая бронза (Скифы)",
			Description:  "tmp",
			EnergyKev:    8.04,
			IntensityCps: 7.21,
			Status:       "published",
			ImageURL:     "http://localhost:9000/media/514770_big.jpg",
			VideoURL:     "http://localhost:9000/media/bronze_scan.mp4",
			Likes:        []int{101, 102, 105},
		},
		{
			ID:           2,
			Title:        "Аурихалк / Латунь (Рим)",
			Description:  "tmp",
			EnergyKev:    8.63,
			IntensityCps: 6.51,
			Status:       "draft",
			ImageURL:     "http://localhost:9000/media/514770_big.jpg",
			VideoURL:     "http://localhost:9000/media/brass_scan.mp4",
			Likes:        []int{103},
		},
		{
			ID:           3,
			Title:        "Электрум (Лидия)",
			Description:  "tmp",
			EnergyKev:    9.71,
			IntensityCps: 2.62,
			Status:       "published",
			ImageURL:     "http://localhost:9000/media/514770_big.jpg",
			VideoURL:     "http://localhost:9000/media/electrum_scan.mp4",
			Likes:        []int{101, 109, 110, 115},
		},
	}
	return &Repository{data: data}, nil
}

func (r *Repository) GetAlloysCatalog(energyQuery string) ([]Alloy, error) {
	var result []Alloy
	for _, alloy := range r.data {
		if alloy.Status == "published" {
			if energyQuery != "" {
				energyStr := strconv.FormatFloat(alloy.EnergyKev, 'f', -1, 64)
				if energyStr == energyQuery {
					result = append(result, alloy)
				}
			} else {
				result = append(result, alloy)
			}
		}
	}
	return result, nil
}

func (r *Repository) GetAlloyDraft() (Alloy, error) {
	for _, alloy := range r.data {
		if alloy.Status == "draft" {
			return alloy, nil
		}
	}
	return Alloy{}, fmt.Errorf("черновик не найден")
}

func (r *Repository) GetAlloyFeed(id int, next bool) (Alloy, error) {
	if id == 0 {
		for _, alloy := range r.data {
			if alloy.Status == "published" {
				return alloy, nil
			}
		}
	}

	for i, alloy := range r.data {
		if alloy.ID == id {
			if next {
				for j := i + 1; j < len(r.data); j++ {
					if r.data[j].Status == "published" {
						return r.data[j], nil
					}
				}
				for _, first := range r.data {
					if first.Status == "published" {
						return first, nil
					}
				}
			}
			if alloy.Status == "published" {
				return alloy, nil
			}
		}
	}
	return Alloy{}, fmt.Errorf("сплав не найден")
}
