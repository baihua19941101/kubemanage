package infra

import (
	"context"
	"fmt"
	"time"

	"kubeManage/backend/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewStore(cfg config.Config) (*Store, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db failed: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql failed: %w", err)
	}

	if err := db.AutoMigrate(&ClusterConnectionRecord{}, &UserRecord{}, &RefreshTokenRecord{}, &AuthProviderRecord{}); err != nil {
		return nil, fmt.Errorf("auto migrate cluster connections failed: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis failed: %w", err)
	}

	return &Store{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func (s *Store) Close() error {
	var firstErr error

	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if s.DB != nil {
		sqlDB, err := s.DB.DB()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
		} else if err := sqlDB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
