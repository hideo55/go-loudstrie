package loudstrie

import (
	"sort"

	"github.com/hideo55/go-sbvector"
	"github.com/oleiade/lane"
)

/*
TrieBuilderData holds information of LOUDS Trie Builder
*/
type TrieBuilderData struct {
	trie *TrieData
}

/*
TrieBuilder is interface of LOUDS Trie Builder
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
NewTrieBuilder returns new LOUDS Trie Builder
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
Build builds LOUDS Trie from keyList.
If useTailTrie is true, compress TAIL array.
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
		if q.Empty() {
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
		for i := prev + 1; ; i++ {
			if i < right && prevC == keyList[i][depth] {
				continue
			}
			trie.edges = append(trie.edges, prevC)
			treeBuilder.PushBack(false)
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
		keyList[tailIdx] = reverseString(tail)
	}
	origKeyList := make([]string, len(origTails))
	copy(origKeyList, keyList)
	tailTrie, _ := vtailTrieBuilder.Build(keyList, false)
	builder.trie.tailTrie = tailTrie
	builder.trie.tailIDSize = lg2(tailTrie.GetNumOfKeys())
	tailIDBuilder := sbvector.NewVectorBuilder()
	for _, tail := range origKeyList {
		id, _ := tailTrie.ExactMatchSearch(tail)
		tailIDBuilder.PushBackBits(id, builder.trie.tailIDSize)
	}
	builder.trie.tailIDs, _ = tailIDBuilder.Build(false, false)
	builder.trie.hasTailTrie = true
	builder.trie.vtails = make([]string, 0)
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

func reverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
