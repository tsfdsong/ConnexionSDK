package model

import (
	"time"
)

type TNftDepositRecord struct {
	ID              uint64    `gorm:"primaryKey;column:id;type:int;not null"`
	AppID           int       `gorm:"index:unq_order_id;column:app_id;type:int"`
	AppOrderID      string    `gorm:"index:unq_order_id;column:app_order_id;type:varchar(80);null;default:''"`
	UID             uint64    `gorm:"index:idx_t_email_binds_uid;column:uid;type:bigint unsigned;not null"`
	Account         string    `gorm:"column:account;type:varchar(100)"`
	GameAssetName   string    `gorm:"column:game_asset_name;type:varchar(45)`
	EquipmentID     string    `gorm:"column:equipment_id;type:varchar(45)"`
	ContractAddress string    `gorm:"column:contract_address;type:varchar(45);not null;default:''"`
	DepositAddress  string    `gorm:"column:deposit_address;type:varchar(45)"`
	TargetAddress   string    `gorm:"column:trease_address;type:varchar(45)"`
	TokenID         string    `gorm:"column:token_id;type:varchar(80);default:''"`
	TxHash          string    `gorm:"column:tx_hash;type:varchar(66)"`
	OrderStatus     int       `gorm:"column:order_status;type:int;default:0"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Timestamp       int64     `gorm:"column:timestamp;type:bigint unsigned;not null"`
	Nonce           string    `gorm:"column:nonce;type:varchar(80);not null;default:''"`
	Height          uint64    `gorm:"column:height;type:bigint unsigned;not null"`
}

// TableName
func (m *TNftDepositRecord) TableName() string {
	return "t_nft_deposit_record"
}
