# FlyLSH

`FlyLSH` is a [Go](http://golang.org/) implementation of the fly’s olfactory circuits algorithm.

`flylsh` is a novel variant of Locality Sensitive Hashing reveals how a fruit fly distinguishes different type odors
and make corresponding and proposed in article ['A neural algorithm for a fundamental computing problem'](http://science.sciencemag.org/cgi/rapidpdf/358/6364/793?ijkey=aX3uts9Y4xqPE&keytype=ref&siteid=sci)
published by [S. Dasgupta](http://cseweb.ucsd.edu/~dasgupta/), [C. F. Stevens](https://www.salk.edu/scientist/charles-f-stevens/), and [S. Navlakha](http://www.snl.salk.edu/~navlakha/) (2017).

According the research, `flylsh` gains 30~50% accurate improvement compares to traditional LSH method.

LSH (Locality Sensitive Hashing) is a kind of algorithm with the useful property that similar samples produce
similar hashes. That's why LSH is widely used in large scale search engine for similarity search, such as Google.


# Installation

```
go get -u github.com/pantsing/flylsh
```

# Usage

Example usage:

```go
package main

import (
	"fmt"
	"github.com/pantsing/flylsh"
)

func main() {
	var docs = [][]byte{
		[]byte("this is a test phrase"),
        []byte("this is a test phrass"),
        []byte("this is one of test phrases"),
        []byte("different test phrase"),
	}

	hashes := make([][]byte, len(docs))
	for i, d := range docs {
		hashes[i] = flylsh.Sum(d)
		fmt.Printf("flylsh of %s: %x\n", d, hashes[i])
	}

	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[1], flylsh.Compare(hashes[0], hashes[1]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[2], flylsh.Compare(hashes[0], hashes[2]))
	fmt.Printf("Comparison of `%s` and `%s`: %d\n", docs[0], docs[3], flylsh.Compare(hashes[0], hashes[3]))
}
```

Output:

```
flylsh of this is a test phrase: 00002008200020000000004000000002040000080100080400004000001008000080000000020400
flylsh of this is a test phrass: 00002008200020000000004000000002040000080100080400004000001008000080000000020400
flylsh of this is one of test phrases: 00802008200000000000004000000002040010080100080400000000001000000080000000020400
flylsh of different test phrase: 00802000200002004000000000200002040010000000080400000000001000000090000000020404
Comparison of `this is a test phrase` and `this is a test phrass`: 0
Comparison of `this is a test phrase` and `this is one of test phrases`: 5
Comparison of `this is a test phrase` and `different test phrase`: 14
```
