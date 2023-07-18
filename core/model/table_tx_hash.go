package model

import "time"

type TTxHash struct {
	ID        uint64    `gorm:"primaryKey;column:id;type:int unsigned;not null"`
	Height    uint64    `gorm:"column:height;type:bigint unsigned;not null"`
	TxHash    string    `gorm:"column:tx_hash;type:varchar(66)"`
	BlockHash string    `gorm:"column:block_hash;type:varchar(66)"`
	Status    uint      `gorm:"column:status;type:int unsigned;not null"`
	OrderType uint      `gorm:"column:order_type;type:int unsigned;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

//TableName
func (m *TTxHash) TableName() string {
	return "t_tx_hash"
}
