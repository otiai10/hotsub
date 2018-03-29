package core

import "fmt"

// Env represent an envirionment variable.
type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Pair returns key=value pair string.
// TODO: better name
func (env Env) Pair() string {
	return fmt.Sprintf("%s=%s", env.Name, env.Value)
}
