package geni

import "net/http"

// AddOption modifies an outgoing request for the profile/union add-*
// endpoints (add-child, add-sibling, add-partner).
type AddOption func(*http.Request)

// WithModifier sets the relationship_modifier query parameter. An empty
// value is a no-op.
func WithModifier(modifier string) AddOption {
	return func(r *http.Request) {
		if modifier == "" {
			return
		}
		q := r.URL.Query()
		q.Set("relationship_modifier", modifier)
		r.URL.RawQuery = q.Encode()
	}
}
