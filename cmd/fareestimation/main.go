package main

import (
	"flag"
	"fmt"

	"github.com/kanzafirov/fareestimation/internal/fs/csv"
	"github.com/kanzafirov/fareestimation/internal/ride"
	"github.com/kanzafirov/fareestimation/internal/ride/financial"
	"github.com/kanzafirov/fareestimation/internal/sync"
)

func main() {
	inpath := flag.String("inpath", "./tests/testdata/test.csv", "path to file with input data")
	outPath := flag.String("outpath", "result.csv", "path to newly created file with result")
	flag.Parse()

	err := runFareEstimation(*inpath, *outPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("complete")
}

func runFareEstimation(inpath string, outPath string) error {
	readChannel, readErrors, err := csv.NewReadChannel(inpath)
	if err != nil {
		return fmt.Errorf("Failed to get read channel; err: %v", err)
	}
	parseChannel, parseErrors := ride.NewParseChannel(readChannel)
	calculationChannel := financial.NewFairFeeCalculationChannel(parseChannel)
	done, writeErrors, err := csv.NewWriteChannel(outPath, calculationChannel)
	if err != nil {
		return fmt.Errorf("Failed to initialize write channel; err: %v", err)
	}

	select {
	case <-done:
	case err := <-sync.MergeErrorChannels(readErrors, parseErrors, writeErrors):
		if err != nil {
			return fmt.Errorf("[runFareEstimation] error from channels: %v", err)
		}
	}

	return nil
}
