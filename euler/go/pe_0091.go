// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=91
https://projecteuler.net/minimal=91

<p>The points $P(x_1, y_1)$ and $Q(x_2, y_2)$ are plotted at integer co-ordinates and are joined to the origin, $O(0,0)$, to form $\triangle OPQ$.</p>

<div class="center">
<img src="resources/images/0091_1.png?1678992052" class="dark_img" alt=""><br></div>

<p>There are exactly fourteen triangles containing a right angle that can be formed when each co-ordinate lies between $0$ and $2$ inclusive; that is, $0 \le x_1, y_1, x_2, y_2 \le 2$.</p>

<div class="center">
<img src="resources/images/0091_2.png?1678992052" alt=""><br></div>

<p>Given that $0 \le x_1, y_1, x_2, y_2 \le 50$, how many right triangles can be formed?</p>

/
*/
/*
	https://en.wikipedia.org/wiki/Right_triangle#Characterizations
	At a glance, none of the identities seem to involve triangle points

	Semiperimiter = 1/2*(a + b +c)
	a*a + b*b = c*c
	(s-a)*(s-b)=s*(s-c)

	Lines
	https://en.wikipedia.org/wiki/Line_(geometry)#Linear_equation
	https://en.wikipedia.org/wiki/Linear_equation#Two-point_form
	Just distance alone...

	p0 = (0,0) which helps with two of the lines.
	L*L	= p1x*p1x + p1y*p1y
	R*R	= p2x*p2x + p2y*p2y
	F*F	= abssub(p2y,p1y)^2 + abssub(p2x,p1x)^2

	If F*F == L*L + R*R then it's a right triangle... unless
	Not a triangle at all if all three points are on the same line:
	p1x == p2x == 0 || p1y == p2y == 0
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

func Euler0091(max uint32) uint32 {
	abssub := func(a, b uint32) uint32 {
		if a < b {
			return b - a
		}
		return a - b
	}
	var ret, LL, RR, FF, temp, p1x, p1y, p2x, p2y uint32
	for p1x = 0; p1x <= max; p1x++ {
		for p2x = 0; p2x <= max; p2x++ {
			if 0 == p1x && 0 == p2x {
				continue
			}
			for p1y = 0; p1y <= max; p1y++ {
				for p2y = 0; p2y <= max; p2y++ {
					// Test for on the Y line		Either Point equal to the fixed 0,0 point				Or both points are the same point
					if (0 == p1y && 0 == p2y) || (0 == p1y && 0 == p1x) || (0 == p2y && 0 == p2x) || (p2y == p1y && p2x == p1x) {
						continue
					}
					// F*F	= abssub(p2y,p1y)^2 + abssub(p2x,p1x)^2
					FF, temp = abssub(p2y, p1y), abssub(p2x, p1x)
					FF = FF*FF + temp*temp
					// p0 = (0,0) which helps with two of the lines.
					LL, RR = p1x*p1x+p1y*p1y, p2x*p2x+p2y*p2y
					if LL+RR == FF || LL+FF == RR || LL == RR+FF {
						// fmt.Printf("91: Right Triangle (%d,%d) (%d,%d) = %d ~ %d ~ %d\n", p1x, p1y, p2x, p2y, LL, RR, FF)
						ret++
					} else {
						// fmt.Printf("91: Reject (%d,%d) (%d,%d) = %d + %d == %d\n", p1x, p1y, p2x, p2y, LL, RR, FF)
					}
				}
			}
		}
	}

	return ret >> 1 // The above iterates ALL points and yields double, since it counts swapped points again.  A cache seems likely to be less effective on a modern CPU for large numbers of points; modifying the iteration start patterns of the inner loops didn't seem quite correct either.
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 91 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

.
*/
func main() {
	var r uint32
	//test
	r = Euler0091(2)
	fmt.Printf("Euler 91: Test: %d\n", r)
	if 14 != r {
		panic("Did not reach expected value.")
	}

	//run
	r = Euler0091(50)
	fmt.Printf("Euler 91: Right Triangles with Integer Coordinates: %d\n", r)
	if 14234 != r {
		panic("Did not reach expected value.")
	}
}
