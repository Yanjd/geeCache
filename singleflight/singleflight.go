package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}
type Group struct {
	mu    sync.Mutex
	calls map[string]*call
}

func (g *Group) Do(key string, handle func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.calls == nil {
		g.calls = make(map[string]*call)
	}
	if v, ok := g.calls[key]; ok {
		g.mu.Unlock()
		v.wg.Wait()
		return v.val, v.err
	}
	c := new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.mu.Unlock()
	c.val, c.err = handle()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()

	return c.val, c.err
}
