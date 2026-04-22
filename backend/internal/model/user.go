package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username" comment:"用户名"`
	PasswordHash string         `gorm:"size:255;not null" json:"-" comment:"密码哈希"`
	Email        string         `gorm:"size:128;not null;uniqueIndex" json:"email" comment:"邮箱"`
	AvatarURL    string         `gorm:"size:255" json:"avatar_url" comment:"头像URL"`
	Status       int8           `gorm:"size:1;default:1" json:"status" comment:"状态：1正常，0禁用"`
}

func (User) TableName() string {
	return "users"
}
