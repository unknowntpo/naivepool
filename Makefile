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
	@go test -v --tags=for_select_go_work --bench=.
