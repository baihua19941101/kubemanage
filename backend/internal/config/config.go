package config

import (
	"os"
	"strconv"
)

type Config struct {
	ListenAddr          string
	MySQLDSN            string
	RedisAddr           string
	RedisPass           string
	RedisDB             int
	K8sAdapterMode      string
	SecretKey           string
	AuthJWTSecret       string
	AuthAccessTTL       int
	AuthRefreshTTL      int
	AuthCompatStageKeep int
}

func Load() Config {
	listenAddr := getenv("KM_LISTEN_ADDR", ":8080")
	mysqlDSN := getenv("KM_MYSQL_DSN", "root:123456@tcp(localhost:3306)/kubemanage?charset=utf8mb4&parseTime=True&loc=Local")
	redisAddr := getenv("KM_REDIS_ADDR", "localhost:6379")
	redisPass := getenv("KM_REDIS_PASS", "")
	redisDB := getenvInt("KM_REDIS_DB", 0)
	k8sAdapterMode := getenv("KM_K8S_ADAPTER_MODE", "live")
	secretKey := getenv("KM_SECRET_KEY", "")
	authJWTSecret := getenv("KM_AUTH_JWT_SECRET", "km-dev-jwt-secret")
	authAccessTTL := getenvInt("KM_AUTH_ACCESS_TTL_SECONDS", 3600)
	authRefreshTTL := getenvInt("KM_AUTH_REFRESH_TTL_SECONDS", 604800)
	authCompatKeep := getenvInt("KM_AUTH_COMPAT_STAGE_KEEP", 1)

	return Config{
		ListenAddr:          listenAddr,
		MySQLDSN:            mysqlDSN,
		RedisAddr:           redisAddr,
		RedisPass:           redisPass,
		RedisDB:             redisDB,
		K8sAdapterMode:      k8sAdapterMode,
		SecretKey:           secretKey,
		AuthJWTSecret:       authJWTSecret,
		AuthAccessTTL:       authAccessTTL,
		AuthRefreshTTL:      authRefreshTTL,
		AuthCompatStageKeep: authCompatKeep,
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
