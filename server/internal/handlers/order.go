package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Motionists/mt-order/internal/models"
	"gorm.io/gorm"
)

type createOrderReq struct {
	MerchantID uint `json:"merchant_id" binding:"required"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req createOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := uid(c)

	// 拿购物车生成订单
	var items []models.CartItem
	if err := h.db.Where("user_id = ? AND merchant_id = ?", uid, req.MerchantID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}
	var dishes []models.Dish
	var total int64
	dishIDs := make([]uint, 0, len(items))
	for _, it := range items {
		dishIDs = append(dishIDs, it.DishID)
	}
	if err := h.db.Where("id IN ?", dishIDs).Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	m := map[uint]models.Dish{}
	for _, d := range dishes {
		m[d.ID] = d
	}
	for _, it := range items {
		total += m[it.DishID].Price * int64(it.Quantity)
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		o := models.Order{
			UserID:      uid,
			MerchantID:  req.MerchantID,
			Status:      "created",
			TotalAmount: total,
			PayAmount:   total,
		}
		if err := tx.Create(&o).Error; err != nil {
			return err
		}
		for _, it := range items {
			d := m[it.DishID]
			oi := models.OrderItem{
				OrderID:  o.ID,
				DishID:   d.ID,
				Name:     d.Name,
				Price:    d.Price,
				Quantity: it.Quantity,
			}
			if err := tx.Create(&oi).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("user_id = ? AND merchant_id = ?", uid, req.MerchantID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		c.JSON(http.StatusOK, gin.H{"order_id": o.ID, "status": o.Status, "pay_amount": o.PayAmount})
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	var o models.Order
	if err := h.db.First(&o, "id = ? AND user_id = ?", id, uid(c)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListOrders(c *gin.Context) {
	var os []models.Order
	if err := h.db.Where("user_id = ?", uid(c)).Order("id desc").Limit(50).Find(&os).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, os)
}

// 模拟支付：直接把订单状态改为 paid，然后异步推进状态
func (h *Handler) PayOrderMock(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Model(&models.Order{}).
		Where("id = ? AND user_id = ?", id, uid(c)).
		Update("status", "paid").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	// 模拟状态推进（生产请用消息队列）
	go func(orderID int) {
		time.Sleep(2 * time.Second)
		h.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "accepted")
		time.Sleep(2 * time.Second)
		h.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "making")
		time.Sleep(3 * time.Second)
		h.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "delivering")
		time.Sleep(3 * time.Second)
		h.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "done")
	}(id)
	c.JSON(http.StatusOK, gin.H{"status": "paid"})
}

// SSE 推送订单状态
func (h *Handler) OrderStatusStream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	uidv := uid(c)
	orderID := c.Param("id")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			var o models.Order
			if err := h.db.First(&o, "id = ? AND user_id = ?", orderID, uidv).Error; err != nil {
				return
			}
			fmt.Fprintf(c.Writer, "event: status\n")
			fmt.Fprintf(c.Writer, "data: {\"status\":\"%s\"}\n\n", o.Status)
			c.Writer.Flush()
			if o.Status == "done" || o.Status == "canceled" {
				return
			}
		}
	}
}