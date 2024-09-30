// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=36

<p>The decimal number, $585 = 1001001001_2$ (binary), is palindromic in both bases.</p>
<p>Find the sum of all numbers, less than one million, which are palindromic in base $10$ and base $2$.</p>
<p class="smaller">(Please note that the palindromic number, in either base, may not include leading zeros.)</p>

If it can't include leading zeros, then it can't include trailing zeros either... All even numbers are out.

0 0b_0000
1 0b_0001 OK
2 0b_0010
3 0b_0011 OK
4 0b_0100
5 0b_0101 OK
6 0b_0110
7 0b_0111 OK
8 0b_1000
9 0b_1001 OK

< 1_000_000 so 6, 5, 4, 3, 2 or 1 digits ; though are the single digits included?  Difference is 25

Next there's the flip around issue, E.G. 1 000? 0001 ; naively 111 (dec) except 111 (dec) in binary is 0b_0110_1111




*/

import (
	// "bufio"
	// "bitvector"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler036(limit uint64) uint64 {
	if 0 == limit {
		limit = 1_000_000
	}
	var sum, ii, even, odd uint64
	for ii = 1; limit > odd; ii += 2 {
		var res uint64
		even = euler.PalindromeMakeDec(ii, false)
		odd = euler.PalindromeMakeDec(ii, true)
		if limit > even {
			res = euler.PalindromeFlipBinary(even)
			if even == res {
				sum += even
				fmt.Printf("Found: %d\t\t%d\n", res, sum)
			}
		}
		res = euler.PalindromeFlipBinary(odd)
		if odd == res {
			sum += odd
			fmt.Printf("Found: %d\t\t%d\n", res, sum)
		}
	}
	return sum
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 36 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Found: 1                1
Found: 3                4
Found: 5                9
Found: 7                16
Found: 9                25
Euler036: TEST: 10 == 25?       25
Found: 1                1
Found: 33               34
Found: 3                37
Found: 5                42
Found: 7                49
Found: 99               148
Found: 9                157
Found: 313              470
Found: 717              1187
Found: 7447             8634
Found: 585              9219
Found: 32223            41442
Found: 53235            94677
Found: 15351            110028
Found: 585585           695613
Found: 73737            769350
Found: 53835            823185
Found: 39993            863178
Euler036: RUN : 1_000_000 ==    863178



*/
func main() {
	//test
	sum := Euler036(10)
	fmt.Printf("Euler036: TEST: 10 == 25?\t%d\n", sum)

	//run
	sum = Euler036(1_000_000)
	fmt.Printf("Euler036: RUN : 1_000_000 ==\t%d\n", sum)

}
