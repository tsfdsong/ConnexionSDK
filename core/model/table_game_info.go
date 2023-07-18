package model

import "time"

type TGameInfo struct {
	ID            uint64    `gorm:"primaryKey;column:id;type:bigint unsigned;not null"`
	AppID         uint8     `gorm:"index:idx_t_game_infos_game_id;column:app_id;type:tinyint unsigned;not null"`
	BaseServerURL string    `gorm:"index:idx_t_game_infos_base_server_url;column:base_server_url;type:varchar(64);not null"`
	AppSecret     string    `gorm:"unique;column:app_secret;type:varchar(64);not null"`
	AppKey        string    `gorm:"unique;column:app_key;type:varchar(64);not null"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

//TableName
func (m *TGameInfo) TableName() string {
	return "t_game_info"
}
