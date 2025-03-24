package geni

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"golang.org/x/oauth2"
)

type cachingTokenSource struct {
	filePath string
	new      oauth2.TokenSource

	mu sync.Mutex
}

// NewCachingTokenSource returns a TokenSource that caches the token from the
// provided TokenSource. It is based on the implementation of
// oauth2.ReuseTokenSource.
func NewCachingTokenSource(filePath string, src oauth2.TokenSource) oauth2.TokenSource {
	return &cachingTokenSource{
		filePath: filePath,
		new:      src,
	}
}

func (s *cachingTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, err := loadTokenFromDisk(s.filePath); err == nil && t.Valid() {
		return t, nil
	}

	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}

	_ = saveTokenToDisk(s.filePath, t)

	return t, nil
}

func loadTokenFromDisk(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var t oauth2.Token
	if json.NewDecoder(f).Decode(&t) != nil {
		return nil, errors.New("failed to parse token file")
	}
	return &t, nil
}

func saveTokenToDisk(p string, t *oauth2.Token) error {
	err := os.MkdirAll(path.Dir(p), os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("failed to create credential cache directory, %w", err)
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(t)
}
