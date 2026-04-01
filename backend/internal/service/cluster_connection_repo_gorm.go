package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kubeManage/backend/internal/infra"

	"gorm.io/gorm"
)

type gormClusterConnectionRepo struct {
	db    *gorm.DB
	codec *sensitiveCodec
}

func NewGormClusterConnectionRepo(db *gorm.DB, secretKey string) ClusterConnectionRepository {
	return &gormClusterConnectionRepo{
		db:    db,
		codec: newSensitiveCodec(secretKey),
	}
}

func (r *gormClusterConnectionRepo) List(ctx context.Context) ([]infra.ClusterConnectionRecord, error) {
	var items []infra.ClusterConnectionRecord
	if err := r.db.WithContext(ctx).Order("id asc").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list cluster connections failed: %w", err)
	}
	if err := r.decryptRecords(items); err != nil {
		return nil, err
	}
	return items, nil
}
func (r *gormClusterConnectionRepo) Create(ctx context.Context, record *infra.ClusterConnectionRecord) error {
	toCreate := *record
	if err := r.encryptRecord(&toCreate); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(&toCreate).Error; err != nil {
		return fmt.Errorf("create cluster connection failed: %w", err)
	}
	*record = toCreate
	if err := r.decryptRecord(record); err != nil {
		return err
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
	if err := r.decryptRecord(&item); err != nil {
		return infra.ClusterConnectionRecord{}, err
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
			return infra.ClusterConnectionRecord{}, ErrNoActiveClusterConnection
		}
		return infra.ClusterConnectionRecord{}, fmt.Errorf("get active cluster connection failed: %w", err)
	}
	if err := r.decryptRecord(&item); err != nil {
		return infra.ClusterConnectionRecord{}, err
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

func (r *gormClusterConnectionRepo) encryptRecord(item *infra.ClusterConnectionRecord) error {
	var err error
	item.KubeconfigContent, err = r.codec.Encrypt(item.KubeconfigContent)
	if err != nil {
		return fmt.Errorf("encrypt kubeconfig content failed: %w", err)
	}
	item.BearerToken, err = r.codec.Encrypt(item.BearerToken)
	if err != nil {
		return fmt.Errorf("encrypt bearer token failed: %w", err)
	}
	item.CACert, err = r.codec.Encrypt(item.CACert)
	if err != nil {
		return fmt.Errorf("encrypt ca cert failed: %w", err)
	}
	return nil
}

func (r *gormClusterConnectionRepo) decryptRecords(items []infra.ClusterConnectionRecord) error {
	for i := range items {
		if err := r.decryptRecord(&items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *gormClusterConnectionRepo) decryptRecord(item *infra.ClusterConnectionRecord) error {
	var err error
	item.KubeconfigContent, err = r.codec.Decrypt(item.KubeconfigContent)
	if err != nil {
		return fmt.Errorf("decrypt kubeconfig content failed: %w", err)
	}
	item.BearerToken, err = r.codec.Decrypt(item.BearerToken)
	if err != nil {
		return fmt.Errorf("decrypt bearer token failed: %w", err)
	}
	item.CACert, err = r.codec.Decrypt(item.CACert)
	if err != nil {
		return fmt.Errorf("decrypt ca cert failed: %w", err)
	}
	return nil
}
