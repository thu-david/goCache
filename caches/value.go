package caches

import (
	"goCache/helpers"
	"sync/atomic"
	"time"
)

const (
	NeverDie = 0
)

type value struct {
	data []byte
	ttl int64
	ctime int64	
}

func newValue(data []byte, ttl int64) *value {
	return &value{
		data: helpers.Copy(data),
		ttl : ttl,
		ctime: time.Now().Unix(),
	}
}

func (v *value) alive() bool {
	return v.ttl == NeverDie || time.Now().Unix() - v.ctime < v.ttl
}

func (v *value) visit() []byte {
	atomic.SwapInt64(&v.ctime, time.Now().Unix())
	return v.data
}