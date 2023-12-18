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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC20ZeroToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeIncreased\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeMoreThan25\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeMoreThan5\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalanceForOnTopFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Mismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoBots\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSwapBack\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TradingAlreadyEnabled\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BuyBackTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sniper\",\"type\":\"address\"}],\"name\":\"CaughtEarlyBuyer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DisableTrading\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EnableTrading\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isExcluded\",\"type\":\"bool\"}],\"name\":\"ExcludeFromFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"MaxTransactionExclusion\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"OwnerForcedSwapBack\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RemovedLimits\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sLiqF\",\"type\":\"uint256\"}],\"name\":\"SetUniswapV2LiquidityPool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensSwapped\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethReceived\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokensIntoLiquidity\",\"type\":\"uint256\"}],\"name\":\"SwapAndLiquify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferForeignToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxBuyAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxSellAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"UpdatedMaxWalletAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newWallet\",\"type\":\"address\"}],\"name\":\"UpdatedOperationsAddress\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"airdropToWallets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"botsCaught\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableTrading\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableTransferDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableTrading\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"excluded\",\"type\":\"bool\"}],\"name\":\"excludeFromFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceSwapBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTradingEnabledBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"limitsInEffect\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"manageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"wallets\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"flag\",\"type\":\"bool\"}],\"name\":\"massManageBoughtEarly\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operationsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"removeLimits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operationsAddress\",\"type\":\"address\"}],\"name\":\"setOperationsAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"routerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sLiqF\",\"type\":\"uint256\"}],\"name\":\"setUniswapV2LiquidityPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokensForLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transferDelayEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"transferForeignToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_sent\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"uniswapV2LiquidityPoolSlots\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"receiver\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pairAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bLiqF\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sLiqF\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uniswapV2LiquidityPools\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawStuckETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040525f600a60016101000a81548160ff0219169083151502179055506001600a60026101000a81548160ff0219169083151502179055506001600c5f6101000a81548160ff0219169083151502179055505f60115534801562000063575f80fd5b50336040518060400160405280600381526020017f49636500000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f49434500000000000000000000000000000000000000000000000000000000008152508160029081620000e2919062000782565b508060039081620000f4919062000782565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036200016a575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401620001619190620008a9565b60405180910390fd5b6200017b816200023b60201b60201c565b503360075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001ee60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16620002fe60201b60201c565b6200022260075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660016200039560201b60201c565b620002353060016200039560201b60201c565b620008fb565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6200030e6200044d60201b60201c565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000381575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401620003789190620008a9565b60405180910390fd5b62000392816200023b60201b60201c565b50565b620003a56200044d60201b60201c565b8060105f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df782604051620004419190620008e0565b60405180910390a25050565b6200045d620004ef60201b60201c565b73ffffffffffffffffffffffffffffffffffffffff1662000483620004f660201b60201c565b73ffffffffffffffffffffffffffffffffffffffff1614620004ed57620004af620004ef60201b60201c565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401620004e49190620008a9565b60405180910390fd5b565b5f33905090565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200059a57607f821691505b602082108103620005b057620005af62000555565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620006147fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620005d7565b620006208683620005d7565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6200066a620006646200065e8462000638565b62000641565b62000638565b9050919050565b5f819050919050565b62000685836200064a565b6200069d620006948262000671565b848454620005e3565b825550505050565b5f90565b620006b3620006a5565b620006c08184846200067a565b505050565b5b81811015620006e757620006db5f82620006a9565b600181019050620006c6565b5050565b601f82111562000736576200070081620005b6565b6200070b84620005c8565b810160208510156200071b578190505b620007336200072a85620005c8565b830182620006c5565b50505b505050565b5f82821c905092915050565b5f620007585f19846008026200073b565b1980831691505092915050565b5f62000772838362000747565b9150826002028217905092915050565b6200078d826200051e565b67ffffffffffffffff811115620007a957620007a862000528565b5b620007b5825462000582565b620007c2828285620006eb565b5f60209050601f831160018114620007f8575f8415620007e3578287015190505b620007ef858262000765565b8655506200085e565b601f1984166200080886620005b6565b5f5b8281101562000831578489015182556001820191506020850194506020810190506200080a565b868310156200085157848901516200084d601f89168262000747565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f620008918262000866565b9050919050565b620008a38162000885565b82525050565b5f602082019050620008be5f83018462000898565b92915050565b5f8115159050919050565b620008da81620008c4565b82525050565b5f602082019050620008f55f830184620008cf565b92915050565b613f1280620009095f395ff3fe608060405260043610610207575f3560e01c806370a0823111610117578063bfd792841161009f578063e800dff71161006e578063e800dff714610735578063e884f2601461075d578063ea4cfe1214610773578063f2fde38b1461079d578063f5648a4f146107c55761020e565b8063bfd792841461066b578063c0246668146106a7578063c876d0b9146106cf578063dd62ed3e146106f95761020e565b80638366e79a116100e65780638366e79a146105895780638a8c523c146105c55780638da5cb5b146105db57806395d89b4114610605578063a9059cbb1461062f5761020e565b806370a08231146104e5578063715018a614610521578063751039fc14610537578063785cca3e1461054d5761020e565b8063296a803c1161019a5780634a62bb65116101695780634a62bb651461042757806351f205e414610451578063588b655f146104675780636b0a894c146104915780636ddd1713146104bb5761020e565b8063296a803c1461036c578063313ce567146103ad5780633d6d11b0146103d7578063499b8394146103ff5761020e565b806318160ddd116101d657806318160ddd146102b45780631a8145bb146102de5780632307b4411461030857806323b872dd146103305761020e565b806306fdde0314610210578063095ea7b31461023a578063130a2c3c1461027657806317700f011461029e5761020e565b3661020e57005b005b34801561021b575f80fd5b506102246107db565b6040516102319190613243565b60405180910390f35b348015610245575f80fd5b50610260600480360381019061025b91906132f8565b61086b565b60405161026d9190613350565b60405180910390f35b348015610281575f80fd5b5061029c600480360381019061029791906133f4565b61088d565b005b3480156102a9575f80fd5b506102b2610936565b005b3480156102bf575f80fd5b506102c86109c8565b6040516102d59190613460565b60405180910390f35b3480156102e9575f80fd5b506102f26109d1565b6040516102ff9190613460565b60405180910390f35b348015610313575f80fd5b5061032e600480360381019061032991906134ce565b6109d7565b005b34801561033b575f80fd5b506103566004803603810190610351919061354c565b610a8b565b6040516103639190613350565b60405180910390f35b348015610377575f80fd5b50610392600480360381019061038d919061359c565b610ab9565b6040516103a4969594939291906135d6565b60405180910390f35b3480156103b8575f80fd5b506103c1610b47565b6040516103ce9190613650565b60405180910390f35b3480156103e2575f80fd5b506103fd60048036038101906103f89190613669565b610b4f565b005b34801561040a575f80fd5b506104256004803603810190610420919061359c565b61103e565b005b348015610432575f80fd5b5061043b611152565b6040516104489190613350565b60405180910390f35b34801561045c575f80fd5b50610465611165565b005b348015610472575f80fd5b5061047b611222565b6040516104889190613460565b60405180910390f35b34801561049c575f80fd5b506104a561122b565b6040516104b29190613460565b60405180910390f35b3480156104c6575f80fd5b506104cf611231565b6040516104dc9190613350565b60405180910390f35b3480156104f0575f80fd5b5061050b6004803603810190610506919061359c565b611244565b6040516105189190613460565b60405180910390f35b34801561052c575f80fd5b5061053561128a565b005b348015610542575f80fd5b5061054b61129d565b005b348015610558575f80fd5b50610573600480360381019061056e91906136f2565b611306565b604051610580919061371d565b60405180910390f35b348015610594575f80fd5b506105af60048036038101906105aa9190613736565b611341565b6040516105bc9190613350565b60405180910390f35b3480156105d0575f80fd5b506105d96114e7565b005b3480156105e6575f80fd5b506105ef61157a565b6040516105fc919061371d565b60405180910390f35b348015610610575f80fd5b506106196115a2565b6040516106269190613243565b60405180910390f35b34801561063a575f80fd5b50610655600480360381019061065091906132f8565b611632565b6040516106629190613350565b60405180910390f35b348015610676575f80fd5b50610691600480360381019061068c919061359c565b611654565b60405161069e9190613350565b60405180910390f35b3480156106b2575f80fd5b506106cd60048036038101906106c89190613774565b611671565b005b3480156106da575f80fd5b506106e361171f565b6040516106f09190613350565b60405180910390f35b348015610704575f80fd5b5061071f600480360381019061071a9190613736565b611731565b60405161072c9190613460565b60405180910390f35b348015610740575f80fd5b5061075b60048036038101906107569190613774565b6117b2565b005b348015610768575f80fd5b50610771611812565b005b34801561077e575f80fd5b50610787611835565b604051610794919061371d565b60405180910390f35b3480156107a8575f80fd5b506107c360048036038101906107be919061359c565b61185a565b005b3480156107d0575f80fd5b506107d96118de565b005b6060600280546107ea906137df565b80601f0160208091040260200160405190810160405280929190818152602001828054610816906137df565b80156108615780601f1061083857610100808354040283529160200191610861565b820191905f5260205f20905b81548152906001019060200180831161084457829003601f168201915b5050505050905090565b5f80610875611953565b905061088281858561195a565b600191505092915050565b61089561196c565b5f5b83839050811015610930578160085f8686858181106108b9576108b861380f565b5b90506020020160208101906108ce919061359c565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550808061092890613869565b915050610897565b50505050565b61093e61196c565b6011545f03610979576040517fb01ef5a300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600a60016101000a81548160ff0219169083151502179055505f6011819055507fec2947b67240d5cf7801542dd5d0fce67e1fa1106caab214ce6092e1d9e2731b60405160405180910390a1565b5f600654905090565b600d5481565b6109df61196c565b818190508484905014610a1e576040517f77a93d8d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8484905090505f604051602081016004356024355f5b86811015610a6f57602460208202018083013581850135875260058652604087208181540181558189019850505050600181019050610a35565b50505050508060065f8282540192505081905550505050505050565b5f80610a95611953565b9050610aa28582856119f3565b610aad858585611a85565b60019150509392505050565b600e602052805f5260405f205f91509050805f015f9054906101000a900460ff1690805f0160019054906101000a900460ff1690805f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905086565b5f6012905090565b610b5761196c565b6019811115610b92576040517f3601c4ac00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6005821115610bcd576040517fbc32fbf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f805b600f80549050811015610c67578673ffffffffffffffffffffffffffffffffffffffff16600f8281548110610c0857610c0761380f565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610c545760019150610c67565b8080610c5f90613869565b915050610bd0565b5080610cd257600f86908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610d97565b600e5f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2060020154831180610d5f5750600e5f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206003015482115b15610d96576040517f5ab576c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5b85600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555084600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f6101000a81548160ff02191690831515021790555083600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160016101000a81548160ff02191690831515021790555086600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555082600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206002018190555081600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600301819055508573ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff167f78c6052d21b64dbee719270b6a9e56ba166f8f88fbaa39a44abf4a11642bae928787878760405161102d94939291906138b0565b60405180910390a350505050505050565b61104661196c565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036110ab576040517f4f41053900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f4efa56652237561d0f1fd31311aeaaa41f3b754a461545ed3cf6ced5876d298260405160405180910390a250565b600a60029054906101000a900460ff1681565b61116d61196c565b61117630611244565b5f036111ae576040517f0b952c3600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600a5f6101000a81548160ff0219169083151502179055506111d0612599565b5f600a5f6101000a81548160ff0219169083151502179055507f1b56c383f4f48fc992e45667ea4eabae777b9cca68b516a9562d8cda78f1bb32426040516112189190613460565b60405180910390a1565b5f601154905090565b60095481565b600a60019054906101000a900460ff1681565b5f60055f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b61129261196c565b61129b5f612647565b565b6112a561196c565b5f600a60026101000a81548160ff0219169083151502179055505f600c5f6101000a81548160ff0219169083151502179055507fa4ffae85e880608d5d4365c2b682786545d136145537788e7e0940dff9f0b98c60405160405180910390a1565b600f8181548110611315575f80fd5b905f5260205f20015f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f61134a61196c565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036113af576040517fdad1a1b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016113e9919061371d565b602060405180830381865afa158015611404573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906114289190613907565b90508373ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84836040518363ffffffff1660e01b8152600401611465929190613932565b6020604051808303815f875af1158015611481573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906114a5919061396d565b91507fdeda980967fcead7b61e78ac46a4da14274af29e894d4d61e8b81ec38ab3e43884826040516114d8929190613932565b60405180910390a15092915050565b6114ef61196c565b6011545f1461152a576040517fd723eaba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600a60016101000a81548160ff021916908315150217905550436011819055507f1d97b7cdf6b6f3405cbe398b69512e5419a0ce78232b6e9c6ffbf1466774bd8d60405160405180910390a1565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6060600380546115b1906137df565b80601f01602080910402602001604051908101604052809291908181526020018280546115dd906137df565b80156116285780601f106115ff57610100808354040283529160200191611628565b820191905f5260205f20905b81548152906001019060200180831161160b57829003601f168201915b5050505050905090565b5f8061163c611953565b9050611649818585611a85565b600191505092915050565b6008602052805f5260405f205f915054906101000a900460ff1681565b61167961196c565b8060105f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff167f9d8f7706ea1113d1a167b526eca956215946dd36cc7df39eb16180222d8b5df7826040516117139190613350565b60405180910390a25050565b600c5f9054906101000a900460ff1681565b5f805f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b6117ba61196c565b8060085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055505050565b61181a61196c565b5f600c5f6101000a81548160ff021916908315150217905550565b60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61186261196c565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036118d2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016118c9919061371d565b60405180910390fd5b6118db81612647565b50565b6118e661196c565b5f3373ffffffffffffffffffffffffffffffffffffffff164760405161190b906139c5565b5f6040518083038185875af1925050503d805f8114611945576040519150601f19603f3d011682016040523d82523d5f602084013e61194a565b606091505b50508091505050565b5f33905090565b611967838383600161270a565b505050565b611974611953565b73ffffffffffffffffffffffffffffffffffffffff1661199261157a565b73ffffffffffffffffffffffffffffffffffffffff16146119f1576119b5611953565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016119e8919061371d565b60405180910390fd5b565b5f6119fe8484611731565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611a7f5781811015611a70578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401611a67939291906139d9565b60405180910390fd5b611a7e84848484035f61270a565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611af557816040517fec442f05000000000000000000000000000000000000000000000000000000008152600401611aec919061371d565b60405180910390fd5b5f8111611b37576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611b2e90613a58565b60405180910390fd5b611b3f6128d8565b611c1e5760105f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1680611bde575060105f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b611c1d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611c1490613ac0565b60405180910390fd5b5b60085f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1680611cb9575060085f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff165b15611cef576040517e61c20e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a60029054906101000a900460ff161561205957611d0c61157a565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614158015611d7a5750611d4a61157a565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015611db257505f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b8015611e05575060105f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b8015611e58575060105f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b1561205857600c5f9054906101000a900460ff1615612057575f611e7a6128e4565b9050806040015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614158015611eec5750806060015173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614155b1561205557600243611efe9190613ade565b600b5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054108015611f915750600243611f519190613ade565b600b5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054105b611fd0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611fc790613b81565b60405180910390fd5b43600b5f3273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555043600b5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055505b505b5b5b5f805b600f8054905081101561214e573373ffffffffffffffffffffffffffffffffffffffff16600e5f600f84815481106120975761209661380f565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361213b57600191505b808061214690613869565b91505061205c565b50600a60019054906101000a900460ff1680156121775750600a5f9054906101000a900460ff16155b80156121cc5750600e5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff16155b80156121d6575080155b8015612229575060105f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b801561227c575060105f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b156122bd576001600a5f6101000a81548160ff0219169083151502179055506122a3612599565b5f600a5f6101000a81548160ff0219169083151502179055505b5f60105f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1615801561235c575060105f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16155b90508015612586575f8061236e6128e4565b9050600e5f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff1680156123cb57505f8160a00151115b1561249e575f6123da87612af1565b90505f8160a001519050606481886123f29190613b9f565b6123fc9190613c0d565b9350808260a001518561240f9190613b9f565b6124199190613c0d565b600d5f8282546124299190613c3d565b92505081905550838761243c9190613c3d565b6124458a611244565b101561247d576040517f136c487600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612488893086612d40565b612493898989612d40565b505050505050612594565b600e5f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f015f9054906101000a900460ff1680156124f957505f8160800151115b15612561575f61250888612af1565b90505f81608001519050606481886125209190613b9f565b61252a9190613c0d565b93508082608001518561253d9190613b9f565b6125479190613c0d565b600d5f8282546125579190613c3d565b9250508190555050505b5f82111561257557612574873084612d40565b5b81856125819190613ade565b945050505b612591858585612d40565b50505b505050565b5f6125a330611244565b90505f600d5490505f811480156125ba57505f8214155b156125d0576125c98247612ec3565b5050612645565b5f811480156125de57505f82145b156125ea575050612645565b5f600282600d54856125fc9190613b9f565b6126069190613c0d565b6126109190613c0d565b905061262681846126219190613ade565b612f6e565b5f4790505f8190505f600d8190555061263f8382612ec3565b50505050505b565b5f60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff160361277a575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401612771919061371d565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036127ea575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016127e1919061371d565b60405180910390fd5b815f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555080156128d2578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516128c99190613460565b60405180910390a35b50505050565b5f6011545f1415905090565b6128ec613159565b6128f4613159565b5f5b600f80549050811015612ae9575f600f82815481106129185761291761380f565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600e5f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f0160019054906101000a900460ff1615612ad557600e5f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200160028201548152602001600382015481525050925050612ae9565b508080612ae190613869565b9150506128f6565b508091505090565b612af9613159565b612b01613159565b5f5b600f80549050811015612d36575f600f8281548110612b2557612b2461380f565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508473ffffffffffffffffffffffffffffffffffffffff16600e5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603612d2257600e5f8273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015f9054906101000a900460ff161515151581526020015f820160019054906101000a900460ff161515151581526020015f820160029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200160028201548152602001600382015481525050925050612d36565b508080612d2e90613869565b915050612b03565b5080915050919050565b5f60055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015612dca578381836040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401612dc1939291906139d9565b60405180910390fd5b81810360055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508160055f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612eb59190613460565b60405180910390a350505050565b5f612ecc6128e4565b9050612edd3082604001518561195a565b806040015173ffffffffffffffffffffffffffffffffffffffff1663f305d7198330865f8030426040518863ffffffff1660e01b8152600401612f2596959493929190613cb2565b60606040518083038185885af1158015612f41573d5f803e3d5ffd5b50505050506040513d601f19601f82011682018060405250810190612f669190613d11565b505050505050565b5f612f776128e4565b90505f600267ffffffffffffffff811115612f9557612f94613d61565b5b604051908082528060200260200182016040528015612fc35781602001602082028036833780820191505090505b50905030815f81518110612fda57612fd961380f565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050816040015173ffffffffffffffffffffffffffffffffffffffff1663ad5c46486040518163ffffffff1660e01b8152600401602060405180830381865afa158015613061573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906130859190613da2565b816001815181106130995761309861380f565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250506130e23083604001518561195a565b816040015173ffffffffffffffffffffffffffffffffffffffff1663791ac947845f8430426040518663ffffffff1660e01b8152600401613127959493929190613e84565b5f604051808303815f87803b15801561313e575f80fd5b505af1158015613150573d5f803e3d5ffd5b50505050505050565b6040518060c001604052805f151581526020015f151581526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff1681526020015f81526020015f81525090565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156131f05780820151818401526020810190506131d5565b5f8484015250505050565b5f601f19601f8301169050919050565b5f613215826131b9565b61321f81856131c3565b935061322f8185602086016131d3565b613238816131fb565b840191505092915050565b5f6020820190508181035f83015261325b818461320b565b905092915050565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6132948261326b565b9050919050565b6132a48161328a565b81146132ae575f80fd5b50565b5f813590506132bf8161329b565b92915050565b5f819050919050565b6132d7816132c5565b81146132e1575f80fd5b50565b5f813590506132f2816132ce565b92915050565b5f806040838503121561330e5761330d613263565b5b5f61331b858286016132b1565b925050602061332c858286016132e4565b9150509250929050565b5f8115159050919050565b61334a81613336565b82525050565b5f6020820190506133635f830184613341565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f84011261338a57613389613369565b5b8235905067ffffffffffffffff8111156133a7576133a661336d565b5b6020830191508360208202830111156133c3576133c2613371565b5b9250929050565b6133d381613336565b81146133dd575f80fd5b50565b5f813590506133ee816133ca565b92915050565b5f805f6040848603121561340b5761340a613263565b5b5f84013567ffffffffffffffff81111561342857613427613267565b5b61343486828701613375565b93509350506020613447868287016133e0565b9150509250925092565b61345a816132c5565b82525050565b5f6020820190506134735f830184613451565b92915050565b5f8083601f84011261348e5761348d613369565b5b8235905067ffffffffffffffff8111156134ab576134aa61336d565b5b6020830191508360208202830111156134c7576134c6613371565b5b9250929050565b5f805f80604085870312156134e6576134e5613263565b5b5f85013567ffffffffffffffff81111561350357613502613267565b5b61350f87828801613375565b9450945050602085013567ffffffffffffffff81111561353257613531613267565b5b61353e87828801613479565b925092505092959194509250565b5f805f6060848603121561356357613562613263565b5b5f613570868287016132b1565b9350506020613581868287016132b1565b9250506040613592868287016132e4565b9150509250925092565b5f602082840312156135b1576135b0613263565b5b5f6135be848285016132b1565b91505092915050565b6135d08161328a565b82525050565b5f60c0820190506135e95f830189613341565b6135f66020830188613341565b61360360408301876135c7565b61361060608301866135c7565b61361d6080830185613451565b61362a60a0830184613451565b979650505050505050565b5f60ff82169050919050565b61364a81613635565b82525050565b5f6020820190506136635f830184613641565b92915050565b5f805f805f8060c0878903121561368357613682613263565b5b5f61369089828a016132b1565b96505060206136a189828a016132b1565b95505060406136b289828a016133e0565b94505060606136c389828a016133e0565b93505060806136d489828a016132e4565b92505060a06136e589828a016132e4565b9150509295509295509295565b5f6020828403121561370757613706613263565b5b5f613714848285016132e4565b91505092915050565b5f6020820190506137305f8301846135c7565b92915050565b5f806040838503121561374c5761374b613263565b5b5f613759858286016132b1565b925050602061376a858286016132b1565b9150509250929050565b5f806040838503121561378a57613789613263565b5b5f613797858286016132b1565b92505060206137a8858286016133e0565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806137f657607f821691505b602082108103613809576138086137b2565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f613873826132c5565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036138a5576138a461383c565b5b600182019050919050565b5f6080820190506138c35f830187613341565b6138d06020830186613341565b6138dd6040830185613451565b6138ea6060830184613451565b95945050505050565b5f81519050613901816132ce565b92915050565b5f6020828403121561391c5761391b613263565b5b5f613929848285016138f3565b91505092915050565b5f6040820190506139455f8301856135c7565b6139526020830184613451565b9392505050565b5f81519050613967816133ca565b92915050565b5f6020828403121561398257613981613263565b5b5f61398f84828501613959565b91505092915050565b5f81905092915050565b50565b5f6139b05f83613998565b91506139bb826139a2565b5f82019050919050565b5f6139cf826139a5565b9150819050919050565b5f6060820190506139ec5f8301866135c7565b6139f96020830185613451565b613a066040830184613451565b949350505050565b7f616d6f756e74206d7573742062652067726561746572207468616e20300000005f82015250565b5f613a42601d836131c3565b9150613a4d82613a0e565b602082019050919050565b5f6020820190508181035f830152613a6f81613a36565b9050919050565b7f54726164696e67206973206e6f74206163746976652e000000000000000000005f82015250565b5f613aaa6016836131c3565b9150613ab582613a76565b602082019050919050565b5f6020820190508181035f830152613ad781613a9e565b9050919050565b5f613ae8826132c5565b9150613af3836132c5565b9250828203905081811115613b0b57613b0a61383c565b5b92915050565b7f5f7472616e736665723a3a205472616e736665722044656c617920656e61626c5f8201527f65642e202054727920616761696e206c617465722e0000000000000000000000602082015250565b5f613b6b6035836131c3565b9150613b7682613b11565b604082019050919050565b5f6020820190508181035f830152613b9881613b5f565b9050919050565b5f613ba9826132c5565b9150613bb4836132c5565b9250828202613bc2816132c5565b91508282048414831517613bd957613bd861383c565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f613c17826132c5565b9150613c22836132c5565b925082613c3257613c31613be0565b5b828204905092915050565b5f613c47826132c5565b9150613c52836132c5565b9250828201905080821115613c6a57613c6961383c565b5b92915050565b5f819050919050565b5f819050919050565b5f613c9c613c97613c9284613c70565b613c79565b6132c5565b9050919050565b613cac81613c82565b82525050565b5f60c082019050613cc55f8301896135c7565b613cd26020830188613451565b613cdf6040830187613ca3565b613cec6060830186613ca3565b613cf960808301856135c7565b613d0660a0830184613451565b979650505050505050565b5f805f60608486031215613d2857613d27613263565b5b5f613d35868287016138f3565b9350506020613d46868287016138f3565b9250506040613d57868287016138f3565b9150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f81519050613d9c8161329b565b92915050565b5f60208284031215613db757613db6613263565b5b5f613dc484828501613d8e565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b613dff8161328a565b82525050565b5f613e108383613df6565b60208301905092915050565b5f602082019050919050565b5f613e3282613dcd565b613e3c8185613dd7565b9350613e4783613de7565b805f5b83811015613e77578151613e5e8882613e05565b9750613e6983613e1c565b925050600181019050613e4a565b5085935050505092915050565b5f60a082019050613e975f830188613451565b613ea46020830187613ca3565b8181036040830152613eb68186613e28565b9050613ec560608301856135c7565b613ed26080830184613451565b969550505050505056fea26469706673582212206a526e331321b2477c54520c447d52d5269a3e2439909ff38367c21f5a0c6bfd64736f6c63430008140033",
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

// OperationsAddress is a free data retrieval call binding the contract method 0xea4cfe12.
//
// Solidity: function operationsAddress() view returns(address)
func (_Coindistribution *CoindistributionCaller) OperationsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Coindistribution.contract.Call(opts, &out, "operationsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperationsAddress is a free data retrieval call binding the contract method 0xea4cfe12.
//
// Solidity: function operationsAddress() view returns(address)
func (_Coindistribution *CoindistributionSession) OperationsAddress() (common.Address, error) {
	return _Coindistribution.Contract.OperationsAddress(&_Coindistribution.CallOpts)
}

// OperationsAddress is a free data retrieval call binding the contract method 0xea4cfe12.
//
// Solidity: function operationsAddress() view returns(address)
func (_Coindistribution *CoindistributionCallerSession) OperationsAddress() (common.Address, error) {
	return _Coindistribution.Contract.OperationsAddress(&_Coindistribution.CallOpts)
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

// DisableTrading is a paid mutator transaction binding the contract method 0x17700f01.
//
// Solidity: function disableTrading() returns()
func (_Coindistribution *CoindistributionTransactor) DisableTrading(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "disableTrading")
}

// DisableTrading is a paid mutator transaction binding the contract method 0x17700f01.
//
// Solidity: function disableTrading() returns()
func (_Coindistribution *CoindistributionSession) DisableTrading() (*types.Transaction, error) {
	return _Coindistribution.Contract.DisableTrading(&_Coindistribution.TransactOpts)
}

// DisableTrading is a paid mutator transaction binding the contract method 0x17700f01.
//
// Solidity: function disableTrading() returns()
func (_Coindistribution *CoindistributionTransactorSession) DisableTrading() (*types.Transaction, error) {
	return _Coindistribution.Contract.DisableTrading(&_Coindistribution.TransactOpts)
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

// SetOperationsAddress is a paid mutator transaction binding the contract method 0x499b8394.
//
// Solidity: function setOperationsAddress(address _operationsAddress) returns()
func (_Coindistribution *CoindistributionTransactor) SetOperationsAddress(opts *bind.TransactOpts, _operationsAddress common.Address) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "setOperationsAddress", _operationsAddress)
}

// SetOperationsAddress is a paid mutator transaction binding the contract method 0x499b8394.
//
// Solidity: function setOperationsAddress(address _operationsAddress) returns()
func (_Coindistribution *CoindistributionSession) SetOperationsAddress(_operationsAddress common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetOperationsAddress(&_Coindistribution.TransactOpts, _operationsAddress)
}

// SetOperationsAddress is a paid mutator transaction binding the contract method 0x499b8394.
//
// Solidity: function setOperationsAddress(address _operationsAddress) returns()
func (_Coindistribution *CoindistributionTransactorSession) SetOperationsAddress(_operationsAddress common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.SetOperationsAddress(&_Coindistribution.TransactOpts, _operationsAddress)
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
// Solidity: function transferForeignToken(address _token, address _to) returns(bool _sent)
func (_Coindistribution *CoindistributionTransactor) TransferForeignToken(opts *bind.TransactOpts, _token common.Address, _to common.Address) (*types.Transaction, error) {
	return _Coindistribution.contract.Transact(opts, "transferForeignToken", _token, _to)
}

// TransferForeignToken is a paid mutator transaction binding the contract method 0x8366e79a.
//
// Solidity: function transferForeignToken(address _token, address _to) returns(bool _sent)
func (_Coindistribution *CoindistributionSession) TransferForeignToken(_token common.Address, _to common.Address) (*types.Transaction, error) {
	return _Coindistribution.Contract.TransferForeignToken(&_Coindistribution.TransactOpts, _token, _to)
}

// TransferForeignToken is a paid mutator transaction binding the contract method 0x8366e79a.
//
// Solidity: function transferForeignToken(address _token, address _to) returns(bool _sent)
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

// CoindistributionDisableTradingIterator is returned from FilterDisableTrading and is used to iterate over the raw logs and unpacked data for DisableTrading events raised by the Coindistribution contract.
type CoindistributionDisableTradingIterator struct {
	Event *CoindistributionDisableTrading // Event containing the contract specifics and raw log

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
func (it *CoindistributionDisableTradingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionDisableTrading)
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
		it.Event = new(CoindistributionDisableTrading)
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
func (it *CoindistributionDisableTradingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionDisableTradingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionDisableTrading represents a DisableTrading event raised by the Coindistribution contract.
type CoindistributionDisableTrading struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDisableTrading is a free log retrieval operation binding the contract event 0xec2947b67240d5cf7801542dd5d0fce67e1fa1106caab214ce6092e1d9e2731b.
//
// Solidity: event DisableTrading()
func (_Coindistribution *CoindistributionFilterer) FilterDisableTrading(opts *bind.FilterOpts) (*CoindistributionDisableTradingIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "DisableTrading")
	if err != nil {
		return nil, err
	}
	return &CoindistributionDisableTradingIterator{contract: _Coindistribution.contract, event: "DisableTrading", logs: logs, sub: sub}, nil
}

// WatchDisableTrading is a free log subscription operation binding the contract event 0xec2947b67240d5cf7801542dd5d0fce67e1fa1106caab214ce6092e1d9e2731b.
//
// Solidity: event DisableTrading()
func (_Coindistribution *CoindistributionFilterer) WatchDisableTrading(opts *bind.WatchOpts, sink chan<- *CoindistributionDisableTrading) (event.Subscription, error) {

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "DisableTrading")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionDisableTrading)
				if err := _Coindistribution.contract.UnpackLog(event, "DisableTrading", log); err != nil {
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

// ParseDisableTrading is a log parse operation binding the contract event 0xec2947b67240d5cf7801542dd5d0fce67e1fa1106caab214ce6092e1d9e2731b.
//
// Solidity: event DisableTrading()
func (_Coindistribution *CoindistributionFilterer) ParseDisableTrading(log types.Log) (*CoindistributionDisableTrading, error) {
	event := new(CoindistributionDisableTrading)
	if err := _Coindistribution.contract.UnpackLog(event, "DisableTrading", log); err != nil {
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
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferForeignToken is a free log retrieval operation binding the contract event 0xdeda980967fcead7b61e78ac46a4da14274af29e894d4d61e8b81ec38ab3e438.
//
// Solidity: event TransferForeignToken(address token, uint256 amount)
func (_Coindistribution *CoindistributionFilterer) FilterTransferForeignToken(opts *bind.FilterOpts) (*CoindistributionTransferForeignTokenIterator, error) {

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "TransferForeignToken")
	if err != nil {
		return nil, err
	}
	return &CoindistributionTransferForeignTokenIterator{contract: _Coindistribution.contract, event: "TransferForeignToken", logs: logs, sub: sub}, nil
}

// WatchTransferForeignToken is a free log subscription operation binding the contract event 0xdeda980967fcead7b61e78ac46a4da14274af29e894d4d61e8b81ec38ab3e438.
//
// Solidity: event TransferForeignToken(address token, uint256 amount)
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

// ParseTransferForeignToken is a log parse operation binding the contract event 0xdeda980967fcead7b61e78ac46a4da14274af29e894d4d61e8b81ec38ab3e438.
//
// Solidity: event TransferForeignToken(address token, uint256 amount)
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

// CoindistributionUpdatedOperationsAddressIterator is returned from FilterUpdatedOperationsAddress and is used to iterate over the raw logs and unpacked data for UpdatedOperationsAddress events raised by the Coindistribution contract.
type CoindistributionUpdatedOperationsAddressIterator struct {
	Event *CoindistributionUpdatedOperationsAddress // Event containing the contract specifics and raw log

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
func (it *CoindistributionUpdatedOperationsAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoindistributionUpdatedOperationsAddress)
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
		it.Event = new(CoindistributionUpdatedOperationsAddress)
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
func (it *CoindistributionUpdatedOperationsAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoindistributionUpdatedOperationsAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoindistributionUpdatedOperationsAddress represents a UpdatedOperationsAddress event raised by the Coindistribution contract.
type CoindistributionUpdatedOperationsAddress struct {
	NewWallet common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUpdatedOperationsAddress is a free log retrieval operation binding the contract event 0x4efa56652237561d0f1fd31311aeaaa41f3b754a461545ed3cf6ced5876d2982.
//
// Solidity: event UpdatedOperationsAddress(address indexed newWallet)
func (_Coindistribution *CoindistributionFilterer) FilterUpdatedOperationsAddress(opts *bind.FilterOpts, newWallet []common.Address) (*CoindistributionUpdatedOperationsAddressIterator, error) {

	var newWalletRule []interface{}
	for _, newWalletItem := range newWallet {
		newWalletRule = append(newWalletRule, newWalletItem)
	}

	logs, sub, err := _Coindistribution.contract.FilterLogs(opts, "UpdatedOperationsAddress", newWalletRule)
	if err != nil {
		return nil, err
	}
	return &CoindistributionUpdatedOperationsAddressIterator{contract: _Coindistribution.contract, event: "UpdatedOperationsAddress", logs: logs, sub: sub}, nil
}

// WatchUpdatedOperationsAddress is a free log subscription operation binding the contract event 0x4efa56652237561d0f1fd31311aeaaa41f3b754a461545ed3cf6ced5876d2982.
//
// Solidity: event UpdatedOperationsAddress(address indexed newWallet)
func (_Coindistribution *CoindistributionFilterer) WatchUpdatedOperationsAddress(opts *bind.WatchOpts, sink chan<- *CoindistributionUpdatedOperationsAddress, newWallet []common.Address) (event.Subscription, error) {

	var newWalletRule []interface{}
	for _, newWalletItem := range newWallet {
		newWalletRule = append(newWalletRule, newWalletItem)
	}

	logs, sub, err := _Coindistribution.contract.WatchLogs(opts, "UpdatedOperationsAddress", newWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoindistributionUpdatedOperationsAddress)
				if err := _Coindistribution.contract.UnpackLog(event, "UpdatedOperationsAddress", log); err != nil {
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

// ParseUpdatedOperationsAddress is a log parse operation binding the contract event 0x4efa56652237561d0f1fd31311aeaaa41f3b754a461545ed3cf6ced5876d2982.
//
// Solidity: event UpdatedOperationsAddress(address indexed newWallet)
func (_Coindistribution *CoindistributionFilterer) ParseUpdatedOperationsAddress(log types.Log) (*CoindistributionUpdatedOperationsAddress, error) {
	event := new(CoindistributionUpdatedOperationsAddress)
	if err := _Coindistribution.contract.UnpackLog(event, "UpdatedOperationsAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
