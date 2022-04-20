package lru

import "container/list"

type cache struct {
	maxBytes        int64
	curBytes        int64
	linkList        *list.List
	data            map[string]*list.Element
	onRemoveHandler func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onRemoved func(string, Value)) *cache {
	return &cache{
		maxBytes:        maxBytes,
		onRemoveHandler: onRemoved,
		linkList:        list.New(),
		data:            make(map[string]*list.Element),
	}
}

func (c *cache) Set(key string, value Value) {
	if v, ok := c.data[key]; ok {
		oldLen := v.Value.(*entry).value.Len()
		c.linkList.MoveToFront(v)
		v.Value = value
		c.curBytes += int64(value.Len()) - int64(oldLen)
	} else {
		v := c.linkList.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.data[key] = v
		c.curBytes += int64(value.Len()) + int64(len(key))
	}

	for c.maxBytes != 0 && c.curBytes > c.maxBytes {
		c.RemoveFurthest()
	}
}

func (c *cache) RemoveFurthest() {
	v := c.linkList.Back()
	if v != nil {
		c.linkList.Remove(v)
		e := v.Value.(*entry)
		delete(c.data, e.key)
		c.curBytes -= int64(e.value.Len()) + int64(len(e.key))
		if c.onRemoveHandler != nil {
			c.onRemoveHandler(e.key, e.value)
		}
	}
}

func (c *cache) Get(key string) (value Value, ok bool) {
	if v, ok := c.data[key]; ok {
		c.linkList.MoveToFront(v)
		e := v.Value.(*entry)
		return e.value, true
	}
	return
}

func (c *cache) Len() int {
	return c.linkList.Len()
}
