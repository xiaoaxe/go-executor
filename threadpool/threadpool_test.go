//test run
//@author: baoqiang
//@time: 2021/11/15 20:18:28
package threadpool

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestThreadPool1(t *testing.T) {
	var (
		num = 10
	)

	p := NewThreadPool(3, num)
	for i := 0; i < num; i++ {
		task := &HelloPrinter{name: strconv.Itoa(i)}
		p.Submit(task)
	}

	// wait finish
	time.Sleep(time.Second)

	p.Stop()
}

func TestThreadPool2(t *testing.T) {
	var (
		num = 10
	)

	p := NewThreadPool(3, num)
	fs := make([]*Future, 0, num) //fuck
	for i := 0; i < num; i++ {
		task := &AddOneCal{name: strconv.Itoa(i), arg: i}
		f := p.ExecuteFuture(task)
		fs = append(fs, f)
	}

	// wait finish
	time.Sleep(time.Second)

	// must wait all done
	for _, f := range fs {
		fmt.Printf("Got resp: %v\n", f.Get())
	}

	p.Stop()
}

// runner
type HelloPrinter struct {
	name string
}

func (p *HelloPrinter) Run() {
	fmt.Printf("my name is %s, current is: %v\n", p.name, time.Now().UnixNano())
}

// callable
type AddOneCal struct {
	name string
	arg  int
}

func (a *AddOneCal) Call(obj interface{}) interface{} {
	i := obj.(int)
	fmt.Printf("my name is %s, input: %v, output: %v\n", a.name, i, i+1)
	return i + 1
}

func (a *AddOneCal) GetArg() interface{} {
	return a.arg
}
