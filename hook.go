package dfa

import "sync"

var Registry = NewRegistry()

type registry struct {
	lock  sync.Mutex
	hooks map[string]func(*Status)
}

func (r *registry) Add(name string, hook func(*Status)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.hooks[name] = hook
}

func (r *registry) Get(name string) (func(*Status), bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	hook, ok := r.hooks[name]
	return hook, ok
}

func NewRegistry() *registry {
	return &registry{
		hooks: make(map[string]func(*Status)),
	}
}
