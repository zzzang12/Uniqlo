package waitgroupcount

import (
	"sync"
	"sync/atomic"
)

type WaitGroupCount struct {
	sync.WaitGroup
	//count int64
	count atomic.Int64
}

func (wg *WaitGroupCount) Add(delta int) {
	//atomic.AddInt64(&wg.count, int64(delta))
	wg.count.Add(int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroupCount) Done() {
	//atomic.AddInt64(&wg.count, -1)
	wg.count.Add(-1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroupCount) Wait() {
	wg.WaitGroup.Wait()
}

func (wg *WaitGroupCount) GetCount() int {
	//return int(atomic.LoadInt64(&wg.count))
	return int(wg.count.Load())
}
