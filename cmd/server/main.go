package main

import (
	"fmt"
	"gra/global"
	"gra/pkg/redis"
	"log"

	"gra/internal/business"
	"gra/internal/router"
	"gra/internal/system"
	"gra/internal/system/menus"
	"gra/internal/system/user"
	"gra/pkg/config"
	"gra/pkg/database"
	"gra/pkg/id"
	"gra/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置通过vIp
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	global.Conf = cfg
	// 2. 初始化日志
	logger.Init(&cfg.Log)

	// 3. 初始化雪花算法
	if err := id.Init(cfg.Snowflake.MachineID); err != nil {
		log.Fatalf("init snowflake: %v", err)
	}

	// 4. 初始化数据库
	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	//初始化 Redis
	global.Rdb, err = redis.Init(&cfg.Redis)
	if err != nil {
		panic(fmt.Sprintf("Redis 连接失败: %v", err))
	}
	defer func() {
		if err := global.Rdb.Close(); err != nil {
			// 打印日志即可，无需 panic
			log.Printf("Redis 释放资源时报错: %v", err)
		}
	}()
	// 5. 自动迁移
	if err := db.AutoMigrate(&user.User{}, &menus.Menus{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	// 6. 依赖注入 — 三行搞定，不管模块多少
	sysHandlers, sysSvc := system.Init(db)
	bizHandlers := business.Init(db, sysSvc)

	// 7. 启动 Gin（用 zap 替代默认日志）
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
