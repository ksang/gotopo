package lldp

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/soniah/gosnmp"
)

var (
	// OIDChassisID is exact OID to get switch ChassisId
	OIDChassisID = "1.0.8802.1.1.2.1.3.2.0"
	// OIDName is exact OID to get switch Name
	OIDName = "1.0.8802.1.1.2.1.3.3.0"
	// OIDDescription is exact OID to get switch Description
	OIDDescription = "1.0.8802.1.1.2.1.3.4.0"
	// OIDLocalPortTableNumber is root OID to get a subtree of port number
	OIDLocalPortTableNumber = "1.0.8802.1.1.2.1.3.7.1.1"
	// OIDLocalPortTableID is root OID to get a subtree of port id
	OIDLocalPortTableID = "1.0.8802.1.1.2.1.3.7.1.3"
	// OIDLocalPortTableDesc is root OID to get a subtree of port description
	OIDLocalPortTableDesc = "1.0.8802.1.1.2.1.3.7.1.4"
)

/*
SNMPEndpoint is an LLDP endpoint using SNMP
*/
type SNMPEndpoint struct {
	Address string
	Port    uint16
	// sleep time interval after collecting data from endpoint
	Interval  time.Duration
	Community string
	// SNMP version used, import "github.com/soniah/gosnmp" type SnmpVersion
	// default v2c
	Version gosnmp.SnmpVersion
	Timeout time.Duration
	Retries int
	MaxOids int
	// private members
	snmp   *gosnmp.GoSNMP
	dataCh chan *Snapshot
	errCh  chan error
	quitCh chan struct{}
}

// NewDefaultSNMPEndpoint create SNMPEndpiont with default configurations
func NewDefaultSNMPEndpoint(addr string) Endpoint {
	ep := &SNMPEndpoint{
		Address:   addr,
		Port:      161,
		Interval:  60 * time.Second,
		Community: "public",
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   3,
		MaxOids:   60,
	}
	return NewSNMPEndpoint(ep)
}

// NewSNMPEndpoint create SNMPEndpiont with provided configurations
func NewSNMPEndpoint(ep *SNMPEndpoint) Endpoint {
	ep.snmp = &gosnmp.GoSNMP{
		Target:    ep.Address,
		Port:      ep.Port,
		Community: ep.Community,
		Version:   ep.Version,
		Timeout:   ep.Timeout,
		Retries:   ep.Retries,
		MaxOids:   ep.MaxOids,
	}
	if ep.Interval == 0 {
		ep.Interval = 10 * time.Second
	}
	ep.dataCh = make(chan *Snapshot, 1)
	ep.errCh = make(chan error, 1)
	ep.quitCh = make(chan struct{}, 1)
	return ep
}

// Start SNMPEndpoint collecting data snapshots
func (s *SNMPEndpoint) Start() (chan *Snapshot, chan error) {
	go func() {
		defer func() {
			close(s.errCh)
			close(s.dataCh)
			s.snmp.Conn.Close()
		}()
		err := s.snmp.Connect()
		if err != nil {
			s.errCh <- err
			return
		}
		for {
			select {
			case <-time.After(s.Interval):
				ss, err := s.getSnapshot()
				if err != nil {
					s.errCh <- err
					continue
				}
				s.dataCh <- ss
			case <-s.quitCh:
				s.errCh <- errors.New("stopped")
				return
			}
		}
	}()
	return s.dataCh, s.errCh
}

// Stop the SNMPEndpoint
func (s *SNMPEndpoint) Stop() error {
	s.quitCh <- struct{}{}
	return nil
}

func (s *SNMPEndpoint) getSnapshot() (*Snapshot, error) {
	exactOids := []string{
		OIDChassisID,
		OIDName,
		OIDDescription,
	}
	var (
		chassisID string
		name      string
		desc      string
	)
	res1, err := s.snmp.Get(exactOids)
	if err != nil {
		return nil, err
	}
	for i, v := range res1.Variables {
		if v.Type == gosnmp.OctetString {
			switch i {
			case 0:
				chassisID = hex.EncodeToString(v.Value.([]byte))
			case 1:
				name = string(v.Value.([]byte))
			case 2:
				desc = string(v.Value.([]byte))
			}
		}
	}
	ret := &Snapshot{
		Local: LocalSystemData{
			ChassisID:   chassisID,
			Name:        name,
			Description: desc,
		},
	}
	res2, err := s.snmp.WalkAll(OIDLocalPortTableNumber)
	if err != nil {
		return ret, err
	}
	localPortTable := []PortTableEntry{}
	for _, v := range res2 {
		if v.Type == gosnmp.Integer {
			localPortTable = append(localPortTable, PortTableEntry{Number: v.Value.(int)})
		}
	}
	ret.Local.PortTable = localPortTable
	res3, err := s.snmp.WalkAll(OIDLocalPortTableID)
	if err != nil {
		return ret, err
	}
	for i, v := range res3 {
		if i >= len(localPortTable) {
			break
		}
		if v.Type == gosnmp.OctetString {
			localPortTable[i].ID = string(v.Value.([]byte))
		}
	}
	res4, err := s.snmp.WalkAll(OIDLocalPortTableDesc)
	if err != nil {
		return ret, err
	}
	for i, v := range res4 {
		if i >= len(localPortTable) {
			break
		}
		if v.Type == gosnmp.OctetString {
			localPortTable[i].Description = string(v.Value.([]byte))
		}
	}
	return ret, nil
}
