package main

import (
	"fmt"
	"time"

	"snmp-pps-get/internal/dev"
)

func refreshOut(updatePeriod time.Duration, targets []dev.Switch) {
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
