package perm

import (
	"ajz.xyz/fastperm/rand"
)

// Item alias for permutation item type
type Item interface{}

// Rand wrapper to rand.Rand type
type Rand rand.Rand

// Slice permute a slice in place
// Avoiding division makes the shuffle twice as fast.
// see math/rand.int31n for library implementation
func (r *Rand) Slice(dst []Item) {
	n := len(dst)
	rg := rand.Rand(*r)
	i := n - 2
	for ; i > 0; i-- {
		v := rg.Uint32()
		j := int((uint64(v) * uint64(i+1)) >> 32)
		dst[i], dst[j] = dst[j], dst[i]
	}
}
