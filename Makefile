all: help 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


## bench/all: benchmark pool 
.PHONY: bench/all
bench/all: bench/for_select_go_work bench/for_select_multiple_workerChans_1 bench/for_select_multiple_workerChans_2 bench/no_jobChan_1 bench/no_jobChan_2 bench/native_goroutine

## bench/for_select_go_work: benchmark pool which uses for_select loop pattern, and create one goroutine for every job.
.PHONY: bench/for_select_go_work
bench/for_select_go_work:
	@echo "\nbenchmark pool which uses for_select loop pattern, and create one goroutine for every job...\n"
	@go test -v --tags=for_select_go_work --bench=./naivepool -benchmem

## bench/for_select_multiple_workerChans_1: benchmark pool which uses for_select pattern, pool: for-select, worker: for-select pattern
.PHONY: bench/for_select_multiple_workerChans_1
bench/for_select_multiple_workerChans_1:
	@echo "\nbenchmark pool which uses for_select pattern, pool: for-select, worker: for-select pattern\n"
	@go test -v --tags=for_select_multiple_workerChans_1 --bench=./naivepool -benchmem

## bench/for_select_multiple_workerChans_2: benchmark pool which uses for_select pattern, pool: for-select, worker: for-range pattern
.PHONY: bench/for_select_multiple_workerChans_2
bench/for_select_multiple_workerChans_2:
	@echo "\nbenchmark pool which uses for_select pattern, pool: for-select, worker: for-range pattern\n"
	@go test -v --tags=for_select_multiple_workerChans_2 --bench=./naivepool -benchmem

## bench/no_jobChan_1: benchmark pool which doesn't have jobChan, pool: for-range, worker: for-range pattern
.PHONY: bench/no_jobChan_1
bench/no_jobChan_1:
	@echo "\nbenchmark pool which doesn't have jobChan, pool: for-range, worker: for-range pattern\n"
	@go test -v --tags=no_jobChan_1 --bench=./naivepool -benchmem

## bench/no_jobChan_2: benchmark pool which doesn't have jobChan, and all workers share a single workerChan: for-range, worker: for-range pattern
.PHONY: bench/no_jobChan_2
bench/no_jobChan_2:
	@echo "\nbenchmark pool which doesn't have jobChan, and all workers share a single workerChan: for-range, worker: for-range pattern\n"
	@go test -v --tags=no_jobChan_2 --bench=./naivepool -benchmem

## bench/native_goroutine: benchmark native goroutine
.PHONY: bench/native_goroutine
bench/native_goroutine:
	@echo "\n benchmark native goroutine\n"
	@go test -v --tags=no_jobChan_2 --bench=./native_goroutine -benchmem
