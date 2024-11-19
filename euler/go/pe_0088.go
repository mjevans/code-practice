// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=88
https://projecteuler.net/minimal=88

<p>A natural number, $N$, that can be written as the sum and product of a given set of at least two natural numbers, $\{a_1, a_2, \dots, a_k\}$ is called a product-sum number: $N = a_1 + a_2 + \cdots + a_k = a_1 \times a_2 \times \cdots \times a_k$.</p>
<p>For example, $6 = 1 + 2 + 3 = 1 \times 2 \times 3$.</p>
<p>For a given set of size, $k$, we shall call the smallest $N$ with this property a minimal product-sum number. The minimal product-sum numbers for sets of size, $k = 2, 3, 4, 5$, and $6$ are as follows.</p>
<ul style="list-style-type:none;">
<li>$k=2$: $4 = 2 \times 2 = 2 + 2$</li>
<li>$k=3$: $6 = 1 \times 2 \times 3 = 1 + 2 + 3$</li>
<li>$k=4$: $8 = 1 \times 1 \times 2 \times 4 = 1 + 1 + 2 + 4$</li>
<li>$k=5$: $8 = 1 \times 1 \times 2 \times 2 \times 2 = 1 + 1 + 2 + 2 + 2$</li><li>$k=6$: $12 = 1 \times 1 \times 1 \times 1 \times 2 \times 6 = 1 + 1 + 1 + 1 + 2 + 6$</li></ul>
<p>Hence for $2 \le k \le 6$, the sum of all the minimal product-sum numbers is $4+6+8+12 = 30$; note that $8$ is only counted once in the sum.</p>
<p>In fact, as the complete set of minimal product-sum numbers for $2 \le k \le 12$ is $\{4, 6, 8, 12, 15, 16\}$, the sum is $61$.</p>
<p>What is the sum of all the minimal product-sum numbers for $2 \le k \le 12000$?</p>

/
*/
/*
	I don't have to store all the terms out in memory like that...
	E.G. the type Factorized I wrote for working with prime power numbers can also work, though most of the special functions assume 'reduced' forms, not an 'improper fraction' like other numbers form.
	The core issue revolves around a balance point between:
	* All 1s 1x1x...x1 = 1 but sum() = k*1
	* Some multiplied term in a few of those slots which EQUALS than the sum of that number's terms and the remaining 1s...

	Clearly the floor cannot be under K, but finding sum/product N which is the smallest seems difficult, more so when factors can be combined in multitudes of ways (proper divisors)....
	Uggh, they're trying to blindside with the coin problem and a limited set of coins which change dynamically as some are utilized.

	The example should be studied more closely for inspiration on algorithms...
	2	1+1 = 2 ne, 1+2 = 3 not factors, 2*2 = 4 -- done
	3	2*2*2 = 8 -- too high.  2 is under sum, 1*2*2 = 4 is under sum... 2*2*2=8 over sum.
	4	1*1*2*2 = 4 vs 6, 1*2*2*2 = 8 sum 7 -- the answer was 1,1,2,4 = 8
	5	1*1*1*2*2 = 4 vs 7, 1*1*2*2*2 = 8 sum 8 -- done
	6	1*1*1*2*2*2 = 8 (9),1*1*2*2*2*2 = 16 (10) too high
	7	12 = 1^5*3*4 ~~ 1*1*1*1*2*2*2 = 8 (10) low, 1*1*1*2*2*2*2 = 16 (11) high, 11 prime, 13 prime, 12 is 2*2*3 (11) or 4*3 () or 2*6
	8	12 is 2*2*3*1^5 (12)		1^5*2^3 = 8 (11), 1^4*2^4 = 16 (12), 9 is 3*3 * 1^6 (12), 10 is 2*5 * 1^6 (13), 11 is prime, 12 is is 2*2*3*1^5 (12)
	9?  15?	3*5*1^7		1^5*2^4 = 8 (13) lower, 1^4*2^5 = 16 (14) higher scan between too low and too high
	12? 16?	16 = 1*1*1 * 1*1*1 * 1*1*2 * 2*2*2 = 16
	So I think 10 and 11 aren't answers
	10	1^7*2^3 = 8 (13), 1^6*2^4 = 16 (14), 9 is 3x3*1^8 (14), 10 is 2*5*1^8 (15), 11 is prime, 12 is is 2*2*3*1^7 (14) or 4*3*1^8 (15) or 2*6*1^8 (16), 13 is prime, 14 is 2*7*1^8 (17), 15 is 3*5*1^8 (16) exhausted
	11	1^8*2^3 = 8 (14), 1^7*2^4 = 16 (15), 9 is 3x3*1^9 (15), 10 is 2*5*1^9 (16), 11 is prime, 12 is is 2*2*3*1^8 (15) or 4*3*1^9 (16) or 2*6*1^9 (17), 13 is prime, 14 is 2*7*1^9 (18), 15 is 3*5*1^9 (17) exhausted

	So far it's reasonable to start with a power of two apart as a set of bounds, but that gets very wide very quick.

	What if I swapped two of the upper bounds 2s with a 1 and a 3?  That should be 75% of the number but the same total.

	11	1^8*2^3 = 8 (14), 1^7*2^4 = 16 (15), 1^8*3*2^2 (15) = 12 that is a better lower bound,
		How can I verify that though?  1^7*3*2^3 (16) = 24 higher than my previous high, not great so far... but this reminds me a bit of how square roots are approximated.
		1^8*3^2*2 (16) = 12
		I feel mostly confident that the '75%' substitution technique (2*2 -> 3*1) won't overlook better solutions since 2 and 3 are both prime, and overall it doesn't change the total, and thus N...
		However once other numbers are introduced to the system either the sum (1,4), result (?), or both (1,5) increase.
		There seems to be a special relationship between the multiplicative identity (1), it's doubling / the first prime number (2), and it's tripling / second prime number (3).  Further in the summation it's also special that (1,3) and (2,2) can swap to change the multiplication result without modifying the sum.

	The initial approach solved the trivial test cases quickly, and I was hopeful to estimate a narrow solution range.  About 10% of the way to the answer and 15 min into running it's clear that there's a LOT of wasted factorization effort, the same numbers tried with slightly different K term numbers.  Plus the estimate is at best 10-20% of the input number, rather than something better like always less than 20.

	I need to take an entirely different approach, but it's 2:30 am.

	I'm waffling on if I was on the wrong track or not, this might be one where I am better off allowing a brute force to run and looking over the discussion to see where my existing knowledge was lacking.

	A quick review of the known set to see if any patterns pop out that I missed:
	K	Answer
	2	4	2,2
	3	6	1,2,3
	4	8	1,1,2,4
	5	8*	1,1,2,2,2
	6	12	1,1,1,1,2,6
	8	12*	1,1,1,1,1,2,2,3	// Generate this one first
	7	12*	1,1,1,1,1,3,4	// annoying
	9	15	1,1,1,1,1,1,1,3,5
	12	16	1,1,1,1,1,1,1,1,2,2,2,2



	Q: Is it possible to have more factors than K slots?
	A: If that were going to happen it'd be with the Power of 2 test, but 2^(k) > 2*k (for k > 1) and 2^(k-n) > 2*k+n (for k >= n >= 0)

	However, it strongly looks like the focus should be on the numbers to _factor_ rather than approximating any limits; given they increase.  That would also greatly reduce duplicated work.
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

func Euler0088(min, max uint32) uint64 {
	// 32 bit int is enough for Euler 88
	var ret uint64
	var nextK, P, k, iter uint32

	uniq := make(map[uint32]uint16)
	check := make(map[uint16]uint32)

	addPSN := func(k, n uint32) {
		K := uint16(k)
		if oldN, exists := check[K]; !exists || oldN > n {
			// ret += uint64(n)
			if exists {
				if oldN > n {
					fmt.Printf("SET: %5d = %d\n", K, n)
				} else {
					return
				}
			}
			check[K] = n
			if _, exists = uniq[n]; !exists {
				ret += uint64(n)
				uniq[n] = K
			}
		}
	}

	// K increases as terms go up, sum and mul decrease as factors are 'used'
	var MinK func(mul uint32) uint32
	MinK = func(P uint32) uint32 {
		var sum, mul, k, idx, ret uint32
		sum, mul = P, P
		for mul >= sum && 1 < mul && idx <= euler.PrimesSmallU8Mx && mul >= uint32(euler.PrimesSmallU8[idx])*uint32(euler.PrimesSmallU8[idx]) {
			if sum == mul && sum == uint32(euler.PrimesSmallU8[idx]) {
				if 0 < k {
					return 0 // Prime
				}
				addPSN(k+1, P)
				return k + 1
			}
			for 1 < mul && sum >= uint32(euler.PrimesSmallU8[idx]) && 0 == mul%uint32(euler.PrimesSmallU8[idx]) {
				k, mul, sum = k+1, mul/uint32(euler.PrimesSmallU8[idx]), sum-uint32(euler.PrimesSmallU8[idx])
				if 0 < k && sum >= mul {
					ret = k + 1 + sum - mul
					addPSN(ret, P)
					// return k + 1 + sum - mul  // This can continue and provide 5 too though...
				}
			}
			idx++
		}
		fmt.Printf("DEBUG: Euler 88: MinK(%d): EXIT: sum: %d\tmul: %d\tf: #%d\n", P, sum, mul, idx)
		if mul > euler.PrimesSmallU8MxValPow2After {
			fmt.Printf("DANGER: Numbers above %d are handled as primes when they could have prime factors greater than 256: mul: %d\n", euler.PrimesSmallU8MxValPow2After, mul)
			panic("Unexpected, extend prime search if use case intended.")
		}
		if sum >= mul && 0 < k {
			ret = k + 1 + sum - mul
			addPSN(ret, P)
			return ret // The overrun of sum can be corrected with sum-mul ones.
		}
		return ret
	}

	for nextK, P = min, min; P < max<<2; P++ {
		iter++
		if 0 == iter&0xff {
			fmt.Printf("Euler 88: iter %d: nextK: %d\tP: %d\n", iter, nextK, P)
		}
		k = MinK(P)
		fmt.Printf("Euler 88: P %d: got: %d\n", P, k)
		if 1 < k && k >= nextK {
			if k <= max {
				// fmt.Printf("DEBUG: Euler 88: add k: %d val: %d\n", nextK, P)
				addPSN(k, P)
				nextK = k + 1
			} else {
				break
			}
		}
	}

	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 88 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

.
*/
func main() {
	var r uint64
	//test
	r = Euler0088(2, 6)
	if 30 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	r = Euler0088(2, 12)
	if 61 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	fmt.Printf("Euler 88: Passed pretests\n")
	return

	//run
	r = Euler0088(2, 12000)
	fmt.Printf("Euler 88: Product-sum Numbers: %d\n", r)
	if 1097343 != r {
		panic("Did not reach expected value.")
	}
}
