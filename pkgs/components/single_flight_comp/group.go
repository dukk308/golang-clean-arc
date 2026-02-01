package single_flight_comp

import (
	"golang.org/x/sync/singleflight"
)

type Group struct {
	g singleflight.Group
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	return g.g.Do(key, fn)
}

func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan singleflight.Result {
	return g.g.DoChan(key, fn)
}

func (g *Group) Forget(key string) {
	g.g.Forget(key)
}
