// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AttributeData is an auto generated low-level Go binding around an user-defined struct.
type AttributeData struct {
	AttrID    *big.Int
	AttrValue *big.Int
}

// GameLootEquipmentMetaData contains all meta data concerning the GameLootEquipment contract.
var GameLootEquipmentMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"revealSVG_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"treasure_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"vault_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"timeLocker_\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_signers\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"cap_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"attrID\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"name\":\"AttributeAttached\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint128[]\",\"name\":\"attrIDs\",\"type\":\"uint128[]\"},{\"indexed\":false,\"internalType\":\"uint128[]\",\"name\":\"values\",\"type\":\"uint128[]\"}],\"name\":\"AttributeAttachedBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint128[]\",\"name\":\"attrIDs\",\"type\":\"uint128[]\"}],\"name\":\"AttributeRemoveBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"attrID\",\"type\":\"uint128\"}],\"name\":\"AttributeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"attrIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"}],\"name\":\"AttributeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"attrIndexes\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint128[]\",\"name\":\"values\",\"type\":\"uint128[]\"}],\"name\":\"AttributeUpdatedBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"Revealed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"attrID_\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"_value\",\"type\":\"uint128\"}],\"name\":\"attach\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint128[]\",\"name\":\"attrIDs_\",\"type\":\"uint128[]\"},{\"internalType\":\"uint128[]\",\"name\":\"_values\",\"type\":\"uint128[]\"}],\"name\":\"attachBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenID\",\"type\":\"uint256\"}],\"name\":\"attributes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"attrID\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"attrValue\",\"type\":\"uint128\"}],\"internalType\":\"structAttributeData[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gameMinter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hasRevealed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"reciever\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"attrIndex_\",\"type\":\"uint256\"}],\"name\":\"remove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"attrIndexes_\",\"type\":\"uint256[]\"}],\"name\":\"removeBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce_\",\"type\":\"uint256\"},{\"internalType\":\"uint128[]\",\"name\":\"attrIDs_\",\"type\":\"uint128[]\"},{\"internalType\":\"uint128[]\",\"name\":\"attrValues_\",\"type\":\"uint128[]\"},{\"internalType\":\"bytes\",\"name\":\"signature_\",\"type\":\"bytes\"}],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"seller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"blindBoxURI_\",\"type\":\"string\"}],\"name\":\"setBlindBoxURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"cap\",\"type\":\"uint256\"}],\"name\":\"setCap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"gameMinter_\",\"type\":\"address\"}],\"name\":\"setGameMinter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"revealSVG_\",\"type\":\"address\"}],\"name\":\"setReveal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"setRevealed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller_\",\"type\":\"address\"}],\"name\":\"setSeller\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isOk\",\"type\":\"bool\"}],\"name\":\"setSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"timeLocker_\",\"type\":\"address\"}],\"name\":\"setTimeLocker\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"treasure_\",\"type\":\"address\"}],\"name\":\"setTreasure\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault_\",\"type\":\"address\"}],\"name\":\"setVault\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeLocker\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"registry\",\"type\":\"address\"}],\"name\":\"tokenURIGame\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasure\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"attrIndex_\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"value_\",\"type\":\"uint128\"}],\"name\":\"update\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"attrIndexes_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint128[]\",\"name\":\"values_\",\"type\":\"uint128[]\"}],\"name\":\"updateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"usedNonce\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vault\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// GameLootEquipmentABI is the input ABI used to generate the binding from.
// Deprecated: Use GameLootEquipmentMetaData.ABI instead.
var GameLootEquipmentABI = GameLootEquipmentMetaData.ABI

// GameLootEquipment is an auto generated Go binding around an Ethereum contract.
type GameLootEquipment struct {
	GameLootEquipmentCaller     // Read-only binding to the contract
	GameLootEquipmentTransactor // Write-only binding to the contract
	GameLootEquipmentFilterer   // Log filterer for contract events
}

// GameLootEquipmentCaller is an auto generated read-only Go binding around an Ethereum contract.
type GameLootEquipmentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GameLootEquipmentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GameLootEquipmentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GameLootEquipmentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GameLootEquipmentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GameLootEquipmentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GameLootEquipmentSession struct {
	Contract     *GameLootEquipment // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// GameLootEquipmentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GameLootEquipmentCallerSession struct {
	Contract *GameLootEquipmentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// GameLootEquipmentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GameLootEquipmentTransactorSession struct {
	Contract     *GameLootEquipmentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GameLootEquipmentRaw is an auto generated low-level Go binding around an Ethereum contract.
type GameLootEquipmentRaw struct {
	Contract *GameLootEquipment // Generic contract binding to access the raw methods on
}

// GameLootEquipmentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GameLootEquipmentCallerRaw struct {
	Contract *GameLootEquipmentCaller // Generic read-only contract binding to access the raw methods on
}

// GameLootEquipmentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GameLootEquipmentTransactorRaw struct {
	Contract *GameLootEquipmentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGameLootEquipment creates a new instance of GameLootEquipment, bound to a specific deployed contract.
func NewGameLootEquipment(address common.Address, backend bind.ContractBackend) (*GameLootEquipment, error) {
	contract, err := bindGameLootEquipment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipment{GameLootEquipmentCaller: GameLootEquipmentCaller{contract: contract}, GameLootEquipmentTransactor: GameLootEquipmentTransactor{contract: contract}, GameLootEquipmentFilterer: GameLootEquipmentFilterer{contract: contract}}, nil
}

// NewGameLootEquipmentCaller creates a new read-only instance of GameLootEquipment, bound to a specific deployed contract.
func NewGameLootEquipmentCaller(address common.Address, caller bind.ContractCaller) (*GameLootEquipmentCaller, error) {
	contract, err := bindGameLootEquipment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentCaller{contract: contract}, nil
}

// NewGameLootEquipmentTransactor creates a new write-only instance of GameLootEquipment, bound to a specific deployed contract.
func NewGameLootEquipmentTransactor(address common.Address, transactor bind.ContractTransactor) (*GameLootEquipmentTransactor, error) {
	contract, err := bindGameLootEquipment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentTransactor{contract: contract}, nil
}

// NewGameLootEquipmentFilterer creates a new log filterer instance of GameLootEquipment, bound to a specific deployed contract.
func NewGameLootEquipmentFilterer(address common.Address, filterer bind.ContractFilterer) (*GameLootEquipmentFilterer, error) {
	contract, err := bindGameLootEquipment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentFilterer{contract: contract}, nil
}

// bindGameLootEquipment binds a generic wrapper to an already deployed contract.
func bindGameLootEquipment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GameLootEquipmentABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GameLootEquipment *GameLootEquipmentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GameLootEquipment.Contract.GameLootEquipmentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GameLootEquipment *GameLootEquipmentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.GameLootEquipmentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GameLootEquipment *GameLootEquipmentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.GameLootEquipmentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GameLootEquipment *GameLootEquipmentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GameLootEquipment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GameLootEquipment *GameLootEquipmentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GameLootEquipment *GameLootEquipmentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.contract.Transact(opts, method, params...)
}

// Attributes is a free data retrieval call binding the contract method 0xd05dcc6a.
//
// Solidity: function attributes(uint256 _tokenID) view returns((uint128,uint128)[])
func (_GameLootEquipment *GameLootEquipmentCaller) Attributes(opts *bind.CallOpts, _tokenID *big.Int) ([]AttributeData, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "attributes", _tokenID)

	if err != nil {
		return *new([]AttributeData), err
	}

	out0 := *abi.ConvertType(out[0], new([]AttributeData)).(*[]AttributeData)

	return out0, err

}

// Attributes is a free data retrieval call binding the contract method 0xd05dcc6a.
//
// Solidity: function attributes(uint256 _tokenID) view returns((uint128,uint128)[])
func (_GameLootEquipment *GameLootEquipmentSession) Attributes(_tokenID *big.Int) ([]AttributeData, error) {
	return _GameLootEquipment.Contract.Attributes(&_GameLootEquipment.CallOpts, _tokenID)
}

// Attributes is a free data retrieval call binding the contract method 0xd05dcc6a.
//
// Solidity: function attributes(uint256 _tokenID) view returns((uint128,uint128)[])
func (_GameLootEquipment *GameLootEquipmentCallerSession) Attributes(_tokenID *big.Int) ([]AttributeData, error) {
	return _GameLootEquipment.Contract.Attributes(&_GameLootEquipment.CallOpts, _tokenID)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GameLootEquipment.Contract.BalanceOf(&_GameLootEquipment.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GameLootEquipment.Contract.BalanceOf(&_GameLootEquipment.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_GameLootEquipment *GameLootEquipmentCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_GameLootEquipment *GameLootEquipmentSession) Decimals() (uint8, error) {
	return _GameLootEquipment.Contract.Decimals(&_GameLootEquipment.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Decimals() (uint8, error) {
	return _GameLootEquipment.Contract.Decimals(&_GameLootEquipment.CallOpts)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns()
func (_GameLootEquipment *GameLootEquipmentCaller) Exists(opts *bind.CallOpts, tokenId *big.Int) error {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "exists", tokenId)

	if err != nil {
		return err
	}

	return err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns()
func (_GameLootEquipment *GameLootEquipmentSession) Exists(tokenId *big.Int) error {
	return _GameLootEquipment.Contract.Exists(&_GameLootEquipment.CallOpts, tokenId)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns()
func (_GameLootEquipment *GameLootEquipmentCallerSession) Exists(tokenId *big.Int) error {
	return _GameLootEquipment.Contract.Exists(&_GameLootEquipment.CallOpts, tokenId)
}

// GameMinter is a free data retrieval call binding the contract method 0x39654844.
//
// Solidity: function gameMinter() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) GameMinter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "gameMinter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GameMinter is a free data retrieval call binding the contract method 0x39654844.
//
// Solidity: function gameMinter() view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) GameMinter() (common.Address, error) {
	return _GameLootEquipment.Contract.GameMinter(&_GameLootEquipment.CallOpts)
}

// GameMinter is a free data retrieval call binding the contract method 0x39654844.
//
// Solidity: function gameMinter() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) GameMinter() (common.Address, error) {
	return _GameLootEquipment.Contract.GameMinter(&_GameLootEquipment.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GameLootEquipment.Contract.GetApproved(&_GameLootEquipment.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GameLootEquipment.Contract.GetApproved(&_GameLootEquipment.CallOpts, tokenId)
}

// GetCap is a free data retrieval call binding the contract method 0x554d578d.
//
// Solidity: function getCap() view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentCaller) GetCap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "getCap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCap is a free data retrieval call binding the contract method 0x554d578d.
//
// Solidity: function getCap() view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentSession) GetCap() (*big.Int, error) {
	return _GameLootEquipment.Contract.GetCap(&_GameLootEquipment.CallOpts)
}

// GetCap is a free data retrieval call binding the contract method 0x554d578d.
//
// Solidity: function getCap() view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentCallerSession) GetCap() (*big.Int, error) {
	return _GameLootEquipment.Contract.GetCap(&_GameLootEquipment.CallOpts)
}

// HasRevealed is a free data retrieval call binding the contract method 0x7f8052be.
//
// Solidity: function hasRevealed(uint256 ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCaller) HasRevealed(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "hasRevealed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRevealed is a free data retrieval call binding the contract method 0x7f8052be.
//
// Solidity: function hasRevealed(uint256 ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentSession) HasRevealed(arg0 *big.Int) (bool, error) {
	return _GameLootEquipment.Contract.HasRevealed(&_GameLootEquipment.CallOpts, arg0)
}

// HasRevealed is a free data retrieval call binding the contract method 0x7f8052be.
//
// Solidity: function hasRevealed(uint256 ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCallerSession) HasRevealed(arg0 *big.Int) (bool, error) {
	return _GameLootEquipment.Contract.HasRevealed(&_GameLootEquipment.CallOpts, arg0)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GameLootEquipment.Contract.IsApprovedForAll(&_GameLootEquipment.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GameLootEquipment.Contract.IsApprovedForAll(&_GameLootEquipment.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GameLootEquipment *GameLootEquipmentCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GameLootEquipment *GameLootEquipmentSession) Name() (string, error) {
	return _GameLootEquipment.Contract.Name(&_GameLootEquipment.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Name() (string, error) {
	return _GameLootEquipment.Contract.Name(&_GameLootEquipment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) Owner() (common.Address, error) {
	return _GameLootEquipment.Contract.Owner(&_GameLootEquipment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Owner() (common.Address, error) {
	return _GameLootEquipment.Contract.Owner(&_GameLootEquipment.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GameLootEquipment.Contract.OwnerOf(&_GameLootEquipment.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GameLootEquipment.Contract.OwnerOf(&_GameLootEquipment.CallOpts, tokenId)
}

// Seller is a free data retrieval call binding the contract method 0x08551a53.
//
// Solidity: function seller() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) Seller(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "seller")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Seller is a free data retrieval call binding the contract method 0x08551a53.
//
// Solidity: function seller() view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) Seller() (common.Address, error) {
	return _GameLootEquipment.Contract.Seller(&_GameLootEquipment.CallOpts)
}

// Seller is a free data retrieval call binding the contract method 0x08551a53.
//
// Solidity: function seller() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Seller() (common.Address, error) {
	return _GameLootEquipment.Contract.Seller(&_GameLootEquipment.CallOpts)
}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCaller) Signers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "signers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentSession) Signers(arg0 common.Address) (bool, error) {
	return _GameLootEquipment.Contract.Signers(&_GameLootEquipment.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Signers(arg0 common.Address) (bool, error) {
	return _GameLootEquipment.Contract.Signers(&_GameLootEquipment.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GameLootEquipment.Contract.SupportsInterface(&_GameLootEquipment.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GameLootEquipment.Contract.SupportsInterface(&_GameLootEquipment.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GameLootEquipment *GameLootEquipmentCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GameLootEquipment *GameLootEquipmentSession) Symbol() (string, error) {
	return _GameLootEquipment.Contract.Symbol(&_GameLootEquipment.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Symbol() (string, error) {
	return _GameLootEquipment.Contract.Symbol(&_GameLootEquipment.CallOpts)
}

// TimeLocker is a free data retrieval call binding the contract method 0xd792fd5c.
//
// Solidity: function timeLocker() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) TimeLocker(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "timeLocker")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TimeLocker is a free data retrieval call binding the contract method 0xd792fd5c.
//
// Solidity: function timeLocker() view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) TimeLocker() (common.Address, error) {
	return _GameLootEquipment.Contract.TimeLocker(&_GameLootEquipment.CallOpts)
}

// TimeLocker is a free data retrieval call binding the contract method 0xd792fd5c.
//
// Solidity: function timeLocker() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) TimeLocker() (common.Address, error) {
	return _GameLootEquipment.Contract.TimeLocker(&_GameLootEquipment.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenID) view returns(string)
func (_GameLootEquipment *GameLootEquipmentCaller) TokenURI(opts *bind.CallOpts, tokenID *big.Int) (string, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "tokenURI", tokenID)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenID) view returns(string)
func (_GameLootEquipment *GameLootEquipmentSession) TokenURI(tokenID *big.Int) (string, error) {
	return _GameLootEquipment.Contract.TokenURI(&_GameLootEquipment.CallOpts, tokenID)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenID) view returns(string)
func (_GameLootEquipment *GameLootEquipmentCallerSession) TokenURI(tokenID *big.Int) (string, error) {
	return _GameLootEquipment.Contract.TokenURI(&_GameLootEquipment.CallOpts, tokenID)
}

// TokenURIGame is a free data retrieval call binding the contract method 0xef1750c0.
//
// Solidity: function tokenURIGame(uint256 tokenID, address registry) view returns(string)
func (_GameLootEquipment *GameLootEquipmentCaller) TokenURIGame(opts *bind.CallOpts, tokenID *big.Int, registry common.Address) (string, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "tokenURIGame", tokenID, registry)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURIGame is a free data retrieval call binding the contract method 0xef1750c0.
//
// Solidity: function tokenURIGame(uint256 tokenID, address registry) view returns(string)
func (_GameLootEquipment *GameLootEquipmentSession) TokenURIGame(tokenID *big.Int, registry common.Address) (string, error) {
	return _GameLootEquipment.Contract.TokenURIGame(&_GameLootEquipment.CallOpts, tokenID, registry)
}

// TokenURIGame is a free data retrieval call binding the contract method 0xef1750c0.
//
// Solidity: function tokenURIGame(uint256 tokenID, address registry) view returns(string)
func (_GameLootEquipment *GameLootEquipmentCallerSession) TokenURIGame(tokenID *big.Int, registry common.Address) (string, error) {
	return _GameLootEquipment.Contract.TokenURIGame(&_GameLootEquipment.CallOpts, tokenID, registry)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentSession) TotalSupply() (*big.Int, error) {
	return _GameLootEquipment.Contract.TotalSupply(&_GameLootEquipment.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GameLootEquipment *GameLootEquipmentCallerSession) TotalSupply() (*big.Int, error) {
	return _GameLootEquipment.Contract.TotalSupply(&_GameLootEquipment.CallOpts)
}

// Treasure is a free data retrieval call binding the contract method 0xe520fc7e.
//
// Solidity: function treasure() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) Treasure(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "treasure")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasure is a free data retrieval call binding the contract method 0xe520fc7e.
//
// Solidity: function treasure() view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) Treasure() (common.Address, error) {
	return _GameLootEquipment.Contract.Treasure(&_GameLootEquipment.CallOpts)
}

// Treasure is a free data retrieval call binding the contract method 0xe520fc7e.
//
// Solidity: function treasure() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Treasure() (common.Address, error) {
	return _GameLootEquipment.Contract.Treasure(&_GameLootEquipment.CallOpts)
}

// UsedNonce is a free data retrieval call binding the contract method 0x9723fb6d.
//
// Solidity: function usedNonce(uint256 ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCaller) UsedNonce(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "usedNonce", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedNonce is a free data retrieval call binding the contract method 0x9723fb6d.
//
// Solidity: function usedNonce(uint256 ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentSession) UsedNonce(arg0 *big.Int) (bool, error) {
	return _GameLootEquipment.Contract.UsedNonce(&_GameLootEquipment.CallOpts, arg0)
}

// UsedNonce is a free data retrieval call binding the contract method 0x9723fb6d.
//
// Solidity: function usedNonce(uint256 ) view returns(bool)
func (_GameLootEquipment *GameLootEquipmentCallerSession) UsedNonce(arg0 *big.Int) (bool, error) {
	return _GameLootEquipment.Contract.UsedNonce(&_GameLootEquipment.CallOpts, arg0)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCaller) Vault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GameLootEquipment.contract.Call(opts, &out, "vault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_GameLootEquipment *GameLootEquipmentSession) Vault() (common.Address, error) {
	return _GameLootEquipment.Contract.Vault(&_GameLootEquipment.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_GameLootEquipment *GameLootEquipmentCallerSession) Vault() (common.Address, error) {
	return _GameLootEquipment.Contract.Vault(&_GameLootEquipment.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Approve(&_GameLootEquipment.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Approve(&_GameLootEquipment.TransactOpts, to, tokenId)
}

// Attach is a paid mutator transaction binding the contract method 0x41be669d.
//
// Solidity: function attach(uint256 tokenID_, uint128 attrID_, uint128 _value) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Attach(opts *bind.TransactOpts, tokenID_ *big.Int, attrID_ *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "attach", tokenID_, attrID_, _value)
}

// Attach is a paid mutator transaction binding the contract method 0x41be669d.
//
// Solidity: function attach(uint256 tokenID_, uint128 attrID_, uint128 _value) returns()
func (_GameLootEquipment *GameLootEquipmentSession) Attach(tokenID_ *big.Int, attrID_ *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Attach(&_GameLootEquipment.TransactOpts, tokenID_, attrID_, _value)
}

// Attach is a paid mutator transaction binding the contract method 0x41be669d.
//
// Solidity: function attach(uint256 tokenID_, uint128 attrID_, uint128 _value) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Attach(tokenID_ *big.Int, attrID_ *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Attach(&_GameLootEquipment.TransactOpts, tokenID_, attrID_, _value)
}

// AttachBatch is a paid mutator transaction binding the contract method 0x6d2468d6.
//
// Solidity: function attachBatch(uint256 tokenID_, uint128[] attrIDs_, uint128[] _values) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) AttachBatch(opts *bind.TransactOpts, tokenID_ *big.Int, attrIDs_ []*big.Int, _values []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "attachBatch", tokenID_, attrIDs_, _values)
}

// AttachBatch is a paid mutator transaction binding the contract method 0x6d2468d6.
//
// Solidity: function attachBatch(uint256 tokenID_, uint128[] attrIDs_, uint128[] _values) returns()
func (_GameLootEquipment *GameLootEquipmentSession) AttachBatch(tokenID_ *big.Int, attrIDs_ []*big.Int, _values []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.AttachBatch(&_GameLootEquipment.TransactOpts, tokenID_, attrIDs_, _values)
}

// AttachBatch is a paid mutator transaction binding the contract method 0x6d2468d6.
//
// Solidity: function attachBatch(uint256 tokenID_, uint128[] attrIDs_, uint128[] _values) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) AttachBatch(tokenID_ *big.Int, attrIDs_ []*big.Int, _values []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.AttachBatch(&_GameLootEquipment.TransactOpts, tokenID_, attrIDs_, _values)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address reciever, uint256 amount) payable returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Mint(opts *bind.TransactOpts, reciever common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "mint", reciever, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address reciever, uint256 amount) payable returns()
func (_GameLootEquipment *GameLootEquipmentSession) Mint(reciever common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Mint(&_GameLootEquipment.TransactOpts, reciever, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address reciever, uint256 amount) payable returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Mint(reciever common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Mint(&_GameLootEquipment.TransactOpts, reciever, amount)
}

// Remove is a paid mutator transaction binding the contract method 0x6526db7a.
//
// Solidity: function remove(uint256 tokenID_, uint256 attrIndex_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Remove(opts *bind.TransactOpts, tokenID_ *big.Int, attrIndex_ *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "remove", tokenID_, attrIndex_)
}

// Remove is a paid mutator transaction binding the contract method 0x6526db7a.
//
// Solidity: function remove(uint256 tokenID_, uint256 attrIndex_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) Remove(tokenID_ *big.Int, attrIndex_ *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Remove(&_GameLootEquipment.TransactOpts, tokenID_, attrIndex_)
}

// Remove is a paid mutator transaction binding the contract method 0x6526db7a.
//
// Solidity: function remove(uint256 tokenID_, uint256 attrIndex_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Remove(tokenID_ *big.Int, attrIndex_ *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Remove(&_GameLootEquipment.TransactOpts, tokenID_, attrIndex_)
}

// RemoveBatch is a paid mutator transaction binding the contract method 0x350ce9e8.
//
// Solidity: function removeBatch(uint256 tokenID_, uint256[] attrIndexes_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) RemoveBatch(opts *bind.TransactOpts, tokenID_ *big.Int, attrIndexes_ []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "removeBatch", tokenID_, attrIndexes_)
}

// RemoveBatch is a paid mutator transaction binding the contract method 0x350ce9e8.
//
// Solidity: function removeBatch(uint256 tokenID_, uint256[] attrIndexes_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) RemoveBatch(tokenID_ *big.Int, attrIndexes_ []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.RemoveBatch(&_GameLootEquipment.TransactOpts, tokenID_, attrIndexes_)
}

// RemoveBatch is a paid mutator transaction binding the contract method 0x350ce9e8.
//
// Solidity: function removeBatch(uint256 tokenID_, uint256[] attrIndexes_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) RemoveBatch(tokenID_ *big.Int, attrIndexes_ []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.RemoveBatch(&_GameLootEquipment.TransactOpts, tokenID_, attrIndexes_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GameLootEquipment *GameLootEquipmentSession) RenounceOwnership() (*types.Transaction, error) {
	return _GameLootEquipment.Contract.RenounceOwnership(&_GameLootEquipment.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GameLootEquipment.Contract.RenounceOwnership(&_GameLootEquipment.TransactOpts)
}

// Reveal is a paid mutator transaction binding the contract method 0x4d95595e.
//
// Solidity: function reveal(uint256 tokenID_, uint256 nonce_, uint128[] attrIDs_, uint128[] attrValues_, bytes signature_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Reveal(opts *bind.TransactOpts, tokenID_ *big.Int, nonce_ *big.Int, attrIDs_ []*big.Int, attrValues_ []*big.Int, signature_ []byte) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "reveal", tokenID_, nonce_, attrIDs_, attrValues_, signature_)
}

// Reveal is a paid mutator transaction binding the contract method 0x4d95595e.
//
// Solidity: function reveal(uint256 tokenID_, uint256 nonce_, uint128[] attrIDs_, uint128[] attrValues_, bytes signature_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) Reveal(tokenID_ *big.Int, nonce_ *big.Int, attrIDs_ []*big.Int, attrValues_ []*big.Int, signature_ []byte) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Reveal(&_GameLootEquipment.TransactOpts, tokenID_, nonce_, attrIDs_, attrValues_, signature_)
}

// Reveal is a paid mutator transaction binding the contract method 0x4d95595e.
//
// Solidity: function reveal(uint256 tokenID_, uint256 nonce_, uint128[] attrIDs_, uint128[] attrValues_, bytes signature_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Reveal(tokenID_ *big.Int, nonce_ *big.Int, attrIDs_ []*big.Int, attrValues_ []*big.Int, signature_ []byte) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Reveal(&_GameLootEquipment.TransactOpts, tokenID_, nonce_, attrIDs_, attrValues_, signature_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SafeTransferFrom(&_GameLootEquipment.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SafeTransferFrom(&_GameLootEquipment.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SafeTransferFrom0(&_GameLootEquipment.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SafeTransferFrom0(&_GameLootEquipment.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetApprovalForAll(&_GameLootEquipment.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetApprovalForAll(&_GameLootEquipment.TransactOpts, operator, approved)
}

// SetBlindBoxURI is a paid mutator transaction binding the contract method 0xb023b315.
//
// Solidity: function setBlindBoxURI(string blindBoxURI_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetBlindBoxURI(opts *bind.TransactOpts, blindBoxURI_ string) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setBlindBoxURI", blindBoxURI_)
}

// SetBlindBoxURI is a paid mutator transaction binding the contract method 0xb023b315.
//
// Solidity: function setBlindBoxURI(string blindBoxURI_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetBlindBoxURI(blindBoxURI_ string) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetBlindBoxURI(&_GameLootEquipment.TransactOpts, blindBoxURI_)
}

// SetBlindBoxURI is a paid mutator transaction binding the contract method 0xb023b315.
//
// Solidity: function setBlindBoxURI(string blindBoxURI_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetBlindBoxURI(blindBoxURI_ string) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetBlindBoxURI(&_GameLootEquipment.TransactOpts, blindBoxURI_)
}

// SetCap is a paid mutator transaction binding the contract method 0x47786d37.
//
// Solidity: function setCap(uint256 cap) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetCap(opts *bind.TransactOpts, cap *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setCap", cap)
}

// SetCap is a paid mutator transaction binding the contract method 0x47786d37.
//
// Solidity: function setCap(uint256 cap) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetCap(cap *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetCap(&_GameLootEquipment.TransactOpts, cap)
}

// SetCap is a paid mutator transaction binding the contract method 0x47786d37.
//
// Solidity: function setCap(uint256 cap) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetCap(cap *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetCap(&_GameLootEquipment.TransactOpts, cap)
}

// SetGameMinter is a paid mutator transaction binding the contract method 0x6ce6cf1e.
//
// Solidity: function setGameMinter(address gameMinter_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetGameMinter(opts *bind.TransactOpts, gameMinter_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setGameMinter", gameMinter_)
}

// SetGameMinter is a paid mutator transaction binding the contract method 0x6ce6cf1e.
//
// Solidity: function setGameMinter(address gameMinter_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetGameMinter(gameMinter_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetGameMinter(&_GameLootEquipment.TransactOpts, gameMinter_)
}

// SetGameMinter is a paid mutator transaction binding the contract method 0x6ce6cf1e.
//
// Solidity: function setGameMinter(address gameMinter_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetGameMinter(gameMinter_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetGameMinter(&_GameLootEquipment.TransactOpts, gameMinter_)
}

// SetReveal is a paid mutator transaction binding the contract method 0x2c035d37.
//
// Solidity: function setReveal(address revealSVG_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetReveal(opts *bind.TransactOpts, revealSVG_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setReveal", revealSVG_)
}

// SetReveal is a paid mutator transaction binding the contract method 0x2c035d37.
//
// Solidity: function setReveal(address revealSVG_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetReveal(revealSVG_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetReveal(&_GameLootEquipment.TransactOpts, revealSVG_)
}

// SetReveal is a paid mutator transaction binding the contract method 0x2c035d37.
//
// Solidity: function setReveal(address revealSVG_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetReveal(revealSVG_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetReveal(&_GameLootEquipment.TransactOpts, revealSVG_)
}

// SetRevealed is a paid mutator transaction binding the contract method 0x2fd5bce2.
//
// Solidity: function setRevealed(uint256 tokenID) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetRevealed(opts *bind.TransactOpts, tokenID *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setRevealed", tokenID)
}

// SetRevealed is a paid mutator transaction binding the contract method 0x2fd5bce2.
//
// Solidity: function setRevealed(uint256 tokenID) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetRevealed(tokenID *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetRevealed(&_GameLootEquipment.TransactOpts, tokenID)
}

// SetRevealed is a paid mutator transaction binding the contract method 0x2fd5bce2.
//
// Solidity: function setRevealed(uint256 tokenID) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetRevealed(tokenID *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetRevealed(&_GameLootEquipment.TransactOpts, tokenID)
}

// SetSeller is a paid mutator transaction binding the contract method 0xe99d2866.
//
// Solidity: function setSeller(address seller_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetSeller(opts *bind.TransactOpts, seller_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setSeller", seller_)
}

// SetSeller is a paid mutator transaction binding the contract method 0xe99d2866.
//
// Solidity: function setSeller(address seller_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetSeller(seller_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetSeller(&_GameLootEquipment.TransactOpts, seller_)
}

// SetSeller is a paid mutator transaction binding the contract method 0xe99d2866.
//
// Solidity: function setSeller(address seller_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetSeller(seller_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetSeller(&_GameLootEquipment.TransactOpts, seller_)
}

// SetSigner is a paid mutator transaction binding the contract method 0x31cb6105.
//
// Solidity: function setSigner(address signer, bool isOk) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetSigner(opts *bind.TransactOpts, signer common.Address, isOk bool) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setSigner", signer, isOk)
}

// SetSigner is a paid mutator transaction binding the contract method 0x31cb6105.
//
// Solidity: function setSigner(address signer, bool isOk) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetSigner(signer common.Address, isOk bool) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetSigner(&_GameLootEquipment.TransactOpts, signer, isOk)
}

// SetSigner is a paid mutator transaction binding the contract method 0x31cb6105.
//
// Solidity: function setSigner(address signer, bool isOk) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetSigner(signer common.Address, isOk bool) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetSigner(&_GameLootEquipment.TransactOpts, signer, isOk)
}

// SetTimeLocker is a paid mutator transaction binding the contract method 0x07aa1043.
//
// Solidity: function setTimeLocker(address timeLocker_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetTimeLocker(opts *bind.TransactOpts, timeLocker_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setTimeLocker", timeLocker_)
}

// SetTimeLocker is a paid mutator transaction binding the contract method 0x07aa1043.
//
// Solidity: function setTimeLocker(address timeLocker_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetTimeLocker(timeLocker_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetTimeLocker(&_GameLootEquipment.TransactOpts, timeLocker_)
}

// SetTimeLocker is a paid mutator transaction binding the contract method 0x07aa1043.
//
// Solidity: function setTimeLocker(address timeLocker_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetTimeLocker(timeLocker_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetTimeLocker(&_GameLootEquipment.TransactOpts, timeLocker_)
}

// SetTreasure is a paid mutator transaction binding the contract method 0x1d3cd2a6.
//
// Solidity: function setTreasure(address treasure_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetTreasure(opts *bind.TransactOpts, treasure_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setTreasure", treasure_)
}

// SetTreasure is a paid mutator transaction binding the contract method 0x1d3cd2a6.
//
// Solidity: function setTreasure(address treasure_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetTreasure(treasure_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetTreasure(&_GameLootEquipment.TransactOpts, treasure_)
}

// SetTreasure is a paid mutator transaction binding the contract method 0x1d3cd2a6.
//
// Solidity: function setTreasure(address treasure_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetTreasure(treasure_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetTreasure(&_GameLootEquipment.TransactOpts, treasure_)
}

// SetVault is a paid mutator transaction binding the contract method 0x6817031b.
//
// Solidity: function setVault(address vault_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) SetVault(opts *bind.TransactOpts, vault_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "setVault", vault_)
}

// SetVault is a paid mutator transaction binding the contract method 0x6817031b.
//
// Solidity: function setVault(address vault_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) SetVault(vault_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetVault(&_GameLootEquipment.TransactOpts, vault_)
}

// SetVault is a paid mutator transaction binding the contract method 0x6817031b.
//
// Solidity: function setVault(address vault_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) SetVault(vault_ common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.SetVault(&_GameLootEquipment.TransactOpts, vault_)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.TransferFrom(&_GameLootEquipment.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.TransferFrom(&_GameLootEquipment.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GameLootEquipment *GameLootEquipmentSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.TransferOwnership(&_GameLootEquipment.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.TransferOwnership(&_GameLootEquipment.TransactOpts, newOwner)
}

// Update is a paid mutator transaction binding the contract method 0x5fbc3938.
//
// Solidity: function update(uint256 tokenID_, uint256 attrIndex_, uint128 value_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Update(opts *bind.TransactOpts, tokenID_ *big.Int, attrIndex_ *big.Int, value_ *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "update", tokenID_, attrIndex_, value_)
}

// Update is a paid mutator transaction binding the contract method 0x5fbc3938.
//
// Solidity: function update(uint256 tokenID_, uint256 attrIndex_, uint128 value_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) Update(tokenID_ *big.Int, attrIndex_ *big.Int, value_ *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Update(&_GameLootEquipment.TransactOpts, tokenID_, attrIndex_, value_)
}

// Update is a paid mutator transaction binding the contract method 0x5fbc3938.
//
// Solidity: function update(uint256 tokenID_, uint256 attrIndex_, uint128 value_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Update(tokenID_ *big.Int, attrIndex_ *big.Int, value_ *big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Update(&_GameLootEquipment.TransactOpts, tokenID_, attrIndex_, value_)
}

// UpdateBatch is a paid mutator transaction binding the contract method 0x499e0547.
//
// Solidity: function updateBatch(uint256 tokenID_, uint256[] attrIndexes_, uint128[] values_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) UpdateBatch(opts *bind.TransactOpts, tokenID_ *big.Int, attrIndexes_ []*big.Int, values_ []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "updateBatch", tokenID_, attrIndexes_, values_)
}

// UpdateBatch is a paid mutator transaction binding the contract method 0x499e0547.
//
// Solidity: function updateBatch(uint256 tokenID_, uint256[] attrIndexes_, uint128[] values_) returns()
func (_GameLootEquipment *GameLootEquipmentSession) UpdateBatch(tokenID_ *big.Int, attrIndexes_ []*big.Int, values_ []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.UpdateBatch(&_GameLootEquipment.TransactOpts, tokenID_, attrIndexes_, values_)
}

// UpdateBatch is a paid mutator transaction binding the contract method 0x499e0547.
//
// Solidity: function updateBatch(uint256 tokenID_, uint256[] attrIndexes_, uint128[] values_) returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) UpdateBatch(tokenID_ *big.Int, attrIndexes_ []*big.Int, values_ []*big.Int) (*types.Transaction, error) {
	return _GameLootEquipment.Contract.UpdateBatch(&_GameLootEquipment.TransactOpts, tokenID_, attrIndexes_, values_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GameLootEquipment.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_GameLootEquipment *GameLootEquipmentSession) Withdraw() (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Withdraw(&_GameLootEquipment.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Withdraw() (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Withdraw(&_GameLootEquipment.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GameLootEquipment *GameLootEquipmentTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GameLootEquipment.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GameLootEquipment *GameLootEquipmentSession) Receive() (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Receive(&_GameLootEquipment.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GameLootEquipment *GameLootEquipmentTransactorSession) Receive() (*types.Transaction, error) {
	return _GameLootEquipment.Contract.Receive(&_GameLootEquipment.TransactOpts)
}

// GameLootEquipmentApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GameLootEquipment contract.
type GameLootEquipmentApprovalIterator struct {
	Event *GameLootEquipmentApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentApproval represents a Approval event raised by the GameLootEquipment contract.
type GameLootEquipmentApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*GameLootEquipmentApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentApprovalIterator{contract: _GameLootEquipment.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentApproval)
				if err := _GameLootEquipment.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseApproval(log types.Log) (*GameLootEquipmentApproval, error) {
	event := new(GameLootEquipmentApproval)
	if err := _GameLootEquipment.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the GameLootEquipment contract.
type GameLootEquipmentApprovalForAllIterator struct {
	Event *GameLootEquipmentApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentApprovalForAll represents a ApprovalForAll event raised by the GameLootEquipment contract.
type GameLootEquipmentApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*GameLootEquipmentApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentApprovalForAllIterator{contract: _GameLootEquipment.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentApprovalForAll)
				if err := _GameLootEquipment.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseApprovalForAll(log types.Log) (*GameLootEquipmentApprovalForAll, error) {
	event := new(GameLootEquipmentApprovalForAll)
	if err := _GameLootEquipment.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentAttributeAttachedIterator is returned from FilterAttributeAttached and is used to iterate over the raw logs and unpacked data for AttributeAttached events raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeAttachedIterator struct {
	Event *GameLootEquipmentAttributeAttached // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentAttributeAttachedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentAttributeAttached)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentAttributeAttached)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentAttributeAttachedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentAttributeAttachedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentAttributeAttached represents a AttributeAttached event raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeAttached struct {
	TokenID *big.Int
	AttrID  *big.Int
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAttributeAttached is a free log retrieval operation binding the contract event 0xaf7478c51afd39672cf93c314e0f346b8018d02d74e430338895168ca2686fe5.
//
// Solidity: event AttributeAttached(uint256 tokenID, uint128 attrID, uint128 value)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterAttributeAttached(opts *bind.FilterOpts) (*GameLootEquipmentAttributeAttachedIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "AttributeAttached")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentAttributeAttachedIterator{contract: _GameLootEquipment.contract, event: "AttributeAttached", logs: logs, sub: sub}, nil
}

// WatchAttributeAttached is a free log subscription operation binding the contract event 0xaf7478c51afd39672cf93c314e0f346b8018d02d74e430338895168ca2686fe5.
//
// Solidity: event AttributeAttached(uint256 tokenID, uint128 attrID, uint128 value)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchAttributeAttached(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentAttributeAttached) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "AttributeAttached")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentAttributeAttached)
				if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeAttached", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAttributeAttached is a log parse operation binding the contract event 0xaf7478c51afd39672cf93c314e0f346b8018d02d74e430338895168ca2686fe5.
//
// Solidity: event AttributeAttached(uint256 tokenID, uint128 attrID, uint128 value)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseAttributeAttached(log types.Log) (*GameLootEquipmentAttributeAttached, error) {
	event := new(GameLootEquipmentAttributeAttached)
	if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeAttached", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentAttributeAttachedBatchIterator is returned from FilterAttributeAttachedBatch and is used to iterate over the raw logs and unpacked data for AttributeAttachedBatch events raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeAttachedBatchIterator struct {
	Event *GameLootEquipmentAttributeAttachedBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentAttributeAttachedBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentAttributeAttachedBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentAttributeAttachedBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentAttributeAttachedBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentAttributeAttachedBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentAttributeAttachedBatch represents a AttributeAttachedBatch event raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeAttachedBatch struct {
	TokenID *big.Int
	AttrIDs []*big.Int
	Values  []*big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAttributeAttachedBatch is a free log retrieval operation binding the contract event 0x3bea8271f01d46ef999dd7066026123b25201f6108ed34a13f3de7ed4b1d765b.
//
// Solidity: event AttributeAttachedBatch(uint256 tokenID, uint128[] attrIDs, uint128[] values)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterAttributeAttachedBatch(opts *bind.FilterOpts) (*GameLootEquipmentAttributeAttachedBatchIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "AttributeAttachedBatch")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentAttributeAttachedBatchIterator{contract: _GameLootEquipment.contract, event: "AttributeAttachedBatch", logs: logs, sub: sub}, nil
}

// WatchAttributeAttachedBatch is a free log subscription operation binding the contract event 0x3bea8271f01d46ef999dd7066026123b25201f6108ed34a13f3de7ed4b1d765b.
//
// Solidity: event AttributeAttachedBatch(uint256 tokenID, uint128[] attrIDs, uint128[] values)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchAttributeAttachedBatch(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentAttributeAttachedBatch) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "AttributeAttachedBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentAttributeAttachedBatch)
				if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeAttachedBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAttributeAttachedBatch is a log parse operation binding the contract event 0x3bea8271f01d46ef999dd7066026123b25201f6108ed34a13f3de7ed4b1d765b.
//
// Solidity: event AttributeAttachedBatch(uint256 tokenID, uint128[] attrIDs, uint128[] values)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseAttributeAttachedBatch(log types.Log) (*GameLootEquipmentAttributeAttachedBatch, error) {
	event := new(GameLootEquipmentAttributeAttachedBatch)
	if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeAttachedBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentAttributeRemoveBatchIterator is returned from FilterAttributeRemoveBatch and is used to iterate over the raw logs and unpacked data for AttributeRemoveBatch events raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeRemoveBatchIterator struct {
	Event *GameLootEquipmentAttributeRemoveBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentAttributeRemoveBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentAttributeRemoveBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentAttributeRemoveBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentAttributeRemoveBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentAttributeRemoveBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentAttributeRemoveBatch represents a AttributeRemoveBatch event raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeRemoveBatch struct {
	TokenID *big.Int
	AttrIDs []*big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAttributeRemoveBatch is a free log retrieval operation binding the contract event 0xb26bb644f8b28b119909f1d697aebd30afd7c141a0d3fb47faf163ee81c02b71.
//
// Solidity: event AttributeRemoveBatch(uint256 tokenID, uint128[] attrIDs)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterAttributeRemoveBatch(opts *bind.FilterOpts) (*GameLootEquipmentAttributeRemoveBatchIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "AttributeRemoveBatch")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentAttributeRemoveBatchIterator{contract: _GameLootEquipment.contract, event: "AttributeRemoveBatch", logs: logs, sub: sub}, nil
}

// WatchAttributeRemoveBatch is a free log subscription operation binding the contract event 0xb26bb644f8b28b119909f1d697aebd30afd7c141a0d3fb47faf163ee81c02b71.
//
// Solidity: event AttributeRemoveBatch(uint256 tokenID, uint128[] attrIDs)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchAttributeRemoveBatch(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentAttributeRemoveBatch) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "AttributeRemoveBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentAttributeRemoveBatch)
				if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeRemoveBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAttributeRemoveBatch is a log parse operation binding the contract event 0xb26bb644f8b28b119909f1d697aebd30afd7c141a0d3fb47faf163ee81c02b71.
//
// Solidity: event AttributeRemoveBatch(uint256 tokenID, uint128[] attrIDs)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseAttributeRemoveBatch(log types.Log) (*GameLootEquipmentAttributeRemoveBatch, error) {
	event := new(GameLootEquipmentAttributeRemoveBatch)
	if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeRemoveBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentAttributeRemovedIterator is returned from FilterAttributeRemoved and is used to iterate over the raw logs and unpacked data for AttributeRemoved events raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeRemovedIterator struct {
	Event *GameLootEquipmentAttributeRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentAttributeRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentAttributeRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentAttributeRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentAttributeRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentAttributeRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentAttributeRemoved represents a AttributeRemoved event raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeRemoved struct {
	TokenID *big.Int
	AttrID  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAttributeRemoved is a free log retrieval operation binding the contract event 0x762f9e50ef18e928c58da3c41eb989fa4f4b2870843eff7043d57cf27fb99435.
//
// Solidity: event AttributeRemoved(uint256 tokenID, uint128 attrID)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterAttributeRemoved(opts *bind.FilterOpts) (*GameLootEquipmentAttributeRemovedIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "AttributeRemoved")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentAttributeRemovedIterator{contract: _GameLootEquipment.contract, event: "AttributeRemoved", logs: logs, sub: sub}, nil
}

// WatchAttributeRemoved is a free log subscription operation binding the contract event 0x762f9e50ef18e928c58da3c41eb989fa4f4b2870843eff7043d57cf27fb99435.
//
// Solidity: event AttributeRemoved(uint256 tokenID, uint128 attrID)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchAttributeRemoved(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentAttributeRemoved) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "AttributeRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentAttributeRemoved)
				if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAttributeRemoved is a log parse operation binding the contract event 0x762f9e50ef18e928c58da3c41eb989fa4f4b2870843eff7043d57cf27fb99435.
//
// Solidity: event AttributeRemoved(uint256 tokenID, uint128 attrID)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseAttributeRemoved(log types.Log) (*GameLootEquipmentAttributeRemoved, error) {
	event := new(GameLootEquipmentAttributeRemoved)
	if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentAttributeUpdatedIterator is returned from FilterAttributeUpdated and is used to iterate over the raw logs and unpacked data for AttributeUpdated events raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeUpdatedIterator struct {
	Event *GameLootEquipmentAttributeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentAttributeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentAttributeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentAttributeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentAttributeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentAttributeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentAttributeUpdated represents a AttributeUpdated event raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeUpdated struct {
	TokenID   *big.Int
	AttrIndex *big.Int
	Value     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAttributeUpdated is a free log retrieval operation binding the contract event 0x26fd96bc3854a9c47a23ec63704691e24a72fb10c15c5689b994764b2606cb3e.
//
// Solidity: event AttributeUpdated(uint256 tokenID, uint256 attrIndex, uint128 value)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterAttributeUpdated(opts *bind.FilterOpts) (*GameLootEquipmentAttributeUpdatedIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "AttributeUpdated")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentAttributeUpdatedIterator{contract: _GameLootEquipment.contract, event: "AttributeUpdated", logs: logs, sub: sub}, nil
}

// WatchAttributeUpdated is a free log subscription operation binding the contract event 0x26fd96bc3854a9c47a23ec63704691e24a72fb10c15c5689b994764b2606cb3e.
//
// Solidity: event AttributeUpdated(uint256 tokenID, uint256 attrIndex, uint128 value)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchAttributeUpdated(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentAttributeUpdated) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "AttributeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentAttributeUpdated)
				if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAttributeUpdated is a log parse operation binding the contract event 0x26fd96bc3854a9c47a23ec63704691e24a72fb10c15c5689b994764b2606cb3e.
//
// Solidity: event AttributeUpdated(uint256 tokenID, uint256 attrIndex, uint128 value)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseAttributeUpdated(log types.Log) (*GameLootEquipmentAttributeUpdated, error) {
	event := new(GameLootEquipmentAttributeUpdated)
	if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentAttributeUpdatedBatchIterator is returned from FilterAttributeUpdatedBatch and is used to iterate over the raw logs and unpacked data for AttributeUpdatedBatch events raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeUpdatedBatchIterator struct {
	Event *GameLootEquipmentAttributeUpdatedBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentAttributeUpdatedBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentAttributeUpdatedBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentAttributeUpdatedBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentAttributeUpdatedBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentAttributeUpdatedBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentAttributeUpdatedBatch represents a AttributeUpdatedBatch event raised by the GameLootEquipment contract.
type GameLootEquipmentAttributeUpdatedBatch struct {
	TokenID     *big.Int
	AttrIndexes []*big.Int
	Values      []*big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAttributeUpdatedBatch is a free log retrieval operation binding the contract event 0xef13e189c87a2b993d3fdec515d7c85b19464ca8a2464d5f6ff7711bc4412c78.
//
// Solidity: event AttributeUpdatedBatch(uint256 tokenID, uint256[] attrIndexes, uint128[] values)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterAttributeUpdatedBatch(opts *bind.FilterOpts) (*GameLootEquipmentAttributeUpdatedBatchIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "AttributeUpdatedBatch")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentAttributeUpdatedBatchIterator{contract: _GameLootEquipment.contract, event: "AttributeUpdatedBatch", logs: logs, sub: sub}, nil
}

// WatchAttributeUpdatedBatch is a free log subscription operation binding the contract event 0xef13e189c87a2b993d3fdec515d7c85b19464ca8a2464d5f6ff7711bc4412c78.
//
// Solidity: event AttributeUpdatedBatch(uint256 tokenID, uint256[] attrIndexes, uint128[] values)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchAttributeUpdatedBatch(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentAttributeUpdatedBatch) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "AttributeUpdatedBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentAttributeUpdatedBatch)
				if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeUpdatedBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAttributeUpdatedBatch is a log parse operation binding the contract event 0xef13e189c87a2b993d3fdec515d7c85b19464ca8a2464d5f6ff7711bc4412c78.
//
// Solidity: event AttributeUpdatedBatch(uint256 tokenID, uint256[] attrIndexes, uint128[] values)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseAttributeUpdatedBatch(log types.Log) (*GameLootEquipmentAttributeUpdatedBatch, error) {
	event := new(GameLootEquipmentAttributeUpdatedBatch)
	if err := _GameLootEquipment.contract.UnpackLog(event, "AttributeUpdatedBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GameLootEquipment contract.
type GameLootEquipmentOwnershipTransferredIterator struct {
	Event *GameLootEquipmentOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentOwnershipTransferred represents a OwnershipTransferred event raised by the GameLootEquipment contract.
type GameLootEquipmentOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GameLootEquipmentOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentOwnershipTransferredIterator{contract: _GameLootEquipment.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentOwnershipTransferred)
				if err := _GameLootEquipment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseOwnershipTransferred(log types.Log) (*GameLootEquipmentOwnershipTransferred, error) {
	event := new(GameLootEquipmentOwnershipTransferred)
	if err := _GameLootEquipment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentRevealedIterator is returned from FilterRevealed and is used to iterate over the raw logs and unpacked data for Revealed events raised by the GameLootEquipment contract.
type GameLootEquipmentRevealedIterator struct {
	Event *GameLootEquipmentRevealed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentRevealedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentRevealed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentRevealed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentRevealedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentRevealedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentRevealed represents a Revealed event raised by the GameLootEquipment contract.
type GameLootEquipmentRevealed struct {
	TokenID *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRevealed is a free log retrieval operation binding the contract event 0x15120e52505e619cbf6c2af910d5cf7f9ee1befa55801b078c33e93880b2d609.
//
// Solidity: event Revealed(uint256 tokenID)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterRevealed(opts *bind.FilterOpts) (*GameLootEquipmentRevealedIterator, error) {

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "Revealed")
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentRevealedIterator{contract: _GameLootEquipment.contract, event: "Revealed", logs: logs, sub: sub}, nil
}

// WatchRevealed is a free log subscription operation binding the contract event 0x15120e52505e619cbf6c2af910d5cf7f9ee1befa55801b078c33e93880b2d609.
//
// Solidity: event Revealed(uint256 tokenID)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchRevealed(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentRevealed) (event.Subscription, error) {

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "Revealed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentRevealed)
				if err := _GameLootEquipment.contract.UnpackLog(event, "Revealed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRevealed is a log parse operation binding the contract event 0x15120e52505e619cbf6c2af910d5cf7f9ee1befa55801b078c33e93880b2d609.
//
// Solidity: event Revealed(uint256 tokenID)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseRevealed(log types.Log) (*GameLootEquipmentRevealed, error) {
	event := new(GameLootEquipmentRevealed)
	if err := _GameLootEquipment.contract.UnpackLog(event, "Revealed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GameLootEquipmentTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GameLootEquipment contract.
type GameLootEquipmentTransferIterator struct {
	Event *GameLootEquipmentTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GameLootEquipmentTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GameLootEquipmentTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GameLootEquipmentTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GameLootEquipmentTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GameLootEquipmentTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GameLootEquipmentTransfer represents a Transfer event raised by the GameLootEquipment contract.
type GameLootEquipmentTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GameLootEquipment *GameLootEquipmentFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*GameLootEquipmentTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _GameLootEquipment.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GameLootEquipmentTransferIterator{contract: _GameLootEquipment.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GameLootEquipment *GameLootEquipmentFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GameLootEquipmentTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _GameLootEquipment.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GameLootEquipmentTransfer)
				if err := _GameLootEquipment.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GameLootEquipment *GameLootEquipmentFilterer) ParseTransfer(log types.Log) (*GameLootEquipmentTransfer, error) {
	event := new(GameLootEquipmentTransfer)
	if err := _GameLootEquipment.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
