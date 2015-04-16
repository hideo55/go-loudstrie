package loudstrie

import (
	"bytes"
	"sort"
	"unsafe"

	"github.com/hideo55/go-sbvector"
	"github.com/oleiade/lane"
)

type LoudsTrieBuilderData struct {
	trie *LoudsTrieData
}

type LoudsTrieBuilder interface {
	Build(keyList []string, useTailTrie bool) (LoudsTrie, error)
}

type rangeNode struct {
	left  uint64
	right uint64
}

func NewLoudsTrieBuilder() LoudsTrieBuilder {
	builder := new(LoudsTrieBuilderData)
	builder.trie = new(LoudsTrieData)
	return builder
}

func lg2(x uint64) uint64 {
	ret := uint64(0)
	for x>>ret != 0 {
		ret++
	}
	return ret
}

func (builder *LoudsTrieBuilderData) Build(keyList []string, useTailTrie bool) (LoudsTrie, error) {
	trie := builder.trie
	sort.Strings(keyList)
	keyList = removeDuplicates(keyList)
	trie.numOfKeys = uint64(len(keyList))

	q := lane.NewQueue()
	nextQ := lane.NewQueue()

	if trie.numOfKeys != 0 {
		q.Enqueue(rangeNode{0, trie.numOfKeys})
	}

	treeBuilder := sbvector.NewVectorBuilder()
	terminalBuilder := sbvector.NewVectorBuilder()
	tailBuilder := sbvector.NewVectorBuilder()

	treeBuilder.PushBack(false)
	treeBuilder.PushBack(true)

	depth := uint64(0)
	for {
		if q.Size() == 0 {
			tmp := q
			q = nextQ
			nextQ = tmp
			depth++
			if q.Empty() {
				break
			}
		}
		rn := (q.Dequeue()).(rangeNode)
		left := rn.left
		right := rn.right
		cur := *(*[]byte)(unsafe.Pointer(&keyList[left]))
		curSize := uint64(len(cur))
		if left+1 == right && depth+1 < curSize {
			treeBuilder.PushBack(true)
			terminalBuilder.PushBack(true)
			tailBuilder.PushBack(true)
			tail := cur[depth : curSize-1]
			trie.vtails = append(trie.vtails, tail)
		} else {
			treeBuilder.PushBack(false)
		}

		newLeft := left
		if depth == curSize {
			terminalBuilder.PushBack(true)
			newLeft++
			if newLeft == right {
				treeBuilder.PushBack(true)
				continue
			}
		} else {
			terminalBuilder.PushBack(false)
		}

		prev := newLeft
		prevC := (*(*[]byte)(unsafe.Pointer(&keyList[prev])))[depth]
		degree := uint64(0)
		for i := prev; ; i++ {
			if i < right && prevC == (*(*[]byte)(unsafe.Pointer(&keyList[i])))[depth] {
				continue
			}
			trie.edges = append(trie.edges, prevC)
			treeBuilder.PushBack(false)
			degree++
			nextQ.Enqueue(rangeNode{prev, i})
			if i == right {
				break
			}
			prev = i
			prevC = (*(*[]byte)(unsafe.Pointer(&keyList[prev])))[depth]
		}
		treeBuilder.PushBack(true)
	}

	treeBuilder.Build(true, true)
	terminalBuilder.Build(true, false)
	tailBuilder.Build(false, false)

	if useTailTrie {
		builder.buildTailTrie()
	}

	return trie, nil
}

func (builder *LoudsTrieBuilderData) buildTailTrie() {
	origTails := builder.trie.vtails
	vtailTrieBuilder := NewLoudsTrieBuilder()
	var keyList []string
	for _, tail := range origTails {
		runes := []rune(bytes.NewBuffer(tail).String())
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		keyList = append(keyList, string(runes))
	}
	tailTrie, _ := vtailTrieBuilder.Build(keyList, false)
	builder.trie.tailIDSize = lg2(tailTrie.GetNumOfKeys())
	tailIDBuilder := sbvector.NewVectorBuilder()
	for _, tail := range keyList {
		id := tailTrie.ExactMatchSearch(tail)
		tailIDBuilder.PushBackBits(id, builder.trie.tailIDSize)
	}
	builder.trie.tailIDs, _ = tailIDBuilder.Build(false, false)
}

func removeDuplicates(a []string) []string {
	var result []string
	var seen map[string]bool

	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = true
		}
	}
	return result
}
