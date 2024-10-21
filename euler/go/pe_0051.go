// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=51
https://projecteuler.net/minimal=51

<p>By replacing the 1<sup>st</sup> digit of the 2-digit number *3, it turns out that six of the nine possible values: 13, 23, 43, 53, 73, and 83, are all prime.</p>
<p>By replacing the 3<sup>rd</sup> and 4<sup>th</sup> digits of 56**3 with the same digit, this 5-digit number is the first example having seven primes among the ten generated numbers, yielding the family: 56003, 56113, 56333, 56443, 56663, 56773, and 56993. Consequently 56003, being the first member of this family, is the smallest prime with this property.</p>
<p>Find the smallest prime which, by replacing part of the number (not necessarily adjacent digits) with the same digit, is part of an eight prime value family.</p>


*/
/*

Examples:

_3 -> 13 23 43 53 73 83 == 6 primes wait 9 possible values?  They exclude zeros? (this would be 7 if '03' were included)
56__3 -> 56003 56113 56333 56443 56663 56773 56993 == 7 primes, this one includes zero.


I considered two main tactics:
1) scan all the primes and try to bucket them based on runs of digits - this... seemed way too expensive for memory, CPU, and complexity.
2) 'generative' and reject numbers.

If I think about the inconsistency in the examples, they stated with 13 not 03.

If the 1s place is replaced, then 50% can't be primes due to even, so exclude it from consideration entirely (less than 80% minimum == 8 / 10 combos).

Q1: What to to if there's more than one digit that repeats in the base prime?

* Reverse-filter, if 0 and 1 or 2 aren't the base-repeat digit then it isn't possible for that to be the prime
* 3 strikes and it's out

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

func Euler051(limit, minrun, base uint64) (uint64, uint64) {
	euler.Primes.Grow(limit)
	var prime, run, bestprime, bestrun, strikes, strikeOut uint64
	strikeOut = base + 1 - minrun
	// prime = 56001 // FIXME 0
	prime = 0
	for prime <= limit {
		prime = euler.Primes.PrimeAfter(prime)
		nv := euler.Uint64ToDigitsUint8(uint64(prime), uint64(base))
		ddmx := uint8(len(nv))
		nvtest := make([]uint8, ddmx)
		// skip the base digit place, given the 7s result for even _3+ as the case and even numbers == strikes
		// fmt.Printf("%d: Testing %v\n", prime, nv)
		for dd := uint8(1); dd < ddmx; dd++ {
			ddnum := nv[dd]
			if strikeOut < uint64(ddnum) {
				// fmt.Printf("%d: digit %d (%d) above %d, SKIP\n", prime, ddnum, dd, strikeOut)
				continue
			}
			ddx := make([]uint8, 0, 4)
			// always include ddnum, 0 or more numbers from ddx will also rotate, tests without ddnum occur in later for dd loops
			// ddx = append(ddx, dd)
			for mm := dd + 1; mm < ddmx; mm++ {
				if ddnum == nv[mm] {
					ddx = append(ddx, mm)
				}
			}
			ddxMx := uint64(len(ddx))
			perms := (uint64(1) << ddxMx) - 1
			// fmt.Printf("%d: digit %d (%d) + digits %v\tperms? %d\n", prime, ddnum, dd, ddx, perms)
			for pp := uint64(0); pp <= perms; pp++ {
				// clean template copy
				copy(nvtest, nv)
				run = 1
				strikes = uint64(ddnum) + 1 - 1
				for mm := ddnum + 1; strikes <= strikeOut && mm < uint8(base); mm++ {
					// Construct this cycle's permutation
					nvtest[dd] = mm
					tbit := uint64(1)
					tnum := tbit
					for tnum <= ddxMx {
						if 0 < pp&tbit {
							nvtest[ddx[tnum-1]] = mm
						}
						tbit <<= 1
						tnum++
					}
					numtest := uint64(euler.Uint8DigitsToUint64(nvtest, uint64(base)))
					numtested := euler.Primes.KnownPrime(numtest)
					// fmt.Printf("Tested %d +%d+ %d = %t\n", prime, run, numtest, numtested)

					// Is it prime?
					if numtested {
						run++
					} else {
						strikes++
					}
				}
				if strikes >= strikeOut {
					// fmt.Printf("%d: perm %d+ %d Strike Out\n", prime, dd, pp)
					continue
				}
				if run > bestrun {
					fmt.Printf("%d: perm %d+ %d New Best Run %d < %d : %d\n", prime, dd, pp, bestrun, run, prime)
					bestprime = prime
					bestrun = run
				}
			}
		}
		if bestrun >= minrun {
			fmt.Printf("%d: Tartet Best Run Reached: %d <= %d : %d\n", prime, minrun, bestrun, prime)
			return bestprime, bestrun
		}
	}
	return bestprime, bestrun
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 51 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

13: perm 1+ 0 New Best Run 0 < 6 : 13
13: Tartet Best Run Reached: 6 <= 6 : 13
Euler 51: test: prime: 13       run: 6  true
56003: perm 1+ 1 New Best Run 0 < 7 : 56003
56003: Tartet Best Run Reached: 7 <= 7 : 56003
Euler 51: test: prime: 56003    run: 7  true
121313: perm 1+ 3 New Best Run 0 < 8 : 121313
121313: Tartet Best Run Reached: 8 <= 8 : 121313
Euler 51: Prime Digit Replacements: 1000000     run: 8  prime: 121313


*/
func main() {
	//test
	var n, t uint64
	n, t = 100, 6
	prime, run := Euler051(n, t, 10)
	fmt.Printf("Euler 51: test: prime: %d\trun: %d\t%t\n", prime, run, 13 == prime && run == t)
	n, t = 100000, 7
	prime, run = Euler051(n, t, 10)
	fmt.Printf("Euler 51: test: prime: %d\trun: %d\t%t\n", prime, run, 56003 == prime && run == t)

	//run
	n, t = 1_000_000, 8
	prime, run = Euler051(n, t, 10)
	fmt.Printf("Euler 51: Prime Digit Replacements: %d\trun: %d\tprime: %d\n", n, run, prime)
}
