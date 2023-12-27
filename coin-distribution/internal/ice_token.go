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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20ZeroToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeIncreased\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeMoreThan5\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForeignTokenSelfTransfer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Mismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoBots\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSwapBack\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAirDropper\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawStuckETH\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BuyBackTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sniper\",\"type\":\"address\"}],\"name\":\"CaughtEarlyBuyer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EnableTrading\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isExcluded\",\"type\":\"bool\"}],\"name\":\"ExcludeFromFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"MaxTransactionExclusion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"theAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"OnSetUniswapRouter\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"OwnerForcedSwapBack\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RemovedLimits\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"}],\"name\":\"SetUniswapV2LiquidityPool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensSwapped\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethReceived\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensIntoLiquidity\",\"type\":\"uint256\"}],\"name\":\"SwapAndLiquify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferForeignToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxBuyAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxWalletAmount\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"airdropToWallets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"first\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"second\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"}],\"name\":\"applySlippage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"botsCaught\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableTransferDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableTrading\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"excludeFromFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceSwapBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAirDropper\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTradingEnabledBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"limitsInEffect\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"manageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"wallets\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"massManageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"removeLimits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAirDropper\",\"type\":\"address\"}],\"name\":\"setAirDropper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"setKnownUniswapRouters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setSwapBackThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"theAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"setUniswapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"}],\"name\":\"setUniswapV2LiquidityPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokensForLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transferDelayEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferForeignToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"uniswapV2LiquidityPoolSlots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uniswapV2LiquidityPools\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawStuckETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052683635c9adc5dea00000600b555f600c5f6101000a81548160ff0219169083151502179055506001600c60016101000a81548160ff0219169083151502179055506001600e5f6101000a81548160ff0219169083151502179055505f6013553480156200006f575f80fd5b50336040518060400160405280600381526020017f49636500000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f49434500000000000000000000000000000000000000000000000000000000008152508160029081620000ee91906200090d565b5080600390816200010091906200090d565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000176575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016200016d919062000a34565b60405180910390fd5b62000187816200020660201b60201c565b503360055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001db336001620002c960201b60201c565b620001ee306001620002c960201b60201c565b6200020060016200038160201b60201c565b62000a86565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b620002d9620004c460201b60201c565b8060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df78260405162000375919062000a6b565b60405180910390a25050565b62000391620004c460201b60201c565b620003b7737a250d5630b4cf539739df2c5dacb4c659f2488d826200056660201b60201c565b620003dd73e592427a0aece92de3edee1f18e0157c05861564826200056660201b60201c565b620004037368b3465833fb72a70ecdf485e0e4c7bd8665fc45826200056660201b60201c565b6200042973eff92a263d31888d860bd50809a8d171709b7b1c826200056660201b60201c565b6200044f731b81d678ffb9c0263b24a97847620c99d213eb14826200056660201b60201c565b620004757313f4ea83d0bd40e75c8222255bc855a974568dd4826200056660201b60201c565b6200049b731b02da8cb0d097eb8d57a175b88c7d8b47997506826200056660201b60201c565b620004c173d9e1ce17f2641f24ae83637ab66a2cca9c378b9f826200056660201b60201c565b50565b620004d46200067a60201b60201c565b73ffffffffffffffffffffffffffffffffffffffff16620004fa6200068160201b60201c565b73ffffffffffffffffffffffffffffffffffffffff16146200056457620005266200067a60201b60201c565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016200055b919062000a34565b60405180910390fd5b565b62000576620004c460201b60201c565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620005dc576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060145f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508015158273ffffffffffffffffffffffffffffffffffffffff167f3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c60405160405180910390a35050565b5f33905090565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200072557607f821691505b6020821081036200073b576200073a620006e0565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026200079f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000762565b620007ab868362000762565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620007f5620007ef620007e984620007c3565b620007cc565b620007c3565b9050919050565b5f819050919050565b6200081083620007d5565b620008286200081f82620007fc565b8484546200076e565b825550505050565b5f90565b6200083e62000830565b6200084b81848462000805565b505050565b5b818110156200087257620008665f8262000834565b60018101905062000851565b5050565b601f821115620008c1576200088b8162000741565b620008968462000753565b81016020851015620008a6578190505b620008be620008b58562000753565b83018262000850565b50505b505050565b5f82821c905092915050565b5f620008e35f1984600802620008c6565b1980831691505092915050565b5f620008fd8383620008d2565b9150826002028217905092915050565b6200091882620006a9565b67ffffffffffffffff811115620009345762000933620006b3565b5b6200094082546200070d565b6200094d82828562000876565b5f60209050601f83116001811462000983575f84156200096e578287015190505b6200097a8582620008f0565b865550620009e9565b601f198416620009938662000741565b5f5b82811015620009bc5784890151825560018201915060208501945060208101905062000995565b86831015620009dc5784890151620009d8601f891682620008d2565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f62000a1c82620009f1565b9050919050565b62000a2e8162000a10565b82525050565b5f60208201905062000a495f83018462000a23565b92915050565b5f8115159050919050565b62000a658162000a4f565b82525050565b5f60208201905062000a805f83018462000a5a565b92915050565b6147b58062000a945f395ff3fe608060405260043610610228575f3560e01c806370a0823111610122578063a9059cbb116100aa578063e800dff71161006e578063e800dff7146107e1578063e884f26014610809578063f2fde38b1461081f578063f317f2c314610847578063f5648a4f1461086f5761022f565b8063a9059cbb146106db578063bfd7928414610717578063c024666814610753578063c876d0b91461077b578063dd62ed3e146107a55761022f565b80638366e79a116100f15780638366e79a146106215780638a8c523c146106495780638da5cb5b1461065f5780638fe62b8a1461068957806395d89b41146106b15761022f565b806370a082311461057d578063715018a6146105b9578063751039fc146105cf578063785cca3e146105e55761022f565b8063313ce567116101b0578063588b655f11610174578063588b655f146104ad5780635db7548b146104d75780635de998ef146105015780636b0a894c146105295780636ddd1713146105535761022f565b8063313ce567146103df5780634370efda1461040957806344d8a785146104455780634a62bb651461046d57806351f205e4146104975761022f565b806318160ddd116101f757806318160ddd146102e75780631a8145bb146103115780632307b4411461033b57806323b872dd14610363578063296a803c1461039f5761022f565b806306fdde0314610231578063072280c31461025b578063095ea7b314610283578063130a2c3c146102bf5761022f565b3661022f57005b005b34801561023c575f80fd5b50610245610885565b6040516102529190613965565b60405180910390f35b348015610266575f80fd5b50610281600480360381019061027c9190613a1c565b610915565b005b34801561028e575f80fd5b506102a960048036038101906102a49190613a8d565b610a20565b6040516102b69190613ada565b60405180910390f35b3480156102ca575f80fd5b506102e560048036038101906102e09190613b54565b610a42565b005b3480156102f2575f80fd5b506102fb610aeb565b6040516103089190613bc0565b60405180910390f35b34801561031c575f80fd5b50610325610af4565b6040516103329190613bc0565b60405180910390f35b348015610346575f80fd5b50610361600480360381019061035c9190613c2e565b610afa565b005b34801561036e575f80fd5b5061038960048036038101906103849190613cac565b610c2c565b6040516103969190613ada565b60405180910390f35b3480156103aa575f80fd5b506103c560048036038101906103c09190613cfc565b610c5a565b6040516103d6959493929190613d36565b60405180910390f35b3480156103ea575f80fd5b506103f3610ce2565b6040516104009190613da2565b60405180910390f35b348015610414575f80fd5b5061042f600480360381019061042a9190613dbb565b610cea565b60405161043c9190613bc0565b60405180910390f35b348015610450575f80fd5b5061046b60048036038101906104669190613e1f565b610efd565b005b348015610478575f80fd5b50610481610ff8565b60405161048e9190613ada565b60405180910390f35b3480156104a2575f80fd5b506104ab61100b565b005b3480156104b8575f80fd5b506104c16110c8565b6040516104ce9190613bc0565b60405180910390f35b3480156104e2575f80fd5b506104eb6110d1565b6040516104f89190613e4a565b60405180910390f35b34801561050c575f80fd5b5061052760048036038101906105229190613cfc565b6110f9565b005b348015610534575f80fd5b5061053d611144565b60405161054a9190613bc0565b60405180910390f35b34801561055e575f80fd5b5061056761114a565b6040516105749190613ada565b60405180910390f35b348015610588575f80fd5b506105a3600480360381019061059e9190613cfc565b61115c565b6040516105b09190613bc0565b60405180910390f35b3480156105c4575f80fd5b506105cd6111a2565b005b3480156105da575f80fd5b506105e36111f4565b005b3480156105f0575f80fd5b5061060b60048036038101906106069190613e63565b61125d565b6040516106189190613e4a565b60405180910390f35b34801561062c575f80fd5b5061064760048036038101906106429190613e8e565b611298565b005b348015610654575f80fd5b5061065d611430565b005b34801561066a575f80fd5b506106736114c2565b6040516106809190613e4a565b60405180910390f35b348015610694575f80fd5b506106af60048036038101906106aa9190613e63565b6114ea565b005b3480156106bc575f80fd5b506106c56114fc565b6040516106d29190613965565b60405180910390f35b3480156106e6575f80fd5b5061070160048036038101906106fc9190613a8d565b61158c565b60405161070e9190613ada565b60405180910390f35b348015610722575f80fd5b5061073d60048036038101906107389190613cfc565b6115ae565b60405161074a9190613ada565b60405180910390f35b34801561075e575f80fd5b5061077960048036038101906107749190613a1c565b6115cb565b005b348015610786575f80fd5b5061078f611679565b60405161079c9190613ada565b60405180910390f35b3480156107b0575f80fd5b506107cb60048036038101906107c69190613e8e565b61168b565b6040516107d89190613bc0565b60405180910390f35b3480156107ec575f80fd5b5061080760048036038101906108029190613a1c565b61170c565b005b348015610814575f80fd5b5061081d61176c565b005b34801561082a575f80fd5b5061084560048036038101906108409190613cfc565b61178f565b005b348015610852575f80fd5b5061086d60048036038101906108689190613ecc565b611813565b005b34801561087a575f80fd5b50610883611cd1565b005b60606002805461089490613f70565b80601f01602080910402602001604051908101604052809291908181526020018280546108c090613f70565b801561090b5780601f106108e25761010080835404028352916020019161090b565b820191905f5260205f20905b8154815290600101906020018083116108ee57829003601f168201915b5050505050905090565b61091d611dbb565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610982576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060145f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508015158273ffffffffffffffffffffffffffffffffffffffff167f3fefb3a0b9178802e3aa79b6dae4164acd27eba06e14c1cb7bed09fb0801f84c60405160405180910390a35050565b5f80610a2a611e42565b9050610a37818585611e49565b600191505092915050565b610a4a611dbb565b5f5b83839050811015610ae5578160085f868685818110610a6e57610a6d613fa0565b5b9050602002016020810190610a839190613cfc565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508080610add90613ffa565b915050610a4c565b50505050565b5f600754905090565b600f5481565b3373ffffffffffffffffffffffffffffffffffffffff1660055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610b80576040517f16ad4feb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818190508484905014610bbf576040517f77a93d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8484905090505f604051602081016004356024355f5b86811015610c1057602460208202018083013581850135875260068652604087208181540181558189019850505050600181019050610bd6565b50505050508060075f8282540192505081905550505050505050565b5f80610c36611e42565b9050610c43858285611e5b565b610c4e858585611eed565b60019150509392505050565b6010602052805f5260405f205f91509050805f015f9054906101000a900460ff1690805f0160019054906101000a900460ff1690805f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154905085565b5f6012905090565b5f808290505f8173ffffffffffffffffffffffffffffffffffffffff1663c45a01556040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d39573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d5d9190614055565b90505f8173ffffffffffffffffffffffffffffffffffffffff1663e6a4390588886040518363ffffffff1660e01b8152600401610d9b929190614080565b602060405180830381865afa158015610db6573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610dda9190614055565b90505f808273ffffffffffffffffffffffffffffffffffffffff16630902f1ac6040518163ffffffff1660e01b81526004016040805180830381865afa158015610e26573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e4a91906140ea565b80925081935050508473ffffffffffffffffffffffffffffffffffffffff1663ad615dec8b846dffffffffffffffffffffffffffff16846dffffffffffffffffffffffffffff166040518463ffffffff1660e01b8152600401610eaf93929190614128565b602060405180830381865afa158015610eca573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610eee9190614171565b95505050505050949350505050565b610f05611dbb565b610f23737a250d5630b4cf539739df2c5dacb4c659f2488d82610915565b610f4173e592427a0aece92de3edee1f18e0157c0586156482610915565b610f5f7368b3465833fb72a70ecdf485e0e4c7bd8665fc4582610915565b610f7d73eff92a263d31888d860bd50809a8d171709b7b1c82610915565b610f9b731b81d678ffb9c0263b24a97847620c99d213eb1482610915565b610fb97313f4ea83d0bd40e75c8222255bc855a974568dd482610915565b610fd7731b02da8cb0d097eb8d57a175b88c7d8b4799750682610915565b610ff573d9e1ce17f2641f24ae83637ab66a2cca9c378b9f82610915565b50565b600c60019054906101000a900460ff1681565b611013611dbb565b61101c3061115c565b5f03611054576040517f0b952c3600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600a5f6101000a81548160ff02191690831515021790555061107661291c565b5f600a5f6101000a81548160ff0219169083151502179055507f1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb32426040516110be9190613bc0565b60405180910390a1565b5f601354905090565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b611101611dbb565b8060055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60095481565b600c5f9054906101000a900460ff1681565b5f60065f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6111aa611dbb565b6111b26129a5565b5f60055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6111fc611dbb565b5f600c60016101000a81548160ff0219169083151502179055505f600e5f6101000a81548160ff0219169083151502179055507fa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c60405160405180910390a1565b6011818154811061126c575f80fd5b905f5260205f20015f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6112a0611dbb565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611305576040517fdad1a1b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361136a576040517f74fc211300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8273ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016113a49190613e4a565b602060405180830381865afa1580156113bf573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906113e39190614171565b90506113f08383836129b8565b7f5661684995ab94d684bfe57a43c4141578f52d3e7374e8cd3250e2f062e13ac18383836040516114239392919061419c565b60405180910390a1505050565b611438611dbb565b6013545f14611473576040517fd723eaba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600c5f6101000a81548160ff021916908315150217905550436013819055507f1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d60405160405180910390a1565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6114f2611dbb565b80600b8190555050565b60606003805461150b90613f70565b80601f016020809104026020016040519081016040528092919081815260200182805461153790613f70565b80156115825780601f1061155957610100808354040283529160200191611582565b820191905f5260205f20905b81548152906001019060200180831161156557829003601f168201915b5050505050905090565b5f80611596611e42565b90506115a3818585611eed565b600191505092915050565b6008602052805f5260405f205f915054906101000a900460ff1681565b6115d3611dbb565b8060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df78260405161166d9190613ada565b60405180910390a25050565b600e5f9054906101000a900460ff1681565b5f805f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b611714611dbb565b8060085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055505050565b611774611dbb565b5f600e5f6101000a81548160ff021916908315150217905550565b611797611dbb565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611807575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016117fe9190613e4a565b60405180910390fd5b61181081612a37565b50565b61181b611dbb565b6005811115611856576040517fbc32fbf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614806118bb57505f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16145b156118f2576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f805b60118054905081101561198c578573ffffffffffffffffffffffffffffffffffffffff166011828154811061192d5761192c613fa0565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603611979576001915061198c565b808061198490613ffa565b9150506118f5565b50806119f757601185908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611a72565b60105f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2060020154821115611a71576040517f5ab576c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5b8460105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508360105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f6101000a81548160ff0219169083151502179055508260105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160016101000a81548160ff0219169083151502179055508560105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508160105f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600201819055508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167ff0be2d87461053944aa20ea0674e2f12c7f14eed47161dc6fc4821dc05fc53b4868686604051611cc1939291906141d1565b60405180910390a3505050505050565b611cd9611dbb565b5f479050805f03611d16576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f3373ffffffffffffffffffffffffffffffffffffffff1682604051611d3b90614233565b5f6040518083038185875af1925050503d805f8114611d75576040519150601f19603f3d011682016040523d82523d5f602084013e611d7a565b606091505b50508091505080611db7576040517f3132169500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b611dc3611e42565b73ffffffffffffffffffffffffffffffffffffffff16611de16114c2565b73ffffffffffffffffffffffffffffffffffffffff1614611e4057611e04611e42565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401611e379190613e4a565b60405180910390fd5b565b5f33905090565b611e568383836001612afa565b505050565b5f611e66848461168b565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611ee75781811015611ed8578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401611ecf93929190614247565b60405180910390fd5b611ee684848484035f612afa565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611f5d57816040517fec442f05000000000000000000000000000000000000000000000000000000008152600401611f549190613e4a565b60405180910390fd5b5f8111611f9f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f96906142c6565b60405180910390fd5b611fa7612cc8565b158015611ffa575060145f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b156120da5760125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff168061209a575060125f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b6120d9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016120d09061432e565b60405180910390fd5b5b60085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1680612175575060085f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b156121ab576040517e61c20e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600c60019054906101000a900460ff1615612515576121c86114c2565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415801561223657506122066114c2565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b801561226e57505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b80156122c1575060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b8015612314575060125f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b1561251457600e5f9054906101000a900460ff1615612513575f612336612cd4565b9050806040015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141580156123a85750806060015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614155b15612511576002436123ba919061434c565b600d5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410801561244d575060024361240d919061434c565b600d5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054105b61248c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612483906143ef565b60405180910390fd5b43600d5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555043600d5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055505b505b5b5b5f805b60118054905081101561260e573373ffffffffffffffffffffffffffffffffffffffff1660105f6011848154811061255357612552613fa0565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036125fb576001915061260e565b808061260690613ffa565b915050612518565b50600c5f9054906101000a900460ff1680156126365750600a5f9054906101000a900460ff16155b801561268b575060105f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff16155b8015612695575080155b80156126e8575060125f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b801561273b575060125f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b80156127515750600b5461274e3061115c565b10155b15612792576001600a5f6101000a81548160ff02191690831515021790555061277861291c565b5f600a5f6101000a81548160ff0219169083151502179055505b5f60125f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16158015612831575060125f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b905080156129095760105f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff1615612908575f61289486612ed7565b90505f81608001511115612906575f60648260800151866128b5919061440d565b6128bf919061447b565b905080600f5f8282546128d291906144ab565b925050819055506128e487308361311c565b80856128f0919061434c565b94506128fd87878761311c565b50505050612917565b505b5b61291485858561311c565b50505b505050565b5f4790505f61292a3061115c565b90505f600f5403612960575f8111801561294357505f82115b1561295957612952814761329f565b50506129a3565b50506129a3565b5f60028261296e919061447b565b9050612984818361297f919061434c565b61343e565b5f4790505f8190505f600f8190555061299d838261329f565b50505050505b565b6129ad611dbb565b6129b65f612a37565b565b612a32838473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85856040516024016129eb9291906144de565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061363d565b505050565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603612b6a575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401612b619190613e4a565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603612bda575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401612bd19190613e4a565b60405180910390fd5b815f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015612cc2578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051612cb99190613bc0565b60405180910390a35b50505050565b5f6013545f1415905090565b612cdc613881565b612ce4613881565b5f5b601180549050811015612ecf575f60118281548110612d0857612d07613fa0565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060105f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160019054906101000a900460ff1615612ebb5760105f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060a00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481525050925050612ecf565b508080612ec790613ffa565b915050612ce6565b508091505090565b612edf613881565b612ee7613881565b5f5b601180549050811015613112575f60118281548110612f0b57612f0a613fa0565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508473ffffffffffffffffffffffffffffffffffffffff1660105f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036130fe5760105f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060a00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481525050925050613112565b50808061310a90613ffa565b915050612ee9565b5080915050919050565b5f60065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050818110156131a6578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161319d93929190614247565b60405180910390fd5b81810360065f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508160065f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516132919190613bc0565b60405180910390a350505050565b5f6132a8612cd4565b90506132b930826040015185611e49565b5f816040015173ffffffffffffffffffffffffffffffffffffffff1663ad5c46486040518163ffffffff1660e01b8152600401602060405180830381865afa158015613307573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061332b9190614055565b90505f61333e8530848660400151610cea565b90505f6133518584308760400151610cea565b90505f805f866040015173ffffffffffffffffffffffffffffffffffffffff1663f305d71989308c898930426040518863ffffffff1660e01b815260040161339e96959493929190614505565b60606040518083038185885af11580156133ba573d5f803e3d5ffd5b50505050506040513d601f19601f820116820180604052508101906133df9190614564565b925092509250828511806133f257508184115b806133fc57505f81145b15613433576040517f0bc488c500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050565b5f613447612cd4565b90505f816040015173ffffffffffffffffffffffffffffffffffffffff1663ad5c46486040518163ffffffff1660e01b8152600401602060405180830381865afa158015613497573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906134bb9190614055565b90505f600267ffffffffffffffff8111156134d9576134d86145b4565b5b6040519080825280602002602001820160405280156135075781602001602082028036833780820191505090505b50905030815f8151811061351e5761351d613fa0565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050818160018151811061356d5761356c613fa0565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250506135b630846040015186611e49565b826040015173ffffffffffffffffffffffffffffffffffffffff1663791ac947856135e78730878960400151610cea565b8430426040518663ffffffff1660e01b815260040161360a959493929190614698565b5f604051808303815f87803b158015613621575f80fd5b505af1158015613633573d5f803e3d5ffd5b5050505050505050565b5f613667828473ffffffffffffffffffffffffffffffffffffffff166136d290919063ffffffff16565b90505f81511415801561368b5750808060200190518101906136899190614704565b155b156136cd57826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016136c49190613e4a565b60405180910390fd5b505050565b60606136df83835f6136e7565b905092915050565b60608147101561372e57306040517fcd7860590000000000000000000000000000000000000000000000000000000081526004016137259190613e4a565b60405180910390fd5b5f808573ffffffffffffffffffffffffffffffffffffffff1684866040516137569190614769565b5f6040518083038185875af1925050503d805f8114613790576040519150601f19603f3d011682016040523d82523d5f602084013e613795565b606091505b50915091506137a58683836137b0565b925050509392505050565b6060826137c5576137c08261383d565b613835565b5f82511480156137eb57505f8473ffffffffffffffffffffffffffffffffffffffff163b145b1561382d57836040517f9996b3150000000000000000000000000000000000000000000000000000000081526004016138249190613e4a565b60405180910390fd5b819050613836565b5b9392505050565b5f8151111561384f5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040518060a001604052805f151581526020015f151581526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f81525090565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156139125780820151818401526020810190506138f7565b5f8484015250505050565b5f601f19601f8301169050919050565b5f613937826138db565b61394181856138e5565b93506139518185602086016138f5565b61395a8161391d565b840191505092915050565b5f6020820190508181035f83015261397d818461392d565b905092915050565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6139b68261398d565b9050919050565b6139c6816139ac565b81146139d0575f80fd5b50565b5f813590506139e1816139bd565b92915050565b5f8115159050919050565b6139fb816139e7565b8114613a05575f80fd5b50565b5f81359050613a16816139f2565b92915050565b5f8060408385031215613a3257613a31613985565b5b5f613a3f858286016139d3565b9250506020613a5085828601613a08565b9150509250929050565b5f819050919050565b613a6c81613a5a565b8114613a76575f80fd5b50565b5f81359050613a8781613a63565b92915050565b5f8060408385031215613aa357613aa2613985565b5b5f613ab0858286016139d3565b9250506020613ac185828601613a79565b9150509250929050565b613ad4816139e7565b82525050565b5f602082019050613aed5f830184613acb565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f840112613b1457613b13613af3565b5b8235905067ffffffffffffffff811115613b3157613b30613af7565b5b602083019150836020820283011115613b4d57613b4c613afb565b5b9250929050565b5f805f60408486031215613b6b57613b6a613985565b5b5f84013567ffffffffffffffff811115613b8857613b87613989565b5b613b9486828701613aff565b93509350506020613ba786828701613a08565b9150509250925092565b613bba81613a5a565b82525050565b5f602082019050613bd35f830184613bb1565b92915050565b5f8083601f840112613bee57613bed613af3565b5b8235905067ffffffffffffffff811115613c0b57613c0a613af7565b5b602083019150836020820283011115613c2757613c26613afb565b5b9250929050565b5f805f8060408587031215613c4657613c45613985565b5b5f85013567ffffffffffffffff811115613c6357613c62613989565b5b613c6f87828801613aff565b9450945050602085013567ffffffffffffffff811115613c9257613c91613989565b5b613c9e87828801613bd9565b925092505092959194509250565b5f805f60608486031215613cc357613cc2613985565b5b5f613cd0868287016139d3565b9350506020613ce1868287016139d3565b9250506040613cf286828701613a79565b9150509250925092565b5f60208284031215613d1157613d10613985565b5b5f613d1e848285016139d3565b91505092915050565b613d30816139ac565b82525050565b5f60a082019050613d495f830188613acb565b613d566020830187613acb565b613d636040830186613d27565b613d706060830185613d27565b613d7d6080830184613bb1565b9695505050505050565b5f60ff82169050919050565b613d9c81613d87565b82525050565b5f602082019050613db55f830184613d93565b92915050565b5f805f8060808587031215613dd357613dd2613985565b5b5f613de087828801613a79565b9450506020613df1878288016139d3565b9350506040613e02878288016139d3565b9250506060613e13878288016139d3565b91505092959194509250565b5f60208284031215613e3457613e33613985565b5b5f613e4184828501613a08565b91505092915050565b5f602082019050613e5d5f830184613d27565b92915050565b5f60208284031215613e7857613e77613985565b5b5f613e8584828501613a79565b91505092915050565b5f8060408385031215613ea457613ea3613985565b5b5f613eb1858286016139d3565b9250506020613ec2858286016139d3565b9150509250929050565b5f805f805f60a08688031215613ee557613ee4613985565b5b5f613ef2888289016139d3565b9550506020613f03888289016139d3565b9450506040613f1488828901613a08565b9350506060613f2588828901613a08565b9250506080613f3688828901613a79565b9150509295509295909350565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680613f8757607f821691505b602082108103613f9a57613f99613f43565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61400482613a5a565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361403657614035613fcd565b5b600182019050919050565b5f8151905061404f816139bd565b92915050565b5f6020828403121561406a57614069613985565b5b5f61407784828501614041565b91505092915050565b5f6040820190506140935f830185613d27565b6140a06020830184613d27565b9392505050565b5f6dffffffffffffffffffffffffffff82169050919050565b6140c9816140a7565b81146140d3575f80fd5b50565b5f815190506140e4816140c0565b92915050565b5f8060408385031215614100576140ff613985565b5b5f61410d858286016140d6565b925050602061411e858286016140d6565b9150509250929050565b5f60608201905061413b5f830186613bb1565b6141486020830185613bb1565b6141556040830184613bb1565b949350505050565b5f8151905061416b81613a63565b92915050565b5f6020828403121561418657614185613985565b5b5f6141938482850161415d565b91505092915050565b5f6060820190506141af5f830186613d27565b6141bc6020830185613d27565b6141c96040830184613bb1565b949350505050565b5f6060820190506141e45f830186613acb565b6141f16020830185613acb565b6141fe6040830184613bb1565b949350505050565b5f81905092915050565b50565b5f61421e5f83614206565b915061422982614210565b5f82019050919050565b5f61423d82614213565b9150819050919050565b5f60608201905061425a5f830186613d27565b6142676020830185613bb1565b6142746040830184613bb1565b949350505050565b7f616d6f756e74206d7573742062652067726561746572207468616e20300000005f82015250565b5f6142b0601d836138e5565b91506142bb8261427c565b602082019050919050565b5f6020820190508181035f8301526142dd816142a4565b9050919050565b7f54726164696e67206973206e6f74206163746976652e000000000000000000005f82015250565b5f6143186016836138e5565b9150614323826142e4565b602082019050919050565b5f6020820190508181035f8301526143458161430c565b9050919050565b5f61435682613a5a565b915061436183613a5a565b925082820390508181111561437957614378613fcd565b5b92915050565b7f5f7472616e736665723a3a205472616e736665722044656c617920656e61626c5f8201527f65642e202054727920616761696e206c617465722e0000000000000000000000602082015250565b5f6143d96035836138e5565b91506143e48261437f565b604082019050919050565b5f6020820190508181035f830152614406816143cd565b9050919050565b5f61441782613a5a565b915061442283613a5a565b925082820261443081613a5a565b9150828204841483151761444757614446613fcd565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61448582613a5a565b915061449083613a5a565b9250826144a05761449f61444e565b5b828204905092915050565b5f6144b582613a5a565b91506144c083613a5a565b92508282019050808211156144d8576144d7613fcd565b5b92915050565b5f6040820190506144f15f830185613d27565b6144fe6020830184613bb1565b9392505050565b5f60c0820190506145185f830189613d27565b6145256020830188613bb1565b6145326040830187613bb1565b61453f6060830186613bb1565b61454c6080830185613d27565b61455960a0830184613bb1565b979650505050505050565b5f805f6060848603121561457b5761457a613985565b5b5f6145888682870161415d565b93505060206145998682870161415d565b92505060406145aa8682870161415d565b9150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b614613816139ac565b82525050565b5f614624838361460a565b60208301905092915050565b5f602082019050919050565b5f614646826145e1565b61465081856145eb565b935061465b836145fb565b805f5b8381101561468b5781516146728882614619565b975061467d83614630565b92505060018101905061465e565b5085935050505092915050565b5f60a0820190506146ab5f830188613bb1565b6146b86020830187613bb1565b81810360408301526146ca818661463c565b90506146d96060830185613d27565b6146e66080830184613bb1565b9695505050505050565b5f815190506146fe816139f2565b92915050565b5f6020828403121561471957614718613985565b5b5f614726848285016146f0565b91505092915050565b5f81519050919050565b5f6147438261472f565b61474d8185614206565b935061475d8185602086016138f5565b80840191505092915050565b5f6147748284614739565b91508190509291505056fea26469706673582212209a8403e9b51b53938e8d412112422cbbbb1c4b12bfa593535d51bd5ee739053c64736f6c63430008140033",
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

// ApplySlippage is a free data retrieval call binding the contract method 0x4370efda.
//
// Solidity: function applySlippage(uint256 amount, address first, address second, address routerAddress) view returns(uint256)
func (_Coindistribution *CoindistributionCaller) ApplySlippage(opts *bind.CallOpts, amount *big.Int, first common.Address, second common.Address, routerAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "applySlippage", amount, first, second, routerAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ApplySlippage is a free data retrieval call binding the contract method 0x4370efda.
//
// Solidity: function applySlippage(uint256 amount, address first, address second, address routerAddress) view returns(uint256)
func (_Coindistribution *CoindistributionSession) ApplySlippage(amount *big.Int, first common.Address, second common.Address, routerAddress common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.ApplySlippage(&_Coindistribution.CallOpts, amount, first, second, routerAddress)
}

// ApplySlippage is a free data retrieval call binding the contract method 0x4370efda.
//
// Solidity: function applySlippage(uint256 amount, address first, address second, address routerAddress) view returns(uint256)
func (_Coindistribution *CoindistributionCallerSession) ApplySlippage(amount *big.Int, first common.Address, second common.Address, routerAddress common.Address) (*big.Int, error) {
	return _Coindistribution.Contract.ApplySlippage(&_Coindistribution.CallOpts, amount, first, second, routerAddress)
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
