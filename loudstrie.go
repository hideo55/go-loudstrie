package loudstrie

import (
	"github.com/hideo55/go-loudstrie"
)

type LoudsTrie interface {
	ExactMatchSearch(key string) uint64
	CommomPrefixSearch(key string, ret map[uint64]uint64)
	PredictiveSearch(key string) []uint64
}

