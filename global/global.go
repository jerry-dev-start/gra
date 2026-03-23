package global

import (
	"gra/pkg/config"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb  *redis.Client
	Conf *config.Config
)
