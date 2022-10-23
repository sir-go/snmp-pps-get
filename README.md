## CLI tool for monitoring PPS on the D-Link switches by ports

### Build

```bash
go mod download && go build -o pps ./cmd/cli
```

### Flags
`<ip0> <ip1> ... ` - target switches IP addresses

### Config
`pps_conf.json` must be located in the same directory. It contains the array of SNMP OIDs for getting metrics
(TX and RX PPS) by the switch model name and default values.

The SNMP community `public` is hardcoded into `switch.go` as constant (todo: make as variable)

### Usage
```bash
pps <ip0> <ip1> ... <ipN>
```

![](pps.gif)
