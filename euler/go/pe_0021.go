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
	"math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler021(start, end int) int64 {
	var ret, oflow int64
	cache := make(map[int]int, int(math.Log2(float64(end)))+1)
	for bb := 64; bb <= end; bb *= 2 {
		ret += int64(bb - 1)
		cache[bb] = bb - 1
	}
	// fmt.Println(cache)
	primes := euler.GetPrimes(nil, 256-8)
	for ii := start; ii <= end; ii++ {
		var iis, iie int
		var ok bool
		if iis, ok = cache[ii]; !ok {
			iis = euler.ListSum(euler.FactorsToDivisors(euler.Factor(primes, ii)))
		} else {
		}
		if iis > end {
			// fmt.Println("OVEREND", iis)
			oflow += int64(iis)
		}
		if iie, ok = cache[iis]; !ok {
			iie = euler.ListSum(euler.FactorsToDivisors(euler.Factor(primes, iis)))
		} else {
		}
		// fmt.Println("Loop", ii, iis, iie)
		if ii == iie {
			fmt.Println("Adding ", ii, " (", iis, ")")
			ret += int64(ii)
		}
	}
	fmt.Println("Euler021", ret, "with Overflow:", ret+oflow)
	return ret
}

/*
... 8192?
Euler021 40284 with Overflow: 8155829
Euler021  40284

Euler021 divisors of 220 :  [1 2 4 5 10 11 20 22 44 55 110]  sum to  284 true
Euler021 divisors of 284 :  [1 2 4 71 142]  sum to  220 true
Adding  6  ( 6 )
Adding  28  ( 28 )
Adding  220  ( 284 )
Adding  284  ( 220 )
Adding  496  ( 496 )
Adding  1184  ( 1210 )
Adding  1210  ( 1184 )
Adding  2620  ( 2924 )
Adding  2924  ( 2620 )
Adding  5020  ( 5564 )
Adding  5564  ( 5020 )
Adding  6232  ( 6368 )
Adding  6368  ( 6232 )
Adding  8128  ( 8128 )
Euler021 56596 with Overflow: 8172141
Euler021  56596
$ factor 8127
8127: 3 3 3 7 43
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
