package api

import (
	"Evolution/internal/app/ds"
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

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
	cards := make([]ds.Cards, 0, len(card_ids_in_move))
	for _, v := range card_ids_in_move {
		v_string := strconv.Itoa(v)
		card_to_append, err := a.repo.GetCardByID(v_string)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cards = append(cards, card_to_append)
	}
	response := schemas.GetMoveResponse{Move: move, Count: len(card_ids_in_move), MoveCards: cards}
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
