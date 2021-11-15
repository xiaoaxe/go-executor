//worker
//@author: baoqiang
//@time: 2021/11/15 20:20:13
package threadpool

import "fmt"

type Worker struct {
	name     string
	workPool chan chan interface{}
	stopChan chan struct{}

	// internal vars
	jobChan chan interface{}
}

func NewWorker(name string, workPool chan chan interface{}, stopChan chan struct{}) *Worker {
	return &Worker{
		name:     name,
		workPool: workPool,
		jobChan:  make(chan interface{}),
		stopChan: stopChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// return self
			w.workPool <- w.jobChan
			select {
			case job := <-w.jobChan:
				w.Do(job)
			case <-w.stopChan:
				return
			}
		}
	}()
}

func (w *Worker) Do(task interface{}) {
	switch f := task.(type) {
	case Runnable:
		fmt.Printf("run in worker: %v\n", w.name)
		f.Run()
	case *callableTask:
		fmt.Printf("run in worker: %v\n", w.name)
		res := f.Task.Call(f.Task.GetArg())
		f.Resp.done = true
		f.Resp.response <- res
	default:
		fmt.Printf("invalid input type: type=%v, task=%v\n", f, task)
	}
}
