package core

import "time"

// Identity ...
type Identity struct {
	Timestamp int64
	Name      string
	Prefix    string
	Index     int
}

// NewIdentity ...
func NewIdentity(prefix string) Identity {
	return Identity{
		Timestamp: time.Now().Unix(),
		Prefix:    prefix,
	}
}
