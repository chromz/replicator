SHELL := /bin/bash

REPLICATOR_SERVER := cmd/rserver/rserver.go
REPLICATOR_CLIENT := cmd/rclient/rclient.go
.PHONY: all

all:
	@go build $(REPLICATOR_SERVER)
