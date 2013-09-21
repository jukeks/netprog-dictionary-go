package dictionary

import (
	"sync"
)

type Dictionary struct {
	l sync.Mutex
	m map[string]string
}

func New() Dictionary {
	d := Dictionary{}
	d.m = make(map[string]string)

	return d
}

func (d *Dictionary) Add(key, value string) bool {
	d.l.Lock()
	defer d.l.Unlock()

	_, ok := d.m[key]
	if !ok {
		d.m[key] = value
	}

	return !ok
}

func (d *Dictionary) Update(key, value string) bool {
	d.l.Lock()
	defer d.l.Unlock()

	_, ok := d.m[key]
	if ok {
		d.m[key] = value
	}

	return ok
}

func (d *Dictionary) Get(key string) (string, bool) {
	d.l.Lock()
	defer d.l.Unlock()

	value, ok := d.m[key]

	return value, ok
}

func (d *Dictionary) Remove(key string) bool {
	d.l.Lock()
	defer d.l.Unlock()

	_, ok := d.m[key]
	if ok {
		delete(d.m, key)
	}

	return ok
}
