// executor
//@author: baoqiang
//@time: 2021/11/09 20:47:31
package executor

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// log
func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rand.Seed(time.Now().UnixNano())
}

type Executor struct {
	preFetchers    []IFetcher
	processorChain []Processor

	// internal vars
	fetcherMap   map[string]IFetcher
	fetcherErrs  sync.Map // fetcher_name => error
	fetcherOnces map[string]*sync.Once
}

// exported funcs
func NewExecutor() *Executor {
	return &Executor{
		fetcherMap:   make(map[string]IFetcher),
		fetcherErrs:  sync.Map{},
		fetcherOnces: make(map[string]*sync.Once),
	}
}

func (e *Executor) Exec(ctx context.Context) error {
	e.iteratorAllFetcher()

	for name, fetcher := range e.fetcherMap {
		fmt.Printf("name: %v, fetcher: %v\n", name, structName(fetcher))
	}

	for _, p := range e.processorChain {
		e.executeProcessor(ctx, p)
	}

	log.Printf("run ok\n")
	return nil
}

func (e *Executor) AppendProcessor(p Processor) {
	e.processorChain = append(e.processorChain, p)
}

// internal funcs
func (e *Executor) iteratorAllFetcher() {
	for _, p := range e.processorChain {
		e.iteratorFetchers(p.GetFetchers())
	}
}

func (e *Executor) iteratorFetchers(fetchers []IFetcher) {
	for _, f := range fetchers {
		e.iteratorFetchers(f.DependFetchers())

		if old, found := e.fetcherMap[f.Name()]; !found {
			e.fetcherMap[f.Name()] = f
			e.fetcherOnces[f.Name()] = &sync.Once{}
		} else {
			if old != f {
				panic(fmt.Sprintf("diff fetchers with same name: %v, %v", old, f))
			}
		}
	}
}

func (e *Executor) executeProcessor(ctx context.Context, p Processor) error {
	// recursive get fetchers
	e.executeFetchers(ctx, p.GetFetchers())
	// run process
	return p.Process(ctx)
}

func (e *Executor) executeFetchers(ctx context.Context, fetchers []IFetcher) error {
	var (
		fs []func()
	)

	for _, fetcher := range fetchers {
		// run deps first
		e.executeFetchers(ctx, fetcher.DependFetchers())

		f := func(ft IFetcher) func() {
			return func() {
				once := e.fetcherOnces[ft.Name()]
				once.Do(func() {
					ft.Fetch(ctx)
				})
			}
		}(fetcher)

		fs = append(fs, f)
	}

	// real run
	ParallelDo(fs...)
	return nil
}
