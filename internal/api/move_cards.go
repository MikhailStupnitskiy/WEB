package api

import (
	"Evolution/internal/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (a *Application) UpdateFoodMoveCardMeal(c *gin.Context) {
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
