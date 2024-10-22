package api

import (
	"Evolution/internal/app/config"
	"Evolution/internal/app/dsn"
	"Evolution/internal/app/repository"

	"github.com/gin-gonic/gin"
	"log"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
}

func (a *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/api/cards", a.GetAllCards)
	r.GET("/api/card/:ID", a.GetCard)
	r.POST("/api/card", a.CreateCard)
	r.DELETE("/api/card/:ID", a.DeleteCard)
	r.PUT("/api/card/:ID", a.UpdateCard)
	r.POST("/api/card_to_move/:ID", a.AddCardToMove)
	r.POST("api/card/change_pic/:ID", a.ChangePic)

	r.GET("/api/move", a.GetAllMovesWithParams)
	r.GET("/api/move/:ID", a.GetMove)
	r.PUT("/api/move/:ID", a.UpdateFieldsMove)
	r.DELETE("/api/move/:ID", a.DeleteMove)
	r.PUT("/api/move/form/:ID", a.FormMove)
	r.PUT("/api/move/finish/:ID", a.FinishMove)

	r.DELETE("/api/move_cards/:ID", a.DeleteCardFromMove)
	r.PUT("/api/move_cards/:ID", a.UpdateFoodMoveCardMeal)

	r.POST("/api/registration", a.CreateUser)

	r.Static("/css", "./resources")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	return &app, nil
}
