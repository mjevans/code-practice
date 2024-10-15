// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=63
https://projecteuler.net/minimal=63

<p>The $5$-digit number, $16807=7^5$, is also a fifth power. Similarly, the $9$-digit number, $134217728=8^9$, is a ninth power.</p>
<p>How many $n$-digit positive integers exist which are also an $n$th power?</p>


*/
/*

The real question: Given the form a**b (a^b) AND base 10 numbers: at what point does increasing (a) outpace the result length, and (b) underflow the result length?

10^n is exactly equal to N+1 ; 11 grows bigger, and 9 grows smaller, but it takes a LONG time.
9^n takes so long to get small it overflows 64bit uints (thanks Kcalc) ; but golang has math.big
8^20 fits, and should have less than 20 digits (Kcalc)
*/

import (
	// "bufio"
	// "euler"
	"fmt"
	"math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler0063() uint64 {
	var count, bb, ee, prev uint64
	// Special case, 1^1 is 10^0
	fmt.Printf("1 digits 1\t= special case, Log10(1.0) is 0. The trick fails for exact E+/- numbers\n")
	count = 1
	for bb = 2; bb < 10; bb++ {
		prev = ^prev
		for ee = 1; ; ee++ {
			num := math.Pow(float64(bb), float64(ee))
			if uint64(num) == prev {
				break
			}
			digits := math.Ceil(math.Log10(num))
			if ee == uint64(digits) {
				fmt.Printf("%d digits %.0f\t= %d^%d\n", ee, num, bb, ee)
				count++
			}
			if ee < uint64(digits) {
				break
			}
			prev = uint64(num)
		}
	}
	return count
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 63 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

18 digits 150094635296999136    = 9^18
19 digits 1350851717672992000   = 9^19
20 digits 12157665459056928768  = 9^20
21 digits 109418989131512365056 = 9^21
Euler 63: Powerful Digit Counts: 49

*/
func main() {
	var a uint64
	//test

	//run
	a = Euler0063()
	fmt.Printf("Euler 63: Powerful Digit Counts: %d\n", a)
}
