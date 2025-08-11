package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Motionists/mt-order/internal/models"
)

func (h *Handler) ListMerchants(c *gin.Context) {
	var ms []models.Merchant
	if err := h.db.Limit(50).Find(&ms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, ms)
}

func (h *Handler) GetMerchant(c *gin.Context) {
	id := c.Param("id")
	var m models.Merchant
	if err := h.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) ListDishes(c *gin.Context) {
	mid, _ := strconv.Atoi(c.Param("id"))
	var ds []models.Dish
	if err := h.db.Where("merchant_id = ?", mid).Find(&ds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, ds)
}