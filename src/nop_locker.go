package main

// A sync.Locker that has no operation
// Use to remove runtime overhead of synchronisation
type NopLocker struct{}

func (NopLocker) Lock()   {}
func (NopLocker) Unlock() {}
