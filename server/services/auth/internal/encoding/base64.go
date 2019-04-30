package encoding

const Base64UrlEncoding = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

var (
	// Decode LUT.
	Base64UrlDecoding [256]byte
)

func init() {
	for i := 0; i < len(Base64UrlDecoding); i++ {
		Base64UrlDecoding[i] = 0xFF
	}

	for i := 0; i < len(Base64UrlEncoding); i++ {
		Base64UrlDecoding[Base64UrlEncoding[i]] = byte(i)
	}
}
