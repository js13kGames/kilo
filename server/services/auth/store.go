package auth

import (
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/js13kgames/kilo/server/services/auth/nonce"
	"github.com/js13kgames/kilo/server/services/auth/token"
)

type Store struct {
	tokens *token.Store
	nonces *nonce.Store
}

func NewStore(db *bolt.DB) *Store {
	return &Store{
		tokens: token.NewStore(db, nil),
		nonces: nonce.NewStore(),
	}
}

func (store *Store) Start() {
	store.nonces.Start()
}

func (store *Store) Stop(deadline *time.Time) {
	store.nonces.Stop(deadline)
}
