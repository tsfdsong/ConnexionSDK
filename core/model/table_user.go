package model

import "time"

type TUser struct {
	ID             uint64    `gorm:"primaryKey;column:id;type:int;not null"`
	AppID          uint8     `gorm:"index:game_email;column:app_id;type:tinyint unsigned;not null"`
	Account        string    `gorm:"index:game_email;column:account;type:varchar(100);not null"`
	Name           string    `gorm:"column:name;type:varchar(32);not null"`
	Image          string    `gorm:"column:image;type:varchar(256);not null"`
	Bio            string    `gorm:"column:bio;type:varchar(256);not null"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	WithdrawSwitch int8      `gorm:"column:withdraw_switch;type:tinyint;default:0"`
}

// TableName
func (m *TUser) TableName() string {
	return "t_user"
}
