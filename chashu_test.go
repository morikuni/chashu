package chashu

import (
	"math/rand"
	"testing"
	"testing/quick"
	"time"
)

func TestResolver(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)

	err := quick.Check(func(a, b, c, d string, keys []string) bool {
		r1 := NewResolver(3, func(i int) string { return []string{a, b, c}[i] })
		r2 := NewResolver(2, func(i int) string { return []string{a, c}[i] })
		r3 := NewResolver(4, func(i int) string { return []string{a, b, c, d}[i] })

		for _, key := range keys {
			idx1 := r1.ResolveIndex(key)
			idx2 := r2.ResolveIndex(key)
			idx3 := r3.ResolveIndex(key)

			equal := func(i1, i2, i3 int) bool {
				if i1 != -1 && idx1 != i1 {
					return false
				}
				if i2 != -1 && idx2 != i2 {
					return false
				}
				if i3 != -1 && idx3 != i3 {
					return false
				}
				return true
			}

			switch {
			case equal(0, 0, 0):
			case equal(1, -1, 1):
			case equal(2, 1, 2):
			case equal(0, 0, 3):
			case equal(1, -1, 3):
			case equal(2, 1, 3):
			default:
				t.Log("ng", idx1, idx2, idx3)
				return false
			}
		}
		return true
	}, &quick.Config{
		Rand: rand.New(rand.NewSource(seed)),
	})
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
