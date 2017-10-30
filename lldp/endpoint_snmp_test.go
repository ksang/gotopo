package lldp

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/soniah/gosnmp"
)

var addr = flag.String("gss", "", "Enable GetSnapshot testing towards a SNMP server, provide server address")

func TestNewSNMPEndpoint(t *testing.T) {
	ep1 := NewDefaultSNMPEndpoint("192.168.1.1")
	t.Logf("NewDefaultSNMPEndpoint, ep1: %#v", ep1)
	ep2 := &SNMPEndpoint{
		Address:   "1.2.3.4",
		Port:      161,
		Community: "public",
		Version:   gosnmp.Version2c,
		Timeout:   3 * time.Second,
		Retries:   3,
		MaxOids:   60,
	}
	ep3 := NewSNMPEndpoint(ep2)
	t.Logf("NewSNMPEndpoint, ep2: %#v, ep3:%#v", ep2, ep3)
}

func TestGetSnapshot(t *testing.T) {
	if len(*addr) <= 0 {
		return
	}
	t.Logf("SNMP server address: %s", *addr)
	ep := NewDefaultSNMPEndpoint(*addr)
	t.Logf("NewDefaultSNMPEndpoint, ep: %#v", ep)
	dataCh, errCh := ep.Start()
	select {
	case ss := <-dataCh:
		fmt.Printf("Snapshot.Local: ChassisId: %s, Name: %s, Description: %s\n",
			ss.Local.ChassisID,
			ss.Local.Name,
			ss.Local.Description,
		)
		fmt.Printf("LocalPortTable:\n")
		for _, entry := range ss.Local.PortTable {
			fmt.Printf("\t%#v\n", entry)
		}
	case err := <-errCh:
		t.Logf("GetSnapshot error: %v", err)
	}
	if err := ep.Stop(); err != nil {
		t.Error("failed to stop endpoint")
	}
}
