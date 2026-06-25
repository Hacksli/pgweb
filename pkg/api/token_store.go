package api

import (
	"sync"
	"time"

	"github.com/tuvistavie/securerandom"
)

// TokenStore holds short-lived mappings of opaque access tokens to database
// connection strings. Tokens are issued via POST /api/token and later
// exchanged for an authenticated session via GET /token/:token.
//
// Connection strings (which may contain credentials) are kept in memory only
// and expire after the configured TTL.
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

// Get returns the connection string for a token. Unknown or expired tokens
// return errTokenInvalid; expired entries are removed on access.
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

	return entry.url, nil
}
