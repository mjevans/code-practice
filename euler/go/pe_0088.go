// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=88
https://projecteuler.net/minimal=88

<p>A natural number, $N$, that can be written as the sum and product of a given set of at least two natural numbers, $\{a_1, a_2, \dots, a_k\}$ is called a product-sum number: $N = a_1 + a_2 + \cdots + a_k = a_1 \times a_2 \times \cdots \times a_k$.</p>
<p>For example, $6 = 1 + 2 + 3 = 1 \times 2 \times 3$.</p>
<p>For a given set of size, $k$, we shall call the smallest $N$ with this property a minimal product-sum number. The minimal product-sum numbers for sets of size, $k = 2, 3, 4, 5$, and $6$ are as follows.</p>
<ul style="list-style-type:none;">
<li>$k=2$: $4 = 2 \times 2 = 2 + 2$</li>
<li>$k=3$: $6 = 1 \times 2 \times 3 = 1 + 2 + 3$</li>
<li>$k=4$: $8 = 1 \times 1 \times 2 \times 4 = 1 + 1 + 2 + 4$</li>
<li>$k=5$: $8 = 1 \times 1 \times 2 \times 2 \times 2 = 1 + 1 + 2 + 2 + 2$</li><li>$k=6$: $12 = 1 \times 1 \times 1 \times 1 \times 2 \times 6 = 1 + 1 + 1 + 1 + 2 + 6$</li></ul>
<p>Hence for $2 \le k \le 6$, the sum of all the minimal product-sum numbers is $4+6+8+12 = 30$; note that $8$ is only counted once in the sum.</p>
<p>In fact, as the complete set of minimal product-sum numbers for $2 \le k \le 12$ is $\{4, 6, 8, 12, 15, 16\}$, the sum is $61$.</p>
<p>What is the sum of all the minimal product-sum numbers for $2 \le k \le 12000$?</p>

/
*/
/*
	I don't have to store all the terms out in memory like that...
	E.G. the type Factorized I wrote for working with prime power numbers can also work, though most of the special functions assume 'reduced' forms, not an 'improper fraction' like other numbers form.
	The core issue revolves around a balance point between:
	* All 1s 1x1x...x1 = 1 but sum() = k*1
	* Some multiplied term in a few of those slots which EQUALS than the sum of that number's terms and the remaining 1s...

	Clearly the floor cannot be under K, but finding sum/product N which is the smallest seems difficult, more so when factors can be combined in multitudes of ways (proper divisors)....
	Uggh, they're trying to blindside with the coin problem and a limited set of coins which change dynamically as some are utilized.

	The example should be studied more closely for inspiration on algorithms...
	2	1+1 = 2 ne, 1+2 = 3 not factors, 2*2 = 4 -- done
	3	2*2*2 = 8 -- too high.  2 is under sum, 1*2*2 = 4 is under sum... 2*2*2=8 over sum.
	4	1*1*2*2 = 4 vs 6, 1*2*2*2 = 8 sum 7 -- the answer was 1,1,2,4 = 8
	5	1*1*1*2*2 = 4 vs 7, 1*1*2*2*2 = 8 sum 8 -- done
	6	1*1*1*2*2*2 = 8 (9),1*1*2*2*2*2 = 16 (10) too high
	7	12 = 1^5*3*4 ~~ 1*1*1*1*2*2*2 = 8 (10) low, 1*1*1*2*2*2*2 = 16 (11) high, 11 prime, 13 prime, 12 is 2*2*3 (11) or 4*3 () or 2*6
	8	12 is 2*2*3*1^5 (12)		1^5*2^3 = 8 (11), 1^4*2^4 = 16 (12), 9 is 3*3 * 1^6 (12), 10 is 2*5 * 1^6 (13), 11 is prime, 12 is is 2*2*3*1^5 (12)
	9?  15?	3*5*1^7		1^5*2^4 = 8 (13) lower, 1^4*2^5 = 16 (14) higher scan between too low and too high
	12? 16?	16 = 1*1*1 * 1*1*1 * 1*1*2 * 2*2*2 = 16
	So I think 10 and 11 aren't answers
	10	1^7*2^3 = 8 (13), 1^6*2^4 = 16 (14), 9 is 3x3*1^8 (14), 10 is 2*5*1^8 (15), 11 is prime, 12 is is 2*2*3*1^7 (14) or 4*3*1^8 (15) or 2*6*1^8 (16), 13 is prime, 14 is 2*7*1^8 (17), 15 is 3*5*1^8 (16) exhausted
	11	1^8*2^3 = 8 (14), 1^7*2^4 = 16 (15), 9 is 3x3*1^9 (15), 10 is 2*5*1^9 (16), 11 is prime, 12 is is 2*2*3*1^8 (15) or 4*3*1^9 (16) or 2*6*1^9 (17), 13 is prime, 14 is 2*7*1^9 (18), 15 is 3*5*1^9 (17) exhausted

	So far it's reasonable to start with a power of two apart as a set of bounds, but that gets very wide very quick.

	What if I swapped two of the upper bounds 2s with a 1 and a 3?  That should be 75% of the number but the same total.

	11	1^8*2^3 = 8 (14), 1^7*2^4 = 16 (15), 1^8*3*2^2 (15) = 12 that is a better lower bound,
		How can I verify that though?  1^7*3*2^3 (16) = 24 higher than my previous high, not great so far... but this reminds me a bit of how square roots are approximated.
		1^8*3^2*2 (16) = 12
		I feel mostly confident that the '75%' substitution technique (2*2 -> 3*1) won't overlook better solutions since 2 and 3 are both prime, and overall it doesn't change the total, and thus N...
		However once other numbers are introduced to the system either the sum (1,4), result (?), or both (1,5) increase.
		There seems to be a special relationship between the multiplicative identity (1), it's doubling / the first prime number (2), and it's tripling / second prime number (3).  Further in the summation it's also special that (1,3) and (2,2) can swap to change the multiplication result without modifying the sum.

	The initial approach solved the trivial test cases quickly, and I was hopeful to estimate a narrow solution range.  About 10% of the way to the answer and 15 min into running it's clear that there's a LOT of wasted factorization effort, the same numbers tried with slightly different K term numbers.  Plus the estimate is at best 10-20% of the input number, rather than something better like always less than 20.

	I need to take an entirely different approach, but it's 2:30 am.

	I'm waffling on if I was on the wrong track or not, this might be one where I am better off allowing a brute force to run and looking over the discussion to see where my existing knowledge was lacking.

	A quick review of the known set to see if any patterns pop out that I missed:
	K	Answer
	2	4	2,2
	3	6	1,2,3
	4	8	1,1,2,4
	5	8*	1,1,2,2,2
	6	12	1,1,1,1,2,6
	7	12*	1,1,1,1,1,2,6
	8	12*	1,1,1,1,1,2,2,3
	9	15	1,1,1,1,1,1,1,3,5
	10	-	--5,2,1^8--
	11	-	--11,1^10--
	12	16	1,1,1,1,1,1,1,1,2,2,2,2

	Q: Is it possible to have more factors than K slots?
	A: If that were going to happen it'd be with the Power of 2 test, but 2^(k) > 2*k (for k > 1) and 2^(k-n) > 2*k+n (for k >= n >= 0)

	However, it strongly looks like the focus should be on the numbers to _factor_ rather than approximating any limits; given they increase.  That would also greatly reduce duplicated work.
/
*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0088(min, max uint64) uint64 {
	// I shouldn't need these, the default base should suffice with Rho and Lenstra guarded by a Probably Prime test pair known to be accurate for all 64 bit numbers.
	// t := euler.SqrtU64(uint64(max))
	// fmt.Printf("Euler 88: Primes.Grow(%d + 8)\n", t)
	// euler.Primes.Grow(t + 8)

	// 16 bit int is enough for Euler 88
	check := make(map[uint16]uint8)

	var ii, u2p, u3p, l2p, l3p, m2p, m3p, ret, tsum, tmul uint64

	addPSN := func(n uint64) {
		if _, exists := check[uint16(n)]; !exists {
			ret += n
			check[uint16(n)] = 1
		}
	}

	// NOTE: ProperDivisors() returns a slice that starts with 1.. so element 0 is useless here

	var SumMul func(k, n, sum, mul, divMx uint64, divs []uint64) bool
	SumMul = func(k, n, sum, mul, divMx uint64, divs []uint64) bool {
		// Every time N is divided the sum is reduced, as a+b < a*b
		if 0 == k {
			return false
		}
		dMx := divMx
		tmul := mul * divs[dMx]
		for tmul > n && dMx > 1 {
			dMx--
			tmul = mul * divs[dMx]
		}
		tsum := sum - 1 + divs[dMx]

		// fmt.Printf("debug SumMul: %t %d == %d (%d, %d, %d, %d * %d)\n", tsum == tmul, tsum, tmul, k, n, sum, mul, divs[dMx])
		// if the answer was found OR if a sub-iteration finds the answer...
		if (tsum == n && tmul == n) || SumMul(k-1, n, tsum, tmul, dMx, divs) {
			return true
		}

		// Otherwise divMx started too big, reduce by one and try again
		for 1 < divMx {
			divMx--
			if SumMul(k, n, sum, mul, divMx, divs) {
				return true
			}
		}
		return false
	}

euler0088outer:
	for ii = min; ii <= max; ii++ {
		if 0 == ii&0xFF {
			fmt.Println("Euler 88: running %6d current total %7d\n", ii, ret)
		}
		// Find the initial upper and lower bound
		for u2p = 0; ii+u2p > 1<<u2p; u2p++ {
		}
		// bp is currently the bit shift of the upper bound
		l2p, l3p, u3p = u2p-1, 0, 0

		// Test and add possible answers: for sum power of two doubles the number and power of 3 triples.  For mul, ii isn't even a factor
		tsum = ii + u2p                   // + u3p + u3p
		if tsum == euler.PowInt(2, u2p) { // * euler.PowInt(3, u3p) {
			// fmt.Printf("debug: %d add high %d\n", ii, tsum)
			addPSN(tsum)
			continue
		}
		tsum = ii + l2p                   // + l3p + l3p
		if tsum == euler.PowInt(2, l2p) { // * euler.PowInt(3, l3p) {
			// fmt.Printf("debug: %d add low  %d\n", ii, tsum)
			addPSN(tsum)
			continue
		}

		// m2p, m3p = ii+l2p+l3p+l3p, ii+u2p+u3p+u3p
		// fmt.Printf("Euler 88 : %d\t range: %d (%d - %d) [%d > %d]\n", ii, m3p-m2p, m3p, m2p, euler.PowInt(2, u2p)*euler.PowInt(3, u3p), euler.PowInt(2, l2p)*euler.PowInt(3, l3p))

		// Use that 75% of the product trick to shrink the range as much as possible...
		for u2p >= 2 {
			m2p, m3p = u2p-2, u3p+1
			tsum, tmul = ii+m2p+m3p+m3p, euler.PowInt(2, m2p)*euler.PowInt(3, m3p)
			if tsum == tmul {
				// fmt.Printf("debug: %d add mid  %d\n", ii, tsum)
				addPSN(tsum)
				continue euler0088outer
			} else if tsum > tmul {
				l2p, l3p = m2p, m3p
				break
			} else {
				// fmt.Printf("Old %d (%d, %d)\n", ii+u2p+u3p+u3p, u2p, u3p)
				u2p, u3p = m2p, m3p
			}
		}

		// The ranges need to use the multiplied values
		m3p, m2p = euler.PowInt(2, u2p)*euler.PowInt(3, u3p), euler.PowInt(2, l2p)*euler.PowInt(3, l3p)
		// if m2p < ii {
		//	m2p = ii
		// }
		fmt.Printf("Euler 88 : %5d\t range: %3d [%5d .. %5d]\n", ii, m3p-m2p, m2p, m3p)

		for m2p = m2p; m2p <= m3p; m2p++ {
			fact := euler.Primes.Factorize(m2p)
			if fact.Lenbase == 1 && fact.Fact[0].Base == uint32(m2p) {
				continue // number was prime
			}

			// this is where it's sort of like the coin purse sort of puzzles, except combining two coins yields a different value which wrecks dynamic programming recursion.
			// It's different enough that making change doesn't quite seem worth it
			divs := *(fact.ProperDivisors())

			tsum = m2p - ii // target sum, less the not used slots; logically the largest sum cannot be greater than this.  No, it's not off by one, at LEAST one prime factor will split out.
			// fmt.Println(divs)
			// fmt.Println(tsum)
			for tmul = uint64(len(divs)) - 1; tsum < divs[tmul] && 0 < tmul; tmul-- {
				// fmt.Printf("%d tsum %d < %d\n", m2p, tsum, divs[tmul])
			}

			// Only the factors that could _possibly_ fit
			divs = divs[:tmul+1]

			//    func(k, n, sum, mul, divMx uint64, divs []uint64) bool
			if SumMul(ii, m2p, ii, 1, tmul, divs) { // returns true / false for if it solves the puzzle
				// fmt.Printf("debug: %d add sm   %d\n", ii, tsum)
				addPSN(m2p)
				continue euler0088outer
			}
			// loop
		}
		//
	}

	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 88 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

.
*/
func main() {
	var r uint64
	//test
	r = Euler0088(2, 6)
	if 30 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	r = Euler0088(2, 12)
	if 61 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	fmt.Printf("Euler 88: Passed pretests\n")

	//run
	r = Euler0088(2, 12000)
	fmt.Printf("Euler 88: Product-sum Numbers: %d\n", r)
	if 1097343 != r {
		panic("Did not reach expected value.")
	}
}
