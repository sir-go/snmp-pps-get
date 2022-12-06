package main

import (
	"io/ioutil"
	"os"
	"time"

	"snmp-pps-get/internal/dev"
)

const (
	updatePeriod  = 1200 * time.Millisecond
	requestPeriod = 1500 * time.Millisecond
)

func main() {
	confBytes, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	cfg, err := LoadConfig(confBytes)
	if err != nil {
		panic(err)
	}

	var oidTx, oidRx string
	targets := make([]dev.Switch, 0)

	for _, ip := range os.Args[1:] {
		sw, err := dev.NewSwitch(ip, cfg.SnmpCommunity)
		if err != nil {
			panic(err)
		}
		oidTx, oidRx = getOids(sw.Model, cfg)
		sw.StartMonPPS(requestPeriod, oidTx, oidRx)
		targets = append(targets, *sw)
	}
	refreshOut(updatePeriod, targets)
}
