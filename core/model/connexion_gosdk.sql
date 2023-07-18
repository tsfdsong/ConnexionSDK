/*
 Navicat Premium Data Transfer

 Source Server         : sdk-test
 Source Server Type    : MySQL
 Source Server Version : 50726
 Source Host           : localhost:3306
 Source Schema         : connexion_gosdk

 Target Server Type    : MySQL
 Target Server Version : 50726
 File Encoding         : 65001

 Date: 27/06/2022 15:32:16
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_attribute
-- ----------------------------
DROP TABLE IF EXISTS `t_attribute`;
CREATE TABLE `t_attribute` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` tinyint(3) unsigned NOT NULL COMMENT '''游戏ID''',
  `attr_id` bigint(20) unsigned NOT NULL COMMENT '''属性ID''',
  `attr_description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '''属性描述''',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'contract address',
  `attr_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'attr name',
  `attr_decimal` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'attr_decimal',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `appid_contract_attrid` (`app_id`,`contract_address`,`attr_id`) USING BTREE,
  KEY `idx_t_attributes_game_id` (`attr_id`,`contract_address`,`attr_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_blindbox_config
-- ----------------------------
DROP TABLE IF EXISTS `t_blindbox_config`;
CREATE TABLE `t_blindbox_config` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `network` varchar(32) NOT NULL DEFAULT '',
  `seller` varchar(42) NOT NULL DEFAULT '' COMMENT '销售(盲盒平台)地址',
  `eq_address` varchar(42) NOT NULL DEFAULT '' COMMENT 'mint(装备)合约地址',
  `merkle_root` varchar(255) NOT NULL DEFAULT '' COMMENT 'merkle_root',
  `file_name` varchar(20) NOT NULL DEFAULT '' COMMENT '属性文件名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_seller` (`seller`) USING BTREE,
  KEY `idx_network` (`network`) USING BTREE,
  KEY `idx_mint` (`eq_address`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for t_blindbox_eq_attr_data
-- ----------------------------
DROP TABLE IF EXISTS `t_blindbox_eq_attr_data`;
CREATE TABLE `t_blindbox_eq_attr_data` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `image_id` varchar(64) NOT NULL DEFAULT '' COMMENT '图片id',
  `seller` varchar(42) NOT NULL DEFAULT '' COMMENT '销售(盲盒平台)地址',
  `equipment_id` varchar(42) NOT NULL DEFAULT '' COMMENT '装备id',
  `bnft_token_id` varchar(255) NOT NULL DEFAULT '' COMMENT '销售平台nft tokenid',
  `eq_token_id` varchar(255) NOT NULL DEFAULT '' COMMENT '装备nft tokenID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `equipment_attr` json NOT NULL COMMENT '属性列表',
  `proofs` json NOT NULL COMMENT 'proofs',
  PRIMARY KEY (`id`),
  KEY `idx_bnft_tokenid` (`bnft_token_id`) USING BTREE,
  KEY `idx_seller` (`seller`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=141 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for t_block_height
-- ----------------------------
DROP TABLE IF EXISTS `t_block_height`;
CREATE TABLE `t_block_height` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `latest_parsed_height` bigint(20) unsigned NOT NULL COMMENT '''最新解析过的高度''',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `switch` tinyint(3) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_t_block_heights_latest_parsed_height` (`latest_parsed_height`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_email_bind
-- ----------------------------
DROP TABLE IF EXISTS `t_email_bind`;
CREATE TABLE `t_email_bind` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` int(10) unsigned NOT NULL COMMENT '''游戏ID''',
  `uid` bigint(20) unsigned NOT NULL COMMENT '''用户ID''',
  `address` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''用户地址''',
  `account` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''用户邮箱''',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_t_email_binds_uid` (`uid`) USING BTREE,
  KEY `game_addr` (`app_id`,`address`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=123 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_ft_contract
-- ----------------------------
DROP TABLE IF EXISTS `t_ft_contract`;
CREATE TABLE `t_ft_contract` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) DEFAULT NULL COMMENT '游戏app id',
  `chain_id` int(11) DEFAULT NULL COMMENT '区块链网络',
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '同质化合约地址',
  `token_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币全称',
  `token_symbol` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币简称',
  `token_supply` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币总量',
  `token_decimal` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '代币精度',
  `game_coin_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '游戏中代币名称',
  `deposit_switch` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '充值开关',
  `withdraw_switch` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '提现开关',
  `treasure` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '金库合约',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `game_decimal` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_ft_deposit_record
-- ----------------------------
DROP TABLE IF EXISTS `t_ft_deposit_record`;
CREATE TABLE `t_ft_deposit_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0' COMMENT '游戏ID',
  `app_order_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'SDK充值订单编号',
  `account` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '充值用户账号',
  `amount` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '充值金额',
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'FT充值代币合约地址',
  `deposit_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '充值地址',
  `target_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '充值代币目标地址',
  `nonce` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'sdk nonce',
  `height` bigint(10) DEFAULT '0',
  `tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '交易Hash',
  `order_status` int(11) DEFAULT '0' COMMENT '订单状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `uid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '玩家UID',
  `game_coin_name` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=527 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_ft_withdraw_record
-- ----------------------------
DROP TABLE IF EXISTS `t_ft_withdraw_record`;
CREATE TABLE `t_ft_withdraw_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0' COMMENT '游戏ID',
  `app_order_id` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '游戏内部订单编号',
  `account` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '提现用户账号ID',
  `height` bigint(10) DEFAULT '0',
  `tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '交易Hash',
  `order_status` int(11) DEFAULT NULL COMMENT '提现订单状态',
  `amount` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '提现金额',
  `uid` bigint(20) unsigned NOT NULL DEFAULT '0',
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '提现FT的代币合约地址',
  `withdraw_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '提现地址',
  `signature` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '签名数据',
  `signature_hash` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '签名机签名hash',
  `nonce` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'sdk nonce',
  `risk_status` int(11) NOT NULL DEFAULT '0' COMMENT '风控审核状态',
  `risk_reviewer` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_coin_name` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `key` (`order_status`) USING BTREE,
  KEY `unq_order_id` (`app_order_id`,`app_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=376 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_game_equipment
-- ----------------------------
DROP TABLE IF EXISTS `t_game_equipment`;
CREATE TABLE `t_game_equipment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0' COMMENT '游戏id',
  `account` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `equipment_id` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token_id` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '装备对应链上NFT的TokenID',
  `status` int(4) unsigned NOT NULL DEFAULT '0' COMMENT '装备状态0初始化1已充值2已提现',
  `equipment_attr` json DEFAULT NULL,
  `image_uri` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '装备NFT合约地址',
  `game_asset_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chain_id` int(11) NOT NULL,
  `withdraw_switch` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uq_eid` (`app_id`,`equipment_id`,`token_id`,`contract_address`) USING BTREE,
  KEY `index` (`equipment_id`,`app_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_game_info
-- ----------------------------
DROP TABLE IF EXISTS `t_game_info`;
CREATE TABLE `t_game_info` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` tinyint(3) unsigned NOT NULL COMMENT '''游戏ID''',
  `base_server_url` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''游戏serverAPI的ip端口地址''',
  `app_secret` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '''erc721金库合约地址''',
  `app_key` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unq_app_id` (`app_id`) USING BTREE,
  KEY `idx_t_game_infos_game_id` (`app_id`,`base_server_url`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_nft_contract
-- ----------------------------
DROP TABLE IF EXISTS `t_nft_contract`;
CREATE TABLE `t_nft_contract` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) DEFAULT NULL COMMENT '游戏app id',
  `chain_id` int(11) DEFAULT NULL COMMENT '区块链网络chained',
  `treasure` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `minter_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Game minter contract address',
  `game_asset_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '游戏内资产名称',
  `contract_address` char(42) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '非同质化资产合约地址',
  `token_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币名称',
  `token_symbol` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币简称',
  `token_supply` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币总量',
  `base_url` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'nft图片基础URL',
  `decimal` int(11) NOT NULL DEFAULT '0' COMMENT '属性精度',
  `deposit_switch` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '充值开关',
  `withdraw_switch` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '提现开关',
  `file_name` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `attr_update_time` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unq` (`app_id`,`contract_address`,`game_asset_name`) USING BTREE,
  KEY `contract_address` (`contract_address`(21),`app_id`,`game_asset_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_nft_deposit_record
-- ----------------------------
DROP TABLE IF EXISTS `t_nft_deposit_record`;
CREATE TABLE `t_nft_deposit_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `equipment_id` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `token_id` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '充值代币ID',
  `order_status` int(11) DEFAULT '0' COMMENT '订单状态',
  `tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '交易Hash',
  `height` bigint(10) DEFAULT '0',
  `app_order_id` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'SDK充值订单编号',
  `account` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '充值用户账户',
  `uid` bigint(20) DEFAULT NULL,
  `app_id` int(11) DEFAULT NULL,
  `nonce` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `deposit_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '充值地址',
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'NFT 充值代币合约地址',
  `trease_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'NFT充值金库合约地址',
  `game_asset_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `index` (`token_id`,`order_status`,`deposit_address`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_nft_withdraw_record
-- ----------------------------
DROP TABLE IF EXISTS `t_nft_withdraw_record`;
CREATE TABLE `t_nft_withdraw_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `equipment_id` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '装备id',
  `token_id` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '装备对应链上NFT的TokenID',
  `order_status` int(11) DEFAULT NULL COMMENT '提现订单状态',
  `app_order_id` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '游戏内部订单编号',
  `height` bigint(10) DEFAULT '0',
  `tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '交易Hash',
  `account` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '提现用户账号',
  `uid` bigint(20) DEFAULT NULL,
  `nonce` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `app_id` int(11) NOT NULL DEFAULT '0' COMMENT '游戏ID',
  `withdraw_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '提现地址',
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '提现NFT的代币合约地址',
  `trease_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '金库合约地址',
  `minter_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Game minter contract address',
  `signature` varchar(150) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '签名数据',
  `signature_hash` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '签名机签名hash',
  `signature_source` json DEFAULT NULL COMMENT '装备对应属性列表',
  `risk_status` int(11) DEFAULT '0',
  `risk_reviewer` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `game_asset_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `ind_key` (`app_id`,`order_status`,`withdraw_address`,`uid`,`app_order_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_risk_control
-- ----------------------------
DROP TABLE IF EXISTS `t_risk_control`;
CREATE TABLE `t_risk_control` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) DEFAULT NULL COMMENT '游戏App ID',
  `token_type` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币类型， erc20: FT； erc721: NFT',
  `contract_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代币合约地址',
  `token_name` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `token_symbol` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount_limit` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '单笔提现额度',
  `count_limit` int(11) DEFAULT '0' COMMENT '连续高频提现次数限制',
  `count_time` int(11) DEFAULT '0' COMMENT '高频提现风控时间限制',
  `total_limit` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '连续提现额度限制',
  `total_time` int(11) DEFAULT '0' COMMENT '连续提现风控时间限制',
  `status` tinyint(4) DEFAULT '0' COMMENT '是否启用风控， 0 ： 不启用； 1: 启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uq_index` (`app_id`,`token_type`,`contract_address`,`token_symbol`,`token_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` tinyint(3) unsigned NOT NULL COMMENT '''游戏ID''',
  `account` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''邮箱''',
  `name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''用户名''',
  `image` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''用户图片''',
  `bio` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '''用户个人简介''',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `withdraw_switch` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '提现开关',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `game_email` (`app_id`,`account`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
