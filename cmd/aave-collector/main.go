package main

import (
	"aavev3-raw-balances-collector/internal/blockfinder"
	"aavev3-raw-balances-collector/internal/datalab"
	"aavev3-raw-balances-collector/internal/events"
	"aavev3-raw-balances-collector/internal/pool"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DailyEtlInputs struct {
	day             time.Time
	region          string
	accessKeyID     string
	secretAccessKey string
	bucket          string
	providerUrl     string
	logsProviderUrl string
	saveToAWS       bool
}

type HourlyEtlInputs struct {
	hour            time.Time
	region          string
	accessKeyID     string
	secretAccessKey string
	bucket          string
	providerUrl     string
	logsProviderUrl string
	saveToAWS       bool
}

func main() {
	fmt.Println("Starting Job...")

	accessKeyID, _ := os.LookupEnv("ACCESS_KEY_ID")
	secretAccessKey, _ := os.LookupEnv("SECRET_ACCESS_KEY")
	region, _ := os.LookupEnv("REGION")
	bucket, isBucketDefined := os.LookupEnv("BUCKET")
	providerUrl := os.Getenv("PROVIDER_URL")
	logsProviderUrl, isLogsProviderUrlDefined := os.LookupEnv("LOGS_PROVIDER_URL")
	runMode := os.Getenv("RUN_MODE")
	startStr, isStartDefined := os.LookupEnv("START_DATE")
	endStr, isEndDefined := os.LookupEnv("END_DATE")
	lagStr, isLagDefined := os.LookupEnv("LAG")

	if !isLogsProviderUrlDefined {
		logsProviderUrl = providerUrl
	}

	if runMode == "archive" {
		fmt.Println("Running in archive mode...")
		if !(isStartDefined && isEndDefined) {
			panic(errors.New("you chose the archive mode, please set START_DATE and END_DATE env variables to the format YYYY-MM-DD"))
		}
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			panic(err)
		}
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			panic(err)
		}
		start = start.UTC()
		end = end.UTC()
		for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
			inputs := DailyEtlInputs{
				day:             day,
				region:          region,
				accessKeyID:     accessKeyID,
				secretAccessKey: secretAccessKey,
				bucket:          bucket,
				providerUrl:     providerUrl,
				logsProviderUrl: logsProviderUrl,
				saveToAWS:       isBucketDefined,
			}
			err := dailyEtl(inputs)
			if err != nil {
				panic(err)
			}
		}
	} else if runMode == "day-lag" {
		fmt.Println("Running in day-lag mode...")
		if !isLagDefined {
			panic(errors.New("you chose the lag mode, please set LAG env variable."))
		}
		lag, err := strconv.Atoi(lagStr)
		if err != nil {
			panic(err)
		}
		ref := time.Now().UTC().AddDate(0, 0, -lag)
		snapshotDay := time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
		inputs := DailyEtlInputs{
			day:             snapshotDay,
			region:          region,
			accessKeyID:     accessKeyID,
			secretAccessKey: secretAccessKey,
			bucket:          bucket,
			providerUrl:     providerUrl,
			logsProviderUrl: logsProviderUrl,
			saveToAWS:       isBucketDefined,
		}
		err = dailyEtl(inputs)
		if err != nil {
			panic(err)
		}
	} else if runMode == "hour-lag" {
		if !isLagDefined {
			panic(errors.New("you chose the lag mode, please set LAG env variable."))
		}
		lag, err := strconv.Atoi(lagStr)
		if err != nil {
			panic(err)
		}
		ref := time.Now().UTC().Add(time.Duration(-lag) * time.Hour)
		snapshotHour := time.Date(ref.Year(), ref.Month(), ref.Day(), ref.Hour(), 0, 0, 0, time.UTC)
		inputs := HourlyEtlInputs{
			hour:            snapshotHour,
			region:          region,
			accessKeyID:     accessKeyID,
			secretAccessKey: secretAccessKey,
			bucket:          bucket,
			providerUrl:     providerUrl,
			logsProviderUrl: logsProviderUrl,
			saveToAWS:       isBucketDefined,
		}
		err = hourlyEtl(inputs)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("invalid RUN_MODE env variable, valid modes are 'archive', 'day-lag' or 'hour-lag'"))
	}
}

func dailyEtl(inputs DailyEtlInputs) error {
	fmt.Printf("Starting Job for day %v\n", inputs.day.String())

	fmt.Printf("STEP 1 - Connecting to logsClient...")
	logsClient, err := ethclient.Dial(inputs.logsProviderUrl)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 2 - Connecting to client...")
	client, err := ethclient.Dial(inputs.providerUrl)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 3 - Setting pool contract...")
	poolCtr, err := pool.NewPool(common.HexToAddress("0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2"), client)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 4 - Finding start block and end block of the day...")
	dayBeginTmstp := time.Date(inputs.day.Year(), inputs.day.Month(), inputs.day.Day(), 0, 0, 0, 0, time.UTC)
	dayEndTmstp := dayBeginTmstp.AddDate(0, 0, 1)

	// references := utils.GetBlockReferences()
	// refBlockNumber := references[time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC).Unix()]
	// fmt.Println(refBlockNumber)

	refBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
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
	logs, err := events.AsyncCollectEvents(logsClient, fromBlock, toBlock, big.NewInt(9), addresses)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 7 - Generating and saving outputs...")
	if inputs.saveToAWS {
		if err = datalab.SaveRecordsAWS(
			inputs.region,
			inputs.accessKeyID,
			inputs.secretAccessKey,
			inputs.bucket,
			"test/daily-raw-events/raw_events_snapshot_date="+fmt.Sprint(inputs.day)[:10]+"/raw_events.json",
			logs,
		); err != nil {
			return err
		}
	} else {
		if err = datalab.SaveRecordsLOCAL(
			"./data/daily-raw-events/raw_events_snapshot_date="+fmt.Sprint(inputs.day)[:10]+"/raw_events.json",
			logs,
		); err != nil {
			return err
		}
	}

	fmt.Println("Done!")
	return nil
}

func hourlyEtl(inputs HourlyEtlInputs) error {
	fmt.Printf("Starting Job for hour %v\n", inputs.hour.String())

	fmt.Printf("STEP 1 - Connecting to logsClient...")
	logsClient, err := ethclient.Dial(inputs.logsProviderUrl)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 2 - Connecting to client...")
	client, err := ethclient.Dial(inputs.providerUrl)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 3 - Setting pool contract...")
	poolCtr, err := pool.NewPool(common.HexToAddress("0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2"), client)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 4 - Finding start block and end block of the hour...")
	hourBeginTmstp := time.Date(inputs.hour.Year(), inputs.hour.Month(), inputs.hour.Day(), inputs.hour.Hour(), 0, 0, 0, time.UTC)
	hourEndTmstp := hourBeginTmstp.Add(time.Hour)

	// references := utils.GetBlockReferences()
	// refBlockNumber := references[time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC).Unix()]
	// fmt.Println(refBlockNumber)

	refBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	refBlock := big.NewInt(int64(refBlockNumber))

	fromBlock, _, err := blockfinder.FindClosestBlocks(client, uint64(hourBeginTmstp.Unix()), refBlock, 7000)
	if err != nil {
		return err
	}
	_, toBlock, err := blockfinder.FindClosestBlocks(client, uint64(hourEndTmstp.Unix()), refBlock, 7000)
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
	logs, err := events.AsyncCollectEvents(logsClient, fromBlock, toBlock, big.NewInt(9), addresses)
	if err != nil {
		return err
	}
	fmt.Printf(" Done!\n")

	fmt.Printf("STEP 7 - Generating and saving outputs...")
	if inputs.saveToAWS {
		if err = datalab.SaveRecordsAWS(
			inputs.region,
			inputs.accessKeyID,
			inputs.secretAccessKey,
			inputs.bucket,
			"test/hourly-raw-events/raw_events_snapshot_date="+fmt.Sprint(inputs.hour)[:10]+"/raw_events_hour="+fmt.Sprint(inputs.hour)[11:13]+".json",
			logs,
		); err != nil {
			return err
		}
	} else {
		if err = datalab.SaveRecordsLOCAL(
			"./data/hourly-raw-events/raw_events_snapshot_date="+fmt.Sprint(inputs.hour)[:10]+"/raw_events_hour="+fmt.Sprint(inputs.hour)[11:13]+".json",
			logs,
		); err != nil {
			return err
		}
	}

	fmt.Println("Done!")
	return nil
}
