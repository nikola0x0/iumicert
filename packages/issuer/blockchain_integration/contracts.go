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
)

// IUMiCertRegistryMetaData contains all meta data concerning the IUMiCertRegistry contract.
var IUMiCertRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getPublishedRoot\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublishedRootsCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTermRootInfo\",\"inputs\":[{\"name\":\"_verkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"publishTermRoot\",\"inputs\":[{\"name\":\"_verkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"publishedRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"termRoots\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyReceiptAnchor\",\"inputs\":[{\"name\":\"_blockchainAnchor\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"termId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"publishedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TermRootPublished\",\"inputs\":[{\"name\":\"verkleRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"termId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"totalStudents\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
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
	parsed, err := abi.JSON(strings.NewReader(IUMiCertRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUMiCertRegistry *IUMiCertRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUMiCertRegistry.Contract.IUMiCertRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// the fallback function (if any).
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

// Transact invokes the (paid) contract method with params as input values.
func (_IUMiCertRegistry *IUMiCertRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUMiCertRegistry.Contract.contract.Transact(opts, method, params...)
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

// GetTermRootInfo is a free data retrieval call binding the contract method 0x6f90e9ea.
//
// Solidity: function getTermRootInfo(bytes32 _verkleRoot) view returns(string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) GetTermRootInfo(opts *bind.CallOpts, _verkleRoot [32]byte) (struct {
	TermId        string
	TotalStudents *big.Int
	PublishedAt   *big.Int
	Exists        bool
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "getTermRootInfo", _verkleRoot)

	outstruct := new(struct {
		TermId        string
		TotalStudents *big.Int
		PublishedAt   *big.Int
		Exists        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TermId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.TotalStudents = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PublishedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// GetTermRootInfo is a free data retrieval call binding the contract method 0x6f90e9ea.
//
// Solidity: function getTermRootInfo(bytes32 _verkleRoot) view returns(string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
func (_IUMiCertRegistry *IUMiCertRegistrySession) GetTermRootInfo(_verkleRoot [32]byte) (struct {
	TermId        string
	TotalStudents *big.Int
	PublishedAt   *big.Int
	Exists        bool
}, error) {
	return _IUMiCertRegistry.Contract.GetTermRootInfo(&_IUMiCertRegistry.CallOpts, _verkleRoot)
}

// GetTermRootInfo is a free data retrieval call binding the contract method 0x6f90e9ea.
//
// Solidity: function getTermRootInfo(bytes32 _verkleRoot) view returns(string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) GetTermRootInfo(_verkleRoot [32]byte) (struct {
	TermId        string
	TotalStudents *big.Int
	PublishedAt   *big.Int
	Exists        bool
}, error) {
	return _IUMiCertRegistry.Contract.GetTermRootInfo(&_IUMiCertRegistry.CallOpts, _verkleRoot)
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

// TermRoots is a free data retrieval call binding the contract method 0x50fe2c36.
//
// Solidity: function termRoots(bytes32 ) view returns(string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
func (_IUMiCertRegistry *IUMiCertRegistryCaller) TermRoots(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TermId        string
	TotalStudents *big.Int
	PublishedAt   *big.Int
	Exists        bool
}, error) {
	var out []interface{}
	err := _IUMiCertRegistry.contract.Call(opts, &out, "termRoots", arg0)

	outstruct := new(struct {
		TermId        string
		TotalStudents *big.Int
		PublishedAt   *big.Int
		Exists        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TermId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.TotalStudents = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PublishedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// TermRoots is a free data retrieval call binding the contract method 0x50fe2c36.
//
// Solidity: function termRoots(bytes32 ) view returns(string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
func (_IUMiCertRegistry *IUMiCertRegistrySession) TermRoots(arg0 [32]byte) (struct {
	TermId        string
	TotalStudents *big.Int
	PublishedAt   *big.Int
	Exists        bool
}, error) {
	return _IUMiCertRegistry.Contract.TermRoots(&_IUMiCertRegistry.CallOpts, arg0)
}

// TermRoots is a free data retrieval call binding the contract method 0x50fe2c36.
//
// Solidity: function termRoots(bytes32 ) view returns(string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
func (_IUMiCertRegistry *IUMiCertRegistryCallerSession) TermRoots(arg0 [32]byte) (struct {
	TermId        string
	TotalStudents *big.Int
	PublishedAt   *big.Int
	Exists        bool
}, error) {
	return _IUMiCertRegistry.Contract.TermRoots(&_IUMiCertRegistry.CallOpts, arg0)
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
	VerkleRoot    [32]byte
	TermId        common.Hash
	TotalStudents *big.Int
	Timestamp     *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTermRootPublished is a free log retrieval operation binding the contract event 0xea6b942a9a27c613b21181fd056514044b334e9a6930a7e0be909ab4b520ee4b.
//
// Solidity: event TermRootPublished(bytes32 indexed verkleRoot, string indexed termId, uint256 totalStudents, uint256 timestamp)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) FilterTermRootPublished(opts *bind.FilterOpts, verkleRoot [][32]byte, termId []string) (*IUMiCertRegistryTermRootPublishedIterator, error) {

	var verkleRootRule []interface{}
	for _, verkleRootItem := range verkleRoot {
		verkleRootRule = append(verkleRootRule, verkleRootItem)
	}
	var termIdRule []interface{}
	for _, termIdItem := range termId {
		termIdRule = append(termIdRule, termIdItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.FilterLogs(opts, "TermRootPublished", verkleRootRule, termIdRule)
	if err != nil {
		return nil, err
	}
	return &IUMiCertRegistryTermRootPublishedIterator{contract: _IUMiCertRegistry.contract, event: "TermRootPublished", logs: logs, sub: sub}, nil
}

// WatchTermRootPublished is a free log subscription operation binding the contract event 0xea6b942a9a27c613b21181fd056514044b334e9a6930a7e0be909ab4b520ee4b.
//
// Solidity: event TermRootPublished(bytes32 indexed verkleRoot, string indexed termId, uint256 totalStudents, uint256 timestamp)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) WatchTermRootPublished(opts *bind.WatchOpts, sink chan<- *IUMiCertRegistryTermRootPublished, verkleRoot [][32]byte, termId []string) (event.Subscription, error) {

	var verkleRootRule []interface{}
	for _, verkleRootItem := range verkleRoot {
		verkleRootRule = append(verkleRootRule, verkleRootItem)
	}
	var termIdRule []interface{}
	for _, termIdItem := range termId {
		termIdRule = append(termIdRule, termIdItem)
	}

	logs, sub, err := _IUMiCertRegistry.contract.WatchLogs(opts, "TermRootPublished", verkleRootRule, termIdRule)
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

// ParseTermRootPublished is a log parse operation binding the contract event 0xea6b942a9a27c613b21181fd056514044b334e9a6930a7e0be909ab4b520ee4b.
//
// Solidity: event TermRootPublished(bytes32 indexed verkleRoot, string indexed termId, uint256 totalStudents, uint256 timestamp)
func (_IUMiCertRegistry *IUMiCertRegistryFilterer) ParseTermRootPublished(log types.Log) (*IUMiCertRegistryTermRootPublished, error) {
	event := new(IUMiCertRegistryTermRootPublished)
	if err := _IUMiCertRegistry.contract.UnpackLog(event, "TermRootPublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployIUMiCertRegistry deploys a new IUMiCertRegistry contract to the blockchain.
func DeployIUMiCertRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, initialOwner common.Address) (common.Address, *types.Transaction, *IUMiCertRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(IUMiCertRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex("608060405234801561000f575f5ffd5b506040516115fc3803806115fc833981810160405281019061003191906101d7565b805f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100999190610211565b60405180910390fd5b6100b1816100b860201b60201c565b505061022a565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101a68261017d565b9050919050565b6101b68161019c565b81146101c0575f5ffd5b50565b5f815190506101d1816101ad565b92915050565b5f602082840312156101ec576101eb610179565b5b5f6101f9848285016101c3565b91505092915050565b61020b8161019c565b82525050565b5f6020820190506102245f830184610202565b92915050565b6113c5806102375f395ff3fe608060405234801561000f575f5ffd5b506004361061009c575f3560e01c8063744a92b611610064578063744a92b6146101725780638da5cb5b146101a2578063a2fbbc58146101c0578063f2fde38b146101de578063fe1e73fb146101fa5761009c565b80632a2dab4a146100a0578063351863fb146100d257806350fe2c36146101025780636f90e9ea14610135578063715018a614610168575b5f5ffd5b6100ba60048036038101906100b591906109e9565b610216565b6040516100c993929190610ab6565b60405180910390f35b6100ec60048036038101906100e79190610b1c565b610312565b6040516100f99190610b56565b60405180910390f35b61011c600480360381019061011791906109e9565b61037e565b60405161012c9493929190610b6f565b60405180910390f35b61014f600480360381019061014a91906109e9565b61043c565b60405161015f9493929190610b6f565b60405180910390f35b610170610540565b005b61018c60048036038101906101879190610b1c565b610553565b6040516101999190610b56565b60405180910390f35b6101aa610573565b6040516101b79190610bf8565b60405180910390f35b6101c861059a565b6040516101d59190610c11565b60405180910390f35b6101f860048036038101906101f39190610c54565b6105a6565b005b610214600480360381019061020f9190610dab565b61062a565b005b5f60605f5f60015f8681526020019081526020015f206040518060800160405290815f8201805461024690610e44565b80601f016020809104026020016040519081016040528092919081815260200182805461027290610e44565b80156102bd5780601f10610294576101008083540402835291602001916102bd565b820191905f5260205f20905b8154815290600101906020018083116102a057829003601f168201915b505050505081526020016001820154815260200160028201548152602001600382015f9054906101000a900460ff16151515158152505090508060600151815f01518260400151935093509350509193909250565b5f600280549050821061035a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035190610ebe565b60405180910390fd5b6002828154811061036e5761036d610edc565b5b905f5260205f2001549050919050565b6001602052805f5260405f205f91509050805f01805461039d90610e44565b80601f01602080910402602001604051908101604052809291908181526020018280546103c990610e44565b80156104145780601f106103eb57610100808354040283529160200191610414565b820191905f5260205f20905b8154815290600101906020018083116103f757829003601f168201915b505050505090806001015490806002015490806003015f9054906101000a900460ff16905084565b60605f5f5f5f60015f8781526020019081526020015f206040518060800160405290815f8201805461046d90610e44565b80601f016020809104026020016040519081016040528092919081815260200182805461049990610e44565b80156104e45780601f106104bb576101008083540402835291602001916104e4565b820191905f5260205f20905b8154815290600101906020018083116104c757829003601f168201915b505050505081526020016001820154815260200160028201548152602001600382015f9054906101000a900460ff1615151515815250509050805f01518160200151826040015183606001519450945094509450509193509193565b610548610856565b6105515f6108dd565b565b60028181548110610562575f80fd5b905f5260205f20015f915090505481565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f600280549050905090565b6105ae610856565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361061e575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016106159190610bf8565b60405180910390fd5b610627816108dd565b50565b610632610856565b5f5f1b8303610676576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161066d90610f53565b60405180910390fd5b5f8251116106b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106b090610fbb565b60405180910390fd5b5f81116106fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106f290611023565b60405180910390fd5b60015f8481526020019081526020015f206003015f9054906101000a900460ff161561075c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107539061108b565b60405180910390fd5b60405180608001604052808381526020018281526020014281526020016001151581525060015f8581526020019081526020015f205f820151815f0190816107a49190611249565b5060208201518160010155604082015181600201556060820151816003015f6101000a81548160ff021916908315150217905550905050600283908060018154018082558091505060019003905f5260205f20015f90919091909150558160405161080f9190611352565b6040518091039020837fea6b942a9a27c613b21181fd056514044b334e9a6930a7e0be909ab4b520ee4b8342604051610849929190611368565b60405180910390a3505050565b61085e61099e565b73ffffffffffffffffffffffffffffffffffffffff1661087c610573565b73ffffffffffffffffffffffffffffffffffffffff16146108db5761089f61099e565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016108d29190610bf8565b60405180910390fd565b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b5f604051905090565b5f5ffd5b5f5ffd5b5f819050919050565b6109c8816109b6565b81146109d2575f5ffd5b50565b5f813590506109e3816109bf565b92915050565b5f602082840312156109fe576109fd6109ae565b5b5f610a0b848285016109d5565b91505092915050565b5f5ffd5b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610cbd82610a56565b810181811067ffffffffffffffff82111715610cdc57610cdb610c87565b5b80604052505050565b5f610cee6109a5565b9050610cfa8282610cb4565b919050565b5f67ffffffffffffffff821115610d1957610d18610c87565b5b610d2282610a56565b9050602081019050919050565b828183375f83830152505050565b5f610d4f610d4a84610cff565b610ce5565b905082815260208101848484011115610d6b57610d6a610c83565b5b610d76848285610d2f565b509392505050565b5f82601f830112610d9257610d91610c7f565b5b8135610da2848260208601610d3d565b91505092915050565b5f5f5f60608486031215610dc257610dc16109ae565b5b5f610dcf868287016109d5565b935050602084013567ffffffffffffffff811115610df057610def6109b2565b5b610dfc86828701610d7e565b9250506040610e0d86828701610b08565b9150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610e5b57607f821691505b602082108103610e6e57610e6d610e17565b5b50919050565b7f496e646578206f7574206f6620626f756e6473000000000000000000000000005f82015250565b5f610ea8601383610a38565b9150610eb382610e74565b602082019050919050565b5f6020820190508181035f830152610ed581610e9c565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f496e76616c6964205665726b6c6520726f6f74000000000000000000000000005f82015250565b5f610f3d601383610a38565b9150610f4882610f09565b602082019050919050565b5f6020820190508181035f830152610f6a81610f31565b9050919050565b7f5465726d204944207265717569726564000000000000000000000000000000005f82015250565b5f610fa5601083610a38565b9150610fb082610f71565b602082019050919050565b5f6020820190508181035f830152610fd281610f99565b9050919050565b7f496e76616c69642073747564656e7420636f756e7400000000000000000000005f82015250565b5f61100d601583610a38565b915061101882610fd9565b602082019050919050565b5f6020820190508181035f83015261103a81611001565b9050919050565b7f526f6f7420616c7265616479207075626c6973686564000000000000000000005f82015250565b5f611075601683610a38565b915061108082611041565b602082019050919050565b5f6020820190508181035f8301526110a281611069565b9050919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026111057fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826110ca565b61110f86836110ca565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61114a61114561114084610a9e565b611127565b610a9e565b9050919050565b5f819050919050565b61116383611130565b61117761116f82611151565b8484546110d6565b825550505050565b5f5f905090565b61118e61117f565b61119981848461115a565b505050565b5b818110156111bc576111b15f82611186565b60018101905061119f565b5050565b601f821115611201576111d2816110a9565b6111db846110bb565b810160208510156111ea578190505b6111fe6111f6856110bb565b83018261119e565b50505b505050565b5f82821c905092915050565b5f6112215f1984600802611206565b1980831691505092915050565b5f6112398383611212565b9150826002028217905092915050565b61125282610a2e565b67ffffffffffffffff81111561126b5761126a610c87565b5b6112758254610e44565b6112808282856111c0565b5f60209050601f8311600181146112b1575f841561129f578287015190505b6112a9858261122e565b865550611310565b601f1984166112bf866110a9565b5f5b828110156112e6578489015182556001820191506020850194506020810190506112c1565b8683101561130357848901516112ff601f891682611212565b8355505b6001600288020188555050505b505050505050565b5f81905092915050565b5f61132c82610a2e565b6113368185611318565b9350611346818560208601610a48565b80840191505092915050565b5f61135d8284611322565b915081905092915050565b5f60408201905061137b5f830185610aa7565b6113886020830184610aa7565b939250505056fea2646970667358221220e277a5bb80b2c97fb9417eddcbd259ba65c7fd5618839aa005dd0c7ab387248d64736f6c634300081e0033"), backend, initialOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	return address, tx, &IUMiCertRegistry{IUMiCertRegistryCaller: IUMiCertRegistryCaller{contract: bind.NewBoundContract(address, parsed, backend, backend, backend)}, IUMiCertRegistryTransactor: IUMiCertRegistryTransactor{contract: bind.NewBoundContract(address, parsed, backend, backend, backend)}, IUMiCertRegistryFilterer: IUMiCertRegistryFilterer{contract: bind.NewBoundContract(address, parsed, backend, backend, backend)}}, nil
}