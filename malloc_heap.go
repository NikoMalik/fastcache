//go:build appengine || windows
// +build appengine windows

package fastcache

func getChunk() []byte {
	return make([]byte, chunkSize)
}

func clearChunks() error {
	// No-op.
	return nil
}

func putChunk(chunk []byte) {
	// No-op.
}
