//callable
//@author: baoqiang
//@time: 2021/11/15 20:50:52
package threadpool

type Callable interface {
	Call(interface{}) interface{}
	GetArg() interface{}
}

type Future struct {
	response chan interface{}
	done     bool
}

func (f *Future) IsDone() bool {
	return f.done
}

func (f *Future) Get() interface{} {
	return <-f.response
}

type callableTask struct {
	Task Callable
	Resp *Future
}
