package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Motionists/mt-order/internal/config"
	"github.com/Motionists/mt-order/internal/handlers"
	"github.com/Motionists/mt-order/internal/middleware"
)

func Register(r *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg config.Config) {
	h := handlers.New(db, rdb, cfg)

	api := r.Group("/api/v1")

	// auth
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)
	api.GET("/me", middleware.Auth(cfg.JWT.Secret), h.Me)

	// merchants & dishes
	api.GET("/merchants", h.ListMerchants)
	api.GET("/merchants/:id", h.GetMerchant)
	api.GET("/merchants/:id/dishes", h.ListDishes)

	// cart (login required)
	cart := api.Group("/cart").Use(middleware.Auth(cfg.JWT.Secret))
	cart.GET("", h.GetCart)
	cart.POST("/items", h.AddCartItem)
	cart.PATCH("/items/:id", h.UpdateCartItem)
	cart.DELETE("/items/:id", h.DeleteCartItem)

	// orders
	orders := api.Group("/orders").Use(middleware.Auth(cfg.JWT.Secret))
	orders.POST("", h.CreateOrder)
	orders.GET("", h.ListOrders)
	orders.GET("/:id", h.GetOrder)
	orders.GET("/:id/stream", h.OrderStatusStream) // SSE demo
	orders.POST("/:id/pay", h.PayOrderMock)
}