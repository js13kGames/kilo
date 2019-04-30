package nonce

import "unsafe"

const Size = 4

type Nonce [Size]byte

func (n Nonce) String() string {
	b := n[:]
	return *(*string)(unsafe.Pointer(&b))
}

func (n Nonce) MarshalText() ([]byte, error) {
	return n[:], nil
}

func (n *Nonce) UnmarshalText(src []byte) error {
	if len(src) != Size {
		return errInvalidInput
	}

	n[0] = src[0]
	n[1] = src[1]
	n[2] = src[2]
	n[3] = src[3]

	return nil
}
