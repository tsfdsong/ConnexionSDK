package model

import "time"

type TFtWithdrawRecord struct {
	ID              uint64    `gorm:"primaryKey;column:id;type:int;not null"`
	UID             uint64    `gorm:"column:uid;type:bigint unsigned;not null"`
	GameCoinName    string    `gorm:"column:game_coin_name;type:varchar(45)`
	AppID           int       `gorm:"index:unq_order_id;column:app_id;type:int;not null;default:0"`
	ChainID         uint64    `gorm:"column:chain_id;type:int;not null;default:0"`
	AppOrderID      string    `gorm:"index:unq_order_id;column:app_order_id;type:varchar(45);not null;default:''"`
	Account         string    `gorm:"column:account;type:varchar(100);not null"`
	ContractAddress string    `gorm:"column:contract_address;type:varchar(45)"`
	WithdrawAddress string    `gorm:"column:withdraw_address;type:varchar(45)"`
	Amount          string    `gorm:"column:amount;type:varchar(80)"`
	TxHash          string    `gorm:"column:tx_hash;type:varchar(66)"`
	OrderStatus     int       `gorm:"column:order_status;type:int"`
	Signature       string    `gorm:"column:signature;type:varchar(150);not null;default:''"`
	SignatureHash   string    `gorm:"column:signature_hash;type:varchar(45);not null;default:''"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Timestamp       int64     `gorm:"column:timestamp;type:bigint unsigned;not null"`
	Nonce           string    `gorm:"column:nonce;type:varchar(80);not null;default:''"` // sdk nonce
	RiskStatus      int       `gorm:"column:risk_status;type:int(11);not null;default:0"`
	RiskReviewer    string    `gorm:"column:risk_reviewer;type:varchar(45)"`
	Height          uint64    `gorm:"column:height;type:bigint unsigned;not null"`
}

// TableName
func (m *TFtWithdrawRecord) TableName() string {
	return "t_ft_withdraw_record"
}
