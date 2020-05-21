# fareestimation

A script, calculating a fare extimation for taxi rides

## Project Overview

Fare Estimation Script loads a file with records of type (id_ride, lat, lng, timestamp), creates the calculation for each id and writes down in file the information in format  (id_ride, fare_estimate). For the implementation is used go v1.14.

## Implementation

Fare Estimation Script is designed to be concurrency efficient in order to process large datasets with drivers’ coordinates. In order to provide that efficiency, channel pipelines are used:

## How to Run

```bash

Usage:
  make runall inpath=<inpath> outpath=<outpath>

Flags:
      inpath                  The path of the csv file with all the location signals ex. ./test/testdata/test.csv
      outpath                 The path of the csv file, where the result is stored ex. ./results/out.csv
```
