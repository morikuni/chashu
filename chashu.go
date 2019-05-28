package chashu

import (
	"sort"
	"strconv"
	"strings"
)

type Resolver interface {
	ResolveIndex(key string) int
	ReHash(n int, f func(i int) string)
}

func NewResolver(n int, f func(i int) string, opts ...Option) Resolver {
	c := defaultConfig
	for _, o := range opts {
		o(&c)
	}
	r := resolver{hash: c.hash, numVNode: c.numVNode}
	r.ReHash(n, f)
	return &r
}

type element struct {
	hash  uint32
	index int
}

type config struct {
	hash     func(string) uint32
	numVNode int
}

var defaultConfig = config{
	hash:     defaultHash,
	numVNode: 100,
}

type resolver struct {
	ring     []element
	hash     func(key string) uint32
	numVNode int
}

func (r *resolver) ReHash(n int, f func(i int) string) {
	type node struct {
		key   string
		index int
	}

	nodes := make([]node, n)
	for i := 0; i < n; i++ {
		nodes[i] = node{f(i), i}
	}
	// sort to ensure same result even if hash conflicts
	sort.Slice(nodes, func(i, j int) bool {
		return strings.Compare(nodes[i].key, nodes[j].key) < 0
	})

	ring := make([]element, 0, n*r.numVNode)
	for _, node := range nodes {
		for i := 0; i < r.numVNode; i++ {
			ring = append(ring, element{r.hash(node.key + strconv.Itoa(i)), node.index})
		}
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
		return r.ring[0].index
	}
	return r.ring[i].index
}
