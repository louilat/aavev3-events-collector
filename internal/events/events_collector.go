package events

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

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

type job struct {
	start *big.Int
	stop  *big.Int
}

type jobResult struct {
	logs []types.Log
	err  error
}

func worker(ctx context.Context, client *ethclient.Client, addresses []common.Address, jobs <-chan job, results chan<- jobResult, ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-jobs:
			if !ok {
				return
			}
			<-ticker.C
			logs, err := ExecuteQuery(client, j.start, j.stop, addresses)
			results <- jobResult{logs: logs, err: err}
		}
	}
}

func AsyncCollectEvents(client *ethclient.Client, fromBlock, toBlock, blockStep *big.Int, addresses []common.Address) ([]types.Log, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ensures context is canceled on function exit

	const maxRequestsPerSecond = 5
	ticker := time.NewTicker(time.Second / maxRequestsPerSecond)
	defer ticker.Stop()

	jobs := make(chan job)          // Input channel
	results := make(chan jobResult) // Output channel
	var wg sync.WaitGroup
	const numWorkers = 5

	// Start workers
	for range numWorkers {
		wg.Add(1)
		go worker(ctx, client, addresses, jobs, results, ticker, &wg)
	}

	go func() {
		wg.Wait() // wait for all workers to end
		close(results)
	}()

	// Fill input channel
	go func() {
		start := new(big.Int).Set(fromBlock)
		stop := new(big.Int).Add(start, blockStep)
		for start.Cmp(toBlock) < 0 {
			if stop.Cmp(toBlock) > 0 {
				stop.Set(toBlock)
			}
			jobs <- job{start: new(big.Int).Set(start), stop: new(big.Int).Set(stop)}
			start.Set(stop)
			stop.Add(stop, blockStep)
		}
		close(jobs)
	}()

	// Collect results
	logs := make([]types.Log, 0)
	for res := range results {
		if res.err != nil {
			cancel()
			return make([]types.Log, 0), res.err
		}
		logs = append(logs, res.logs...)
	}
	return logs, nil
}
