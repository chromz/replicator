SHELL := /bin/bash

REPLICATOR_CLIENT := cmd/rclient/rclient.go
.PHONY: all

all:
	@go build $(REPLICATOR_CLIENT)
