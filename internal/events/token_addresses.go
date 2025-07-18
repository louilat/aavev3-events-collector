package events

import (
	"aavev3-raw-balances-collector/internal/pool"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

type TokenAddress struct {
	UnderlyingAsset common.Address
	ATokenAddress   common.Address
	VTokenAddress   common.Address
}

func FindTokenAddresses(pool *pool.Pool, blockNumber *big.Int) ([]TokenAddress, error) {

	output := make([]TokenAddress, 0)

	addresses, err := pool.GetReservesList(&bind.CallOpts{BlockNumber: blockNumber})
	if err != nil {
		return make([]TokenAddress, 0), nil
	}

	for _, a := range addresses {
		reserveData, err := pool.GetReserveData(&bind.CallOpts{BlockNumber: blockNumber}, a)
		if err != nil {
			return make([]TokenAddress, 0), nil
		}

		output = append(
			output,
			TokenAddress{
				UnderlyingAsset: a,
				ATokenAddress:   reserveData.ATokenAddress,
				VTokenAddress:   reserveData.VariableDebtTokenAddress,
			},
		)
	}
	return output, nil
}

func ProvideAddressesToQuery(pool *pool.Pool, blockNumber *big.Int) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	addresses = append(addresses, common.HexToAddress("0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2"))

	tokenAddresses, err := FindTokenAddresses(pool, blockNumber)
	if err != nil {
		return make([]common.Address, 0), err
	}
	for _, a := range tokenAddresses {
		addresses = append(addresses, a.ATokenAddress)
		addresses = append(addresses, a.VTokenAddress)
	}
	return addresses, nil
}
