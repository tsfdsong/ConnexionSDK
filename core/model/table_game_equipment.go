package model

import (
	"time"

	"gorm.io/datatypes"
)

type TGameEquipment struct {
	ID              uint64         `gorm:"primaryKey;column:id;type:int;not null"`
	AppID           int            `gorm:"uniqueIndex:aap_eqe_id;column:app_id;type:int;not null;default:0" json:"app_id"`
	GameAssetName   string         `gorm:"column:game_asset_name;type:varchar(45)" json:"game_asset_name"`
	Account         string         `gorm:"column:account;type:varchar(100)"`
	EquipmentID     string         `gorm:"uniqueIndex:aap_eqe_id;column:equipment_id;type:varchar(45)" json:"equipment_id"`
	ContractAddress string         `gorm:"column:contract_address;type:varchar(45)"`
	TokenID         string         `gorm:"column:token_id;type:varchar(80)" json:"token_id"`
	ChainID         int            `gorm:"column:chain_id;type:int"`
	Status          uint32         `gorm:"column:status;type:int(4) unsigned;not null;default:0"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	WithdrawSwitch  int8           `gorm:"column:withdraw_switch;type:tinyint"`
	EquipmentAttr   datatypes.JSON `gorm:"column:equipment_attr;type:json"`
	ImageURI        string         `gorm:"column:image_uri;varchar(256)"`
}

// TableName
func (m *TGameEquipment) TableName() string {
	return "t_game_equipment"
}
