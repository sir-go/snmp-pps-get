package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/alouca/gosnmp"
)

const (
	snmpBulkMax   = 64
	snmpCommunity = "public"
	snmpTimeout   = 5

	OIDsysDescr = ".1.3.6.1.2.1.1.1.0"
)

type PPSinfo struct {
	TX, RX []int
}

type Switch struct {
	Ip            string
	Model         string
	Snmp          *gosnmp.GoSNMP
	PPSinfo       *PPSinfo
	PortsAmount   int
	PortsRangeMon string
}

func NewSwitch(addrPorts string) *Switch {
	var ip, ports string

	if strings.ContainsRune(addrPorts, ':') {
		addr := strings.SplitN(addrPorts, ":", 3)
		ip = addr[0]
		ports = addr[1]
	} else {
		ip = addrPorts
	}

	snmpInterface, err := gosnmp.NewGoSNMP(ip, snmpCommunity, gosnmp.Version2c, snmpTimeout)
	if err != nil {
		log.Error(err)
	}

	sw := &Switch{
		Ip:            ip,
		Snmp:          snmpInterface,
		PPSinfo:       &PPSinfo{},
		PortsAmount:   -1,
		PortsRangeMon: ports,
	}

	res, err := snmpInterface.Get(OIDsysDescr)
	if err != nil {
		log.Fatal("snmp get", sw.Ip, "error", err.Error())
	}
	if len(res.Variables) > 0 {
		sw.Model = strings.ToUpper(res.Variables[0].Value.(string))
		log.Info(sw.Ip, sw.Model)
	}

	return sw
}

func (sw *Switch) getPortsPPS(oidTx string, oidRx string) error {

	pRx, err := sw.Snmp.BulkWalk(snmpBulkMax, oidRx)
	if err != nil {
		return err
	}

	if sw.PortsAmount < 0 {
		sw.PortsAmount = len(pRx)
	}

	sw.PPSinfo.RX = make([]int, sw.PortsAmount+1)
	sw.PPSinfo.TX = make([]int, sw.PortsAmount+1)

	for idx, pdu := range pRx {
		if sw.PortsRangeMon == "" || PortIsInRange(idx+1, sw.PortsRangeMon) {
			sw.PPSinfo.RX[idx+1] = pdu.Value.(int)
		} else {
			sw.PPSinfo.RX[idx+1] = -1
		}
	}

	pTx, err := sw.Snmp.BulkWalk(snmpBulkMax, oidTx)
	if err != nil {
		return err
	}
	for idx, pdu := range pTx {
		if sw.PortsRangeMon == "" || PortIsInRange(idx+1, sw.PortsRangeMon) {
			val := pdu.Value.(int)
			sw.PPSinfo.TX[idx+1] = val
		} else {
			sw.PPSinfo.TX[idx+1] = -1
		}
	}

	return nil
}

func (sw *Switch) StartMonPPS(td time.Duration, oidTx string, oidRx string) {
	go func() {
		tickTime := time.NewTicker(td)
		for {
			if err := sw.getPortsPPS(oidTx, oidRx); err != nil {
				log.Error(err)
			}
			<-tickTime.C
		}
	}()
}

// PortIsInRange (13, "4,6,12-15") -> true
func PortIsInRange(pNum int, pRange string) bool {

	if pNum < 1 || len(pRange) < 1 {
		return false
	}

	// "4,6,12-15" -> "4", "6", "12-15"
	for _, part := range strings.Split(pRange, ",") {
		if !strings.ContainsRune(part, '-') {
			// "4", "6"
			if digit, err := strconv.Atoi(part); err == nil {
				if digit == pNum {
					return true
				}
			}
		} else {
			//"12-15" | "15-12" | "15-" | "-12"
			edges := strings.SplitN(part, "-", 3)

			e0, _ := strconv.Atoi(edges[0])
			e1, _ := strconv.Atoi(edges[1])

			if e0 == e1 {
				if e0 == 0 {
					// "_-_" ? -> next part
					continue
				}
				if e0 == pNum {
					// "15-15"; pNum == 15 ? else -> next part
					return true
				} else {
					continue
				}
			}

			// ........p..........e1
			if e0 == 0 {
				if pNum <= e1 {
					return true
				} else {
					continue
				}
			}

			// e0.......p..........
			if e1 == 0 {
				if pNum >= e0 {
					return true
				} else {
					continue
				}
			}

			// swap values if reversed "15-12" -> "12-15"
			if e0 > e1 {
				e0, e1 = e1, e0
			}

			// e0......p.........e1
			if pNum <= e1 && pNum >= e0 {
				return true
			}
		}
	}

	return false
}
