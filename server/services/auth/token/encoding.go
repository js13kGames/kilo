package token

import "github.com/js13kgames/kilo/server/services/auth/internal/encoding"

// Tokens are base64 encoded in transit using the standard URL-safe variant of base64.
// As our sizes are known, static and don't need padding, enc/dec are unrolled and have all
// sanity checks removed.
//
// BenchmarkEncodeStd-8    50000000               25.9 ns/op             0 B/op          0 allocs/op
// BenchmarkEncodeDev-8    100000000              18.4 ns/op             0 B/op          0 allocs/op
//
// BenchmarkDecodeStd-8    20000000               65.3 ns/op             0 B/op          0 allocs/op
// BenchmarkDecodeDev-8    100000000              22.9 ns/op             0 B/op          0 allocs/op
//
// Where std is encoding/base64 in the std lib and dev is our unrolled version, with decoding
// performance being the more important factor to consider. Results apply on 64-bit platforms only
// due to the use of uint64s to hold intermediate values.

var (
	enc = encoding.Base64UrlEncoding
	dec = encoding.Base64UrlDecoding
)

func (src *Token) encode(dst []byte) {
	// BCE hints.
	_ = src[23]
	_ = dst[31]

	var v uint64

	// 6 octets into 8 characters, 4 times in total. The intermediary additions into uint64s can be skipped
	// and the respective shifts applied directly but the bit fiddling code makes this even more convoluted.
	v = uint64(src[0])<<40 | uint64(src[1])<<32 | uint64(src[2])<<24 | uint64(src[3])<<16 | uint64(src[4])<<8 | uint64(src[5])

	dst[0] = enc[v>>42&0x3F]
	dst[1] = enc[v>>36&0x3F]
	dst[2] = enc[v>>30&0x3F]
	dst[3] = enc[v>>24&0x3F]
	dst[4] = enc[v>>18&0x3F]
	dst[5] = enc[v>>12&0x3F]
	dst[6] = enc[v>>6&0x3F]
	dst[7] = enc[v&0x3F]

	v = uint64(src[6])<<40 | uint64(src[7])<<32 | uint64(src[8])<<24 | uint64(src[9])<<16 | uint64(src[10])<<8 | uint64(src[11])

	dst[8] = enc[v>>42&0x3F]
	dst[9] = enc[v>>36&0x3F]
	dst[10] = enc[v>>30&0x3F]
	dst[11] = enc[v>>24&0x3F]
	dst[12] = enc[v>>18&0x3F]
	dst[13] = enc[v>>12&0x3F]
	dst[14] = enc[v>>6&0x3F]
	dst[15] = enc[v&0x3F]

	v = uint64(src[12])<<40 | uint64(src[13])<<32 | uint64(src[14])<<24 | uint64(src[15])<<16 | uint64(src[16])<<8 | uint64(src[17])

	dst[16] = enc[v>>42&0x3F]
	dst[17] = enc[v>>36&0x3F]
	dst[18] = enc[v>>30&0x3F]
	dst[19] = enc[v>>24&0x3F]
	dst[20] = enc[v>>18&0x3F]
	dst[21] = enc[v>>12&0x3F]
	dst[22] = enc[v>>6&0x3F]
	dst[23] = enc[v&0x3F]

	v = uint64(src[18])<<40 | uint64(src[19])<<32 | uint64(src[20])<<24 | uint64(src[21])<<16 | uint64(src[22])<<8 | uint64(src[23])

	dst[24] = enc[v>>42&0x3F]
	dst[25] = enc[v>>36&0x3F]
	dst[26] = enc[v>>30&0x3F]
	dst[27] = enc[v>>24&0x3F]
	dst[28] = enc[v>>18&0x3F]
	dst[29] = enc[v>>12&0x3F]
	dst[30] = enc[v>>6&0x3F]
	dst[31] = enc[v&0x3F]
}

func (dst *Token) decode(src []byte) {
	// BCE hints.
	_ = dst[23]
	_ = src[31]

	var v uint64

	v = uint64(dec[src[0]])<<58 | uint64(dec[src[1]])<<52 | uint64(dec[src[2]])<<46 | uint64(dec[src[3]])<<40 |
		uint64(dec[src[4]])<<34 | uint64(dec[src[5]])<<28 | uint64(dec[src[6]])<<22 | uint64(dec[src[7]])<<16

	dst[0] = byte(v >> 56)
	dst[1] = byte(v >> 48)
	dst[2] = byte(v >> 40)
	dst[3] = byte(v >> 32)
	dst[4] = byte(v >> 24)
	dst[5] = byte(v >> 16)

	v = uint64(dec[src[8]])<<58 | uint64(dec[src[9]])<<52 | uint64(dec[src[10]])<<46 | uint64(dec[src[11]])<<40 |
		uint64(dec[src[12]])<<34 | uint64(dec[src[13]])<<28 | uint64(dec[src[14]])<<22 | uint64(dec[src[15]])<<16

	dst[6] = byte(v >> 56)
	dst[7] = byte(v >> 48)
	dst[8] = byte(v >> 40)
	dst[9] = byte(v >> 32)
	dst[10] = byte(v >> 24)
	dst[11] = byte(v >> 16)

	v = uint64(dec[src[16]])<<58 | uint64(dec[src[17]])<<52 | uint64(dec[src[18]])<<46 | uint64(dec[src[19]])<<40 |
		uint64(dec[src[20]])<<34 | uint64(dec[src[21]])<<28 | uint64(dec[src[22]])<<22 | uint64(dec[src[23]])<<16

	dst[12] = byte(v >> 56)
	dst[13] = byte(v >> 48)
	dst[14] = byte(v >> 40)
	dst[15] = byte(v >> 32)
	dst[16] = byte(v >> 24)
	dst[17] = byte(v >> 16)

	v = uint64(dec[src[24]])<<58 | uint64(dec[src[25]])<<52 | uint64(dec[src[26]])<<46 | uint64(dec[src[27]])<<40 |
		uint64(dec[src[28]])<<34 | uint64(dec[src[29]])<<28 | uint64(dec[src[30]])<<22 | uint64(dec[src[31]])<<16

	dst[18] = byte(v >> 56)
	dst[19] = byte(v >> 48)
	dst[20] = byte(v >> 40)
	dst[21] = byte(v >> 32)
	dst[22] = byte(v >> 24)
	dst[23] = byte(v >> 16)
}
