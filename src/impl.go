package main

import (
	"errors"
	"strings"
)

type Implementation int

// An enum that represents the different GoL implementations
const (
	ImplementationSerial = iota
)

const ImplementationDefault Implementation = ImplementationSerial

// This is case insensitive
func implementationFromString(s string) (Implementation, error) {
	s = strings.ToLower(s)
	switch s {
	case "serial":
		return ImplementationSerial, nil
	default:
		return 0, errors.New("invalid implementation string")
	}
}

func (i Implementation) function() func(*World) {
	switch i {
	case ImplementationSerial:
		return updateWorldSerially
	}
	panic("unmatched case")
}

func (i Implementation) toString() string {
	switch i {
	case ImplementationSerial:
		return "Serial"
	}
	panic("unmatched case")
}
