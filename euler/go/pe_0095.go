// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=95
https://projecteuler.net/minimal=95

<p>The proper divisors of a number are all the divisors excluding the number itself. For example, the proper divisors of $28$ are $1$, $2$, $4$, $7$, and $14$. As the sum of these divisors is equal to $28$, we call it a perfect number.</p>
<p>Interestingly the sum of the proper divisors of $220$ is $284$ and the sum of the proper divisors of $284$ is $220$, forming a chain of two numbers. For this reason, $220$ and $284$ are called an amicable pair.</p>
<p>Perhaps less well known are longer chains. For example, starting with $12496$, we form a chain of five numbers:
$$12496 \to 14288 \to 15472 \to 14536 \to 14264 (\to 12496 \to \cdots)$$</p>
<p>Since this chain returns to its starting point, it is called an amicable chain.</p>
<p>Find the smallest member of the longest amicable chain with no element exceeding one million.</p>

/
*/
/*
	Between a refactored ProperDivisors(Sum) that uses Dynamic Programming and the Cache for Factorize the speed improved by almost 2/3rds of the prior execution time; but still seems slow.
	Caching the chain length rather than just if seen helped to correctly identify the looped parts of the chain, rather than precursors that might lead there.
/
*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0095(max, maxChain uint32) (uint32, uint32) {
	chain := make([]uint32, maxChain)
	cl := make([]uint8, 1+max)
	var blMin, blLen, curMin, curLen, ii, cur, jj uint32
	for ii = 2; ii <= max; ii++ {
		if 0 == ii&0xFFFF {
			fmt.Println(ii)
		}
		if 0 != cl[ii] {
			continue // If already seen; next
		}
		chain[0], cur, curLen = ii, ii, 1
		for {
			cur = uint32(euler.Primes.Factorize(uint64(cur)).ProperDivisorsSum())
			if max < cur {
				for 0 < curLen {
					curLen--
					cl[chain[curLen]] = 0xFF // -1 -- rejected value
				}
				curLen = 0xFF // -1 // Overwrite for the outer-loop
				break         // 1 // Abort: exceeded allowed range
			}
			if 0 != cl[cur] {
				jj, curLen = curLen, uint32(cl[cur])
				for 0 < jj {
					jj--
					cl[chain[jj]] = uint8(curLen)
				}
				break // 1 // Found a new entry to a known Loop
			}
			for jj = curLen - 1; 0 < jj && cur != chain[jj]; jj-- {
			}
			if cur == chain[jj] {
				chainLen := uint8(curLen - jj)
				curMin = cur
				for 0 < curLen {
					curLen--
					cl[chain[curLen]] = chainLen
					if curLen > jj && curMin > chain[curLen] {
						curMin = chain[curLen]
					}
				}
				curLen = uint32(chainLen)
				if (curLen > blLen) || (curLen == blLen && blMin > curMin) {
					fmt.Printf("Found new best length (from %d) (%d) %d and min (%d) %d \n", ii, blLen, curLen, blMin, curMin)
					blMin, blLen = curMin, curLen
				}
				break
			}
			chain[curLen] = cur
			curLen++
		}
		cl[ii] = uint8(curLen)
	}

	return blMin, blLen
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 95 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Found new best length (from 2) (0) 2 and min (0) 0
Found new best length (from 9464) (2) 5 and min (0) 12496
Found new best length (from 2) (0) 2 and min (0) 0
Found new best length (from 5916) (2) 28 and min (0) 14316
65536
131072
196608
262144
327680
393216
458752
524288
589824
655360
720896
786432
851968
917504
983040
Euler 95: Amicable Chains: (28) 14316

real    0m12.098s
user    0m12.123s
sys     0m0.102s
.
*/
func main() {
	var mn, ln uint32
	euler.Primes.FactSetCache(1_000_000) // This is a LITTLE faster
	//test
	mn, ln = Euler0095(15472, 60)
	if 12496 != mn || 5 != ln {
		panic(fmt.Sprintf("Euler 95: Test Case failed with values: (%d) %d\n", ln, mn))
	}

	//run
	mn, ln = Euler0095(1_000_000, 90)
	fmt.Printf("Euler 95: Amicable Chains: (%d) %d\n", ln, mn)
	if 14316 != mn {
		panic("Did not reach expected value.")
	}
}
