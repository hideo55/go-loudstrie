loudstrie
=========

[![Build Status](https://travis-ci.org/hideo55/go-loudstrie.svg?branch=master)](https://travis-ci.org/hideo55/go-loudstrie)
[![Godoc](https://godoc.org/github.com/hideo55/go-loudstrie?status.png)](https://godoc.org/github.com/hideo55/go-loudstrie)
[![Coverage Status](https://coveralls.io/repos/hideo55/go-loudstrie/badge.svg?branch=master)](https://coveralls.io/r/hideo55/go-loudstrie?branch=master)

Description
-----------

LOUDS(Level-Order Unary Degree Sequence) Trie implementation for Go

Installation
------------

This package can be installed with the go get command:

    go get github.com/hideo55/go-loudstrie

Usage
------

```go
import (
    "fmt"

    "github.com/hideo55/go-loudstrie"
)

func main() {
    keyList := []string{
        "bbc",
        "able",
        "abc",
        "abcde",
        "can",
    }

    trie, err := loudstrie.NewTrie(keyList, true)
    if err != nil {
        // Failed to build trie.
    }

    // Common prefix search
    searchKey := "abcde"
    result := trie.CommonPrefixSearch(searchkey, 0)
    for _, item := range result {
        // item has two menbers, ID and Length.
        // ID: ID of the key.
        // Length: Length of the key string.
        key, _ := trie.DecodeKey(item.ID)// key == searchKey[:item.Length]
        fmt.Printf("ID:%d, key:%s, len:%d\n", item.ID, key, item.Length)
    }
}
```

Documentation
-------------

[API documentation](http://godoc.org/github.com/hideo55/go-loudstrie)

Supported version
-----------------

Go 1.4 or later

License
--------

MIT License
