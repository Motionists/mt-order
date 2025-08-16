package main

import (
	"log"

	"github.com/Motionists/mt-order/server/internal/config"
	"github.com/Motionists/mt-order/server/internal/database"
	"github.com/Motionists/mt-order/server/internal/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db := database.Init(cfg)

	// 设置路由
	r := router.SetupRouter(cfg, db)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
