// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=50
https://projecteuler.net/minimal=50

<p>The prime $41$, can be written as the sum of six consecutive primes:</p>
$$41 = 2 + 3 + 5 + 7 + 11 + 13.$$
<p>This is the longest sum of consecutive primes that adds to a prime below one-hundred.</p>
<p>The longest sum of consecutive primes below one-thousand that adds to a prime, contains $21$ terms, and is equal to $953$.</p>
<p>Which prime, below one-million, can be written as the sum of the most consecutive primes?</p>


*/
/*

Find runs of consecutive prime numbers that add to a prime number, less than a limit

P returns the _gretest chain_

P(100) = 41, 6
P(1000) = 953, 21
P(1_000_000) = ?, ?

Sliding window?  I thought about trying to transform the state, but all the previous tests would need to be re-run without the first number...

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

func Euler050(limit uint) (uint, uint) {
	euler.Primes.Grow(limit)
	var start, next, sum, run, bestsum, bestrun uint
	for {
		start = euler.Primes.PrimeAfter(start)
		run = 1
		sum = start
		next = start
		if start*bestrun > limit {
			break
		}
		for {
			next = euler.Primes.PrimeAfter(next)
			sum += next
			if sum > limit {
				break
			}
			run++
			if run > bestrun && euler.Primes.KnownPrime(sum) {
				// fmt.Printf("New Best Run: %d > %d (%d ?? %d)\n", run, bestrun, sum, bestsum)
				bestrun, bestsum = run, sum
			}
		}
	}
	return bestsum, bestrun
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 50 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
Euler 50: test: 100 => sum: 41  run: 6  true
Euler 50: test: 1000 => sum: 953        run: 21 true
Euler 50: Consecutive Prime Sum: 1000000        run: 543        sum: 997651


*/
func main() {
	//test
	n := uint(100)
	sum, run := Euler050(n)
	fmt.Printf("Euler 50: test: %d => sum: %d\trun: %d\t%t\n", n, sum, run, 41 == sum && run == 6)
	n = 1000
	sum, run = Euler050(n)
	fmt.Printf("Euler 50: test: %d => sum: %d\trun: %d\t%t\n", n, sum, run, 953 == sum && run == 21)

	//run
	n = 1_000_000
	sum, run = Euler050(n)
	fmt.Printf("Euler 50: Consecutive Prime Sum: %d\trun: %d\tsum: %d\n", n, run, sum)
}
