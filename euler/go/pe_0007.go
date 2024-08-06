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

func Euler006(start, end int) [3]int {
	var sumsq, sum int
	for ; start <= end; start++ {
		sum += start
		sumsq += start * start
	}
	sum *= sum
	return [3]int{sum - sumsq, sumsq, sum}
}

func main() {
	// Tests
	fmt.Println(euler.Factor(nil, 200))
	answer := Euler006(1, 10)
	fmt.Println(answer[0] == 2640, answer[1] == 385, answer[2] == 3025, ": Tests for 1..10 passed?", answer)
	// Q
	answer = Euler006(1, 100)
	fmt.Println("Euler006 for 1..100:\t", answer)

}
