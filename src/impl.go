package main

import (
	"errors"
	"strings"
)

type Implementation interface {
	Init(world World, threads int)
	NextTurn()
	GetWorld() World
}

type ImplementationEnum int

// An enum that represents the different GoL implementations
const (
	ImplementationSerial = iota
	ImplementationParallel
	ImplementationParallelShared
	ImplementationHalo
)

const ImplementationDefault ImplementationEnum = ImplementationParallel

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
	default:
		return 0, errors.New("invalid implementation name")
	}
}

func (i ImplementationEnum) new() Implementation {
	switch i {
	case ImplementationSerial:
		return &Serial{}
	case ImplementationParallel:
		return &Parallel{}
	case ImplementationParallelShared:
		return &ParallelShared{}
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
	}
	panic("unmatched case")
}
