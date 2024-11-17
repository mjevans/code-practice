// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=86
https://projecteuler.net/minimal=86

<p>A spider, S, sits in one corner of a cuboid room, measuring $6$ by $5$ by $3$, and a fly, F, sits in the opposite corner. By travelling on the surfaces of the room the shortest "straight line" distance from S to F is $10$ and the path is shown on the diagram.</p>
<div class="center">
<img src="resources/images/0086.png?1678992052" class="dark_img" alt=""><br></div>
<p>However, there are up to three "shortest" path candidates for any given cuboid and the shortest route doesn't always have integer length.</p>
<p>It can be shown that there are exactly $2060$ distinct cuboids, ignoring rotations, with integer dimensions, up to a maximum size of $M$ by $M$ by $M$, for which the shortest route has integer length when $M = 100$. This is the least value of $M$ for which the number of solutions first exceeds two thousand; the number of solutions when $M = 99$ is $1975$.</p>
<p>Find the least value of $M$ such that the number of solutions first exceeds one million.</p>


/
*/
/*

NOPE - for M=1..100 there are '2060 distinct integer cuboids, ignoring rotations'. it's unclear if there are any collisions of solution length between 2000 and 2060
NOPE - for M=1..99 there are 1975 solutions

After some initial tests I realized they were asking a different question:

NOPE - "Find the (lowest value) of the integer hypotenuse which is the solution for at least one million 'cubes'"

Oh, that's... a bit more difficult...

Problem 75 was supposed to introduce us to...
	https://en.wikipedia.org/wiki/Pythagorean_triple

	...
	Euclid's Formula

	For two
		* coprime (GCD = 1) integers
		* with m > n > 0 ;
		* EXACTLY ONE is even
		k is any integer (for this problem, find the highest integer that fits the perimeter)
	If the resulting (base) a is even, exchange with b

	a, b, c = (m*m - n*n) , (2*m*n) , (m*m + n*n)
	scale the triangle by k*(a,b,c) for 0 < k < max perimeter

https://en.wikipedia.org/wiki/Pythagorean_triple#Not_exchanging_a_and_b
	a, b, c = (m*m - n*n) , (2*m*n) , (m*m + n*n)
	a, b, c = 2*z*m*n, z*(m*m - n*n), (m*m + n*n) with z = 1/2 when both are odd, and 1 if a is even

Still NOPE, but closer.

"It can be shown that there are exactly $2060$ distinct cuboids, ignoring rotations, with integer dimensions, up to a maximum size of $M$ by $M$ by $M$, for which the shortest route has integer length when $M = 100$. This is the least value of $M$ for which the number of solutions first exceeds two thousand; the number of solutions when $M = 99$ is $1975$.
Find the least value of $M$ such that the number of solutions first exceeds one million."

This is written in an intentionally really obtuse way.

When the maximum cuboid side (M) is 100, there are 2060 where shortest path solutions that is an integer.
When M is 99 there are 1975 solutions.

I had to take a break and look at another perspective for more inspiration on this problem.
https://www.hackerrank.com/contests/projecteuler/challenges/euler086/problem

A solution should be optimized towards something that will fulfill HR's output criteria, which changes the approach I'll take a little.

My test case numbers are still off, but the cleaner problem description, a fresh day's pair of eyes, inspired me about where things might have been going wrong.
H must also be the shortest H for a given cuboid.

For binning / sorting the cuboids the longest edge must categorize (so A or B from the triangle should be used when they fit)
B may be the only time that side is seen as a size so it can't be ignored

As solutions are binned by the longest side, those will be covered by a combination of B's splits (from b0 to b1+b2 ... both of which must must be shorter than A to bin there) and lower B's of future A sizes.

Hacker Rank constraints: Number of testcases (IDC), maximum side length 400000, that's plausible

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

func Euler0086(maxSide, required uint32) (uint32, uint32) {
	var a, b, c, n, m, bb, aa, ret uint32
	_, _, _, _ = a, b, c, ret
	var key uint64
	var countAtSize []uint32
	countAtSize = make([]uint32, maxSide+1)
	tfound := make(map[uint64]uint8)
	_ = tfound

	if 0 == required {
		required = ^required
	}

	cfiter := func(aa, jjkk uint32) uint32 {
		// As solutions are binned by the longest side, those will be covered by a combination of B's splits (from b0 to b1+b2 ... both of which must must be shorter than A to bin there) and lower B's of future A sizes.
		// if jj or kk has to be longer than aa, then this isn't the correct place to count the triangle, filter it out
		// if aa > max side this isn't a valid cube either -- but must not even be recorded if that's the case.
		if aa<<1 < jjkk {
			return 0
		}
		// The smaller side of the triangle can split into half as many cuboid edges, a value of 1 becomes 0 which is also valid (3,4,5 is the smallest anyway *shrug*)
		// 2 -> 1,1 ;; 3 -> 1,2 ;; 4 -> 1,3 2,2
		if aa >= jjkk {
			return jjkk >> 1
		}
		// else aa < jjkk but 2x aa >= jjkk ; or equivocally aa >= jjkk / 2
		// aa still has to be the longest side
		// 5,12 gets filtered above, too tall against too short
		// 6,12 (actual is 5,12,13 but for hypothetical since a needs a nudge and this is so close)
		// 6 & 6,6 but not -11,1-, -10,2-, -9,3-, -8,4-, -7,5-
		// 3,4 => 3 1,3 2,2
		// jjkk must be split, so /2 ( >>1 ) ; however that's not enough, 3,4,5 it works, but plausible 6,12 fails one short at 0
		// 8,15,17 => 8 8,7 and -9,6- etc again...
		// flooring division works for both the even split and one shy of that cases
		return aa - (jjkk-1)>>1
	}

Euler0086_outer:
	for n = 1; ; n++ {
		for m = n + 1; ; m++ {
			if (1 != (n+m)&1) || (1 != euler.GCDbin(n, m)) {
				continue
			}
			a, b, c = (m*m - n*n), (2 * m * n), (m*m + n*n)
			//c = m*m + n*n
			if a > b {
				a, b = b, a
			}
			if a > maxSide || b > maxSide<<1 {
				if m == n+1 {
					break Euler0086_outer // 2 to outer loop
				}
				break
			}
			//for 0 == a&1 && 0 == b&1 && 0 == c&1 {
			//	a, b, c = a>>1, b>>1, c>>1
			//}
			//if c*c != a*a + b*b {
			//	fmt.Printf("%9d\tERROR triangle?: %d ,\t%d ,\t%d\n", ret, a, b, c)
			//}
			// aa = euler.GCDbin(a, b)
			// aa = euler.GCDbin(aa, c)
			// if 1 != aa {
			//	fmt.Printf("%9d\tWARNING triangle? factor?: %d ,\t%d ,\t%d // %d\n", ret, a, b, c, aa)
			// }
			key = uint64(a)<<32 | uint64(b)
			if _, exists := tfound[key]; exists {
				tfound[key]++
				fmt.Printf("SKIP triangle: %d ,\t%d ,\t%d\n", a, b, c)
				continue
			}
			tfound[key] = 1
			bb, aa = b, a
			for aa <= maxSide && bb <= maxSide<<1 {
				countAtSize[aa] += cfiter(aa, bb)
				if bb <= maxSide {
					countAtSize[bb] += cfiter(bb, aa)
				}
				aa, bb = aa+a, bb+b
			}
		}
	}

	for aa = 0; aa <= maxSide && ret < required; aa++ {
		ret += countAtSize[aa]
	}
	aa--	// correct the loop overshoot

	if 6 == maxSide {
		// fmt.Println(countAtSize)
	}

	return aa, ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 86 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 86: : 1000457 @ 1818

real    0m0.193s
user    0m0.220s
sys     0m0.062s
.
*/
func main() {
	var ii, r uint32
	//test
	_, r = Euler0086(6, 0)
	if 6 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	_, r = Euler0086(99, 0)
	if 1975 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}
	_, r = Euler0086(100, 0)
	if 2060 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	// ii, r = Euler0086(400_000, 1_000_000)
	ii, r = Euler0086(1886, 1_000_000)
	fmt.Printf("Euler 86: : %d @ %d\n", r, ii)
	if 1818 != ii {
		panic("Did not reach expected value.")
	}
}
