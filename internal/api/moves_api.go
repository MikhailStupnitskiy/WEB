package api

import (
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary Получить все заявки на ходы с параметрами
// @Description Получить список ходов с возможностью фильтрации по статусу и датам
// @Tags moves
// @Accept json
// @Produce json
// @Param status query string false "Статус хода" Enum(1,2,3)
// @Param is_status query string false "Наличие статуса" Enum(true,false)
// @Param from_date query string false "Дата от" format(date)
// @Param to_date query string false "Дата до" format(date)
// @Success 200 {object} schemas.GetAllMovesWithParamsResponse
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move [get]
// @Security BearerAuth
func (a *Application) GetAllMovesWithParams(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}
	var request schemas.GetAllMovesWithParamsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	moves, err := a.repo.GetAllMovesWithFilters(request.Status, userID.(float64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.GetAllMovesWithParamsResponse{Moves: moves}
	c.JSON(http.StatusOK, response)
}

// @Summary Получить заявку на ход по ID
// @Description Получить детализированную информацию о ходе, включая карты, участвующие в ходе, и их статус
// @Tags moves
// @Accept json
// @Produce json
// @Param ID path string true "ID хода"
// @Success 200 {object} schemas.GetMoveResponse
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 404 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move/{ID} [get]
// @Security BearerAuth
func (a *Application) GetMove(c *gin.Context) {
	var request schemas.GetMoveRequest
	request.ID = c.Param("ID")
	id_int, err := strconv.Atoi(request.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("error was there")
		return
	}
	move, err := a.repo.GetMoveByID(id_int)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	card_ids_in_move, err := a.repo.GetCardsIDsByMoveID(move.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	CardsInMove := []schemas.InfoForMove{}
	var curr_card schemas.InfoForMove
	for _, v := range card_ids_in_move {
		v_string := strconv.Itoa(v)
		card_to_append, err := a.repo.GetCardByID(v_string)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		food, err := a.repo.GetFoodByCardID(id_int, v)
		curr_card.Card = card_to_append
		curr_card.Food = food
		CardsInMove = append(CardsInMove, curr_card)
		log.Println(len(CardsInMove))
	}
	result := map[string]interface{}{
		"ID":           move.ID, // Строка
		"Status":       move.Status,
		"DateCreate":   move.DateCreate,
		"DateUpdate":   move.DateUpdate,
		"DateFinish":   move.DateFinish,
		"CreatorID":    move.CreatorID,
		"ModeratorID":  move.ModeratorID,
		"CreatorLogin": move.Creator.Login,
		"Moderator":    move.Moderator,
		"Player":       move.Player,
		"Stage":        move.Stage,
		"Cube":         move.Cube,
	}
	response := schemas.GetMoveResponse{Move: result, Status: move.Status, MoveCards: CardsInMove}
	c.JSON(http.StatusOK, response)
}

// @Summary Обновить поля хода
// @Description Обновить информацию о ходе, включая поля игрока и этапа
// @Tags moves
// @Accept json
// @Produce json
// @Param ID path string true "ID хода"
// @Param player body string true "Имя игрока"
// @Param stage body string true "Этап хода"
// @Success 200 {string} string "Fields was updated"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move/{ID} [put]
// @Security BearerAuth
func (a *Application) UpdateFieldsMove(c *gin.Context) {
	var request schemas.UpdateFieldsMoveRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateFieldsMove(request)
	if request.Stage == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Fields was updated")
}

// @Summary Удалить ход
// @Description Удалить заявку на ход по ID
// @Tags moves
// @Accept json
// @Produce json
// @Param ID path string true "ID хода"
// @Success 200 {string} string "Move was deleted"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move/{ID} [delete]
// @Security BearerAuth
func (a *Application) DeleteMove(c *gin.Context) {
	var request schemas.DeleteMoveRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = a.repo.DeleteMove(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Move was deleted")
}

// @Summary Сформировать ход
// @Description Сформировать ход по его ID
// @Tags moves
// @Accept json
// @Produce json
// @Param ID path string true "ID хода"
// @Success 200 {string} string "Move was Formed"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move/form/{ID} [put]
// @Security BearerAuth
func (a *Application) FormMove(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}
	var request schemas.FormMoveRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.FormMove(id, userID.(float64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Move was Formed")
}

// @Summary Завершить ход
// @Description Завершить ход по его ID, обновив его статус
// @Tags moves
// @Accept json
// @Produce json
// @Param ID path string true "ID хода"
// @Param status body int true "Статус завершенного хода" Enum(1,2,3)
// @Success 200 {string} string "Move was Finished"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move/finish/{ID} [put]
// @Security BearerAuth
func (a *Application) FinishMove(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}
	var request schemas.FinishMoveRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.FinishMove(id, request.Status, userID.(float64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Move was Finished")
}
