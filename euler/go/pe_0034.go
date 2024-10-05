// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=34
https://projecteuler.net/minimal=34


<p>$145$ is a curious number, as $1! + 4! + 5! = 1 + 24 + 120 = 145$.</p>
<p>Find the sum of all numbers which are equal to the sum of the factorial of their digits.</p>
<p class="smaller">Note: As $1! = 1$ and $2! = 2$ are not sums they are not included.</p>

*/
/*

9! == 362880 soooo things can get big fast.
9!*6 == 2_177_280
9!*7 == 2_540_160
9!*8 == 2_903_040

So... 0! == 1 according to some convention about identity multiplicand results
https://en.wikipedia.org/wiki/Factorial

However It'd be great if the PROBLEM page included that since it isn't easy to infer from the example.  It could have been provided as part of a match fail (but factorial example) number.

This, of course, makes my entire wheel in wheels iterator fail, since leading zeros now 'count' as 1.

*/

import (
	// "bufio"
	// "bitvector"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler034() uint {
	const limit = 2_540_160
	fa := [10]uint{1}
	fac := uint(1)
	for f := uint(1); f < 10; f++ {
		fac *= f
		fa[f] = fac
	}
	fmt.Println(fa, fa[1], fa[4], fa[5], fa[1]+fa[4]+fa[5])
	var sum, faSum, iiSum uint
	for ii := uint(9); ii < limit; ii++ {
		iiSum = ii
		faSum = 0
		for 0 < iiSum {
			faSum += fa[iiSum%10]
			iiSum /= 10
		}
		if ii == faSum {
			sum += faSum
			fmt.Printf("Found: %d\n", faSum)
		}
	}
	return sum
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 34 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

[1 1 2 6 24 120 720 5040 40320 362880] 1 24 120 145
Found: 145
Found: 40585
Euler034: Sum of Digit Factorials: 40730
*/
func main() {
	//test

	//run
	sum := Euler034()
	fmt.Printf("Euler034: Sum of Digit Factorials: %d\n", sum)
}
