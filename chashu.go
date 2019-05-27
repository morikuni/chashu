package chashu

import (
	"sort"
)

type Resolver interface {
	ResolveIndex(key string) int
}

func NewResolver(nodes []string) Resolver {
	hash := fnvHash
	ring := make([]element, len(nodes))
	for i := range nodes {
		ring[i] = element{hash(nodes[i]), i}
	}
	sort.Slice(ring, func(i, j int) bool {
		return ring[i].hash < ring[j].hash
	})

	return resolver{ring, hash}
}

type element struct {
	hash  uint64
	index int
}

type resolver struct {
	ring []element
	hash func(key string) uint64
}

func (r resolver) ResolveIndex(key string) int {
	h := r.hash(key)
	i := sort.Search(len(r.ring), func(i int) bool {
		return r.ring[i].hash >= h
	})
	if i >= len(r.ring) {
		return 0
	}
	return r.ring[i].index
}
