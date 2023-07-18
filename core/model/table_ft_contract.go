package model

import "time"

type TFtContract struct {
	ID              uint64    `gorm:"primaryKey;column:id;type:int;not null"`
	AppID           int       `gorm:"column:app_id;type:int"`
	ChainID         int       `gorm:"column:chain_id;type:int"`
	ContractAddress string    `gorm:"column:contract_address;type:varchar(45)"`
	TokenName       string    `gorm:"column:token_name;type:varchar(45)"`
	TokenSymbol     string    `gorm:"column:token_symbol;type:varchar(45)"`
	TokenSupply     string    `gorm:"column:token_supply;type:varchar(80)"`
	TokenDecimal    int       `gorm:"column:token_decimal;type:int"`
	GameDecimal     int       `gorm:"column:game_decimal;type:int"`
	GameCoinName    string    `gorm:"column:game_coin_name;type:varchar(45)"`
	DepositSwitch   int8      `gorm:"column:deposit_switch;type:tinyint;default:0"`
	WithdrawSwitch  int8      `gorm:"column:withdraw_switch;type:tinyint;default:0"`
	Treasure        string    `gorm:"column:treasure;type:varchar(45);not null;default:''"`
	DepositTreasure string    `gorm:"column:deposit_treasure;type:varchar(45);not null;default:''"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// TableName
func (m *TFtContract) TableName() string {
	return "t_ft_contract"
}

// Some
func (m *TFtContract) Some() (list *[]TFtContract, err error) {
	list = &[]TFtContract{}
	err = querySome(m, list)

	return
}
