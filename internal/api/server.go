package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.Static("/css", "./resources")

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		cardsquery := c.Query("CardItem")
		cards := CardFunc()

		var filteredCards []Card
		for _, card := range cards {
			if strings.Contains(strings.ToLower(card.TitleRu), strings.ToLower(cardsquery)) || strings.Contains(strings.ToLower(card.TitleEn), strings.ToLower(cardsquery)) {
				filteredCards = append(filteredCards, card)
			}
		}

		if len(filteredCards) == 0 {
			filteredCards = cards
		}

		c.HTML(http.StatusOK, "cards.html", gin.H{
			"title":         "Колода карт",
			"cards_data":    CardFunc(),
			"filteredCards": filteredCards,
			"searchQuery":   cardsquery,
			"cart_ID":       MoveFunc()[0].ID,
			"card_count":    len(CardFunc()),
		})
	})

	r.GET("/card_detailed/:id", func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из URL
		index, err := strconv.Atoi(id)

		if err != nil || index < 0 || index > len(CardFunc()) {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		c.HTML(http.StatusOK, "card_detailed.html", gin.H{
			"title":     "Информация о карте",
			"card_data": CardFunc()[index-1],
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

		c.HTML(http.StatusOK, "move.html", gin.H{
			"title":      "Заявка на ход",
			"cards_data": MoveFunc(),
			"move_id":    id,
			"options":    options,
			"stage":      stage,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
