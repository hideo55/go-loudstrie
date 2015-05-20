package loudstrie

import (
	"testing"
)

func TestBuild(t *testing.T) {
	builder := NewTrieBuilder()
	keyList := []string{
		"foo",
		"bar",
	}
	trie, err := builder.Build(keyList, true)
	if err != nil {
		t.Error()
	}
	_ = trie

	id := trie.ExactMatchSearch("foo")
	if id == NotFound {
		t.Error()
	}
	t.Logf("%d", id)
	id = trie.ExactMatchSearch("bar")
	t.Logf("%d", id)
}
