package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Motionists/mt-order/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

type CreateOrderRequest struct {
	MerchantID uint   `json:"merchant_id" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Remark     string `json:"remark"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetFloat64("userID")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItems []models.CartItem
	if err := h.db.Where("user_id = ? AND merchant_id = ?", uint(userID), req.MerchantID).Preload("Dish").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items in cart for this merchant"})
		return
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Dish.Price * float64(item.Quantity)
	}

	orderNumber := fmt.Sprintf("ORD-%s", uuid.New().String()[:8])

	order := models.Order{
		UserID:      uint(userID),
		MerchantID:  req.MerchantID,
		OrderNumber: orderNumber,
		TotalAmount: totalAmount,
		Status:      "pending",
		Address:     req.Address,
		Phone:       req.Phone,
		Remark:      req.Remark,
	}

	tx := h.db.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:  order.ID,
			DishID:   item.DishID,
			Quantity: item.Quantity,
			Price:    item.Dish.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Where("user_id = ? AND merchant_id = ?", uint(userID), req.MerchantID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetFloat64("userID")

	var orders []models.Order
	if err := h.db.Where("user_id = ?", uint(userID)).Preload("Merchant").Preload("OrderItems.Dish").Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetFloat64("userID")
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := h.db.Where("id = ? AND user_id = ?", orderID, uint(userID)).Preload("Merchant").Preload("OrderItems.Dish").First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}
