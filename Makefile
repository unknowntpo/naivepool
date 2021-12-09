all: help 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## plot: plot the data
.PHONY: plot
plot:
	@echo "Post processing the data..."
	#@cat go-work
	#BenchmarkNaivepool/simple_task/1K_tasks-4
	@sed -n -E 's/BenchmarkNaivepool\/(.*)-[0-9]+/\1/p' go-work.dat | column -t -s ' ' | awk '{print $$1, $$3}' > go-work-final.dat
	@sed -n -E 's/BenchmarkNaivepool\/(.*)-[0-9]+/\1/p' for-select-for-select.dat | column -t -s ' ' | awk '{print $$1, $$3}' > for-select-for-select-final.dat
	@sed -n -E 's/BenchmarkNaivepool\/(.*)-[0-9]+/\1/p' for-select-for-range.dat | column -t -s ' ' | awk '{print $$1, $$3}' > for-select-for-range-final.dat
	@sed -n -E 's/BenchmarkNaivepool\/(.*)-[0-9]+/\1/p' no-jobChan-for-range-worker.dat | column -t -s ' ' | awk '{print $$1, $$3}' > no-jobChan-for-range-worker-final.dat
	@echo "Done post processing.\n"
	@echo "plotting data with gnuplot..."
	#TODO: plot with gnuplot

## bench/naivepool/different-impl: benchmark all tests of naivepool between different implementation.
.PHONY: bench/naivepool/different-impl
bench/naivepool/different-impl:
	@echo "Running benchmark between different implementation of naivepool..."
	@go test -v --bench=BenchmarkNormalGoroutine > normal-goroutine
	@go test -v --bench=BenchmarkPond > pond
	@git checkout feat-pool-go-work && make bench/naivepool/all > go-work.dat
	@git checkout feat-dispatcher-for-select-worker-for-select && make bench/naivepool/all > for-select-for-select.dat
	@git checkout feat-dispatcher-for-select-worker-for-range && make bench/naivepool/all > for-select-for-range.dat
	@git checkout feat-no-jobChan-with-for-range-worker && make bench/naivepool/all > no-jobChan-for-range-worker.dat
	@git checkout master
	@benchstat -html no-jobChan-for-range-worker for-select-for-range for-select-for-select go-work > output.html
	@benchstat -html no-jobChan-for-range-worker normal-goroutine pond > cmp-normal-pond.html


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

## bench/naivepool/print: benchmark naivepool with IO bound function - print
.PHONY: bench/naivepool/print
bench/naivepool/print:
	@go test -v --bench=BenchmarkNaivepool/print -benchmem
