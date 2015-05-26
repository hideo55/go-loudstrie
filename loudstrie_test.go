package loudstrie

import (
	"testing"
)

func TestBuild(t *testing.T) {
	builder := NewTrieBuilder()
	keyList := []string{
		"bbc",
		"able",
		"abc",
		"abcde",
		"can",
	}
	trie, err := builder.Build(keyList, false)
	if err != nil {
		t.Error("Build error")
	}
	if trie.GetNumOfKeys() != uint64(len(keyList)) {
		t.Error("")
	}
}

func TestExactMatchSearch(t *testing.T) {
	builder := NewTrieBuilder()
	keyList := []string{
		"bbc",
		"able",
		"abc",
		"abcde",
		"can",
	}
	trie, _ := builder.Build(keyList, true)
	for _, key := range keyList {
		id := trie.ExactMatchSearch(key)
		decode := trie.DecodeKey(id)
		if key != decode {
			t.Error("Expected", key, "got", decode)
		}
	}
}

func TestCommonPrefixSearch(t *testing.T) {
	builder := NewTrieBuilder()
	keyList := []string{
		"bbc",
		"able",
		"abc",
		"abcde",
		"can",
	}
	trie, _ := builder.Build(keyList, true)
	results := make([]Result, 0)
	trie.CommonPrefixSearch("abcde", &results, 100)
	if len(results) != 2 {
		t.Error(results)
	}
	str := trie.DecodeKey(results[0].ID)
	if str != "abc" {
		t.Error(str)
	}
	str = trie.DecodeKey(results[1].ID)
	if str != "abcde" {
		t.Error(str)
	}

}

func TestPredictiveSearch(t *testing.T) {
	builder := NewTrieBuilder()
	keyList := []string{
		"bbc",
		"able",
		"abc",
		"abcde",
		"can",
	}
	trie, _ := builder.Build(keyList, true)
	results := make([]Result, 0)
	trie.PredictiveSearch("ab", &results, 100)
	if len(results) != 3 {
		t.Error(len(results))
	}
}
