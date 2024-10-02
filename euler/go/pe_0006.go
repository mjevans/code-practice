// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=6
https://projecteuler.net/minimal=6

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
