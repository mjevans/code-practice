// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=21
<p>Let $d(n)$ be defined as the sum of proper divisors of $n$ (numbers less than $n$ which divide evenly into $n$).<br>
If $d(a) = b$ and $d(b) = a$, where $a \ne b$, then $a$ and $b$ are an amicable pair and each of $a$ and $b$ are called amicable numbers.</p>
<p>For example, the proper divisors of $220$ are $1, 2, 4, 5, 10, 11, 20, 22, 44, 55$ and $110$; therefore $d(220) = 284$. The proper divisors of $284$ are $1, 2, 4, 71$ and $142$; so $d(284) = 220$.</p>
<p>Evaluate the sum of all the amicable numbers under $10000$.</p>




*/

import (
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

func Euler021(start, end int) int64 {
	var ret, oflow int64
	primes := euler.GetPrimes(nil, 256-8)
	for ii := start; ii <= end; ii++ {
		iis := euler.ListSum(euler.FactorsToDivisors(euler.Factor(primes, ii)))
		if iis > end {
			fmt.Println("OVEREND", iis)
			oflow += int64(iis)
		}
		if ii == euler.ListSum(euler.FactorsToDivisors(euler.Factor(primes, iis))) {
			fmt.Println("Adding ", ii, " (", iis, ")")
			ret += int64(ii)
		}
	}
	fmt.Println("Euler021", ret, "with Overflow:", ret+oflow)
	return ret
}

/*
Euler021 40284 with Overflow: 8155829
Euler021  40284

 */
func main() {
	// fmt.Println(grid)
	//test
	dv := euler.FactorsToDivisors(euler.Factor(nil, 220))
	fmt.Println("Euler021 divisors of 220 : ", dv, " sum to ", 284, 284 == euler.ListSum(dv))
	dv = euler.FactorsToDivisors(euler.Factor(nil, 284))
	fmt.Println("Euler021 divisors of 284 : ", dv, " sum to ", 220, 220 == euler.ListSum(dv))

	//run
	fmt.Println("Euler021 ", Euler021(1, 10000))
}
