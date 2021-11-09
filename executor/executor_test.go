// run test
//@author: baoqiang
//@time: 2021/11/09 20:53:24
package executor

import (
	"context"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func Test_Executor(t *testing.T) {
	var (
		ctx = context.Background()
		// fetchers
		one   = NewOneFetcher()
		two   = NewTwoFetcher()
		three = NewThreeFetcher(one, two)
		np    = NewNumberProcessor(three)
	)

	e := NewExecutor()

	// build
	e.AppendProcessor(np)

	e.Exec(ctx)
}

// helpers
// process
type NumberProcessor struct {
	three *ThreeFetcher
}

func NewNumberProcessor(three *ThreeFetcher) *NumberProcessor {
	return &NumberProcessor{
		three: three,
	}
}

func (p *NumberProcessor) Name() string {
	return reflect.TypeOf(p).Elem().Name()
}

func (p *NumberProcessor) Process(ctx context.Context) error {
	var final = p.three.GetData().(int64)
	log.Printf("Got Final Data: %v\n", final)
	return nil
}

func (p *NumberProcessor) GetFetchers() []IFetcher {
	return []IFetcher{
		p.three,
	}
}

// fetchers
// ONE
type OneFetcher struct {
	// internal vars
	data int64
}

func NewOneFetcher() *OneFetcher {
	return &OneFetcher{}
}
func (f *OneFetcher) Fetch(ctx context.Context) error {
	log.Printf("fetcher run: %v", f.Name())
	Sleep()
	f.data = 1
	return nil
}
func (f *OneFetcher) DependFetchers() []IFetcher {
	return []IFetcher{}
}
func (f *OneFetcher) Name() string {
	return reflect.TypeOf(f).Elem().Name()
}
func (f *OneFetcher) GetData() interface{} {
	return f.data
}

//TWO
type TwoFetcher struct {
	// internal vars
	data int64
}

func NewTwoFetcher() *TwoFetcher {
	return &TwoFetcher{}
}
func (f *TwoFetcher) Fetch(ctx context.Context) error {
	log.Printf("fetcher run: %v", f.Name())
	Sleep()
	f.data = 2
	return nil
}
func (f *TwoFetcher) DependFetchers() []IFetcher {
	return []IFetcher{}
}
func (f *TwoFetcher) Name() string {
	return reflect.TypeOf(f).Elem().Name()
}
func (f *TwoFetcher) GetData() interface{} {
	return f.data
}

//THREE
type ThreeFetcher struct {
	one *OneFetcher
	two *TwoFetcher

	// internal vars
	data int64
}

func NewThreeFetcher(one *OneFetcher, two *TwoFetcher) *ThreeFetcher {
	return &ThreeFetcher{
		one: one,
		two: two,
	}
}
func (f *ThreeFetcher) Fetch(ctx context.Context) error {
	log.Printf("fetcher run: %v", f.Name())
	Sleep()
	f.data = 3
	f.data += f.one.GetData().(int64) + f.two.GetData().(int64)
	return nil
}
func (f *ThreeFetcher) DependFetchers() []IFetcher {
	return []IFetcher{
		f.one,
		f.two,
	}
}
func (f *ThreeFetcher) Name() string {
	return reflect.TypeOf(f).Elem().Name()
}
func (f *ThreeFetcher) GetData() interface{} {
	return f.data
}

// utils
func Sleep() {
	// time.Sleep(time.Second * time.Duration(rand.Float64()))
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
}
