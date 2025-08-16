package handlers

import (
	"net/http"
	"strconv"

	"github.com/Motionists/mt-order/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DishHandler struct {
	db *gorm.DB
}

func NewDishHandler(db *gorm.DB) *DishHandler {
	return &DishHandler{db: db}
}

func (h *DishHandler) GetDishesByMerchant(c *gin.Context) {
	merchantID, err := strconv.Atoi(c.Param("merchantId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	var dishes []models.Dish
	if err := h.db.Where("merchant_id = ? AND status = ?", merchantID, "available").Preload("Merchant").Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dishes": dishes})
}

func (h *DishHandler) GetDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID"})
		return
	}

	var dish models.Dish
	if err := h.db.Preload("Merchant").First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dish": dish})
}
