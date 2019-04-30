package nonce

const (
	// 1 in 16 777 216 collision chance assuming Nonces only get used to verify
	// a single transaction (like OAuth2 state, in our case) in a CSRF context.
	Size    = 4
	Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
)

type Nonce [Size]byte
