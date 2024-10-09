// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=53
https://projecteuler.net/minimal=53

<p>There are exactly ten ways of selecting three from five, 12345:</p>
<p class="center">123, 124, 125, 134, 135, 145, 234, 235, 245, and 345</p>
<p>In combinatorics, we use the notation, $\displaystyle \binom 5 3 = 10$.</p>
<p>In general, $\displaystyle \binom n r = \dfrac{n!}{r!(n-r)!}$, where $r \le n$, $n! = n \times (n-1) \times ... \times 3 \times 2 \times 1$, and $0! = 1$.
</p>
<p>It is not until $n = 23$, that a value exceeds one-million: $\displaystyle \binom {23} {10} = 1144066$.</p>
<p>How many, not necessarily distinct, values of $\displaystyle \binom n r$ for $1 \le n \le 100$, are greater than one-million?</p>

*/
/*

I'm going to re-write it into text friendlier format...
Cmb(5 3) = 10
Cmb(n r) = (n!) / (r! * (n-r)!)  REQ: r < n

Looking at this, the 'worst' case, the smallest denominator, is when r ~= n / 2 since n-r will be ~n/2 and the floor will be on either side...
R will round up to the next whole number, so r = floor((n+1)/2.0)

So, n! will be 'cut' by ((n+1)/2)! ; which is why R isn't given, the worst case is driven by R's relation to N.

Kcalc
	22!/10!/12! = 646646
	22!/11!/11! = 705432

	23!/ 9!/14! = 817190
4
	23!/10!/13! = 1144066
	23!/11!/12! = 1352078
	23!/12!/11! = 1352078
	23!/13!/10! = 1144066
7
	24!/ 9!/15! = 1307504
	24!/10!/14! = 1961256
	24!/11!/13! = 2496144
	24!/12!/12! = 2704156
	24!/13!/11! = 2496144
	24!/14!/10! = 1961256
	24!/15!/ 9! = 1307504
10
	25!/17!/ 8! = 1081575 16:9 15:10 14:11 13:12 ||
11
	26!/18!/ 8! = 1562275 17:9 16:10 15:11 14:12 13||13
...
	100!/96!/4! GT || LT 97 98 99 100


20! < 2^64, 21! > 2^64

How many values of the combinatorial 1..n are greater than 1_000_000

*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler053(nMx, gt uint64) uint64 {
	if 2 > nMx {
		return 0
	}
	var ret, ii, left, pos, right uint64
	cut := big.NewInt(int64(gt))
	for ii = 2; ii <= nMx; ii++ {
		right = ii
		left = ii >> 1
		pos = left + left>>1
		var cmp, prev int
		for {
			bg := euler.FactorialDivFactU64toBig(ii, pos)
			bg.Div(bg, euler.FactorialDivFactU64toBig(ii-pos, 1))

			prev = cmp
			cmp = bg.Cmp(cut)
			if left >= right || pos >= right {
				// fmt.Printf("%d = pos = %d (cmp = %d) \t %s\n", ii, pos, cmp, bg)
				break
			}
			// fmt.Printf("%d : pos = %d (cmp = %d) \t %s\n", ii, pos, cmp, bg)
			// cmp ~= bg - cut
			// fmt.Printf("%d lpr %d: %d %d %d\n", ii, cmp, left, pos, right)
			if 0 < cmp {
				left = pos - 1
				pos += (right + 1 - pos) >> 1
			} else if 0 > cmp {
				right = pos - 1 // Reach the high edge, might overshoot by one, hence prev and the break after the compare to figure out where it got.
				pos -= (pos + 1 - left) >> 1
			} else {
				// fmt.Printf("%d === pos = %d (cmp = %d) \t %s\n", ii, pos, cmp, bg)
				break
			}

		}
		// fmt.Printf("%d : pos = %d (cmp = %d (%d))\n", ii, pos, cmp, prev)
		// cmp || prev -1 is under the cut +1 is over the cut
		// if cmp == 1 // pos == The right most side of the > target
		// if cmp == -1 but prev == 1, overshot by 1, pull back
		if 0 == cmp || (-1 == cmp && 1 == prev) {
			pos--
			cmp = 1
			// fmt.Printf("Forced\t%d\n", pos)
		}
		// are any values greater than the target?
		if 1 == cmp {
			ok := (pos << 1) - ii + 1
			ret += ok
			// fmt.Printf("%d\tadded = %d\t\t%d\n", ii, ok, ret)
		}
	}
	return ret
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 53 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 53: Combinatoric Selections: 4075


*/
func main() {
	//test
	// Euler053(27, 1_000_000)

	//run
	fmt.Printf("Euler 53: Combinatoric Selections: %d\n", Euler053(100, 1_000_000))
}
