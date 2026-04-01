package infra

import "time"

type ClusterConnectionRecord struct {
	ID                uint   `gorm:"primaryKey"`
	Name              string `gorm:"size:128;uniqueIndex;not null"`
	Mode              string `gorm:"size:32;not null"`
	APIServer         string `gorm:"size:512"`
	KubeconfigContent string `gorm:"type:longtext"`
	BearerToken       string `gorm:"type:longtext"`
	CACert            string `gorm:"type:longtext"`
	SkipTLSVerify     bool   `gorm:"not null;default:false"`
	IsDefault         bool   `gorm:"not null;default:false"`
	Status            string `gorm:"size:32;not null;default:'unknown'"`
	LastCheckedAt     *time.Time
	LastError         string `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (ClusterConnectionRecord) TableName() string {
	return "cluster_connections"
}
