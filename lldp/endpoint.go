package lldp

/*
Endpoint is the abtraction of a LLDP endpoint
Typically it is a SNMP agent
*/
type Endpoint interface {
	Start() (chan *Snapshot, chan error)
	Stop() error
}
