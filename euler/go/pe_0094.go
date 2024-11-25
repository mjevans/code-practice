// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=94
https://projecteuler.net/minimal=94

<p>It is easily proved that no equilateral triangle exists with integral length sides and integral area. However, the <dfn>almost equilateral triangle</dfn> $5$-$5$-$6$ has an area of $12$ square units.</p>
<p>We shall define an <dfn>almost equilateral triangle</dfn> to be a triangle for which two sides are equal and the third differs by no more than one unit.</p>
<p>Find the sum of the perimeters of all <dfn>almost equilateral triangles</dfn> with integral side lengths and area and whose perimeters do not exceed one billion ($1\,000\,000\,000$).</p>


/
*/
/*
	https://en.wikipedia.org/wiki/Acute_and_obtuse_triangles
	From the description, It can't be a right angled triangle, and it isn't an obtuse triangle. So it must be Acute
	a*a + b*b > c*c	&&	b*b + c*c > a*a 	&&	a*a + c*c > b*b

	https://en.wikipedia.org/wiki/Isosceles_triangle
	H = Sqrt(a*a - (b*b)/4)
	The 'altitude' line (vertex of the same length sides to mid-base line of the almost equal side) also forms a pair of ring angle triangles.
	Area
	A(t) = (b/4) * Sqrt(4*a*a - b*b)

	Answer: Sum of perimeters (2a+b) of all integer sided, and integer area Isosceles_triangle's with a +/- 1 = b; where the total perimeter is <= 1E9

	If the square root of the numbers inside is a multiple of 4, then b can otherwise be any integer.  Otherwise B must be a multiple of 4 (or 2 and 2 from the sqrt)

	When in doubt plan it out:
	A	B	Sqrt?
	1	2	0	Err, 0.5
	2	1	15
	2	3	7
	3	2	32
	3	4	20
	4	3	55
	4	5	39
	5	4	84
	5	6	64	8	Winner!
	6	5	119
	6	7	95
	7	6	160
	7	8	132
	8	7	207
	8	9	175
	9	8	260
	9	10	224
	10	9	319
	10	11	279
	11	10	384
	11	12	340
	12	11	455
	12	13	407
	13	12	532
	13	14	480
	14	13	-
	14	15	-	etc
	17	16		30	-- Is not a multiple of 4, but between this and B it is a winner.

	Maybe, if the first ~10-20 positives were brute forced a pattern for doing it faster might emerge.

5, 6    (16)    12
17, 16  (50)    120
65, 66  (196)   1848
241, 240        (722)   25080
901, 902        (2704)  351780
3361, 3360      (10082) 4890480
12545, 12546    (37636) 68149872
46817, 46816    (140450)        949077360
174725, 174726  (524176)        13219419708
652081, 652080  (1956242)       184120982760
2433601, 2433602        (7300804)       2564481115560
9082321, 9082320        (27246962)      35718589344360
33895685, 33895686      (101687056)     497495864091732
126500417, 126500416    (379501250)     6929223155685600

Provided these are correct and complete, I don't see any obvious pattern...


Given what I know now, that A is always odd and B is always even; a faster program with Pythagorean Triples (like one of the earlier programs and mentioned on the Euler forum thread) would be way faster than the blind brute force approach.

Reflecting on this further after sleeping on it:

I didn't recall the Area identity of Right_triangles offhand:
https://en.wikipedia.org/wiki/Right_triangle#Characterizations
Area = a*b / 2 (Right Triangle)
Ra, Rb, Rc == ITh, ITb / 2, ITa

The Area as an Integer value combined with ITb ~ 2*Rb as an integer value restricts Ra~ITh to either an integer or half an integer value.
Then the difference of at most 1 for a and b (I'm unsure exactly how, but logically it's the only likely fence) constrains the values to integer solutions only for everything.

With that set of constraints in place it's logical to expect a very narrow range of Pythagorean Triples as the only possible answers.

I didn't exhaustively examine the super fast C/C++ code on Euler's forum; I wonder if it's small stack happens to implement the above algorithm in a more optimal way?



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

func Euler0094(min, max uint64) uint64 {
	var sumP, ii, inside, area uint64
	var seen uint8
	if min < 6 {
		min = 6
	}
	for ii = min / 3; ; ii++ {
		// perimeter test
		if max < (ii<<1)+ii-1 {
			break
		}
		//	A(t) = (b/4) * Sqrt(4*a*a - b*b)
		inside = ((ii * ii) << 2) - (ii-1)*(ii-1)
		area = euler.SqrtU64(inside)
		if area*area == inside {
			area *= (ii - 1)
			if 0 == area&0b11 {
				sumP += (ii << 1) + ii - 1
				if 20 > seen {
					area, seen = area>>2, seen+1
					fmt.Printf("%d, %d\t(%d)\t%d\n", ii, ii-1, (ii<<1)+ii-1, area)
				}
			}
		}

		// perimeter test
		if max < (ii<<1)+ii+1 {
			continue
		}
		//	A(t) = (b/4) * Sqrt(4*a*a - b*b)
		inside = ((ii * ii) << 2) - (ii+1)*(ii+1)
		area = euler.SqrtU64(inside)
		if area*area == inside {
			area *= (ii + 1)
			if 0 == area&0b11 {
				sumP += (ii << 1) + ii + 1
				if 20 > seen {
					area, seen = area>>2, seen+1
					fmt.Printf("%d, %d\t(%d)\t%d\n", ii, ii+1, (ii<<1)+ii+1, area)
				}
			}
		}
	}

	return sumP
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 94 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

5, 6    (16)    12
17, 16  (50)    120	// Test run
5, 6    (16)    12
17, 16  (50)    120
65, 66  (196)   1848
241, 240        (722)   25080
901, 902        (2704)  351780
3361, 3360      (10082) 4890480
12545, 12546    (37636) 68149872
46817, 46816    (140450)        949077360
174725, 174726  (524176)        13219419708
652081, 652080  (1956242)       184120982760
2433601, 2433602        (7300804)       2564481115560
9082321, 9082320        (27246962)      35718589344360
33895685, 33895686      (101687056)     497495864091732
126500417, 126500416    (379501250)     6929223155685600
Euler 94: Almost Equilateral Triangles: 518408346

real    3m39.818s
user    3m39.623s
sys     0m0.319s
.
*/
func main() {
	var r uint64
	//test
	r = Euler0094(1, 50)
	if 66 != r {
		panic("Did not reach expected value.")
	}

	//run
	r = Euler0094(1, 1_000_000_000)
	fmt.Printf("Euler 94: Almost Equilateral Triangles: %d\n", r)
	if 518408346 != r {
		panic("Did not reach expected value.")
	}
}
