package fastcache

import "unsafe"

func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
