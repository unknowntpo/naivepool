all: help 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## bench/naivepool/different-impl: benchmark all tests of naivepool between different implementation.
.PHONY: bench/naivepool/different-impl
bench/naivepool/different-impl:
	@mkdir -p /tmp/naivepool
	@echo "Running benchmark between different implementation of naivepool..."
	@git checkout feat-pool-go-work && make bench/naivepool/all > /tmp/naivepool/go-work.txt
	@git checkout feat-dispatcher-for-select-worker-for-select && make bench/naivepool/all > /tmp/naivepool/for-select-for-select.txt
	@git checkout feat-dispatcher-for-select-worker-for-range && make bench/naivepool/all > /tmp/naivepool/for-select-for-range.txt
	@git checkout feat-no-jobChan-with-for-range-worker && make bench/naivepool/all > /tmp/naivepool/no-jobChan-for-range-worker.txt
	@git checkout master
	@benchstat -html /tmp/naivepool/no-jobChan-for-range-worker.txt /tmp/naivepool/for-select-for-range.txt /tmp/naivepool/for-select-for-select.txt /tmp/naivepool/go-work.txt > output.html

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
