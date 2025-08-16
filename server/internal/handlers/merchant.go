package handlers

import (
	"net/http"
	"strconv"

	"github.com/Motionists/mt-order/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MerchantHandler struct {
	db *gorm.DB
}

func NewMerchantHandler(db *gorm.DB) *MerchantHandler {
	return &MerchantHandler{db: db}
}

func (h *MerchantHandler) GetMerchants(c *gin.Context) {
	var merchants []models.Merchant
	if err := h.db.Where("status = ?", "active").Find(&merchants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"merchants": merchants})
}

func (h *MerchantHandler) GetMerchant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	var merchant models.Merchant
	if err := h.db.First(&merchant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"merchant": merchant})
}
