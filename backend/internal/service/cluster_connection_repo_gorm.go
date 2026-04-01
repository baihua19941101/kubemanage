package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kubeManage/backend/internal/infra"

	"gorm.io/gorm"
)

type gormClusterConnectionRepo struct{ db *gorm.DB }

func NewGormClusterConnectionRepo(db *gorm.DB) ClusterConnectionRepository {
	return &gormClusterConnectionRepo{db: db}
}

func (r *gormClusterConnectionRepo) List(ctx context.Context) ([]infra.ClusterConnectionRecord, error) {
	var items []infra.ClusterConnectionRecord
	if err := r.db.WithContext(ctx).Order("id asc").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list cluster connections failed: %w", err)
	}
	return items, nil
}
func (r *gormClusterConnectionRepo) Create(ctx context.Context, record *infra.ClusterConnectionRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create cluster connection failed: %w", err)
	}
	return nil
}
func (r *gormClusterConnectionRepo) Get(ctx context.Context, id uint) (infra.ClusterConnectionRecord, error) {
	var item infra.ClusterConnectionRecord
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return infra.ClusterConnectionRecord{}, fmt.Errorf("cluster connection not found: %d", id)
		}
		return infra.ClusterConnectionRecord{}, fmt.Errorf("get cluster connection failed: %w", err)
	}
	return item, nil
}
func (r *gormClusterConnectionRepo) SetActive(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&infra.ClusterConnectionRecord{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return fmt.Errorf("clear active cluster connection failed: %w", err)
		}
		result := tx.Model(&infra.ClusterConnectionRecord{}).Where("id = ?", id).Updates(map[string]any{"is_default": true, "updated_at": time.Now()})
		if result.Error != nil {
			return fmt.Errorf("set active cluster connection failed: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("cluster connection not found: %d", id)
		}
		return nil
	})
}
func (r *gormClusterConnectionRepo) GetActive(ctx context.Context) (infra.ClusterConnectionRecord, error) {
	var item infra.ClusterConnectionRecord
	if err := r.db.WithContext(ctx).Where("is_default = ?", true).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return infra.ClusterConnectionRecord{}, errors.New("no active cluster connection")
		}
		return infra.ClusterConnectionRecord{}, fmt.Errorf("get active cluster connection failed: %w", err)
	}
	return item, nil
}
func (r *gormClusterConnectionRepo) UpdateStatus(ctx context.Context, id uint, status string, checkedAt time.Time, lastError string) error {
	result := r.db.WithContext(ctx).Model(&infra.ClusterConnectionRecord{}).Where("id = ?", id).Updates(map[string]any{"status": status, "last_checked_at": checkedAt, "last_error": lastError, "updated_at": time.Now()})
	if result.Error != nil {
		return fmt.Errorf("update cluster connection status failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("cluster connection not found: %d", id)
	}
	return nil
}
