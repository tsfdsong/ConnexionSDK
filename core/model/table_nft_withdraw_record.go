package model

import (
	"time"

	"gorm.io/datatypes"
)

type TNftWithdrawRecord struct {
	ID                uint64         `gorm:"primaryKey;column:id;type:int;not null"`
	AppID             int            `gorm:"index:unq_order_id;column:app_id;type:int;not null;default:0"`
	UID               uint64         `gorm:"index:idx_t_email_binds_uid;column:uid;type:bigint unsigned;not null"`
	AppOrderID        string         `gorm:"index:unq_order_id;column:app_order_id;type:varchar(45);not null;default:''"`
	EquipmentID       string         `gorm:"column:equipment_id;type:varchar(45);not null;default:0"`
	GameAssetName     string         `gorm:"column:game_asset_name;type:varchar(45)`
	Account           string         `gorm:"column:account;type:varchar(100);not null"`
	ContractAddress   string         `gorm:"column:contract_address;type:varchar(45)"`
	WithdrawAddress   string         `gorm:"column:withdraw_address;type:varchar(45)"`
	TreaseAddress     string         `gorm:"column:trease_address;type:varchar(45)"`
	GameMinterAddress string         `gorm:"column:minter_address;type:varchar(45)"`
	TokenID           string         `gorm:"column:token_id;type:varchar(80);not null;default:''"`
	TxHash            string         `gorm:"column:tx_hash;type:varchar(66)"`
	OrderStatus       int            `gorm:"column:order_status;type:int"`
	Signature         string         `gorm:"column:signature;type:varchar(150);not null;default:''"`
	SignatureHash     string         `gorm:"column:signature_hash;type:varchar(45);not null;default:''"`
	SignatureSrc      datatypes.JSON `gorm:"column:signature_source;type:json"`
	CreatedAt         time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Timestamp         int64          `gorm:"column:timestamp;type:bigint unsigned;not null"`
	Nonce             string         `gorm:"column:nonce;type:varchar(80);not null;default:''"`
	RiskStatus        int            `gorm:"column:risk_status;type:int(11);not null;default:0" `
	RiskReviewer      string         `gorm:"column:risk_reviewer;type:varchar(45)" `
	Height            uint64         `gorm:"column:height;type:bigint unsigned;not null"`
}

// TableName
func (m *TNftWithdrawRecord) TableName() string {
	return "t_nft_withdraw_record"
}

func (m *TNftWithdrawRecord) Copy() *TNftWithdrawRecord {
	return &TNftWithdrawRecord{
		ID:                m.ID,
		AppID:             m.AppID,
		UID:               m.UID,
		AppOrderID:        m.AppOrderID,
		EquipmentID:       m.EquipmentID,
		GameAssetName:     m.GameAssetName,
		Account:           m.Account,
		ContractAddress:   m.ContractAddress,
		WithdrawAddress:   m.WithdrawAddress,
		TreaseAddress:     m.TreaseAddress,
		GameMinterAddress: m.GameMinterAddress,
		TokenID:           m.TokenID,
		TxHash:            m.TxHash,
		OrderStatus:       m.OrderStatus,
		Signature:         m.Signature,
		SignatureHash:     m.SignatureHash,
		SignatureSrc:      m.SignatureSrc,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		Nonce:             m.Nonce,
		RiskStatus:        m.RiskStatus,
		RiskReviewer:      "",
		Height:            0,
	}
}
