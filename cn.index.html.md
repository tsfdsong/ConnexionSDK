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
# 介绍

***ConnexionSDK*** ***API v1.0.0***

# JSON HTTP 请求

* 为了接入ConnexionSDK, 开发者使用JSON with HTTP实现本文档如下接口
* 游戏方应该限制访问IP(只允许SDK部署的服务器IP)

## 请求信息

* 请求头Content-Type为application/json
* API接口请求参数中包含签名sign,params,appId字段

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在SDK中的标识ID   |
| sign                | true               | string         |                | 签名信息              |
| params              | true               | array,object   |                | 请求参数              |

* 签名用于与ConnexionSDK之间的安全通信，开发者需要校验参数与签名数据的一致性。签名方法如下:

> Golang sign 示例:

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

1. 将入参params转为json字符串
2. 待签名字符串为:appId={入参appId}&params={第一步json字符串}&secretKey={游戏在管理后台设置的SecketKey}
3. 对待签名字符串md5生成对应的sign(大写)

## 响应信息

响应如下：

| **Parameter** | **Required** | **Type** | **Enum** | **Description**    |
| ------------------- | ------------------ | -------------- | -------------- | ------------------------ |
| code                | true               | number         |                | 0(success),others(error) |
| message             | true               | string         |                | 错误信息                 |
| data                | true               | array,object   |                | 返回数据                 |

# 资产

## 查询用户在游戏内的FT资产

### **描述**

* 查询用户在游戏内的FT资产

### 路径

* /ftToken/getAssets?appId={%d}&uid={%d}

### 方法

* GET

> 请求示例:

```json
appId=1&uid=1001
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| uid                 | true               | number         |                | 游戏用户ID            |

> 返回json数据如下:

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

### **响应**

| Name                     | Type   | Enum | Description              |
| ------------------------ | ------ | ---- | ------------------------ |
| data                     | array  |      | 游戏用户FT资产列表       |
| &gt; game_coin_name      | string |      | 游戏资产名               |
| &gt; coin_balance        | string |      | 游戏内该资产可用资产数量 |
| &gt; coin_frozen_balance | string |      | 游戏内该资产冻结资产数量 |

## FT资产充值

### 注意

<font color=#FF0000 >由开发者记录充值订单信息</font>

<font color=#FF0000 >推荐-处理充值订单时，强制用户下线</font>

### **描述**

* FT资产充值

### 路径

* /ftToken/deposit

### 方法

* POST

> 请求示例:

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

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | arrry          |                | FT充值信息列表        |

params 字段

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| game_coin_name      | true               | string         |                | 游戏资产名            |
| amount              | true               | string         |                | FT资产充值数量        |
| tx_hash             | true               | string         |                | 链上交易hash          |
| uid                 | true               | number         |                | 游戏用户ID            |

> 返回json数据如下:

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

### **响应**

| Name                | Type   | Enum                  | Description        |
| ------------------- | ------ | --------------------- | ------------------ |
| data                | array  |                       | 游戏FT充值处理列表 |
| &gt; game_coin_name | string |                       | 游戏资产名         |
| &gt; app_order_id   | string |                       | 游戏充值订单号     |
| &gt; tx_hash        | string |                       | 链上交易hash       |
| &gt; status         | number | 1(success),2(failure) | 状态码             |

## FT资产预提现

### 注意

<font color=#FF0000 >如果成功 则应冻结用户相应资产</font>

<font color=#FF0000 >由开发者记录提现订单信息</font>

<font color=#FF0000 >推荐-处理预提现请求时，强制用户下线</font>

### **描述**

* FT资产预提现

### 路径

* /ftToken/preWithdraw

### 方法

* POST

> 请求示例:

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

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | object         |                | FT预提现信息          |

params 字段

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| game_coin_name      | true               | string         |                | 游戏内资产名          |
| amount              | true               | string         |                | FT资产预提现数量      |
| uid                 | true               | number         |                | 游戏用户ID            |
| nonce               | true               | string         |                | SDK提现订单号         |

> 返回json数据如下:

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

### **响应**

| Name              | Type   | Enum                  | Description          |
| ----------------- | ------ | --------------------- | -------------------- |
| data              | object |                       | 游戏FT预提现处理结果 |
| &gt; app_order_id | string |                       | 游戏FT提现订单号     |
| &gt; status       | number | 1(success),2(failure) | 状态码               |

## FT资产正式提现

<font color=#FF0000 >如果成功 则彻底删除或者恢复用户相应资产</font>

<font color=#FF0000 >开发者更新提现订单信息</font>

<font color=#FF0000 >推荐-处理正式提现请求时，强制用户下线</font>

### **描述**

* FT资产正式提现

### 路径

* /ftToken/withdraw

### 方法

* POST

> 请求示例:

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

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | array          |                | FT正式提现列表        |

params 字段

| **Parameter** | **Required** | **Type** | **Enum**       | **Description** |
| ------------------- | ------------------ | -------------- | -------------------- | --------------------- |
| game_coin_name      | true               | string         |                      | 游戏资产名            |
| app_order_id        | true               | string         |                      | 游戏FT提现订单号      |
| nonce               | true               | string         |                      | 提现订单号            |
| operation           | true               | number         | 1(delete),2(recover) | 操作                  |
| uid                 | true               | number         |                      | 游戏用户ID            |

> 返回json数据如下:

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

### **响应**

| Name                | Type   | Enum                  | Description            |
| ------------------- | ------ | --------------------- | ---------------------- |
| data                | array  |                       | 游戏FT正式提现处理结果 |
| &gt; game_coin_name | string |                       | 游戏资产名             |
| &gt; app_order_id   | string |                       | 游戏FT提现订单号       |
| &gt; nonce          | string |                       | SDK FT提现订单号       |
| &gt; status         | number | 1(success),2(failure) | 状态码                 |

## 查询用户在游戏内的NFT资产

### **描述**

* 查询用户在游戏内的NFT资产

### 路径

* /nftToken/getAssets?appId={%d}&uid={%d}&page={%d}&pageSize={%d}

### 方法

* GET

> 请求示例:

```json
appId=1&uid=1001&page=1&pageSize=30
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**                    |
| ------------------- | ------------------ | -------------- | -------------- | ---------------------------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID                        |
| uid                 | true               | number         |                | 游戏用户ID                               |
| page                | true               | number         |                | 资产列表页码                             |
| pageSize            | true               | number         |                | 资产列表单页数量                         |
| assetName           | true               | string         |                | 游戏内资产名称，为空时，表示查询全部资产 |

> 返回json数据如下:

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
      "attribute_name": "attack +1"
      "attribute_value": "100000000000000"
    }]
  }]
}
```

### **响应**

| Name                 | Type   | Enum | Description                                |
| -------------------- | ------ | ---- | ------------------------------------------ |
| total                | number |      | 游戏用户NFT资产总数                        |
| data                 | array  |      | 游戏用户NFT资产列表                        |
| &gt; game_asset_name | string |      | 游戏资产名                                 |
| &gt; token_id        | string |      | NFT资产对应的合约内TokenID(可以为空字符串) |
| &gt; equipment_id    | string |      | 游戏内资产对应的装备ID                     |
| &gt; frozen          | bool   |      | 该资产是否被冻结                           |
| &gt; image           | string |      | 该资产(token)对应的图片地址                |
| &gt; attrs           | array  |      | 该资产(token)对应的属性列表                |

attrs 字段

| Name            | Type   | Enum | Description |
| --------------- | ------ | ---- | ----------- |
| attribute_id    | number |      | 属性id      |
| attribute_value | string |      | 属性值      |
| attribute_name  | string |      | 属性名称    |

## 查询单个NFT资产详情

### **描述**

* 查询游戏内单个NFT资产详情

### 路径

* /nftToken/getAssetDetail?appId={%d}&equipment_id={%s}

### 方法

* GET

> 请求示例:

```json
appId=1&equipment_id=24&skip_validate=1
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum**                | **Description**  |
| ------------------- | ------------------ | -------------- | ----------------------------- | ---------------------- |
| appId               | true               | number         |                               | 游戏在sdk中标识ID      |
| equipment_id        | true               | string         |                               | 游戏内资产对应的装备ID |
| skip_validate       | false              | number         | 1(跳过校验) 0或者不传(不跳过) | 是否跳过游戏校验       |

> 返回json数据如下:

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
      "attribute_value": "100000000000000"
    }]
  }
}
```

### **响应**

| Name                 | Type   | Enum | Description                     |
| -------------------- | ------ | ---- | ------------------------------- |
| data                 | object |      | 游戏用户NFT资产信息             |
| &gt; game_asset_name | string |      | 游戏资产名称                    |
| &gt; equipment_id    | string |      | 游戏内资产对应的装备ID          |
| &gt; frozen          | bool   |      | 该资产是否被冻结                |
| &gt; image           | string |      | 该资产(token)对应的图片名或地址 |
| &gt; attrs           | array  |      | 该资产(token)对应的属性列表     |

attrs 字段

| Name            | Type   | Enum | Description  |
| --------------- | ------ | ---- | ------------ |
| attribute_name  | string |      | 属性显示名称 |
| attribute_id    | number |      | 属性id       |
| attribute_value | string |      | 属性值       |

## NFT资产充值

<font color=#FF0000 >由开发者记录充值订单信息</font>

<font color=#FF0000 >推荐-处理充值订单时，强制用户下线</font>

### **描述**

* NFT资产充值

### 路径

* /nftToken/deposit

### 方法

* POST

> 请求示例:

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

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | arrry          |                | NFT充值信息列表       |

params 字段

| **Parameter** | **Required** | **Type** | **Enum** | **Description**      |
| ------------------- | ------------------ | -------------- | -------------- | -------------------------- |
| game_asset_name     | true               | string         |                | 游戏资产名                 |
| token_id            | true               | string         |                | NFT资产对应的合约内TokenID |
| equipment_id        | true               | string         |                | 资产对应游戏的最新装备ID   |
| tx_hash             | true               | string         |                | 链上交易hash               |
| uid                 | true               | number         |                | 游戏用户ID                 |
| attrs               | true               | array          |                | 属性列表                   |

attrs 字段

| Name            | Type   | Enum | Description |
| --------------- | ------ | ---- | ----------- |
| attribute_id    | number |      | 属性id      |
| attribute_value | string |      | 属性值      |

> 返回json数据如下:

```json
{
  "code": 0,
  "message": "success",
  "data": [{
    "game_asset_name": "MonsterWeapon",
    "app_order_id": "7239895723",
    "tx_hash": "0xa0a99d138a84934d9c324eb3d669416c41ac89d723648f022fea3de0ba41d080",
    "token_id": "1111",
    "equipment_id": 222,
    "status": 1,
  }]
}
```

### **响应**

| Name                 | Type   | Enum                  | Description                |
| -------------------- | ------ | --------------------- | -------------------------- |
| data                 | array  |                       | 游戏NFT充值处理结果列表    |
| &gt; game_asset_name | string |                       | 游戏资产名                 |
| &gt; app_order_id    | string |                       | 游戏NFT充值订单号          |
| &gt; tx_hash         | string |                       | 链上交易hash               |
| &gt; token_id        | string |                       | NFT资产对应的合约内TokenID |
| &gt; equipment_id    | string |                       | 资产对应游戏的最新装备ID   |
| &gt; status          | number | 1(success),2(failure) | 状态码                     |

## NFT资产预提现

<font color=#FF0000 >如果成功 则应冻结用户相应资产</font>

<font color=#FF0000 >由开发者记录提现订单信息</font>

<font color=#FF0000 >推荐-处理与提现订单时，强制用户下线</font>

### **描述**

* NFT资产预提现

### 路径

* /nftToken/preWithdraw

### 方法

* POST

> 请求示例:

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

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | object         |                | NFT预提现信息列表     |

param 字段

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| game_asset_name     | true               | string         |                | 游戏资产名            |
| uid                 | true               | number         |                | 游戏用户ID            |
| equipment_id        | true               | string         |                | 资产对应的最新装备ID  |
| nonce               | true               | string         |                | SDK提现订单号         |

> 返回json数据如下:

```json
{
  "code": 0,
  "message": "success",
  "data": [{
    "game_asset_name":"MonsterWeapon",
    "nonce":"24162707295030857711765346710133083490650876608514302741168749846684398572756",
    "app_order_id": "2351235812",
    "uid": 1001,
    "equipment_id": "347",
    "status": 1
  }]
}
```

### **响应**

| Name                 | Type   | Enum                  | Description               |
| -------------------- | ------ | --------------------- | ------------------------- |
| data                 | array  |                       | 游戏NFT预提现处理结果列表 |
| &gt; game_asset_name | string |                       | 游戏资产名                |
| &gt; nonce           | string |                       | SDK提现订单号             |
| &gt; app_order_id    | string |                       | 游戏NFT提现订单号         |
| &gt; uid             | number |                       | 游戏用户ID                |
| &gt; equipment_id    | string |                       | 资产对应的最新装备ID      |
| &gt; status          | number | 1(success),2(failure) | 状态码                    |

## NFT资产正式提现

<font color=#FF0000 >如果成功 则彻底删除或者恢复用户相应资产</font>

<font color=#FF0000 >由开发者更新提现订单状态</font>

<font color=#FF0000 >推荐-处理体现订单时，强制用户下线处理</font>

### 描述

* NFT资产正式提现

### 路径

* /nftToken/withdraw

### 方法

* POST

> 请求示例:

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

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | object         |                | NFT正式提现列表       |

param 字段

| **Parameter** | **Required** | **Type** | **Enum**       | **Description** |
| ------------------- | ------------------ | -------------- | -------------------- | --------------------- |
| game_asset_name     | true               | string         |                      | 游戏资产名            |
| app_order_id        | true               | string         |                      | 游戏NFT提现订单号     |
| nonce               | true               | string         |                      | SDK NFT提现订单号     |
| uid                 | true               | number         |                      | 游戏用户ID            |
| operation           | true               | number         | 1(delete),2(recover) | 操作                  |

> 返回json数据如下:

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

### **响应**

| Name                 | Type   | Enum                  | Description                 |
| -------------------- | ------ | --------------------- | --------------------------- |
| data                 | array  |                       | 游戏NFT正式提现处理结果列表 |
| &gt; game_asset_name | string |                       | 游戏资产名                  |
| &gt; app_order_id    | string |                       | 游戏NFT提现订单号           |
| &gt; nonce           | string |                       | SDK NFT提现订单号           |
| &gt; uid             | number |                       | 游戏用户ID                  |
| &gt; status          | number | 1(success),2(failure) | 状态码                      |

## 客户端支付游戏商品充值

### **描述**

* ApplePay或GooglePay支付成功之后给用户充值游戏资产接口

### 路径

* /payOrder/deposit

### 方法

* POST

> 请求示例:

```json
{
    "appId": 1,
    "sign": "2YA43T0...55YGF5AH",
    "params": {
        "productId": "1234",
        "orderId": "0455513284",
        "quantity": 1,
        "uid": 1001
    }
}

```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| appId               | true               | number         |                | 游戏在sdk中标识ID     |
| sign                | true               | string         |                | 请求附带的签名        |
| params              | true               | object         |                | 参数列表              |

params 字段

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| productId           | true               | string         |                | 商品id                |
| orderId             | true               | string         |                | 充值订单号            |
| quantity            | true               | number         |                | 数量                  |
| uid                 | true               | number         |                | 游戏用户ID            |

> 返回json数据如下:

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "productId": "1234",
        "orderId": "0455513284",
        "quantity": 1,
        "uid": 1001,
        "status": 1
    }
}
```

### **响应**

| Name           | Type   | Enum                  | Description |
| -------------- | ------ | --------------------- | ----------- |
| data           | object |                       | 处理结果    |
| &gt; productId | string |                       | 商品id      |
| &gt; orderId   | string |                       | 充值订单号  |
| &gt; quantity  | number |                       | 数量        |
| &gt; uid       | number |                       | 游戏用户ID  |
| &gt; status    | number | 1(success),2(failure) | 状态码      |

# 用户

## 查询用户在游戏内角色信息

### **描述**

* 查询用户在游戏内角色信息

### 路径

* /roleInfo

### 方法

* POST

> 请求示例:

```json
{
  "uids": [1001, 1002]
}
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| uids                | true               | array          |                | 游戏用户ID数组        |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": {
      	"uid": 2237800,                    //游戏用户ID
        "roleName": "name",                //角色名称
        "roleType": "warrior",             //角色类型
        "level": "33",                     //等级
        "militaryCapability": "23580",     //战力
        "ranking": "10",                   //战力排行
        "extraData":  {                    //其他额外数据
            "XXX": 1,                      //字段1
            "XXXX": "extraData1",          //字段2
            ...                            //字段...
        }
    }
}
```

### **响应**

| Name                 | Type   | Enum | Description  |
| -------------------- | ------ | ---- | ------------ |
| data                 | array  |      | 游戏用户信息 |
| > uid                | number |      | 游戏用户ID   |
| &gt; roleName        | string |      | 角色名称     |
| &gt; roleType        | string |      | 角色类型     |
| &gt; level           | string |      | 等级         |
| > militaryCapability | string |      | 战力         |
| > ranking            | string |      | 战力排名     |
| > extraData          | array  |      | 其他信息     |

extraData 字段

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | 字段1       |
| field2 | string |      | 字段2       |
| ...    | string |      |             |

# 统计

## 查询资产总产出及总空投

### **描述**

* 查询资产总产出及总空投量

### 路径

* /statistics/assetsTotal?type={%s}

### 方法

* GET

> 请求示例:

```json
type="FT"
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| type                | true               | string         | FT/NFT         | 资产类型              |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
        {
          "assetName": "ALG",              //游戏内资产名称
          "totalSupply": "310662",         //总产出量
          "totalAirdrop": "655490",        //总空投量
        },
        {
          "assetName": "ALT",              //游戏内资产名称
          "totalSupply": "636",            //总产出量
          "totalAirdrop": "0",             //总空投量
        },
        ...                                //资产N
    ]
}
```

### **响应**

| Name              | Type   | Enum | Description      |
| ----------------- | ------ | ---- | ---------------- |
| data              | array  |      | 总产出及总空投量 |
| &gt; assetName    | string |      | 游戏内资产名称   |
| &gt; totalSupply  | string |      | 总产出量         |
| &gt; totalAirdrop | string |      | 总空投量         |

## 查询FT资产统计数据

### **描述**

* 查询24个刻度单位的统计数据，按时间正序排列。

  例如刻度为1H，查询当前时间往前推24H的FT资产的产出及消耗数据。

  如刻度为6H，以每天0，6，12，18时数据为基础，取24个刻度单位数据。

  刻度为12H及1D同理。

### 路径

* /statistics/FTToken?scale={%s}

### 方法

* GET

> 请求示例:

```json
scale="1H"
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum**                                  | **Description** |
| ------------------- | ------------------ | -------------- | ----------------------------------------------- | --------------------- |
| scale               | true               | string         | "1H"(1小时),"6H"(6小时),"12H"(12小时),"1D"(1天) | 时间刻度              |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": {
      	//X时间轴
        "xAxis": [1667541600000, 1667563200000, 1667584800000, 1667606400000, 1667628000000, 1667649600000, 1667671200000, 1667692800000, 1667714400000, 1667736000000, 1667757600000, 1667779200000, 1667800800000, 1667822400000, 1667844000000, 1667865600000, 1667887200000, 1667908800000, 1667930400000, 1667952000000, 1667973600000, 1667995200000, 1668016800000, 1668038400000],
        "series": [
            {
              "name": "ALG",		//游戏内资产名称
              "production": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],							 //产出量
              "consumption": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]						 //消耗量
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

### **响应**

| Name        | Type  | Enum | Description              |
| ----------- | ----- | ---- | ------------------------ |
| data        | array |      | 游戏内FT资产产出及消耗量 |
| &gt; xAxis  | array |      | X时间轴，数组范围24      |
| &gt; series | array |      | 总量，数组范围24         |

series 字段

| Name        | Type   | Enum | Description        |
| ----------- | ------ | ---- | ------------------ |
| name        | string |      | 游戏内资产名称     |
| production  | array  |      | 产出量，数组范围24 |
| consumption | array  |      | 消耗量，数组范围24 |

## 查询NFT资产统计数据

### **描述**

* 查询24个刻度单位的统计数据，按时间正序排列。

  例如刻度为1H，查询当前时间往前推24H的NFT资产的产出及消耗数据。

  如刻度为6H，以每天0，6，12，18时数据为基础，取24个刻度单位数据。

  刻度为12H及1D同理。

### 路径

* /statistics/NFTToken?assetName={%s}&scale={%s}

### 方法

* GET

> 请求示例:

```json
assetName="chest"&scale="1H"
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum**                                  | **Description** |
| ------------------- | ------------------ | -------------- | ----------------------------------------------- | --------------------- |
| assetName           | false              | string         |                                                 | 该参数不传则查询所有  |
| scale               | true               | string         | "1H"(1小时),"6H"(6小时),"12H"(12小时),"1D"(1天) | 时间刻度              |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": {
      	//X时间轴
        "xAxis": [1667541600000, 1667563200000, 1667584800000, 1667606400000, 1667628000000, 1667649600000, 1667671200000, 1667692800000, 1667714400000, 1667736000000, 1667757600000, 1667779200000, 1667800800000, 1667822400000, 1667844000000, 1667865600000, 1667887200000, 1667908800000, 1667930400000, 1667952000000, 1667973600000, 1667995200000, 1668016800000, 1668038400000],
        "series": [
            {
              "name": "CommonChest",		//游戏内资产名称
              "production": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],							 				 //产出量
              "consumption": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]						         //消耗量
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

### **响应**

| Name        | Type  | Enum | Description              |
| ----------- | ----- | ---- | ------------------------ |
| data        | array |      | 游戏内FT资产产出及消耗量 |
| &gt; xAxis  | array |      | X时间轴，数组范围24      |
| &gt; series | array |      | 总量，数组范围24         |

series 字段

| Name        | Type   | Enum | Description        |
| ----------- | ------ | ---- | ------------------ |
| name        | string |      | 游戏内资产名称     |
| production  | array  |      | 产出量，数组范围24 |
| consumption | array  |      | 消耗量，数组范围24 |

## 查询排行榜数据

### **描述**

* 查询战力，等级，资产排行榜数据信息

### 路径

* /statistics/rankingList?type={%s}&assetName={%s}&count={%d}&sort={%s}

### 方法

* GET

> 请求示例:

```json
type=2&assetName="potion"&count=100&sort="desc"
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum**           | **Description**                   |
| ------------------- | ------------------ | -------------- | ------------------------ | --------------------------------------- |
| type                | true               | number         | 0(战力),1(等级),2(资产)  | 类型                                    |
| assetName           | false              | string         |                          | 游戏内资产名称，type传2时必须传入该参数 |
| count               | false              | number         |                          | 默认100                                 |
| sort                | false              | string         | "asc"(升序),"desc"(降序) | 默认"desc"                              |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
        {
          "uid": 12345678,                   //游戏用户ID
          "roleName": "name",                //角色名称
          "roleType": "master",              //角色类型
          "level": "33",                     //等级
          "militaryCapability": "23580",     //战力
          "ranking": "10",                   //战力排行
          "asset": "3245",                   //资产总量
          "extraData": {                     //其他额外数据
              "XXX": 1,                      //字段1
              "XXXX": "extraData1",          //字段2
              ...                            //字段...
          }
        },
        {
          "uid": 12345678,
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

### **响应**

| Name                 | Type   | Enum | Description  |
| -------------------- | ------ | ---- | ------------ |
| data                 | array  |      | 游戏用户信息 |
| > uid                | number |      | 游戏用户ID   |
| &gt; roleName        | string |      | 角色名称     |
| &gt; roleType        | string |      | 角色类型     |
| &gt; level           | string |      | 等级         |
| > militaryCapability | string |      | 战力         |
| > ranking            | string |      | 战力排名     |
| > extraData          | array  |      | 其他信息     |

extraData 字段

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | 字段1       |
| field2 | string |      | 字段2       |
| ...    | string |      |             |

## 查询用户累积消耗数据

### **描述**

* 查询一定时间内累积消耗资产数据信息

### 路径

* /statistics/accumulateConsumption

### 方法

* POST

> 请求示例:

```json
{
  "assetName": "ALT",
  "beginTime": "1668065809607",
  "endTime": "1668065975683",
  "uids": [1001,1002]
}
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description** |
| ------------------- | ------------------ | -------------- | -------------- | --------------------- |
| assetName           | true               | string         |                | 游戏内资产名称        |
| beginTime           | false              | string         |                | 开始时间时间戳        |
| endTime             | false              | string         |                | 截止时间时间戳        |
| uids                | true               | array          |                | 游戏用户ID数组        |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
      	{
            "uid": 1001,											 //游戏用户ID
            "totalAmount": "1220.1234",        //累积消耗总量
            "consumptionAmount": "10.0000",    //待结算消耗量
            "roleName": "nameABC",             //角色名称
            "roleType": "master",              //角色类型
            "level": "50",                     //等级
            "militaryCapability": "13987",     //战力
            "ranking": "20",                   //战力排行
            "asset": "1000",                   //资产总量
            "extraData": {                     //其他额外数据
                "XXX": 1,                      //字段1
                "XXXX": "extraData1",          //字段2
                ...                            //字段...
            }
        },
      	{
            "uid": 1001,
            "totalAmount": "1220.1234",
            "consumptionAmount": "1220.1234",
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

### **响应**

| Name                 | Type   | Enum | Description      |
| -------------------- | ------ | ---- | ---------------- |
| data                 | array  |      | 用户累积消耗数据 |
| &gt; uid             | number |      | 游戏用户ID       |
| &gt; totalAmount     | string |      | 累积消耗总量     |
| > consumptionAmount  | string |      | 待结算消耗量     |
| > roleName           | string |      | 角色名称         |
| &gt; roleType        | string |      | 角色类型         |
| &gt; level           | string |      | 等级             |
| > militaryCapability | string |      | 战力             |
| > ranking            | string |      | 战力排名         |
| > extraData          | array  |      | 其他信息         |

extraData 字段

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | 字段1       |
| field2 | string |      | 字段2       |
| ...    | string |      |             |

## 结算通知

### **描述**

* 工会申请结算级审核结果通知接口
* 申请结算，统计未结算金额，并持续记录累积消耗
* 审核通过，则完成该结算申请
* 审核拒绝，则应恢复未结算金额

### 路径

* /statistics/settlement

### 方法

* POST

> 请求示例:

```json
{
  "assetName": "ALT",
  "type": "0",
  "uids": [1001,1002]
}
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum**                      | **Description** |
| ------------------- | ------------------ | -------------- | ----------------------------------- | --------------------- |
| assetName           | true               | string         |                                     | 游戏内资产名称        |
| type                | true               | string         | 0(申请结算),1(审核通过),2(审核拒绝) | 通知类型              |
| uids                | true               | array          |                                     | 游戏用户ID数组        |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0
}
```

## 查询资产流向数据

### **描述**

* 查询一定时间内新增资产流向数据，按资产数量倒序排列。

  如当前查询时间为2022-11-14 16:23:30，则1H查询范围是2022-11-14 15:00:00～2022-11-14 16:00:00。

  6H：0点，6点，12点，18点

  12H：0点，12点

  1D：0点

  当前查询时间之前最近的刻度时间点往前一个刻度时间的数据

### 路径

* /statistics/asset/flow?assetName={%s}&scale={%s}&count={%d}

### 方法

* GET

> 请求示例:

```json
{
  "assetName": "ALT",
  "beginTime": "1668065809607",
  "endTime": "1668065975683"
}
```

### **参数**

| **Parameter** | **Required** | **Type** | **Enum** | **Description**     |
| ------------------- | ------------------ | -------------- | -------------- | ------------------------- |
| assetName           | true               | string         |                | 游戏内资产名称            |
| scale               | true               | string         |                | 时间刻度：1H，6H，12H，1D |
| count               | false              | number         |                | 默认100                   |

> 返回json数据如下:

```json
{
    "msg": "success",
    "code": 0,
    "data": [
        {
          "uid": 11115678,                   //游戏用户ID
          "roleName": "nameABC",             //角色名称
          "roleType": "master",              //角色类型
          "level": "50",                     //等级
          "militaryCapability": "13987",     //战力
          "ranking": "20"，                  //战力排行
          "asset": "1000",                   //资产总量
          "extraData": {                     //其他额外数据
              "XXX": 1,                      //字段1
              "XXXX": "extraData1",          //字段2
              ...                            //字段...
          }
          "amount": "1101.1234"              //数量
        },
        {
          "uid": 22225678,                   //游戏用户ID
          "roleName": "nameABC",             //角色名称
          "roleType": "master",              //角色类型
          "level": "50",                     //等级
          "militaryCapability": "13987",     //战力
          "ranking": "20",                   //战力排行
          "asset": "1000",                   //资产总量
          "extraData": {                     //其他额外数据
              "XXX": 1,                      //字段1
              "XXXX": "extraData1",          //字段2
              ...                            //字段...
          }
          "amount": "2202.5678"              //数量
        },
        ...
    ]
}
```

### **响应**

| Name                 | Type   | Enum | Description  |
| -------------------- | ------ | ---- | ------------ |
| data                 | array  |      | 资产流向数据 |
| &gt; uid             | number |      | 游戏用户ID   |
| &gt; amount          | string |      | 资产数量     |
| > roleName           | string |      | 角色名称     |
| &gt; roleType        | string |      | 角色类型     |
| &gt; level           | string |      | 等级         |
| > militaryCapability | string |      | 战力         |
| > ranking            | string |      | 战力排名     |
| > extraData          | array  |      | 其他信息     |

extraData 字段

| Name   | Type   | Enum | Description |
| ------ | ------ | ---- | ----------- |
| field1 | string |      | 字段1       |
| field2 | string |      | 字段2       |
| ...    | string |      |             |

# 后台管理

## 用户状态变更

### 描述

* 用户状态变更

### 路径

* /userStatus

### 方法

* POST

### 请求头Content-Type

* application/x-www-form-urlencoded

> 请求示例:

```curl
curl --location -g --request POST '${gameServerApiUrl}/userStatus' \--header 'Content-Type: application/x-www-form-urlencoded' \--data-urlencode 'appID=1' \--data-urlencode 'uid=1' \--data-urlencode 'status=10' \--data-urlencode 'timestamp=1645415419000' \--data-urlencode 'sign=xxxxxxxxxxxxxxx'
```

### 参数

| **Parameter** | **Required** | **Type** | **Enum**                                                         | **Description**                              |
| ------------------- | ------------------ | -------------- | ---------------------------------------------------------------------- | -------------------------------------------------- |
| appId               | true               | string         |                                                                        | Game ID in Connexion SDK                           |
| uid                 | true               | string         |                                                                        | Game user ID                                       |
| status              | true               | string         | 0-pause withdraw，1-resume withdraw，2-ban，10-force offline，11-unban | changing user status                               |
| timestamp           | true               | string         |                                                                        | Timestamp at the time of communication, unit in ms |
| sign                | true               | string         |                                                                        | Signature                                          |

> 返回数据如下

```
SUCCESS
```

### sign生成方法

1. 除了sign字段本身， 将其余字段按照字段key值的升序排列， 然后按照key=value&key=value.....拼接成字符串
2. 将上述生成的字符串，最后附加&secretKey={SDK Key参数}格式生成待签名字符串。{SDK Key参数}替换为分配的AppSecret参数
3. 对上面得到的待签名字符串，做md5(32位大写)计算，生成对应的sign

> Java 签名示例:

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

### 响应

* 游戏服务器收到请求之后， 校验签名并且状态变更正确， 给SDK Server返回文本字符串：SUCCESS;否则返回文本字符串：FAIL

# 移动客户端

* ConnexionSDK为Mobile Client提供*.framework 和 *.aar等([下载地址](https://github.com/Connector-Gamefi/mobile-sdk)),函数接口说明如下
* 对于Android.XPlatform.getInstance() 返回一个单例对象，可使用该单例对象调用aar包提供的各种方法

## IOS 部署说明

### 部署说明1

* iOS需要删除SceneDelegate,具体为:

> AppDelegate中删除

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

> Info.plist中删除

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

1. 删除AppDelegate中的UISceneSession lifecycle相关代码
2. 删除SceneDelegate.h 和SceneDelegate.m文件
3. 删除Info.plist中的Application Scene Manifest

### 部署说明2

* target -> Build Settings -> Linking ->Other Linker Flags的值设置为 -all_load

## IOS-SDK 初始化

### 描述

* IOS SDK 初始化

### 方法

* initWithAppID

> 示例:

```
UG_SDKParams* sdkParams = [[UG_SDKParams alloc] initWithAppID:kAppID appKey:kAppKey orientation:@"landscape";
```

### 参数

| **Parameter** | **Required** | **Type** | **Enum**         | **Description** |
| ------------------- | ------------------ | -------------- | ---------------------- | --------------------- |
| appID               | true               | string         |                        | 游戏在SDK中的标识ID   |
| appKey              | true               | string         |                        | appKey                |
| orientation         | true               | string         | "portrait","landscape" | 横竖屏                |

## IOS-SDK 登陆

### 描述

* IOS 登陆

### 方法

* sharedInstance.login

> 示例:

```
[[UGSDKPlatform sharedInstance] login];
```

## IOS-SDK 登出

### 描述

* IOS 登出

### 方法

* sharedInstance.logout

> 示例:

```
[[UGSDKPlatform sharedInstance] logout];
```

## IOS-SDK 数据上报

### 描述

* IOS 数据上报

### 方法

* sharedInstance.submitGameData

> 示例:

```
[[UGSDKPlatform sharedInstance] submitGameData:role];
```

### 参数

| **Parameter** | **Required** | **Type** | **Enum**                                             | **Description** |
| ------------------- | ------------------ | -------------- | ---------------------------------------------------------- | --------------------- |
| role                | true               | *UG_GameRole   |                                                            |                       |
| &gt;type            | true               | string         | 1-创建角色,2-进入游戏,3-等级升级,4-退出游戏,others其他类型 | 上报类型              |
| &gt;roleID          | true               | string         |                                                            | 角色ID                |
| &gt;roleName        | true               | string         |                                                            | 角色名称              |
| &gt;roleLevel       | true               | string         |                                                            | 角色等级              |
| &gt;serverID        | true               | string         |                                                            | 服务器ID              |
| &gt;serverName      | true               | string         |                                                            | 服务器名称            |
| &gt;vip             | true               | string         |                                                            | VIP等级               |
| &gt;createTime      | true               | string         |                                                            | 时间戳,Unix秒         |
| &gt;lastLevelUpTime | true               | string         |                                                            | 最后升级时间,Unix秒   |
| &gt;extraData       | true               | string         |                                                            | 其他扩展数据          |

## Android-SDK 初始化

### 描述

* 初始化,需要在游戏启动的时候调用，一般在游戏Activity的onCreate中调用

### 方法

* init

> 示例:

```java
public void init(Activity context, UInitParams params, IInitListener listener)
```

## Android-SDK SetLogoutListener

### 描述

* 设置从悬浮窗中点击切换账号之后的回调，游戏层收到该回调，需要返回到游戏登陆界面，引导玩家重新登陆

### 方法

* setLogoutListener

> 示例:

```java
public void setLogoutListener(ILogoutListener logoutListener)
```

## Android-SDK 登陆

### 描述

* 登陆接口,调用之后,弹出SDK的登陆界面

### 方法

* login

> 示例:

```java
public void login(Activity activity, ILoginListener listener)
```

## Android-SDK 登出

### 描述

* 登出接口,调用之后,清除当前登录状态

### 方法

* logout

> 示例:

```java
public void logout(Activity activity)
```

## Android-SDK 数据上报

### 描述

* 数据上报,需要在创建角色、进入游戏、等级升级、退出游戏等几个地方调用

### 方法

* submit

> 示例:

```java
public void submit(Activity activity, URoleData roleData)
```

### 参数

| **Parameter** | **Required** | **Type** | **Enum**                                             | **Description**      |
| ------------------- | ------------------ | -------------- | ---------------------------------------------------------- | -------------------------- |
| roleData            | true               | URoleData      |                                                            |                            |
| &gt;type            | true               | string         | 1-创建角色,2-进入游戏,3-等级升级,4-退出游戏,others其他类型 | 上报类型                   |
| &gt;roleID          | true               | string         |                                                            | 角色ID                     |
| &gt;roleName        | true               | string         |                                                            | 角色名称                   |
| &gt;roleLevel       | true               | string         |                                                            | 角色等级                   |
| &gt;serverID        | true               | string         |                                                            | 服务器ID                   |
| &gt;serverName      | true               | string         |                                                            | 服务器名称                 |
| &gt;vip             | true               | string         |                                                            | VIP等级                    |
| &gt;createTime      | true               | string         |                                                            | 时间戳,Unix秒              |
| &gt;lastLevelUpTime | true               | string         |                                                            | 最后升级时间,Unix秒        |
| &gt;extraData       | true               | string         |                                                            | json格式,扩展数据,暂时无用 |

## Android-SDK 退出游戏

### 描述

* 退出游戏，弹出退出确认框

### 方法

* exit

> 示例:

```java
public void exit(final Activity context)
```

## Android-SDK 生命周期方法

* Android-SDK 生命周期方法 onStart onStop   onResume   onPause   onDestroy
