package main

import "strings"

// There is an error in yaegi when importing packages in benchmarks.
// The workaround we found is to execute two different evals, one for the import,
// one for the rest of the code. The YAEGY comment below is used to cut the code
// for the second eval.

// YAEGI

const size = 100

func main() {
	var s string
	for r := rune(0); r < size*2; r++ {
		if r%2 == 0 {
			s += string(r)
		}
	}
	n := 0
	for r := rune(0); r < size*2; r++ {
		if strings.ContainsRune(s, r) {
			n++
		}
	}
}
