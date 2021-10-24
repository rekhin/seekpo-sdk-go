package snmp

import "github.com/rekhin/seekpo-sdk-go"

type Tag struct {
	seekpo.Tag
	Oid     string
	Command Command
}

type Command uint8

const (
	CommangSnmpGet Command = iota
	CommandSnmpWalk
	CommandSnmpBulkWalk
)
