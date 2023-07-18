package model

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"testing"
)

const configPath = "../"
const log = "./sdk.log"

func init_() {
	logger.Init(log)
	config.LoadConf(configPath)
	mysql.GetDB()
}

func TestTNftContract_Create(t *testing.T) {
	init_()

	m := TNftContract{
		AppID:           0,
		ChainID:         0,
		ContractAddress: "0x05A63E7155cBe76CF42E9f6CdE5DcECf6469E14d",
		TokenName:       "Bggggggggg",
		TokenSymbol:     "MMM",
		TokenSupply:     "10000000000000000000",
		GameAssetName:   "leg",
		DepositSwitch:   0,
		WithdrawSwitch:  0,
		Decimal:         0,
		BaseURL:         "https://www.55555.com/",
	}

	err := m.Create()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTNftContract_Some(t *testing.T) {
	init_()

	m := TNftContract{
		ContractAddress: "0xbCcC2073ADfC46421308f62cfD9868dF00D339a8",
	}

	list, err := m.Some()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(list)
}

func TestTNftContract_One(t *testing.T) {
	init_()

	m := TNftContract{
		ID: 1,
	}

	one, err := m.One()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(one)
}

func TestTNftContract_Update(t *testing.T) {
	init_()

	m := TNftContract{
		ID: 1,
	}

	one, err := m.One()
	if err != nil {
		t.Error(err)
		return
	}

	one.ContractAddress = "0xbCcC2073ADfC46421308f62cfD9868dF00D339a8"

	err = one.Update()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTNftContract_Delete(t *testing.T) {
	init_()

	m := TNftContract{
		ID: 1,
	}
	err := m.Delete()
	if err != nil {
		t.Error(err)
		return
	}
}
