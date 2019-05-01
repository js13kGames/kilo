package token

import "github.com/js13kgames/kilo/server/services/auth/internal/encoding"

// Tokens are base64 encoded in transit using the standard URL-safe variant of base64.
// As our sizes are known, static and don't need padding, enc/dec are unrolled and have all
// sanity checks removed.
//
// BenchmarkEncodeStd   50000000               25.9 ns/op             0 B/op          0 allocs/op
// BenchmarkEncodeDev   100000000              16.6 ns/op             0 B/op          0 allocs/op
//
// BenchmarkDecodeStd    20000000              65.3 ns/op             0 B/op          0 allocs/op
// BenchmarkDecodeDev    100000000             17.5 ns/op             0 B/op          0 allocs/op
//
// Where std is encoding/base64 in the std lib and dev is our unrolled version, with decoding
// performance being the more important factor to consider in our case.

const enc = encoding.Base64UrlEncoding

var dec = encoding.Base64UrlDecoding

func (src *Token) encode(dst []byte) {
	// BCE hint.
	_ = dst[31]

	dst[0] = enc[(src[0]>>2)&0x3F]
	dst[1] = enc[(src[0]<<4|src[1]>>4)&0x3F]
	dst[2] = enc[(src[1]<<2|src[2]>>6)&0x3F]
	dst[3] = enc[src[2]&0x3F]
	dst[4] = enc[(src[3]>>2)&0x3F]
	dst[5] = enc[(src[3]<<4|src[4]>>4)&0x3F]
	dst[6] = enc[(src[4]<<2|src[5]>>6)&0x3F]
	dst[7] = enc[src[5]&0x3F]

	dst[8] = enc[(src[6]>>2)&0x3F]
	dst[9] = enc[(src[6]<<4|src[7]>>4)&0x3F]
	dst[10] = enc[(src[7]<<2|src[8]>>6)&0x3F]
	dst[11] = enc[src[8]&0x3F]
	dst[12] = enc[(src[9]>>2)&0x3F]
	dst[13] = enc[(src[9]<<4|src[10]>>4)&0x3F]
	dst[14] = enc[(src[10]<<2|src[11]>>6)&0x3F]
	dst[15] = enc[src[11]&0x3F]

	dst[16] = enc[(src[12]>>2)&0x3F]
	dst[17] = enc[(src[12]<<4|src[13]>>4)&0x3F]
	dst[18] = enc[(src[13]<<2|src[14]>>6)&0x3F]
	dst[19] = enc[src[14]&0x3F]
	dst[20] = enc[(src[15]>>2)&0x3F]
	dst[21] = enc[(src[15]<<4|src[16]>>4)&0x3F]
	dst[22] = enc[(src[16]<<2|src[17]>>6)&0x3F]
	dst[23] = enc[src[17]&0x3F]

	dst[24] = enc[(src[18]>>2)&0x3F]
	dst[25] = enc[(src[18]<<4|src[19]>>4)&0x3F]
	dst[26] = enc[(src[19]<<2|src[20]>>6)&0x3F]
	dst[27] = enc[src[20]&0x3F]
	dst[28] = enc[(src[21]>>2)&0x3F]
	dst[29] = enc[(src[21]<<4|src[22]>>4)&0x3F]
	dst[30] = enc[(src[22]<<2|src[23]>>6)&0x3F]
	dst[31] = enc[src[23]&0x3F]
}

func (dst *Token) decode(src []byte) {
	// BCE hint.
	_ = src[31]

	dst[0] = byte(dec[src[0]]<<2 | dec[src[1]]>>4)
	dst[1] = byte(dec[src[1]]<<4 | dec[src[2]]>>2)
	dst[2] = byte(dec[src[2]]<<6 | dec[src[3]])
	dst[3] = byte(dec[src[4]]<<2 | dec[src[5]]>>4)
	dst[4] = byte(dec[src[5]]<<4 | dec[src[6]]>>2)
	dst[5] = byte(dec[src[6]]<<6 | dec[src[7]])

	dst[6] = byte(dec[src[8]]<<2 | dec[src[9]]>>4)
	dst[7] = byte(dec[src[9]]<<4 | dec[src[10]]>>2)
	dst[8] = byte(dec[src[10]]<<6 | dec[src[11]])
	dst[9] = byte(dec[src[12]]<<2 | dec[src[13]]>>4)
	dst[10] = byte(dec[src[13]]<<4 | dec[src[14]]>>2)
	dst[11] = byte(dec[src[14]]<<6 | dec[src[15]])

	dst[12] = byte(dec[src[16]]<<2 | dec[src[17]]>>4)
	dst[13] = byte(dec[src[17]]<<4 | dec[src[18]]>>2)
	dst[14] = byte(dec[src[18]]<<6 | dec[src[19]])
	dst[15] = byte(dec[src[20]]<<2 | dec[src[21]]>>4)
	dst[16] = byte(dec[src[21]]<<4 | dec[src[22]]>>2)
	dst[17] = byte(dec[src[22]]<<6 | dec[src[23]])

	dst[18] = byte(dec[src[24]]<<2 | dec[src[25]]>>4)
	dst[19] = byte(dec[src[25]]<<4 | dec[src[26]]>>2)
	dst[20] = byte(dec[src[26]]<<6 | dec[src[27]])
	dst[21] = byte(dec[src[28]]<<2 | dec[src[29]]>>4)
	dst[22] = byte(dec[src[29]]<<4 | dec[src[30]]>>2)
	dst[23] = byte(dec[src[30]]<<6 | dec[src[31]])
}
