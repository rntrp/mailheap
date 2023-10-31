package idsrc

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrEpochExhausted = errors.New("the epoch is exhausted; cannot generate new IDs")
var ErrInvalidID = errors.New("this ID is invalid")

const epochOffset = int64(1_690_000_000_000)
const maxEpochMilli = int64(1_099_511_627_775) // upper 40 bits ≈ 35 years
const cntOffset = 8
const maxCnt = int64(255) // lower 8 bits

type IdSrc interface {
	Gen() (int64, error)
}

type idSrcData struct {
	mtx sync.Mutex
	t   atomic.Int64
	cnt atomic.Int64
}

func New() IdSrc {
	return new(idSrcData)
}

func Decode(id int64) (time.Time, uint8, error) {
	if t := id >> cntOffset; t >= 0 && t <= maxEpochMilli {
		return time.UnixMilli(t), uint8(id & maxCnt), nil
	}
	return time.Time{}, 0, ErrInvalidID
}

func (d *idSrcData) Gen() (int64, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	nowUnixMilli := time.Now().UnixMilli()
	now := nowUnixMilli - epochOffset
	if now > maxEpochMilli {
		return 0, ErrEpochExhausted
	}
	t := d.t.Load()
	if t < now {
		d.t.Store(now)
		d.cnt.Store(0)
		return now << cntOffset, nil
	}
	c := d.cnt.Load()
	if c < maxCnt {
		d.cnt.Add(1)
		return (t << cntOffset) | (c + 1), nil
	}
	for t >= now {
		nextUnixMilli := time.UnixMilli(nowUnixMilli + 1)
		nextUnixMilli = <-time.After(time.Until(nextUnixMilli)) // blocking
		now = nextUnixMilli.UnixMilli() - epochOffset
		if now > maxEpochMilli {
			return 0, ErrEpochExhausted
		}
	}
	d.t.Store(now)
	d.cnt.Store(0)
	return now << cntOffset, nil
}
