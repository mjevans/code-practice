// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=3
<p>The prime factors of $13195$ are $5, 7, 13$ and $29$.</p>
<p>What is the largest prime factor of the number $600851475143$?</p>

*/

import (
	"fmt"
	"strings"
	// "os" // os.Stdout
)

func Factor(primes []int, num int) []int {
	// Public school factoring algorithm from memory...

	// With a list of known primes, the largest number that can be factored is Pn * Pn
	for ; nil == primes || num > primes[len(primes)-1]*primes[len(primes)-1]; primes = GetPrimes(primes) {
		// fmt.Println(len(primes), primes[len(primes)-1])
	}

	ret := []int{}
	for _, prime := range primes {
		for ; 0 == num%prime; num /= prime {
			ret = append(ret, prime)
		}
		if num < prime*prime {
			break
		} // break if no more prime factors are possible
	}
	if 0 < num {
		ret = append(ret, num)
	}
	return ret
}

func GetPrimes(primes []int) []int {
	if nil == primes {
		return []int{2, 3, 5, 7, 11, 13, 17, 19}
	}
	// Semi-arbitrary expansion target, find 8 more primes (8, 16, 32, 64 it'll fit within the append growth algo)
PrimeHunt:
	for primehunt := 8; 0 < primehunt; primehunt-- {
		for ii := primes[len(primes)-1] + 1; ; ii++ {
			result := Factor(primes, ii)
			if 1 == len(result) && primes[len(primes)-1] < result[0] {
				//fmt.Println("Found Prime:\t", result[0])
				primes = append(primes, result[0])
				continue PrimeHunt // I could break once, but this documents the intent
			}
		}
	}
	return primes
}

func PrintFactors(factors []int) {
	// Join only takes []string s? fff
	strFact := make([]string, len(factors), len(factors))
	for ii, val := range factors {
		strFact[ii] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strFact, ", "))
}

func main() {
	// Test
	// PrintFactors(Factor(nil, 13195))
	// matches factor
	PrintFactors(Factor(nil, 600851475143))
}
