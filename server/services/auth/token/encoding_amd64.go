package token

import "github.com/klauspost/cpuid"

// Proof of concept base64 encoding/decoding (using the URL-safe alphabet) with SSSE3 SIMD instructions
// as per Wojciech Muła based on his C implementation described at:
// - http://0x80.pl/notesen/2016-01-17-sse-base64-decoding.html (BSD 2-clause)
//
// Note that this implementation does not currently check validity nor does it handle padding -
// and is just an experiment tailored to this particular codebase.
// An AVX2 implementation could handle our payloads - ie. 32 bytes per call (Token) in a single pass.
// This one currently makes 2 passes but is still several times faster than a scalar approach.

// go:noescape
func encodeSSE3(dst []byte, src *Token)

// go:noescape
func decodeSSE3(dst *Token, src []byte)

func init() {
	if cpuid.CPU.SSSE3() {
		encode = encodeSSE3
		decode = decodeSSE3
	}
}
