#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -o go-node-out/app -ldflags -s -a -installsuffix cgo ./go-node
