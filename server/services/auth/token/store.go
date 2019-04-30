package token

import (
	"crypto/rand"
	"sync"

	bolt "go.etcd.io/bbolt"

	"github.com/js13kgames/kilo/server/types"
)

type Store struct {
	db             *bolt.DB
	dbBucketKey    []byte
	items          map[Token]*TokenMetadata
	itemsByActorId map[types.ActorId][]Token
	mu             sync.RWMutex
}

func NewStore(db *bolt.DB, dbBucketKey []byte) *Store {
	if dbBucketKey == nil {
		dbBucketKey = []byte("auth.tokens")
	}

	store := &Store{
		db:             db,
		dbBucketKey:    dbBucketKey,
		items:          make(map[Token]*TokenMetadata),
		itemsByActorId: make(map[types.ActorId][]Token),
	}

	if err := loadTokens(store, db); err != nil {
		panic(err)
	}

	return store
}

// loadTokens hydrates the store with data persisted in db, building relevant indices
// in the process.
func loadTokens(store *Store, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(store.dbBucketKey)
		if err != nil {
			return err
		}

		var (
			token   Token
			actorId types.ActorId
			c       = bucket.Cursor()
		)

		for k, v := c.First(); k != nil; k, v = c.Next() {
			copy(token[:], k)
			copy(actorId[:], v)

			store.items[token] = &TokenMetadata{
				ActorId: actorId,
			}
			store.itemsByActorId[actorId] = append(store.itemsByActorId[actorId], token)
		}

		return nil
	})
}

//
func (store *Store) Generate(actorId types.ActorId) (Token, error) {
	var token Token

	_, err := rand.Read(token[:])
	if err != nil {
		return token, err
	}

	store.mu.Lock()

	// In the off-chance of a collision, repeat.
	if _, exists := store.items[token]; exists {
		store.mu.Unlock()
		return store.Generate(actorId)
	}

	if err = store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(store.dbBucketKey).Put(token[:], actorId[:])
	}); err != nil {
		store.mu.Unlock()
		return zero, err
	}

	metadata := &TokenMetadata{
		ActorId: actorId,
	}

	store.items[token] = metadata
	store.itemsByActorId[actorId] = append(store.itemsByActorId[actorId], token)
	store.mu.Unlock()

	return token, err
}

// Get returns the metadata associated with the given Token and true if the Token
// exists in the Store. Otherwise it returns nil and false respectively.
func (store *Store) Get(token Token) (*TokenMetadata, bool) {
	store.mu.RLock()
	t, ok := store.items[token]
	store.mu.RUnlock()

	return t, ok
}

// List returns a list of all Tokens owned by the given actor along with their metadata.
func (store *Store) List(actorId types.ActorId) []tokenListItem {
	store.mu.RLock()

	var (
		tokens = store.itemsByActorId[actorId]
		out    = make([]tokenListItem, len(tokens))
	)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		out[i] = tokenListItem{
			Token:         token,
			TokenMetadata: store.items[token],
		}
	}

	store.mu.RUnlock()

	return out
}

type tokenListItem struct {
	*TokenMetadata
	ActorId types.ActorId `json:"-"`
	Token   Token         `json:"token"`
}

//
func (store *Store) RevokeOwned(token Token, actorId types.ActorId) (bool, error) {
	store.mu.Lock()

	var (
		tokens = store.itemsByActorId[actorId]
		j      = -1
		n      = len(tokens)
		err    error
	)

	if n == 1 && tokens[0] == token {
		j = 0
	} else {
		for i := 0; i < len(tokens); i++ {
			if tokens[i] == token {
				j = i
				break
			}
		}
	}

	if j == -1 {
		store.mu.Unlock()
		return false, nil
	}

	if err = store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(store.dbBucketKey).Delete(token[:])
	}); err != nil {
		store.mu.Unlock()
		return true, err
	}

	if j == 0 && n == 1 {
		// Most common use case is a single token per user so we can release
		// the backing array entirely if the token gets revoked.
		store.itemsByActorId[actorId] = nil
	} else {
		copy(tokens[j:], tokens[j+1:])
		tokens[n-1] = zero

		store.itemsByActorId[actorId] = tokens[:n-1]
	}

	delete(store.items, token)

	store.mu.Unlock()

	return true, nil
}
