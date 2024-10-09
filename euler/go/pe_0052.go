// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=52
https://projecteuler.net/minimal=52

<p>It can be seen that the number, $125874$, and its double, $251748$, contain exactly the same digits, but in a different order.</p>
<p>Find the smallest positive integer, $x$, such that $2x$, $3x$, $4x$, $5x$, and $6x$, contain the same digits.</p>

*/
/*

Criteria:
* Smallest Positive Integer
* 'all' the digits must be in...
x 1 2 3 4 5 6
??? Can additional digits be in the other numbers?

For the moment, assume they all must be reused.

*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler052(limit, muls, base uint64) uint64 {
	var ii, mm uint64
	for ii = 0; ii < limit; ii++ {
		iiv := append(euler.Uint64ToDigitsUint8(ii, base), 1)
		iin := euler.Uint8DigitsToUint64(iiv, base)
		iisrt := euler.Uint8CopyInsertSort(iiv)
		for mm = 2; mm <= muls; mm++ {
			iimn := iin * mm
			iimv := euler.Uint64ToDigitsUint8(iimn, base)
			iimsrt := euler.Uint8CopyInsertSort(iimv)
			if 0 != euler.Uint8Compare(iisrt, iimsrt) {
				mm = 0
				break
			}
		}
		if mm >= muls {
			return iin
		}
	}
	return 0
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 52 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done



*/
func main() {
	var l, m, b, res uint64
	//test

	//run
	l, m, b = 1_000_000, 6, 10
	res = Euler052(l, m, b)
	fmt.Printf("Euler 52: Permuted Multiples: base: %d\tmuls: %d\tbase: %d\n", b, m, res)
}
