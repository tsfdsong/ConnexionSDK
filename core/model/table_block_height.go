package model

import "time"

type TBlockHeight struct {
	ID                 uint64    `gorm:"primaryKey;column:id;type:int unsigned;not null"`
	AppID              int       `gorm:"column:app_id;type:int;not null;default:0"`
	LatestParsedHeight uint64    `gorm:"index:idx_t_block_heights_latest_parsed_height;column:latest_parsed_height;type:bigint unsigned;not null"`
	Switch             int8      `gorm:"column:switch;type:tinyint;default:0"`
	CreatedAt          time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// TableName
func (m *TBlockHeight) TableName() string {
	return "t_block_height"
}
