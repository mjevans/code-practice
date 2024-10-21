// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=37
https://projecteuler.net/minimal=37

<p>The number $3797$ has an interesting property. Being prime itself, it is possible to continuously remove digits from left to right, and remain prime at each stage: $3797$, $797$, $97$, and $7$. Similarly we can work from right to left: $3797$, $379$, $37$, and $3$.</p>
<p>Find the sum of the only eleven primes that are both truncatable from left to right and right to left.</p>
<p class="smaller">NOTE: $2$, $3$, $5$, and $7$ are not considered to be truncatable primes.</p>



*/
/*





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

func Euler037(limit uint64) uint64 {
	if 0 == limit {
		limit = 1_000_000
	}
	p := euler.NewBVPrimes()
	p.Grow(uint64(limit))
	var sum uint64
	found := uint8(0)  // Magic numbers, the problem said there were 11 to find
	prime := uint64(7) // the first prime after 7 is a 2 digit number
Euler037Outer:
	for found < 11 && prime < uint64(limit) {
		prime = p.PrimeAfter(prime)
		// From the right, and also track how big it is.
		digits := uint64(1)
		for test := prime / 10; 0 < test; test /= 10 {
			if false == p.KnownPrime(test) {
				continue Euler037Outer
			}
			digits *= 10
		}
		for test := prime % digits; 0 < test; test %= digits {
			if false == p.KnownPrime(test) {
				continue Euler037Outer
			}
			digits /= 10
		}
		sum += prime
		found++
		fmt.Printf("Euler037: match %d\t%d\t\t%d\n", found, prime, sum)
	}

	return uint64(sum)
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 37 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler037: match 1       23              23
Euler037: match 2       37              60
Euler037: match 3       53              113
Euler037: match 4       73              186
Euler037: match 5       313             499
Euler037: match 6       317             816
Euler037: match 7       373             1189
Euler037: match 8       797             1986
Euler037: match 9       3137            5123
Euler037: match 10      3797            8920
Euler037: match 11      739397          748317
Euler037: Sum of truncatable primes : 1_000_000 ==      748317



*/
func main() {
	//test

	//run
	sum := Euler037(1_000_000)
	fmt.Printf("Euler037: Sum of truncatable primes : 1_000_000 ==\t%d\n", sum)

}
