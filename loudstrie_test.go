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
	id := trie.ExactMatchSearch("aaa")
	if id != NotFound {
		t.Error("Expected", NotFound, "got", id)
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
	results := trie.CommonPrefixSearch("abcde", 0)
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

	results = trie.CommonPrefixSearch("abcde", 1)
	if len(results) != 1 {
		t.Error(results)
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

	results := trie.PredictiveSearch("ab", 0)
	if len(results) != 3 {
		t.Error(results)
	}

	results = trie.PredictiveSearch("ab", 1)
	if len(results) != 1 {
		t.Error(results)
	}
	results = trie.PredictiveSearch("can", 0)
	if len(results) != 1 {
		t.Error(results)
	}
	results = trie.PredictiveSearch("d", 0)
	if len(results) != 0 {
		t.Error(results)
	}
	results = trie.PredictiveSearch("cas", 0)
	if len(results) != 0 {
		t.Error(results)
	}
}
