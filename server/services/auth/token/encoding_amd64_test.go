package token

import (
	"testing"

	"github.com/klauspost/cpuid"
)

var hasSSE3 = cpuid.CPU.SSSE3()

func TestEncodeSSE3(t *testing.T) {
	if !hasSSE3 {
		return
	}

	var (
		tok Token
		dst = make([]byte, 32)
	)

	copy(tok[:], decoded)
	encodeSSE3(dst, &tok)

	if string(dst) != encoded {
		t.Errorf("Expected %s, got %s instead.", encoded, dst)
	}
}

func TestDecodeSSE3(t *testing.T) {
	if !hasSSE3 {
		return
	}

	var (
		tok Token
	)

	decodeSSE3(&tok, []byte(encoded))

	if string(tok[:]) != decoded {
		t.Errorf("Expected %s, got %s instead.", decoded, tok)
	}
}

func BenchmarkEncodeSSE3(b *testing.B) {
	if !hasSSE3 {
		return
	}

	var (
		dst = make([]byte, SizeEncoded)
		src = &Token{}
	)

	copy(src[:], decoded)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encodeSSE3(dst, src)
		}
	})
}

func BenchmarkDecodeSSE3(b *testing.B) {
	if !hasSSE3 {
		return
	}

	var (
		dst = &Token{}
		src = []byte(encoded)
	)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			decodeSSE3(dst, src)
		}
	})
}
