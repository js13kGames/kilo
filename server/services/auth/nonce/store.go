package nonce

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/js13kgames/kilo/server/types"
)

const (
	ttl             = 300
	cleanupInterval = 600
)

// Store holds Nonces and manages their lifecycle.
// Nonces are short-lived, single-use tokens which expire after a default TTL measured by a
// monotonic clock. They are kept in-memory only and not persisted, making them volatile
// to process shutdowns.
type Store struct {
	items map[Nonce]data
	mu    sync.Mutex
	done  chan struct{}
}

func New() *Store {
	return &Store{
		items: make(map[Nonce]data),
	}
}

type data struct {
	exp     uint32
	actorId *types.ActorId
}

// Generate creates and stores a Nonce, optionally storing an actorId alongside it. The actorId
// will get returned once the Nonce gets consumed.
//
// TODO(alcore) Could be decoupled from our own models and take arbitrary metadata instead.
func (store *Store) Generate(actorId *types.ActorId) (*Nonce, error) {
	var (
		n     Nonce
		d     data
		mnsec int64
	)

	if _, err := rand.Read(n[:]); err != nil {
		return nil, err
	}

	for i := 0; i < Size; i++ {
		n[i] = Charset[n[i]>>2]
	}

	// mnsec is monotonic nanoseconds since process start.
	_, _, mnsec = now()
	d = data{
		exp:     uint32(mnsec/1e9) + ttl,
		actorId: actorId,
	}

	store.mu.Lock()
	store.items[n] = d
	store.mu.Unlock()

	return &n, nil
}

// Consume takes a Nonce and returns the data associated with it if the Nonce
// is known to the Store and has not yet expired in the process (ie. is valid).
// The Nonce gets invalidated in the process.
func (store *Store) Consume(n Nonce) (*types.ActorId, bool) {
	_, _, mnsec := now()
	store.mu.Lock()

	data, ok := store.items[n]
	if !ok {
		store.mu.Unlock()
		return nil, false
	}

	// This happens every time, even if t has expired, in which case we simply offload
	// work from the cleanup routine by removing the item right away.
	delete(store.items, n)
	store.mu.Unlock()

	if uint32(mnsec/1e9) > data.exp {
		return nil, false
	}

	return data.actorId, true
}

func (store *Store) Start() {
	if store.done != nil {
		panic("already running")
	}

	store.done = make(chan struct{})
	ticker := time.NewTicker(cleanupInterval)

	for {
		select {
		case <-ticker.C:
			_, _, mnsec := now()
			mono := uint32(mnsec) / 1e9
			store.mu.Lock()

			for k, md := range store.items {
				if mono > md.exp {
					delete(store.items, k)
				}
			}

			store.mu.Unlock()

		case <-store.done:
			ticker.Stop()
			return
		}
	}
}

func (store *Store) Stop(deadline *time.Time) {
	if store.done == nil {
		return
	}

	close(store.done)
	store.done = nil
}
