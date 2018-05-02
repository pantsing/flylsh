package flylsh

import (
	"fmt"
	"testing"
)

func TestFlyLSH_Sum(t *testing.T) {
	var docs = [][]byte{
		[]byte("this is a test phrase"),
		[]byte("this is a test phrass"),
		[]byte("this is one of test phrases"),
		[]byte("different test phrase"),
	}

	hashes := make([][]byte, len(docs))

	for i, d := range docs {
		hashes[i] = Sum(d)
		fmt.Printf("flylsh of %s: %x\n", d, hashes[i])
	}

	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[1], Compare(hashes[0], hashes[1]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[2], Compare(hashes[0], hashes[2]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[3], Compare(hashes[0], hashes[3]))
}
