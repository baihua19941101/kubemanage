package infra

import "time"

type AuthProviderRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64;not null;uniqueIndex"`
	Type      string `gorm:"size:32;not null"`
	IsEnabled bool   `gorm:"not null;default:true"`
	IsDefault bool   `gorm:"not null;default:false"`
	Config    string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (AuthProviderRecord) TableName() string {
	return "auth_providers"
}
