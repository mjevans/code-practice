// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=44
https://projecteuler.net/minimal=44

<p>Pentagonal numbers are generated by the formula, $P_n=n(3n-1)/2$. The first ten pentagonal numbers are:
$$1, 5, 12, 22, 35, 51, 70, 92, 117, 145, \dots$$</p>
<p>It can be seen that $P_4 + P_7 = 22 + 70 = 92 = P_8$. However, their difference, $70 - 22 = 48$, is not pentagonal.</p>
<p>Find the pair of pentagonal numbers, $P_j$ and $P_k$, for which their sum and difference are pentagonal and $D = |P_k - P_j|$ is minimised; what is the value of $D$?</p>

*/
/*

func P(n) { return (n(3n-1))/2 }

something like...

c > b > a > s
P(s) = (b(3b-1))/2 - (a(3a-1))/2
P(c) = (b(3b-1))/2 + (a(3a-1))/2

(s(3s-1))/2 = (b(3b-1))/2 - (a(3a-1))/2
(c(3c-1))/2 = (b(3b-1))/2 + (a(3a-1))/2

Get rid of those divisions
*2

s(3s-1) = b(3b-1) - a(3a-1)
c(3c-1) = b(3b-1) + a(3a-1)

s(3s-1) = b(3b-1) - a(3a-1)
c(3c-1) - a(3a-1) = b(3b-1)

s(3s-1) = c(3c-1) - a(3a-1) - a(3a-1)

Might not be right and doesn't seem like a strong path forward either.

func P(n) { return (n(3n-1))/2 }

Maybe a quick test?
Result = (n(3n-1))/2

The larger the series gets the more are between every pair, irrespective of spacing; but a wider diff is a lower addition and closer sub...

An absolute bound to beat for distance would help with evaluating lower tests at wider windows.

R = (n(3n-1))/2
2R = n(3n-1)
... I need an estimate
2R ~~ 3*n*n
2R/3 ~~ n*n

Use "math" for a float64 square root instead of the digital approximation


""Find the pair of pentagonal numbers P(a) and P(b),
for which P(a) + P(b) && P(b) - P(b) are pentagonal
and D = |P(b) - P(a)| is minimised;
what is the value of $D$""


*/

import (
	// "bufio"
	// "euler"
	"fmt"
	"math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func PentagonNumber(n uint64) uint64 {
	// Euler 44
	return (n * (3*n - 1)) >> 1
}

func IsPentagonNumber(n uint64) bool {
	if 0 == n {
		return false
	}
	// Estimate floor using: sqrt ( 2n / 3 )
	floor := uint64(math.Sqrt(float64(n<<1) / float64(3)))
	var test uint64
	for ii := floor; test < n; ii++ {
		test = PentagonNumber(ii)
	}
	return test == n
}

func ReversePentagonNumber(n uint64) uint64 {
	if 0 == n {
		return 1
	}
	// Estimate floor using: sqrt ( 2n / 3 )
	floor := uint64(math.Sqrt(float64(n<<1) / float64(3)))
	var test, ii uint64
	for ii = floor; test < n; ii++ {
		test = PentagonNumber(ii)
	}
	if test == n {
		return ii - 1
	}
	return 0
}

func Euler044(limit uint64) (uint64, uint64) {
	var prevPB, b, pb, mindiff, minA, minB, gap uint64
	if 0 == limit {
		limit = uint64(0xFFFF_FFFF_FFFF)
	}
	mindiff = limit
	pb = PentagonNumber(1)
	b = 1
	for limit == mindiff {
		prevPB = pb
		b++
		pb = PentagonNumber(b)
		if mindiff < pb-prevPB {
			fmt.Printf("lowest diff <= cur-last for %d <= %d\t( %d - %d )\n", mindiff, pb-prevPB, pb, prevPB)
			// search space exhausted
			break
		}
		for gap = 1; gap < b; gap++ {
			pa := PentagonNumber(b - gap)
			if mindiff < pb-pa {
				fmt.Printf("lowest diff <= cur-last for %d <= %d\t( %d - %d )\n", mindiff, pb-pa, pb, pa)
				// next highest start
				break
			}
			if IsPentagonNumber(pb-pa) && IsPentagonNumber(pb+pa) {
				minA = b - gap
				minB = b
				mindiff = pb - pa
				fmt.Printf("Found span %d ( lower bound %d and %d ): %d < %d +/- %d < %d\n", gap, b-gap, b, pb-pa, pa, pb, pb+pa)
			}
		}
	}
	return minA, minB
}

// This was my first attempt.  It works, but it's too slow because some guesses I made about the nature of the solution were wrong.  Mostly I expected a big-ish span for a gap distance of 1, but to be found quickly.
func Euler044_tooslow(limit uint64) (uint64, uint64) {
	var gap, ii, cur, last, minA, minB, lowest uint64
	gap = 0
	// Around when I added growing the gap for the _first_ solution is about when I might have wanted to step back and re-evaluate...
Euler044FindBound:
	for spin := uint64(0); spin < limit; spin++ {
		gap++
		ii = 1
		cur = 0
		// This is also a problem, it needs a limit at all, instead of self-calibrating
		for cur < limit {
			last = PentagonNumber(ii)
			cur = PentagonNumber(ii + gap)
			if IsPentagonNumber(cur-last) && IsPentagonNumber(cur+last) {
				fmt.Printf("Found span %d (first lower bound): %d < %d +/- %d < %d\n", gap, cur-last, last, cur, cur+last)
				lowest = cur - last
				minA = ii
				minB = ii + gap
				break Euler044FindBound
			}
			ii++
		}
		ii--
		// fmt.Printf("Failed to find a match along span %d ending in (%d) %d\t,\t(%d) %d\n", gap, ii, last, ii+gap, cur)
		// if 5 < gap {
		// panic("debug")
		// }
	}
	if 0 == lowest {
		fmt.Printf("Did not converge\n")
	}
	// This looks closer to the more optimal algorithm I found later, though the check against the prior run is moved to an outer loop condition
	for {
		gap++
		ii = 1
		for {
			last = PentagonNumber(ii)
			cur = PentagonNumber(ii + gap)
			if 1 == ii {
				// fmt.Printf("%d GAP started %d\t,\t%d\n", gap, last, cur)
				if lowest < cur-last {
					// fmt.Printf("%d lowest <= cur-last for %d <= %d\t( %d - %d )\n", gap, lowest, cur-last, cur, last)
					return minA, minB
				}
			}
			if lowest <= cur-last {
				// fmt.Printf("%d + %d lowest <= cur-last for %d <= %d\t( %d - %d )\n", gap, ii, lowest, cur-last, cur, last)
				break // 1 to gap++
			}
			// guarded cur - last < lowest
			if IsPentagonNumber(cur-last) && IsPentagonNumber(cur+last) {
				fmt.Printf("Found span %d (lower bound): %d < %d +/- %d < %d\n", gap, cur-last, last, cur, cur+last)
				lowest = cur - last
				minA = ii - gap
				minB = ii
				break // 1 to gap++
			}
			ii++
		}
	}
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 44 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 44: tests
Found span 1147 ( lower bound 1020 and 2167 ): 5482660 < 1560090 +/- 7042750 < 8602840
lowest diff <= cur-last for 5482660 <= 5485718  ( 7042750 - 1557032 )
Euler 44: Pentagon Numbers: D = >>> 5482660 <<<          P 1020 and P 2167 ( 1560090  ,  7042750 )


*/
func main() {
	//test
	fmt.Printf("Euler 44: tests\n")
	for ii := uint64(1); ii < 65; ii++ {
		if ii != ReversePentagonNumber(PentagonNumber(ii)) || false == IsPentagonNumber(PentagonNumber(ii)) {
			fmt.Printf("FAILED: %d expected but got %d and %t\n", ii, ReversePentagonNumber(PentagonNumber(ii)), IsPentagonNumber(PentagonNumber(ii)))
		}
	}

	//run
	a, b := Euler044(uint64(0xFF_FFFF_FFFF))
	fmt.Printf("Euler 44: Pentagon Numbers: D = >>> %d <<< \t P %d and P %d ( %d  ,  %d )\n", PentagonNumber(b)-PentagonNumber(a), a, b, PentagonNumber(a), PentagonNumber(b))
}
