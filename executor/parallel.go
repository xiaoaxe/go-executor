//parallel do
//@author: baoqiang
//@time: 2021/11/09 21:07:20
package executor

import "sync"

func ParallelDo(fs ...func()) {
	var wg sync.WaitGroup

	for _, f := range fs {
		wg.Add(1)
		go func(g func()) {
			defer wg.Done()
			g()
		}(f)
	}

	wg.Wait()
}
