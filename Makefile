all: help 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## bench/for_select_go_work: benchmark pool which uses for_select loop pattern, and create one goroutine for every job.
.PHONY: bench/for_select_go_work
bench/for_select_go_work:
	@go test -v --tags=for_select_go_work --bench=. -benchmem

## bench/for_select_multiple_workerChans_1: benchmark pool which uses for_select pattern, pool: for-select, worker: for-select pattern
.PHONY: bench/for_select_multiple_workerChans_1
bench/for_select_multiple_workerChans_1:
	@go test -v --tags=for_select_multiple_workerChans_1 --bench=. -benchmem

## bench/for_select_multiple_workerChans_2: benchmark pool which uses for_select pattern, pool: for-select, worker: for-range pattern
.PHONY: bench/for_select_multiple_workerChans_2
bench/for_select_multiple_workerChans_2:
	@go test -v --tags=for_select_multiple_workerChans_2 --bench=. -benchmem
