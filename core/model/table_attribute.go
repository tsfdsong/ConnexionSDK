package model

import "time"

type TAttribute struct {
	ID              uint64    `gorm:"primaryKey;column:id;type:bigint unsigned;not null"`
	AppID           uint8     `gorm:"index:idx_t_attributes_game_id;column:app_id;type:tinyint unsigned;not null"`
	AttrID          uint64    `gorm:"index:idx_t_attributes_attr_id;column:attr_id;type:bigint unsigned;not null"`
	AttrDecimal     uint64    `gorm:"column:attr_decimal;type:bigint;unsigned;not null"`
	AttrName        string    `gorm:"index:idx_t_attributes_attr_name;column:attr_name;type:varchar(256);not null"`
	AttrDescription string    `gorm:"column:attr_description;type:varchar(256);not null"`
	ContractAddress string    `gorm:"column:contract_address;type:varchar(45);not null"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

//TableName
func (m *TAttribute) TableName() string {
	return "t_attribute"
}
