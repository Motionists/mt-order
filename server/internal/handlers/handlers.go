package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Motionists/mt-order/internal/config"
	"github.com/Motionists/mt-order/internal/models"
)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Client
	cfg config.Config
}

func New(db *gorm.DB, rdb *redis.Client, cfg config.Config) *Handler {
	return &Handler{db: db, rdb: rdb, cfg: cfg}
}

// Health example
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pong": time.Now().Unix()})
}

// Helpers
func uid(c *gin.Context) uint {
	v, _ := c.Get("uid")
	return v.(uint)
}

var ctx = context.Background()

// AutoMigrate helper (call once in init script or admin endpoint)
func (h *Handler) AutoMigrate() error {
	return h.db.AutoMigrate(
		&models.User{},
		&models.Merchant{},
		&models.Dish{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
}