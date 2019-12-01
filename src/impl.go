package main

import (
	"errors"
	"strings"
)

type implementation interface {
	nextTurn()
	getWorld() world
	close()
}

type implementationEnum int

type implementationInitFn func(world world, threads int) implementation

// An enum that represents the different GoL implementations
const (
	implementationSerial = iota
	implementationParallel
	implementationParallelShared
	implementationHalo
	implementationRust
)

const implementationDefault implementationEnum = implementationHalo

// This is case insensitive
func implementationFromName(s string) (implementationEnum, error) {
	s = strings.ToLower(s)
	switch s {
	case "serial":
		return implementationSerial, nil
	case "parallel":
		return implementationParallel, nil
	case "parallelshared":
		return implementationParallelShared, nil
	case "halo":
		return implementationHalo, nil
	case "rust":
		return implementationRust, nil
	default:
		return 0, errors.New("invalid implementation name")
	}
}

func (i implementationEnum) initFn() implementationInitFn {
	switch i {
	case implementationSerial:
		return initSerial
	case implementationParallel:
		return initParallel
	case implementationParallelShared:
		return initParallelShared
	case implementationHalo:
		return initHalo
	case implementationRust:
		return initRust
	}
	panic("unmatched case")
}

func (i implementationEnum) name() string {
	switch i {
	case implementationSerial:
		return "Serial"
	case implementationParallel:
		return "Parallel"
	case implementationParallelShared:
		return "ParallelShared"
	case implementationHalo:
		return "Halo"
	case implementationRust:
		return "Rust"
	}
	panic("unmatched case")
}
