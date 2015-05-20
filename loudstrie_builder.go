package loudstrie

import (
	"sort"

	"github.com/hideo55/go-sbvector"
	"github.com/oleiade/lane"
)

/*
TrieBuilderData holds
*/
type TrieBuilderData struct {
	trie *TrieData
}

/*
TrieBuilder is
*/
type TrieBuilder interface {
	Build(keyList []string, useTailTrie bool) (Trie, error)
	buildTailTrie()
}

type rangeNode struct {
	left  uint64
	right uint64
}

/*
NewTrieBuilder is
*/
func NewTrieBuilder() TrieBuilder {
	builder := &TrieBuilderData{}
	builder.trie = &TrieData{}
	return builder
}

func lg2(x uint64) uint64 {
	ret := uint64(0)
	for x>>ret != 0 {
		ret++
	}
	return ret
}

/*
Build is
*/
func (builder *TrieBuilderData) Build(keyList []string, useTailTrie bool) (Trie, error) {
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
		cur := keyList[left]
		curSize := uint64(len(cur))
		if left+1 == right && depth+1 < curSize {
			treeBuilder.PushBack(true)
			terminalBuilder.PushBack(true)
			tailBuilder.PushBack(true)
			tail := cur[depth:curSize]
			trie.vtails = append(trie.vtails, tail)
			continue
		} else {
			tailBuilder.PushBack(false)
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
		prevC := keyList[prev][depth]
		degree := uint64(0)
		for i := prev; ; i++ {
			if i < right && prevC == keyList[i][depth] {
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
			prevC = keyList[prev][depth]
		}
		treeBuilder.PushBack(true)
	}
	trie.louds, _ = treeBuilder.Build(true, true)
	trie.terminal, _ = terminalBuilder.Build(true, false)
	trie.tail, _ = tailBuilder.Build(false, false)

	if useTailTrie {
		builder.buildTailTrie()
	}

	return trie, nil
}

func (builder *TrieBuilderData) buildTailTrie() {
	origTails := builder.trie.vtails
	vtailTrieBuilder := NewTrieBuilder()
	keyList := make([]string, len(origTails))
	for tailIdx, tail := range origTails {
		runes := []rune(tail)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		keyList[tailIdx] = string(runes)
	}

	tailTrie, _ := vtailTrieBuilder.Build(keyList, false)
	builder.trie.tailTrie = tailTrie
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
	seen := make(map[string]bool)

	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = true
		}
	}
	return result
}
