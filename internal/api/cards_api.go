package api

import (
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (a *Application) GetAllCards(c *gin.Context) {
	var request schemas.GetAllCardsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cards, err := a.repo.GetAllCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cards_cnt := len(cards)
	curr_move, err := a.repo.GetCurrMove()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var id int
	if len(curr_move) == 0 {
		id = 0
	} else {
		id = curr_move[0].ID
	}
	response := schemas.GetAllCardsResponse{ID: id, Count: cards_cnt, Cards: cards}
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetCard(c *gin.Context) {
	var request schemas.GetCardRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := a.repo.GetCardByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.GetCardResponse{Card: card}
	c.JSON(http.StatusOK, response)
}

func (a *Application) CreateCard(c *gin.Context) {
	var request schemas.CreateCardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.CreateCard(request.Cards)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Meal was created")
}

func (a *Application) DeleteCard(c *gin.Context) {
	var request schemas.GetCardRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.DeleteCardByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Card was deleted")
}

func (a *Application) UpdateCard(c *gin.Context) {
	var request schemas.UpdateCardRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request.Cards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request.Cards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateCardByID(request.ID, request.Cards)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal was updated")
}

func (a *Application) AddCardToMove(c *gin.Context) {
	var request schemas.AddCardToMoveRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	new_move, err := a.repo.CreateMove()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	new_move_id := new_move.ID
	card_id, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(new_move_id)
	err = a.repo.AddToMove(new_move_id, card_id)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Card was added")
}

func (a *Application) ChangePic(c *gin.Context) {
	var request schemas.ChangePicRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(request.ID, request.ImageUrl)
	err := a.repo.ChangePicByID(request.ID, request.ImageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Meal Pic was updated")
}
