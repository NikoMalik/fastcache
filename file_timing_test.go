package fastcache

import (
	"fmt"
	"os"
	"sync"
	"testing"

	xxhash "github.com/cespare/xxhash/v2"
	wyhash "github.com/orisano/wyhash/v4"
	wyhashv1 "github.com/zeebo/wyhash"
)

func BenchmarkSaveToFile(b *testing.B) {
	for _, concurrency := range []int{1, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf("concurrency_%d", concurrency), func(b *testing.B) {
			benchmarkSaveToFile(b, concurrency)
		})
	}
}

func BenchmarkHash(b *testing.B) {
	key := []byte("testkey123")
	seed := uint64(12345)
	b.Run("wyhash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wyhash.Sum64(seed, key)
		}
	})

	b.Run("wyhash1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wyhashv1.Hash(key, seed)
		}
	})

	b.Run("xxhash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xxhash.Sum64(key)
		}
	})
}

func benchmarkSaveToFile(b *testing.B, concurrency int) {
	filePath := fmt.Sprintf("BencharkSaveToFile.%d.fastcache", concurrency)
	defer os.RemoveAll(filePath)
	c := newBenchCache()

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(benchCacheSize)
	for i := 0; i < b.N; i++ {
		if err := c.SaveToFileConcurrent(filePath, concurrency); err != nil {
			b.Fatalf("unexpected error when saving to file: %s", err)
		}
	}
}

func BenchmarkLoadFromFile(b *testing.B) {
	for _, concurrency := range []int{1, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf("concurrency_%d", concurrency), func(b *testing.B) {
			benchmarkLoadFromFile(b, concurrency)
		})
	}
}

func benchmarkLoadFromFile(b *testing.B, concurrency int) {
	filePath := fmt.Sprintf("BenchmarkLoadFromFile.%d.fastcache", concurrency)
	defer os.RemoveAll(filePath)

	c := newBenchCache()
	if err := c.SaveToFileConcurrent(filePath, concurrency); err != nil {
		b.Fatalf("cannot save cache to file: %s", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(benchCacheSize)
	for i := 0; i < b.N; i++ {
		c, err := LoadFromFile(filePath, 0)
		if err != nil {
			b.Fatalf("cannot load cache from file: %s", err)
		}
		var s Stats
		c.UpdateStats(&s)
		if s.EntriesCount == 0 {
			b.Fatalf("unexpected zero entries")
		}
	}
}

var (
	benchCache     *Cache
	benchCacheOnce sync.Once
)

func newBenchCache() *Cache {
	benchCacheOnce.Do(func() {
		c := New(benchCacheSize)
		itemsCount := benchCacheSize / 20
		for i := 0; i < itemsCount; i++ {
			k := []byte(fmt.Sprintf("key %d", i))
			v := []byte(fmt.Sprintf("value %d", i))
			c.Set(k, v)
		}
		benchCache = c
	})
	return benchCache
}

const benchCacheSize = bucketsCount * chunkSize
