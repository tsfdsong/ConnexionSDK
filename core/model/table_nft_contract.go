package model

import (
	"errors"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"time"
)

type TNftContract struct {
	ID              uint64    `gorm:"primaryKey;column:id;type:int;not null"`
	AppID           int       `gorm:"column:app_id;type:int"`
	ChainID         int       `gorm:"column:chain_id;type:int"`
	Treasure        string    `gorm:"column:treasure;type:varchar(45);not null;default:''"`
	ContractAddress string    `gorm:"column:contract_address;type:char(42)"`
	MinterAddress   string    `gorm:"column:minter_address;type:char(42)"`
	TokenName       string    `gorm:"column:token_name;type:varchar(45)"`
	TokenSymbol     string    `gorm:"column:token_symbol;type:varchar(45)"`
	TokenSupply     string    `gorm:"column:token_supply;type:varchar(45)"`
	GameAssetName   string    `gorm:"column:game_asset_name;type:varchar(45)"`
	DepositSwitch   int8      `gorm:"column:deposit_switch;type:tinyint;default:0"`
	WithdrawSwitch  int8      `gorm:"column:withdraw_switch;type:tinyint;default:0"`
	Decimal         int       `gorm:"column:decimal;type:int;not null;default:0"`
	BaseURL         string    `gorm:"column:base_url;type:varchar(256);not null"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	FileName        string    `gorm:"column:file_name;type:varchar(20);not null"`
	AttrUpdateTime  string    `gorm:"column:attr_update_time;type:varchar(20);not null`
}

//	TableName
func (m *TNftContract) TableName() string {
	return "t_nft_contract"
}

//	One
func (m *TNftContract) One() (one *TNftContract, err error) {
	one = &TNftContract{}
	err = queryOne(m, one)
	return
}

//	Some
func (m *TNftContract) Some() (list *[]TNftContract, err error) {
	list = &[]TNftContract{}
	err = querySome(m, list)

	return
}

//	Update
func (m *TNftContract) Update() (err error) {
	where := TNftContract{ID: m.ID}
	m.ID = 0

	return update(m, where)
}

//	Create
func (m *TNftContract) Create() (err error) {
	m.ID = 0

	return mysql.GetDB().Create(m).Error
}

//	Delete
func (m *TNftContract) Delete() (err error) {
	if m.ID == 0 {
		return errors.New("resource must not be zero value")
	}
	return delete(m)
}

func (m *TNftContract) Refresh() (err error) {

	return nil
}
func (m *TNftContract) RefreshAssets() error {
	return nil
}
