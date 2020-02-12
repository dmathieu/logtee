package metrics

import (
	"sync"
	"sync/atomic"
)

var globalRegistry = defaultRegistry()

// RegistryProvider returns the globally shared registry.
func RegistryProvider() *Registry {
	return globalRegistry.Load().(*Registry)
}

// Registry allows storing and retrieving all metrics
type Registry struct {
	counters map[string]*Counter
	mux      sync.Mutex
}

// Counter creates and retrieves the Counter with the specified key
func (r *Registry) Counter(key string) *Counter {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.counters[key]; !ok {
		r.counters[key] = newCounter(key)
	}

	return r.counters[key]
}

func newRegistry() *Registry {
	return &Registry{counters: map[string]*Counter{}}
}

func defaultRegistry() *atomic.Value {
	var v atomic.Value
	v.Store(newRegistry())
	return &v
}
