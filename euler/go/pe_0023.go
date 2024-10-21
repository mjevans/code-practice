// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=23
https://projecteuler.net/minimal=23

<p>A perfect number is a number for which the sum of its proper divisors is exactly equal to the number. For example, the sum of the proper divisors of $28$ would be $1 + 2 + 4 + 7 + 14 = 28$, which means that $28$ is a perfect number.</p>
<p>A number $n$ is called deficient if the sum of its proper divisors is less than $n$ and it is called abundant if this sum exceeds $n$.</p>

<p>As $12$ is the smallest abundant number, $1 + 2 + 3 + 4 + 6 = 16$, the smallest number that can be written as the sum of two abundant numbers is $24$. By mathematical analysis, it can be shown that all integers greater than $28123$ can be written as the sum of two abundant numbers. However, this upper limit cannot be reduced any further by analysis even though it is known that the greatest number that cannot be expressed as the sum of two abundant numbers is less than this limit.</p>
<p>Find the sum of all the positive integers which cannot be written as the sum of two abundant numbers.</p>

*/
/*

Another Primes revisit, I wasn't happy with the runtime, so I wanted to see if it improved with the updated Primes / Factorization classes.

Just one run, but...
real    0m12.928s
user    0m12.963s
sys     0m0.080s



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
	euler.Primes.Grow(28123)
	//sum := uint64(1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10 + 11) //constant+doc
	var sum uint64
	abund := []uint16{12}
	fmt.Printf("Euler023 ... Scan for Abundant Numbers...\n")
	for ii := uint16(13); ii < 28123; ii++ {
		iisum := uint16(euler.ListSumUint64(*(euler.Primes.Factorize(uint64(ii)).ProperDivisors())))
		if iisum > ii {
			abund = append(abund, ii)
			// fmt.Println(ii, "\tabundant")
		}
	}
	fmt.Printf("Euler023 ... Scan for combos...\n")
	aLim := len(abund)
Euler023outer:
	for ii := uint16(0); ii < 28123; ii++ {
		for x := 0; x < aLim && ii > abund[x]; x++ {
			xTarget := ii - abund[x]
			for y := 0; y < len(abund) && xTarget >= abund[y]; y++ {
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
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 23 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler023 ... Scan for Abundant Numbers...
Euler023 ... Scan for combos...
4179871         the sum of all the positive integers which cannot be written as the sum of two abundant numbers.

real    0m7.858s
user    0m7.900s
sys     0m0.071s

4179859         the sum of all the positive integers which cannot be written as the sum of two abundant numbers.
*/
func main() {
	// fmt.Println(grid)
	//test

	//run
	fmt.Println(Euler023(), "\tthe sum of all the positive integers which cannot be written as the sum of two abundant numbers.")
}
