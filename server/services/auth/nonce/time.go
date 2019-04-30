package nonce

import (
	_ "unsafe" // required to use //go:linkname
)

//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64)
