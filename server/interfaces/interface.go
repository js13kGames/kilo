package interfaces

import "github.com/js13kgames/kilo/server"

type Tag uint8

const (
	API Tag = 1 << iota
	Auth
)

type Interface interface {
	server.Process
	GetKind() string
	HasTags(Tag) bool
}
