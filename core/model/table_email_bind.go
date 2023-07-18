package model

import "time"

type TEmailBind struct {
	ID         uint64    `gorm:"primaryKey;column:id;type:int unsigned;not null"`
	AppID      uint      `gorm:"index:game_addr;column:app_id;type:int unsigned;not null"`
	UID        uint64    `gorm:"index:idx_t_email_binds_uid;column:uid;type:bigint unsigned;not null"`
	Address    string    `gorm:"index:game_addr;column:address;type:varchar(64);not null"`
	ZKSAddress string    `gorm:"column:zks_address;type:varchar(64);not null"`
	Account    string    `gorm:"column:account;type:varchar(100);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// TableName
func (m *TEmailBind) TableName() string {
	return "t_email_bind"
}
