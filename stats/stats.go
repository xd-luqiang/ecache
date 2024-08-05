package stats

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/xd-luqiang/ecache"
)

var m sync.Map

type StatsNode struct {
	// don't reorder them or add field between them
	Evicted, Updated, Added, GetMiss, GetHit, DelMiss, DelHit uint64
}

// HitRate
func (s *StatsNode) HitRate() float64 {
	if s.GetHit == 0 && s.GetMiss == 0 {
		return 0.0
	}
	return float64(s.GetHit) / float64(s.GetHit+s.GetMiss)
}

// Bind - to stats a cache
// `pool` can be used to classify instances that store same items
// `caches` is cache instances to be binded
func Bind(pool string, caches ...*ecache.Cache) error {
	v, _ := m.LoadOrStore(pool, &StatsNode{})
	for _, c := range caches {
		c.Inspect(func(action int, _ string, _ *interface{}, _ []byte, status int) {
			// very, very, very low-cost for stats
			atomic.AddUint64((*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(v.(*StatsNode)))+uintptr(status+action*2-1)*unsafe.Sizeof(&status))), 1)
		})
	}
	return nil
}

// Stats - get the result like follows
//
// `k` is categoy, type is string
// `v` is node, type is `*stats.StatsNode`
//
//		stats.Stats().Range(func(k, v interface{}) bool {
//	    	fmt.Println("stats:", k, v)
//	    	return true
//		})
func Stats() *sync.Map {
	return &m
}
