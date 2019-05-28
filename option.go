package chashu

type Option func(*config)

func HashFunc(f func(string) uint32) Option {
	return func(c *config) {
		c.hash = f
	}
}

func VirtualNode(n int) Option {
	return func(c *config) {
		c.numVNode = n
	}
}
