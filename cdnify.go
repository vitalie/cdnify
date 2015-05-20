package cdnify

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultPrefix = "/assets/"
)

var (
	DefaultTTL = 7 * 24 * time.Hour
)

// SetPrefix configures the prefix.
func SetPrefix(prefix string) func(*Cdnify) {
	return func(m *Cdnify) {
		m.prefix = prefix
	}
}

// SetTTL configures TTL.
func SetTTL(ttl time.Duration) func(*Cdnify) {
	return func(m *Cdnify) {
		m.ttl = ttl
	}
}

// Cdnify represents a Negroni middleware which sets
// caching headers for a prefix to ttl.
type Cdnify struct {
	prefix string
	ttl    time.Duration
	isDev  bool
}

// New returns a new middleware.
func New(isDev bool, opts ...func(*Cdnify)) *Cdnify {
	m := &Cdnify{
		prefix: DefaultPrefix,
		ttl:    DefaultTTL,
		isDev:  isDev,
	}

	// Apply options.
	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Cdnify) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Handle only `GET` requests.
	if r.Method != "GET" {
		next(w, r)
		return
	}

	// Set `Cache-Control` header.
	if !m.isDev && strings.HasPrefix(r.URL.Path, m.prefix) {
		w.Header().Set("Cache-Control",
			fmt.Sprintf("public, max-age=%d",
				int(m.ttl.Seconds())),
		)
	}

	next(w, r)
}
