---
title: ConnexionSDK API Reference

# language_tabs: # must be one of https://git.io/vQNgJ
#   - shell
#   - ruby
#   - python
#   - javascript

toc_footers:
  - <a href='https://connexion.games/'>To the Connexion website</a>
  - <a href='https://github.com/slatedocs/slate'>Documentation Powered by Slate</a>

# includes:
#   - errors

search: true

code_clipboard: true

meta:
  - name: description
    content: Documentation for the ConnexionSDK API
---

# Introduction

***ConnexionSDK*** ***API v1.0.0***


# JSON with HTTP

* To deploy ConnexionSDK, developer use JSON with HTTP to realize documentation files API as below


* <font color=#FF0000 >Game developers should restrict IP address access (only allowing developer consented SDK server ServerIP to access）</font>


## Request messages

* HTTP head Content-Type format is application/json
* API request parameters should include phrases such as sign,params,appId

| **Parameter** | **Required** | **Type** | **Enum** | **Description**   |
| ------------------- | ------------------ | -------------- | -------------- | ----------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK |
| sign                | true               | string         |                | Signature Info          |
| params              | true               | array,object   |                | Request Parameter       |

* Signature is used to establish secured communication with ConnexionSDK, developers need to check to make sure that the parameters agree with the signed data. Signature rules are as such:

> Golang sign example

```golang
func GameSign(appSecret string, appId int, params interface{}) (string, error) {
	bytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	rawStr := fmt.Sprintf("appId=%d&params=%s&secretKey=%s", appId, string(bytes), appSecret)
	res := md5.Sum([]byte(rawStr))
	mdHash := fmt.Sprintf("%x", res)
	return strings.ToUpper(mdHash), nil
}

```

1. Convert request parameters 'params' to json string
2. Create a signature sting:appId={appId(from request parameters)}&params={json siganature string from step 1}&secretKey={AppSecret set up in backend dashboard}
3. Calculate MD5 checksum with sigature string from step2 (All Capitals)

## Response messages

Request corresponding return parameters as below：

| **Parameter** | **Required** | **Type** | **Enum** | **Description**    |
| ------------------- | ------------------ | -------------- | -------------- | ------------------------ |
| code                | true               | number         |                | 0(success),others(error) |
| message             | true               | string         |                | Message Error            |
| data                | true               | array,object   |                | Return Data              |

# Assets

## Get FT Assets in Game

### **Description**

* Search user's in game FT assets

### Path

* /ftToken/getAssets?appId={%d}&uid={%d}

### Method

* GET

> Request example

```json
appId=1&uid=1001
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**   |
| ------------------- | ------------------ | -------------- | -------------- | ----------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK |
| uid                 | true               | number         |                | Game User ID            |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": [{
        "game_coin_name": "gold",
        "coin_balance": "10056.3355",
        "coin_frozen_balance": "5033.3735",
    }]
}
```

### **Response**

| Name                     | Type   | Enum | Description                |
| ------------------------ | ------ | ---- | -------------------------- |
| data                     | array  |      | Balance of User's FT asset |
| &gt; game_coin_name      | string |      | Game FT asset name         |
| &gt; coin_balance        | string |      | In game FT balance         |
| &gt; coin_frozen_balance | string |      | In game frozen FT balance  |

## FT Deposit

### NOTICE

<font color=#FF0000 >Devlopers to record deposit orders info</font>

<font color=#FF0000 >Recommedation - when processing deposit orders, force users offline</font>

### **Description**

* Inform Game FT deposit

### Path

* /ftToken/deposit

### Method

* POST

> Request example

```json
{
    "appId": 1,
    "sign":"617070...ECF8427E",
    "params": [{
        "game_coin_name": "gold",
        "amount": "7746.1847",
        "tx_hash": "0x1e1f51d74d03f7e9f9419f7d29325350730d9581f0714ae516c149227d6829ab",
        "uid": 1001
    }]
}
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**       |
| ------------------- | ------------------ | -------------- | -------------- | --------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK     |
| sign                | true               | string         |                | Signature string of request |
| params              | true               | arrry          |                | FT deposit info             |

params fields

| **Parameter** | **Required** | **Type** | **Enum** | **Description**     |
| ------------------- | ------------------ | -------------- | -------------- | ------------------------- |
| game_coin_name      | true               | string         |                | In game asset name        |
| amount              | true               | string         |                | FT deposit amount         |
| tx_hash             | true               | string         |                | On chain transaction hash |
| uid                 | true               | number         |                | Game user ID              |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": [{
        "game_coin_name": "gold",
        "app_order_id": "3855182303",
        "tx_hash": "0x1e1f51d74d03f7e9f9419f7d29325350730d9581f0714ae516c149227d6829ab",
        "status": 1
    }]
}
```

### **Response**

| Name                | Type   | Enum                  | Description                        |
| ------------------- | ------ | --------------------- | ---------------------------------- |
| data                | array  |                       | In game FT deposit processing list |
| &gt; game_coin_name | string |                       | In game asset name                 |
| &gt; app_order_id   | string |                       | Deposit order number               |
| &gt; tx_hash        | string |                       | On chain transaction hash          |
| &gt; status         | number | 1(success),2(failure) | Status code                        |

## FT PreWithdraw

### NOTICE

<font color=#FF0000 >If successfuly, freeze user's corresponding assets</font>

<font color=#FF0000 >Developers to record withraw orders info</font>

<font color=#FF0000 >Recommendation - when processing withdraw orders, force user offline</font>

### **Description**

* FT asset pre-withdraw

### Path

* /ftToken/preWithdraw

### Method

* POST

> Request example

```json
{
    "appId": 1,
    "sign": "324460...BC78249AE",
    "params": {
        "game_coin_name": "gold",
        "amount": "10000.2374",
        "uid": 1001,
        "nonce":"13505205326163198829256798168532583699875537401379248958576017862581953029599"
    }
}
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**       |
| ------------------- | ------------------ | -------------- | -------------- | --------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK     |
| sign                | true               | string         |                | Signature string of request |
| params              | true               | object         |                | FT pre-withdraw info        |

params fields

| **Parameter** | **Required** | **Type** | **Enum** | **Description**  |
| ------------------- | ------------------ | -------------- | -------------- | ---------------------- |
| game_coin_name      | true               | string         |                | In game asset name     |
| amount              | true               | string         |                | FT pre-withdraw amount |
| uid                 | true               | number         |                | Game user ID           |
| nonce               | true               | string         |                | ConnexionSDK FT withdraw order number |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "app_order_id": "2865971757",
        "status": 1
    }
}
```

### **Response**

| Name              | Type   | Enum                  | Description                                |
| ----------------- | ------ | --------------------- | ------------------------------------------ |
| data              | object |                       | In game FT pre-withdraw processing results |
| &gt; app_order_id | string |                       | In game FT pre-withdraw order number       |
| &gt; status       | number | 1(success),2(failure) | Status code                                |

## FT Withdraw

<font color=#FF0000 >If successful, permanently delete or storage coressponding assets</font>

<font color=#FF0000 >Developers to update withdraw orders info</font>

<font color=#FF0000 >Recommendation - when processing withdraw orders, force users offline</font>

### **Description**

* FT asset official withdraw

### Path

* /ftToken/withdraw

### Method

* POST

> Request example

```json
{
    "appId": 1,
    "sign": "2YA43T0...55YGF5AH",
    "params": [{
        "game_coin_name": "gold",
        "app_order_id": "0455513284",
        "nonce": "13505205326163198829256798168532583699875537401379248958576017862581953029599",
        "operation": 1,
        "uid": 1001
    }]
}

```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**      |
| ------------------- | ------------------ | -------------- | -------------- | -------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK    |
| sign                | true               | string         |                | Signature sting of request |
| params              | true               | array          |                | FT official withdraw list  |

params fields

| **Parameter** | **Required** | **Type** | **Enum**       | **Description**                 |
| ------------------- | ------------------ | -------------- | -------------------- | ------------------------------------- |
| game_coin_name      | true               | string         |                      | In game asset name                    |
| app_order_id        | true               | string         |                      | In game FT withdraw order number      |
| nonce               | true               | string         |                      | ConnexionSDK FT withdraw order number |
| operation           | true               | number         | 1(delete),2(recover) | action                                |
| uid                 | true               | number         |                      | Game user ID                          |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": [{
        "game_coin_name": "gold",
        "app_order_id": "0455513284",
        "nonce": "13505205326163198829256798168532583699875537401379248958576017862581953029599",
        "status": 1
    }]
}
```

### **Response**

| Name                | Type   | Enum                  | Description                                     |
| ------------------- | ------ | --------------------- | ----------------------------------------------- |
| data                | array  |                       | In game FT official withdraw processing results |
| &gt; game_coin_name | string |                       | In game asset name                              |
| &gt; app_order_id   | string |                       | In game FT withdraw order number                |
| &gt; nonce          | string |                       | SDK FT withdraw order number                    |
| &gt; status         | number | 1(success),2(failure) | Status code                                     |

## Get NFT Assets in Game

### **Description**

* Search user's in game NFT asset

### Path

* /nftToken/getAssets?appId={%d}&uid={%d}&page={%d}&pageSize={%d}

### Method

* GET

> Request example

```json
appId=1&uid=1001&page=1&pageSize=30
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**   |
| ------------------- | ------------------ | -------------- | -------------- | ----------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK |
| uid                 | true               | number         |                | Game user ID            |
| page                | true               | number         |                | asset list page number  |
| pageSize            | true               | number         |                | number of NFTs listed on a single page|

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "total":100,
    "data": [{
        "game_asset_name": "MonsterWeapon",
        "token_id": "333",
        "equipment_id": "24",
        "frozen": false,
        "image": "https://hadesathena.github.io/images/333.png",
        "attrs": [{
            "attribute_id": 1,
            "attribute_value": "100000000000000"
        }]
    }]
}
```

### **Response**

| Name                 | Type   | Enum | Description                                     |
| -------------------- | ------ | ---- | ----------------------------------------------- |
| total                | number |      | total NFT number owned by the user              |
| data                 | array  |      | Game user NFT asset list                        |
| &gt; game_asset_name | string |      | In game NFT asset name                          |
| &gt; token_id        | string |      | NFT asset corresponding TokenID, blank if N/A   |
| &gt; equipment_id    | string |      | In game asset corresponding gear ID             |
| &gt; frozen          | bool   |      | Is the asset frozen or not                      |
| &gt; image           | string |      | This asset(token) corresponding image URL       |
| &gt; attrs           | array  |      | This asset(token) corresponding attributes list |

attrs fileds

| Name            | Type   | Enum | Description     |
| --------------- | ------ | ---- | --------------- |
| attribute_id    | number |      | Attribute ID    |
| attribute_value | string |      | Attribute value |

## Get Single NFT Asset Detail

### **Description**

* To Search details on a specific in game NFT asset of a user

### Path

* /nftToken/getAssetDetail?appId={%d}&equipment_id={%s}

### Method

* GET

> Request example

```json
appId=1&equipment_id=24
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**                         |
| ------------------- | ------------------ | -------------- | -------------- | --------------------------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK                       |
| equipment_id        | true               | string         |                | In game corresponding gear ID                 |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "game_asset_name": "MonsterWeapon",
        "equipment_id": "24",
        "frozen": false,
        "image": "https://hadesathena.github.io/images/333.png",
        "attrs": [{
            "attribute_id": 1,
            "attribute_name":"attrack",
            "attribute_value": "100000000000000"
        }]
    }
}
```

### **Response**

| Name                 | Type   | Enum | Description                                      |
| -------------------- | ------ | ---- | ------------------------------------------------ |
| data                 | object |      | Game user's NFT asset info                       |
| &gt; game_asset_name | string |      | In game asset name                               |
| &gt; equipment_id    | string |      | In game asset corresponding gear ID              |
| &gt; frozen          | bool   |      | is this asset frozen or not                      |
| &gt; image           | string |      | This asset (token) corresponding image URL       |
| &gt; attrs           | array  |      | This asset (token) corresponding attributes list |

attrs fileds

| Name            | Type   | Enum | Description     |
| --------------- | ------ | ---- | --------------- |
| attribute_id    | number |      | Attribute ID    |
| attribute_name  | string |      | Attribute name    |
| attribute_value | string |      | Attribute value |

## NFT Deposit

<font color=#FF0000 >Developer to record NFT desopit info</font>

<font color=#FF0000 >Recommendation - when processing deposit orders, force user offline</font>

### **Description**

* Inform on in game NFT asset deposit

### Path

* /nftToken/deposit

### Method

* POST

> Request example

```json
{
    "appId": 1,
    "sign": "CCYT65...GFR763HYG",
    "params": [{
        "game_asset_name": "MonsterWeapon",
        "token_id": "1111",
        "equipment_id": "222",
        "tx_hash": "0xa0a99d138a84934d9c324eb3d669416c41ac89d723648f022fea3de0ba41d080",
        "uid": 1001,
        "attrs": [{
            "attribute_id": 1,
            "attribute_value": "100000000000000"
        }]
    }]
}

```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**       |
| ------------------- | ------------------ | -------------- | -------------- | --------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK     |
| sign                | true               | string         |                | Signature string of request |
| params              | true               | arrry          |                | NFT deposit info list       |

params fields

| **Parameter** | **Required** | **Type** | **Enum** | **Description**               |
| ------------------- | ------------------ | -------------- | -------------- | ----------------------------------- |
| game_asset_name     | true               | string         |                | In game asset name                  |
| token_id            | true               | string         |                | NFT asset corresponding TokenID     |
| equipment_id        | true               | string         |                | Asset corresponding in game gear ID |
| tx_hash             | true               | string         |                | On chain transaction hash           |
| uid                 | true               | number         |                | Game user ID                        |
| attrs               | true               | array          |                | Attributes list                     |

attrs fileds

| Name            | Type   | Enum | Description     |
| --------------- | ------ | ---- | --------------- |
| attribute_id    | number |      | Attribute ID    |
| attribute_value | string |      | Attribute value |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": [{
        "game_asset_name": "MonsterWeapon",
        "app_order_id": "7239895723",
        "tx_hash": "0xa0a99d138a84934d9c324eb3d669416c41ac89d723648f022fea3de0ba41d080",
        "token_id": "1111",
        "equipment_id": "222",
        "status": 1,
    }]
}
```

### **Response**

| Name                 | Type   | Enum                  | Description                            |
| -------------------- | ------ | --------------------- | -------------------------------------- |
| data                 | array  |                       | In game NFT deposit processing results |
| &gt; game_asset_name | string |                       | In game asset name                     |
| &gt; app_order_id    | string |                       | In game NFT deposit order number       |
| &gt; tx_hash         | string |                       | On chain transaction hash              |
| &gt; token_id        | string |                       | NFT asset corresponding TokenID        |
| &gt; equipment_id    | string |                       | Asset corresponding in game gear ID    |
| &gt; status          | number | 1(success),2(failure) | Status code                            |

## NFT PreWithdraw

<font color=#FF0000 >If successful, will freeze user's corresponding asset</font>

<font color=#FF0000 >Developer to record withdraw order details</font>

<font color=#FF0000 >Recommendation - when processing withdraw orders, force user offline</font>

### **Description**

* NFT asset pre-withdraw

### Path

* /nftToken/preWithdraw

### Method

* POST

> Request example

```json
{
    "appId": 1,
    "sign": "KSA733B...UI907T5FA",
    "params": [{
        "game_asset_name":"MonsterWeapon",
        "uid": 1001,
        "equipment_id": "347",
        "nonce":"24162707295030857711765346710133083490650876608514302741168749846684398572756"
    }]
}

```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**       |
| ------------------- | ------------------ | -------------- | -------------- | --------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK     |
| sign                | true               | string         |                | Signature string of request |
| params              | true               | object         |                | NFT pre-withdraw info list  |

params fields

| **Parameter** | **Required** | **Type** | **Enum** | **Description**              |
| ------------------- | ------------------ | -------------- | -------------- | ---------------------------------- |
| game_asset_name     | true               | string         |                | In game asset name                 |
| uid                 | true               | number         |                | Game user ID                       |
| equipment_id        | true               | string         |                | Asset corrsponding in game gear ID |
| nonce               | true               | string         |                | SDK NFT withdraw order number      |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": [{
        "game_asset_name":"MonsterWeapon",
        "nonce":"24162707295030857711765346710133083490650876608514302741168749846684398572756",
        "app_order_id": "2351235812",
        "uid": 1001,
        "equipment_id": 347,
        "status": 1
    }]
}
```

### **Response**

| Name                 | Type   | Enum                  | Description                                 |
| -------------------- | ------ | --------------------- | ------------------------------------------- |
| data                 | array |                       | In game NFT pre-withdraw processing results |
| &gt; game_asset_name | string |                       | In game NFT asset name                      |
| &gt; nonce           | string |                       | SDK NFT withdraw order number               |
| &gt; app_order_id    | string |                       | Withdraw order number from game             |
| &gt; uid             | number |                       | Game user ID                                |
| &gt; equipment_id    | string |                       | Asset corresponding in game gear ID         |
| &gt; status          | number | 1(success),2(failure) | Status code                                 |

## NFT Withdraw

<font color=#FF0000 >If successful, the corresponding asset will be permanently deleted or restored</font>

<font color=#FF0000 >Developer to update withdraw order status</font>

<font color=#FF0000 >Recommendation - when processing withdraw orders, force user offline</font>


### Description

* NFT asset official withdraw

### Path

* /nftToken/withdraw

### Method

* POST

> Request example

```json
{
    "appId": 1,
    "sign": "KS54W7...08245HGT6",
    "params": [{
        "game_asset_name":"MonsterWeapon",
        "app_order_id": "2351235812",
        "nonce":"24162707295030857711765346710133083490650876608514302741168749846684398572756",
        "uid": 1001,
        "operation": 1
    }]
}
```

### Parameters

| **Parameter** | **Required** | **Type** | **Enum** | **Description**           |
| ------------------- | ------------------ | -------------- | -------------- | ------------------------------- |
| appId               | true               | number         |                | Game ID in ConnexionSDK         |
| sign                | true               | string         |                | Signature string of request     |
| params              | true               | object         |                | NFT offcial withdraw order list |

params fields

| **Parameter** | **Required** | **Type** | **Enum**       | **Description**                         |
| ------------------- | ------------------ | -------------- | -------------------- | --------------------------------------------- |
| game_asset_name     | true               | string         |                      | In game asset name                            |
| app_order_id        | true               | string         |                      | In game NFT withdraw order number             |
| nonce               | true               | string         |                      | SDK NFT withdraw order number                 |
| uid                 | true               | number         |                      | Game user ID                                  |
| operation           | true               | number         | 1(delete),2(recover) | Delete or restore corresponding in game asset |

> The above command returns JSON structured like this

```json
{
    "code": 0,
    "message": "success",
    "data": [{
        "game_asset_name":"MonsterWeapon",
        "app_order_id": "2351235812",
        "nonce":"24162707295030857711765346710133083490650876608514302741168749846684398572756",
        "uid": 1001,
        "status": 1
    }]
}
```

### Response

| Name                 | Type   | Enum                  | Description                                      |
| -------------------- | ------ | --------------------- | ------------------------------------------------ |
| data                 | array  |                       | In game NFT official withdraw processing results |
| &gt; game_asset_name | string |                       | In game asset name                               |
| &gt; app_order_id    | string |                       | In game NFT withdraw order number                |
| &gt; nonce           | string |                       | SDK NFT withdraw order number                    |
| &gt; uid             | number |                       | Game user ID                                     |
| &gt; status          | number | 1(success),2(failure) | Status code                                      |

# Users

## Search for users' in game character information

### **Description**

* Use to search for in game character information of certain players

### Path

* /roleInfo

### Method

* POST

> Sample request:

```json
{
  "uids": [1001,1002]
}
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**          |
| ------------- | ------------ | -------- | -------- | ------------------------ |
| uids          | true         | array    |          | array of in game user ID |

> Output json data as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": {
      	"uid": "2237800",                  //in game user ID
        "roleName": "name",                //character name
        "roleType": "warrior",             //character type
        "level": "33",                     //level
        "militaryCapability": "23580",     //fighting capacity
        "ranking": "10",                   //fighting capacity ranking
        "extraData":  {                    //other information
            "XXX": 1,                      //field 1
            "XXXX": "extraData1",          //field 2
            ...                            //feild ...
        }
    }
}
```

### **Output**

| Name                 | Type   | Enum | Description               |
| -------------------- | ------ | ---- | ------------------------- |
| data                 | array  |      | in game user ID           |
| &gt; roleName        | string |      | character name            |
| &gt; roleType        | string |      | character type            |
| &gt; level           | string |      | level                     |
| > militaryCapability | string |      | fighting capacity         |
| > ranking            | string |      | fighting capacity ranking |
| > extraData          | array  |      | other information         |

extraData field

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | field1      |
| field2 | string |      | Field2      |
| ...    | string |      |             |

# Statistics

## Search for total supply and total airdrop amount of assets

### **Description**

* Use to search for total supply and total airdrop amount of certain assets

### Path

* /statistics/assetsTotal?type={%s}

### Method

* GET

> Sample request:

```json
type="FT"
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------- | ------------ | -------- | -------- | --------------- |
| type          | true         | string   | FT/NFT   | asset type      |

> Output json data as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
        {
          "assetName": "ALG",              //in game asset name
          "totalSupply": "310662",         //total supply
          "totalAirdrop": "655490",        //total airdrop amount
        },
        {
          "assetName": "ALT",              //in game asset name
          "totalSupply": "636",            //total supply
          "totalAirdrop": "0",             //total airdrop amount
        },
        ...                                //asset N
    ]
}
```

### **Output**

| Name              | Type   | Enum | Description                           |
| ----------------- | ------ | ---- | ------------------------------------- |
| data              | array  |      | total supply and total airdrop amount |
| &gt; assetName    | string |      | in game asset name                    |
| &gt; totalSupply  | string |      | total supply amount                   |
| &gt; totalAirdrop | string |      | total airdrop amount                  |

## Search for FT asset statistics

### **Description**

* Search for 24 scale unit of data, arranging in chronological order.

  For instance, when scale unit is 1H, search for the past 24H of FT asset supply and consumption data.

  If the scale unit is 6H, search for 24 scale units of data, with 0, 6, 12, 18 o'clock as benchmark

  Same for scale unit 12H and 1D

### Path

* /statistics/FTToken?scale={%s}

### Method

* GET

> Sample request:

```json
scale="1H"
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum**                              | **Description** |
| ------------- | ------------ | -------- | ------------------------------------- | --------------- |
| scale         | true         | string   | "1H"(1H),"6H"(6H),"12H"(12H),"1D"(1D) | time scale unit |

> output son data as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": {
      	//X time axis
        "xAxis": [1667541600000, 1667563200000, 1667584800000, 1667606400000, 1667628000000, 1667649600000, 1667671200000, 1667692800000, 1667714400000, 1667736000000, 1667757600000, 1667779200000, 1667800800000, 1667822400000, 1667844000000, 1667865600000, 1667887200000, 1667908800000, 1667930400000, 1667952000000, 1667973600000, 1667995200000, 1668016800000, 1668038400000],
        "series": [
            {
              "name": "ALG",		//in game asset name
              "production": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],							 //total production amount
              "consumption": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]						 //total consumption amount
            },
            {
              "name": "ALT",
              "production": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
              "consumption": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
            }
        ]
    }
}
```

### **Output**

| Name        | Type  | Enum | Description                                 |
| ----------- | ----- | ---- | ------------------------------------------- |
| data        | array |      | In game FT asset production and consumption |
| &gt; xAxis  | array |      | X Axis, data array range 24                 |
| &gt; series | array |      | total amount，data array range 24           |

series field

| Name        | Type   | Enum | Description                                   |
| ----------- | ------ | ---- | --------------------------------------------- |
| name        | string |      | in game asset name                            |
| production  | array  |      | total production amount, data array range 24  |
| consumption | array  |      | total consumption amount, data array range 24 |

## Search for NFT asset statistics

### **Description**

* Search for 24 scale unit of data, arranging in chronological order.

  For instance, when scale unit is 1H, search for the past 24H of FT asset supply and consumption data.

  If the scale unit is 6H, search for 24 scale units of data, with 0, 6, 12, 18 o'clock as benchmark

  Same for scale unit 12H and 1D

### Path

* /statistics/NFTToken?assetName={%s}&scale={%s}

### Method

* GET

> Sample request:

```json
assetName="chest"&scale="1H"
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum**                              | **Description**                                     |
| ------------- | ------------ | -------- | ------------------------------------- | --------------------------------------------------- |
| assetName     | false        | string   |                                       | If do not submit this parameter then search for all |
| scale         | true         | string   | "1H"(1H),"6H"(6H),"12H"(12H),"1D"(1D) | time scale unit                                     |

> Output json data as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": {
      	//X time axis
        "xAxis": [1667541600000, 1667563200000, 1667584800000, 1667606400000, 1667628000000, 1667649600000, 1667671200000, 1667692800000, 1667714400000, 1667736000000, 1667757600000, 1667779200000, 1667800800000, 1667822400000, 1667844000000, 1667865600000, 1667887200000, 1667908800000, 1667930400000, 1667952000000, 1667973600000, 1667995200000, 1668016800000, 1668038400000],
        "series": [
            {
              "name": "CommonChest",		//in game asset name
              "production": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],							 				 //total production amount
              "consumption": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]						         //total consumption amount
            },
            {
              "name": "BossChest",
              "production": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
              "consumption": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
            }
        ]
    }
}
```

### **Output**

| Name        | Type  | Enum | Description                                         |
| ----------- | ----- | ---- | --------------------------------------------------- |
| data        | array |      | In game NFT total production and consumption amount |
| &gt; xAxis  | array |      | X axis，data array range 24                         |
| &gt; series | array |      | total amount，data array range 24                   |

series field

| Name        | Type   | Enum | Description                                   |
| ----------- | ------ | ---- | --------------------------------------------- |
| name        | string |      | In game asset name                            |
| production  | array  |      | total production amount, data array range 24  |
| consumption | array  |      | total consumption amount, data array range 24 |

## Search for leaderboard data

### **Description**

* Search for fighting capacity, level and asset leaderboard information

### Path

* /statistics/rankingList?type={%s}&assetName={%s}&count={%d}&sort={%s}

### Method

* GET

> Sample request:

```json
type=2&assetName="potion"&count=100&sort="desc"
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum**                                        | **Description**                                           |
| ------------- | ------------ | -------- | ----------------------------------------------- | --------------------------------------------------------- |
| type          | true         | number   | 0(fighting capacity),1(level),2(asset)          | type                                                      |
| assetName     | false        | string   |                                                 | in game asset name, required input parameter if type is 2 |
| count         | false        | number   |                                                 | default 100                                               |
| sort          | false        | string   | "asc"(ascending order),"desc"(descending order) | dafault "desc"                                            |

> outpu json data as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
        {
          "uid": "12345678",                 //in game user ID
          "roleName": "name",                //character name
          "roleType": "master",              //character type
          "level": "33",                     //level
          "militaryCapability": "23580",     //fighting capacity
          "ranking": "10",                   //fighting capacity ranking
          "asset": "3245",                   //total asset amount
          "extraData": {                     //other data
              "XXX": 1,                      //field1
              "XXXX": "extraData1",          //field2
              ...                            //field...
          }
        },
        {
          "uid": "12345678",
          "roleName": "nameABC",
          "roleType": "master",
          "level": "50",
          "militaryCapability": "13987",
          "ranking": "20",
          "asset": "1000",
          "extraData": {
              "XXX": 1,
              "XXXX": "extraData1",
              ...
          }
        },
        ...
    ]
}
```

### **Output**

| Name                 | Type   | Enum | Description               |
| -------------------- | ------ | ---- | ------------------------- |
| data                 | array  |      | in game user data info    |
| &gt; roleName        | string |      | character name            |
| &gt; roleType        | string |      | character type            |
| &gt; level           | string |      | level                     |
| > militaryCapability | string |      | fighting capacity         |
| > ranking            | string |      | fighting capacity ranking |
| > extraData          | array  |      | other data                |

extraData field

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | field1      |
| field2 | string |      | field2      |
| ...    | string |      |             |

## Search for users' accumulated consumption data

### **Description**

* Search for accumulated consumption data of certain asset in a given timespan

### Path

* /statistics/accumulateConsumption

### Method

* POST

> Sample request:

```json
{
  "assetName": "ALT",
  "beginTime": "1668065809607",
  "endTime": "1668065975683",
  "uids": [1001,1002]
}
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**    |
| ------------- | ------------ | -------- | -------- | ------------------ |
| assetName     | true         | string   |          | in game asset name |
| beginTime     | false        | string   |          | begining timestamp |
| endTime       | false        | string   |          | ending timestamp   |
| uids          | true         | array    |          | user id array      |

> Output json as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
      	{
            "uid": 1001,											 //in game user ID
            "amount": "1220.1234",             //accumulated consumption amount
            "roleName": "nameABC",             //character name
            "roleType": "master",              //character type
            "level": "50",                     //level
            "militaryCapability": "13987",     //fighting capacity
            "ranking": "20",                   //fighting capacity ranking
            "asset": "1000",                   //total amount of asset
            "extraData": {                     //other data
                "XXX": 1,                      //field1
                "XXXX": "extraData1",          //field2
                ...                            //field...
            }
        },
      	{
            "uid": 1001,
            "amount": "1220.1234",
            "roleName": "nameABC",
            "roleType": "master",
            "level": "50",
            "militaryCapability": "13987",
            "ranking": "20",
            "asset": "1000",
            "extraData": {
                "XXX": 1,
                "XXXX": "extraData1",
                ...
            }
        }
      ...
    ]
}
```

### **Output**

| Name                 | Type   | Enum | Description                   |
| -------------------- | ------ | ---- | ----------------------------- |
| data                 | array  |      | accumulated comsumption data  |
| &gt; uid             | string |      | user id                       |
| &gt; amount          | string |      | total accumulated consumption |
| > roleName           | string |      | in game character name        |
| &gt; roleType        | string |      | character type                |
| &gt; level           | string |      | level                         |
| > militaryCapability | string |      | fighting capacity             |
| > ranking            | string |      | fighting capacity ranking     |
| > extraData          | array  |      | other data                    |

extraData fields

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | field1      |
| field2 | string |      | field2      |
| ...    | string |      |             |

## Search for asset flow data

### **Description**

* Search for new asset flow data in a given timespan, arranging in descending order of asset amount.

  for instance, if the target search timespan is 2022-11-14 16:23:30，then 1H is 2022-11-14 15:00:00～2022-11-14 16:00:00。

  6H：12AM，6AM，12PM，6PM

  12H：12AM，12PM

  1D：12AM

  search for the data between most recent time scale unit and one unit backward

### Path

* /statistics/asset/flow?assetName={%s}&scale={%s}&count={%d}

### Method

* GET

> Sample request:

```json
{
  "assetName": "ALT",
  "beginTime": "1668065809607",
  "endTime": "1668065975683"
}
```

### **Parameters**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**                  |
| ------------- | ------------ | -------- | -------- | -------------------------------- |
| assetName     | true         | string   |          | in game asset name               |
| scale         | true         | string   |          | time scale unit：1H，6H，12H，1D |
| count         | false        | number   |          | default 100                      |

> Output json data as follow:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
        {
          "uid": "11115678",                 //in game user ID
          "roleName": "nameABC",             //character name
          "roleType": "master",              //character type
          "level": "50",                     //level
          "militaryCapability": "13987",     //fighting capacity
          "ranking": "20"，                  //fighting capacity ranking
          "asset": "1000",                   //total asset amount
          "extraData": {                     //other data
              "XXX": 1,                      //field1
              "XXXX": "extraData1",          //field2
              ...                            //field...
          }
          "amount": "1101.1234"              //amount
        },
        {
          "uid": "22225678",                 //in game user ID
          "roleName": "nameABC",             //character name
          "roleType": "master",              //character type
          "level": "50",                     //level
          "militaryCapability": "13987",     //fighting capacity
          "ranking": "20",                   //fighting capacity ranking
          "asset": "1000",                   //total asset amount
          "extraData": {                     //other data
              "XXX": 1,                      //field1
              "XXXX": "extraData1",          //field2
              ...                            //field...
          }
          "amount": "2202.5678"              //amount
        },
        ...
    ]
}
```

### **Output**

| Name                 | Type   | Enum | Description               |
| -------------------- | ------ | ---- | ------------------------- |
| data                 | array  |      | asset flow data           |
| &gt; uid             | string |      | user id                   |
| &gt; amount          | string |      | amount of asset           |
| > roleName           | string |      | character name            |
| &gt; roleType        | string |      | character type            |
| &gt; level           | string |      | level                     |
| > militaryCapability | string |      | fighting capacity         |
| > ranking            | string |      | fighting capacity ranking |
| > extraData          | array  |      | other data                |

extraData fields

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | field1      |
| field2 | string |      | field2      |
| ...    | string |      |             |

# Backend Management

## User Status Change

### Description

* Changing user's status

### PATH

* /userStatus

### Method

* POST

### Content-Type

* application/x-www-form-urlencoded

> Request example

```curl
curl --location -g --request POST '${gameServerApiUrl}/userStatus' \--header 'Content-Type: application/x-www-form-urlencoded' \--data-urlencode 'appID=1' \--data-urlencode 'uid=1' \--data-urlencode 'status=10' \--data-urlencode 'timestamp=1645415419000' \--data-urlencode 'sign=xxxxxxxxxxxxxxx'
```

### Parameters

| **Parameter** | **Required** | **Type** | **Enum**                                                         | **Description**                              |
| ------------------- | ------------------ | -------------- | ---------------------------------------------------------------------- | -------------------------------------------------- |
| appId               | true               | string         |                                                                        | Game ID in Connexion SDK                           |
| uid                 | true               | string         |                                                                        | Game user ID                                       |
| status              | true               | string         | 0-pause withdraw，1-resume withdraw，2-ban，10-force offline，11-unban | changing user status                               |
| timestamp           | true               | string         |                                                                        | Timestamp at the time of communication, unit in ms |
| sign                | true               | string         |                                                                        | Signature                                          |

> The above command returns JSON structured like this

```
SUCCESS
```

### sign generation method

1. Other than the sign field, rank other fields in increasing order of key value， as key=value&key=value.....to create a string
2. Add &secretKey={SDK Keyparams} to the end of the string generated in step 1 to generate string to-be-signed. Substitute {SDK Keyparams} with AppSecret Parameters
3. Conduct md5 calculation(32 digits capital) to the string to-be-signed in step 2, generating corresponding sign

> Java sign example

```java
public static String sign(Map<String, String> params, String secret) {
        Map<String, String> sortedParams = new TreeMap<>(params);

        StringBuilder sb = StringUtils.getStringBuilder();
        for (Map.Entry<String, String> param : sortedParams.entrySet()) {
            String value = param.getValue();
            String key = param.getKey();
            if ("sign".equalsIgnoreCase(key)) {
                continue;
            }

            if (value != null && value.length() > 0) {
                sb.append(param.getKey()).append("=").append(value).append("&");
            }
        }
        sb.append("secretKey").append("=").append(secret);

        String signStr = sb.toString();

        logger.debug("SignString: {}", signStr);

        String sign = EncryptUtils.md5(signStr).toUpperCase();

        return sign;
    }
```

### Response

* When game server receive request, confirmed signature and status change, then return text string: SUCCESS; otherwise, return text string: FAIL

# Mobile Client

* ConnexionSDK offers Mobile Client function interfaces such as *.framework and *.aar([Download](https://github.com/Connector-Gamefi/mobile-sdk)),details are as such.

* For Android: XPlatform.getInstance() will return a singleton which could be used to call functions provided by aar package.


## IOS Setup illustrations

### Illustration 1

* iOS need to delete SceneDelegate,detailed methods are:

> Delete from AppDelegate

```
pragma mark - UISceneSession lifecycle

(UISceneConfiguration *)application:(UIApplication *)application configurationForConnectingSceneSession:(UISceneSession *)connectingSceneSession options:(UISceneConnectionOptions *)options {
// Called when a new scene session is being created.
// Use this method to select a configuration to create the new scene with.
return [[UISceneConfiguration alloc] initWithName:@"Default Configuration" sessionRole:connectingSceneSession.role];
}

(void)application:(UIApplication *)application didDiscardSceneSessions:(NSSet<UISceneSession *> *)sceneSessions {
// Called when the user discards a scene session.
// If any sessions were discarded while the application was not running, this will be called shortly after application:didFinishLaunchingWithOptions.
// Use this method to release any resources that were specific to the discarded scenes, as they will not return.
}
```


> Delete from  Info.plist

```
`<key>`UIApplicationSceneManifest`</key>`
`<dict>`
    `<key>`UIApplicationSupportsMultipleScenes`</key>`
    `<false/>`
    `<key>`UISceneConfigurations`</key>`
    `<dict>`
        `<key>`UIWindowSceneSessionRoleApplication`</key>`
        `<array>`
            `<dict>`
                `<key>`UISceneConfigurationName`</key>`
                `<string>`Default Configuration`</string>`
                `<key>`UISceneDelegateClassName`</key>`
                `<string>`SceneDelegate`</string>`
                `<key>`UISceneStoryboardFile`</key>`
                `<string>`Main`</string>`
            `</dict>`
        `</array>`
    `</dict>`
`</dict>`
```



  1. Delete relevant codes for UISceneSession lifecyle within AppDelegate
  2. Delete both SceneDelegate.h and SceneDelegate.m files
  3. Detele Application Scene Manifest within Info.plist

### Illustration 2

* set target -> Build Settings -> Linking ->Other Linker Flags -all_load

## IOS-SDK Init

### Description

* IOS SDK Initialization

### Method

* initWithAppID

> example

```
UG_SDKParams* sdkParams = [[UG_SDKParams alloc] initWithAppID:kAppID appKey:kAppKey orientation:@"landscape";
```

### Parameters

| **Parameter** | **Required** | **Type** | **Enum**         | **Description**   |
| ------------------- | ------------------ | -------------- | ---------------------- | ----------------------- |
| appID               | true               | string         |                        | Game ID in ConnexionSDK |
| appKey              | true               | string         |                        | appKey                  |
| orientation         | true               | string         | "portrait","landscape" | portrait or landscape   |

## IOS-SDK login

### Description

* IOS login

### Method

* sharedInstance.login

> example

```
[[UGSDKPlatform sharedInstance] login];
```

## IOS-SDK logOut

### Description

* IOS logout

### Method

* sharedInstance.logout

> example

```
[[UGSDKPlatform sharedInstance] logout];
```

## IOS-SDK submit

### Description

* IOS submit

### Method

* sharedInstance.submitGameData

> example

```
[[UGSDKPlatform sharedInstance] submitGameData:role];
```

### Parameters

| **Parameter** | **Required** | **Type** | **Enum**   | **Description**     |
| ------------------- | ------------------ | -------------- | ---------------- | ------------------------- |
| role                | true               | *UG_GameRole   |                  |                           |
| &gt;type            | true               | string         | 1-create character,2-enter game,3-level up,4-exit game,others | report submission type                |
| &gt;roleID          | true               | string         |                  | Chracter ID               |
| &gt;roleName        | true               | string         |                  | Character name            |
| &gt;roleLevel       | true               | string         |                  | Character level           |
| &gt;serverID        | true               | string         |                  | Server ID                 |
| &gt;serverName      | true               | string         |                  | Server name               |
| &gt;vip             | true               | string         |                  | VIP level                 |
| &gt;createTime      | true               | string         |                  | Timestamp,Unit s          |
| &gt;lastLevelUpTime | true               | string         |                  | Last level up time,Unit s |
| &gt;extraData       | true               | string         |                  | Other extended data       |

## Android-SDK Init

### Description

* Initialization process should be called as game starts, usually called from onCreate in game Activity

### Method

* init

> example

```java
public void init(Activity context, UInitParams params, IInitListener listener)
```

## Android-SDK SetLogoutListener

### Description

* To set up click to switch account function in floating window, game will receive this callback, need to go back to game login interface and guide player to re-login

### Method

* setLogoutListener

> example

```java
public void setLogoutListener(ILogoutListener logoutListener)
```

## Android-SDK login

### Description

* Login interface, once called, popping out SDK login interface

### Method

* login

> example

```java
public void login(Activity activity, ILoginListener listener)
```

## Android-SDK logout

### Description

* Logout interface, once called, clear current login status.

### Method

* logout

> example

```java
public void logout(Activity activity)
```

## Android-SDK submit

### Description

* Submit data, should be called in events such as character creation, enter game, level up, exiting game.

### Method

* submit

> example

```java
public void submit(Activity activity, URoleData roleData)
```

### Parameters

| **Parameter** | **Required** | **Type** | **Enum**                                                | **Description**                 |
| ------------------- | ------------------ | -------------- | ------------------------------------------------------------- | ------------------------------------- |
| roleData            | true               | URoleData      |                                                               |                                       |
| &gt;type            | true               | string         | 1-create character,2-enter game,3-level up,4-exit game,others | report submission type                |
| &gt;roleID          | true               | string         |                                                               | Character ID                          |
| &gt;roleName        | true               | string         |                                                               | Character name                        |
| &gt;roleLevel       | true               | string         |                                                               | Character level                       |
| &gt;serverID        | true               | string         |                                                               | Game server ID                        |
| &gt;serverName      | true               | string         |                                                               | Game server name                      |
| &gt;vip             | true               | string         |                                                               | VIP level                             |
| &gt;createTime      | true               | string         |                                                               | Timestamp, Unit s                     |
| &gt;lastLevelUpTime | true               | string         |                                                               | Last level up time, Unit s            |
| &gt;extraData       | true               | string         |                                                               | json format,extended data,current N/A |

## Android-SDK exit

### Description

* When exiting game, pop out exit confirmation box

### Method

* exit

> example

```java
public void exit(final Activity context)
```

## Android-SDK life-cycle approach

* onStart onStop   onResume   onPause   onDestroy
