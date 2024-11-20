package api

import (
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Удалить карту из запроса на ход
// @Description Удаляет карту из запроса на ход по ID запроса и CardID
// @Tags moves_cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ID path string true "ID заявки"
// @Param body body schemas.DeleteCardFromMoveRequest true "Delete card from move"
// @Success 200 {string} string "Card was deleted from move"
// @Failure 400 {object} schemas.ResponseMessage
// @Failure 500 {object} schemas.ResponseMessage
// @Router /api/move_cards/{ID} [delete]
func (a *Application) DeleteCardFromMove(c *gin.Context) {
	var request schemas.DeleteCardFromMoveRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.DeleteCardFromMove(request.ID, request.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Сard was deleted from move")
}

// @Summary Обновить количество корма для карты в ходе
// @Description Обновляет количество корма для карты в заявке
// @Tags moves_cards
// @Accept json
// @Produce json
// @Param ID path string true "Move ID"
// @Param body body schemas.UpdateFoodMoveCardRequest true "Update amount of food in move"
// @Success 200 {string} string "Food was updated"
// @Failure 400 {object} schemas.ResponseMessage
// @Router /api/move_cards/{ID} [put]
func (a *Application) UpdateFoodMoveCard(c *gin.Context) {
	var request schemas.UpdateFoodMoveCardRequest
	request.ID = c.Param("ID")
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.repo.UpdateFoodMoveCard(request.ID, request.CardID, request.Food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Card food was changed in move")
}
