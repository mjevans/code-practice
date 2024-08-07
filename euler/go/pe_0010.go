// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=10
<p>The sum of the primes below $10$ is $2 + 3 + 5 + 7 = 17$.</p>
<p>Find the sum of all the primes below two million.</p>



*/

import (
	"euler"
	"fmt"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler010(lim int) int {
	var primes []int
	for primes = euler.GetPrimes(nil, 0); lim > primes[len(primes)-1]; primes = euler.GetPrimes(primes, 0) {
	}
	ii := len(primes) - 1
	for {
		if lim < primes[ii-1] {
			ii--
		} else {
			break
		}
	}
	// fmt.Println(ii, primes[:ii])
	return euler.ListSum(primes[:ii])
}

func main() {
	//test
	fmt.Println("Euler010:\t", Euler010(10), Euler010(10) == 17)
	//run
	fmt.Println("Euler010:\t", Euler010(2000000))
}
