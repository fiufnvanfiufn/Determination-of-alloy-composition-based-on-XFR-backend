package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"alloyDetermination/internal/app/config"
	"alloyDetermination/internal/app/dsn"
	"alloyDetermination/internal/app/handler"
	"alloyDetermination/internal/app/redis"
	"alloyDetermination/internal/app/repository"
)

// @title Alloy Determination API
// @version 1.0
// @description API для определения состава сплавов
// @host localhost:3000
// @BasePath /api

func main() {
	_ = godotenv.Load()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient, err := redis.New(context.Background(), redis.Config{
		Host:        os.Getenv("REDIS_HOST"),
		Port:        redisPort,
		Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          redisDB,
		DialTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
	})
	if err != nil {
		logrus.Fatalf("cant connect to redis: %v", err)
	}
	defer redisClient.Close()

	rep, err := repository.New(&repository.Settings{
		PostgresDSN:    dsn.FromEnv(),
		MinioEndpoint:  os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket:    os.Getenv("MINIO_BUCKET_NAME"),
	})
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret"
	}

	hand := handler.NewHandler(rep, jwtSecret, redisClient)

	router := gin.Default()
	hand.RegisterStatic(router)
	hand.RegisterHandler(router)

	addr := conf.ServiceHost + ":" + strconv.Itoa(conf.ServicePort)
	logrus.Infof("server starting at %s", addr)
	if err := router.Run(addr); err != nil {
		logrus.Fatal(err)
	}
}
