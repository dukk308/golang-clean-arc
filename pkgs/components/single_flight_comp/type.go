package single_flight_comp

import (
	"golang.org/x/sync/singleflight"
)

type Result struct {
	Value  interface{}
	Err    error
	Shared bool
}

type Singleflight interface {
	Do(key string, fn func() (interface{}, error)) (interface{}, error, bool)
	DoChan(key string, fn func() (interface{}, error)) <-chan singleflight.Result
	Forget(key string)
}
