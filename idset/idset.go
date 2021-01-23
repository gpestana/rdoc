// Package idset implements a idempotent set of string IDs and helpers to use
// in the context of `rdoc`
package idset

// Set is a set of IDs
type Set struct {
	ids map[string]struct{}
}

// New returns an empty set
func New() *Set {
	return &Set{ids: map[string]struct{}{}}
}

// GetIDs returns a slice with all ids in the set
func (set *Set) GetIDs() []string {
	keys := []string{}
	for k := range set.ids {
		keys = append(keys, k)
	}
	return keys
}

// Add inserts a new id to the set
func (set *Set) Add(id string) {
	set.ids[id] = struct{}{}
}

// Exists checks whether an ID exists in the set
func (set Set) Exists(id string) bool {
	_, exists := set.ids[id]
	return exists
}

// Diff returns all strings in the base doc that do not exist in the Doc
func (set Set) Diff(ids []string) []string {
	var diff []string
	for _, id := range ids {
		if !set.Exists(id) {
			diff = append(diff, id)
		}
	}
	return diff
}
