SHELL := /bin/bash

REPLICATOR_CLIENT := cmd/rclient/rclient.go
REPLICATOR_SERVER := cmd/rserver/rserver.go
.PHONY: all

all:
	@go build $(REPLICATOR_CLIENT)
	@go build $(REPLICATOR_SERVER)

rclient: .FORCECLIENT

.FORCECLIENT:
	@go build $(REPLICATOR_CLIENT)

rserver: .FORCESERVER

.FORCESERVER:
	@go build $(REPLICATOR_SERVER)
