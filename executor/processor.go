//processor interface
//@author: baoqiang
//@time: 2021/11/09 20:48:46
package executor

import "context"

type Processor interface {
	Process(ctx context.Context) error
	GetFetchers() []IFetcher
	Name() string
}
