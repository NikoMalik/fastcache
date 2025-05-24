package fastcache

import (
	"math/bits"
	"sync/atomic"
)

const (
	_wyp0 = 0xa0761d6478bd642f
	_wyp1 = 0xe7037ed1a0b428db
	_wyp2 = 0x8ebc6af09c88c6e3
	_wyp3 = 0x589965cc75374cc3
	_wyp4 = 0x1d8e4e27c47d124f
)

// SRNG is a thread-safe random number generator.
// The zero value is valid.
type SRNG uint64

func (s SRNG) Uint() uint64 { return uint64(s) }

func _wymum(A, B uint64) uint64 {
	hi, lo := bits.Mul64(A, B)
	return hi ^ lo
}

// Int returns a random positive int.
func (r *SRNG) Int() int {
	return int(uint(r.Uint64()) >> 1)
}

// Intn returns an int uniformly in [0, n).
func (r *SRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.Uint64n(uint64(n)))
}

// Uint64 returns a random uint64.
func (r *SRNG) Uint64() uint64 {
	return _wymum(atomic.AddUint64((*uint64)(r), _wyp0)^_wyp1, uint64(*r))
}

// Uint64n returns a uint64 uniformly in [0, n).
func (r *SRNG) Uint64n(n uint64) uint64 {
	if n == 0 {
		return 0
	}

	x := r.Uint64()
	h, l := bits.Mul64(x, n)

	if l < n {
		t := -n
		if t >= n {
			t -= n
			if t >= n {
				t = t % n
			}
		}

	again:
		if l < t {
			x = r.Uint64()
			h, l = bits.Mul64(x, n)
			goto again
		}
	}

	return h
}

// Float64 returns a float64 uniformly in [0, 1).
func (r *SRNG) Float64() (v float64) {
again:
	v = float64(r.Uint64()>>11) / (1 << 53)
	if v == 1 {
		goto again
	}
	return
}
