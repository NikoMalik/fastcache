package fastcache

import (
	"bytes"
	"strconv"
	"testing"

	wyhash "github.com/orisano/wyhash/v4"
)

func TestGenerationOverflow(t *testing.T) {
	c := New(1)

	key1, key2 := findKeysForBucket(0)
	if key1 == nil || key2 == nil {
		t.Fatalf("Failed to find keys for the same bucket")
	}

	h1 := wyhash.Sum64(0, key1)
	bucketIdx := h1 % bucketsCount

	bigVal1 := make([]byte, 1024)
	for i := range bigVal1 {
		bigVal1[i] = 1
	}
	bigVal2 := make([]byte, 1024)
	for i := range bigVal2 {
		bigVal2[i] = 2
	}

	genVal := func(t *testing.T, c *Cache, bucketIdx uint64, expected uint64) {
		t.Helper()
		actual := c.buckets[bucketIdx].gen
		if actual != expected {
			t.Fatalf("Expected generation to be %d found %d instead for bucket %d", expected, actual, bucketIdx)
		}
	}

	c.Set(key1, bigVal1)
	getVal(t, c, key1, bigVal1)
	c.Set(key2, bigVal2)
	getVal(t, c, key2, bigVal2)
	genVal(t, c, bucketIdx, 1) // 1

	c.buckets[bucketIdx].gen = (1 << 24) - 2

	c.Set(key1, bigVal1)
	getVal(t, c, key1, bigVal1)
	c.Set(key2, bigVal2)
	getVal(t, c, key2, bigVal2)
	genVal(t, c, bucketIdx, (1<<24)-2)

	c.buckets[bucketIdx].gen = maxGen // 16777215 for genSizeBits=24
	c.Set(key1, bigVal1)
	getVal(t, c, key1, bigVal1)
	c.Set(key2, bigVal2)
	getVal(t, c, key2, bigVal2)
	genVal(t, c, bucketIdx, maxGen)
}
func getVal(t *testing.T, c *Cache, key, expected []byte) {
	t.Helper()
	get := c.Get(nil, key)
	if !bytes.Equal(get, expected) {
		t.Errorf("Expected value (%v) was not returned from the cache, instead got %v", expected[:10], get)
	}
}

func genVal(t *testing.T, c *Cache, bucketIdx uint64, expected uint64) {
	t.Helper()
	actual := c.buckets[bucketIdx].gen
	if actual != expected {
		t.Fatalf("Expected generation to be %d found %d instead for bucket %d", expected, actual, bucketIdx)
	}
}
func findKeysForBucket(seed uint64) (key1, key2 []byte) {
	for i := 0; i < 1000; i++ {
		k1 := []byte(strconv.Itoa(i))
		for j := i + 1; j < 1000; j++ {
			k2 := []byte(strconv.Itoa(j))
			if wyhash.Sum64(seed, k1)%bucketsCount == wyhash.Sum64(seed, k2)%bucketsCount {
				return k1, k2
			}
		}
	}
	return nil, nil
}
