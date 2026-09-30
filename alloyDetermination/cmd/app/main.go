package main

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"alloyDetermination/internal/app/config"
	"alloyDetermination/internal/app/dsn"
	"alloyDetermination/internal/app/handler"
	"alloyDetermination/internal/app/repository"
)

func main() {
	_ = godotenv.Load()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

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

	hand := handler.NewHandler(rep)

	router := gin.Default()
	hand.RegisterStatic(router)
	hand.RegisterHandler(router)

	addr := conf.ServiceHost + ":" + strconv.Itoa(conf.ServicePort)
	logrus.Infof("server starting at %s", addr)
	if err := router.Run(addr); err != nil {
		logrus.Fatal(err)
	}
}
