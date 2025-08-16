package handlers

import (
	"net/http"
	"strconv"

	"github.com/Motionists/mt-order/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

type AddToCartRequest struct {
	MerchantID uint `json:"merchant_id" binding:"required"`
	DishID     uint `json:"dish_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,min=1"`
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetFloat64("userID")

	var cartItems []models.CartItem
	if err := h.db.Where("user_id = ?", uint(userID)).Preload("Dish").Preload("Merchant").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": cartItems})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetFloat64("userID")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingItem models.CartItem
	err := h.db.Where("user_id = ? AND dish_id = ?", uint(userID), req.DishID).First(&existingItem).Error

	if err == nil {
		existingItem.Quantity += req.Quantity
		if err := h.db.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully"})
		return
	}

	cartItem := models.CartItem{
		UserID:     uint(userID),
		MerchantID: req.MerchantID,
		DishID:     req.DishID,
		Quantity:   req.Quantity,
	}

	if err := h.db.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetFloat64("userID")
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItem models.CartItem
	if err := h.db.Where("id = ? AND user_id = ?", itemID, uint(userID)).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	cartItem.Quantity = req.Quantity
	if err := h.db.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetFloat64("userID")
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := h.db.Where("id = ? AND user_id = ?", itemID, uint(userID)).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
