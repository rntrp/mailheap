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
	for i := 0; i < routines; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				a[i][j], _ = src.Gen()
				// ignore errors, those will pop up as zeros later
			}
		}()
	}
	wg.Wait()
	m := make(map[int64]bool, routines*iterations)
	for i := 0; i < routines; i++ {
		for j := 0; j < iterations; j++ {
			id := a[i][j]
			if id == 0 || m[id] {
				println(id, m[id])
				t.FailNow() // error or duplicate
			} else {
				m[id] = true
			}
		}
	}
}
