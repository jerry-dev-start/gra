package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"gra/internal/business"
	"gra/internal/router"
	"gra/internal/system"
	"gra/internal/system/user"
	"gra/pkg/config"
	"gra/pkg/database"
	"gra/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 2. 初始化日志
	logger.Init(&cfg.Log)

	// 3. 初始化数据库
	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	// 4. 自动迁移
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	// 5. 依赖注入 — 三行搞定，不管模块多少
	sysHandlers, sysSvc := system.Init(db)
	bizHandlers := business.Init(db, sysSvc)

	// 6. 启动 Gin（用 zap 替代默认日志）
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery())
	router.Setup(r, sysHandlers, bizHandlers)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Infof("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
