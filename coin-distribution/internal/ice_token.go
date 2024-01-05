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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20ZeroToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeIncreased\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeMoreThan5\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForeignTokenSelfTransfer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Mismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoBots\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSwapBack\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAirDropper\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawStuckETH\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BuyBackTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sniper\",\"type\":\"address\"}],\"name\":\"CaughtEarlyBuyer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EnableTrading\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isExcluded\",\"type\":\"bool\"}],\"name\":\"ExcludeFromFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"MaxTransactionExclusion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"theAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"OnSetUniswapRouter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"OwnerForcedSwapBack\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RemovedLimits\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"}],\"name\":\"SetUniswapV2LiquidityPool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensSwapped\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethReceived\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensIntoLiquidity\",\"type\":\"uint256\"}],\"name\":\"SwapAndLiquify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferForeignToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxBuyAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxWalletAmount\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"airdropToWallets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"first\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"second\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeElapsed\",\"type\":\"uint256\"}],\"name\":\"applySlippage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"botsCaught\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableTransferDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableTrading\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"excludeFromFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceSwapBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAirDropper\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTradingEnabledBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"limitsInEffect\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"manageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"wallets\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"massManageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"first\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"second\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"}],\"name\":\"priceUpdateTimeElapsedForPair\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"removeLimits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAirDropper\",\"type\":\"address\"}],\"name\":\"setAirDropper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"setKnownUniswapRouters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setSwapBackThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"theAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"setUniswapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"}],\"name\":\"setUniswapV2LiquidityPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokensForLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transferDelayEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferForeignToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"uniswapV2LiquidityPoolSlots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uniswapV2LiquidityPools\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawStuckETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052683635c9adc5dea00000600b555f600c5f6101000a81548160ff0219169083151502179055506001600c60016101000a81548160ff0219169083151502179055506001600e5f6101000a81548160ff0219169083151502179055505f6013553480156200006f575f80fd5b50336040518060400160405280600381526020017f49636500000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f49434500000000000000000000000000000000000000000000000000000000008152508160029081620000ee91906200090d565b5080600390816200010091906200090d565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000176575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016200016d919062000a34565b60405180910390fd5b62000187816200020660201b60201c565b503360055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001db336001620002c960201b60201c565b620001ee306001620002c960201b60201c565b6200020060016200038160201b60201c565b62000a86565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b620002d9620004c460201b60201c565b8060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df78260405162000375919062000a6b565b60405180910390a25050565b62000391620004c460201b60201c565b620003b7737a250d5630b4cf539739df2c5dacb4c659f2488d826200056660201b60201c565b620003dd73e592427a0aece92de3edee1f18e0157c05861564826200056660201b60201c565b620004037368b3465833fb72a70ecdf485e0e4c7bd8665fc45826200056660201b60201c565b6200042973eff92a263d31888d860bd50809a8d171709b7b1c826200056660201b60201c565b6200044f731b81d678ffb9c0263b24a97847620c99d213eb14826200056660201b60201c565b620004757313f4ea83d0bd40e75c8222255bc855a974568dd4826200056660201b60201c565b6200049b731b02da8cb0d097eb8d57a175b88c7d8b47997506826200056660201b60201c565b620004c173d9e1ce17f2641f24ae83637ab66a2cca9c378b9f826200056660201b60201c565b50565b620004d46200067a60201b60201c565b73ffffffffffffffffffffffffffffffffffffffff16620004fa6200068160201b60201c565b73ffffffffffffffffffffffffffffffffffffffff16146200056457620005266200067a60201b60201c565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016200055b919062000a34565b60405180910390fd5b565b62000576620004c460201b60201c565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620005dc576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060145f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508015158273ffffffffffffffffffffffffffffffffffffffff167f3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c60405160405180910390a35050565b5f33905090565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200072557607f821691505b6020821081036200073b576200073a620006e0565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026200079f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000762565b620007ab868362000762565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620007f5620007ef620007e984620007c3565b620007cc565b620007c3565b9050919050565b5f819050919050565b6200081083620007d5565b620008286200081f82620007fc565b8484546200076e565b825550505050565b5f90565b6200083e62000830565b6200084b81848462000805565b505050565b5b818110156200087257620008665f8262000834565b60018101905062000851565b5050565b601f821115620008c1576200088b8162000741565b620008968462000753565b81016020851015620008a6578190505b620008be620008b58562000753565b83018262000850565b50505b505050565b5f82821c905092915050565b5f620008e35f1984600802620008c6565b1980831691505092915050565b5f620008fd8383620008d2565b9150826002028217905092915050565b6200091882620006a9565b67ffffffffffffffff811115620009345762000933620006b3565b5b6200094082546200070d565b6200094d82828562000876565b5f60209050601f83116001811462000983575f84156200096e578287015190505b6200097a8582620008f0565b865550620009e9565b601f198416620009938662000741565b5f5b82811015620009bc5784890151825560018201915060208501945060208101905062000995565b86831015620009dc5784890151620009d8601f891682620008d2565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f62000a1c82620009f1565b9050919050565b62000a2e8162000a10565b82525050565b5f60208201905062000a495f83018462000a23565b92915050565b5f8115159050919050565b62000a658162000a4f565b82525050565b5f60208201905062000a805f83018462000a5a565b92915050565b614d338062000a945f395ff3fe608060405260043610610233575f3560e01c806370a082311161012d578063a9059cbb116100aa578063e800dff71161006e578063e800dff714610828578063e884f26014610850578063f2fde38b14610866578063f317f2c31461088e578063f5648a4f146108b65761023a565b8063a9059cbb14610722578063bfd792841461075e578063c02466681461079a578063c876d0b9146107c2578063dd62ed3e146107ec5761023a565b80638366e79a116100f15780638366e79a146106685780638a8c523c146106905780638da5cb5b146106a65780638fe62b8a146106d057806395d89b41146106f85761023a565b806370a0823114610588578063715018a6146105c45780637371dfff146105da578063751039fc14610616578063785cca3e1461062c5761023a565b8063313ce567116101bb578063588b655f1161017f578063588b655f146104b85780635db7548b146104e25780635de998ef1461050c5780636b0a894c146105345780636ddd17131461055e5761023a565b8063313ce567146103ea57806344d8a785146104145780634a62bb651461043c57806351f205e41461046657806352f93b5e1461047c5761023a565b806318160ddd1161020257806318160ddd146102f25780631a8145bb1461031c5780632307b4411461034657806323b872dd1461036e578063296a803c146103aa5761023a565b806306fdde031461023c578063072280c314610266578063095ea7b31461028e578063130a2c3c146102ca5761023a565b3661023a57005b005b348015610247575f80fd5b506102506108cc565b60405161025d9190613dfc565b60405180910390f35b348015610271575f80fd5b5061028c60048036038101906102879190613eb3565b61095c565b005b348015610299575f80fd5b506102b460048036038101906102af9190613f24565b610a67565b6040516102c19190613f71565b60405180910390f35b3480156102d5575f80fd5b506102f060048036038101906102eb9190613feb565b610a89565b005b3480156102fd575f80fd5b50610306610b32565b6040516103139190614057565b60405180910390f35b348015610327575f80fd5b50610330610b3b565b60405161033d9190614057565b60405180910390f35b348015610351575f80fd5b5061036c600480360381019061036791906140c5565b610b41565b005b348015610379575f80fd5b50610394600480360381019061038f9190614143565b610c73565b6040516103a19190613f71565b60405180910390f35b3480156103b5575f80fd5b506103d060048036038101906103cb9190614193565b610ca1565b6040516103e19594939291906141cd565b60405180910390f35b3480156103f5575f80fd5b506103fe610d29565b60405161040b9190614239565b60405180910390f35b34801561041f575f80fd5b5061043a60048036038101906104359190614252565b610d31565b005b348015610447575f80fd5b50610450610e2c565b60405161045d9190613f71565b60405180910390f35b348015610471575f80fd5b5061047a610e3f565b005b348015610487575f80fd5b506104a2600480360381019061049d919061427d565b610efc565b6040516104af9190614057565b60405180910390f35b3480156104c3575f80fd5b506104cc611075565b6040516104d99190614057565b60405180910390f35b3480156104ed575f80fd5b506104f661107e565b60405161050391906142cd565b60405180910390f35b348015610517575f80fd5b50610532600480360381019061052d9190614193565b6110a6565b005b34801561053f575f80fd5b506105486110f1565b6040516105559190614057565b60405180910390f35b348015610569575f80fd5b506105726110f7565b60405161057f9190613f71565b60405180910390f35b348015610593575f80fd5b506105ae60048036038101906105a99190614193565b611109565b6040516105bb9190614057565b60405180910390f35b3480156105cf575f80fd5b506105d861114f565b005b3480156105e5575f80fd5b5061060060048036038101906105fb91906142e6565b6111a1565b60405161060d9190614057565b60405180910390f35b348015610621575f80fd5b5061062a611662565b005b348015610637575f80fd5b50610652600480360381019061064d919061435d565b6116cb565b60405161065f91906142cd565b60405180910390f35b348015610673575f80fd5b5061068e60048036038101906106899190614388565b611706565b005b34801561069b575f80fd5b506106a461189e565b005b3480156106b1575f80fd5b506106ba611930565b6040516106c791906142cd565b60405180910390f35b3480156106db575f80fd5b506106f660048036038101906106f1919061435d565b611958565b005b348015610703575f80fd5b5061070c61196a565b6040516107199190613dfc565b60405180910390f35b34801561072d575f80fd5b5061074860048036038101906107439190613f24565b6119fa565b6040516107559190613f71565b60405180910390f35b348015610769575f80fd5b50610784600480360381019061077f9190614193565b611a1c565b6040516107919190613f71565b60405180910390f35b3480156107a5575f80fd5b506107c060048036038101906107bb9190613eb3565b611a39565b005b3480156107cd575f80fd5b506107d6611ae7565b6040516107e39190613f71565b60405180910390f35b3480156107f7575f80fd5b50610812600480360381019061080d9190614388565b611af9565b60405161081f9190614057565b60405180910390f35b348015610833575f80fd5b5061084e60048036038101906108499190613eb3565b611b7a565b005b34801561085b575f80fd5b50610864611bda565b005b348015610871575f80fd5b5061088c60048036038101906108879190614193565b611bfd565b005b348015610899575f80fd5b506108b460048036038101906108af91906143c6565b611c81565b005b3480156108c1575f80fd5b506108ca61213f565b005b6060600280546108db9061446a565b80601f01602080910402602001604051908101604052809291908181526020018280546109079061446a565b80156109525780601f1061092957610100808354040283529160200191610952565b820191905f5260205f20905b81548152906001019060200180831161093557829003601f168201915b5050505050905090565b610964612229565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036109c9576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060145f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508015158273ffffffffffffffffffffffffffffffffffffffff167f3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c60405160405180910390a35050565b5f80610a716122b0565b9050610a7e8185856122b7565b600191505092915050565b610a91612229565b5f5b83839050811015610b2c578160085f868685818110610ab557610ab461449a565b5b9050602002016020810190610aca9190614193565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508080610b24906144f4565b915050610a93565b50505050565b5f600754905090565b600f5481565b3373ffffffffffffffffffffffffffffffffffffffff1660055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610bc7576040517f16ad4feb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818190508484905014610c06576040517f77a93d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8484905090505f604051602081016004356024355f5b86811015610c5757602460208202018083013581850135875260068652604087208181540181558189019850505050600181019050610c1d565b50505050508060075f8282540192505081905550505050505050565b5f80610c7d6122b0565b9050610c8a8582856122c9565b610c9585858561235b565b60019150509392505050565b6010602052805f5260405f205f91509050805f015f9054906101000a900460ff1690805f0160019054906101000a900460ff1690805f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154905085565b5f6012905090565b610d39612229565b610d57737a250d5630b4cf539739df2c5dacb4c659f2488d8261095c565b610d7573e592427a0aece92de3edee1f18e0157c058615648261095c565b610d937368b3465833fb72a70ecdf485e0e4c7bd8665fc458261095c565b610db173eff92a263d31888d860bd50809a8d171709b7b1c8261095c565b610dcf731b81d678ffb9c0263b24a97847620c99d213eb148261095c565b610ded7313f4ea83d0bd40e75c8222255bc855a974568dd48261095c565b610e0b731b02da8cb0d097eb8d57a175b88c7d8b479975068261095c565b610e2973d9e1ce17f2641f24ae83637ab66a2cca9c378b9f8261095c565b50565b600c60019054906101000a900460ff1681565b610e47612229565b610e5030611109565b5f03610e88576040517f0b952c3600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600a5f6101000a81548160ff021916908315150217905550610eaa612d8a565b5f600a5f6101000a81548160ff0219169083151502179055507f1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb3242604051610ef29190614057565b60405180910390a1565b5f808273ffffffffffffffffffffffffffffffffffffffff1663c45a01556040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f47573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f6b919061454f565b73ffffffffffffffffffffffffffffffffffffffff1663e6a4390586866040518363ffffffff1660e01b8152600401610fa592919061457a565b602060405180830381865afa158015610fc0573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fe4919061454f565b90505f8173ffffffffffffffffffffffffffffffffffffffff16630902f1ac6040518163ffffffff1660e01b8152600401606060405180830381865afa158015611030573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611054919061461d565b925050508063ffffffff164261106a919061466d565b925050509392505050565b5f601354905090565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6110ae612229565b8060055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60095481565b600c5f9054906101000a900460ff1681565b5f60065f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b611157612229565b61115f613038565b5f60055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b5f808373ffffffffffffffffffffffffffffffffffffffff1663c45a01556040518163ffffffff1660e01b8152600401602060405180830381865afa1580156111ec573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611210919061454f565b73ffffffffffffffffffffffffffffffffffffffff1663e6a4390587876040518363ffffffff1660e01b815260040161124a92919061457a565b602060405180830381865afa158015611265573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611289919061454f565b90505f8173ffffffffffffffffffffffffffffffffffffffff16630dfe16816040518163ffffffff1660e01b8152600401602060405180830381865afa1580156112d5573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906112f9919061454f565b73ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff161461139d578173ffffffffffffffffffffffffffffffffffffffff16635a3d54936040518163ffffffff1660e01b8152600401602060405180830381865afa158015611374573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061139891906146b4565b61140b565b8173ffffffffffffffffffffffffffffffffffffffff16635909c0d56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156113e6573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061140a91906146b4565b5b90505f808373ffffffffffffffffffffffffffffffffffffffff16630902f1ac6040518163ffffffff1660e01b8152600401606060405180830381865afa158015611458573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061147c919061461d565b5091509150825f10156114b0575f8684611496919061470c565b9050808b6114a4919061470c565b95505050505050611659565b8373ffffffffffffffffffffffffffffffffffffffff16630dfe16816040518163ffffffff1660e01b8152600401602060405180830381865afa1580156114f9573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061151d919061454f565b73ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff16036115d6578673ffffffffffffffffffffffffffffffffffffffff1663ad615dec8b84846040518463ffffffff1660e01b815260040161158c93929190614775565b602060405180830381865afa1580156115a7573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906115cb91906146b4565b945050505050611659565b8673ffffffffffffffffffffffffffffffffffffffff1663ad615dec8b83856040518463ffffffff1660e01b815260040161161393929190614775565b602060405180830381865afa15801561162e573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061165291906146b4565b9450505050505b95945050505050565b61166a612229565b5f600c60016101000a81548160ff0219169083151502179055505f600e5f6101000a81548160ff0219169083151502179055507fa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c60405160405180910390a1565b601181815481106116da575f80fd5b905f5260205f20015f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61170e612229565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611773576040517fdad1a1b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036117d8576040517f74fc211300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8273ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161181291906142cd565b602060405180830381865afa15801561182d573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061185191906146b4565b905061185e83838361304b565b7f5661684995ab94d684bfe57a43c4141578f52d3e7374e8cd3250e2f062e13ac1838383604051611891939291906147aa565b60405180910390a1505050565b6118a6612229565b6013545f146118e1576040517fd723eaba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600c5f6101000a81548160ff021916908315150217905550436013819055507f1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d60405160405180910390a1565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b611960612229565b80600b8190555050565b6060600380546119799061446a565b80601f01602080910402602001604051908101604052809291908181526020018280546119a59061446a565b80156119f05780601f106119c7576101008083540402835291602001916119f0565b820191905f5260205f20905b8154815290600101906020018083116119d357829003601f168201915b5050505050905090565b5f80611a046122b0565b9050611a1181858561235b565b600191505092915050565b6008602052805f5260405f205f915054906101000a900460ff1681565b611a41612229565b8060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df782604051611adb9190613f71565b60405180910390a25050565b600e5f9054906101000a900460ff1681565b5f805f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b611b82612229565b8060085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055505050565b611be2612229565b5f600e5f6101000a81548160ff021916908315150217905550565b611c05612229565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611c75575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401611c6c91906142cd565b60405180910390fd5b611c7e816130ca565b50565b611c89612229565b6005811115611cc4576040517fbc32fbf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161480611d2957505f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16145b15611d60576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f805b601180549050811015611dfa578573ffffffffffffffffffffffffffffffffffffffff1660118281548110611d9b57611d9a61449a565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603611de75760019150611dfa565b8080611df2906144f4565b915050611d63565b5080611e6557601185908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611ee0565b60105f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2060020154821115611edf576040517f5ab576c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5b8460105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508360105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f6101000a81548160ff0219169083151502179055508260105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160016101000a81548160ff0219169083151502179055508560105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508160105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600201819055508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167ff0be2d87461053944aa20ea0674e2f12c7f14eed47161dc6fc4821dc05fc53b486868660405161212f939291906147df565b60405180910390a3505050505050565b612147612229565b5f479050805f03612184576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f3373ffffffffffffffffffffffffffffffffffffffff16826040516121a990614841565b5f6040518083038185875af1925050503d805f81146121e3576040519150601f19603f3d011682016040523d82523d5f602084013e6121e8565b606091505b50508091505080612225576040517f3132169500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b6122316122b0565b73ffffffffffffffffffffffffffffffffffffffff1661224f611930565b73ffffffffffffffffffffffffffffffffffffffff16146122ae576122726122b0565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016122a591906142cd565b60405180910390fd5b565b5f33905090565b6122c4838383600161318d565b505050565b5f6122d48484611af9565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146123555781811015612346578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161233d93929190614855565b60405180910390fd5b61235484848484035f61318d565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036123cb57816040517fec442f050000000000000000000000000000000000000000000000000000000081526004016123c291906142cd565b60405180910390fd5b5f811161240d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612404906148d4565b60405180910390fd5b61241561335b565b158015612468575060145f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b156125485760125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1680612508575060125f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b612547576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161253e9061493c565b60405180910390fd5b5b60085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16806125e3575060085f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b15612619576040517e61c20e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c60019054906101000a900460ff161561298357612636611930565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141580156126a45750612674611930565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156126dc57505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b801561272f575060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b8015612782575060125f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b1561298257600e5f9054906101000a900460ff1615612981575f6127a4613367565b9050806040015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141580156128165750806060015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614155b1561297f57600243612828919061466d565b600d5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541080156128bb575060024361287b919061466d565b600d5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054105b6128fa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016128f1906149ca565b60405180910390fd5b43600d5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555043600d5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055505b505b5b5b5f805b601180549050811015612a7c573373ffffffffffffffffffffffffffffffffffffffff1660105f601184815481106129c1576129c061449a565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603612a695760019150612a7c565b8080612a74906144f4565b915050612986565b50600c5f9054906101000a900460ff168015612aa45750600a5f9054906101000a900460ff16155b8015612af9575060105f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff16155b8015612b03575080155b8015612b56575060125f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b8015612ba9575060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b8015612bbf5750600b54612bbc30611109565b10155b15612c00576001600a5f6101000a81548160ff021916908315150217905550612be6612d8a565b5f600a5f6101000a81548160ff0219169083151502179055505b5f60125f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16158015612c9f575060125f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b90508015612d775760105f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff1615612d76575f612d028661356a565b90505f81608001511115612d74575f6064826080015186612d2391906149e8565b612d2d919061470c565b905080600f5f828254612d409190614a29565b92505081905550612d528730836137af565b8085612d5e919061466d565b9450612d6b8787876137af565b50505050612d85565b505b5b612d828585856137af565b50505b505050565b5f4790505f612d9830611109565b90505f612da3613367565b90505f816040015173ffffffffffffffffffffffffffffffffffffffff1663ad5c46486040518163ffffffff1660e01b8152600401602060405180830381865afa158015612df3573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612e17919061454f565b90505f612e2930838560400151610efc565b9050805f03612e3c575050505050613036565b5f600f5403612e77575f84118015612e5357505f85115b15612e6d57612e63844783613932565b5050505050613036565b5050505050613036565b5f600285612e85919061470c565b90505f8186612e94919061466d565b90505f600267ffffffffffffffff811115612eb257612eb1614a5c565b5b604051908082528060200260200182016040528015612ee05781602001602082028036833780820191505090505b50905030815f81518110612ef757612ef661449a565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508481600181518110612f4657612f4561449a565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050612f8f308760400151846122b7565b856040015173ffffffffffffffffffffffffffffffffffffffff1663791ac94783612fc185308a8c604001518b6111a1565b8430426040518663ffffffff1660e01b8152600401612fe4959493929190614b40565b5f604051808303815f87803b158015612ffb575f80fd5b505af115801561300d573d5f803e3d5ffd5b50505050505f4790505f8190505f600f8190555061302c848287613932565b5050505050505050505b565b613040612229565b6130495f6130ca565b565b6130c5838473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb858560405160240161307e929190614b98565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050613ad4565b505050565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036131fd575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016131f491906142cd565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361326d575f6040517f94280d6200000000000000000000000000000000000000000000000000000000815260040161326491906142cd565b60405180910390fd5b815f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015613355578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161334c9190614057565b60405180910390a35b50505050565b5f6013545f1415905090565b61336f613d18565b613377613d18565b5f5b601180549050811015613562575f6011828154811061339b5761339a61449a565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060105f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160019054906101000a900460ff161561354e5760105f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060a00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481525050925050613562565b50808061355a906144f4565b915050613379565b508091505090565b613572613d18565b61357a613d18565b5f5b6011805490508110156137a5575f6011828154811061359e5761359d61449a565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508473ffffffffffffffffffffffffffffffffffffffff1660105f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036137915760105f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060a00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820154815250509250506137a5565b50808061379d906144f4565b91505061357c565b5080915050919050565b5f60065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015613839578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161383093929190614855565b60405180910390fd5b81810360065f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508160065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516139249190614057565b60405180910390a350505050565b5f61393b613367565b905061394c308260400151866122b7565b5f816040015173ffffffffffffffffffffffffffffffffffffffff1663ad5c46486040518163ffffffff1660e01b8152600401602060405180830381865afa15801561399a573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906139be919061454f565b90505f6139d28630848660400151886111a1565b90505f6139e68684308760400151896111a1565b90505f805f866040015173ffffffffffffffffffffffffffffffffffffffff1663f305d7198a308d898930426040518863ffffffff1660e01b8152600401613a3396959493929190614bbf565b60606040518083038185885af1158015613a4f573d5f803e3d5ffd5b50505050506040513d601f19601f82011682018060405250810190613a749190614c1e565b92509250925082851180613a8757508184115b80613a9157505f81145b15613ac8576040517f0bc488c500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50505050505050505050565b5f613afe828473ffffffffffffffffffffffffffffffffffffffff16613b6990919063ffffffff16565b90505f815114158015613b22575080806020019051810190613b209190614c82565b155b15613b6457826040517f5274afe7000000000000000000000000000000000000000000000000000000008152600401613b5b91906142cd565b60405180910390fd5b505050565b6060613b7683835f613b7e565b905092915050565b606081471015613bc557306040517fcd786059000000000000000000000000000000000000000000000000000000008152600401613bbc91906142cd565b60405180910390fd5b5f808573ffffffffffffffffffffffffffffffffffffffff168486604051613bed9190614ce7565b5f6040518083038185875af1925050503d805f8114613c27576040519150601f19603f3d011682016040523d82523d5f602084013e613c2c565b606091505b5091509150613c3c868383613c47565b925050509392505050565b606082613c5c57613c5782613cd4565b613ccc565b5f8251148015613c8257505f8473ffffffffffffffffffffffffffffffffffffffff163b145b15613cc457836040517f9996b315000000000000000000000000000000000000000000000000000000008152600401613cbb91906142cd565b60405180910390fd5b819050613ccd565b5b9392505050565b5f81511115613ce65780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060a001604052805f151581526020015f151581526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f81525090565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015613da9578082015181840152602081019050613d8e565b5f8484015250505050565b5f601f19601f8301169050919050565b5f613dce82613d72565b613dd88185613d7c565b9350613de8818560208601613d8c565b613df181613db4565b840191505092915050565b5f6020820190508181035f830152613e148184613dc4565b905092915050565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f613e4d82613e24565b9050919050565b613e5d81613e43565b8114613e67575f80fd5b50565b5f81359050613e7881613e54565b92915050565b5f8115159050919050565b613e9281613e7e565b8114613e9c575f80fd5b50565b5f81359050613ead81613e89565b92915050565b5f8060408385031215613ec957613ec8613e1c565b5b5f613ed685828601613e6a565b9250506020613ee785828601613e9f565b9150509250929050565b5f819050919050565b613f0381613ef1565b8114613f0d575f80fd5b50565b5f81359050613f1e81613efa565b92915050565b5f8060408385031215613f3a57613f39613e1c565b5b5f613f4785828601613e6a565b9250506020613f5885828601613f10565b9150509250929050565b613f6b81613e7e565b82525050565b5f602082019050613f845f830184613f62565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f840112613fab57613faa613f8a565b5b8235905067ffffffffffffffff811115613fc857613fc7613f8e565b5b602083019150836020820283011115613fe457613fe3613f92565b5b9250929050565b5f805f6040848603121561400257614001613e1c565b5b5f84013567ffffffffffffffff81111561401f5761401e613e20565b5b61402b86828701613f96565b9350935050602061403e86828701613e9f565b9150509250925092565b61405181613ef1565b82525050565b5f60208201905061406a5f830184614048565b92915050565b5f8083601f84011261408557614084613f8a565b5b8235905067ffffffffffffffff8111156140a2576140a1613f8e565b5b6020830191508360208202830111156140be576140bd613f92565b5b9250929050565b5f805f80604085870312156140dd576140dc613e1c565b5b5f85013567ffffffffffffffff8111156140fa576140f9613e20565b5b61410687828801613f96565b9450945050602085013567ffffffffffffffff81111561412957614128613e20565b5b61413587828801614070565b925092505092959194509250565b5f805f6060848603121561415a57614159613e1c565b5b5f61416786828701613e6a565b935050602061417886828701613e6a565b925050604061418986828701613f10565b9150509250925092565b5f602082840312156141a8576141a7613e1c565b5b5f6141b584828501613e6a565b91505092915050565b6141c781613e43565b82525050565b5f60a0820190506141e05f830188613f62565b6141ed6020830187613f62565b6141fa60408301866141be565b61420760608301856141be565b6142146080830184614048565b9695505050505050565b5f60ff82169050919050565b6142338161421e565b82525050565b5f60208201905061424c5f83018461422a565b92915050565b5f6020828403121561426757614266613e1c565b5b5f61427484828501613e9f565b91505092915050565b5f805f6060848603121561429457614293613e1c565b5b5f6142a186828701613e6a565b93505060206142b286828701613e6a565b92505060406142c386828701613e6a565b9150509250925092565b5f6020820190506142e05f8301846141be565b92915050565b5f805f805f60a086880312156142ff576142fe613e1c565b5b5f61430c88828901613f10565b955050602061431d88828901613e6a565b945050604061432e88828901613e6a565b935050606061433f88828901613e6a565b925050608061435088828901613f10565b9150509295509295909350565b5f6020828403121561437257614371613e1c565b5b5f61437f84828501613f10565b91505092915050565b5f806040838503121561439e5761439d613e1c565b5b5f6143ab85828601613e6a565b92505060206143bc85828601613e6a565b9150509250929050565b5f805f805f60a086880312156143df576143de613e1c565b5b5f6143ec88828901613e6a565b95505060206143fd88828901613e6a565b945050604061440e88828901613e9f565b935050606061441f88828901613e9f565b925050608061443088828901613f10565b9150509295509295909350565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061448157607f821691505b6020821081036144945761449361443d565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6144fe82613ef1565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036145305761452f6144c7565b5b600182019050919050565b5f8151905061454981613e54565b92915050565b5f6020828403121561456457614563613e1c565b5b5f6145718482850161453b565b91505092915050565b5f60408201905061458d5f8301856141be565b61459a60208301846141be565b9392505050565b5f6dffffffffffffffffffffffffffff82169050919050565b6145c3816145a1565b81146145cd575f80fd5b50565b5f815190506145de816145ba565b92915050565b5f63ffffffff82169050919050565b6145fc816145e4565b8114614606575f80fd5b50565b5f81519050614617816145f3565b92915050565b5f805f6060848603121561463457614633613e1c565b5b5f614641868287016145d0565b9350506020614652868287016145d0565b925050604061466386828701614609565b9150509250925092565b5f61467782613ef1565b915061468283613ef1565b925082820390508181111561469a576146996144c7565b5b92915050565b5f815190506146ae81613efa565b92915050565b5f602082840312156146c9576146c8613e1c565b5b5f6146d6848285016146a0565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61471682613ef1565b915061472183613ef1565b925082614731576147306146df565b5b828204905092915050565b5f819050919050565b5f61475f61475a614755846145a1565b61473c565b613ef1565b9050919050565b61476f81614745565b82525050565b5f6060820190506147885f830186614048565b6147956020830185614766565b6147a26040830184614766565b949350505050565b5f6060820190506147bd5f8301866141be565b6147ca60208301856141be565b6147d76040830184614048565b949350505050565b5f6060820190506147f25f830186613f62565b6147ff6020830185613f62565b61480c6040830184614048565b949350505050565b5f81905092915050565b50565b5f61482c5f83614814565b91506148378261481e565b5f82019050919050565b5f61484b82614821565b9150819050919050565b5f6060820190506148685f8301866141be565b6148756020830185614048565b6148826040830184614048565b949350505050565b7f616d6f756e74206d7573742062652067726561746572207468616e20300000005f82015250565b5f6148be601d83613d7c565b91506148c98261488a565b602082019050919050565b5f6020820190508181035f8301526148eb816148b2565b9050919050565b7f54726164696e67206973206e6f74206163746976652e000000000000000000005f82015250565b5f614926601683613d7c565b9150614931826148f2565b602082019050919050565b5f6020820190508181035f8301526149538161491a565b9050919050565b7f5f7472616e736665723a3a205472616e736665722044656c617920656e61626c5f8201527f65642e202054727920616761696e206c617465722e0000000000000000000000602082015250565b5f6149b4603583613d7c565b91506149bf8261495a565b604082019050919050565b5f6020820190508181035f8301526149e1816149a8565b9050919050565b5f6149f282613ef1565b91506149fd83613ef1565b9250828202614a0b81613ef1565b91508282048414831517614a2257614a216144c7565b5b5092915050565b5f614a3382613ef1565b9150614a3e83613ef1565b9250828201905080821115614a5657614a556144c7565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b614abb81613e43565b82525050565b5f614acc8383614ab2565b60208301905092915050565b5f602082019050919050565b5f614aee82614a89565b614af88185614a93565b9350614b0383614aa3565b805f5b83811015614b33578151614b1a8882614ac1565b9750614b2583614ad8565b925050600181019050614b06565b5085935050505092915050565b5f60a082019050614b535f830188614048565b614b606020830187614048565b8181036040830152614b728186614ae4565b9050614b8160608301856141be565b614b8e6080830184614048565b9695505050505050565b5f604082019050614bab5f8301856141be565b614bb86020830184614048565b9392505050565b5f60c082019050614bd25f8301896141be565b614bdf6020830188614048565b614bec6040830187614048565b614bf96060830186614048565b614c0660808301856141be565b614c1360a0830184614048565b979650505050505050565b5f805f60608486031215614c3557614c34613e1c565b5b5f614c42868287016146a0565b9350506020614c53868287016146a0565b9250506040614c64868287016146a0565b9150509250925092565b5f81519050614c7c81613e89565b92915050565b5f60208284031215614c9757614c96613e1c565b5b5f614ca484828501614c6e565b91505092915050565b5f81519050919050565b5f614cc182614cad565b614ccb8185614814565b9350614cdb818560208601613d8c565b80840191505092915050565b5f614cf28284614cb7565b91508190509291505056fea2646970667358221220485ce3a0f209d3cc7f7273be386ea27ecf35a90c3b9674a7d47ced1b487e51d864736f6c63430008140033",
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

// ApplySlippage is a free data retrieval call binding the contract method 0x7371dfff.
//
// Solidity: function applySlippage(uint256 amount, address first, address second, address routerAddress, uint256 timeElapsed) view returns(uint256)
func (_Coindistribution *CoindistributionCaller) ApplySlippage(opts *bind.CallOpts, amount *big.Int, first common.Address, second common.Address, routerAddress common.Address, timeElapsed *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "applySlippage", amount, first, second, routerAddress, timeElapsed)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ApplySlippage is a free data retrieval call binding the contract method 0x7371dfff.
//
// Solidity: function applySlippage(uint256 amount, address first, address second, address routerAddress, uint256 timeElapsed) view returns(uint256)
func (_Coindistribution *CoindistributionSession) ApplySlippage(amount *big.Int, first common.Address, second common.Address, routerAddress common.Address, timeElapsed *big.Int) (*big.Int, error) {
	return _Coindistribution.Contract.ApplySlippage(&_Coindistribution.CallOpts, amount, first, second, routerAddress, timeElapsed)
}

// ApplySlippage is a free data retrieval call binding the contract method 0x7371dfff.
//
// Solidity: function applySlippage(uint256 amount, address first, address second, address routerAddress, uint256 timeElapsed) view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) ApplySlippage(amount *big.Int, first common.Address, second common.Address, routerAddress common.Address, timeElapsed *big.Int) (*big.Int, error) {
	return _Coindistribution.Contract.ApplySlippage(&_Coindistribution.CallOpts, amount, first, second, routerAddress, timeElapsed)
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

// PriceUpdateTimeElapsedForPair is a free data retrieval call binding the contract method 0x52f93b5e.
//
// Solidity: function priceUpdateTimeElapsedForPair(address first, address second, address routerAddress) view returns(uint256)
func (_Coindistribution *CoindistributionCaller) PriceUpdateTimeElapsedForPair(opts *bind.CallOpts, first common.Address, second common.Address, routerAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "priceUpdateTimeElapsedForPair", first, second, routerAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PriceUpdateTimeElapsedForPair is a free data retrieval call binding the contract method 0x52f93b5e.
//
// Solidity: function priceUpdateTimeElapsedForPair(address first, address second, address routerAddress) view returns(uint256)
func (_Coindistribution *CoindistributionSession) PriceUpdateTimeElapsedForPair(first common.Address, second common.Address, routerAddress common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.PriceUpdateTimeElapsedForPair(&_Coindistribution.CallOpts, first, second, routerAddress)
}

// PriceUpdateTimeElapsedForPair is a free data retrieval call binding the contract method 0x52f93b5e.
//
// Solidity: function priceUpdateTimeElapsedForPair(address first, address second, address routerAddress) view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) PriceUpdateTimeElapsedForPair(first common.Address, second common.Address, routerAddress common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.PriceUpdateTimeElapsedForPair(&_Coindistribution.CallOpts, first, second, routerAddress)
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
// Solidity: function uniswapV2LiquidityPoolSlots(address ) view returns(bool enabled, bool receiver, address router, address pairAddress, uint256 bLiqF)
func (_Coindistribution *CoindistributionCaller) UniswapV2LiquidityPoolSlots(opts *bind.CallOpts, arg0 common.Address) (struct {
	Enabled     bool
	Receiver    bool
	Router      common.Address
	PairAddress common.Address
	BLiqF       *big.Int
}, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "uniswapV2LiquidityPoolSlots", arg0)

	outstruct := new(struct {
		Enabled     bool
		Receiver    bool
		Router      common.Address
		PairAddress common.Address
		BLiqF       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Enabled = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Receiver = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Router = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.PairAddress = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.BLiqF = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UniswapV2LiquidityPoolSlots is a free data retrieval call binding the contract method 0x296a803c.
//
// Solidity: function uniswapV2LiquidityPoolSlots(address ) view returns(bool enabled, bool receiver, address router, address pairAddress, uint256 bLiqF)
func (_Coindistribution *CoindistributionSession) UniswapV2LiquidityPoolSlots(arg0 common.Address) (struct {
	Enabled     bool
	Receiver    bool
	Router      common.Address
	PairAddress common.Address
	BLiqF       *big.Int
}, error) {
	return _Coindistribution.Contract.UniswapV2LiquidityPoolSlots(&_Coindistribution.CallOpts, arg0)
}

// UniswapV2LiquidityPoolSlots is a free data retrieval call binding the contract method 0x296a803c.
//
// Solidity: function uniswapV2LiquidityPoolSlots(address ) view returns(bool enabled, bool receiver, address router, address pairAddress, uint256 bLiqF)
func (_Coindistribution *CoindistributionCallerSession) UniswapV2LiquidityPoolSlots(arg0 common.Address) (struct {
	Enabled     bool
	Receiver    bool
	Router      common.Address
	PairAddress common.Address
	BLiqF       *big.Int
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

// SetUniswapV2LiquidityPool is a paid mutator transaction binding the contract method 0xf317f2c3.
//
// Solidity: function setUniswapV2LiquidityPool(address routerAddress, address pairAddress, bool enabled, bool receiver, uint256 bLiqF) returns()
func (_Coindistribution *CoindistributionTransactor) SetUniswapV2LiquidityPool(opts *bind.TransactOpts, routerAddress common.Address, pairAddress common.Address, enabled bool, receiver bool, bLiqF *big.Int) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setUniswapV2LiquidityPool", routerAddress, pairAddress, enabled, receiver, bLiqF)
}

// SetUniswapV2LiquidityPool is a paid mutator transaction binding the contract method 0xf317f2c3.
//
// Solidity: function setUniswapV2LiquidityPool(address routerAddress, address pairAddress, bool enabled, bool receiver, uint256 bLiqF) returns()
func (_Coindistribution *CoindistributionSession) SetUniswapV2LiquidityPool(routerAddress common.Address, pairAddress common.Address, enabled bool, receiver bool, bLiqF *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetUniswapV2LiquidityPool(&_Coindistribution.TransactOpts, routerAddress, pairAddress, enabled, receiver, bLiqF)
}

// SetUniswapV2LiquidityPool is a paid mutator transaction binding the contract method 0xf317f2c3.
//
// Solidity: function setUniswapV2LiquidityPool(address routerAddress, address pairAddress, bool enabled, bool receiver, uint256 bLiqF) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetUniswapV2LiquidityPool(routerAddress common.Address, pairAddress common.Address, enabled bool, receiver bool, bLiqF *big.Int) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetUniswapV2LiquidityPool(&_Coindistribution.TransactOpts, routerAddress, pairAddress, enabled, receiver, bLiqF)
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
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetUniswapV2LiquidityPool is a free log retrieval operation binding the contract event 0xf0be2d87461053944aa20ea0674e2f12c7f14eed47161dc6fc4821dc05fc53b4.
//
// Solidity: event SetUniswapV2LiquidityPool(address indexed router, address indexed pair, bool enabled, bool receiver, uint256 bLiqF)
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

// WatchSetUniswapV2LiquidityPool is a free log subscription operation binding the contract event 0xf0be2d87461053944aa20ea0674e2f12c7f14eed47161dc6fc4821dc05fc53b4.
//
// Solidity: event SetUniswapV2LiquidityPool(address indexed router, address indexed pair, bool enabled, bool receiver, uint256 bLiqF)
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

// ParseSetUniswapV2LiquidityPool is a log parse operation binding the contract event 0xf0be2d87461053944aa20ea0674e2f12c7f14eed47161dc6fc4821dc05fc53b4.
//
// Solidity: event SetUniswapV2LiquidityPool(address indexed router, address indexed pair, bool enabled, bool receiver, uint256 bLiqF)
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
