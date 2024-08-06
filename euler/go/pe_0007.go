// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=7
<p>By listing the first six prime numbers: $2, 3, 5, 7, 11$, and $13$, we can see that the $6$th prime is $13$.</p>
<p>What is the $10\,001$st prime number?</p>



*/

import (
	"euler"
	"fmt"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "os" // os.Stdout
)

func main() {
	// Tests
	// fmt.Println(euler.Factor(nil, 200))
	primes := euler.GetPrimes(nil, 10001-8)
	fmt.Println(primes[6-1])
	fmt.Println(primes[10001-1])
}
