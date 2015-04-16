package loudstrie

import (
	"github.com/hideo55/go-sbvector"
)

type LoudsTrieData struct {
	louds sbvector.SuccinctBitVector
	terminal sbvector.SuccinctBitVector
	tail sbvector.SuccinctBitVector
	vtails [][]byte
	edges []byte
	numOfKeys uint64
	tailTrie LoudsTrie
	tailIDs sbvector.SuccinctBitVector
	tailIDSize uint64
}

type LoudsTrie interface {
	ExactMatchSearch(key string) uint64
	CommonPrefixSearch(key string, res map[uint64]uint64)
	PredictiveSearch(key string, res map[uint64]uint64)
	Traverse(key string, nodePos *uint64, zeros *uint64, keyPos *uint64) uint64
	DecodeKey(id uint64) string
	GetNumOfKeys() uint64
}

func (trie *LoudsTrieData) ExactMatchSearch(key string) uint64 {
	id := uint64(0)
	return id
}

func (trie *LoudsTrieData) CommonPrefixSearch(key string, res map[uint64]uint64){
}

func (trie *LoudsTrieData) PredictiveSearch(key string, res map[uint64]uint64) {
}

func (trie *LoudsTrieData) Traverse(key string, nodePos *uint64, zeros *uint64, keyPos *uint64) uint64 {
	id := uint64(0)
	return id
}

func (trie *LoudsTrieData) DecodeKey(id uint64) string {
	key := ""
	return key
}

func (trie *LoudsTrieData) GetNumOfKeys() uint64 {
	return trie.numOfKeys
}
