package api

import (
	"Evolution/docs"
	"Evolution/internal/app/config"
	"Evolution/internal/app/ds"
	"Evolution/internal/app/dsn"
	"Evolution/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (a *Application) Run() {
	log.Println("Server start up")

	docs.SwaggerInfo.Title = "Evolution"
	docs.SwaggerInfo.Description = "API SERVER"
	docs.SwaggerInfo.Version = "1.1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	r := gin.Default()

	r.GET("/api/cards", a.OptionalAuthMiddleware(), a.GetAllCards)
	r.GET("/api/card/:ID", a.GetCard)
	r.GET("/api/card_by_name/:name", a.GetCardByName)

	r.POST("/api/card", a.RoleMiddleware(ds.Users{IsModerator: true}), a.CreateCard)
	r.DELETE("/api/card/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.DeleteCard)
	r.PUT("/api/card/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.UpdateCard)
	r.POST("/api/card_to_move/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.AddCardToMove)
	r.POST("api/card/change_pic/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.ChangePic)

	r.GET("/api/move", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetAllMovesWithParams)
	r.GET("/api/move/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.GetMove)
	r.PUT("/api/move/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.UpdateFieldsMove)
	r.DELETE("/api/move/:ID", a.RoleMiddleware(ds.Users{IsModerator: false}, ds.Users{IsModerator: true}), a.DeleteMove)
	r.PUT("/api/move/form/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.FormMove)
	r.PUT("/api/move/finish/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.FinishMove)

	r.DELETE("/api/move_cards/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.DeleteCardFromMove)
	r.PUT("/api/move_cards/:ID", a.RoleMiddleware(ds.Users{IsModerator: true}), a.UpdateFoodMoveCard)

	r.POST("/api/register_user", a.RegisterUser)
	r.POST("/api/login_user", a.LoginUser)
	r.POST("/api/logout", a.LogoutUser)

	r.GET("/protected", a.RoleMiddleware(ds.Users{IsModerator: true}), func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)
		c.JSON(http.StatusOK, gin.H{"message": "Пользователь авторизован с правами модератора", "userID": userID})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
