// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=68
https://projecteuler.net/minimal=68

<p>Consider the following "magic" 3-gon ring, filled with the numbers 1 to 6, and each line adding to nine.</p>
<div class="center">
<img src="resources/images/0068_1.png?1678992052" class="dark_img" alt=""><br></div>
<p>Working <b>clockwise</b>, and starting from the group of three with the numerically lowest external node (4,3,2 in this example), each solution can be described uniquely. For example, the above solution can be described by the set: 4,3,2; 6,2,1; 5,1,3.</p>
<p>It is possible to complete the ring with four different totals: 9, 10, 11, and 12. There are eight solutions in total.</p>
<div class="center">
<table width="400" cellspacing="0" cellpadding="0"><tr><td width="100"><b>Total</b></td><td width="300"><b>Solution Set</b></td>
</tr><tr><td>9</td><td>4,2,3; 5,3,1; 6,1,2</td>
</tr><tr><td>9</td><td>4,3,2; 6,2,1; 5,1,3</td>
</tr><tr><td>10</td><td>2,3,5; 4,5,1; 6,1,3</td>
</tr><tr><td>10</td><td>2,5,3; 6,3,1; 4,1,5</td>
</tr><tr><td>11</td><td>1,4,6; 3,6,2; 5,2,4</td>
</tr><tr><td>11</td><td>1,6,4; 5,4,2; 3,2,6</td>
</tr><tr><td>12</td><td>1,5,6; 2,6,4; 3,4,5</td>
</tr><tr><td>12</td><td>1,6,5; 3,5,4; 2,4,6</td>
</tr></table></div>
<p>By concatenating each group it is possible to form 9-digit strings; the maximum string for a 3-gon ring is 432621513.</p>
<p>Using the numbers 1 to 10, and depending on arrangements, it is possible to form 16- and 17-digit strings. What is the maximum <b>16-digit</b> string for a "magic" 5-gon ring?</p>
<div class="center">
<img src="resources/images/0068_2.png?1678992052" class="dark_img" alt=""><br></div>


*/
/*

Re-sorting by the number they wanted to construct

Total	Solution Set
9	4,3,2; 6,2,1; 5,1,3
9	4,2,3; 5,3,1; 6,1,2
10	2,5,3; 6,3,1; 4,1,5
10	2,3,5; 4,5,1; 6,1,3
12	1,6,5; 3,5,4; 2,4,6
11	1,6,4; 5,4,2; 3,2,6
12	1,5,6; 2,6,4; 3,4,5
11	1,4,6; 3,6,2; 5,2,4

Observations:

Given the internal structure of one strand's third number being the (clockwise given the shape they want) strand's second number.

Starting with the lowest 'outer' number, then rotating clockwise, means that the third number of that strand controls if the next (to the right) number is highest or next lowest.  The smallest number possible in the third slot determines the outcome key.

The example is a '3-gon' with 3 inner and 3 outer numbers, from 1..6

Two groups must be utilized to form 3 totals, inner and outer.

Gon	Line	Ring	Outer	Pool	(2I+O)/G
3	9	6	15	21	9
3	10	9	12	21	10
3	11	12	9	21	11
3	12	15	6	21	12

Given the rules (starting with the lowest outer number, making the biggest possible concat number); it should be clear that the greatest possible pool is if the outer numbers have the maximum possible numbers and thus also outer sum.  This drives the smallest numbers inside as well.

The addition of that 10 throws a wrench into the works, except it will never be in the first three digits and thus the number of digits to the right of those is fixed.

A 17 digit string is possible if the 10 is used inside, which complicates the algorithm since a suddenly the sets of numbers where 10 is inside gain an extra leading digit and thus win.

The 16 digit sums only rule is expressly asking for numbers with the only 10 inside to be ignored.

Q: Which numbers will be in the outer spikes?

H: [ 2N .. (2N - gon + 1) ]
N == gon ; so really
[ 2*gon .. (gon + 1) ]
For a 3 gon that's 6 and 4 ; for a 5 gone that's 10 and 6 ?

Q: What is the sum of the 'best' line length going to be?

H: The pool sums of values given the earlier hypothesis are:
Outer sum((gon + 1)..2*gon) and Inner sum(1..gon) ; which also total to 1..(2*gon)

For a 5-gon the total pool is that N(N+1)/2 which is 10*11 / 2 == 55
Pool 55 (correct); hypothesis: Inner 15 Outer 40

Guess, numbers inside are used twice... does that relate to the possible line total?

(2I + O)/gon == total ; thus also, must be a multiple of gon if there's a valid solution

That rule seems to hold for 3 gons...

14 = ( 15*2 + 40 ) / 5

Experimental test:

14	10 + 3 + 1
14	 9 + 1 + 4
14	 8 + 4 + 2
14	 7 + 2 + 5
14	 6 + 5 + 3

Theory: Biggest + Smallest number, middle is the difference from the desired total, cross-weave the numbers in decreasing order.

The math works out that the smallest arm's second number is always 'gon'

*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	"strconv"
	"strings"
)

func Euler0068(gon int) string {
	pool := (2 * gon * (2*gon + 1)) / 2
	inner := (gon * (gon + 1)) / 2
	outer := pool - inner
	linet := (2*inner + outer)
	if 0 != linet%gon {
		return fmt.Sprintf("Magic is gone?  %d %% %d is not 0\n", linet, gon)
	}
	linet /= gon

	// make n-gon and setup the outer ring
	iiLim := 3 * gon
	sl := make([]int, iiLim)
	num := linet - 2*gon - 1
	// Setup the starting state, second number written back to the final number
	sl[0], sl[1], sl[iiLim-1] = gon+1, gon, gon
	num = gon + 2
	for ii := iiLim - 3; 3 <= ii; ii -= 3 {
		gap := linet - num - sl[ii+2]
		sl[ii], sl[ii+1], sl[ii-1] = num, gap, gap
		num++
	}
	var sb strings.Builder
	sb.Grow(16) // I know it's 16 chars and there is a function from an earlier Euler problem to figure this out, but it doesn't make sense to calculate it for toy problems like this.
	for ii := 0; ii < iiLim; ii++ {
		sb.WriteString(strconv.FormatInt(int64(sl[ii]), 10))
	}
	return sb.String()
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 68 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

65310 31914 84272 5

Euler 68: Magic 5-gon Ring TEST: 432621513
Euler 68: Magic 5-gon Ring: 6531031914842725
*/
func main() {
	//test
	fmt.Printf("Euler 68: Magic 5-gon Ring TEST: %s\n", Euler0068(3))

	//run
	fmt.Printf("Euler 68: Magic 5-gon Ring: %s\n", Euler0068(5))
}
