package api

import (
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// @Summary Получить все заявки на ходы с параметрами
// @Description Получить список ходов с возможностью фильтрации по статусу и датам
// @Tags moves
// @Accept json
// @Produce json
// @Param status query string false "Статус хода"
// @Param is_status query string false "Наличие статуса"
// @Success 200 {object} schemas.GetAllMovesWithParamsResponse
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move [get]
// @Security BearerAuth
func (a *Application) GetAllMovesWithParams(c *gin.Context) {
	var request schemas.GetAllMovesWithParamsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.FromDate.IsZero() {
		request.FromDate = time.Date(2000, time.January, 1, 0, 0, 0, 396641000, time.UTC)
	}
	if request.ToDate.IsZero() {
		request.ToDate = time.Now()
	}
	if request.Status == 3 {
		c.JSON(http.StatusNotFound, "Moves deleted")
		return
	}
	moves, err := a.repo.GetAllMovesWithFilters(request.Status, request.HavingStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := schemas.GetAllMovesWithParamsResponse{Moves: moves}
	c.JSON(http.StatusOK, response)
}

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
	response := schemas.GetMoveResponse{Move: result, MoveCards: CardsInMove}
	c.JSON(http.StatusOK, response)
}

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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Fields was updated")
}

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

func (a *Application) FormMove(c *gin.Context) {
	var request schemas.FormMoveRequest
	id := c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.FormMove(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Move was Formed")
}

func (a *Application) FinishMove(c *gin.Context) {
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
	err := a.repo.FinishMove(id, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Move was Finished")
}
