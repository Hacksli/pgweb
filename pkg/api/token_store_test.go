package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenStore_CreateAndGet(t *testing.T) {
	store := NewTokenStore(time.Hour)

	token, err := store.Create("postgres://localhost/db")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	url, err := store.Get(token)
	assert.NoError(t, err)
	assert.Equal(t, "postgres://localhost/db", url)
}

func TestTokenStore_InvalidToken(t *testing.T) {
	store := NewTokenStore(time.Hour)

	_, err := store.Get("does-not-exist")
	assert.Equal(t, errTokenInvalid, err)
}

func TestTokenStore_SlidingExpiration(t *testing.T) {
	store := NewTokenStore(60 * time.Millisecond)

	token, err := store.Create("postgres://localhost/db")
	assert.NoError(t, err)

	// Keep accessing before expiry; each access extends the lifetime.
	for i := 0; i < 5; i++ {
		time.Sleep(40 * time.Millisecond)
		_, err := store.Get(token)
		assert.NoError(t, err, "token should still be valid on access %d", i)
	}
}

func TestTokenStore_ExpiresOnAccess(t *testing.T) {
	store := NewTokenStore(30 * time.Millisecond)

	token, err := store.Create("postgres://localhost/db")
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	_, err = store.Get(token)
	assert.Equal(t, errTokenInvalid, err)
}

func TestTokenStore_CleanupRemovesUnaccessed(t *testing.T) {
	store := NewTokenStore(30 * time.Millisecond)

	token, err := store.Create("postgres://localhost/db")
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	// Token is removed by the reaper even though it was never accessed again.
	assert.Equal(t, 1, store.Cleanup())

	_, err = store.Get(token)
	assert.Equal(t, errTokenInvalid, err)
}
