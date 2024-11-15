// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=85
https://projecteuler.net/minimal=85

<p>By counting carefully it can be seen that a rectangular grid measuring $3$ by $2$ contains eighteen rectangles:</p>
<div class="center">
<img src="resources/images/0085.png?1678992052" class="dark_img" alt=""></div>
<p>Although there exists no rectangular grid that contains exactly two million rectangles, find the area of the grid with the nearest solution.</p>

/
*/
/*
	They want the area, which is great since the order of the sides doesn't matter.
	Anywhere from 1xN to just over half N.

	For 1xN the answer is (N * (N+1))/2

	1 x 1 : 1
	1 x 2 : 3
	1 x 3 : 6		x2
	1 x 4 : 10
	1 x 5 : 15		x3
	1 x 6 : 21
	1 x 7 : 28		x4
	1 x 8 : 36
	1 x 9 : 45		x5
	1 x 10 :        55
	1 x 11 :        66	x6
	1 x 12 :        78
	1 x 13 :        91	x7
	1 x 14 :        105
	1 x 15 :        120	x8

	2 x 2 : 9
	2 x 3 : 18
	2 x 4 : 30
	2 x 5 : 45

	3 x 3 : 36
	3 x 4 : 60
	3 x 5 : 90

	4 x 4 : 100
	4 x 5 : 150

	5 x 5 : 225	(25^2) 625 so well under the square root

	For 1xN the answer is (N * (N+1))/2
	Solve for N?: I/2 = N*N + N
	I/2 = N * (N+1)
	The rhs^2 term isn't negative so it doesn't decompose, but that square root is in there as an approximation

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

func Euler0085gridsize(x, y uint64) uint64 {
	var ret, ii, jj uint64
	for ii = 1; ii <= x; ii++ {
		for jj = 1; jj <= y; jj++ {
			// rectangle of ii x jj but it's area doesn't matter, just where it can be shifted
			ret += (x - ii + 1) * (y - jj + 1)
		}
	}
	return ret
}

func Euler0085bestfit(target, x, yEst uint64) (uint64, uint64) {
	var r, p uint64
	r = Euler0085gridsize(x, yEst)
	if target == r {
		return 0, yEst
	} else if target < r {
		for yEst--; target < r; yEst-- {
			p, r = r, Euler0085gridsize(x, yEst)
		}
		d0, d1 := target-r, p-target
		if d0 < d1 {
			return d0, yEst + 1
		} else {
			return d1, yEst + 2
		}
	} else {
		for yEst++; target > r; yEst++ {
			p, r = r, Euler0085gridsize(x, yEst)
		}
		d1, d0 := target-p, r-target
		if d0 < d1 {
			return d0, yEst - 1
		} else {
			return d1, yEst - 2
		}
	}
}

func Euler0085(target uint64) (uint64, uint64) {
	var left, right, ii, r uint64
	left, right = target, 1
	for left > right {
		left >>= 1
		right <<= 1
	}
	_ = ii
	left = right

	// not quite https://en.wikipedia.org/wiki/Binary_search
	for left+1 <= right {
		ii = (left + right) >> 1
		r = Euler0085gridsize(1, ii)
		if r < target {
			left = ii + 1
		} else if r > target {
			right = ii - 1
		} else {
			return r, ii // exact fit found in 1 x ii == target
		}
	}

	var bestDiff, bestA, diff, y, iiMx, x1y uint64
	bestDiff, bestA = Euler0085bestfit(target, 1, right)
	iiMx, x1y = (bestA+1)>>1, bestA

	fmt.Printf("Euler 85: initial fit: 1x%d fit %d real %d\n", bestA, bestDiff, Euler0085gridsize(1, bestA))
	for ii = 2; ii <= iiMx; ii++ {
		diff, y = Euler0085bestfit(target, ii, x1y/ii)
		if diff < bestDiff {
			fmt.Printf("Euler 85: better fit: %dx%d fit %d real %d\n", ii, y, diff, Euler0085gridsize(ii, y))
			bestDiff, bestA = diff, y*ii
		}
	}

	return bestDiff, bestA
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 85 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 85: initial fit: 1x2000 fit 1000 real 2001000
Euler 85: better fit: 2x1154 fit 695 real 1999305
Euler 85: better fit: 3x816 fit 16 real 2000016
Euler 85: better fit: 36x77 fit 2 real 1999998
Euler 85: Counting Rectangles (off by 2): 2772

real    0m0.125s
user    0m0.153s
sys     0m0.072s
.
*/
func main() {
	var area, r uint64
	//test
	r = Euler0085gridsize(3, 2)
	if 18 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	area, r = Euler0085(2_000_000)
	fmt.Printf("Euler 85: Counting Rectangles (off by %d): %d\n", area, r)
	if 2772 != r {
		panic("Did not reach expected value.")
	}
}
