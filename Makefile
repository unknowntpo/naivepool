all: help 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## bench/naivepool/all: benchmark all tests of naivepool
.PHONY: bench/naivepool/all
bench/naivepool/all:
	@go test -v --bench=BenchmarkNaivepool -benchmem

## bench/naivepool/simple_task: benchmark naivepool with simple task
.PHONY: bench/naivepool/simple_task
bench/naivepool/simple_task:
	@go test -v --bench=BenchmarkNaivepool/simple_task -benchmem

## bench/naivepool/long-running_task: benchmark naivepool with long-running task
.PHONY: bench/naivepool/long-running_task
bench/naivepool/long-running_task:
	@go test -v --bench=BenchmarkNaivepool/long-running_task -benchmem

## bench/naivepool/print: benchmark naivepool with IO bound function: print
.PHONY: bench/naivepool/print
bench/naivepool/print:
	@go test -v --bench=BenchmarkNaivepool/print -benchmem
