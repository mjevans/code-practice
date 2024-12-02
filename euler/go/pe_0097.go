// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=97
https://projecteuler.net/minimal=97

<p>The first known prime found to exceed one million digits was discovered in 1999, and is a Mersenne prime of the form $2^{6972593} - 1$; it contains exactly $2\,098\,960$ digits. Subsequently other Mersenne primes, of the form $2^p - 1$, have been found which contain more digits.</p>
<p>However, in 2004 there was found a massive non-Mersenne prime which contains $2\,357\,207$ digits: $28433 \times 2^{7830457} + 1$.</p>
<p>Find the last ten digits of this prime number.</p>


/
*/
/*
	My first, half a moment thought "Really really big numbers won't normally come up in a real job; when they do the proper library is the correct answer."
	However I can do 'better'.
	Literally, just the last 10 (decimal) digits.  That's approximately 34 bits of information.
	Given it's powers of 2, the mul-mod is a dead-easy special case; I just have to keep multiplying it and have part trimmed far enough to the left that it doesn't enter the required precision area.
	28433 * 2^(7830457) + 1 ... wait, where's that plus one happen?
	28433 * ( 2^(7830457) + 1 ) << I'm going to try this form first given the raw text above.
	7830457 / 30 ~= 261015.23
	Bzzt the format was wrong:
	(28433 * 2^(7830457)) + 1  << correct number
/
*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0097() uint64 {
	var b2, p2 uint64
	const bits = 24
	b2, p2 = 1, 7830457
	for {
		if bits < p2 {
			p2, b2 = p2-bits, (b2<<bits)%10_000_000_000
		} else {
			b2 = (b2 << p2) % 10_000_000_000
			break
		}
	}
	return ((28433 * b2) + 1) % 10_000_000_000
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 97 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 97: Large Non-Mersenne Prime (2004): 8739992577

real    0m0.099s
user    0m0.136s
sys     0m0.059s
.
*/
func main() {
	var r uint64
	//test

	//run
	r = Euler0097()
	fmt.Printf("Euler 97: Large Non-Mersenne Prime (2004): %d\n", r)
	if 8_739_992_577 != r {
		panic("Did not reach expected value.")
	}
}
