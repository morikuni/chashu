package chashu_test

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"testing/quick"
	"time"

	"github.com/morikuni/chashu"
)

func TestResolver(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %v", seed)
	r := rand.New(rand.NewSource(seed))

	err := quick.Check(func(ia, ib, ic, id int) bool {
		a, b, c, d := strconv.Itoa(ia), strconv.Itoa(ib), strconv.Itoa(ic), strconv.Itoa(id)
		r1 := chashu.NewResolver(3, func(i int) string { return []string{a, b, c}[i] })
		r2 := chashu.NewResolver(2, func(i int) string { return []string{c, a}[i] })
		r3 := chashu.NewResolver(4, func(i int) string { return []string{a, c, d, b}[i] })

		err := quick.Check(func(ikey int) bool {
			key := strconv.Itoa(ikey)
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
			case equal(0, 1, 0):
			case equal(1, -1, 3):
			case equal(2, 0, 1):
			case equal(0, 1, 2):
			case equal(1, -1, 2):
			case equal(2, 0, 2):
			default:
				t.Log("ng", idx1, idx2, idx3)
				return false
			}
			return true
		}, &quick.Config{Rand: r})
		if err != nil {
			t.Errorf("error: %v", err)
			return false
		}
		return true
	}, &quick.Config{Rand: r})
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func ExampleNewResolver() {
	ips := []net.IP{{192, 168, 10, 2}, {192, 168, 10, 3}, {192, 168, 10, 4}}
	r := chashu.NewResolver(len(ips), func(i int) string {
		return ips[i].String()
	})
	fmt.Println(ips[r.ResolveIndex("1")].String())
	fmt.Println(ips[r.ResolveIndex("2")].String())
	fmt.Println(ips[r.ResolveIndex("3")].String())

	ips = append(ips[:1], ips[2]) // remove 2nd IP (192.168.10.3)
	r.ReHash(len(ips), func(i int) string {
		return ips[i].String()
	})
	fmt.Println("=== re-hash ===")
	fmt.Println(ips[r.ResolveIndex("1")].String()) // same node
	fmt.Println(ips[r.ResolveIndex("2")].String()) // same node
	fmt.Println(ips[r.ResolveIndex("3")].String()) // moved to other nodes, because 192.168.10.3 was removed
	// Output:
	// 192.168.10.2
	// 192.168.10.4
	// 192.168.10.3
	// === re-hash ===
	// 192.168.10.2
	// 192.168.10.4
	// 192.168.10.4
}
