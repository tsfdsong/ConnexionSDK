package db_operation

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/model"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"strings"
)

func Insert() {
	addressList := []string{
		"0x7720cA66122701469FC0c99C99841cB35Caef3B3",
		"0x545B534Fe1FB4E3BA601CCe7cc38891C09a0881f",
		"0xA728CB0B29b624ecd3C11FDc7Eff0e232A9b2a4F",
		"0x252F63491a1AFf9f7d4725A54101A3D85Ec1971E",
		"0xec7F342d42CF0C6B1Ba1bD956C44568c5b95C4f7",
		"0xC1c6cfBbDA2F3ED56927608feAA41a4CC37b0ED8",
		"0x6A36779D33335448Bd2093bfa992EB7Daf97b8d1",
		"0xd1A2A9DaB2E9f3cE32bfF6D9cfa6B13fdAa8fc89",
		"0xcFF830Aa513247883bB4c17247bd9111E863CCc5",
		"0xAbB4A6F586914A2BD34f37902F95735fF355e208",
		"0x6354068D1A1238e7f64c2c5C08215C1A231f94D1",
		"0x53F61e6660698f99fFF02B14f4347ce24E4B0Eb6",
		"0xd02F274d9bd49179d8579A1A4D1ac0671aE83518",
		"0xA55a8CaF953fC04a061f142CeB97001df9d1DE62",
		"0x6270C3cffb9CdF0Cbb6736b4d732506ee081F176",
		"0x4AD7966ac5798a1B2B7EE4Bd0d7bDa0850D74613",
		"0x955c60b23830b8dA9f4585aA369ab90cb3e8F8eD",
		"0x1a3efd8Ddb7BA5a5293dfe3Ca022975e1c4de8cC",
		"0xbDd66DB90f4110bbCF4106f9f24573f2b0208Bb3",
		"0x54Bb2232F0E35067A19217cD97Bcf3E02bD3f4eE",
	}
	uidList := []uint64{
		6488706401748828160,
		6488706401748828163,
		6488706401748828166,
		6488706401748828169,
		6488706401748828172,
		6488706401748828175,
		6488706401748828178,
		6488706401748828181,
		6488706401748828184,
		6488706401748828187,
		6488706401748828190,
		6488706401748828193,
		6488706401748828196,
		6488706401748828199,
		6488706401748828202,
		6488706401748828205,
		6488706401748828208,
		6488706401748828211,
		6488706401748828214,
		6488706401748828217,
	}
	accountList := []string{
		"test1@archloot.com",
		"test2@archloot.com",
		"test3@archloot.com",
		"test4@archloot.com",
		"test5@archloot.com",
		"test6@archloot.com",
		"test7@archloot.com",
		"test8@archloot.com",
		"test9@archloot.com",
		"test10@archloot.com",
		"test11@archloot.com",
		"test12@archloot.com",
		"test13@archloot.com",
		"test14@archloot.com",
		"test15@archloot.com",
		"test16@archloot.com",
		"test17@archloot.com",
		"test18@archloot.com",
		"test19@archloot.com",
		"test20@archloot.com",
	}

	nameList := []string{
		"test1@archloot.com_2",
		"test2@archloot.com_2",
		"test3@archloot.com_2",
		"test4@archloot.com_2",
		"test5@archloot.com_2",
		"test6@archloot.com_2",
		"test7@archloot.com_2",
		"test8@archloot.com_2",
		"test9@archloot.com_2",
		"test10@archloot.com_2",
		"test11@archloot.com_2",
		"test12@archloot.com_2",
		"test13@archloot.com_2",
		"test14@archloot.com_2",
		"test15@archloot.com_2",
		"test16@archloot.com_2",
		"test17@archloot.com_2",
		"test18@archloot.com_2",
		"test19@archloot.com_2",
		"test20@archloot.com_2",
	}
	bindList := []model.TEmailBind{}
	userList := []model.TUser{}
	for k, e := range accountList {
		item := model.TEmailBind{
			AppID:   2,
			UID:     uidList[k],
			Address: strings.ToLower(addressList[k]),
			Account: e,
		}
		bindList = append(bindList, item)

		userItem := model.TUser{
			AppID:   2,
			Account: e,
			Name:    nameList[k],
		}
		userList = append(userList, userItem)
	}

	err := mysql.WrapInsertBatch(model.TableEmailBind, bindList, len(bindList))
	if err != nil {
		fmt.Printf("insert batch bind email failed:%+v\n", err)
	}

	err = mysql.WrapInsertBatch(model.TableUser, userList, len(userList))
	if err != nil {
		fmt.Printf("insert batch user failed:%+v\n", err)
	}

}
