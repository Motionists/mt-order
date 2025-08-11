package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Motionists/mt-order/internal/models"
)

type addCartReq struct {
	MerchantID uint `json:"merchant_id" binding:"required"`
	DishID     uint `json:"dish_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,min=1"`
}

func (h *Handler) GetCart(c *gin.Context) {
	var items []models.CartItem
	if err := h.db.Where("user_id = ?", uid(c)).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) AddCartItem(c *gin.Context) {
	var req addCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := models.CartItem{
		UserID:     uid(c),
		MerchantID: req.MerchantID,
		DishID:     req.DishID,
		Quantity:   req.Quantity,
	}
	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateCartItem(c *gin.Context) {
	type reqT struct{ Quantity int `json:"quantity" binding:"required,min=1"` }
	var req reqT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	if err := h.db.Model(&models.CartItem{}).
		Where("id = ? AND user_id = ?", id, uid(c)).
		Update("quantity", req.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Where("id = ? AND user_id = ?", id, uid(c)).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}