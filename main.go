package main

import (
	"aavev3-raw-balances-collector/internal/blockfinder"
	"aavev3-raw-balances-collector/internal/datalab"
	"aavev3-raw-balances-collector/internal/events"
	"aavev3-raw-balances-collector/internal/pool"
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Starting Job...")
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	providerUrl := os.Getenv("PROVIDER_URL")
	logsProviderUrl := os.Getenv("LOGS_PROVIDER_URL")

	// start := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	// end := time.Date(2025, 8, 14, 0, 0, 0, 0, time.UTC)
	// for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
	// 	err := DailyEtl(day, accessKeyID, secretAccessKey, providerUrl, logsProviderUrl)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	ref := time.Now().UTC().AddDate(0, 0, -21)
	snapshotDay := time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
	err := DailyEtl(snapshotDay, accessKeyID, secretAccessKey, providerUrl, logsProviderUrl)
	if err != nil {
		panic(err)
	}
}

func DailyEtl(day time.Time, accessKeyID, secretAccessKey, providerUrl, logsProviderUrl string) error {
	fmt.Printf("Starting Job for day %v\n", day.String())

	fmt.Printf("STEP 1 - Connecting to logsClient...")
	logsClient, err := ethclient.Dial(logsProviderUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 2 - Connecting to client...")
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 3 - Setting pool contract...")
	poolCtr, err := pool.NewPool(common.HexToAddress("0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2"), client)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 4 - Finding start block and end block of the day...")
	dayBeginTmstp := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	dayEndTmstp := dayBeginTmstp.AddDate(0, 0, 1)

	// references := utils.GetBlockReferences()
	// refBlockNumber := references[time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC).Unix()]
	// fmt.Println(refBlockNumber)

	refBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	refBlock := big.NewInt(int64(refBlockNumber))

	fromBlock, _, err := blockfinder.FindClosestBlocks(client, uint64(dayBeginTmstp.Unix()), refBlock, 7000)
	if err != nil {
		return err
	}
	_, toBlock, err := blockfinder.FindClosestBlocks(client, uint64(dayEndTmstp.Unix()), refBlock, 7000)
	if err != nil {
		return err
	}
	fmt.Printf("   start block = %v, end block = %v, Done!\n", fromBlock, toBlock)

	fmt.Printf("STEP 5 - Collecting addresses to query...")
	addresses, err := events.ProvideAddressesToQuery(poolCtr, toBlock)
	if err != nil {
		return err
	}
	fmt.Printf("   Found %v target addresses. Done!\n", len(addresses))

	fmt.Printf("STEP 6 - Collecting events...\n")
	logs, err := events.CollectEvents(logsClient, fromBlock, toBlock, big.NewInt(100), addresses)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 7 - Generating and saving outputs...")
	err = datalab.SaveRecords(
		"minio-simple.lab.groupe-genes.fr",
		accessKeyID,
		secretAccessKey,
		logs,
		"projet-datalab-group-jprat",
		"aavev3-raw-datasource/daily-raw-events/raw_events_snapshot_date="+fmt.Sprint(day)[:10]+"/raw_events.json",
	)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")
	return nil
}
