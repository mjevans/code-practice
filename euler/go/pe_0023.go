// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=23
<p>A perfect number is a number for which the sum of its proper divisors is exactly equal to the number. For example, the sum of the proper divisors of $28$ would be $1 + 2 + 4 + 7 + 14 = 28$, which means that $28$ is a perfect number.</p>
<p>A number $n$ is called deficient if the sum of its proper divisors is less than $n$ and it is called abundant if this sum exceeds $n$.</p>

<p>As $12$ is the smallest abundant number, $1 + 2 + 3 + 4 + 6 = 16$, the smallest number that can be written as the sum of two abundant numbers is $24$. By mathematical analysis, it can be shown that all integers greater than $28123$ can be written as the sum of two abundant numbers. However, this upper limit cannot be reduced any further by analysis even though it is known that the greatest number that cannot be expressed as the sum of two abundant numbers is less than this limit.</p>
<p>Find the sum of all the positive integers which cannot be written as the sum of two abundant numbers.</p>





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

func Euler023() uint64 {
	primes := euler.GetPrimes(nil, 512-8)
	sum := uint64(1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10 + 11) //constant+doc
	abund := []uint16{12}
Euler023outer:
	for ii := uint16(13); ii < 28123; ii++ {
		iisum := uint16(euler.ListSum(*(euler.FactorsToProperDivisors(euler.Factor(primes, int(ii))))))
		if iisum > ii {
			abund = append(abund, ii)
			// fmt.Println(ii, "\tabundant")
		}
		for x := 0; x < len(abund); x++ {
			for y := 0; y < len(abund); y++ {
				if ii == abund[x]+abund[y] {
					// fmt.Println(ii, "\t2a")
					continue Euler023outer
				}
			}
		}
		sum += uint64(ii)
	}
	return sum
}

/*
4179859         the sum of all the positive integers which cannot be written as the sum of two abundant numbers.
 */
func main() {
	// fmt.Println(grid)
	//test

	//run
	fmt.Println(Euler023(), "\tthe sum of all the positive integers which cannot be written as the sum of two abundant numbers.")
}
