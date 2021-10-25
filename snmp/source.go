package snmp

import "github.com/rekhin/seekpo-sdk-go"

type Source struct {
	seekpo.Source
	Uri     string
	Version Version
}

type Version uint8

const (
	Version1 Version = iota
	Version2c
	Version3
)
