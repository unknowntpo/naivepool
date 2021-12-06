all: help 

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## bench/naivepool: benchmark naivepool
.PHONY: bench/naivepool
bench/naivepool:
	@go test -v --bench=. -benchmem
