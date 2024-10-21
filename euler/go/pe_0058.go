// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=58
https://projecteuler.net/minimal=58

<p>Starting with $1$ and spiralling anticlockwise in the following way, a square spiral with side length $7$ is formed.</p>
<p class="center monospace"><span class="red"><b>37</b></span> 36 35 34 33 32 <span class="red"><b>31</b></span><br>
38 <span class="red"><b>17</b></span> 16 15 14 <span class="red"><b>13</b></span> 30<br>
39 18 <span class="red"> <b>5</b></span>  4 <span class="red"> <b>3</b></span> 12 29<br>
40 19  6  1  2 11 28<br>
41 20 <span class="red"> <b>7</b></span>  8  9 10 27<br>
42 21 22 23 24 25 26<br><span class="red"><b>43</b></span> 44 45 46 47 48 49</p>
<p>It is interesting to note that the odd squares lie along the bottom right diagonal, but what is more interesting is that $8$ out of the $13$ numbers lying along both diagonals are prime; that is, a ratio of $8/13 \approx 62\%$.</p>
<p>If one complete new layer is wrapped around the spiral above, a square spiral with side length $9$ will be formed. If this process is continued, what is the side length of the square spiral for which the ratio of primes along both diagonals first falls below $10\%$?</p>

*/
/*

Spiral, Euler 28?

Upper Right == 'grid' * 'grid'
'Corners' == 4 ('grid' * 'grid') - 6 * ('grid' - 1)  // Where's the 6 come from?  That's the 1, 2 and 3 'back' the other number of elements on each side

This problem rotated the problem such that the greatest number is in the lower right, rather than the upper right; that doesn't change the tests.

UR := 'grid' * 'grid'
UL := UR - 1 * ('grid' - 1)
LL := UL - 1 * ('grid' - 1) == UR - 2 * ('grid' - 1)
LL := LL - 1 * ('grid' - 1) == UR - 3 * ('grid' - 1)

The phrasing is a bit vague / obtuse.  Do both diagonals individually need to fail the 10% test, or is it the count of primes along both diagonals as a set that falls beneath 10%?  Probably the latter.

PrimesSeen * 10 < 1 + 'grid' * 4

Also, this seems to be wheel factorization applied.

~

This runs _terriably_ with wheel factorization and maybe almost as bad with the Pollar derived 'make SURE it factors' auto-limited test (based on my experiments up to 1 million (~500K tested primes))

I let it run all day while I was at a family event and it's still going.  It's so bad I want to collect a record of just how bad.

About TWO DAYS bad

Side: 26177     685235329       5240/52353      false
Side: 26241     688590081       5248/52481      true
Euler 58: Spiral Primes: 26241

real    2927m11.858s
user    2915m37.092s
sys     6m51.301s

Compared to... less than a second with a 2020+ ProbablyPrime library function from math.big (uses https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test )

Euler 58: Spiral Primes: 26241

real    0m0.426s
user    0m0.473s
sys     0m0.063s



https://en.wikipedia.org/wiki/Primality_test#Miller%E2%80%93Rabin_and_Solovay%E2%80%93Strassen_primality_test

	// Given an integer n, choose some positive integer a < n.
	// FIRST: find: 2^(s)*d = n − 1 ; where d is odd.
	// If BOTH
	// a^d ≢ ± 1 ( mod n )
	// AND
	// a^(d*2^r) ≢ − 1 ( mod n )
	// {\displaystyle a^{2^{r}d}\not \equiv -1{\pmod {n}}}
	// FOR ALL
	// 0 ≤ r ≤ s − 1
	// THEN (if it's true) n is a witness that N is composite
	// ELSE N might or might not be prime

I want to better understand how this test works, and gain confidence that I am reading the formula correctly.

n for 25 (not prime), and 23 (prime)

	// FIRST: find: 2^(s)*d = n − 1 ; where d is odd.
	// 2 ^ s  * d = (n - 1)
	// D must be an odd Integer (whole number), 2 ^ s clearly <= (n - 1)
	// d = (n - 1) / 2 ^ s
24 / 2 = 12
24 / 4 = 6
24 / 8 = 3 == ODD
d = 3 ; s = 3
a = 2 (just picked a num; however the BPSW always picks 2, to test in Base 2.  Base 1 and base n-1  always think N is pseudo-prime)

	// a^d ≢ ± 1 ( mod n )
	// If a is large a ^ d may be greater than N
( 2 ^ 3 (+ / -) 1 ) % n != N (trivial here as 8 => 7,9 which both != N)
	// 0 ≤ r ≤ s − 1 // Re-2-flate?
2^(3*2^[0,1,2]) => 2^3 2^6 2^12 => 8 64 4096 => fail, fail, fail (which is pass)
25 is composite with 2 as witness


How did Wikipedia's code algorithm differ from the math?

Let s, d ... That's all the finding D that's trivial and makes sense.

Repeat?  How much confidence (K) based on random tries

x finds the first equation and ignores the -1 part

repeat S times takes that X (less the = 1 modN) == a^d and squares it, mod N every step ... which is the power communicative version of the second equation.  Though I didn't know the modulus on the other side could just be discarded as well.




How does the 'better' test use this?

https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test

A is the base?  Though the BPSW test only uses MR with Base 2 mode anyway.

== 1 ==
They 'wheel factorize filter' for some small list of N

https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Testing_against_small_sets_of_bases
MR can use rounds of 'base' prime up to 37 for extremely high confidence in numbers < 2^64

With 41 showing extremely high confidence.  2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41

== 2 ==

MR with base 2

== 3 ==

NOTE: Fails when Num is a perfect square, past a limit perform such a check.

WTF is a https://en.wikipedia.org/wiki/Jacobi_symbol

D = Primes: start at 5, every other prime is negative.
n == the Number
Find (first): -1 = Jacobi(D/num)

I haven't done math that high level in years, I could probably learn it again but how often is this used in most jobs? (Semi-serious on that, I don't expect to develop new crypto or compression algorithms anywhere near the first few years of work in the field!)

https://en.wikipedia.org/wiki/Jacobi_symbol#Implementation_in_Lua

Based on this:  P = 1 ; Q = (1 − D) / 4

== 4 == (draw the owl)
Perform a STRONG https://en.wikipedia.org/wiki/Lucas_pseudoprime#Strong_Lucas_pseudoprimes



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

func Euler058() uint64 {
	var ii, iim, sq, pseen, mulseen, ptarget, multarget uint64
	pseen, mulseen, ptarget, multarget = 0, 10, 1, 1
	// Gear 1 for up to...
	for ii = 3; sq < 4_000; ii += 2 {
		iim = ii - 1
		sq = ii * ii
		// if 0 == ii&0b11_1110 {
		//	fmt.Printf("Side: %d\t%d\t", ii, sq)
		// }
		euler.Primes.Grow(sq)
		// This is a square, obviously it's not prime
		for qq := 0; qq < 3; qq++ {
			sq -= iim
			if euler.Primes.KnownPrime(sq) {
				pseen++
			}
		}
		ptarget += 4
		// if 0 == ii&0b11_1110 {
		//	fmt.Printf("%d/%d\t%t\n", pseen, ptarget, pseen*mulseen < ptarget*multarget)
		// }
		if pseen*mulseen < ptarget*multarget {
			break
		}
	}
	// Gear 2
	for {
		iim = ii - 1
		sq = ii * ii
		// if 0 == ii&0b11_1110 {
		//	fmt.Printf("Side: %d\t%d\t", ii, sq)
		// }
		// euler.Primes.Grow(uint(sq))
		// This is a square, obviously it's not prime
		for qq := 0; qq < 3; qq++ {
			sq -= iim
			// This just makes a big.Int and runs the library function; which is why Golang is the correct tool for this job.
			if euler.ProbablyPrimeI64(int64(sq), 8) {
				pseen++
			}
		}
		ptarget += 4
		// if 0 == ii&0b11_1110 {
		//	fmt.Printf("%d/%d\t%t\n", pseen, ptarget, pseen*mulseen < ptarget*multarget)
		// }
		if pseen*mulseen < ptarget*multarget {
			break
		}
		ii += 2
	}
	return ii
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 58 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 58: Spiral Primes: 26241

/
*/
func main() {
	//test

	//run
	fmt.Printf("Euler 58: Spiral Primes: %d\n", Euler058())
}
