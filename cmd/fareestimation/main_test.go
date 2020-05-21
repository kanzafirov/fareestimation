package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kanzafirov/fareestimation/internal/fs/csv"
)

func TestMain(t *testing.T) {
	before()
	err := runFareEstimation("../../tests/testdata/test.csv", "../../e2eresult.csv")
	if err != nil {
		t.Errorf("[TestMain] Failed: program executaion failed; err: %v", err)
	}
	areEqual, err := csv.DeepCompare("../../e2eresult.csv", "../../tests/testdata/testresult.csv")
	if err != nil {
		t.Errorf("[TestMain] Failed: failed comparing the csv; err: %v", err)
	}
	if !areEqual {
		t.Errorf("[TestMain] Failed: e2e test failed - output is not as expected")
	}
}

func before() error {
	resultFile := "../../e2eresult.csv"
	fullpath, err := filepath.Abs(resultFile)
	if err != nil {
		log.Panicf("Invalid file path: %v; err: %v ", resultFile, err)
	}

	// check if result file already exists and delete it
	if _, err := os.Stat(resultFile); err == nil {
		err = os.Remove(fullpath)
		if err != nil {
			log.Fatalf("[TestMain-before] Failed: test precondition failed with err: %v", err)
		}
	}
	return nil
}
