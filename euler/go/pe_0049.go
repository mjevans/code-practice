// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=49
https://projecteuler.net/minimal=49

<p>The arithmetic sequence, $1487, 4817, 8147$, in which each of the terms increases by $3330$, is unusual in two ways: (i) each of the three terms are prime, and, (ii) each of the $4$-digit numbers are permutations of one another.</p>
<p>There are no arithmetic sequences made up of three $1$-, $2$-, or $3$-digit primes, exhibiting this property, but there is one other $4$-digit increasing sequence.</p>
<p>What $12$-digit number do you form by concatenating the three terms in this sequence?</p>



*/
/*

4 digit prime
permutations
'increasing sequence' (so permu is prime AND absolute diff is either 2x the span (halfway number is a perm and also prime) or 1x span and ...)
NOT [ 1487, 4817, 8147 ]

ANSWER logical concat 12 digit number (or just print gapless) in order from lowest to greatest.

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

func Euler049(start, limit uint64) [3]uint64 {
	p := euler.NewBVPrimes()
	p.Grow(limit)
	ii := start - 1
	for {
		ii = p.PrimeAfter(ii)
		if ii > limit {
			break
		}
		sld := euler.Uint64ToDigitsUint8(uint64(ii), 10)
		pl := euler.FactorialUint64(uint64(len(sld))) - 1
		comboWanted := make(map[uint64]uint64)
		for pp := uint64(0); pp < pl; pp++ {
			iipp := uint64(euler.Uint8DigitsToUint64(euler.PermutationSlUint8(pp, sld), 10))
			// if 2969 == ii {
			// fmt.Printf("Combo: %d [%d] => %d\n", ii, pp, iipp)
			// }
			if iipp < start {
				continue
			}
			if other, ok := comboWanted[iipp]; ok {
				// three number sort
				a, b, c := ii, iipp, other
				if b < a {
					a, b = b, a
				}
				if c < a {
					a, c = c, a
				}
				if c < b {
					b, c = c, b
				}
				return [3]uint64{a, b, c}
			}
			if false == p.KnownPrime(iipp) {
				continue // pp
			}
			a, b := ii, iipp
			if b < a {
				a, b = b, a
			}
			ppl, ppg := a+(b-a)>>1, a+(b-a)<<1
			if ppg != ii && ppg != iipp && ppg < limit && p.KnownPrime(ppg) {
				comboWanted[ppg] = iipp
			}
			if ppl != ii && ppl != iipp && p.KnownPrime(ppl) {
				comboWanted[ppl] = iipp
			}
		}
	}
	return [3]uint64{}
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 49 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 49: test: 148748178147
Euler 49: Prime Permutations: 296962999629

*/
func main() {
	//test
	t := Euler049(1000, 10000)
	fmt.Printf("Euler 49: test: %04d%04d%04d\n", t[0], t[1], t[2])

	//run
	a := Euler049(1489, 10000)
	fmt.Printf("Euler 49: Prime Permutations: %04d%04d%04d\n", a[0], a[1], a[2])
}
