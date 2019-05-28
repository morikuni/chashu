package chashu

type Option func(*config)

func WithHashFunc(f func(string) uint64) Option {
	return func(c *config) {
		c.hash = f
	}
}
