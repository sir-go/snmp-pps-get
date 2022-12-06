# CLI tool for monitoring PPS on the D-Link switches by ports
[![Go](https://github.com/sir-go/snmp-pps-get/actions/workflows/go.yml/badge.svg)](https://github.com/sir-go/snmp-pps-get/actions/workflows/go.yml)

CLI utility periodically gets packets-per-second statistic of given devices and print it as a table.

## Build
```bash
go test ./...
gosec ./...
go mod download && go build -o pps ./cmd/cli
```

## Flags
`<ip0> <ip1>:<posts> ... ` - target switches IP addresses, can be an IP or IP:ports-range, i.e.:
`10.10.0.1 10.15.0.12:1-5 192.168.28.26:14,17,20-25`

## Config
`config.yml` must be located in the same directory. It contains the array of SNMP OIDs for getting metrics
(TX and RX PPS) by the switch model name and default values.

## Usage
```bash
pps <ip0> <ip1>:<ports> ... <ipN>
```

![](pps.gif)
