// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=55
https://projecteuler.net/minimal=55

<p>If we take $47$, reverse and add, $47 + 74 = 121$, which is palindromic.</p>
<p>Not all numbers produce palindromes so quickly. For example,</p>
\begin{align}
349 + 943 &amp;= 1292\\
1292 + 2921 &amp;= 4213\\
4213 + 3124 &amp;= 7337
\end{align}
<p>That is, $349$ took three iterations to arrive at a palindrome.</p>
<p>Although no one has proved it yet, it is thought that some numbers, like $196$, never produce a palindrome. A number that never forms a palindrome through the reverse and add process is called a Lychrel number. Due to the theoretical nature of these numbers, and for the purpose of this problem, we shall assume that a number is Lychrel until proven otherwise. In addition you are given that for every number below ten-thousand, it will either (i) become a palindrome in less than fifty iterations, or, (ii) no one, with all the computing power that exists, has managed so far to map it to a palindrome. In fact, $10677$ is the first number to be shown to require over fifty iterations before producing a palindrome: $4668731596684224866951378664$ ($53$ iterations, $28$-digits).</p>
<p>Surprisingly, there are palindromic numbers that are themselves Lychrel numbers; the first example is $4994$.</p>
<p>How many Lychrel numbers are there below ten-thousand?</p>
<p class="smaller">NOTE: Wording was modified slightly on 24 April 2007 to emphasise the theoretical nature of Lychrel numbers.</p>


*/
/*

To start with, 0 clearly can't go through the process, so they _must_ mean
* positive integers 1..N

Bitvector + record if numbers under the inspected limit become palindromes (not if they ARE) within 50 iterations

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

func Euler055(ceil, iterlim uint64) uint64 {
	bvbound := ((ceil + 1) >> 3) + 1
	bv := make([]uint8, bvbound)
	var Lychrel, iter, test, ii uint64
	for ii = ceil; 0 < ii; ii-- {
		test = ii + euler.Uint8DigitsToUint64(euler.Uint8Reverse(euler.Uint64ToDigitsUint8(ii, 10)), 10)
		iter = iterlim - 1
		for {
			if test < ceil {
				idx, bit := test>>3, test&0b_111
				if 0 < bv[idx]&(uint8(1)<<bit) {
					break // Sequence lead to a number that becomes a palindrome
				}
			}
			if 0 == iter || euler.IsPalindrome(test, 10) {
				break // Done either way
			}
			test += euler.Uint8DigitsToUint64(euler.Uint8Reverse(euler.Uint64ToDigitsUint8(test, 10)), 10)
			iter--
		}
		if 0 != iter {
			idx, bit := ii>>3, ii&0b_111
			bv[idx] |= uint8(1) << bit
		} else {
			Lychrel++
			// fmt.Printf("Found (probable) Lychrel number at %d\ttotal: %d\n", ii, Lychrel)
		}
	}
	return Lychrel
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 55 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 55: Lychrel Numbers 249
*/
func main() {
	//test
	// Euler055()

	//run
	ln := Euler055(9_999, 50)
	fmt.Printf("Euler 55: Lychrel Numbers %d\n", ln)
}
