// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coindistribution

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

// CoindistributionMetaData contains all meta data concerning the Coindistribution contract.
var CoindistributionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20ZeroToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeIncreased\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeMoreThan25\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeMoreThan5\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForeignTokenSelfTransfer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalanceForOnTopFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Mismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoBots\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSwapBack\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAirDropper\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawStuckETH\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BuyBackTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sniper\",\"type\":\"address\"}],\"name\":\"CaughtEarlyBuyer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EnableTrading\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isExcluded\",\"type\":\"bool\"}],\"name\":\"ExcludeFromFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"MaxTransactionExclusion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"theAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"OnSetUniswapRouter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"OwnerForcedSwapBack\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RemovedLimits\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sLiqF\",\"type\":\"uint256\"}],\"name\":\"SetUniswapV2LiquidityPool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensSwapped\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethReceived\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensIntoLiquidity\",\"type\":\"uint256\"}],\"name\":\"SwapAndLiquify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferForeignToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxBuyAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxSellAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxWalletAmount\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"airdropToWallets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"applySlippage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"botsCaught\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableTransferDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableTrading\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"excludeFromFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceSwapBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAirDropper\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTradingEnabledBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"limitsInEffect\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"manageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"wallets\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"massManageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"removeLimits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAirDropper\",\"type\":\"address\"}],\"name\":\"setAirDropper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"setKnownUniswapRouters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setSlippage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setSwapBackThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"theAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"setUniswapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sLiqF\",\"type\":\"uint256\"}],\"name\":\"setUniswapV2LiquidityPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokensForLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transferDelayEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferForeignToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"uniswapV2LiquidityPoolSlots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sLiqF\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uniswapV2LiquidityPools\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawStuckETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052683635c9adc5dea00000600b556003600c555f600d5f6101000a81548160ff0219169083151502179055506001600d60016101000a81548160ff0219169083151502179055506001600f5f6101000a81548160ff0219169083151502179055505f60145534801562000074575f80fd5b50336040518060400160405280600381526020017f49636500000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f49434500000000000000000000000000000000000000000000000000000000008152508160029081620000f3919062000912565b50806003908162000105919062000912565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036200017b575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040162000172919062000a39565b60405180910390fd5b6200018c816200020b60201b60201c565b503360055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001e0336001620002ce60201b60201c565b620001f3306001620002ce60201b60201c565b6200020560016200038660201b60201c565b62000a8b565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b620002de620004c960201b60201c565b8060135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df7826040516200037a919062000a70565b60405180910390a25050565b62000396620004c960201b60201c565b620003bc737a250d5630b4cf539739df2c5dacb4c659f2488d826200056b60201b60201c565b620003e273e592427a0aece92de3edee1f18e0157c05861564826200056b60201b60201c565b620004087368b3465833fb72a70ecdf485e0e4c7bd8665fc45826200056b60201b60201c565b6200042e73eff92a263d31888d860bd50809a8d171709b7b1c826200056b60201b60201c565b62000454731b81d678ffb9c0263b24a97847620c99d213eb14826200056b60201b60201c565b6200047a7313f4ea83d0bd40e75c8222255bc855a974568dd4826200056b60201b60201c565b620004a0731b02da8cb0d097eb8d57a175b88c7d8b47997506826200056b60201b60201c565b620004c673d9e1ce17f2641f24ae83637ab66a2cca9c378b9f826200056b60201b60201c565b50565b620004d96200067f60201b60201c565b73ffffffffffffffffffffffffffffffffffffffff16620004ff6200068660201b60201c565b73ffffffffffffffffffffffffffffffffffffffff161462000569576200052b6200067f60201b60201c565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040162000560919062000a39565b60405180910390fd5b565b6200057b620004c960201b60201c565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620005e1576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060155f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508015158273ffffffffffffffffffffffffffffffffffffffff167f3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c60405160405180910390a35050565b5f33905090565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200072a57607f821691505b60208210810362000740576200073f620006e5565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620007a47fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000767565b620007b0868362000767565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620007fa620007f4620007ee84620007c8565b620007d1565b620007c8565b9050919050565b5f819050919050565b6200081583620007da565b6200082d620008248262000801565b84845462000773565b825550505050565b5f90565b6200084362000835565b620008508184846200080a565b505050565b5b8181101562000877576200086b5f8262000839565b60018101905062000856565b5050565b601f821115620008c657620008908162000746565b6200089b8462000758565b81016020851015620008ab578190505b620008c3620008ba8562000758565b83018262000855565b50505b505050565b5f82821c905092915050565b5f620008e85f1984600802620008cb565b1980831691505092915050565b5f620009028383620008d7565b9150826002028217905092915050565b6200091d82620006ae565b67ffffffffffffffff811115620009395762000938620006b8565b5b62000945825462000712565b620009528282856200087b565b5f60209050601f83116001811462000988575f841562000973578287015190505b6200097f8582620008f5565b865550620009ee565b601f198416620009988662000746565b5f5b82811015620009c1578489015182556001820191506020850194506020810190506200099a565b86831015620009e15784890151620009dd601f891682620008d7565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f62000a2182620009f6565b9050919050565b62000a338162000a15565b82525050565b5f60208201905062000a4e5f83018462000a28565b92915050565b5f8115159050919050565b62000a6a8162000a54565b82525050565b5f60208201905062000a855f83018462000a5f565b92915050565b6146748062000a995f395ff3fe608060405260043610610233575f3560e01c80636ddd17131161012d578063a9059cbb116100aa578063e800dff71161006e578063e800dff714610815578063e884f2601461083d578063f0fa55a914610853578063f2fde38b1461087b578063f5648a4f146108a35761023a565b8063a9059cbb1461070f578063bfd792841461074b578063c024666814610787578063c876d0b9146107af578063dd62ed3e146107d95761023a565b80638366e79a116100f15780638366e79a146106555780638a8c523c1461067d5780638da5cb5b146106935780638fe62b8a146106bd57806395d89b41146106e55761023a565b80636ddd17131461058757806370a08231146105b1578063715018a6146105ed578063751039fc14610603578063785cca3e146106195761023a565b8063296a803c116101bb57806351f205e41161017f57806351f205e4146104cb578063588b655f146104e15780635db7548b1461050b5780635de998ef146105355780636b0a894c1461055d5761023a565b8063296a803c146103e6578063313ce567146104275780633d6d11b01461045157806344d8a785146104795780634a62bb65146104a15761023a565b806318160ddd1161020257806318160ddd146102f25780631a8145bb1461031c5780632307b4411461034657806323b872dd1461036e57806327ac2d68146103aa5761023a565b806306fdde031461023c578063072280c314610266578063095ea7b31461028e578063130a2c3c146102ca5761023a565b3661023a57005b005b348015610247575f80fd5b506102506108b9565b60405161025d9190613937565b60405180910390f35b348015610271575f80fd5b5061028c600480360381019061028791906139ee565b610949565b005b348015610299575f80fd5b506102b460048036038101906102af9190613a5f565b610a54565b6040516102c19190613aac565b60405180910390f35b3480156102d5575f80fd5b506102f060048036038101906102eb9190613b26565b610a76565b005b3480156102fd575f80fd5b50610306610b1f565b6040516103139190613b92565b60405180910390f35b348015610327575f80fd5b50610330610b28565b60405161033d9190613b92565b60405180910390f35b348015610351575f80fd5b5061036c60048036038101906103679190613c00565b610b2e565b005b348015610379575f80fd5b50610394600480360381019061038f9190613c7e565b610c60565b6040516103a19190613aac565b60405180910390f35b3480156103b5575f80fd5b506103d060048036038101906103cb9190613cce565b610c8e565b6040516103dd9190613b92565b60405180910390f35b3480156103f1575f80fd5b5061040c60048036038101906104079190613cf9565b610cbc565b60405161041e96959493929190613d33565b60405180910390f35b348015610432575f80fd5b5061043b610d4a565b6040516104489190613dad565b60405180910390f35b34801561045c575f80fd5b5061047760048036038101906104729190613dc6565b610d52565b005b348015610484575f80fd5b5061049f600480360381019061049a9190613e4f565b6112dd565b005b3480156104ac575f80fd5b506104b56113d8565b6040516104c29190613aac565b60405180910390f35b3480156104d6575f80fd5b506104df6113eb565b005b3480156104ec575f80fd5b506104f56114a8565b6040516105029190613b92565b60405180910390f35b348015610516575f80fd5b5061051f6114b1565b60405161052c9190613e7a565b60405180910390f35b348015610540575f80fd5b5061055b60048036038101906105569190613cf9565b6114d9565b005b348015610568575f80fd5b50610571611524565b60405161057e9190613b92565b60405180910390f35b348015610592575f80fd5b5061059b61152a565b6040516105a89190613aac565b60405180910390f35b3480156105bc575f80fd5b506105d760048036038101906105d29190613cf9565b61153c565b6040516105e49190613b92565b60405180910390f35b3480156105f8575f80fd5b50610601611582565b005b34801561060e575f80fd5b506106176115d4565b005b348015610624575f80fd5b5061063f600480360381019061063a9190613cce565b61163d565b60405161064c9190613e7a565b60405180910390f35b348015610660575f80fd5b5061067b60048036038101906106769190613e93565b611678565b005b348015610688575f80fd5b50610691611810565b005b34801561069e575f80fd5b506106a76118a2565b6040516106b49190613e7a565b60405180910390f35b3480156106c8575f80fd5b506106e360048036038101906106de9190613cce565b6118ca565b005b3480156106f0575f80fd5b506106f96118dc565b6040516107069190613937565b60405180910390f35b34801561071a575f80fd5b5061073560048036038101906107309190613a5f565b61196c565b6040516107429190613aac565b60405180910390f35b348015610756575f80fd5b50610771600480360381019061076c9190613cf9565b61198e565b60405161077e9190613aac565b60405180910390f35b348015610792575f80fd5b506107ad60048036038101906107a891906139ee565b6119ab565b005b3480156107ba575f80fd5b506107c3611a59565b6040516107d09190613aac565b60405180910390f35b3480156107e4575f80fd5b506107ff60048036038101906107fa9190613e93565b611a6b565b60405161080c9190613b92565b60405180910390f35b348015610820575f80fd5b5061083b600480360381019061083691906139ee565b611aec565b005b348015610848575f80fd5b50610851611b4c565b005b34801561085e575f80fd5b5061087960048036038101906108749190613cce565b611b6f565b005b348015610886575f80fd5b506108a1600480360381019061089c9190613cf9565b611b81565b005b3480156108ae575f80fd5b506108b7611c05565b005b6060600280546108c890613efe565b80601f01602080910402602001604051908101604052809291908181526020018280546108f490613efe565b801561093f5780601f106109165761010080835404028352916020019161093f565b820191905f5260205f20905b81548152906001019060200180831161092257829003601f168201915b5050505050905090565b610951611cef565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036109b6576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060155f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508015158273ffffffffffffffffffffffffffffffffffffffff167f3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c60405160405180910390a35050565b5f80610a5e611d76565b9050610a6b818585611d7d565b600191505092915050565b610a7e611cef565b5f5b83839050811015610b19578160085f868685818110610aa257610aa1613f2e565b5b9050602002016020810190610ab79190613cf9565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508080610b1190613f88565b915050610a80565b50505050565b5f600754905090565b60105481565b3373ffffffffffffffffffffffffffffffffffffffff1660055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610bb4576040517f16ad4feb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818190508484905014610bf3576040517f77a93d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8484905090505f604051602081016004356024355f5b86811015610c4457602460208202018083013581850135875260068652604087208181540181558189019850505050600181019050610c0a565b50505050508060075f8282540192505081905550505050505050565b5f80610c6a611d76565b9050610c77858285611d8f565b610c82858585611e21565b60019150509392505050565b5f6064600c546064610ca09190613fcf565b83610cab9190614002565b610cb59190614070565b9050919050565b6011602052805f5260405f205f91509050805f015f9054906101000a900460ff1690805f0160019054906101000a900460ff1690805f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905086565b5f6012905090565b610d5a611cef565b6019811115610d95576040517f3601c4ac00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005821115610dd0576040517fbc32fbf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161480610e3557505f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16145b15610e6c576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f805b601280549050811015610f06578673ffffffffffffffffffffffffffffffffffffffff1660128281548110610ea757610ea6613f2e565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610ef35760019150610f06565b8080610efe90613f88565b915050610e6f565b5080610f7157601286908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611036565b60115f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2060020154831180610ffe575060115f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206003015482115b15611035576040517f5ab576c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5b8560115f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508460115f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f6101000a81548160ff0219169083151502179055508360115f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160016101000a81548160ff0219169083151502179055508660115f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508260115f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600201819055508160115f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600301819055508573ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff167f78c6052d21b64dbee719270b6a9e56ba166f8f88fbaa39a44abf4a11642bae92878787876040516112cc94939291906140a0565b60405180910390a350505050505050565b6112e5611cef565b611303737a250d5630b4cf539739df2c5dacb4c659f2488d82610949565b61132173e592427a0aece92de3edee1f18e0157c0586156482610949565b61133f7368b3465833fb72a70ecdf485e0e4c7bd8665fc4582610949565b61135d73eff92a263d31888d860bd50809a8d171709b7b1c82610949565b61137b731b81d678ffb9c0263b24a97847620c99d213eb1482610949565b6113997313f4ea83d0bd40e75c8222255bc855a974568dd482610949565b6113b7731b02da8cb0d097eb8d57a175b88c7d8b4799750682610949565b6113d573d9e1ce17f2641f24ae83637ab66a2cca9c378b9f82610949565b50565b600d60019054906101000a900460ff1681565b6113f3611cef565b6113fc3061153c565b5f03611434576040517f0b952c3600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600a5f6101000a81548160ff021916908315150217905550611456612963565b5f600a5f6101000a81548160ff0219169083151502179055507f1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb324260405161149e9190613b92565b60405180910390a1565b5f601454905090565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6114e1611cef565b8060055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60095481565b600d5f9054906101000a900460ff1681565b5f60065f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b61158a611cef565b6115926129ec565b5f60055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6115dc611cef565b5f600d60016101000a81548160ff0219169083151502179055505f600f5f6101000a81548160ff0219169083151502179055507fa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c60405160405180910390a1565b6012818154811061164c575f80fd5b905f5260205f20015f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b611680611cef565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036116e5576040517fdad1a1b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361174a576040517f74fc211300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8273ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016117849190613e7a565b602060405180830381865afa15801561179f573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906117c391906140f7565b90506117d08383836129ff565b7f5661684995ab94d684bfe57a43c4141578f52d3e7374e8cd3250e2f062e13ac183838360405161180393929190614122565b60405180910390a1505050565b611818611cef565b6014545f14611853576040517fd723eaba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600d5f6101000a81548160ff021916908315150217905550436014819055507f1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d60405160405180910390a1565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6118d2611cef565b80600b8190555050565b6060600380546118eb90613efe565b80601f016020809104026020016040519081016040528092919081815260200182805461191790613efe565b80156119625780601f1061193957610100808354040283529160200191611962565b820191905f5260205f20905b81548152906001019060200180831161194557829003601f168201915b5050505050905090565b5f80611976611d76565b9050611983818585611e21565b600191505092915050565b6008602052805f5260405f205f915054906101000a900460ff1681565b6119b3611cef565b8060135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df782604051611a4d9190613aac565b60405180910390a25050565b600f5f9054906101000a900460ff1681565b5f805f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b611af4611cef565b8060085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055505050565b611b54611cef565b5f600f5f6101000a81548160ff021916908315150217905550565b611b77611cef565b80600c8190555050565b611b89611cef565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611bf9575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401611bf09190613e7a565b60405180910390fd5b611c0281612a7e565b50565b611c0d611cef565b5f479050805f03611c4a576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f3373ffffffffffffffffffffffffffffffffffffffff1682604051611c6f90614184565b5f6040518083038185875af1925050503d805f8114611ca9576040519150601f19603f3d011682016040523d82523d5f602084013e611cae565b606091505b50508091505080611ceb576040517f3132169500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b611cf7611d76565b73ffffffffffffffffffffffffffffffffffffffff16611d156118a2565b73ffffffffffffffffffffffffffffffffffffffff1614611d7457611d38611d76565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401611d6b9190613e7a565b60405180910390fd5b565b5f33905090565b611d8a8383836001612b41565b505050565b5f611d9a8484611a6b565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611e1b5781811015611e0c578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401611e0393929190614198565b60405180910390fd5b611e1a84848484035f612b41565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611e9157816040517fec442f05000000000000000000000000000000000000000000000000000000008152600401611e889190613e7a565b60405180910390fd5b5f8111611ed3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611eca90614217565b60405180910390fd5b611edb612d0f565b158015611f2e575060155f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b1561200e5760135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1680611fce575060135f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b61200d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016120049061427f565b60405180910390fd5b5b60085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16806120a9575060085f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b156120df576040517e61c20e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600d60019054906101000a900460ff1615612449576120fc6118a2565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561216a575061213a6118a2565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156121a257505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156121f5575060135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b8015612248575060135f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b1561244857600f5f9054906101000a900460ff1615612447575f61226a612d1b565b9050806040015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141580156122dc5750806060015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614155b15612445576002436122ee9190613fcf565b600e5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410801561238157506002436123419190613fcf565b600e5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054105b6123c0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016123b79061430d565b60405180910390fd5b43600e5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555043600e5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055505b505b5b5b5f805b601280549050811015612542573373ffffffffffffffffffffffffffffffffffffffff1660115f6012848154811061248757612486613f2e565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361252f5760019150612542565b808061253a90613f88565b91505061244c565b50600d5f9054906101000a900460ff16801561256a5750600a5f9054906101000a900460ff16155b80156125bf575060115f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff16155b80156125c9575080155b801561261c575060135f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b801561266f575060135f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b80156126855750600b546126823061153c565b10155b156126c6576001600a5f6101000a81548160ff0219169083151502179055506126ac612963565b5f600a5f6101000a81548160ff0219169083151502179055505b5f60135f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16158015612765575060135f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b905080156129505760115f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff161561287f575f6127c885612f28565b90505f8160a001511115612879575f60648260a00151866127e99190614002565b6127f39190614070565b90508060105f828254612806919061432b565b925050819055508085612819919061432b565b6128228861153c565b101561285a576040517f136c487600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612865873083613177565b612870878787613177565b5050505061295e565b5061294f565b60115f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff161561294e575f6128da86612f28565b90505f8160800151111561294c575f60648260800151866128fb9190614002565b6129059190614070565b90508060105f828254612918919061432b565b9250508190555061292a873083613177565b80856129369190613fcf565b9450612943878787613177565b5050505061295e565b505b5b5b61295b858585613177565b50505b505050565b5f4790505f6129713061153c565b90505f601054036129a7575f8111801561298a57505f82115b156129a05761299981476132fa565b50506129ea565b50506129ea565b5f6002826129b59190614070565b90506129cb81836129c69190613fcf565b613416565b5f4790505f8190505f6010819055506129e483826132fa565b50505050505b565b6129f4611cef565b6129fd5f612a7e565b565b612a79838473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8585604051602401612a3292919061435e565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050613609565b505050565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603612bb1575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401612ba89190613e7a565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603612c21575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401612c189190613e7a565b60405180910390fd5b815f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015612d09578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051612d009190613b92565b60405180910390a35b50505050565b5f6014545f1415905090565b612d2361384d565b612d2b61384d565b5f5b601280549050811015612f20575f60128281548110612d4f57612d4e613f2e565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060115f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160019054906101000a900460ff1615612f0c5760115f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200160028201548152602001600382015481525050925050612f20565b508080612f1890613f88565b915050612d2d565b508091505090565b612f3061384d565b612f3861384d565b5f5b60128054905081101561316d575f60128281548110612f5c57612f5b613f2e565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508473ffffffffffffffffffffffffffffffffffffffff1660115f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036131595760115f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820154815260200160038201548152505092505061316d565b50808061316590613f88565b915050612f3a565b5080915050919050565b5f60065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015613201578381836040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016131f893929190614198565b60405180910390fd5b81810360065f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508160065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516132ec9190613b92565b60405180910390a350505050565b5f613303612d1b565b905061331430826040015185611d7d565b5f61331e84610c8e565b90505f61332a84610c8e565b90505f805f856040015173ffffffffffffffffffffffffffffffffffffffff1663f305d71988308b898930426040518863ffffffff1660e01b815260040161337796959493929190614385565b60606040518083038185885af1158015613393573d5f803e3d5ffd5b50505050506040513d601f19601f820116820180604052508101906133b891906143e4565b925092509250828511806133cb57508184115b806133d557505f81145b1561340c576040517f0bc488c500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050565b5f61341f612d1b565b90505f600267ffffffffffffffff81111561343d5761343c614434565b5b60405190808252806020026020018201604052801561346b5781602001602082028036833780820191505090505b50905030815f8151811061348257613481613f2e565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050816040015173ffffffffffffffffffffffffffffffffffffffff1663ad5c46486040518163ffffffff1660e01b8152600401602060405180830381865afa158015613509573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061352d9190614475565b8160018151811061354157613540613f2e565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505061358a30836040015185611d7d565b816040015173ffffffffffffffffffffffffffffffffffffffff1663791ac947846135b486610c8e565b8430426040518663ffffffff1660e01b81526004016135d7959493929190614557565b5f604051808303815f87803b1580156135ee575f80fd5b505af1158015613600573d5f803e3d5ffd5b50505050505050565b5f613633828473ffffffffffffffffffffffffffffffffffffffff1661369e90919063ffffffff16565b90505f81511415801561365757508080602001905181019061365591906145c3565b155b1561369957826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016136909190613e7a565b60405180910390fd5b505050565b60606136ab83835f6136b3565b905092915050565b6060814710156136fa57306040517fcd7860590000000000000000000000000000000000000000000000000000000081526004016136f19190613e7a565b60405180910390fd5b5f808573ffffffffffffffffffffffffffffffffffffffff1684866040516137229190614628565b5f6040518083038185875af1925050503d805f811461375c576040519150601f19603f3d011682016040523d82523d5f602084013e613761565b606091505b509150915061377186838361377c565b925050509392505050565b6060826137915761378c82613809565b613801565b5f82511480156137b757505f8473ffffffffffffffffffffffffffffffffffffffff163b145b156137f957836040517f9996b3150000000000000000000000000000000000000000000000000000000081526004016137f09190613e7a565b60405180910390fd5b819050613802565b5b9392505050565b5f8151111561381b5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060c001604052805f151581526020015f151581526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f81526020015f81525090565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156138e45780820151818401526020810190506138c9565b5f8484015250505050565b5f601f19601f8301169050919050565b5f613909826138ad565b61391381856138b7565b93506139238185602086016138c7565b61392c816138ef565b840191505092915050565b5f6020820190508181035f83015261394f81846138ff565b905092915050565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6139888261395f565b9050919050565b6139988161397e565b81146139a2575f80fd5b50565b5f813590506139b38161398f565b92915050565b5f8115159050919050565b6139cd816139b9565b81146139d7575f80fd5b50565b5f813590506139e8816139c4565b92915050565b5f8060408385031215613a0457613a03613957565b5b5f613a11858286016139a5565b9250506020613a22858286016139da565b9150509250929050565b5f819050919050565b613a3e81613a2c565b8114613a48575f80fd5b50565b5f81359050613a5981613a35565b92915050565b5f8060408385031215613a7557613a74613957565b5b5f613a82858286016139a5565b9250506020613a9385828601613a4b565b9150509250929050565b613aa6816139b9565b82525050565b5f602082019050613abf5f830184613a9d565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f840112613ae657613ae5613ac5565b5b8235905067ffffffffffffffff811115613b0357613b02613ac9565b5b602083019150836020820283011115613b1f57613b1e613acd565b5b9250929050565b5f805f60408486031215613b3d57613b3c613957565b5b5f84013567ffffffffffffffff811115613b5a57613b5961395b565b5b613b6686828701613ad1565b93509350506020613b79868287016139da565b9150509250925092565b613b8c81613a2c565b82525050565b5f602082019050613ba55f830184613b83565b92915050565b5f8083601f840112613bc057613bbf613ac5565b5b8235905067ffffffffffffffff811115613bdd57613bdc613ac9565b5b602083019150836020820283011115613bf957613bf8613acd565b5b9250929050565b5f805f8060408587031215613c1857613c17613957565b5b5f85013567ffffffffffffffff811115613c3557613c3461395b565b5b613c4187828801613ad1565b9450945050602085013567ffffffffffffffff811115613c6457613c6361395b565b5b613c7087828801613bab565b925092505092959194509250565b5f805f60608486031215613c9557613c94613957565b5b5f613ca2868287016139a5565b9350506020613cb3868287016139a5565b9250506040613cc486828701613a4b565b9150509250925092565b5f60208284031215613ce357613ce2613957565b5b5f613cf084828501613a4b565b91505092915050565b5f60208284031215613d0e57613d0d613957565b5b5f613d1b848285016139a5565b91505092915050565b613d2d8161397e565b82525050565b5f60c082019050613d465f830189613a9d565b613d536020830188613a9d565b613d606040830187613d24565b613d6d6060830186613d24565b613d7a6080830185613b83565b613d8760a0830184613b83565b979650505050505050565b5f60ff82169050919050565b613da781613d92565b82525050565b5f602082019050613dc05f830184613d9e565b92915050565b5f805f805f8060c08789031215613de057613ddf613957565b5b5f613ded89828a016139a5565b9650506020613dfe89828a016139a5565b9550506040613e0f89828a016139da565b9450506060613e2089828a016139da565b9350506080613e3189828a01613a4b565b92505060a0613e4289828a01613a4b565b9150509295509295509295565b5f60208284031215613e6457613e63613957565b5b5f613e71848285016139da565b91505092915050565b5f602082019050613e8d5f830184613d24565b92915050565b5f8060408385031215613ea957613ea8613957565b5b5f613eb6858286016139a5565b9250506020613ec7858286016139a5565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680613f1557607f821691505b602082108103613f2857613f27613ed1565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f613f9282613a2c565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203613fc457613fc3613f5b565b5b600182019050919050565b5f613fd982613a2c565b9150613fe483613a2c565b9250828203905081811115613ffc57613ffb613f5b565b5b92915050565b5f61400c82613a2c565b915061401783613a2c565b925082820261402581613a2c565b9150828204841483151761403c5761403b613f5b565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61407a82613a2c565b915061408583613a2c565b92508261409557614094614043565b5b828204905092915050565b5f6080820190506140b35f830187613a9d565b6140c06020830186613a9d565b6140cd6040830185613b83565b6140da6060830184613b83565b95945050505050565b5f815190506140f181613a35565b92915050565b5f6020828403121561410c5761410b613957565b5b5f614119848285016140e3565b91505092915050565b5f6060820190506141355f830186613d24565b6141426020830185613d24565b61414f6040830184613b83565b949350505050565b5f81905092915050565b50565b5f61416f5f83614157565b915061417a82614161565b5f82019050919050565b5f61418e82614164565b9150819050919050565b5f6060820190506141ab5f830186613d24565b6141b86020830185613b83565b6141c56040830184613b83565b949350505050565b7f616d6f756e74206d7573742062652067726561746572207468616e20300000005f82015250565b5f614201601d836138b7565b915061420c826141cd565b602082019050919050565b5f6020820190508181035f83015261422e816141f5565b9050919050565b7f54726164696e67206973206e6f74206163746976652e000000000000000000005f82015250565b5f6142696016836138b7565b915061427482614235565b602082019050919050565b5f6020820190508181035f8301526142968161425d565b9050919050565b7f5f7472616e736665723a3a205472616e736665722044656c617920656e61626c5f8201527f65642e202054727920616761696e206c617465722e0000000000000000000000602082015250565b5f6142f76035836138b7565b91506143028261429d565b604082019050919050565b5f6020820190508181035f830152614324816142eb565b9050919050565b5f61433582613a2c565b915061434083613a2c565b925082820190508082111561435857614357613f5b565b5b92915050565b5f6040820190506143715f830185613d24565b61437e6020830184613b83565b9392505050565b5f60c0820190506143985f830189613d24565b6143a56020830188613b83565b6143b26040830187613b83565b6143bf6060830186613b83565b6143cc6080830185613d24565b6143d960a0830184613b83565b979650505050505050565b5f805f606084860312156143fb576143fa613957565b5b5f614408868287016140e3565b9350506020614419868287016140e3565b925050604061442a868287016140e3565b9150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f8151905061446f8161398f565b92915050565b5f6020828403121561448a57614489613957565b5b5f61449784828501614461565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b6144d28161397e565b82525050565b5f6144e383836144c9565b60208301905092915050565b5f602082019050919050565b5f614505826144a0565b61450f81856144aa565b935061451a836144ba565b805f5b8381101561454a57815161453188826144d8565b975061453c836144ef565b92505060018101905061451d565b5085935050505092915050565b5f60a08201905061456a5f830188613b83565b6145776020830187613b83565b818103604083015261458981866144fb565b90506145986060830185613d24565b6145a56080830184613b83565b9695505050505050565b5f815190506145bd816139c4565b92915050565b5f602082840312156145d8576145d7613957565b5b5f6145e5848285016145af565b91505092915050565b5f81519050919050565b5f614602826145ee565b61460c8185614157565b935061461c8185602086016138c7565b80840191505092915050565b5f61463382846145f8565b91508190509291505056fea2646970667358221220de731a4cbfe48c9fd9044c529a8fed5622c6a0cc4dbd415909772686d658162564736f6c63430008140033",
}

// CoindistributionABI is the input ABI used to generate the binding from.
// Deprecated: Use CoindistributionMetaData.ABI instead.
var CoindistributionABI = CoindistributionMetaData.ABI

// CoindistributionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CoindistributionMetaData.Bin instead.
var CoindistributionBin = CoindistributionMetaData.Bin

// DeployCoindistribution deploys a new Ethereum contract, binding an instance of Coindistribution to it.
func DeployCoindistribution(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Coindistribution, error) {
	parsed, err := CoindistributionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CoindistributionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Coindistribution{CoindistributionCaller: CoindistributionCaller{contract: contract}, CoindistributionTransactor: CoindistributionTransactor{contract: contract}, CoindistributionFilterer: CoindistributionFilterer{contract: contract}}, nil
}

// Coindistribution is an auto generated Go binding around an Ethereum contract.
type Coindistribution struct {
	CoindistributionCaller     // Read-only binding to the contract
	CoindistributionTransactor // Write-only binding to the contract
	CoindistributionFilterer   // Log filterer for contract events
}

// CoindistributionCaller is an auto generated read-only Go binding around an Ethereum contract.
type CoindistributionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoindistributionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CoindistributionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoindistributionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CoindistributionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoindistributionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CoindistributionSession struct {
	Contract     *Coindistribution // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CoindistributionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CoindistributionCallerSession struct {
	Contract *CoindistributionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// CoindistributionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CoindistributionTransactorSession struct {
	Contract     *CoindistributionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// CoindistributionRaw is an auto generated low-level Go binding around an Ethereum contract.
type CoindistributionRaw struct {
	Contract *Coindistribution // Generic contract binding to access the raw methods on
}

// CoindistributionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CoindistributionCallerRaw struct {
	Contract *CoindistributionCaller // Generic read-only contract binding to access the raw methods on
}

// CoindistributionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CoindistributionTransactorRaw struct {
	Contract *CoindistributionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCoindistribution creates a new instance of Coindistribution, bound to a specific deployed contract.
func NewCoindistribution(address common.Address, backend bind.ContractBackend) (*Coindistribution, error) {
	contract, err := bindCoindistribution(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Coindistribution{CoindistributionCaller: CoindistributionCaller{contract: contract}, CoindistributionTransactor: CoindistributionTransactor{contract: contract}, CoindistributionFilterer: CoindistributionFilterer{contract: contract}}, nil
}

// NewCoindistributionCaller creates a new read-only instance of Coindistribution, bound to a specific deployed contract.
func NewCoindistributionCaller(address common.Address, caller bind.ContractCaller) (*CoindistributionCaller, error) {
	contract, err := bindCoindistribution(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CoindistributionCaller{contract: contract}, nil
}

// NewCoindistributionTransactor creates a new write-only instance of Coindistribution, bound to a specific deployed contract.
func NewCoindistributionTransactor(address common.Address, transactor bind.ContractTransactor) (*CoindistributionTransactor, error) {
	contract, err := bindCoindistribution(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CoindistributionTransactor{contract: contract}, nil
}

// NewCoindistributionFilterer creates a new log filterer instance of Coindistribution, bound to a specific deployed contract.
func NewCoindistributionFilterer(address common.Address, filterer bind.ContractFilterer) (*CoindistributionFilterer, error) {
	contract, err := bindCoindistribution(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CoindistributionFilterer{contract: contract}, nil
}

// bindCoindistribution binds a generic wrapper to an already deployed contract.
func bindCoindistribution(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CoindistributionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Coindistribution *CoindistributionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Coindistribution.Contract.CoindistributionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Coindistribution *CoindistributionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.Contract.CoindistributionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Coindistribution *CoindistributionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Coindistribution.Contract.CoindistributionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Coindistribution *CoindistributionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Coindistribution.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Coindistribution *CoindistributionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Coindistribution *CoindistributionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Coindistribution.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Coindistribution *CoindistributionCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Coindistribution *CoindistributionSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.Allowance(&_Coindistribution.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.Allowance(&_Coindistribution.CallOpts, owner, spender)
}

// ApplySlippage is a free data retrieval call binding the contract method 0x27ac2d68.
//
// Solidity: function applySlippage(uint256 value) view returns(uint256)
func (_Coindistribution *CoindistributionCaller) ApplySlippage(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "applySlippage", value)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ApplySlippage is a free data retrieval call binding the contract method 0x27ac2d68.
//
// Solidity: function applySlippage(uint256 value) view returns(uint256)
func (_Coindistribution *CoindistributionSession) ApplySlippage(value *big.Int) (*big.Int, error) {
	return _Coindistribution.Contract.ApplySlippage(&_Coindistribution.CallOpts, value)
}

// ApplySlippage is a free data retrieval call binding the contract method 0x27ac2d68.
//
// Solidity: function applySlippage(uint256 value) view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) ApplySlippage(value *big.Int) (*big.Int, error) {
	return _Coindistribution.Contract.ApplySlippage(&_Coindistribution.CallOpts, value)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Coindistribution *CoindistributionCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Coindistribution *CoindistributionSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.BalanceOf(&_Coindistribution.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.BalanceOf(&_Coindistribution.CallOpts, account)
}

// Bots is a free data retrieval call binding the contract method 0xbfd79284.
//
// Solidity: function bots(address ) view returns(bool)
func (_Coindistribution *CoindistributionCaller) Bots(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "bots", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Bots is a free data retrieval call binding the contract method 0xbfd79284.
//
// Solidity: function bots(address ) view returns(bool)
func (_Coindistribution *CoindistributionSession) Bots(arg0 common.Address) (bool, error) {
	return _Coindistribution.Contract.Bots(&_Coindistribution.CallOpts, arg0)
}

// Bots is a free data retrieval call binding the contract method 0xbfd79284.
//
// Solidity: function bots(address ) view returns(bool)
func (_Coindistribution *CoindistributionCallerSession) Bots(arg0 common.Address) (bool, error) {
	return _Coindistribution.Contract.Bots(&_Coindistribution.CallOpts, arg0)
}

// BotsCaught is a free data retrieval call binding the contract method 0x6b0a894c.
//
// Solidity: function botsCaught() view returns(uint256)
func (_Coindistribution *CoindistributionCaller) BotsCaught(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "botsCaught")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BotsCaught is a free data retrieval call binding the contract method 0x6b0a894c.
//
// Solidity: function botsCaught() view returns(uint256)
func (_Coindistribution *CoindistributionSession) BotsCaught() (*big.Int, error) {
	return _Coindistribution.Contract.BotsCaught(&_Coindistribution.CallOpts)
}

// BotsCaught is a free data retrieval call binding the contract method 0x6b0a894c.
//
// Solidity: function botsCaught() view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) BotsCaught() (*big.Int, error) {
	return _Coindistribution.Contract.BotsCaught(&_Coindistribution.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Coindistribution *CoindistributionCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Coindistribution *CoindistributionSession) Decimals() (uint8, error) {
	return _Coindistribution.Contract.Decimals(&_Coindistribution.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Coindistribution *CoindistributionCallerSession) Decimals() (uint8, error) {
	return _Coindistribution.Contract.Decimals(&_Coindistribution.CallOpts)
}

// GetAirDropper is a free data retrieval call binding the contract method 0x5db7548b.
//
// Solidity: function getAirDropper() view returns(address)
func (_Coindistribution *CoindistributionCaller) GetAirDropper(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "getAirDropper")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAirDropper is a free data retrieval call binding the contract method 0x5db7548b.
//
// Solidity: function getAirDropper() view returns(address)
func (_Coindistribution *CoindistributionSession) GetAirDropper() (common.Address, error) {
	return _Coindistribution.Contract.GetAirDropper(&_Coindistribution.CallOpts)
}

// GetAirDropper is a free data retrieval call binding the contract method 0x5db7548b.
//
// Solidity: function getAirDropper() view returns(address)
func (_Coindistribution *CoindistributionCallerSession) GetAirDropper() (common.Address, error) {
	return _Coindistribution.Contract.GetAirDropper(&_Coindistribution.CallOpts)
}

// GetTradingEnabledBlock is a free data retrieval call binding the contract method 0x588b655f.
//
// Solidity: function getTradingEnabledBlock() view returns(uint256)
func (_Coindistribution *CoindistributionCaller) GetTradingEnabledBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "getTradingEnabledBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTradingEnabledBlock is a free data retrieval call binding the contract method 0x588b655f.
//
// Solidity: function getTradingEnabledBlock() view returns(uint256)
func (_Coindistribution *CoindistributionSession) GetTradingEnabledBlock() (*big.Int, error) {
	return _Coindistribution.Contract.GetTradingEnabledBlock(&_Coindistribution.CallOpts)
}

// GetTradingEnabledBlock is a free data retrieval call binding the contract method 0x588b655f.
//
// Solidity: function getTradingEnabledBlock() view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) GetTradingEnabledBlock() (*big.Int, error) {
	return _Coindistribution.Contract.GetTradingEnabledBlock(&_Coindistribution.CallOpts)
}

// LimitsInEffect is a free data retrieval call binding the contract method 0x4a62bb65.
//
// Solidity: function limitsInEffect() view returns(bool)
func (_Coindistribution *CoindistributionCaller) LimitsInEffect(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "limitsInEffect")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// LimitsInEffect is a free data retrieval call binding the contract method 0x4a62bb65.
//
// Solidity: function limitsInEffect() view returns(bool)
func (_Coindistribution *CoindistributionSession) LimitsInEffect() (bool, error) {
	return _Coindistribution.Contract.LimitsInEffect(&_Coindistribution.CallOpts)
}

// LimitsInEffect is a free data retrieval call binding the contract method 0x4a62bb65.
//
// Solidity: function limitsInEffect() view returns(bool)
func (_Coindistribution *CoindistributionCallerSession) LimitsInEffect() (bool, error) {
	return _Coindistribution.Contract.LimitsInEffect(&_Coindistribution.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Coindistribution *CoindistributionCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Coindistribution *CoindistributionSession) Name() (string, error) {
	return _Coindistribution.Contract.Name(&_Coindistribution.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Coindistribution *CoindistributionCallerSession) Name() (string, error) {
	return _Coindistribution.Contract.Name(&_Coindistribution.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Coindistribution *CoindistributionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Coindistribution *CoindistributionSession) Owner() (common.Address, error) {
	return _Coindistribution.Contract.Owner(&_Coindistribution.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Coindistribution *CoindistributionCallerSession) Owner() (common.Address, error) {
	return _Coindistribution.Contract.Owner(&_Coindistribution.CallOpts)
}

// SwapEnabled is a free data retrieval call binding the contract method 0x6ddd1713.
//
// Solidity: function swapEnabled() view returns(bool)
func (_Coindistribution *CoindistributionCaller) SwapEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "swapEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SwapEnabled is a free data retrieval call binding the contract method 0x6ddd1713.
//
// Solidity: function swapEnabled() view returns(bool)
func (_Coindistribution *CoindistributionSession) SwapEnabled() (bool, error) {
	return _Coindistribution.Contract.SwapEnabled(&_Coindistribution.CallOpts)
}

// SwapEnabled is a free data retrieval call binding the contract method 0x6ddd1713.
//
// Solidity: function swapEnabled() view returns(bool)
func (_Coindistribution *CoindistributionCallerSession) SwapEnabled() (bool, error) {
	return _Coindistribution.Contract.SwapEnabled(&_Coindistribution.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Coindistribution *CoindistributionCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Coindistribution *CoindistributionSession) Symbol() (string, error) {
	return _Coindistribution.Contract.Symbol(&_Coindistribution.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Coindistribution *CoindistributionCallerSession) Symbol() (string, error) {
	return _Coindistribution.Contract.Symbol(&_Coindistribution.CallOpts)
}

// TokensForLiquidity is a free data retrieval call binding the contract method 0x1a8145bb.
//
// Solidity: function tokensForLiquidity() view returns(uint256)
func (_Coindistribution *CoindistributionCaller) TokensForLiquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "tokensForLiquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokensForLiquidity is a free data retrieval call binding the contract method 0x1a8145bb.
//
// Solidity: function tokensForLiquidity() view returns(uint256)
func (_Coindistribution *CoindistributionSession) TokensForLiquidity() (*big.Int, error) {
	return _Coindistribution.Contract.TokensForLiquidity(&_Coindistribution.CallOpts)
}

// TokensForLiquidity is a free data retrieval call binding the contract method 0x1a8145bb.
//
// Solidity: function tokensForLiquidity() view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) TokensForLiquidity() (*big.Int, error) {
	return _Coindistribution.Contract.TokensForLiquidity(&_Coindistribution.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Coindistribution *CoindistributionCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Coindistribution *CoindistributionSession) TotalSupply() (*big.Int, error) {
	return _Coindistribution.Contract.TotalSupply(&_Coindistribution.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) TotalSupply() (*big.Int, error) {
	return _Coindistribution.Contract.TotalSupply(&_Coindistribution.CallOpts)
}

// TransferDelayEnabled is a free data retrieval call binding the contract method 0xc876d0b9.
//
// Solidity: function transferDelayEnabled() view returns(bool)
func (_Coindistribution *CoindistributionCaller) TransferDelayEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "transferDelayEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TransferDelayEnabled is a free data retrieval call binding the contract method 0xc876d0b9.
//
// Solidity: function transferDelayEnabled() view returns(bool)
func (_Coindistribution *CoindistributionSession) TransferDelayEnabled() (bool, error) {
	return _Coindistribution.Contract.TransferDelayEnabled(&_Coindistribution.CallOpts)
}

// TransferDelayEnabled is a free data retrieval call binding the contract method 0xc876d0b9.
//
// Solidity: function transferDelayEnabled() view returns(bool)
func (_Coindistribution *CoindistributionCallerSession) TransferDelayEnabled() (bool, error) {
	return _Coindistribution.Contract.TransferDelayEnabled(&_Coindistribution.CallOpts)
}

// UniswapV2LiquidityPoolSlots is a free data retrieval call binding the contract method 0x296a803c.
//
// Solidity: function uniswapV2LiquidityPoolSlots(address ) view returns(bool enabled, bool receiver, address router, address pairAddress, uint256 bLiqF, uint256 sLiqF)
func (_Coindistribution *CoindistributionCaller) UniswapV2LiquidityPoolSlots(opts *bind.CallOpts, arg0 common.Address) (struct {
	Enabled     bool
	Receiver    bool
	Router      common.Address
	PairAddress common.Address
	BLiqF       *big.Int
	SLiqF       *big.Int
}, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "uniswapV2LiquidityPoolSlots", arg0)

	outstruct := new(struct {
		Enabled     bool
		Receiver    bool
		Router      common.Address
		PairAddress common.Address
		BLiqF       *big.Int
		SLiqF       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Enabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Receiver = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.PairAddress = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.BLiqF = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.SLiqF = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UniswapV2LiquidityPoolSlots is a free data retrieval call binding the contract method 0x296a803c.
//
// Solidity: function uniswapV2LiquidityPoolSlots(address ) view returns(bool enabled, bool receiver, address router, address pairAddress, uint256 bLiqF, uint256 sLiqF)
func (_Coindistribution *CoindistributionSession) UniswapV2LiquidityPoolSlots(arg0 common.Address) (struct {
	Enabled     bool
	Receiver    bool
	Router      common.Address
	PairAddress common.Address
	BLiqF       *big.Int
	SLiqF       *big.Int
}, error) {
	return _Coindistribution.Contract.UniswapV2LiquidityPoolSlots(&_Coindistribution.CallOpts, arg0)
}

// UniswapV2LiquidityPoolSlots is a free data retrieval call binding the contract method 0x296a803c.
//
// Solidity: function uniswapV2LiquidityPoolSlots(address ) view returns(bool enabled, bool receiver, address router, address pairAddress, uint256 bLiqF, uint256 sLiqF)
func (_Coindistribution *CoindistributionCallerSession) UniswapV2LiquidityPoolSlots(arg0 common.Address) (struct {
	Enabled     bool
	Receiver    bool
	Router      common.Address
	PairAddress common.Address
	BLiqF       *big.Int
	SLiqF       *big.Int
}, error) {
	return _Coindistribution.Contract.UniswapV2LiquidityPoolSlots(&_Coindistribution.CallOpts, arg0)
}

// UniswapV2LiquidityPools is a free data retrieval call binding the contract method 0x785cca3e.
//
// Solidity: function uniswapV2LiquidityPools(uint256 ) view returns(address)
func (_Coindistribution *CoindistributionCaller) UniswapV2LiquidityPools(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "uniswapV2LiquidityPools", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniswapV2LiquidityPools is a free data retrieval call binding the contract method 0x785cca3e.
//
// Solidity: function uniswapV2LiquidityPools(uint256 ) view returns(address)
func (_Coindistribution *CoindistributionSession) UniswapV2LiquidityPools(arg0 *big.Int) (common.Address, error) {
	return _Coindistribution.Contract.UniswapV2LiquidityPools(&_Coindistribution.CallOpts, arg0)
}

// UniswapV2LiquidityPools is a free data retrieval call binding the contract method 0x785cca3e.
//
// Solidity: function uniswapV2LiquidityPools(uint256 ) view returns(address)
func (_Coindistribution *CoindistributionCallerSession) UniswapV2LiquidityPools(arg0 *big.Int) (common.Address, error) {
	return _Coindistribution.Contract.UniswapV2LiquidityPools(&_Coindistribution.CallOpts, arg0)
}

// AirdropToWallets is a paid mutator transaction binding the contract method 0x2307b441.
//
// Solidity: function airdropToWallets(address[] recipients, uint256[] amounts) returns()
func (_Coindistribution *CoindistributionTransactor) AirdropToWallets(opts *bind.TransactOpts, recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "airdropToWallets", recipients, amounts)
}

// AirdropToWallets is a paid mutator transaction binding the contract method 0x2307b441.
//
// Solidity: function airdropToWallets(address[] recipients, uint256[] amounts) returns()
func (_Coindistribution *CoindistributionSession) AirdropToWallets(recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.AirdropToWallets(&_Coindistribution.TransactOpts, recipients, amounts)
}

// AirdropToWallets is a paid mutator transaction binding the contract method 0x2307b441.
//
// Solidity: function airdropToWallets(address[] recipients, uint256[] amounts) returns()
func (_Coindistribution *CoindistributionTransactorSession) AirdropToWallets(recipients []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.AirdropToWallets(&_Coindistribution.TransactOpts, recipients, amounts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.Approve(&_Coindistribution.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.Approve(&_Coindistribution.TransactOpts, spender, value)
}

// DisableTransferDelay is a paid mutator transaction binding the contract method 0xe884f260.
//
// Solidity: function disableTransferDelay() returns()
func (_Coindistribution *CoindistributionTransactor) DisableTransferDelay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "disableTransferDelay")
}

// DisableTransferDelay is a paid mutator transaction binding the contract method 0xe884f260.
//
// Solidity: function disableTransferDelay() returns()
func (_Coindistribution *CoindistributionSession) DisableTransferDelay() (*types.Transaction, error) {
	return _Coindistribution.Contract.DisableTransferDelay(&_Coindistribution.TransactOpts)
}

// DisableTransferDelay is a paid mutator transaction binding the contract method 0xe884f260.
//
// Solidity: function disableTransferDelay() returns()
func (_Coindistribution *CoindistributionTransactorSession) DisableTransferDelay() (*types.Transaction, error) {
	return _Coindistribution.Contract.DisableTransferDelay(&_Coindistribution.TransactOpts)
}

// EnableTrading is a paid mutator transaction binding the contract method 0x8a8c523c.
//
// Solidity: function enableTrading() returns()
func (_Coindistribution *CoindistributionTransactor) EnableTrading(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "enableTrading")
}

// EnableTrading is a paid mutator transaction binding the contract method 0x8a8c523c.
//
// Solidity: function enableTrading() returns()
func (_Coindistribution *CoindistributionSession) EnableTrading() (*types.Transaction, error) {
	return _Coindistribution.Contract.EnableTrading(&_Coindistribution.TransactOpts)
}

// EnableTrading is a paid mutator transaction binding the contract method 0x8a8c523c.
//
// Solidity: function enableTrading() returns()
func (_Coindistribution *CoindistributionTransactorSession) EnableTrading() (*types.Transaction, error) {
	return _Coindistribution.Contract.EnableTrading(&_Coindistribution.TransactOpts)
}

// ExcludeFromFees is a paid mutator transaction binding the contract method 0xc0246668.
//
// Solidity: function excludeFromFees(address account, bool excluded) returns()
func (_Coindistribution *CoindistributionTransactor) ExcludeFromFees(opts *bind.TransactOpts, account common.Address, excluded bool) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "excludeFromFees", account, excluded)
}

// ExcludeFromFees is a paid mutator transaction binding the contract method 0xc0246668.
//
// Solidity: function excludeFromFees(address account, bool excluded) returns()
func (_Coindistribution *CoindistributionSession) ExcludeFromFees(account common.Address, excluded bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.ExcludeFromFees(&_Coindistribution.TransactOpts, account, excluded)
}

// ExcludeFromFees is a paid mutator transaction binding the contract method 0xc0246668.
//
// Solidity: function excludeFromFees(address account, bool excluded) returns()
func (_Coindistribution *CoindistributionTransactorSession) ExcludeFromFees(account common.Address, excluded bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.ExcludeFromFees(&_Coindistribution.TransactOpts, account, excluded)
}

// ForceSwapBack is a paid mutator transaction binding the contract method 0x51f205e4.
//
// Solidity: function forceSwapBack() returns()
func (_Coindistribution *CoindistributionTransactor) ForceSwapBack(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "forceSwapBack")
}

// ForceSwapBack is a paid mutator transaction binding the contract method 0x51f205e4.
//
// Solidity: function forceSwapBack() returns()
func (_Coindistribution *CoindistributionSession) ForceSwapBack() (*types.Transaction, error) {
	return _Coindistribution.Contract.ForceSwapBack(&_Coindistribution.TransactOpts)
}

// ForceSwapBack is a paid mutator transaction binding the contract method 0x51f205e4.
//
// Solidity: function forceSwapBack() returns()
func (_Coindistribution *CoindistributionTransactorSession) ForceSwapBack() (*types.Transaction, error) {
	return _Coindistribution.Contract.ForceSwapBack(&_Coindistribution.TransactOpts)
}

// ManageBoughtEarly is a paid mutator transaction binding the contract method 0xe800dff7.
//
// Solidity: function manageBoughtEarly(address wallet, bool flag) returns()
func (_Coindistribution *CoindistributionTransactor) ManageBoughtEarly(opts *bind.TransactOpts, wallet common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "manageBoughtEarly", wallet, flag)
}

// ManageBoughtEarly is a paid mutator transaction binding the contract method 0xe800dff7.
//
// Solidity: function manageBoughtEarly(address wallet, bool flag) returns()
func (_Coindistribution *CoindistributionSession) ManageBoughtEarly(wallet common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.ManageBoughtEarly(&_Coindistribution.TransactOpts, wallet, flag)
}

// ManageBoughtEarly is a paid mutator transaction binding the contract method 0xe800dff7.
//
// Solidity: function manageBoughtEarly(address wallet, bool flag) returns()
func (_Coindistribution *CoindistributionTransactorSession) ManageBoughtEarly(wallet common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.ManageBoughtEarly(&_Coindistribution.TransactOpts, wallet, flag)
}

// MassManageBoughtEarly is a paid mutator transaction binding the contract method 0x130a2c3c.
//
// Solidity: function massManageBoughtEarly(address[] wallets, bool flag) returns()
func (_Coindistribution *CoindistributionTransactor) MassManageBoughtEarly(opts *bind.TransactOpts, wallets []common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "massManageBoughtEarly", wallets, flag)
}

// MassManageBoughtEarly is a paid mutator transaction binding the contract method 0x130a2c3c.
//
// Solidity: function massManageBoughtEarly(address[] wallets, bool flag) returns()
func (_Coindistribution *CoindistributionSession) MassManageBoughtEarly(wallets []common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.MassManageBoughtEarly(&_Coindistribution.TransactOpts, wallets, flag)
}

// MassManageBoughtEarly is a paid mutator transaction binding the contract method 0x130a2c3c.
//
// Solidity: function massManageBoughtEarly(address[] wallets, bool flag) returns()
func (_Coindistribution *CoindistributionTransactorSession) MassManageBoughtEarly(wallets []common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.MassManageBoughtEarly(&_Coindistribution.TransactOpts, wallets, flag)
}

// RemoveLimits is a paid mutator transaction binding the contract method 0x751039fc.
//
// Solidity: function removeLimits() returns()
func (_Coindistribution *CoindistributionTransactor) RemoveLimits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "removeLimits")
}

// RemoveLimits is a paid mutator transaction binding the contract method 0x751039fc.
//
// Solidity: function removeLimits() returns()
func (_Coindistribution *CoindistributionSession) RemoveLimits() (*types.Transaction, error) {
	return _Coindistribution.Contract.RemoveLimits(&_Coindistribution.TransactOpts)
}

// RemoveLimits is a paid mutator transaction binding the contract method 0x751039fc.
//
// Solidity: function removeLimits() returns()
func (_Coindistribution *CoindistributionTransactorSession) RemoveLimits() (*types.Transaction, error) {
	return _Coindistribution.Contract.RemoveLimits(&_Coindistribution.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Coindistribution *CoindistributionTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Coindistribution *CoindistributionSession) RenounceOwnership() (*types.Transaction, error) {
	return _Coindistribution.Contract.RenounceOwnership(&_Coindistribution.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Coindistribution *CoindistributionTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Coindistribution.Contract.RenounceOwnership(&_Coindistribution.TransactOpts)
}

// SetAirDropper is a paid mutator transaction binding the contract method 0x5de998ef.
//
// Solidity: function setAirDropper(address newAirDropper) returns()
func (_Coindistribution *CoindistributionTransactor) SetAirDropper(opts *bind.TransactOpts, newAirDropper common.Address) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setAirDropper", newAirDropper)
}

// SetAirDropper is a paid mutator transaction binding the contract method 0x5de998ef.
//
// Solidity: function setAirDropper(address newAirDropper) returns()
func (_Coindistribution *CoindistributionSession) SetAirDropper(newAirDropper common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetAirDropper(&_Coindistribution.TransactOpts, newAirDropper)
}

// SetAirDropper is a paid mutator transaction binding the contract method 0x5de998ef.
//
// Solidity: function setAirDropper(address newAirDropper) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetAirDropper(newAirDropper common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetAirDropper(&_Coindistribution.TransactOpts, newAirDropper)
}

// SetKnownUniswapRouters is a paid mutator transaction binding the contract method 0x44d8a785.
//
// Solidity: function setKnownUniswapRouters(bool flag) returns()
func (_Coindistribution *CoindistributionTransactor) SetKnownUniswapRouters(opts *bind.TransactOpts, flag bool) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setKnownUniswapRouters", flag)
}

// SetKnownUniswapRouters is a paid mutator transaction binding the contract method 0x44d8a785.
//
// Solidity: function setKnownUniswapRouters(bool flag) returns()
func (_Coindistribution *CoindistributionSession) SetKnownUniswapRouters(flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetKnownUniswapRouters(&_Coindistribution.TransactOpts, flag)
}

// SetKnownUniswapRouters is a paid mutator transaction binding the contract method 0x44d8a785.
//
// Solidity: function setKnownUniswapRouters(bool flag) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetKnownUniswapRouters(flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetKnownUniswapRouters(&_Coindistribution.TransactOpts, flag)
}

// SetSlippage is a paid mutator transaction binding the contract method 0xf0fa55a9.
//
// Solidity: function setSlippage(uint256 value) returns()
func (_Coindistribution *CoindistributionTransactor) SetSlippage(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setSlippage", value)
}

// SetSlippage is a paid mutator transaction binding the contract method 0xf0fa55a9.
//
// Solidity: function setSlippage(uint256 value) returns()
func (_Coindistribution *CoindistributionSession) SetSlippage(value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetSlippage(&_Coindistribution.TransactOpts, value)
}

// SetSlippage is a paid mutator transaction binding the contract method 0xf0fa55a9.
//
// Solidity: function setSlippage(uint256 value) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetSlippage(value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetSlippage(&_Coindistribution.TransactOpts, value)
}

// SetSwapBackThreshold is a paid mutator transaction binding the contract method 0x8fe62b8a.
//
// Solidity: function setSwapBackThreshold(uint256 value) returns()
func (_Coindistribution *CoindistributionTransactor) SetSwapBackThreshold(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setSwapBackThreshold", value)
}

// SetSwapBackThreshold is a paid mutator transaction binding the contract method 0x8fe62b8a.
//
// Solidity: function setSwapBackThreshold(uint256 value) returns()
func (_Coindistribution *CoindistributionSession) SetSwapBackThreshold(value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetSwapBackThreshold(&_Coindistribution.TransactOpts, value)
}

// SetSwapBackThreshold is a paid mutator transaction binding the contract method 0x8fe62b8a.
//
// Solidity: function setSwapBackThreshold(uint256 value) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetSwapBackThreshold(value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetSwapBackThreshold(&_Coindistribution.TransactOpts, value)
}

// SetUniswapRouter is a paid mutator transaction binding the contract method 0x072280c3.
//
// Solidity: function setUniswapRouter(address theAddress, bool flag) returns()
func (_Coindistribution *CoindistributionTransactor) SetUniswapRouter(opts *bind.TransactOpts, theAddress common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setUniswapRouter", theAddress, flag)
}

// SetUniswapRouter is a paid mutator transaction binding the contract method 0x072280c3.
//
// Solidity: function setUniswapRouter(address theAddress, bool flag) returns()
func (_Coindistribution *CoindistributionSession) SetUniswapRouter(theAddress common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetUniswapRouter(&_Coindistribution.TransactOpts, theAddress, flag)
}

// SetUniswapRouter is a paid mutator transaction binding the contract method 0x072280c3.
//
// Solidity: function setUniswapRouter(address theAddress, bool flag) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetUniswapRouter(theAddress common.Address, flag bool) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetUniswapRouter(&_Coindistribution.TransactOpts, theAddress, flag)
}

// SetUniswapV2LiquidityPool is a paid mutator transaction binding the contract method 0x3d6d11b0.
//
// Solidity: function setUniswapV2LiquidityPool(address routerAddress, address pairAddress, bool enabled, bool receiver, uint256 bLiqF, uint256 sLiqF) returns()
func (_Coindistribution *CoindistributionTransactor) SetUniswapV2LiquidityPool(opts *bind.TransactOpts, routerAddress common.Address, pairAddress common.Address, enabled bool, receiver bool, bLiqF *big.Int, sLiqF *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setUniswapV2LiquidityPool", routerAddress, pairAddress, enabled, receiver, bLiqF, sLiqF)
}

// SetUniswapV2LiquidityPool is a paid mutator transaction binding the contract method 0x3d6d11b0.
//
// Solidity: function setUniswapV2LiquidityPool(address routerAddress, address pairAddress, bool enabled, bool receiver, uint256 bLiqF, uint256 sLiqF) returns()
func (_Coindistribution *CoindistributionSession) SetUniswapV2LiquidityPool(routerAddress common.Address, pairAddress common.Address, enabled bool, receiver bool, bLiqF *big.Int, sLiqF *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetUniswapV2LiquidityPool(&_Coindistribution.TransactOpts, routerAddress, pairAddress, enabled, receiver, bLiqF, sLiqF)
}

// SetUniswapV2LiquidityPool is a paid mutator transaction binding the contract method 0x3d6d11b0.
//
// Solidity: function setUniswapV2LiquidityPool(address routerAddress, address pairAddress, bool enabled, bool receiver, uint256 bLiqF, uint256 sLiqF) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetUniswapV2LiquidityPool(routerAddress common.Address, pairAddress common.Address, enabled bool, receiver bool, bLiqF *big.Int, sLiqF *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetUniswapV2LiquidityPool(&_Coindistribution.TransactOpts, routerAddress, pairAddress, enabled, receiver, bLiqF, sLiqF)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.Transfer(&_Coindistribution.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.Transfer(&_Coindistribution.TransactOpts, to, value)
}

// TransferForeignToken is a paid mutator transaction binding the contract method 0x8366e79a.
//
// Solidity: function transferForeignToken(address _token, address _to) returns()
func (_Coindistribution *CoindistributionTransactor) TransferForeignToken(opts *bind.TransactOpts, _token common.Address, _to common.Address) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "transferForeignToken", _token, _to)
}

// TransferForeignToken is a paid mutator transaction binding the contract method 0x8366e79a.
//
// Solidity: function transferForeignToken(address _token, address _to) returns()
func (_Coindistribution *CoindistributionSession) TransferForeignToken(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferForeignToken(&_Coindistribution.TransactOpts, _token, _to)
}

// TransferForeignToken is a paid mutator transaction binding the contract method 0x8366e79a.
//
// Solidity: function transferForeignToken(address _token, address _to) returns()
func (_Coindistribution *CoindistributionTransactorSession) TransferForeignToken(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferForeignToken(&_Coindistribution.TransactOpts, _token, _to)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferFrom(&_Coindistribution.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Coindistribution *CoindistributionTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferFrom(&_Coindistribution.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Coindistribution *CoindistributionTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Coindistribution *CoindistributionSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferOwnership(&_Coindistribution.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Coindistribution *CoindistributionTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferOwnership(&_Coindistribution.TransactOpts, newOwner)
}

// WithdrawStuckETH is a paid mutator transaction binding the contract method 0xf5648a4f.
//
// Solidity: function withdrawStuckETH() returns()
func (_Coindistribution *CoindistributionTransactor) WithdrawStuckETH(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "withdrawStuckETH")
}

// WithdrawStuckETH is a paid mutator transaction binding the contract method 0xf5648a4f.
//
// Solidity: function withdrawStuckETH() returns()
func (_Coindistribution *CoindistributionSession) WithdrawStuckETH() (*types.Transaction, error) {
	return _Coindistribution.Contract.WithdrawStuckETH(&_Coindistribution.TransactOpts)
}

// WithdrawStuckETH is a paid mutator transaction binding the contract method 0xf5648a4f.
//
// Solidity: function withdrawStuckETH() returns()
func (_Coindistribution *CoindistributionTransactorSession) WithdrawStuckETH() (*types.Transaction, error) {
	return _Coindistribution.Contract.WithdrawStuckETH(&_Coindistribution.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Coindistribution *CoindistributionTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Coindistribution.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Coindistribution *CoindistributionSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Coindistribution.Contract.Fallback(&_Coindistribution.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Coindistribution *CoindistributionTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Coindistribution.Contract.Fallback(&_Coindistribution.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Coindistribution *CoindistributionTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Coindistribution *CoindistributionSession) Receive() (*types.Transaction, error) {
	return _Coindistribution.Contract.Receive(&_Coindistribution.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Coindistribution *CoindistributionTransactorSession) Receive() (*types.Transaction, error) {
	return _Coindistribution.Contract.Receive(&_Coindistribution.TransactOpts)
}

// CoindistributionApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Coindistribution contract.
type CoindistributionApprovalIterator struct {
	Event *CoindistributionApproval // Event containing the contract specifics and raw log

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
func (it *CoindistributionApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionApproval)
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
		it.Event = new(CoindistributionApproval)
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
func (it *CoindistributionApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionApproval represents a Approval event raised by the Coindistribution contract.
type CoindistributionApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Coindistribution *CoindistributionFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CoindistributionApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionApprovalIterator{contract: _Coindistribution.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Coindistribution *CoindistributionFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CoindistributionApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionApproval)
				if err := _Coindistribution.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Coindistribution *CoindistributionFilterer) ParseApproval(log types.Log) (*CoindistributionApproval, error) {
	event := new(CoindistributionApproval)
	if err := _Coindistribution.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionBuyBackTriggeredIterator is returned from FilterBuyBackTriggered and is used to iterate over the raw logs and unpacked data for BuyBackTriggered events raised by the Coindistribution contract.
type CoindistributionBuyBackTriggeredIterator struct {
	Event *CoindistributionBuyBackTriggered // Event containing the contract specifics and raw log

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
func (it *CoindistributionBuyBackTriggeredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionBuyBackTriggered)
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
		it.Event = new(CoindistributionBuyBackTriggered)
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
func (it *CoindistributionBuyBackTriggeredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionBuyBackTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionBuyBackTriggered represents a BuyBackTriggered event raised by the Coindistribution contract.
type CoindistributionBuyBackTriggered struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBuyBackTriggered is a free log retrieval operation binding the contract event 0xa017c1567cfcdd2d750a8c01e39fe2a846bcebc293c7d078477014d684820568.
//
// Solidity: event BuyBackTriggered(uint256 amount)
func (_Coindistribution *CoindistributionFilterer) FilterBuyBackTriggered(opts *bind.FilterOpts) (*CoindistributionBuyBackTriggeredIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "BuyBackTriggered")
	if err != nil {
		return nil, err
	}
	return &CoindistributionBuyBackTriggeredIterator{contract: _Coindistribution.contract, event: "BuyBackTriggered", logs: logs, sub: sub}, nil
}

// WatchBuyBackTriggered is a free log subscription operation binding the contract event 0xa017c1567cfcdd2d750a8c01e39fe2a846bcebc293c7d078477014d684820568.
//
// Solidity: event BuyBackTriggered(uint256 amount)
func (_Coindistribution *CoindistributionFilterer) WatchBuyBackTriggered(opts *bind.WatchOpts, sink chan<- *CoindistributionBuyBackTriggered) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "BuyBackTriggered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionBuyBackTriggered)
				if err := _Coindistribution.contract.UnpackLog(event, "BuyBackTriggered", log); err != nil {
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

// ParseBuyBackTriggered is a log parse operation binding the contract event 0xa017c1567cfcdd2d750a8c01e39fe2a846bcebc293c7d078477014d684820568.
//
// Solidity: event BuyBackTriggered(uint256 amount)
func (_Coindistribution *CoindistributionFilterer) ParseBuyBackTriggered(log types.Log) (*CoindistributionBuyBackTriggered, error) {
	event := new(CoindistributionBuyBackTriggered)
	if err := _Coindistribution.contract.UnpackLog(event, "BuyBackTriggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionCaughtEarlyBuyerIterator is returned from FilterCaughtEarlyBuyer and is used to iterate over the raw logs and unpacked data for CaughtEarlyBuyer events raised by the Coindistribution contract.
type CoindistributionCaughtEarlyBuyerIterator struct {
	Event *CoindistributionCaughtEarlyBuyer // Event containing the contract specifics and raw log

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
func (it *CoindistributionCaughtEarlyBuyerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionCaughtEarlyBuyer)
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
		it.Event = new(CoindistributionCaughtEarlyBuyer)
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
func (it *CoindistributionCaughtEarlyBuyerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionCaughtEarlyBuyerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionCaughtEarlyBuyer represents a CaughtEarlyBuyer event raised by the Coindistribution contract.
type CoindistributionCaughtEarlyBuyer struct {
	Sniper common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCaughtEarlyBuyer is a free log retrieval operation binding the contract event 0x55678e47d0a699d3ab99b0184c4ff14f2246ba80522deb921aa0c8823578c44a.
//
// Solidity: event CaughtEarlyBuyer(address sniper)
func (_Coindistribution *CoindistributionFilterer) FilterCaughtEarlyBuyer(opts *bind.FilterOpts) (*CoindistributionCaughtEarlyBuyerIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "CaughtEarlyBuyer")
	if err != nil {
		return nil, err
	}
	return &CoindistributionCaughtEarlyBuyerIterator{contract: _Coindistribution.contract, event: "CaughtEarlyBuyer", logs: logs, sub: sub}, nil
}

// WatchCaughtEarlyBuyer is a free log subscription operation binding the contract event 0x55678e47d0a699d3ab99b0184c4ff14f2246ba80522deb921aa0c8823578c44a.
//
// Solidity: event CaughtEarlyBuyer(address sniper)
func (_Coindistribution *CoindistributionFilterer) WatchCaughtEarlyBuyer(opts *bind.WatchOpts, sink chan<- *CoindistributionCaughtEarlyBuyer) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "CaughtEarlyBuyer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionCaughtEarlyBuyer)
				if err := _Coindistribution.contract.UnpackLog(event, "CaughtEarlyBuyer", log); err != nil {
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

// ParseCaughtEarlyBuyer is a log parse operation binding the contract event 0x55678e47d0a699d3ab99b0184c4ff14f2246ba80522deb921aa0c8823578c44a.
//
// Solidity: event CaughtEarlyBuyer(address sniper)
func (_Coindistribution *CoindistributionFilterer) ParseCaughtEarlyBuyer(log types.Log) (*CoindistributionCaughtEarlyBuyer, error) {
	event := new(CoindistributionCaughtEarlyBuyer)
	if err := _Coindistribution.contract.UnpackLog(event, "CaughtEarlyBuyer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionEnableTradingIterator is returned from FilterEnableTrading and is used to iterate over the raw logs and unpacked data for EnableTrading events raised by the Coindistribution contract.
type CoindistributionEnableTradingIterator struct {
	Event *CoindistributionEnableTrading // Event containing the contract specifics and raw log

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
func (it *CoindistributionEnableTradingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionEnableTrading)
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
		it.Event = new(CoindistributionEnableTrading)
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
func (it *CoindistributionEnableTradingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionEnableTradingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionEnableTrading represents a EnableTrading event raised by the Coindistribution contract.
type CoindistributionEnableTrading struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEnableTrading is a free log retrieval operation binding the contract event 0x1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d.
//
// Solidity: event EnableTrading()
func (_Coindistribution *CoindistributionFilterer) FilterEnableTrading(opts *bind.FilterOpts) (*CoindistributionEnableTradingIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "EnableTrading")
	if err != nil {
		return nil, err
	}
	return &CoindistributionEnableTradingIterator{contract: _Coindistribution.contract, event: "EnableTrading", logs: logs, sub: sub}, nil
}

// WatchEnableTrading is a free log subscription operation binding the contract event 0x1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d.
//
// Solidity: event EnableTrading()
func (_Coindistribution *CoindistributionFilterer) WatchEnableTrading(opts *bind.WatchOpts, sink chan<- *CoindistributionEnableTrading) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "EnableTrading")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionEnableTrading)
				if err := _Coindistribution.contract.UnpackLog(event, "EnableTrading", log); err != nil {
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

// ParseEnableTrading is a log parse operation binding the contract event 0x1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d.
//
// Solidity: event EnableTrading()
func (_Coindistribution *CoindistributionFilterer) ParseEnableTrading(log types.Log) (*CoindistributionEnableTrading, error) {
	event := new(CoindistributionEnableTrading)
	if err := _Coindistribution.contract.UnpackLog(event, "EnableTrading", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionExcludeFromFeesIterator is returned from FilterExcludeFromFees and is used to iterate over the raw logs and unpacked data for ExcludeFromFees events raised by the Coindistribution contract.
type CoindistributionExcludeFromFeesIterator struct {
	Event *CoindistributionExcludeFromFees // Event containing the contract specifics and raw log

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
func (it *CoindistributionExcludeFromFeesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionExcludeFromFees)
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
		it.Event = new(CoindistributionExcludeFromFees)
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
func (it *CoindistributionExcludeFromFeesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionExcludeFromFeesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionExcludeFromFees represents a ExcludeFromFees event raised by the Coindistribution contract.
type CoindistributionExcludeFromFees struct {
	Account    common.Address
	IsExcluded bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterExcludeFromFees is a free log retrieval operation binding the contract event 0x9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df7.
//
// Solidity: event ExcludeFromFees(address indexed account, bool isExcluded)
func (_Coindistribution *CoindistributionFilterer) FilterExcludeFromFees(opts *bind.FilterOpts, account []common.Address) (*CoindistributionExcludeFromFeesIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "ExcludeFromFees", accountRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionExcludeFromFeesIterator{contract: _Coindistribution.contract, event: "ExcludeFromFees", logs: logs, sub: sub}, nil
}

// WatchExcludeFromFees is a free log subscription operation binding the contract event 0x9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df7.
//
// Solidity: event ExcludeFromFees(address indexed account, bool isExcluded)
func (_Coindistribution *CoindistributionFilterer) WatchExcludeFromFees(opts *bind.WatchOpts, sink chan<- *CoindistributionExcludeFromFees, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "ExcludeFromFees", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionExcludeFromFees)
				if err := _Coindistribution.contract.UnpackLog(event, "ExcludeFromFees", log); err != nil {
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

// ParseExcludeFromFees is a log parse operation binding the contract event 0x9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df7.
//
// Solidity: event ExcludeFromFees(address indexed account, bool isExcluded)
func (_Coindistribution *CoindistributionFilterer) ParseExcludeFromFees(log types.Log) (*CoindistributionExcludeFromFees, error) {
	event := new(CoindistributionExcludeFromFees)
	if err := _Coindistribution.contract.UnpackLog(event, "ExcludeFromFees", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionMaxTransactionExclusionIterator is returned from FilterMaxTransactionExclusion and is used to iterate over the raw logs and unpacked data for MaxTransactionExclusion events raised by the Coindistribution contract.
type CoindistributionMaxTransactionExclusionIterator struct {
	Event *CoindistributionMaxTransactionExclusion // Event containing the contract specifics and raw log

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
func (it *CoindistributionMaxTransactionExclusionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionMaxTransactionExclusion)
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
		it.Event = new(CoindistributionMaxTransactionExclusion)
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
func (it *CoindistributionMaxTransactionExclusionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionMaxTransactionExclusionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionMaxTransactionExclusion represents a MaxTransactionExclusion event raised by the Coindistribution contract.
type CoindistributionMaxTransactionExclusion struct {
	Address  common.Address
	Excluded bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMaxTransactionExclusion is a free log retrieval operation binding the contract event 0x6b4f1be9103e6cbcd38ca4a922334f2c3109b260130a6676a987f94088fd6746.
//
// Solidity: event MaxTransactionExclusion(address _address, bool excluded)
func (_Coindistribution *CoindistributionFilterer) FilterMaxTransactionExclusion(opts *bind.FilterOpts) (*CoindistributionMaxTransactionExclusionIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "MaxTransactionExclusion")
	if err != nil {
		return nil, err
	}
	return &CoindistributionMaxTransactionExclusionIterator{contract: _Coindistribution.contract, event: "MaxTransactionExclusion", logs: logs, sub: sub}, nil
}

// WatchMaxTransactionExclusion is a free log subscription operation binding the contract event 0x6b4f1be9103e6cbcd38ca4a922334f2c3109b260130a6676a987f94088fd6746.
//
// Solidity: event MaxTransactionExclusion(address _address, bool excluded)
func (_Coindistribution *CoindistributionFilterer) WatchMaxTransactionExclusion(opts *bind.WatchOpts, sink chan<- *CoindistributionMaxTransactionExclusion) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "MaxTransactionExclusion")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionMaxTransactionExclusion)
				if err := _Coindistribution.contract.UnpackLog(event, "MaxTransactionExclusion", log); err != nil {
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

// ParseMaxTransactionExclusion is a log parse operation binding the contract event 0x6b4f1be9103e6cbcd38ca4a922334f2c3109b260130a6676a987f94088fd6746.
//
// Solidity: event MaxTransactionExclusion(address _address, bool excluded)
func (_Coindistribution *CoindistributionFilterer) ParseMaxTransactionExclusion(log types.Log) (*CoindistributionMaxTransactionExclusion, error) {
	event := new(CoindistributionMaxTransactionExclusion)
	if err := _Coindistribution.contract.UnpackLog(event, "MaxTransactionExclusion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionOnSetUniswapRouterIterator is returned from FilterOnSetUniswapRouter and is used to iterate over the raw logs and unpacked data for OnSetUniswapRouter events raised by the Coindistribution contract.
type CoindistributionOnSetUniswapRouterIterator struct {
	Event *CoindistributionOnSetUniswapRouter // Event containing the contract specifics and raw log

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
func (it *CoindistributionOnSetUniswapRouterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionOnSetUniswapRouter)
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
		it.Event = new(CoindistributionOnSetUniswapRouter)
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
func (it *CoindistributionOnSetUniswapRouterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionOnSetUniswapRouterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionOnSetUniswapRouter represents a OnSetUniswapRouter event raised by the Coindistribution contract.
type CoindistributionOnSetUniswapRouter struct {
	TheAddress common.Address
	Flag       bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOnSetUniswapRouter is a free log retrieval operation binding the contract event 0x3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c.
//
// Solidity: event OnSetUniswapRouter(address indexed theAddress, bool indexed flag)
func (_Coindistribution *CoindistributionFilterer) FilterOnSetUniswapRouter(opts *bind.FilterOpts, theAddress []common.Address, flag []bool) (*CoindistributionOnSetUniswapRouterIterator, error) {

	var theAddressRule []interface{}
	for _, theAddressItem := range theAddress {
		theAddressRule = append(theAddressRule, theAddressItem)
	}
	var flagRule []interface{}
	for _, flagItem := range flag {
		flagRule = append(flagRule, flagItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "OnSetUniswapRouter", theAddressRule, flagRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionOnSetUniswapRouterIterator{contract: _Coindistribution.contract, event: "OnSetUniswapRouter", logs: logs, sub: sub}, nil
}

// WatchOnSetUniswapRouter is a free log subscription operation binding the contract event 0x3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c.
//
// Solidity: event OnSetUniswapRouter(address indexed theAddress, bool indexed flag)
func (_Coindistribution *CoindistributionFilterer) WatchOnSetUniswapRouter(opts *bind.WatchOpts, sink chan<- *CoindistributionOnSetUniswapRouter, theAddress []common.Address, flag []bool) (event.Subscription, error) {

	var theAddressRule []interface{}
	for _, theAddressItem := range theAddress {
		theAddressRule = append(theAddressRule, theAddressItem)
	}
	var flagRule []interface{}
	for _, flagItem := range flag {
		flagRule = append(flagRule, flagItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "OnSetUniswapRouter", theAddressRule, flagRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionOnSetUniswapRouter)
				if err := _Coindistribution.contract.UnpackLog(event, "OnSetUniswapRouter", log); err != nil {
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

// ParseOnSetUniswapRouter is a log parse operation binding the contract event 0x3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c.
//
// Solidity: event OnSetUniswapRouter(address indexed theAddress, bool indexed flag)
func (_Coindistribution *CoindistributionFilterer) ParseOnSetUniswapRouter(log types.Log) (*CoindistributionOnSetUniswapRouter, error) {
	event := new(CoindistributionOnSetUniswapRouter)
	if err := _Coindistribution.contract.UnpackLog(event, "OnSetUniswapRouter", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionOwnerForcedSwapBackIterator is returned from FilterOwnerForcedSwapBack and is used to iterate over the raw logs and unpacked data for OwnerForcedSwapBack events raised by the Coindistribution contract.
type CoindistributionOwnerForcedSwapBackIterator struct {
	Event *CoindistributionOwnerForcedSwapBack // Event containing the contract specifics and raw log

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
func (it *CoindistributionOwnerForcedSwapBackIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionOwnerForcedSwapBack)
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
		it.Event = new(CoindistributionOwnerForcedSwapBack)
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
func (it *CoindistributionOwnerForcedSwapBackIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionOwnerForcedSwapBackIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionOwnerForcedSwapBack represents a OwnerForcedSwapBack event raised by the Coindistribution contract.
type CoindistributionOwnerForcedSwapBack struct {
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnerForcedSwapBack is a free log retrieval operation binding the contract event 0x1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb32.
//
// Solidity: event OwnerForcedSwapBack(uint256 timestamp)
func (_Coindistribution *CoindistributionFilterer) FilterOwnerForcedSwapBack(opts *bind.FilterOpts) (*CoindistributionOwnerForcedSwapBackIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "OwnerForcedSwapBack")
	if err != nil {
		return nil, err
	}
	return &CoindistributionOwnerForcedSwapBackIterator{contract: _Coindistribution.contract, event: "OwnerForcedSwapBack", logs: logs, sub: sub}, nil
}

// WatchOwnerForcedSwapBack is a free log subscription operation binding the contract event 0x1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb32.
//
// Solidity: event OwnerForcedSwapBack(uint256 timestamp)
func (_Coindistribution *CoindistributionFilterer) WatchOwnerForcedSwapBack(opts *bind.WatchOpts, sink chan<- *CoindistributionOwnerForcedSwapBack) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "OwnerForcedSwapBack")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionOwnerForcedSwapBack)
				if err := _Coindistribution.contract.UnpackLog(event, "OwnerForcedSwapBack", log); err != nil {
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

// ParseOwnerForcedSwapBack is a log parse operation binding the contract event 0x1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb32.
//
// Solidity: event OwnerForcedSwapBack(uint256 timestamp)
func (_Coindistribution *CoindistributionFilterer) ParseOwnerForcedSwapBack(log types.Log) (*CoindistributionOwnerForcedSwapBack, error) {
	event := new(CoindistributionOwnerForcedSwapBack)
	if err := _Coindistribution.contract.UnpackLog(event, "OwnerForcedSwapBack", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Coindistribution contract.
type CoindistributionOwnershipTransferredIterator struct {
	Event *CoindistributionOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CoindistributionOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionOwnershipTransferred)
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
		it.Event = new(CoindistributionOwnershipTransferred)
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
func (it *CoindistributionOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionOwnershipTransferred represents a OwnershipTransferred event raised by the Coindistribution contract.
type CoindistributionOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Coindistribution *CoindistributionFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CoindistributionOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionOwnershipTransferredIterator{contract: _Coindistribution.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Coindistribution *CoindistributionFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CoindistributionOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionOwnershipTransferred)
				if err := _Coindistribution.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Coindistribution *CoindistributionFilterer) ParseOwnershipTransferred(log types.Log) (*CoindistributionOwnershipTransferred, error) {
	event := new(CoindistributionOwnershipTransferred)
	if err := _Coindistribution.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionRemovedLimitsIterator is returned from FilterRemovedLimits and is used to iterate over the raw logs and unpacked data for RemovedLimits events raised by the Coindistribution contract.
type CoindistributionRemovedLimitsIterator struct {
	Event *CoindistributionRemovedLimits // Event containing the contract specifics and raw log

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
func (it *CoindistributionRemovedLimitsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionRemovedLimits)
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
		it.Event = new(CoindistributionRemovedLimits)
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
func (it *CoindistributionRemovedLimitsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionRemovedLimitsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionRemovedLimits represents a RemovedLimits event raised by the Coindistribution contract.
type CoindistributionRemovedLimits struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRemovedLimits is a free log retrieval operation binding the contract event 0xa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c.
//
// Solidity: event RemovedLimits()
func (_Coindistribution *CoindistributionFilterer) FilterRemovedLimits(opts *bind.FilterOpts) (*CoindistributionRemovedLimitsIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "RemovedLimits")
	if err != nil {
		return nil, err
	}
	return &CoindistributionRemovedLimitsIterator{contract: _Coindistribution.contract, event: "RemovedLimits", logs: logs, sub: sub}, nil
}

// WatchRemovedLimits is a free log subscription operation binding the contract event 0xa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c.
//
// Solidity: event RemovedLimits()
func (_Coindistribution *CoindistributionFilterer) WatchRemovedLimits(opts *bind.WatchOpts, sink chan<- *CoindistributionRemovedLimits) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "RemovedLimits")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionRemovedLimits)
				if err := _Coindistribution.contract.UnpackLog(event, "RemovedLimits", log); err != nil {
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

// ParseRemovedLimits is a log parse operation binding the contract event 0xa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c.
//
// Solidity: event RemovedLimits()
func (_Coindistribution *CoindistributionFilterer) ParseRemovedLimits(log types.Log) (*CoindistributionRemovedLimits, error) {
	event := new(CoindistributionRemovedLimits)
	if err := _Coindistribution.contract.UnpackLog(event, "RemovedLimits", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionSetUniswapV2LiquidityPoolIterator is returned from FilterSetUniswapV2LiquidityPool and is used to iterate over the raw logs and unpacked data for SetUniswapV2LiquidityPool events raised by the Coindistribution contract.
type CoindistributionSetUniswapV2LiquidityPoolIterator struct {
	Event *CoindistributionSetUniswapV2LiquidityPool // Event containing the contract specifics and raw log

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
func (it *CoindistributionSetUniswapV2LiquidityPoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionSetUniswapV2LiquidityPool)
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
		it.Event = new(CoindistributionSetUniswapV2LiquidityPool)
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
func (it *CoindistributionSetUniswapV2LiquidityPoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionSetUniswapV2LiquidityPoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionSetUniswapV2LiquidityPool represents a SetUniswapV2LiquidityPool event raised by the Coindistribution contract.
type CoindistributionSetUniswapV2LiquidityPool struct {
	Router   common.Address
	Pair     common.Address
	Enabled  bool
	Receiver bool
	BLiqF    *big.Int
	SLiqF    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetUniswapV2LiquidityPool is a free log retrieval operation binding the contract event 0x78c6052d21b64dbee719270b6a9e56ba166f8f88fbaa39a44abf4a11642bae92.
//
// Solidity: event SetUniswapV2LiquidityPool(address indexed router, address indexed pair, bool enabled, bool receiver, uint256 bLiqF, uint256 sLiqF)
func (_Coindistribution *CoindistributionFilterer) FilterSetUniswapV2LiquidityPool(opts *bind.FilterOpts, router []common.Address, pair []common.Address) (*CoindistributionSetUniswapV2LiquidityPoolIterator, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}
	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "SetUniswapV2LiquidityPool", routerRule, pairRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionSetUniswapV2LiquidityPoolIterator{contract: _Coindistribution.contract, event: "SetUniswapV2LiquidityPool", logs: logs, sub: sub}, nil
}

// WatchSetUniswapV2LiquidityPool is a free log subscription operation binding the contract event 0x78c6052d21b64dbee719270b6a9e56ba166f8f88fbaa39a44abf4a11642bae92.
//
// Solidity: event SetUniswapV2LiquidityPool(address indexed router, address indexed pair, bool enabled, bool receiver, uint256 bLiqF, uint256 sLiqF)
func (_Coindistribution *CoindistributionFilterer) WatchSetUniswapV2LiquidityPool(opts *bind.WatchOpts, sink chan<- *CoindistributionSetUniswapV2LiquidityPool, router []common.Address, pair []common.Address) (event.Subscription, error) {

	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}
	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "SetUniswapV2LiquidityPool", routerRule, pairRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionSetUniswapV2LiquidityPool)
				if err := _Coindistribution.contract.UnpackLog(event, "SetUniswapV2LiquidityPool", log); err != nil {
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

// ParseSetUniswapV2LiquidityPool is a log parse operation binding the contract event 0x78c6052d21b64dbee719270b6a9e56ba166f8f88fbaa39a44abf4a11642bae92.
//
// Solidity: event SetUniswapV2LiquidityPool(address indexed router, address indexed pair, bool enabled, bool receiver, uint256 bLiqF, uint256 sLiqF)
func (_Coindistribution *CoindistributionFilterer) ParseSetUniswapV2LiquidityPool(log types.Log) (*CoindistributionSetUniswapV2LiquidityPool, error) {
	event := new(CoindistributionSetUniswapV2LiquidityPool)
	if err := _Coindistribution.contract.UnpackLog(event, "SetUniswapV2LiquidityPool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionSwapAndLiquifyIterator is returned from FilterSwapAndLiquify and is used to iterate over the raw logs and unpacked data for SwapAndLiquify events raised by the Coindistribution contract.
type CoindistributionSwapAndLiquifyIterator struct {
	Event *CoindistributionSwapAndLiquify // Event containing the contract specifics and raw log

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
func (it *CoindistributionSwapAndLiquifyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionSwapAndLiquify)
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
		it.Event = new(CoindistributionSwapAndLiquify)
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
func (it *CoindistributionSwapAndLiquifyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionSwapAndLiquifyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionSwapAndLiquify represents a SwapAndLiquify event raised by the Coindistribution contract.
type CoindistributionSwapAndLiquify struct {
	TokensSwapped       *big.Int
	EthReceived         *big.Int
	TokensIntoLiquidity *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSwapAndLiquify is a free log retrieval operation binding the contract event 0x17bbfb9a6069321b6ded73bd96327c9e6b7212a5cd51ff219cd61370acafb561.
//
// Solidity: event SwapAndLiquify(uint256 tokensSwapped, uint256 ethReceived, uint256 tokensIntoLiquidity)
func (_Coindistribution *CoindistributionFilterer) FilterSwapAndLiquify(opts *bind.FilterOpts) (*CoindistributionSwapAndLiquifyIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "SwapAndLiquify")
	if err != nil {
		return nil, err
	}
	return &CoindistributionSwapAndLiquifyIterator{contract: _Coindistribution.contract, event: "SwapAndLiquify", logs: logs, sub: sub}, nil
}

// WatchSwapAndLiquify is a free log subscription operation binding the contract event 0x17bbfb9a6069321b6ded73bd96327c9e6b7212a5cd51ff219cd61370acafb561.
//
// Solidity: event SwapAndLiquify(uint256 tokensSwapped, uint256 ethReceived, uint256 tokensIntoLiquidity)
func (_Coindistribution *CoindistributionFilterer) WatchSwapAndLiquify(opts *bind.WatchOpts, sink chan<- *CoindistributionSwapAndLiquify) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "SwapAndLiquify")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionSwapAndLiquify)
				if err := _Coindistribution.contract.UnpackLog(event, "SwapAndLiquify", log); err != nil {
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

// ParseSwapAndLiquify is a log parse operation binding the contract event 0x17bbfb9a6069321b6ded73bd96327c9e6b7212a5cd51ff219cd61370acafb561.
//
// Solidity: event SwapAndLiquify(uint256 tokensSwapped, uint256 ethReceived, uint256 tokensIntoLiquidity)
func (_Coindistribution *CoindistributionFilterer) ParseSwapAndLiquify(log types.Log) (*CoindistributionSwapAndLiquify, error) {
	event := new(CoindistributionSwapAndLiquify)
	if err := _Coindistribution.contract.UnpackLog(event, "SwapAndLiquify", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Coindistribution contract.
type CoindistributionTransferIterator struct {
	Event *CoindistributionTransfer // Event containing the contract specifics and raw log

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
func (it *CoindistributionTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionTransfer)
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
		it.Event = new(CoindistributionTransfer)
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
func (it *CoindistributionTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionTransfer represents a Transfer event raised by the Coindistribution contract.
type CoindistributionTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Coindistribution *CoindistributionFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CoindistributionTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionTransferIterator{contract: _Coindistribution.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Coindistribution *CoindistributionFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CoindistributionTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionTransfer)
				if err := _Coindistribution.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Coindistribution *CoindistributionFilterer) ParseTransfer(log types.Log) (*CoindistributionTransfer, error) {
	event := new(CoindistributionTransfer)
	if err := _Coindistribution.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionTransferForeignTokenIterator is returned from FilterTransferForeignToken and is used to iterate over the raw logs and unpacked data for TransferForeignToken events raised by the Coindistribution contract.
type CoindistributionTransferForeignTokenIterator struct {
	Event *CoindistributionTransferForeignToken // Event containing the contract specifics and raw log

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
func (it *CoindistributionTransferForeignTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionTransferForeignToken)
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
		it.Event = new(CoindistributionTransferForeignToken)
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
func (it *CoindistributionTransferForeignTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionTransferForeignTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionTransferForeignToken represents a TransferForeignToken event raised by the Coindistribution contract.
type CoindistributionTransferForeignToken struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferForeignToken is a free log retrieval operation binding the contract event 0x5661684995ab94d684bfe57a43c4141578f52d3e7374e8cd3250e2f062e13ac1.
//
// Solidity: event TransferForeignToken(address token, address to, uint256 amount)
func (_Coindistribution *CoindistributionFilterer) FilterTransferForeignToken(opts *bind.FilterOpts) (*CoindistributionTransferForeignTokenIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "TransferForeignToken")
	if err != nil {
		return nil, err
	}
	return &CoindistributionTransferForeignTokenIterator{contract: _Coindistribution.contract, event: "TransferForeignToken", logs: logs, sub: sub}, nil
}

// WatchTransferForeignToken is a free log subscription operation binding the contract event 0x5661684995ab94d684bfe57a43c4141578f52d3e7374e8cd3250e2f062e13ac1.
//
// Solidity: event TransferForeignToken(address token, address to, uint256 amount)
func (_Coindistribution *CoindistributionFilterer) WatchTransferForeignToken(opts *bind.WatchOpts, sink chan<- *CoindistributionTransferForeignToken) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "TransferForeignToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionTransferForeignToken)
				if err := _Coindistribution.contract.UnpackLog(event, "TransferForeignToken", log); err != nil {
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

// ParseTransferForeignToken is a log parse operation binding the contract event 0x5661684995ab94d684bfe57a43c4141578f52d3e7374e8cd3250e2f062e13ac1.
//
// Solidity: event TransferForeignToken(address token, address to, uint256 amount)
func (_Coindistribution *CoindistributionFilterer) ParseTransferForeignToken(log types.Log) (*CoindistributionTransferForeignToken, error) {
	event := new(CoindistributionTransferForeignToken)
	if err := _Coindistribution.contract.UnpackLog(event, "TransferForeignToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionUpdatedMaxBuyAmountIterator is returned from FilterUpdatedMaxBuyAmount and is used to iterate over the raw logs and unpacked data for UpdatedMaxBuyAmount events raised by the Coindistribution contract.
type CoindistributionUpdatedMaxBuyAmountIterator struct {
	Event *CoindistributionUpdatedMaxBuyAmount // Event containing the contract specifics and raw log

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
func (it *CoindistributionUpdatedMaxBuyAmountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionUpdatedMaxBuyAmount)
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
		it.Event = new(CoindistributionUpdatedMaxBuyAmount)
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
func (it *CoindistributionUpdatedMaxBuyAmountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionUpdatedMaxBuyAmountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionUpdatedMaxBuyAmount represents a UpdatedMaxBuyAmount event raised by the Coindistribution contract.
type CoindistributionUpdatedMaxBuyAmount struct {
	NewAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUpdatedMaxBuyAmount is a free log retrieval operation binding the contract event 0xfcc0366804aaa8dbf88a2924100c733b70dec8445957a5d5f8ff92898de41009.
//
// Solidity: event UpdatedMaxBuyAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) FilterUpdatedMaxBuyAmount(opts *bind.FilterOpts) (*CoindistributionUpdatedMaxBuyAmountIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "UpdatedMaxBuyAmount")
	if err != nil {
		return nil, err
	}
	return &CoindistributionUpdatedMaxBuyAmountIterator{contract: _Coindistribution.contract, event: "UpdatedMaxBuyAmount", logs: logs, sub: sub}, nil
}

// WatchUpdatedMaxBuyAmount is a free log subscription operation binding the contract event 0xfcc0366804aaa8dbf88a2924100c733b70dec8445957a5d5f8ff92898de41009.
//
// Solidity: event UpdatedMaxBuyAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) WatchUpdatedMaxBuyAmount(opts *bind.WatchOpts, sink chan<- *CoindistributionUpdatedMaxBuyAmount) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "UpdatedMaxBuyAmount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionUpdatedMaxBuyAmount)
				if err := _Coindistribution.contract.UnpackLog(event, "UpdatedMaxBuyAmount", log); err != nil {
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

// ParseUpdatedMaxBuyAmount is a log parse operation binding the contract event 0xfcc0366804aaa8dbf88a2924100c733b70dec8445957a5d5f8ff92898de41009.
//
// Solidity: event UpdatedMaxBuyAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) ParseUpdatedMaxBuyAmount(log types.Log) (*CoindistributionUpdatedMaxBuyAmount, error) {
	event := new(CoindistributionUpdatedMaxBuyAmount)
	if err := _Coindistribution.contract.UnpackLog(event, "UpdatedMaxBuyAmount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionUpdatedMaxSellAmountIterator is returned from FilterUpdatedMaxSellAmount and is used to iterate over the raw logs and unpacked data for UpdatedMaxSellAmount events raised by the Coindistribution contract.
type CoindistributionUpdatedMaxSellAmountIterator struct {
	Event *CoindistributionUpdatedMaxSellAmount // Event containing the contract specifics and raw log

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
func (it *CoindistributionUpdatedMaxSellAmountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionUpdatedMaxSellAmount)
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
		it.Event = new(CoindistributionUpdatedMaxSellAmount)
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
func (it *CoindistributionUpdatedMaxSellAmountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionUpdatedMaxSellAmountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionUpdatedMaxSellAmount represents a UpdatedMaxSellAmount event raised by the Coindistribution contract.
type CoindistributionUpdatedMaxSellAmount struct {
	NewAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUpdatedMaxSellAmount is a free log retrieval operation binding the contract event 0x53c4eb831d8cfeb750f1c62590d8cd30f4c6f0380d29a05caa09f0d92588560e.
//
// Solidity: event UpdatedMaxSellAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) FilterUpdatedMaxSellAmount(opts *bind.FilterOpts) (*CoindistributionUpdatedMaxSellAmountIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "UpdatedMaxSellAmount")
	if err != nil {
		return nil, err
	}
	return &CoindistributionUpdatedMaxSellAmountIterator{contract: _Coindistribution.contract, event: "UpdatedMaxSellAmount", logs: logs, sub: sub}, nil
}

// WatchUpdatedMaxSellAmount is a free log subscription operation binding the contract event 0x53c4eb831d8cfeb750f1c62590d8cd30f4c6f0380d29a05caa09f0d92588560e.
//
// Solidity: event UpdatedMaxSellAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) WatchUpdatedMaxSellAmount(opts *bind.WatchOpts, sink chan<- *CoindistributionUpdatedMaxSellAmount) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "UpdatedMaxSellAmount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionUpdatedMaxSellAmount)
				if err := _Coindistribution.contract.UnpackLog(event, "UpdatedMaxSellAmount", log); err != nil {
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

// ParseUpdatedMaxSellAmount is a log parse operation binding the contract event 0x53c4eb831d8cfeb750f1c62590d8cd30f4c6f0380d29a05caa09f0d92588560e.
//
// Solidity: event UpdatedMaxSellAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) ParseUpdatedMaxSellAmount(log types.Log) (*CoindistributionUpdatedMaxSellAmount, error) {
	event := new(CoindistributionUpdatedMaxSellAmount)
	if err := _Coindistribution.contract.UnpackLog(event, "UpdatedMaxSellAmount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CoindistributionUpdatedMaxWalletAmountIterator is returned from FilterUpdatedMaxWalletAmount and is used to iterate over the raw logs and unpacked data for UpdatedMaxWalletAmount events raised by the Coindistribution contract.
type CoindistributionUpdatedMaxWalletAmountIterator struct {
	Event *CoindistributionUpdatedMaxWalletAmount // Event containing the contract specifics and raw log

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
func (it *CoindistributionUpdatedMaxWalletAmountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionUpdatedMaxWalletAmount)
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
		it.Event = new(CoindistributionUpdatedMaxWalletAmount)
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
func (it *CoindistributionUpdatedMaxWalletAmountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionUpdatedMaxWalletAmountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionUpdatedMaxWalletAmount represents a UpdatedMaxWalletAmount event raised by the Coindistribution contract.
type CoindistributionUpdatedMaxWalletAmount struct {
	NewAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUpdatedMaxWalletAmount is a free log retrieval operation binding the contract event 0xefc9add9a9b7382de284ef5ad69d8ea863e2680492b21a81948c2d5f04a442bc.
//
// Solidity: event UpdatedMaxWalletAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) FilterUpdatedMaxWalletAmount(opts *bind.FilterOpts) (*CoindistributionUpdatedMaxWalletAmountIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "UpdatedMaxWalletAmount")
	if err != nil {
		return nil, err
	}
	return &CoindistributionUpdatedMaxWalletAmountIterator{contract: _Coindistribution.contract, event: "UpdatedMaxWalletAmount", logs: logs, sub: sub}, nil
}

// WatchUpdatedMaxWalletAmount is a free log subscription operation binding the contract event 0xefc9add9a9b7382de284ef5ad69d8ea863e2680492b21a81948c2d5f04a442bc.
//
// Solidity: event UpdatedMaxWalletAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) WatchUpdatedMaxWalletAmount(opts *bind.WatchOpts, sink chan<- *CoindistributionUpdatedMaxWalletAmount) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "UpdatedMaxWalletAmount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionUpdatedMaxWalletAmount)
				if err := _Coindistribution.contract.UnpackLog(event, "UpdatedMaxWalletAmount", log); err != nil {
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

// ParseUpdatedMaxWalletAmount is a log parse operation binding the contract event 0xefc9add9a9b7382de284ef5ad69d8ea863e2680492b21a81948c2d5f04a442bc.
//
// Solidity: event UpdatedMaxWalletAmount(uint256 newAmount)
func (_Coindistribution *CoindistributionFilterer) ParseUpdatedMaxWalletAmount(log types.Log) (*CoindistributionUpdatedMaxWalletAmount, error) {
	event := new(CoindistributionUpdatedMaxWalletAmount)
	if err := _Coindistribution.contract.UnpackLog(event, "UpdatedMaxWalletAmount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
