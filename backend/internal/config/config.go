package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigFile          string
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
	LogLevel            string
	LogFormat           string
	LogOutput           string
}

type fileConfig struct {
	Server struct {
		ListenAddr string `yaml:"listen_addr"`
	} `yaml:"server"`
	Database struct {
		MySQLDSN string `yaml:"mysql_dsn"`
	} `yaml:"database"`
	Redis struct {
		Addr string `yaml:"addr"`
		Pass string `yaml:"pass"`
		DB   *int   `yaml:"db"`
	} `yaml:"redis"`
	K8s struct {
		AdapterMode string `yaml:"adapter_mode"`
	} `yaml:"k8s"`
	Security struct {
		SecretKey string `yaml:"secret_key"`
	} `yaml:"security"`
	Auth struct {
		JWTSecret       string `yaml:"jwt_secret"`
		AccessTTL       *int   `yaml:"access_ttl_seconds"`
		RefreshTTL      *int   `yaml:"refresh_ttl_seconds"`
		CompatStageKeep *int   `yaml:"compat_stage_keep"`
	} `yaml:"auth"`
	Log struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
	} `yaml:"log"`
}

func Load() (Config, error) {
	cfg := defaultConfig()
	cfg.ConfigFile = getenv("KM_CONFIG_FILE", "configs/config.yaml")

	if err := applyFileConfig(&cfg, cfg.ConfigFile); err != nil {
		return Config{}, err
	}
	applyEnvOverrides(&cfg)

	return cfg, nil
}

func defaultConfig() Config {
	return Config{
		ConfigFile:          "configs/config.yaml",
		ListenAddr:          ":8080",
		MySQLDSN:            "root:123456@tcp(localhost:3306)/kubemanage?charset=utf8mb4&parseTime=True&loc=Local",
		RedisAddr:           "localhost:6379",
		RedisPass:           "",
		RedisDB:             0,
		K8sAdapterMode:      "live",
		SecretKey:           "",
		AuthJWTSecret:       "km-dev-jwt-secret",
		AuthAccessTTL:       3600,
		AuthRefreshTTL:      604800,
		AuthCompatStageKeep: 1,
		LogLevel:            "info",
		LogFormat:           "text",
		LogOutput:           "stdout",
	}
}

func applyFileConfig(cfg *Config, filePath string) error {
	path := strings.TrimSpace(filePath)
	if path == "" {
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read config file failed (%s): %w", path, err)
	}

	var fc fileConfig
	if err := yaml.Unmarshal(raw, &fc); err != nil {
		return fmt.Errorf("parse config file failed (%s): %w", path, err)
	}

	if value := strings.TrimSpace(fc.Server.ListenAddr); value != "" {
		cfg.ListenAddr = value
	}
	if value := strings.TrimSpace(fc.Database.MySQLDSN); value != "" {
		cfg.MySQLDSN = value
	}
	if value := strings.TrimSpace(fc.Redis.Addr); value != "" {
		cfg.RedisAddr = value
	}
	if value := strings.TrimSpace(fc.Redis.Pass); value != "" {
		cfg.RedisPass = value
	}
	if fc.Redis.DB != nil {
		cfg.RedisDB = *fc.Redis.DB
	}
	if value := strings.TrimSpace(fc.K8s.AdapterMode); value != "" {
		cfg.K8sAdapterMode = value
	}
	if value := strings.TrimSpace(fc.Security.SecretKey); value != "" {
		cfg.SecretKey = value
	}
	if value := strings.TrimSpace(fc.Auth.JWTSecret); value != "" {
		cfg.AuthJWTSecret = value
	}
	if fc.Auth.AccessTTL != nil {
		cfg.AuthAccessTTL = *fc.Auth.AccessTTL
	}
	if fc.Auth.RefreshTTL != nil {
		cfg.AuthRefreshTTL = *fc.Auth.RefreshTTL
	}
	if fc.Auth.CompatStageKeep != nil {
		cfg.AuthCompatStageKeep = *fc.Auth.CompatStageKeep
	}
	if value := strings.TrimSpace(fc.Log.Level); value != "" {
		cfg.LogLevel = value
	}
	if value := strings.TrimSpace(fc.Log.Format); value != "" {
		cfg.LogFormat = value
	}
	if value := strings.TrimSpace(fc.Log.Output); value != "" {
		cfg.LogOutput = value
	}
	return nil
}

func applyEnvOverrides(cfg *Config) {
	cfg.ListenAddr = getenv("KM_LISTEN_ADDR", cfg.ListenAddr)
	cfg.MySQLDSN = getenv("KM_MYSQL_DSN", cfg.MySQLDSN)
	cfg.RedisAddr = getenv("KM_REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisPass = getenv("KM_REDIS_PASS", cfg.RedisPass)
	cfg.RedisDB = getenvInt("KM_REDIS_DB", cfg.RedisDB)
	cfg.K8sAdapterMode = getenv("KM_K8S_ADAPTER_MODE", cfg.K8sAdapterMode)
	cfg.SecretKey = getenv("KM_SECRET_KEY", cfg.SecretKey)
	cfg.AuthJWTSecret = getenv("KM_AUTH_JWT_SECRET", cfg.AuthJWTSecret)
	cfg.AuthAccessTTL = getenvInt("KM_AUTH_ACCESS_TTL_SECONDS", cfg.AuthAccessTTL)
	cfg.AuthRefreshTTL = getenvInt("KM_AUTH_REFRESH_TTL_SECONDS", cfg.AuthRefreshTTL)
	cfg.AuthCompatStageKeep = getenvInt("KM_AUTH_COMPAT_STAGE_KEEP", cfg.AuthCompatStageKeep)
	cfg.LogLevel = getenv("KM_LOG_LEVEL", cfg.LogLevel)
	cfg.LogFormat = getenv("KM_LOG_FORMAT", cfg.LogFormat)
	cfg.LogOutput = getenv("KM_LOG_OUTPUT", cfg.LogOutput)
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
