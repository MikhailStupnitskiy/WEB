package api

import (
	"Evolution/internal/app/config"
	"Evolution/internal/app/ds"
	"Evolution/internal/app/dsn"
	"Evolution/internal/app/repository"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Card struct {
	ID          int
	Multiplier  string
	TitleEn     string
	TitleRu     string
	ImageUrl    string
	Description string
}

type Move_M2M struct {
	ID          int
	CardElement []Card
	Food        int
}

type Move struct {
	ID     int
	Player string
	Stage  string
	Info   Move_M2M
}

type Option struct {
	Value string
	Label string
}

type InfoForMove struct {
	Card ds.Cards
	Food int
}

func MoveFunc() []Move {
	MoveArr := []Move{
		{ID: 1, Player: "Игрок 1", Stage: "Кормление", Info: Move_M2M{1, CardFunc(), 1}},
	}
	return MoveArr
}

func CardFunc() []Card {
	CardsInfo := []Card{
		{1, "+1", "HIGH BODY WEIGHT",
			"БОЛЬШОЙ", "http://127.0.0.1:9000/test/image%201.jpg", "Данное животное может быть съедено только “БОЛЬШИМ” хищником",
		},
		{2, "", "RUNNING",
			"БЫСТРОЕ", "http://127.0.0.1:9000/test/running.jpg", "Когда ваше животное с этим свойством атаковано, бросьте кубик",
		},
		{3, "", "POISONOUS",
			"ЯДОВИТОЕ", "http://127.0.0.1:9000/test/poisonous.jpg", "Хищник, съевший это животное, в фазу вымирания погибает",
		},
		{4, " ", "BURROWING",
			"НОРНОЕ", "http://127.0.0.1:9000/test/burrowing.jpg", "Когда животное НАКОРМЛЕНО, оно не может быть атаковано хищником",
		},
		{5, "+1", "HIGH BODY WEIGHT",
			"БОЛЬШОЙ", "http://127.0.0.1:9000/test/image%201.jpg", "Данное животное может быть съедено только “БОЛЬШИМ” хищником",
		}}
	return CardsInfo
}

type Application struct {
	repo   *repository.Repository
	config *config.Config
}

func (a *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	var err error

	r.Static("/css", "./resources")

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		cardsquery := c.Query("CardItem")

		var filteredCards []ds.Cards

		if cardsquery == "" {
			filteredCards, err = a.repo.GetAllCards()
			if err != nil {
				log.Println("unable to get all cards")
				c.Error(err)
				return
			}
		} else {
			filteredCards, err = a.repo.GetCardByInfo(cardsquery)
			if err != nil {
				log.Println("unable to get cards by info")
				filteredCards = []ds.Cards{}
			}
		}

		var move_len int
		var move_ID int
		curr_move, err := a.repo.GetCurrMove()
		if err != nil {
			log.Println("unable to get current move")
		}
		if len(curr_move) == 0 {
			move_len = 0
			move_ID = 0

		} else {
			cards_in_req, err := a.repo.GetCardsIDsByMoveID(curr_move[0].ID)
			if err != nil {
				log.Println("unable to get cards ids by move")
			}
			move_len = len(cards_in_req)
			move_ID = curr_move[0].ID
		}

		c.HTML(http.StatusOK, "cards.html", gin.H{
			"title":         "Колода карт",
			"cards_data":    CardFunc(),
			"filteredCards": filteredCards,
			"searchQuery":   cardsquery,
			"move_ID":       move_ID,
			"card_count":    move_len,
		})
	})

	r.POST("/home", func(c *gin.Context) {

		id := c.PostForm("add")
		card_ID, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Printf("cant transform ind", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		move_wrk, err := a.repo.GetCurrMove()
		var move_ID int
		if len(move_wrk) == 0 {
			new_move, err := a.repo.CreateMove()
			if err != nil {
				log.Println("unable to create move")
			}
			move_ID = new_move[0].ID
		} else {
			move_ID = move_wrk[0].ID
		}

		a.repo.AddToMove(move_ID, card_ID)

		c.Redirect(301, "/home")

	})

	r.GET("/card_detailed/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL
		index, err := strconv.Atoi(id)

		if err != nil { // если не получилось
			log.Printf("cant get card by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		card, err := a.repo.GetCardByID(index)
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		c.HTML(http.StatusOK, "card_detailed.html", gin.H{
			"title":     "Информация о карте",
			"card_data": card,
		})
	})

	r.GET("/move/:id", func(c *gin.Context) {
		options := []Option{
			{Value: "1", Label: "Игрок 1"},
			{Value: "2", Label: "Игрок 2"},
			{Value: "3", Label: "Игрок 3"},
			{Value: "4", Label: "Игрок 4"},
		}

		stage := []Option{
			{Value: "1", Label: "Создание"},
			{Value: "2", Label: "Кормление"},
		}

		id := c.Param("id")
		index, err := strconv.Atoi(id)
		if err != nil { // если не получилось
			log.Printf("cant get move by id %v", err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		move_status, err := a.repo.GetMoveStatusByID(index)
		if err != nil {
			log.Printf("cant get move by id %v", err)
		}
		if move_status == 3 {
			c.Redirect(301, "/home")
		}

		CardsIDs, err := a.repo.GetCardsIDsByMoveID(index)
		if err != nil {
			log.Println("unable to get cardsIDsByMoveID")
			c.Error(err)
			return
		}

		//CardsInMove := []ds.Cards{}
		CardsInMove := []InfoForMove{}
		var curr_card InfoForMove
		for _, v := range CardsIDs {
			card_temp, err := a.repo.GetCardsByID(v.CardID)
			if err != nil {
				c.Error(err)
				return
			}
			food, err := a.repo.GetFoodByCardID(index, v.CardID)
			curr_card.Card = card_temp[0]
			curr_card.Food = food
			CardsInMove = append(CardsInMove, curr_card)
			log.Println(len(CardsInMove))
		}

		c.HTML(http.StatusOK, "move.html", gin.H{
			"title":      "Заявка на ход",
			"cards_data": CardsInMove,
			"move_id":    index,
			"options":    options,
			"stage":      stage,
		})
	})

	r.POST("/move/:id", func(c *gin.Context) {

		id := c.Param("id")
		index, err := strconv.Atoi(id)
		if err != nil { // если не получилось
			log.Printf("cant get cart by id %v", err)
			c.Error(err)
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}
		a.repo.DeleteMove(index)
		c.Redirect(301, "/home")

	})

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
