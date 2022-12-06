package dev

import (
	"strings"
	"time"

	"github.com/alouca/gosnmp"
)

const (
	snmpBulkMax = 64
	snmpTimeout = 5

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

func NewSwitch(target string, community string) (*Switch, error) {
	ip, ports := parseTarget(target)

	snmpInterface, err := gosnmp.NewGoSNMP(ip, community, gosnmp.Version2c, snmpTimeout)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if len(res.Variables) > 0 {
		sw.Model = strings.ToUpper(res.Variables[0].Value.(string))
	}

	return sw, nil
}

func (sw *Switch) updatePPS(oidTx string, oidRx string) {

	pRx, err := sw.Snmp.BulkWalk(snmpBulkMax, oidRx)
	if err != nil {
		return
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
		return
	}
	for idx, pdu := range pTx {
		if sw.PortsRangeMon == "" || PortIsInRange(idx+1, sw.PortsRangeMon) {
			val := pdu.Value.(int)
			sw.PPSinfo.TX[idx+1] = val
		} else {
			sw.PPSinfo.TX[idx+1] = -1
		}
	}
}

func (sw *Switch) StartMonPPS(td time.Duration, oidTx string, oidRx string) {
	go func() {
		tickTime := time.NewTicker(td)
		for {
			sw.updatePPS(oidTx, oidRx)
			<-tickTime.C
		}
	}()
}
