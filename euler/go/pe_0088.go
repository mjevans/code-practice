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

	I'm a little unsure where to go from here; the low unit tests pass.  However the check value on Euler says 32679757 is incorrect...
	Probably need to lookup the number sequence.
	https://oeis.org/search?q=4+6+8+8+12+12+12+15&go=Search
	https://oeis.org/A104173 "a(n) is the smallest integer equal to the sum and the product of the same n positive integers: a(n) = i(1) + i(2) + ... + i(n) = i(1)*i(2)*...*i(n)."
	     4, 6, 8, 8,	 5
	12, 12, 12, 15, 16,	10
	16, 16, 18, 20, 24,	15
	24, 24, 24, 24, 28,	20
	27, 32, 30, 48, 32,	25
	32, 32, 36, 36, 36,	30
	42, 40, 40, 48, 48,	35
	48, 45, 48, 48, 48,	40
	48, 48, 54, 60, 54,	45
	56, 54, 60, 63, 60,	50
	60, 60, 63, 64, 64,	55
	64, 64, 64, 70, 72,	60
	72, 72, 72, 72, 72,	65
	84, 80, 80, 81, 80, 80

	71 (well 72 with the 1 included) numbers?  A nightmare to total... 1082 = 4+6+8+12+15+16+18+20+24+28+27+32+30+48+36+42+40+48+54+60+63+64+70+72+84+80+81

add:    28 =     40	<< Incorrect, due to missed factorization of 36?
add:    31 =     42
add:    32 =     44	<< WRONG should be 40
add:    37 =     45
add:    36 =     48	++ Wow lots of factors of 48
add:    41 =     50	<< WRONG
add:    38 =     52	<< WRONG
add:    45 =     54
add:    43 =     56	<< 54 not 56
add:    46 =     60	<< 56 not 60
add:    49 =     63
add:    50 =     64	<< 60 not 64
add:    51 =     66	<< 60 not 66
add:    56 =     72	<< 64 not 72
add:    65 =     75	<< 72 not 75
add:    64 =     80	<< 72 not 80
add:    69 =     81
add:    67 =     84	<< 80 not 84
add:    71 =     88	<< 80 not 88
add:    63 =     96	<< 72 not 96
add:    60 =    120	<< 72 not 120

40 is the smallest *shrug*
2,2,2,5 = 11 (~4) 1^17
36
2,2,3,3 = 10 + 1^18


/
*/

import (
	// "bufio"
	// "euler"
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
	_, _, _ = iter, nextK, k

	uniq := make(map[uint32]uint16)  // Required
	check := make(map[uint16]uint32) // oldN > n guard

	addPSN := func(k, n uint32) {
		if k > max {
			return
		}
		K := uint16(k)
		if oldN, exists := check[K]; !exists || oldN > n {
			// ret += uint64(n)
			check[K] = n
			if _, exists = uniq[n]; !exists {
				// fmt.Printf("add: %5d = %6d\n", K, n)
				ret += uint64(n)
				uniq[n] = K
			}
		}
	}

	// mxCompat := func[INT ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](a, b INT) INT {
	// Can't do Generics with functions as variables?
	mxCompat := func(a, b uint32) uint32 {
		if a < b {
			return b
		}
		return a
	}

	// K increases as terms go up, sum and mul decrease as factors are 'used'
	var minK func(P, sum, mul, k, fact uint32) uint32
	minK = func(P, sum, mul, k, fact uint32) uint32 {
		var ret uint32					// I'm not using ret at all, but this code might get reused in a HackerRank test I've not looked at yet, and I might need it then.

		// what-if: some number needs 6,6,... it could be a square root and series of ones
		if mul > fact*fact {
			ret = minK(P, sum, mul, k, fact+1) // Try larger factors first E.G. 2,2,3 -> 3,4
		}
		for ; mul >= fact*fact; fact++ {
			if sum == mul && sum == fact {
				if 0 < k {
					return ret // Prime, but the other side(s) might have found a composite
				}
				// fmt.Printf("add A: %5d = %6d\n", k+1, P)
				addPSN(k+1, P)
				return mxCompat(k+1, ret)
			}
			if 1 < mul && sum >= fact && 0 == mul%fact {
				for 1 < mul && sum >= fact && 0 == mul%fact {
					k, mul, sum = k+1, mul/fact, sum-fact
					if 0 < k && sum >= mul {
						rthis := k + 1 + sum - mul
						ret = mxCompat(ret, rthis)
						// fmt.Printf("add B: %5d = %6d\n", rthis, P)
						addPSN(rthis, P)
					}
					ret = mxCompat(ret, minK(P, sum, mul, k, fact+1))
				}
			}
		}
		// fmt.Printf("DEBUG: Euler 88: minK(%d): sum: %d\tmul: %d\tk: %d\tf: %d\n", P, sum, mul, k, fact)
		// Mul must be prime, so if it fits AND isn't the only factor
		if 0 < k && sum >= mul && mul+1 != fact {
			rthis := k + 1 + sum - mul
			// fmt.Printf("add C: %5d = %6d\n", rthis, P)
			addPSN(rthis, P)
			ret = mxCompat(ret, rthis) // The overrun of sum can be corrected with sum-mul ones.
		}
		return ret
	}

	MinK := func(P uint32) uint32 { return minK(P, P, P, 0, 2) }

	for nextK, P = min, min; P <= max<<1; P++ {
		MinK(P)
	}

	// fmt.Println(check)
	if 71 == max {
		test := []uint8{0, 1, 4, 6, 8, 8, 12, 12, 12, 15, 16, 16, 16, 18, 20, 24, 24, 24, 24, 24, 28, 27, 32, 30, 48, 32, 32, 32, 36, 36, 36, 42, 40, 40, 48, 48, 48, 45, 48, 48, 48, 48, 48, 54, 60, 54, 56, 54, 60, 63, 60, 60, 60, 63, 64, 64, 64, 64, 64, 70, 72, 72, 72, 72, 72, 72, 84, 80, 80, 81, 80, 80}
		for iter = min; iter <= max; iter++ {
			if uint32(test[iter]) != check[uint16(iter)] {
				fmt.Printf("Failed oeis.org/A104173 [%d] expected %d got %d\n", iter, test[iter], check[uint16(iter)])
			}
		}
	}
	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 88 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 88: Passed pretests
Euler 88: Product-sum Numbers: 7587457

real    0m1.507s
user    0m1.527s
sys     0m0.071s
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
	r = Euler0088(2, 71)
	if 1135 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	fmt.Printf("Euler 88: Passed pretests\n")

	//run
	r = Euler0088(2, 12000)
	fmt.Printf("Euler 88: Product-sum Numbers: %d\n", r)
	if 7587457 != r {
		panic("Did not reach expected value.")
	}
}
