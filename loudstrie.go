package loudstrie

import (
	"bytes"
	"github.com/hideo55/go-sbvector"
)

type LoudsTrieData struct {
	louds       sbvector.SuccinctBitVector
	terminal    sbvector.SuccinctBitVector
	tail        sbvector.SuccinctBitVector
	vtails      []string
	edges       []byte
	numOfKeys   uint64
	hasTailTrie bool
	tailTrie    LoudsTrie
	tailIDs     sbvector.SuccinctBitVector
	tailIDSize  uint64
}

type LoudsTrie interface {
	ExactMatchSearch(key string) uint64
	CommonPrefixSearch(key string, res map[uint64]uint64)
	PredictiveSearch(key string, res map[uint64]uint64)
	Traverse(key string, nodePos *uint64, zeros *uint64, keyPos *uint64) uint64
	DecodeKey(id uint64) string
	GetNumOfKeys() uint64
}

const (
	// NotFound indicates `value is not found`
	NotFound uint64 = 0xFFFFFFFFFFFFFFFF
	// CanNotTraverse indicate `subsequent node is not found`
	CanNotTraverse uint64 = 0xFFFFFFFFFFFFFFFE
)

func (trie *LoudsTrieData) ExactMatchSearch(key string) uint64 {
	id := uint64(0)
	return id
}

func (trie *LoudsTrieData) CommonPrefixSearch(key string, res map[uint64]uint64) {
}

func (trie *LoudsTrieData) PredictiveSearch(key string, res map[uint64]uint64) {
}

func (trie *LoudsTrieData) Traverse(key string, nodePos *uint64, zeros *uint64, keyPos *uint64) uint64 {
	id := NotFound
	if *nodePos == NotFound {
		return CanNotTraverse
	}
	return id
}

func (trie *LoudsTrieData) isLeaf(pos uint64) bool {
	val, _ := trie.louds.Get(pos)
	return val
}

func (trie *LoudsTrieData) getParent(c *byte, pos *uint64, zeros *uint64) {
	*zeros = *pos - *zeros + uint64(1)
	*pos, _ = trie.louds.Select0(*zeros - uint64(1))
	if *zeros < uint64(2) {
		return
	}
	*c = trie.edges[*zeros-uint64(2)]
}

func (trie *LoudsTrieData) getChild(c byte, pos *uint64, zeros *uint64) {
	for {
		if trie.isLeaf(*pos) {
			*pos = NotFound
			break
		}
		if c == trie.edges[*zeros-uint64(2)] {
			*pos, _ = trie.louds.Select1(*zeros - uint64(1))
			*pos++
			*zeros = *pos - *zeros + uint64(1)
			break
		}
		*pos++
		*zeros++
	}
}

func (trie *LoudsTrieData) enumerateAll(pos uint64, zeros uint64, retIDs []uint64, limit uint64) {
	ones := pos - zeros
	term, _ := trie.terminal.Get(ones)
	if term {
		rank, _ := trie.terminal.Rank1(ones)
		retIDs = append(retIDs, rank)
	}
	for i := uint64(0); uint64(len(retIDs)) < limit; i++ {
		if ok, _ := trie.louds.Get(pos + i); !ok {
			break
		}
		nextPos, _ := trie.louds.Select1(zeros + 1 - uint64(1))
		nextPos++
		trie.enumerateAll(nextPos, nextPos-zeros-i+uint64(1), retIDs, limit)
	}
}

func (trie *LoudsTrieData) tailMatch(str string, strlen uint64, depth uint64, tailID uint64, retLen *uint64) bool {
	tail := trie.getTail(tailID)
	tailLen := uint64(len(tail))
	if tailLen > (strlen - depth) {
		return false
	}
	for i := uint64(0); i < tailLen; i++ {
		if str[i+depth] != tail[i] {
			return false
		}
	}
	*retLen = tailLen
	return true
}

func (trie *LoudsTrieData) getTail(tailID uint64) string {
	if trie.hasTailTrie {
		id, _ := trie.tailIDs.GetBits(trie.tailIDSize*tailID, trie.tailIDSize)
		tail := trie.tailTrie.DecodeKey(id)
		runes := []rune(tail)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	} else {
		return trie.vtails[tailID]
	}
}

func (trie *LoudsTrieData) DecodeKey(id uint64) string {
	nodeID, _ := trie.terminal.Select1(id)
	pos, _ := trie.louds.Select1(nodeID)
	pos++
	zeros := pos - nodeID
	var keyBuf []byte
	for {
		c := byte(0)
		trie.getParent(&c, &pos, &zeros)
		if pos == 0 {
			break
		}
		keyBuf = append([]byte{c}, keyBuf...)
	}
	key := bytes.NewBuffer(keyBuf).String()
	return key
}

func (trie *LoudsTrieData) GetNumOfKeys() uint64 {
	return trie.numOfKeys
}
