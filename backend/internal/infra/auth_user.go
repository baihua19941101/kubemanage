package infra

import "time"

type UserRecord struct {
	ID                uint   `gorm:"primaryKey"`
	Username          string `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash      string `gorm:"size:255;not null"`
	Role              string `gorm:"size:32;not null"`
	AllowedNamespaces string `gorm:"size:1024;not null;default:''"`
	IsActive          bool   `gorm:"not null;default:true"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (UserRecord) TableName() string {
	return "users"
}

type RefreshTokenRecord struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index;not null"`
	TokenHash  string `gorm:"size:128;uniqueIndex;not null"`
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (RefreshTokenRecord) TableName() string {
	return "refresh_tokens"
}
