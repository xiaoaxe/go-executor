//fetcher interface
//@author: baoqiang
//@time: 2021/11/09 20:48:18
package executor

import "context"

type IFetcher interface {
	DependFetchers() []IFetcher
	Fetch(ctx context.Context) error
	GetData() interface{}
	Name() string
}
