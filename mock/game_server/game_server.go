package game_server

import (
	"encoding/json"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/core/ingame"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type MResponse struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var balanceMap map[string]string
var nftAssets map[int64]map[string]common.RGameServerERC721Asset
var nftRecoverAssets map[int64]map[string]common.RGameServerERC721Asset
var nftOrderID map[string]string

var ImageList = []string{"00007", "00012", "00024", "00042", "00158", "00179", "00237", "00291", "00334", "00365", "00399", "00427", "00484", "00527", "00656", "00682", "20018", "20186", "20279", "20372", "20528", "20655", "21157"}

func GetFTAsset(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("GetFTAsset==== %+v\n", string(readBytes))

	values := req.URL.Query()
	fmt.Printf("GetFTAsset  values==== %+v\n", values)
	uid := values.Get("uid")
	fmt.Printf("GetFTAsset  uid==== %+v\n", uid)

	balance, ok := balanceMap[uid]
	if !ok {
		balance = "0"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	erc20Ret := MResponse{
		Code:    0,
		Message: "",
		Data: []common.RGameServerERC20Asset{
			{
				AppCoinName:   "gold",
				Balance:       balance,
				FrozenBalance: "1000000.333",
			},
			// data.RGameServerERC20Asset{
			// 	AppCoinName: "test_token-2",
			// 	Name:        "GameERC20Token-2",
			// 	Symbol:      "G20-2",
			// 	Decimal:     18,
			// 	Contract:    "0x69ce0e042165c175209de9503243553a867cb65c",
			// 	//Balance:       "3000000000000000000003",
			// 	Balance:       "30000003000000000000",
			// 	FrozenBalance: "1000000000000000000000",
			// },
		},
	}

	resp, _ := json.Marshal(erc20Ret)
	w.Write(resp)
}

func GetNFTAsset(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()

	uidString := values.Get("uid")
	pageString := values.Get("page")
	pageSizeString := values.Get("pageSize")

	fmt.Printf("input: %s %s %s", uidString, pageString, pageSizeString)

	// uid, err := strconv.ParseInt(uidString, 0, 64)
	// if err != nil {
	// 	fmt.Printf("GetNFTAsset Decode uid error: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	imageList := []string{"1725", "1721", "1302", "1216", "1175", "1606", "2122", "2127", "4148", "1511", "1201", "1058", "1286", "1232", "1061", "2118", "1110", "1250", "2136", "1309"}

	equipIDList := []string{"558225867460837377", "559249199253684225", "559053546405756928", "567108400994844674", "567749021229121536", "568384521770106880", "568047538636062720", "568397282117943300", "568594897858199552", "568182293234974720", "568402758201245696", "567602889261842432", "568385921929445376", "567788075366744064", "567788058186874880", "568772795403599872", "568047585880702976", "568777129025601536", "569141277827792896", "569141342252302336"}

	erc721Ret := common.RNFTAssetResponse{
		Code:    0,
		Message: "",
		Total:   100,
		Data:    []common.RGameServerERC721Asset{},
	}

	dataRes := make([]common.RGameServerERC721Asset, 0)
	for i := 0; i < 20; i++ {

		item := common.RGameServerERC721Asset{
			GameAssetName: "ArchLootPart",
			TokenID:       fmt.Sprintf("%d", i),
			EquipmentID:   equipIDList[i],
			Frozen:        false,
			Image:         imageList[i],
			Attrs: []commdata.EquipmentAttr{
				{
					AttributeID:    0,
					AttributeName:  "",
					AttributeValue: "1",
				},
				{
					AttributeID:    1,
					AttributeName:  fmt.Sprintf("Attack +%d", i+2),
					AttributeValue: fmt.Sprintf("%d", i+2),
				},
				{
					AttributeID:    2,
					AttributeName:  fmt.Sprintf("Defense +%d", i+3),
					AttributeValue: fmt.Sprintf("%d", i+3),
				},
			},
		}

		dataRes = append(dataRes, item)
	}

	erc721Ret.Data = dataRes
	erc721Ret.Total = int64(len(erc721Ret.Data))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	resp, _ := json.Marshal(erc721Ret)
	w.Write(resp)
}

// func GetNFTAsset(w http.ResponseWriter, req *http.Request) {
// 	values := req.URL.Query()

// 	uidString := values.Get("uid")
// 	pageString := values.Get("page")
// 	pageSizeString := values.Get("pageSize")

// 	uid, err := strconv.ParseInt(uidString, 0, 64)
// 	if err != nil {
// 		fmt.Printf("GetNFTAsset Decode uid error: %v\n", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	erc721Ret := common.RNFTAssetResponse{
// 		Code:    0,
// 		Message: "",
// 		Total:   100,
// 		Data:    []common.RGameServerERC721Asset{},
// 	}

// 	list, ok := nftAssets[uid]
// 	if !ok {
// 		list = make(map[string]common.RGameServerERC721Asset, 20)
// 		for i := 1; i <= 20; i++ {
// 			image := ImageList[i-1]

// 			eid := fmt.Sprintf("%d", i)
// 			list[eid] = common.RGameServerERC721Asset{
// 				GameAssetName: "ArchLoot",
// 				TokenID:       "",
// 				EquipmentID:   fmt.Sprintf("%d", i),
// 				Frozen:        false,
// 				Image:         image,
// 				Attrs: []commdata.EquipmentAttr{
// 					{
// 						AttributeID:    1,
// 						AttributeName:  fmt.Sprintf("name_%d", 1000*i+1),
// 						AttributeValue: fmt.Sprintf("%d", 1000*i+1),
// 					},
// 					{
// 						AttributeID:    2,
// 						AttributeName:  fmt.Sprintf("name_%d", 1000*i+2),
// 						AttributeValue: fmt.Sprintf("%d", 1000*i+2),
// 					},
// 					{
// 						AttributeID:    3,
// 						AttributeName:  fmt.Sprintf("name_%d", 1000*i+3),
// 						AttributeValue: fmt.Sprintf("%d", 1000*i+3),
// 					},
// 				},
// 			}
// 		}

// 		nftAssets[uid] = list
// 	}

// 	fmt.Printf("\n GetNFTAsset==== %v\n", uid)
// 	DisplayNFTAssets(uid)

// 	datas := make([]common.RGameServerERC721Asset, 0)
// 	for _, v := range list {
// 		datas = append(datas, v)
// 	}

// 	page, err := strconv.ParseInt(pageString, 0, 64)
// 	if err != nil {
// 		fmt.Printf("GetNFTAsset Decode page error: %v %v\n", pageString, err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	pageSize, err := strconv.ParseInt(pageSizeString, 0, 64)
// 	if err != nil {
// 		fmt.Printf("GetNFTAsset Decode pageSize error: %v %v\n", pageSizeString, err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	start := int((page - 1) * pageSize)
// 	end := int(page * pageSize)
// 	length := len(datas)

// 	if end > length {
// 		end = length
// 	}

// 	erc721Ret.Data = datas[start:end]
// 	erc721Ret.Total = int64(len(erc721Ret.Data))
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	resp, _ := json.Marshal(erc721Ret)
// 	w.Write(resp)
// }

func GetNFTAssetDetail(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	eqIDString := values.Get("equipment_id")

	i, err := strconv.ParseInt(eqIDString, 0, 64)
	if err != nil {
		fmt.Printf("GetNFTAssetDetails Decode equipmentid error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	image := ImageList[i-1]

	erc721Ret := MResponse{
		Code:    0,
		Message: "",
		Data: common.RGameServerERC721Asset{
			GameAssetName: "ArchLoot",
			TokenID:       "",
			EquipmentID:   fmt.Sprintf("%d", i),
			Frozen:        false,
			Image:         image,
			Attrs: []commdata.EquipmentAttr{
				{
					AttributeID:    1,
					AttributeName:  fmt.Sprintf("name_%d", 1000*i+1),
					AttributeValue: fmt.Sprintf("%d", 1000*i+1),
				},
				{
					AttributeID:    2,
					AttributeName:  fmt.Sprintf("name_%d", 1000*i+2),
					AttributeValue: fmt.Sprintf("%d", 1000*i+2),
				},
				{
					AttributeID:    3,
					AttributeName:  fmt.Sprintf("name_%d", 1000*i+3),
					AttributeValue: fmt.Sprintf("%d", 1000*i+3),
				},
			},
		},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp, _ := json.Marshal(erc721Ret)
	w.Write(resp)
}

func FTDeposit(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("FTDeposit==== %+v\n", string(readBytes))
	s := new(ingame.NotifyFTDeposits)
	json.Unmarshal(readBytes, s)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type Resp struct {
		Code    int64                     `json:"code"`
		Message string                    `json:"message"`
		Data    []ingame.GameFTDepositRes `json:"data"`
	}
	data := []ingame.GameFTDepositRes{}
	for _, e := range s.Params {
		item := ingame.GameFTDepositRes{
			GameCoinName: s.Params[0].GameCoinName,
			AppOrderID:   tools.GenCode(10),
			TxHash:       s.Params[0].TxHash,
			Status:       1,
		}
		data = append(data, item)
		key := fmt.Sprintf("%d", e.Uid)
		before, ok := balanceMap[key]
		if !ok {
			balanceMap[key] = e.Amount
		} else {
			amount, _ := big.NewFloat(0).SetString(e.Amount)
			beforeAmount, _ := big.NewFloat(0).SetString(before)
			after := big.NewFloat(0).Add(amount, beforeAmount)

			balanceMap[key] = after.String()
			fmt.Printf("amount:%+v,beforeAmount:%+v,after:%+v,value:%+v\n", amount, beforeAmount, after, balanceMap[key])
		}
	}

	erc20Ret := Resp{
		Code:    0,
		Message: "success",
		Data:    data,
	}

	resp, _ := json.Marshal(erc20Ret)
	w.Write(resp)
}

func NFTDeposit(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)

	type PNotifyNFTDeposit struct {
		GameAssetName string                   `json:"game_asset_name"`
		TokenID       string                   `json:"token_id"`
		EquipmentID   string                   `json:"equipment_id"`
		TxHash        string                   `json:"tx_hash"`
		Uid           int64                    `json:"uid"`
		Attrs         []commdata.EquipmentAttr `json:"attrs"`
	}

	type PNotifyNFTDeposits struct {
		Params   []PNotifyNFTDeposit `json:"params"`
		AppID    int                 `json:"appId"`
		SignHash string              `json:"sign"`
	}

	s := new(PNotifyNFTDeposits)
	err := json.Unmarshal(readBytes, s)
	if err != nil {
		fmt.Printf("NFTDeposit Decode input error: %v\n", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	type RNFTDepositRes struct {
		GameAssetName string `json:"game_asset_name"`
		AppOrderID    string `json:"app_order_id"`
		TokenId       string `json:"token_id"`
		EquipmentId   string `json:"equipment_id"`
		TxHash        string `json:"tx_hash"`
		Status        int    `json:"status"`
	}

	datas := []RNFTDepositRes{}
	for _, e := range s.Params {
		eid := e.EquipmentID
		if eid == "" {
			//generate equipment id
			eid = fmt.Sprintf("%d", time.Now().Unix()%100+20)
		}

		item := RNFTDepositRes{
			GameAssetName: e.GameAssetName,
			AppOrderID:    tools.GenCode(10),
			TokenId:       e.TokenID,
			EquipmentId:   eid,
			TxHash:        e.TxHash,
			Status:        const_def.GAMESERVER_PASS,
		}
		datas = append(datas, item)

		uid := e.Uid

		fmt.Printf("deposit add equipment id: %v %v \n", uid, eid)

		image := ""
		enum, _ := strconv.ParseInt(eid, 10, 64)
		if enum > 20 {
			index := enum%3 + 20
			image = ImageList[index]
		} else if 0 < enum && enum <= 20 {
			image = ImageList[enum-1]
		}

		newitem := common.RGameServerERC721Asset{
			GameAssetName: e.GameAssetName,
			TokenID:       e.TokenID,
			EquipmentID:   eid,
			Frozen:        false,
			Image:         image,
			Attrs: []commdata.EquipmentAttr{
				{
					AttributeID:    1,
					AttributeName:  fmt.Sprintf("name_%d", 1000*enum+1),
					AttributeValue: fmt.Sprintf("%d", 1000*enum+1),
				},
				{
					AttributeID:    2,
					AttributeName:  fmt.Sprintf("name_%d", 1000*enum+2),
					AttributeValue: fmt.Sprintf("%d", 1000*enum+2),
				},
				{
					AttributeID:    3,
					AttributeName:  fmt.Sprintf("name_%d", 1000*enum+3),
					AttributeValue: fmt.Sprintf("%d", 1000*enum+3),
				},
			},
		}

		newvalue, ok := nftAssets[uid]
		if !ok {
			neweid := make(map[string]common.RGameServerERC721Asset, 1)
			neweid[eid] = newitem
			nftAssets[uid] = neweid
		} else {
			newvalue[eid] = newitem
		}

		fmt.Printf("deposit add item: %v \n", newitem)

		DisplayNFTAssets(uid)
	}

	type RNotifyNFTDeposit struct {
		Code    int64            `json:"code"`
		Message string           `json:"message"`
		Data    []RNFTDepositRes `json:"data"`
	}

	erc721Ret := RNotifyNFTDeposit{
		Code:    0,
		Message: "success",
		Data:    datas,
	}

	fmt.Printf("NFT deposit return: %v\n", erc721Ret)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp, _ := json.Marshal(erc721Ret)
	w.Write(resp)
}
func ERC20PreWithdraw(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("ERC20PreWithdraw==== %+v\n", string(readBytes))

	s := new(common.PGameERC20PreWithdraw)
	json.Unmarshal(readBytes, s)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	key := fmt.Sprintf("%d", s.Params.Uid)
	before, ok := balanceMap[key]
	if ok {
		amount, _ := big.NewFloat(0).SetString(s.Params.Amount)
		beforeAmount, _ := big.NewFloat(0).SetString(before)
		after := big.NewFloat(0).Sub(beforeAmount, amount)
		balanceMap[key] = after.String()
	}

	orderid := fmt.Sprintf("%+v", time.Now().Unix())

	erc20PreWithdraw := common.RERC20PreWithdraw{
		Code:    0,
		Message: "success",
		Data: common.RERC20PreWithdrawData{
			AppOrderID: orderid,
			Status:     1,
		},
	}

	resp, _ := json.Marshal(erc20PreWithdraw)
	w.Write(resp)
}
func ERC20Withdraw(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("ERC20Withdraw==== %+v\n", string(readBytes))

	s := new(common.PGameERC20WithdrawComfirm)
	json.Unmarshal(readBytes, s)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	d := []common.RERC20WithdrawComfirmStatus{}
	for _, e := range s.Params {
		item := common.RERC20WithdrawComfirmStatus{
			AppOrderID:   tools.GenCode(10),
			Nonce:        e.Nonce,
			GameCoinName: e.GameCoinName,
			Status:       1,
		}
		d = append(d, item)

		key := fmt.Sprintf("%d", e.Uid)
		before, ok := balanceMap[key]
		if ok && e.Operation == const_def.NOTI_GAMESERVER_RECOVER {
			amount, _ := big.NewFloat(0).SetString("1000")
			beforeAmount, _ := big.NewFloat(0).SetString(before)
			after := big.NewFloat(0).Add(beforeAmount, amount)
			balanceMap[key] = after.String()
		}

	}

	erc20Withdraw := common.RERC20WithdrawComfirm{
		Code:    0,
		Message: "Success",
		Data:    d,
	}

	resp, _ := json.Marshal(erc20Withdraw)
	w.Write(resp)
}

func DisplayNFTAssets(uid int64) {
	item, ok := nftAssets[uid]
	if !ok {
		fmt.Printf("%v not found\n", uid)
	}
	fmt.Printf("\n *********** uid=%v ****\n", uid)
	for k, v := range item {
		fmt.Printf("EquipmenID => TokenID, %v  %v\n", k, v.TokenID)
	}
	fmt.Print("\n *********** end ****\n")
}

func ERC721PreWithdraw(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)

	type NftPreWithdrawData struct {
		AppID    int                         `json:"appId"`
		Params   []ingame.PrewithdrawDataRes `json:"params"`
		SignHash string                      `json:"sign"`
	}

	var p NftPreWithdrawData
	err := json.Unmarshal(readBytes, &p)
	if err != nil {
		fmt.Printf("ERC721PreWithdraw Decode input error: %v \n %v\n", err, string(readBytes))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	type Resp struct {
		Code    int64                       `json:"code"`
		Message string                      `json:"message"`
		Data    []ingame.PrewithdrawDataRes `json:"data"`
	}

	fmt.Printf("nftAssets before prewithdraw : %v \n", nftAssets)

	fmt.Printf("\nprewithdraw input info: %v \n", p)

	datas := make([]ingame.PrewithdrawDataRes, 0)

	for _, v := range p.Params {
		uid := v.UID
		eID := v.EquipmentID

		orderid := fmt.Sprintf("%s_%d", eID, time.Now().UnixNano())

		fmt.Printf("prewithdraw: %v %v\n", uid, eID)
		DisplayNFTAssets(int64(uid))

		itemValue := ingame.PrewithdrawDataRes{
			GameAssetName: v.GameAssetName,
			UID:           uid,
			AppOrderID:    orderid,
			EquipmentID:   eID,
			Status:        const_def.GAMESERVER_PASS,
		}

		item, ok := nftAssets[int64(uid)][eID]
		if !ok {
			itemValue.Status = const_def.GAMESERVER_REJECT
			fmt.Printf("ERC721PreWithdraw prewithdraw not found: %v %v\n", uid, eID)
		} else {
			itemValue.Status = const_def.GAMESERVER_PASS

			nftOrderID[orderid] = item.EquipmentID
			fmt.Printf("ERC721PreWithdraw add app order id: %v %v\n", orderid, item.EquipmentID)

			newValue := make(map[string]common.RGameServerERC721Asset, 1)
			newValue[item.EquipmentID] = item
			nftRecoverAssets[int64(uid)] = newValue
			delete(nftAssets[int64(uid)], item.EquipmentID)
		}
		datas = append(datas, itemValue)
	}

	nftPreWithdraw := Resp{
		Code:    0,
		Message: "success",
		Data:    datas,
	}
	fmt.Printf("NFT prewithdraw return: %v\n", nftPreWithdraw)

	resp, err := json.Marshal(nftPreWithdraw)
	if err != nil {
		fmt.Printf("ERC721PreWithdraw Decode prewithdraw error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(resp)
}

func ERC721Withdraw(w http.ResponseWriter, req *http.Request) {
	readBytes, _ := ioutil.ReadAll(req.Body)

	type NftWithdrawData struct {
		AppID    int                         `json:"appId"`
		Params   []ingame.CommitWithdrawData `json:"params"`
		SignHash string                      `json:"sign"`
	}

	var p NftWithdrawData
	err := json.Unmarshal(readBytes, &p)
	if err != nil {
		fmt.Printf("Decode input error: %v\n %v\n", err, string(readBytes))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	data := make([]ingame.CommitwithdrawDataRes, 0)
	for _, v := range p.Params {
		item := ingame.CommitwithdrawDataRes{
			GameAssetName: v.GameAssetName,
			Nonce:         v.Nonce,
			UID:           v.UID,
			AppOrderID:    v.AppOrderID,
			Status:        const_def.GAMESERVER_PASS,
		}

		if v.Operate == const_def.NOTI_GAMESERVER_RECOVER {
			uid := int64(v.UID)
			eid, ok := nftOrderID[v.AppOrderID]
			if !ok {
				err := fmt.Errorf("orderID record not found: %v %v", uid, v.AppOrderID)
				fmt.Printf("withdraw recover game assets: %v\n", err)

				item.Status = const_def.GAMESERVER_REJECT
				continue
			}

			value, ok := nftRecoverAssets[uid][eid]
			if !ok {
				err := fmt.Errorf("nftRecoverAssets record not found: %v %v", uid, eid)
				fmt.Printf("withdraw recover game assets: %v\n", err)

				item.Status = const_def.GAMESERVER_REJECT
				continue
			}

			fmt.Printf("withdraw recover game assets, %v %v %v\n", uid, eid, v.AppOrderID)
			nftAssets[uid][eid] = value
			DisplayNFTAssets(uid)
		}

		if v.Operate == const_def.NOTI_GAMESERVER_DELETE {
			uid := int64(v.UID)
			eid, ok := nftOrderID[v.AppOrderID]
			if !ok {
				err := fmt.Errorf(" record not found: %v %v", uid, v.AppOrderID)
				fmt.Printf("withdraw delete game assets: %v\n", err)

				item.Status = const_def.GAMESERVER_REJECT
				continue
			}

			fmt.Printf("withdraw delete game assets, %v %v %v\n", v.UID, eid, v.AppOrderID)
			delete(nftOrderID, v.AppOrderID)
			DisplayNFTAssets(uid)
		}

		data = append(data, item)
	}

	type Resp struct {
		Code    int64                          `json:"code"`
		Message string                         `json:"message"`
		Data    []ingame.CommitwithdrawDataRes `json:"data"`
	}

	erc721Withdraw := Resp{
		Code:    0,
		Message: "Success",
		Data:    data,
	}

	fmt.Printf("NFT claim return: %v\n", erc721Withdraw)

	resp, err := json.Marshal(erc721Withdraw)
	if err != nil {
		fmt.Printf("Decode marshal claim: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(resp)
}

func GameServerMock() {
	balanceMap = map[string]string{}
	nftAssets = make(map[int64]map[string]common.RGameServerERC721Asset, 0)
	nftRecoverAssets = make(map[int64]map[string]common.RGameServerERC721Asset, 0)
	nftOrderID = make(map[string]string, 0)

	//FT
	http.HandleFunc("/ftToken/getAssets", GetFTAsset)

	http.HandleFunc("/ftToken/deposit", FTDeposit)
	http.HandleFunc("/ftToken/preWithdraw", ERC20PreWithdraw)
	http.HandleFunc("/ftToken/withdraw", ERC20Withdraw)

	//NFT
	http.HandleFunc("/nftToken/getAssets", GetNFTAsset)
	http.HandleFunc("/nftToken/getAssetDetail", GetNFTAssetDetail)

	http.HandleFunc("/nftToken/deposit", NFTDeposit)
	http.HandleFunc("/nftToken/preWithdraw", ERC721PreWithdraw)
	http.HandleFunc("/nftToken/withdraw", ERC721Withdraw)
	log.Fatal(http.ListenAndServe("0.0.0.0:8889", nil))
}
