package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/op/go-logging"
)

const (
	//logFormat = `%{time: Jan/02 15:04:05.000} %{level:.1s} %{shortfunc} > %{message}`
	logFormat     = `%{time: 15:04:05.000} : %{message}`
	updatePeriod  = "1200ms"
	requestPeriod = "1500ms"

	confFile = "/opt/pps_conf.json"
)

var (
	log     *logging.Logger
	targets []*Switch
)

func init() {
	log = logging.MustGetLogger("gsnmp")
	formatter := logging.MustStringFormatter(logFormat)
	lb := logging.NewLogBackend(os.Stdout, "", 0)
	lbf := logging.NewBackendFormatter(lb, formatter)
	lbl := logging.AddModuleLevel(lbf)
	logging.SetBackend(lbl)
}

func interaptHandler(c chan os.Signal) {
	for range c {
		log.Info("-- stop --")
		os.Exit(137)
	}
}

func refreshOut(updatePeriod time.Duration) {
	var (
		outStr   string
		strCount int
	)

	tickTime := time.NewTicker(updatePeriod)
	for {
		<-tickTime.C

		outStr = "\n"
		strCount = 0

		for _, target := range targets {
			if strCount < target.PortsAmount {
				strCount = target.PortsAmount + 1
			}

			outStr += fmt.Sprintf("| %-19s", target.Ip)
		}

		fmt.Println(outStr)
		outStr = ""

		for i := 1; i < strCount; i++ {
			for _, target := range targets {
				if i <= target.PortsAmount {
					if target.PPSinfo.RX[i] != -1 {
						outStr += fmt.Sprintf("| %2d: %6d  %6d ", i, target.PPSinfo.RX[i], target.PPSinfo.TX[i])
					} else {
						outStr += fmt.Sprintf("| %2d: %6s  %6s ", i, "-", "-")
					}
				} else {
					outStr += fmt.Sprintf("%21s", " ")
				}
			}
			fmt.Println(outStr)
			outStr = ""
		}
	}
}

func getOids(sw *Switch, cfg *Cfg) (oidTx string, oidRx string) {
	for _, oidCfg := range *cfg.Oids {
		for _, model := range *oidCfg.Models {
			if strings.Contains(sw.Model, model) {
				return oidCfg.Tx, oidCfg.Rx
			}
		}
	}
	return cfg.Default.Tx, cfg.Default.Rx
}

func main() {
	/*
		log.Info(nswitch.PortIsInRange(2, "4,6,12-15,18-") == false)
		log.Info(nswitch.PortIsInRange(4, "4,6,12-15,18-") == true)
		log.Info(nswitch.PortIsInRange(18, "4,6,12-15,18-") == true)
		log.Info(nswitch.PortIsInRange(17, "4,6,12-15,18-") == false)
		log.Info(nswitch.PortIsInRange(25, "4,6,12-15,18-") == true)
		log.Info(nswitch.PortIsInRange(4, "2-15") == true)

		log.Info(nswitch.PortIsInRange(2, "4,6,-18,12-15") == true)
		log.Info(nswitch.PortIsInRange(4, "4,6,-18,12-15") == true)
		log.Info(nswitch.PortIsInRange(18, "4,6,-18,12-15") == true)
		log.Info(nswitch.PortIsInRange(17, "4,6,-18,12-15") == true)
		log.Info(nswitch.PortIsInRange(4, "1,6,8-10") == false)
	*/

	if len(os.Args) == 1 || os.Args[1] == "-h" {
		log.Info("usage: pps <ip0> <ip1> <ip2> ...")
		log.Infof("config file: %s", confFile)
		os.Exit(0)
	}

	cfg, err := LoadConfig(confFile)
	if err != nil {
		log.Panic(err)
	}

	log.Info("-- start --")
	log.Infof("requests timeout: %s, out refresh timeout: %s", requestPeriod, updatePeriod)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go interaptHandler(c)

	tupd, err := time.ParseDuration(updatePeriod)
	if err != nil {
		log.Error(err)
	}

	treq, err := time.ParseDuration(requestPeriod)
	if err != nil {
		log.Error(err)
	}

	var oidTx, oidRx string

	for _, ip := range os.Args[1:] {
		sw := NewSwitch(ip)
		oidTx, oidRx = getOids(sw, cfg)
		sw.StartMonPPS(treq, oidTx, oidRx)
		targets = append(targets, sw)
	}
	refreshOut(tupd)
}
