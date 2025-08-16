package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Motionists/mt-order/server/internal/config"
	"github.com/Motionists/mt-order/server/internal/handlers"
	"github.com/Motionists/mt-order/server/internal/middleware"
	"github.com/Motionists/mt-order/server/internal/services"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// 初始化服务
	authService := services.NewAuthService(db, cfg)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	merchantHandler := handlers.NewMerchantHandler(db)
	dishHandler := handlers.NewDishHandler(db)
	cartHandler := handlers.NewCartHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	// 公开路由
	api := r.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 商家相关
		merchants := api.Group("/merchants")
		{
			merchants.GET("", merchantHandler.GetMerchants)
			merchants.GET("/:id", merchantHandler.GetMerchant)
			merchants.GET("/:merchantId/dishes", dishHandler.GetDishesByMerchant)
		}

		// 菜品相关
		dishes := api.Group("/dishes")
		{
			dishes.GET("/:id", dishHandler.GetDish)
		}
	}

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// 购物车相关
		cart := protected.Group("/cart")
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("", cartHandler.AddToCart)
			cart.PUT("/:id", cartHandler.UpdateCartItem)
			cart.DELETE("/:id", cartHandler.RemoveFromCart)
		}

		// 订单相关
		orders := protected.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrder)
		}
	}

	return r
}
