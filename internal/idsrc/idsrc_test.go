package idsrc

import (
	"sync"
	"testing"
	"time"
)

func TestGenCurrentTime(t *testing.T) {
	src := New()
	now := (time.Now().UnixMilli() - epochOffset) << cntOffset
	if id, err := src.Gen(); err != nil || id < now {
		t.Fail()
	}
}

func TestGenUnique(t *testing.T) {
	src := New()
	var wg sync.WaitGroup
	const routines, iterations = 100, 2000
	var a [routines][iterations]int64
	wg.Add(routines)
	for routine := range routines {
		go func(routine int) {
			defer wg.Done()
			for iteration := range iterations {
				a[routine][iteration], _ = src.Gen()
				// ignore errors, those will pop up as zeros later
			}
		}(routine)
	}
	wg.Wait()
	m := make(map[int64]bool, routines*iterations)
	for routine := range routines {
		for iteration := range iterations {
			id := a[routine][iteration]
			if id == 0 || m[id] {
				println(id, m[id])
				t.FailNow() // error or duplicate
			} else {
				m[id] = true
			}
		}
	}
}
