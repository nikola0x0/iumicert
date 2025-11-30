// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

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
	_ = abi.ConvertType
)

// IUMiCertRegistryMetaData contains all meta data concerning the IUMiCertRegistry contract.
var IUMiCertRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkRootStatus\",\"inputs\":[{\"name\":\"_verkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"latestRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"message\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestRoot\",\"inputs\":[{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"version\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublishedRoot\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublishedRootsCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublishedTerm\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublishedTermsCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTermHistory\",\"inputs\":[{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"versions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"roots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVersionInfo\",\"inputs\":[{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_version\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isSuperseded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"supersededBy\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"supersessionReason\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isTermPublished\",\"inputs\":[{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestVersion\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"publishTermRoot\",\"inputs\":[{\"name\":\"_verkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"publishedRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"publishedTerms\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rootToTerm\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rootToVersion\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supersedeTerm\",\"inputs\":[{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_newVerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_newTotalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"termVersions\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"version\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isSuperseded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"supersededBy\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"supersessionReason\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyReceiptAnchor\",\"inputs\":[{\"name\":\"_blockchainAnchor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TermRootPublished\",\"inputs\":[{\"name\":\"termId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"verkleRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"version\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TermRootSuperseded\",\"inputs\":[{\"name\":\"termId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"oldVersion\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"oldRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newVersion\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// IUMiCertRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IUMiCertRegistryMetaData.ABI instead.
var IUMiCertRegistryABI = IUMiCertRegistryMetaData.ABI

// IUMiCertRegistry is an auto generated Go binding around an Ethereum contract.
type IUMiCertRegistry struct {
	IUMiCertRegistryCaller     // Read-only binding to the contract
	IUMiCertRegistryTransactor // Write-only binding to the contract
	IUMiCertRegistryFilterer   // Log filterer for contract events
}

// IUMiCertRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUMiCertRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUMiCertRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUMiCertRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUMiCertRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUMiCertRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUMiCertRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUMiCertRegistrySession struct {
	Contract     *IUMiCertRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUMiCertRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUMiCertRegistryCallerSession struct {
	Contract *IUMiCertRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IUMiCertRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUMiCertRegistryTransactorSession struct {
	Contract     *IUMiCertRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IUMiCertRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUMiCertRegistryRaw struct {
	Contract *IUMiCertRegistry // Generic contract binding to access the raw methods on
}

// IUMiCertRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUMiCertRegistryCallerRaw struct {
	Contract *IUMiCertRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IUMiCertRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUMiCertRegistryTransactorRaw struct {
	Contract *IUMiCertRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUMiCertRegistry creates a new instance of IUMiCertRegistry, bound to a specific deployed contract.
func NewIUMiCertRegistry(address common.Address, backend bind.ContractBackend) (*IUMiCertRegistry, error) {
	contract, err := bindIUMiCertRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistry{IUMiCertRegistryCaller: IUMiCertRegistryCaller{contract: contract}, IUMiCertRegistryTransactor: IUMiCertRegistryTransactor{contract: contract}, IUMiCertRegistryFilterer: IUMiCertRegistryFilterer{contract: contract}}, nil
}

// NewIUMiCertRegistryCaller creates a new read-only instance of IUMiCertRegistry, bound to a specific deployed contract.
func NewIUMiCertRegistryCaller(address common.Address, caller bind.ContractCaller) (*IUMiCertRegistryCaller, error) {
	contract, err := bindIUMiCertRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryCaller{contract: contract}, nil
}

// NewIUMiCertRegistryTransactor creates a new write-only instance of IUMiCertRegistry, bound to a specific deployed contract.
func NewIUMiCertRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IUMiCertRegistryTransactor, error) {
	contract, err := bindIUMiCertRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryTransactor{contract: contract}, nil
}

// NewIUMiCertRegistryFilterer creates a new log filterer instance of IUMiCertRegistry, bound to a specific deployed contract.
func NewIUMiCertRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IUMiCertRegistryFilterer, error) {
	contract, err := bindIUMiCertRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryFilterer{contract: contract}, nil
}

// bindIUMiCertRegistry binds a generic wrapper to an already deployed contract.
func bindIUMiCertRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUMiCertRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUMiCertRegistry *IUMiCertRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUMiCertRegistry.Contract.IUMiCertRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUMiCertRegistry *IUMiCertRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.IUMiCertRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUMiCertRegistry *IUMiCertRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.IUMiCertRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUMiCertRegistry *IUMiCertRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUMiCertRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUMiCertRegistry *IUMiCertRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUMiCertRegistry *IUMiCertRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.contract.Transact(opts, method, params...)
}

// CheckRootStatus is a free data retrieval call binding the contract method 0x92e1292a.
//
// Solidity: function checkRootStatus(bytes32 _verkleRoot) view returns(uint8 status, string termId, uint256 version, bytes32 latestRoot, string message)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) CheckRootStatus(opts *bind.CallOpts, _verkleRoot [32]byte) (struct {
	Status     uint8
	TermId     string
	Version    *big.Int
	LatestRoot [32]byte
	Message    string
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "checkRootStatus", _verkleRoot)

	outstruct := new(struct {
		Status     uint8
		TermId     string
		Version    *big.Int
		LatestRoot [32]byte
		Message    string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Status = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.TermId = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LatestRoot = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Message = *abi.ConvertType(out[4], new(string)).(*string)

	return *outstruct, err

}

// CheckRootStatus is a free data retrieval call binding the contract method 0x92e1292a.
//
// Solidity: function checkRootStatus(bytes32 _verkleRoot) view returns(uint8 status, string termId, uint256 version, bytes32 latestRoot, string message)
func (_IUMiCertRegistry *IUMiCertRegistrySession) CheckRootStatus(_verkleRoot [32]byte) (struct {
	Status     uint8
	TermId     string
	Version    *big.Int
	LatestRoot [32]byte
	Message    string
}, error) {
	return _IUMiCertRegistry.Contract.CheckRootStatus(&_IUMiCertRegistry.CallOpts, _verkleRoot)
}

// CheckRootStatus is a free data retrieval call binding the contract method 0x92e1292a.
//
// Solidity: function checkRootStatus(bytes32 _verkleRoot) view returns(uint8 status, string termId, uint256 version, bytes32 latestRoot, string message)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) CheckRootStatus(_verkleRoot [32]byte) (struct {
	Status     uint8
	TermId     string
	Version    *big.Int
	LatestRoot [32]byte
	Message    string
}, error) {
	return _IUMiCertRegistry.Contract.CheckRootStatus(&_IUMiCertRegistry.CallOpts, _verkleRoot)
}

// GetLatestRoot is a free data retrieval call binding the contract method 0xf8e47cee.
//
// Solidity: function getLatestRoot(string _termId) view returns(bytes32 rootHash, uint256 version, uint256 totalStudents, uint256 publishedAt)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetLatestRoot(opts *bind.CallOpts, _termId string) (struct {
	RootHash      [32]byte
	Version       *big.Int
	TotalStudents *big.Int
	PublishedAt   *big.Int
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getLatestRoot", _termId)

	outstruct := new(struct {
		RootHash      [32]byte
		Version       *big.Int
		TotalStudents *big.Int
		PublishedAt   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RootHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Version = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalStudents = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PublishedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetLatestRoot is a free data retrieval call binding the contract method 0xf8e47cee.
//
// Solidity: function getLatestRoot(string _termId) view returns(bytes32 rootHash, uint256 version, uint256 totalStudents, uint256 publishedAt)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetLatestRoot(_termId string) (struct {
	RootHash      [32]byte
	Version       *big.Int
	TotalStudents *big.Int
	PublishedAt   *big.Int
}, error) {
	return _IUMiCertRegistry.Contract.GetLatestRoot(&_IUMiCertRegistry.CallOpts, _termId)
}

// GetLatestRoot is a free data retrieval call binding the contract method 0xf8e47cee.
//
// Solidity: function getLatestRoot(string _termId) view returns(bytes32 rootHash, uint256 version, uint256 totalStudents, uint256 publishedAt)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetLatestRoot(_termId string) (struct {
	RootHash      [32]byte
	Version       *big.Int
	TotalStudents *big.Int
	PublishedAt   *big.Int
}, error) {
	return _IUMiCertRegistry.Contract.GetLatestRoot(&_IUMiCertRegistry.CallOpts, _termId)
}

// GetPublishedRoot is a free data retrieval call binding the contract method 0x351863fb.
//
// Solidity: function getPublishedRoot(uint256 _index) view returns(bytes32)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetPublishedRoot(opts *bind.CallOpts, _index *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getPublishedRoot", _index)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPublishedRoot is a free data retrieval call binding the contract method 0x351863fb.
//
// Solidity: function getPublishedRoot(uint256 _index) view returns(bytes32)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetPublishedRoot(_index *big.Int) ([32]byte, error) {
	return _IUMiCertRegistry.Contract.GetPublishedRoot(&_IUMiCertRegistry.CallOpts, _index)
}

// GetPublishedRoot is a free data retrieval call binding the contract method 0x351863fb.
//
// Solidity: function getPublishedRoot(uint256 _index) view returns(bytes32)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetPublishedRoot(_index *big.Int) ([32]byte, error) {
	return _IUMiCertRegistry.Contract.GetPublishedRoot(&_IUMiCertRegistry.CallOpts, _index)
}

// GetPublishedRootsCount is a free data retrieval call binding the contract method 0xa2fbbc58.
//
// Solidity: function getPublishedRootsCount() view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetPublishedRootsCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getPublishedRootsCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPublishedRootsCount is a free data retrieval call binding the contract method 0xa2fbbc58.
//
// Solidity: function getPublishedRootsCount() view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetPublishedRootsCount() (*big.Int, error) {
	return _IUMiCertRegistry.Contract.GetPublishedRootsCount(&_IUMiCertRegistry.CallOpts)
}

// GetPublishedRootsCount is a free data retrieval call binding the contract method 0xa2fbbc58.
//
// Solidity: function getPublishedRootsCount() view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetPublishedRootsCount() (*big.Int, error) {
	return _IUMiCertRegistry.Contract.GetPublishedRootsCount(&_IUMiCertRegistry.CallOpts)
}

// GetPublishedTerm is a free data retrieval call binding the contract method 0x471af597.
//
// Solidity: function getPublishedTerm(uint256 _index) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetPublishedTerm(opts *bind.CallOpts, _index *big.Int) (string, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getPublishedTerm", _index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetPublishedTerm is a free data retrieval call binding the contract method 0x471af597.
//
// Solidity: function getPublishedTerm(uint256 _index) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetPublishedTerm(_index *big.Int) (string, error) {
	return _IUMiCertRegistry.Contract.GetPublishedTerm(&_IUMiCertRegistry.CallOpts, _index)
}

// GetPublishedTerm is a free data retrieval call binding the contract method 0x471af597.
//
// Solidity: function getPublishedTerm(uint256 _index) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetPublishedTerm(_index *big.Int) (string, error) {
	return _IUMiCertRegistry.Contract.GetPublishedTerm(&_IUMiCertRegistry.CallOpts, _index)
}

// GetPublishedTermsCount is a free data retrieval call binding the contract method 0x35d43b4c.
//
// Solidity: function getPublishedTermsCount() view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetPublishedTermsCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getPublishedTermsCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPublishedTermsCount is a free data retrieval call binding the contract method 0x35d43b4c.
//
// Solidity: function getPublishedTermsCount() view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetPublishedTermsCount() (*big.Int, error) {
	return _IUMiCertRegistry.Contract.GetPublishedTermsCount(&_IUMiCertRegistry.CallOpts)
}

// GetPublishedTermsCount is a free data retrieval call binding the contract method 0x35d43b4c.
//
// Solidity: function getPublishedTermsCount() view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetPublishedTermsCount() (*big.Int, error) {
	return _IUMiCertRegistry.Contract.GetPublishedTermsCount(&_IUMiCertRegistry.CallOpts)
}

// GetTermHistory is a free data retrieval call binding the contract method 0x24f8e323.
//
// Solidity: function getTermHistory(string _termId) view returns(uint256[] versions, bytes32[] roots)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetTermHistory(opts *bind.CallOpts, _termId string) (struct {
	Versions []*big.Int
	Roots    [][32]byte
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getTermHistory", _termId)

	outstruct := new(struct {
		Versions []*big.Int
		Roots    [][32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Versions = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.Roots = *abi.ConvertType(out[1], new([][32]byte)).(*[][32]byte)

	return *outstruct, err

}

// GetTermHistory is a free data retrieval call binding the contract method 0x24f8e323.
//
// Solidity: function getTermHistory(string _termId) view returns(uint256[] versions, bytes32[] roots)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetTermHistory(_termId string) (struct {
	Versions []*big.Int
	Roots    [][32]byte
}, error) {
	return _IUMiCertRegistry.Contract.GetTermHistory(&_IUMiCertRegistry.CallOpts, _termId)
}

// GetTermHistory is a free data retrieval call binding the contract method 0x24f8e323.
//
// Solidity: function getTermHistory(string _termId) view returns(uint256[] versions, bytes32[] roots)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetTermHistory(_termId string) (struct {
	Versions []*big.Int
	Roots    [][32]byte
}, error) {
	return _IUMiCertRegistry.Contract.GetTermHistory(&_IUMiCertRegistry.CallOpts, _termId)
}

// GetVersionInfo is a free data retrieval call binding the contract method 0xc2601afa.
//
// Solidity: function getVersionInfo(string _termId, uint256 _version) view returns(bytes32 rootHash, uint256 totalStudents, uint256 publishedAt, bool isSuperseded, bytes32 supersededBy, string supersessionReason)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetVersionInfo(opts *bind.CallOpts, _termId string, _version *big.Int) (struct {
	RootHash           [32]byte
	TotalStudents      *big.Int
	PublishedAt        *big.Int
	IsSuperseded       bool
	SupersededBy       [32]byte
	SupersessionReason string
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getVersionInfo", _termId, _version)

	outstruct := new(struct {
		RootHash           [32]byte
		TotalStudents      *big.Int
		PublishedAt        *big.Int
		IsSuperseded       bool
		SupersededBy       [32]byte
		SupersessionReason string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RootHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.TotalStudents = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PublishedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.IsSuperseded = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.SupersededBy = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.SupersessionReason = *abi.ConvertType(out[5], new(string)).(*string)

	return *outstruct, err

}

// GetVersionInfo is a free data retrieval call binding the contract method 0xc2601afa.
//
// Solidity: function getVersionInfo(string _termId, uint256 _version) view returns(bytes32 rootHash, uint256 totalStudents, uint256 publishedAt, bool isSuperseded, bytes32 supersededBy, string supersessionReason)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetVersionInfo(_termId string, _version *big.Int) (struct {
	RootHash           [32]byte
	TotalStudents      *big.Int
	PublishedAt        *big.Int
	IsSuperseded       bool
	SupersededBy       [32]byte
	SupersessionReason string
}, error) {
	return _IUMiCertRegistry.Contract.GetVersionInfo(&_IUMiCertRegistry.CallOpts, _termId, _version)
}

// GetVersionInfo is a free data retrieval call binding the contract method 0xc2601afa.
//
// Solidity: function getVersionInfo(string _termId, uint256 _version) view returns(bytes32 rootHash, uint256 totalStudents, uint256 publishedAt, bool isSuperseded, bytes32 supersededBy, string supersessionReason)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetVersionInfo(_termId string, _version *big.Int) (struct {
	RootHash           [32]byte
	TotalStudents      *big.Int
	PublishedAt        *big.Int
	IsSuperseded       bool
	SupersededBy       [32]byte
	SupersessionReason string
}, error) {
	return _IUMiCertRegistry.Contract.GetVersionInfo(&_IUMiCertRegistry.CallOpts, _termId, _version)
}

// IsTermPublished is a free data retrieval call binding the contract method 0x5365e9a4.
//
// Solidity: function isTermPublished(string _termId) view returns(bool)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) IsTermPublished(opts *bind.CallOpts, _termId string) (bool, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "isTermPublished", _termId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTermPublished is a free data retrieval call binding the contract method 0x5365e9a4.
//
// Solidity: function isTermPublished(string _termId) view returns(bool)
func (_IUMiCertRegistry *IUMiCertRegistrySession) IsTermPublished(_termId string) (bool, error) {
	return _IUMiCertRegistry.Contract.IsTermPublished(&_IUMiCertRegistry.CallOpts, _termId)
}

// IsTermPublished is a free data retrieval call binding the contract method 0x5365e9a4.
//
// Solidity: function isTermPublished(string _termId) view returns(bool)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) IsTermPublished(_termId string) (bool, error) {
	return _IUMiCertRegistry.Contract.IsTermPublished(&_IUMiCertRegistry.CallOpts, _termId)
}

// LatestVersion is a free data retrieval call binding the contract method 0x070f8f5a.
//
// Solidity: function latestVersion(string ) view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) LatestVersion(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "latestVersion", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestVersion is a free data retrieval call binding the contract method 0x070f8f5a.
//
// Solidity: function latestVersion(string ) view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistrySession) LatestVersion(arg0 string) (*big.Int, error) {
	return _IUMiCertRegistry.Contract.LatestVersion(&_IUMiCertRegistry.CallOpts, arg0)
}

// LatestVersion is a free data retrieval call binding the contract method 0x070f8f5a.
//
// Solidity: function latestVersion(string ) view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) LatestVersion(arg0 string) (*big.Int, error) {
	return _IUMiCertRegistry.Contract.LatestVersion(&_IUMiCertRegistry.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IUMiCertRegistry *IUMiCertRegistrySession) Owner() (common.Address, error) {
	return _IUMiCertRegistry.Contract.Owner(&_IUMiCertRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) Owner() (common.Address, error) {
	return _IUMiCertRegistry.Contract.Owner(&_IUMiCertRegistry.CallOpts)
}

// PublishedRoots is a free data retrieval call binding the contract method 0x744a92b6.
//
// Solidity: function publishedRoots(uint256 ) view returns(bytes32)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) PublishedRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "publishedRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PublishedRoots is a free data retrieval call binding the contract method 0x744a92b6.
//
// Solidity: function publishedRoots(uint256 ) view returns(bytes32)
func (_IUMiCertRegistry *IUMiCertRegistrySession) PublishedRoots(arg0 *big.Int) ([32]byte, error) {
	return _IUMiCertRegistry.Contract.PublishedRoots(&_IUMiCertRegistry.CallOpts, arg0)
}

// PublishedRoots is a free data retrieval call binding the contract method 0x744a92b6.
//
// Solidity: function publishedRoots(uint256 ) view returns(bytes32)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) PublishedRoots(arg0 *big.Int) ([32]byte, error) {
	return _IUMiCertRegistry.Contract.PublishedRoots(&_IUMiCertRegistry.CallOpts, arg0)
}

// PublishedTerms is a free data retrieval call binding the contract method 0x66ef8c1c.
//
// Solidity: function publishedTerms(uint256 ) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) PublishedTerms(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "publishedTerms", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PublishedTerms is a free data retrieval call binding the contract method 0x66ef8c1c.
//
// Solidity: function publishedTerms(uint256 ) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistrySession) PublishedTerms(arg0 *big.Int) (string, error) {
	return _IUMiCertRegistry.Contract.PublishedTerms(&_IUMiCertRegistry.CallOpts, arg0)
}

// PublishedTerms is a free data retrieval call binding the contract method 0x66ef8c1c.
//
// Solidity: function publishedTerms(uint256 ) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) PublishedTerms(arg0 *big.Int) (string, error) {
	return _IUMiCertRegistry.Contract.PublishedTerms(&_IUMiCertRegistry.CallOpts, arg0)
}

// RootToTerm is a free data retrieval call binding the contract method 0xa37b8da8.
//
// Solidity: function rootToTerm(bytes32 ) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) RootToTerm(opts *bind.CallOpts, arg0 [32]byte) (string, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "rootToTerm", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RootToTerm is a free data retrieval call binding the contract method 0xa37b8da8.
//
// Solidity: function rootToTerm(bytes32 ) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistrySession) RootToTerm(arg0 [32]byte) (string, error) {
	return _IUMiCertRegistry.Contract.RootToTerm(&_IUMiCertRegistry.CallOpts, arg0)
}

// RootToTerm is a free data retrieval call binding the contract method 0xa37b8da8.
//
// Solidity: function rootToTerm(bytes32 ) view returns(string)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) RootToTerm(arg0 [32]byte) (string, error) {
	return _IUMiCertRegistry.Contract.RootToTerm(&_IUMiCertRegistry.CallOpts, arg0)
}

// RootToVersion is a free data retrieval call binding the contract method 0xcbb1c002.
//
// Solidity: function rootToVersion(bytes32 ) view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) RootToVersion(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "rootToVersion", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RootToVersion is a free data retrieval call binding the contract method 0xcbb1c002.
//
// Solidity: function rootToVersion(bytes32 ) view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistrySession) RootToVersion(arg0 [32]byte) (*big.Int, error) {
	return _IUMiCertRegistry.Contract.RootToVersion(&_IUMiCertRegistry.CallOpts, arg0)
}

// RootToVersion is a free data retrieval call binding the contract method 0xcbb1c002.
//
// Solidity: function rootToVersion(bytes32 ) view returns(uint256)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) RootToVersion(arg0 [32]byte) (*big.Int, error) {
	return _IUMiCertRegistry.Contract.RootToVersion(&_IUMiCertRegistry.CallOpts, arg0)
}

// TermVersions is a free data retrieval call binding the contract method 0x9645b359.
//
// Solidity: function termVersions(string , uint256 ) view returns(bytes32 rootHash, uint256 version, uint256 totalStudents, uint256 publishedAt, bool isSuperseded, bytes32 supersededBy, string supersessionReason)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) TermVersions(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	RootHash           [32]byte
	Version            *big.Int
	TotalStudents      *big.Int
	PublishedAt        *big.Int
	IsSuperseded       bool
	SupersededBy       [32]byte
	SupersessionReason string
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "termVersions", arg0, arg1)

	outstruct := new(struct {
		RootHash           [32]byte
		Version            *big.Int
		TotalStudents      *big.Int
		PublishedAt        *big.Int
		IsSuperseded       bool
		SupersededBy       [32]byte
		SupersessionReason string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RootHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Version = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalStudents = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PublishedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.IsSuperseded = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.SupersededBy = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.SupersessionReason = *abi.ConvertType(out[6], new(string)).(*string)

	return *outstruct, err

}

// TermVersions is a free data retrieval call binding the contract method 0x9645b359.
//
// Solidity: function termVersions(string , uint256 ) view returns(bytes32 rootHash, uint256 version, uint256 totalStudents, uint256 publishedAt, bool isSuperseded, bytes32 supersededBy, string supersessionReason)
func (_IUMiCertRegistry *IUMiCertRegistrySession) TermVersions(arg0 string, arg1 *big.Int) (struct {
	RootHash           [32]byte
	Version            *big.Int
	TotalStudents      *big.Int
	PublishedAt        *big.Int
	IsSuperseded       bool
	SupersededBy       [32]byte
	SupersessionReason string
}, error) {
	return _IUMiCertRegistry.Contract.TermVersions(&_IUMiCertRegistry.CallOpts, arg0, arg1)
}

// TermVersions is a free data retrieval call binding the contract method 0x9645b359.
//
// Solidity: function termVersions(string , uint256 ) view returns(bytes32 rootHash, uint256 version, uint256 totalStudents, uint256 publishedAt, bool isSuperseded, bytes32 supersededBy, string supersessionReason)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) TermVersions(arg0 string, arg1 *big.Int) (struct {
	RootHash           [32]byte
	Version            *big.Int
	TotalStudents      *big.Int
	PublishedAt        *big.Int
	IsSuperseded       bool
	SupersededBy       [32]byte
	SupersessionReason string
}, error) {
	return _IUMiCertRegistry.Contract.TermVersions(&_IUMiCertRegistry.CallOpts, arg0, arg1)
}

// VerifyReceiptAnchor is a free data retrieval call binding the contract method 0x2a2dab4a.
//
// Solidity: function verifyReceiptAnchor(bytes32 _blockchainAnchor) view returns(bool isValid, string termId, uint256 publishedAt)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) VerifyReceiptAnchor(opts *bind.CallOpts, _blockchainAnchor [32]byte) (struct {
	IsValid     bool
	TermId      string
	PublishedAt *big.Int
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "verifyReceiptAnchor", _blockchainAnchor)

	outstruct := new(struct {
		IsValid     bool
		TermId      string
		PublishedAt *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsValid = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.TermId = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.PublishedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// VerifyReceiptAnchor is a free data retrieval call binding the contract method 0x2a2dab4a.
//
// Solidity: function verifyReceiptAnchor(bytes32 _blockchainAnchor) view returns(bool isValid, string termId, uint256 publishedAt)
func (_IUMiCertRegistry *IUMiCertRegistrySession) VerifyReceiptAnchor(_blockchainAnchor [32]byte) (struct {
	IsValid     bool
	TermId      string
	PublishedAt *big.Int
}, error) {
	return _IUMiCertRegistry.Contract.VerifyReceiptAnchor(&_IUMiCertRegistry.CallOpts, _blockchainAnchor)
}

// VerifyReceiptAnchor is a free data retrieval call binding the contract method 0x2a2dab4a.
//
// Solidity: function verifyReceiptAnchor(bytes32 _blockchainAnchor) view returns(bool isValid, string termId, uint256 publishedAt)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) VerifyReceiptAnchor(_blockchainAnchor [32]byte) (struct {
	IsValid     bool
	TermId      string
	PublishedAt *big.Int
}, error) {
	return _IUMiCertRegistry.Contract.VerifyReceiptAnchor(&_IUMiCertRegistry.CallOpts, _blockchainAnchor)
}

// PublishTermRoot is a paid mutator transaction binding the contract method 0xfe1e73fb.
//
// Solidity: function publishTermRoot(bytes32 _verkleRoot, string _termId, uint256 _totalStudents) returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactor) PublishTermRoot(opts *bind.TransactOpts, _verkleRoot [32]byte, _termId string, _totalStudents *big.Int) (*types.Transaction, error) {
	return _IUMiCertRegistry.contract.Transact(opts, "publishTermRoot", _verkleRoot, _termId, _totalStudents)
}

// PublishTermRoot is a paid mutator transaction binding the contract method 0xfe1e73fb.
//
// Solidity: function publishTermRoot(bytes32 _verkleRoot, string _termId, uint256 _totalStudents) returns()
func (_IUMiCertRegistry *IUMiCertRegistrySession) PublishTermRoot(_verkleRoot [32]byte, _termId string, _totalStudents *big.Int) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.PublishTermRoot(&_IUMiCertRegistry.TransactOpts, _verkleRoot, _termId, _totalStudents)
}

// PublishTermRoot is a paid mutator transaction binding the contract method 0xfe1e73fb.
//
// Solidity: function publishTermRoot(bytes32 _verkleRoot, string _termId, uint256 _totalStudents) returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactorSession) PublishTermRoot(_verkleRoot [32]byte, _termId string, _totalStudents *big.Int) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.PublishTermRoot(&_IUMiCertRegistry.TransactOpts, _verkleRoot, _termId, _totalStudents)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUMiCertRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IUMiCertRegistry *IUMiCertRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.RenounceOwnership(&_IUMiCertRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.RenounceOwnership(&_IUMiCertRegistry.TransactOpts)
}

// SupersedeTerm is a paid mutator transaction binding the contract method 0xf9b1d621.
//
// Solidity: function supersedeTerm(string _termId, bytes32 _newVerkleRoot, uint256 _newTotalStudents, string _reason) returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactor) SupersedeTerm(opts *bind.TransactOpts, _termId string, _newVerkleRoot [32]byte, _newTotalStudents *big.Int, _reason string) (*types.Transaction, error) {
	return _IUMiCertRegistry.contract.Transact(opts, "supersedeTerm", _termId, _newVerkleRoot, _newTotalStudents, _reason)
}

// SupersedeTerm is a paid mutator transaction binding the contract method 0xf9b1d621.
//
// Solidity: function supersedeTerm(string _termId, bytes32 _newVerkleRoot, uint256 _newTotalStudents, string _reason) returns()
func (_IUMiCertRegistry *IUMiCertRegistrySession) SupersedeTerm(_termId string, _newVerkleRoot [32]byte, _newTotalStudents *big.Int, _reason string) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.SupersedeTerm(&_IUMiCertRegistry.TransactOpts, _termId, _newVerkleRoot, _newTotalStudents, _reason)
}

// SupersedeTerm is a paid mutator transaction binding the contract method 0xf9b1d621.
//
// Solidity: function supersedeTerm(string _termId, bytes32 _newVerkleRoot, uint256 _newTotalStudents, string _reason) returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactorSession) SupersedeTerm(_termId string, _newVerkleRoot [32]byte, _newTotalStudents *big.Int, _reason string) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.SupersedeTerm(&_IUMiCertRegistry.TransactOpts, _termId, _newVerkleRoot, _newTotalStudents, _reason)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IUMiCertRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IUMiCertRegistry *IUMiCertRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.TransferOwnership(&_IUMiCertRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IUMiCertRegistry *IUMiCertRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.TransferOwnership(&_IUMiCertRegistry.TransactOpts, newOwner)
}

// IUMiCertRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IUMiCertRegistry contract.
type IUMiCertRegistryOwnershipTransferredIterator struct {
	Event *IUMiCertRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IUMiCertRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUMiCertRegistryOwnershipTransferred)
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
		it.Event = new(IUMiCertRegistryOwnershipTransferred)
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
func (it *IUMiCertRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUMiCertRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUMiCertRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the IUMiCertRegistry contract.
type IUMiCertRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IUMiCertRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryOwnershipTransferredIterator{contract: _IUMiCertRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IUMiCertRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUMiCertRegistryOwnershipTransferred)
				if err := _IUMiCertRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*IUMiCertRegistryOwnershipTransferred, error) {
	event := new(IUMiCertRegistryOwnershipTransferred)
	if err := _IUMiCertRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUMiCertRegistryTermRootPublishedIterator is returned from FilterTermRootPublished and is used to iterate over the raw logs and unpacked data for TermRootPublished events raised by the IUMiCertRegistry contract.
type IUMiCertRegistryTermRootPublishedIterator struct {
	Event *IUMiCertRegistryTermRootPublished // Event containing the contract specifics and raw log

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
func (it *IUMiCertRegistryTermRootPublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUMiCertRegistryTermRootPublished)
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
		it.Event = new(IUMiCertRegistryTermRootPublished)
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
func (it *IUMiCertRegistryTermRootPublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUMiCertRegistryTermRootPublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUMiCertRegistryTermRootPublished represents a TermRootPublished event raised by the IUMiCertRegistry contract.
type IUMiCertRegistryTermRootPublished struct {
	TermId        common.Hash
	VerkleRoot    [32]byte
	Version       *big.Int
	TotalStudents *big.Int
	Timestamp     *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTermRootPublished is a free log retrieval operation binding the contract event 0x62fc1dad9c7e5a08066a48e92398689a7b0d7fda9a4518c90ed24cba8b51191a.
//
// Solidity: event TermRootPublished(string indexed termId, bytes32 indexed verkleRoot, uint256 version, uint256 totalStudents, uint256 timestamp)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) FilterTermRootPublished(opts *bind.FilterOpts, termId []string, verkleRoot [][32]byte) (*IUMiCertRegistryTermRootPublishedIterator, error) {

	var termIdRule []interface{}
	for _, termIdItem := range termId {
		termIdRule = append(termIdRule, termIdItem)
	}
	var verkleRootRule []interface{}
	for _, verkleRootItem := range verkleRoot {
		verkleRootRule = append(verkleRootRule, verkleRootItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.FilterLogs(opts, "TermRootPublished", termIdRule, verkleRootRule)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryTermRootPublishedIterator{contract: _IUMiCertRegistry.contract, event: "TermRootPublished", logs: logs, sub: sub}, nil
}

// WatchTermRootPublished is a free log subscription operation binding the contract event 0x62fc1dad9c7e5a08066a48e92398689a7b0d7fda9a4518c90ed24cba8b51191a.
//
// Solidity: event TermRootPublished(string indexed termId, bytes32 indexed verkleRoot, uint256 version, uint256 totalStudents, uint256 timestamp)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) WatchTermRootPublished(opts *bind.WatchOpts, sink chan<- *IUMiCertRegistryTermRootPublished, termId []string, verkleRoot [][32]byte) (event.Subscription, error) {

	var termIdRule []interface{}
	for _, termIdItem := range termId {
		termIdRule = append(termIdRule, termIdItem)
	}
	var verkleRootRule []interface{}
	for _, verkleRootItem := range verkleRoot {
		verkleRootRule = append(verkleRootRule, verkleRootItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.WatchLogs(opts, "TermRootPublished", termIdRule, verkleRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUMiCertRegistryTermRootPublished)
				if err := _IUMiCertRegistry.contract.UnpackLog(event, "TermRootPublished", log); err != nil {
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

// ParseTermRootPublished is a log parse operation binding the contract event 0x62fc1dad9c7e5a08066a48e92398689a7b0d7fda9a4518c90ed24cba8b51191a.
//
// Solidity: event TermRootPublished(string indexed termId, bytes32 indexed verkleRoot, uint256 version, uint256 totalStudents, uint256 timestamp)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) ParseTermRootPublished(log types.Log) (*IUMiCertRegistryTermRootPublished, error) {
	event := new(IUMiCertRegistryTermRootPublished)
	if err := _IUMiCertRegistry.contract.UnpackLog(event, "TermRootPublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUMiCertRegistryTermRootSupersededIterator is returned from FilterTermRootSuperseded and is used to iterate over the raw logs and unpacked data for TermRootSuperseded events raised by the IUMiCertRegistry contract.
type IUMiCertRegistryTermRootSupersededIterator struct {
	Event *IUMiCertRegistryTermRootSuperseded // Event containing the contract specifics and raw log

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
func (it *IUMiCertRegistryTermRootSupersededIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUMiCertRegistryTermRootSuperseded)
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
		it.Event = new(IUMiCertRegistryTermRootSuperseded)
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
func (it *IUMiCertRegistryTermRootSupersededIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUMiCertRegistryTermRootSupersededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUMiCertRegistryTermRootSuperseded represents a TermRootSuperseded event raised by the IUMiCertRegistry contract.
type IUMiCertRegistryTermRootSuperseded struct {
	TermId     common.Hash
	OldVersion *big.Int
	OldRoot    [32]byte
	NewVersion *big.Int
	NewRoot    [32]byte
	Reason     string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTermRootSuperseded is a free log retrieval operation binding the contract event 0x4a97ba7e0786a95b88bc5a9576a42ba3b36dcdad508f3266d97002d7f65d5cf2.
//
// Solidity: event TermRootSuperseded(string indexed termId, uint256 oldVersion, bytes32 indexed oldRoot, uint256 newVersion, bytes32 indexed newRoot, string reason)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) FilterTermRootSuperseded(opts *bind.FilterOpts, termId []string, oldRoot [][32]byte, newRoot [][32]byte) (*IUMiCertRegistryTermRootSupersededIterator, error) {

	var termIdRule []interface{}
	for _, termIdItem := range termId {
		termIdRule = append(termIdRule, termIdItem)
	}

	var oldRootRule []interface{}
	for _, oldRootItem := range oldRoot {
		oldRootRule = append(oldRootRule, oldRootItem)
	}

	var newRootRule []interface{}
	for _, newRootItem := range newRoot {
		newRootRule = append(newRootRule, newRootItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.FilterLogs(opts, "TermRootSuperseded", termIdRule, oldRootRule, newRootRule)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryTermRootSupersededIterator{contract: _IUMiCertRegistry.contract, event: "TermRootSuperseded", logs: logs, sub: sub}, nil
}

// WatchTermRootSuperseded is a free log subscription operation binding the contract event 0x4a97ba7e0786a95b88bc5a9576a42ba3b36dcdad508f3266d97002d7f65d5cf2.
//
// Solidity: event TermRootSuperseded(string indexed termId, uint256 oldVersion, bytes32 indexed oldRoot, uint256 newVersion, bytes32 indexed newRoot, string reason)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) WatchTermRootSuperseded(opts *bind.WatchOpts, sink chan<- *IUMiCertRegistryTermRootSuperseded, termId []string, oldRoot [][32]byte, newRoot [][32]byte) (event.Subscription, error) {

	var termIdRule []interface{}
	for _, termIdItem := range termId {
		termIdRule = append(termIdRule, termIdItem)
	}

	var oldRootRule []interface{}
	for _, oldRootItem := range oldRoot {
		oldRootRule = append(oldRootRule, oldRootItem)
	}

	var newRootRule []interface{}
	for _, newRootItem := range newRoot {
		newRootRule = append(newRootRule, newRootItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.WatchLogs(opts, "TermRootSuperseded", termIdRule, oldRootRule, newRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUMiCertRegistryTermRootSuperseded)
				if err := _IUMiCertRegistry.contract.UnpackLog(event, "TermRootSuperseded", log); err != nil {
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

// ParseTermRootSuperseded is a log parse operation binding the contract event 0x4a97ba7e0786a95b88bc5a9576a42ba3b36dcdad508f3266d97002d7f65d5cf2.
//
// Solidity: event TermRootSuperseded(string indexed termId, uint256 oldVersion, bytes32 indexed oldRoot, uint256 newVersion, bytes32 indexed newRoot, string reason)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) ParseTermRootSuperseded(log types.Log) (*IUMiCertRegistryTermRootSuperseded, error) {
	event := new(IUMiCertRegistryTermRootSuperseded)
	if err := _IUMiCertRegistry.contract.UnpackLog(event, "TermRootSuperseded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
