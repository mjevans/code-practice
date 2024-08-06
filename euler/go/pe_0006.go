// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=6
<p>The sum of the squares of the first ten natural numbers is,</p>
$$1^2 + 2^2 + ... + 10^2 = 385.$$
<p>The square of the sum of the first ten natural numbers is,</p>
$$(1 + 2 + ... + 10)^2 = 55^2 = 3025.$$
<p>Hence the difference between the sum of the squares of the first ten natural numbers and the square of the sum is $3025 - 385 = 2640$.</p>
<p>Find the difference between the sum of the squares of the first one hundred natural numbers and the square of the sum.</p>


*/

import (
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
	answer := Euler006(1, 10)
	fmt.Println(answer[0] == 2640, answer[1] == 385, answer[2] == 3025, ": Tests for 1..10 passed?", answer)
	// Q
	answer = Euler006(1, 100)
	fmt.Println("Euler006 for 1..100:\t", answer)

}
