#!/bin/bash

# go get github.com/golang/mock/mockgen
# its possible to use //go:generate as well but adding it here for demonstration pursposes
mockgen -destination=./repository/mock_repository/repository.go -package=mock_repository github.com/eran-levy/tokenizer-gophercon/repository Persistence
mockgen -destination=./cache/mock_cache/cache.go -package=mock_cache github.com/eran-levy/tokenizer-gophercon/cache Cache

