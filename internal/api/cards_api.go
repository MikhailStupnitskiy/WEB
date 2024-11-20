package api

import (
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary Получить все карты
// @Description Возвращает список всех карт.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} schemas.GetAllCardsResponse "List of cards retrieved successfully"
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/cards [get]
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

	curr_move, err := a.repo.GetCurrMove()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var id int
	var cards_cnt int
	if len(curr_move) == 0 {
		id = 0
		cards_cnt = 0
	} else {
		id = curr_move[0].ID
		curr_cards := []int{}
		curr_cards, err = a.repo.GetCardsIDsByMoveID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cards_cnt = len(curr_cards)
	}
	response := schemas.GetAllCardsResponse{ID: id, Count: cards_cnt, Cards: cards}
	c.JSON(http.StatusOK, response)
}

// @Summary Получить карту по ID
// @Description Получить информацию о карте по ее ID
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Card ID"
// @Success 200 {object} schemas.GetCardResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/card/{ID} [get]
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

// @Summary Создать карту
// @Description Создать карту со свойствами
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body schemas.CreateCardRequest true "Card data"
// @Success 201 {object} schemas.CreateCardResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/card [post]
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
	c.JSON(http.StatusCreated, "Card was created")
}

// @Summary Удалить карту по ID
// @Description Удаляет карту по ее ID
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Card ID"
// @Success 200 {object} schemas.DeleteCardResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/card/{ID} [delete]
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

// @Summary Обновить карту по ID
// @Description Обновить карту по ее ID с параметрами
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Card ID"
// @Param body body schemas.UpdateCardRequest true "Update card data"
// @Success 200 {object} schemas.UpdateCardResponse
// @Failure 400 {object} schemas.ResponseMessage "Invalid request body"
// @Failure 500 {object} schemas.ResponseMessage "Internal server error"
// @Router /api/card/{ID} [put]
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
	c.JSON(http.StatusOK, "Card was updated")
}

// @Summary Добавить карту в ход
// @Description Этот эндпойнт позволяет добавить карту в ход по ее ID
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Card ID"
// @Param request query schemas.AddCardToMoveRequest true "AddCardToMoveRequest"
// @Success 200 {object} schemas.AddCardToMoveResponse "Card added successfully"
// @Failure 400 {object} schemas.ResponseMessage "Bad Request"
// @Failure 500 {object} schemas.ResponseMessage "Internal Server Error"
// @Router /api/card_to_move/{ID} [post]
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

// @Summary Поменять картинку по ID
// @Description Поменять картинку используя ее ID
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "Card ID"
// @Param image formData file true "File"
// @Success 200 {object} schemas.ResponseMessage "Picture was changed sucessfully"
// @Router /api/card/change_pic/{ID} [post]
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
	c.JSON(http.StatusOK, "Card Pic was updated")
}
