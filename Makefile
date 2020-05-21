mainprogram=run

ifndef inpath
override inpath = ./tests/testdata/test.csv
endif

ifndef outpath
override outpath = result.csv
endif

.PHONY: build
build:
	@go build \
		-o $(mainprogram) cmd/fareestimation/main.go

.PHONY: run
run:
	make build 
	@./$(mainprogram) \
		-inpath=${inpath} \
		-outpath=${outpath}

.PHONY: test
test:
	@go test -v -cover ./...

.PHONY: runall
runall:
	make test
	make run

.PHONY: help
help:
	@echo ""
	@echo "Available tasks:"
	@echo "    build"
	@echo "        code compilation"
	@echo "    test"
	@echo "        run all tests"
	@echo "    run inpath=./tests/testdata/file.csv outpath=result.csv"
	@echo "        build and run"
	@echo "    runall inpath=./tests/testdata/file.csv outpath=result.csv"
	@echo "        test, build and run"
	@echo ""