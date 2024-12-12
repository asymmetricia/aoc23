package aoc

type cache[K comparable, V any] struct {
	Fn    func(K) V
	cache map[K]V
}

func Cache[K comparable, V any](fn func(K) V) func(K) V {
	cache := &cache[K, V]{Fn: fn, cache: map[K]V{}}
	return cache.do
}

func Cache2[K comparable, K2 comparable, V any](fn func(K, K2) V) func(K, K2) V {
	type kt struct {
		k  K
		k2 K2
	}

	cache := &cache[kt, V]{Fn: func(kt kt) V {
		return fn(kt.k, kt.k2)
	}, cache: map[kt]V{}}

	return func(k K, k2 K2) V {
		return cache.do(kt{k, k2})
	}
}

func (c *cache[K, V]) do(k K) V {
	if c.cache == nil {
		c.cache = map[K]V{}
	}
	if _, ok := c.cache[k]; !ok {
		c.cache[k] = c.Fn(k)
	}
	return c.cache[k]
}
