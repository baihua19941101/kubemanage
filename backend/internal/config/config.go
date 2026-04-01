package config

import (
	"os"
	"strconv"
)

type Config struct {
	ListenAddr     string
	MySQLDSN       string
	RedisAddr      string
	RedisPass      string
	RedisDB        int
	K8sAdapterMode string
	SecretKey      string
}

func Load() Config {
	listenAddr := getenv("KM_LISTEN_ADDR", ":8080")
	mysqlDSN := getenv("KM_MYSQL_DSN", "root:123456@tcp(localhost:3306)/kubemanage?charset=utf8mb4&parseTime=True&loc=Local")
	redisAddr := getenv("KM_REDIS_ADDR", "localhost:6379")
	redisPass := getenv("KM_REDIS_PASS", "")
	redisDB := getenvInt("KM_REDIS_DB", 0)
	k8sAdapterMode := getenv("KM_K8S_ADAPTER_MODE", "live")
	secretKey := getenv("KM_SECRET_KEY", "")

	return Config{
		ListenAddr:     listenAddr,
		MySQLDSN:       mysqlDSN,
		RedisAddr:      redisAddr,
		RedisPass:      redisPass,
		RedisDB:        redisDB,
		K8sAdapterMode: k8sAdapterMode,
		SecretKey:      secretKey,
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
