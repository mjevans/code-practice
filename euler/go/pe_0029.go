// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=29
https://projecteuler.net/minimal=29


<p>Consider all integer combinations of $a^b$ for $2 \le a \le 5$ and $2 \le b \le 5$:</p>
\begin{matrix}
2^2=4, &amp;2^3=8, &amp;2^4=16, &amp;2^5=32\\
3^2=9, &amp;3^3=27, &amp;3^4=81, &amp;3^5=243\\
4^2=16, &amp;4^3=64, &amp;4^4=256, &amp;4^5=1024\\
5^2=25, &amp;5^3=125, &amp;5^4=625, &amp;5^5=3125
\end{matrix}
<p>If they are then placed in numerical order, with any repeats removed, we get the following sequence of $15$ distinct terms:
$$4, 8, 9, 16, 25, 27, 32, 64, 81, 125, 243, 256, 625, 1024, 3125.$$</p>
<p>How many distinct terms are in the sequence generated by $a^b$ for $2 \le a \le 100$ and $2 \le b \le 100$?</p>





	?	PP	?
	2^2	t	2		4
	2^3	t	2		8
	2^4	f	B		16
	2^5	t	2		32
	3^2	t	2		9
	3^3	t	2		27
	3^4	f	1		81
	3^5	t	2		243
	4^2	f	E	2^4	16
	4^3	f	E	2^6	64	**Corrected Power > limit, count
	4^4	f	0	2^8	256	**Corrected Power > limit, count
	4^5	f	E	2^10	1024	**Corrected Power > limit, count
	5^2	t	2		25
	5^3	t	2		125
	5^4	f	B		625
	5^5	t	2		3125


			Test if a smaller base could generate this value in bounds? ProperDivisors{{baseGCD(cached)}}*exp
	2^6	4^3	8^2		64
	2^9	-	8^3
	2^12	4^6	8^4		4096
	2^15	-	8^5
	2^18	4^9	8^6		262144
	...
			8^n eventually larger than 2^l limit, factor returns to base / becomes stride mark?
			Eventually _all_ factors return to base, the rest all new
			However this also becomes very annoying for E.G. GCD == Fact{2^2,3^2,5} given 2, 3, 4, 5, 6, 10, 12, 15, 18, ...

			Options?  Sieve style sweep of the range?  Iterate the range and invalidate dupes as identified? << Less memory, less setup hassle


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

func Euler029(baseMin, baseMax, powMin, powMax uint32) uint32 {
	var ret uint32
	for bb := baseMin; bb <= baseMax; bb++ {
		bbfact := euler.Primes.Factorize(uint64(bb))
		bbRedu, bbExp := bbfact.ExtractPower()
		//
		// Do this, but subtract ProperFactors iterations from ret after adding the expected uniques.
		// FIXME: Factorized.ProperFactors() []uint?
		if 1 == bbExp {
			ret += powMax - powMin + 1
			// fmt.Printf("\t%d\t%d\n", bb, powMax-powMin+1)
			continue
		}
		// retOld := ret
	E29_filter_low_powers:
		for pp := powMin; pp <= powMax; pp++ {
			ppTarget := bbExp * uint32(pp)
			ppDivs := *(euler.Primes.Factorize(uint64(ppTarget)).ProperDivisors())
			ppDLen := len(ppDivs)
			// ppDivs[0] = uint64(ppTarget) // 1 := ppTarget
			for ii := 0; ii < ppDLen; ii++ {

				// ppShifted is the effective power, ppDivs[ii] is the base multiple
				ppShifted := ppTarget / uint32(ppDivs[ii])
				// fmt.Printf("Eval: %d^%d (%d)\t%d\t^%d\t", bb, pp, ppTarget, ppDivs[ii], ppShifted)

				// GUARD isPowerInLimit?
				if ppShifted > powMax || ppShifted < powMin {
					// fmt.Printf("Seek more, ^%d outside of ^%d .. ^%d (%d)\n", ppShifted, powMin, powMax, ppDivs[ii])
					continue
				}

				// Base above floor? Then it's already seen
				bbShifted := bbRedu.Pow(uint32(ppDivs[ii])).Uint64()
				if uint64(bb) > bbShifted && bbShifted >= uint64(baseMin) {
					// fmt.Printf("Shadowed: %d^%d == %d^%d (%d)\n", bb, pp, bbShifted, ppShifted, ppDivs[ii])
					continue E29_filter_low_powers
				}
			}
			// New number
			// fmt.Printf("New: %d^%d\n", bbfact.Uint64(), pp)
			ret++
		}
		// fmt.Printf("%d : %d\t%d\n", bb, bbExp, ret-retOld)
	}
	return ret
}

func Euler029_BruteForce(baseMin, baseMax, powMin, powMax uint32) uint32 {
	var ret uint32
	seen := make(map[string]int8, 0)
	for bb := baseMin; bb <= baseMax; bb++ {
		factbase := euler.Primes.Factorize(uint64(bb))
		for pp := powMin; pp <= powMax; pp++ {
			q := factbase.Pow(uint32(pp)).BigInt().String() // bb^pp
			// fmt.Printf("\t%d", q)
			if _, found := seen[q]; false == found {
				seen[q] = 1
				ret++
			} else {
				seen[q]++
			}
		}
	}
	// fmt.Printf("\n")
	return ret
}

/*
	for ii in */ /*.go ; do go fmt "$ii" ; done ; for ii in 29 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:       true   Brute Force:  15        Quick:  15
Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:       true   Brute Force:  23        Quick:  23
Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:       true   Brute Force:  34        Quick:  34
Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:       true   Brute Force:  44        Quick:  44
Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:       true   Brute Force:  54        Quick:  54
Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:       true   Brute Force:  69        Quick:  69
Euler029: Result ..100 (quick):          9183
Euler029: Result ..100 (brute):          9183

*/
func main() {
	//test
	var q, bf uint32
	for ii := uint32(5); ii <= 10; ii++ {
		bf = Euler029_BruteForce(2, ii, 2, ii)
		q = Euler029(2, ii, 2, ii)
		fmt.Println("Euler029: test SELECT COUNT DISTINCT(a^b) a[2..5] b[2..5]:\t", bf == q, "\tBrute Force: ", bf, "\tQuick: ", q)
	}

	//run
	q = Euler029(2, 100, 2, 100)
	fmt.Println("Euler029: Result ..100 (quick):\t\t", q)
	bf = Euler029_BruteForce(2, 100, 2, 100)
	fmt.Println("Euler029: Result ..100 (brute):\t\t", bf)

}
