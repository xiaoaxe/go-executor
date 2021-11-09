//utils
//@author: baoqiang
//@time: 2021/11/09 21:18:33
package executor

import "reflect"

func structName(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
