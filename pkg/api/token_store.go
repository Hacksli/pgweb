package api

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tuvistavie/securerandom"
)

// TokenStore holds short-lived mappings of opaque access tokens to database
// connection strings. Tokens are issued via POST /api/token and later
// exchanged for an authenticated session via GET /token/:token.
//
// Connection strings (which may contain credentials) are kept in the running
// process memory only — there is no database backing this store. Each token
// expires after the configured TTL; the lifetime slides forward by a full TTL
// on every successful access, and a background worker removes tokens once they
// expire even if they are never accessed again.
var TokenStore = NewTokenStore(24 * time.Hour)

type tokenEntry struct {
	url       string
	expiresAt time.Time
}

type tokenStore struct {
	mu     sync.Mutex
	tokens map[string]tokenEntry
	ttl    time.Duration
}

// NewTokenStore creates a token store that expires entries after ttl.
func NewTokenStore(ttl time.Duration) *tokenStore {
	return &tokenStore{
		tokens: map[string]tokenEntry{},
		ttl:    ttl,
	}
}

// Create stores a connection string and returns a new opaque token for it.
func (s *tokenStore) Create(url string) (string, error) {
	token, err := securerandom.Uuid()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = tokenEntry{
		url:       url,
		expiresAt: time.Now().Add(s.ttl),
	}

	return token, nil
}

// Get returns the connection string for a token and renews its lifetime for
// another full TTL window (sliding expiration). Unknown or already-expired
// tokens return errTokenInvalid; an expired token is deleted on access.
func (s *tokenStore) Get(token string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.tokens[token]
	if !ok {
		return "", errTokenInvalid
	}

	if time.Now().After(entry.expiresAt) {
		delete(s.tokens, token)
		return "", errTokenInvalid
	}

	// Sliding expiration: each successful access extends the token by a full TTL.
	entry.expiresAt = time.Now().Add(s.ttl)
	s.tokens[token] = entry

	return entry.url, nil
}

// Cleanup removes all expired tokens and returns the number removed.
func (s *tokenStore) Cleanup() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	removed := 0
	for token, entry := range s.tokens {
		if now.After(entry.expiresAt) {
			delete(s.tokens, token)
			removed++
		}
	}

	return removed
}

// RunPeriodicCleanup deletes expired tokens on a fixed interval so that a token
// is removed once its TTL elapses even if it is never accessed again.
func (s *tokenStore) RunPeriodicCleanup(logger *logrus.Logger) {
	for range time.Tick(time.Minute) {
		if removed := s.Cleanup(); removed > 0 && logger != nil {
			logger.WithField("removed", removed).Debug("removed expired access tokens")
		}
	}
}
