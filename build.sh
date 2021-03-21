#!/bin/bash

env GOOS=linux GOARCH=amd64 go build -v -o bin/binance_exporter-linux-amd64 binance_exporter.go util.go
env GOOS=windows GOARCH=amd64 go build -v -o bin/binance_exporter-windows-amd64 binance_exporter.go util.go
env GOOS=darwin GOARCH=amd64 go build -v -o bin/binance_exporter-darwin-amd64 binance_exporter.go util.go