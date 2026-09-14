package api

import (
	"alloyDetermination/internal/app/handler"
	"alloyDetermination/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	// Подгружаем шаблоны и статику
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// 3 GET запроса по заданию
	r.GET("/feed", h.GetFeed)
	r.GET("/draft", h.GetDraft)
	r.GET("/catalog", h.GetCatalog)

	r.Run(":3000")
	log.Println("Server down")
}
