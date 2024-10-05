// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=36
https://projecteuler.net/minimal=36

<p>The decimal number, $585 = 1001001001_2$ (binary), is palindromic in both bases.</p>
<p>Find the sum of all numbers, less than one million, which are palindromic in base $10$ and base $2$.</p>
<p class="smaller">(Please note that the palindromic number, in either base, may not include leading zeros.)</p>

*/
/*

585 = 0b10_0100_1001

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


Revisit:  Oops, leading zeros... yeah, I can do that.

However it was SO UGLY and convoluted that I just wrote it again from scratch, and tested every decimal number as a palindrome first.

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

func Euler036__works_but_ugly(limit uint64) uint64 {
	if 0 == limit {
		limit = 1_000_000
	}
	var sum, ii, even, odd, res uint64
	ii = 1
	for {
		// Add 00
		if ii < 100 {
			even = euler.PalindromeMakeDec(ii, 1, false)
			odd = euler.PalindromeMakeDec(ii, 1, true)
			if limit >= even {
				res = euler.PalindromeFlipBinary(even)
				if even == res {
					sum += even
					fmt.Printf("Found: %d\t\t%d\tEven\n", res, sum)
				}
			}
			res = euler.PalindromeFlipBinary(odd)
			if limit >= odd {
				if odd == res {
					sum += odd
					fmt.Printf("Found: %d\t\t%d\tOdd\n", res, sum)
				}
			} else {
				break
			}
		}

		// Add 0000
		if ii < 10 {
			even = euler.PalindromeMakeDec(ii, 2, false)
			if limit >= even {
				res = euler.PalindromeFlipBinary(even)
				if even == res {
					sum += even
					fmt.Printf("Found: %d\t\t%d\tEven\n", res, sum)
				}
			}
		}

		// As is
		even = euler.PalindromeMakeDec(ii, 0, false)
		odd = euler.PalindromeMakeDec(ii, 0, true)
		if limit >= even {
			res = euler.PalindromeFlipBinary(even)
			if even == res {
				sum += even
				fmt.Printf("Found: %d\t\t%d\tEven\n", res, sum)
			}
		}
		res = euler.PalindromeFlipBinary(odd)
		if limit >= odd {
			if odd == res {
				sum += odd
				fmt.Printf("Found: %d\t\t%d\tOdd\n", res, sum)
			}
		} else {
			break
		}

		ii += 2
	}
	fmt.Printf("Exit conditions: ii = %d\todd = %d\teven = %d\n", ii, odd, even)
	return sum
}

func Euler036(limit uint64) uint64 {
	if 0 == limit {
		limit = 1_000_000
	}
	var sum, ii uint64
	for ii = 1; ii < limit; ii += 2 {
		if euler.IsPalindrome(ii, 10) && ii == euler.PalindromeFlipBinary(ii) {
			sum += ii
			fmt.Printf("Found: %d\t\t%d\n", ii, sum)
		}
	}
	fmt.Printf("Exit conditions: ii = %d\t%d\n", ii, sum)
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
Exit conditions: ii = 11        25
Euler036: TEST: 10 == 25?       25
Found: 1                1
Found: 3                4
Found: 5                9
Found: 7                16
Found: 9                25
Found: 33               58
Found: 99               157
Found: 313              470
Found: 585              1055
Found: 717              1772
Found: 7447             9219
Found: 9009             18228
Found: 15351            33579
Found: 32223            65802
Found: 39993            105795
Found: 53235            159030
Found: 53835            212865
Found: 73737            286602
Found: 585585           872187
Exit conditions: ii = 1000001   872187
Euler036: RUN : 1_000_000 ==    872187


*/
func main() {
	//test
	sum := Euler036(10)
	fmt.Printf("Euler036: TEST: 10 == 25?\t%d\n", sum)

	//run
	sum = Euler036(1_000_000)
	fmt.Printf("Euler036: RUN : 1_000_000 ==\t%d\n", sum)
}
