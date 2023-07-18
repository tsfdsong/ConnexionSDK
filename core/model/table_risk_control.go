package model

import "time"

type TRiskControl struct {
	ID              uint64    `gorm:"primaryKey;<-:create;column:id;type:int;not null"`
	AppID           int       `gorm:"column:app_id;type:int"`
	TokenType       string    `gorm:"column:token_type;type:varchar(16)"`
	ContractAddress string    `gorm:"column:contract_address;type:varchar(45)"`
	TokenName       string    `gorm:"column:token_name;type:varchar(45)"`
	TokenSymbol     string    `gorm:"column:token_symbol;type:varchar(45)"`
	GameCoinName    string    `gorm:"column:game_coin_name;type:varchar(45)"`
	AmountLimit     string    `gorm:"column:amount_limit;type:varchar(45);default:0"`
	CountLimit      int       `gorm:"column:count_limit;type:int;default:0"`
	CountTime       int       `gorm:"column:count_time;type:int;default:0"`
	TotalLimit      string    `gorm:"column:total_limit;type:varchar(45);default:0"`
	TotalTime       int       `gorm:"column:total_time;type:int;default:0"`
	Status          int8      `gorm:"column:status;type:tinyint;default:0"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// TableName
func (m *TRiskControl) TableName() string {
	return "t_risk_control"
}
