package chashu

import (
	"sort"
)

type Resolver interface {
	ResolveIndex(key string) int
	ReHash(n int, f func(i int) string)
}

func NewResolver(n int, f func(i int) string) Resolver {
	r := resolver{hash: fnvHash}
	r.ReHash(n, f)
	return &r
}

type element struct {
	hash  uint64
	index int
}

type resolver struct {
	ring []element
	hash func(key string) uint64
}

func (r *resolver) ReHash(n int, f func(i int) string) {
	ring := make([]element, n)
	for i := 0; i < n; i++ {
		ring[i] = element{r.hash(f(i)), i}
	}
	sort.Slice(ring, func(i, j int) bool {
		return ring[i].hash < ring[j].hash
	})
	r.ring = ring
}

func (r *resolver) ResolveIndex(key string) int {
	h := r.hash(key)
	i := sort.Search(len(r.ring), func(i int) bool {
		return r.ring[i].hash >= h
	})
	if i >= len(r.ring) {
		return 0
	}
	return r.ring[i].index
}
