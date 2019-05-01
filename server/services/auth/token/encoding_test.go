package token

import (
	"testing"
)

const (
	decoded = "1234567890abcdefghijklmn"
	encoded = "CJ8pD3KsDpWvC65YOsHbPcTeQMfhR6rk"
)

var tok Token

func TestEncode(t *testing.T) {
	var (
		tok Token
		dst = make([]byte, 32)
	)

	copy(tok[:], decoded)
	tok.encode(dst)

	if string(dst) != encoded {
		t.Errorf("Expected %s, got %s instead.", encoded, dst)
	}
}

func TestDecode(t *testing.T) {
	var (
		tok Token
	)

	tok.decode([]byte(encoded))

	if string(tok[:]) != decoded {
		t.Errorf("Expected %s, got %s instead.", decoded, tok)
	}
}

func BenchmarkEncode(b *testing.B) {
	dst := make([]byte, 32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tok.encode(dst)
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	src := []byte(encoded)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tok.decode(src)
		}
	})
}
