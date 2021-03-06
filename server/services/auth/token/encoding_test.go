package token

import (
	"testing"
)

const (
	decoded = "1234567890abcdefghijklmn"
	encoded = "MTIzNDU2Nzg5MGFiY2RlZmdoaWprbG1u"
)

func TestEncodeScalar(t *testing.T) {
	var (
		tok Token
		dst = make([]byte, 32)
	)

	copy(tok[:], decoded)
	encodeScalar(dst, &tok)

	if string(dst) != encoded {
		t.Errorf("Expected %s, got %s instead.", encoded, dst)
	}
}

func TestDecodeScalar(t *testing.T) {
	var (
		tok Token
	)

	decodeScalar(&tok, []byte(encoded))

	if string(tok[:]) != decoded {
		t.Errorf("Expected %s, got %s instead.", decoded, tok)
	}
}

func BenchmarkEncodeScalar(b *testing.B) {
	var (
		dst = make([]byte, SizeEncoded)
		src = &Token{}
	)

	copy(src[:], decoded)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encodeScalar(dst, src)
		}
	})
}

func BenchmarkDecodeScalar(b *testing.B) {
	var (
		dst = &Token{}
		src = []byte(encoded)
	)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			decodeScalar(dst, src)
		}
	})
}
