package main

import (
	"errors"
	"strings"
)

type Implementation interface {
	NextTurn()
	GetWorld() World
	Close()
}

type ImplementationEnum int

type ImplementationInitFn func(world World, threads int) Implementation

// An enum that represents the different GoL implementations
const (
	ImplementationSerial = iota
	ImplementationParallel
	ImplementationParallelShared
	ImplementationHalo
	ImplementationRust
)

const ImplementationDefault ImplementationEnum = ImplementationRust

// This is case insensitive
func implementationFromName(s string) (ImplementationEnum, error) {
	s = strings.ToLower(s)
	switch s {
	case "serial":
		return ImplementationSerial, nil
	case "parallel":
		return ImplementationParallel, nil
	case "parallelshared":
		return ImplementationParallelShared, nil
	case "halo":
		return ImplementationHalo, nil
	case "rust":
		return ImplementationRust, nil
	default:
		return 0, errors.New("invalid implementation name")
	}
}

func (i ImplementationEnum) initFn() ImplementationInitFn {
	switch i {
	case ImplementationSerial:
		return InitSerial
	case ImplementationParallel:
		return InitParallel
	case ImplementationParallelShared:
		return InitParallelShared
	case ImplementationRust:
		return InitRust
	}
	panic("unmatched case")
}

func (i ImplementationEnum) name() string {
	switch i {
	case ImplementationSerial:
		return "Serial"
	case ImplementationParallel:
		return "Parallel"
	case ImplementationParallelShared:
		return "ParallelShared"
	case ImplementationHalo:
		return "Halo"
	case ImplementationRust:
		return "Rust"
	}
	panic("unmatched case")
}
