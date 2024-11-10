package redis

import (
	"fmt"
	"location-service/bin/config"
	"location-service/bin/pkg/utils"
	"strings"
)

type AppConfig struct {
	UseRedis bool
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RedisClusterConfig struct {
	Hosts    []string
	Password string
}

var (
	AppConfigData          AppConfig
	RedisConfigData        RedisConfig
	RedisClusterConfigData RedisClusterConfig
)

func LoadConfig() {

	AppConfigData = AppConfig{
		UseRedis: config.GetConfig().UseRedis,
	}

	redisDb := config.GetConfig().RedisDB
	redisHost := config.GetConfig().RedisHost
	redisPort := config.GetConfig().RedisPort
	redisPass := config.GetConfig().RedisPassword

	RedisConfigData = RedisConfig{
		Host:     fmt.Sprintf("%v", redisHost),
		Port:     fmt.Sprintf("%v", redisPort),
		Password: fmt.Sprintf("%v", redisPass),
		DB:       utils.ConvertInt(redisDb),
	}

	clusterHost := strings.Split(config.GetConfig().RedisClusterNode, ";")
	RedisClusterConfigData = RedisClusterConfig{
		Hosts:    clusterHost,
		Password: config.GetConfig().RedisClusterPassword,
	}
}
