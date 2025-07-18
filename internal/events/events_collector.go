package events

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExecuteQuery(client *ethclient.Client, fromBlock, toBlock *big.Int, addresses []common.Address) ([]types.Log, error) {
	fmt.Printf("   --> Collecting logs from blocks %v to %v,", fromBlock, toBlock)
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return make([]types.Log, 0), err
	}
	fmt.Printf("   Found %v events\n", len(logs))
	return logs, nil
}

func CollectEvents(client *ethclient.Client, fromBlock, toBlock, blockStep *big.Int, addresses []common.Address) ([]types.Log, error) {

	logs := make([]types.Log, 0)

	start := new(big.Int)
	stop := new(big.Int)

	start.Set(fromBlock)
	stop.Add(fromBlock, blockStep)

	for start.Cmp(toBlock) < 0 {
		currentLogs, err := ExecuteQuery(client, start, stop, addresses)
		if err != nil {
			return make([]types.Log, 0), err
		}
		logs = append(logs, currentLogs...)
		start.Set(stop)
		stop.Add(stop, blockStep)
		if stop.Cmp(toBlock) > 0 {
			stop.Set(toBlock)
		}
	}
	return logs, nil
}

// func CollectEvents(client *ethclient.Client, fromBlock, toBlock, blockStep *big.Int, addresses []common.Address) ([]types.Log, error) {

// 	logs := make([]types.Log, 0)

// 	start := new(big.Int)
// 	stop := new(big.Int)

// 	start.Set(fromBlock)
// 	stop.Add(fromBlock, blockStep)

// 	for stop.Cmp(toBlock) <= 0 {
// 		time.Sleep(time.Second)
// 		fmt.Println("Collecting events from", fromBlock, "to", toBlock)
// 		for _, contract := range addresses {
// 			currentLogs, err := ExecuteQuery(client, start, stop, []common.Address{contract})
// 			if err != nil {
// 				return make([]types.Log, 0), err
// 			}
// 			logs = append(logs, currentLogs...)
// 		}
// 		start.Set(stop)
// 		stop.Add(stop, blockStep)
// 	}
// 	return logs, nil
// }
