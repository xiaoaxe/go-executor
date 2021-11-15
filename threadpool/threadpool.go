//threadpool impl
//@author: baoqiang
//@time: 2021/11/15 20:18:15
package threadpool

import (
	"strconv"
)

type TreadPool struct {
	numTaskers int
	jobSize    int

	// internal vars
	workerPool chan chan interface{}
	jobQueue   chan interface{}
	stopChan   chan struct{}
}

func NewThreadPool(numTaskers int, jobSize int) *TreadPool {
	p := &TreadPool{
		numTaskers: numTaskers,
		jobSize:    jobSize,
		workerPool: make(chan chan interface{}, numTaskers),
		jobQueue:   make(chan interface{}, jobSize),
		stopChan:   make(chan struct{}),
	}

	p.createPool()
	return p
}

func (p *TreadPool) Submit(r Runnable) {
	p.jobQueue <- r
	return
}

func (p *TreadPool) ExecuteFuture(c Callable) *Future {
	f := &Future{
		response: make(chan interface{}),
	}
	t := &callableTask{
		Task: c,
		Resp: f,
	}
	p.jobQueue <- t
	return f
}

func (p *TreadPool) Stop() {
	close(p.jobQueue)
	close(p.workerPool)
	close(p.stopChan)
	return
}

// inner impl
func (p *TreadPool) createPool() {
	for i := 0; i < p.numTaskers; i++ {
		w := NewWorker(strconv.Itoa(i), p.workerPool, p.stopChan)
		w.Start()
	}

	go p.dispatch()
}

func (p *TreadPool) dispatch() {
	for {
		select {
		case j := <-p.jobQueue:
			// get a free worker
			jq := <-p.workerPool
			// send task
			jq <- j
		case <-p.stopChan:
			return
		}
	}
}
