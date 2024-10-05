// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=21
https://projecteuler.net/minimal=21

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
	for ii := start; ii <= end; ii++ {
		var iis, iie int
		var ok bool
		if iis, ok = cache[ii]; !ok {
			iis = int(euler.ListSumUint64(*(euler.Primes.Factorize(uint(ii)).ProperDivisors())))
			cache[ii] = iis
		} else {
			continue
		}
		if iie, ok = cache[iis]; !ok {
			iie = int(euler.ListSumUint64(*(euler.Primes.Factorize(uint(iis)).ProperDivisors())))
			cache[iis] = iie
		}
		// fmt.Println("Loop", ii, iis, iie)
		if ii != iis && ii == iie {
			if iis > end || iie > end {
				// fmt.Println("OVEREND", iis)
				oflow += int64(iis) + int64(ii)
				continue
			}
			fmt.Println("Adding ", ii, " (", iis, ")")
			ret += int64(ii) + int64(iis)
		}
	}
	fmt.Println("Euler021", ret, "with Overflow:", ret+oflow)
	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 21 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler021 divisors of 220 :  [1 2 4 5 10 11 20 22 44 55 110]  sum to  284 true
Euler021 divisors of 284 :  [1 2 4 71 142]  sum to  220 true
Adding  220  ( 284 )
Adding  1184  ( 1210 )
Adding  2620  ( 2924 )
Adding  5020  ( 5564 )
Adding  6232  ( 6368 )
Euler021 47938 with Overflow: 47938
Euler021  47938
*/
func main() {
	// fmt.Println(grid)
	//test
	dv := euler.FactorsToProperDivisors(euler.Factor(nil, 220))
	fmt.Println("Euler021 divisors of 220 : ", *dv, " sum to ", 284, 284 == euler.ListSum(*dv))
	dv = euler.FactorsToProperDivisors(euler.Factor(nil, 284))
	fmt.Println("Euler021 divisors of 284 : ", *dv, " sum to ", 220, 220 == euler.ListSum(*dv))

	//run
	fmt.Println("Euler021 ", Euler021(1, 10000))
	fmt.Printf("\t\t\t*** TODO correct answer before new Euler problems ***\n")
}
