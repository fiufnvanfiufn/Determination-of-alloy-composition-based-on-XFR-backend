package repository

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	minio  *minio.Client
	bucket string
}

type Settings struct {
	PostgresDSN    string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
}

func New(s *Settings) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(s.PostgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mc, err := minio.New(s.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.MinioAccessKey, s.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:     db,
		minio:  mc,
		bucket: s.MinioBucket,
	}, nil
}
