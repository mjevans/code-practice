// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=75
https://projecteuler.net/minimal=75

<p>It turns out that $\pu{12 cm}$ is the smallest length of wire that can be bent to form an integer sided right angle triangle in exactly one way, but there are many more examples.</p>
<ul style="list-style-type:none;">
<li>$\pu{\mathbf{12} \mathbf{cm}}$: $(3,4,5)$</li>
<li>$\pu{\mathbf{24} \mathbf{cm}}$: $(6,8,10)$</li>
<li>$\pu{\mathbf{30} \mathbf{cm}}$: $(5,12,13)$</li>
<li>$\pu{\mathbf{36} \mathbf{cm}}$: $(9,12,15)$</li>
<li>$\pu{\mathbf{40} \mathbf{cm}}$: $(8,15,17)$</li>
<li>$\pu{\mathbf{48} \mathbf{cm}}$: $(12,16,20)$</li></ul>
<p>In contrast, some lengths of wire, like $\pu{20 cm}$, cannot be bent to form an integer sided right angle triangle, and other lengths allow more than one solution to be found; for example, using $\pu{120 cm}$ it is possible to form exactly three different integer sided right angle triangles.</p>
<ul style="list-style-type:none;">
<li>$\pu{\mathbf{120} \mathbf{cm}}$: $(30,40,50)$, $(20,48,52)$, $(24,45,51)$</li></ul>

<p>Given that $L$ is the length of the wire, for how many values of $L \le 1\,500\,000$ can exactly one integer sided right angle triangle be formed?</p>

/
*/
/*
/	Offhand, I remember there were a series of right angle triangle equations on Wikipedia, most of which involved at least one angle, but it had another with three sides.
	a*a + b*b = c*c
	It therefore might be useful to have a map of known large whole number to integer square roots, which would be cheaper to construct in reverse. (From roots to squared)

	https://en.wikipedia.org/wiki/Right_triangle#Characterizations

	(s - a)*(s - b) = s*(s-c)

	Problem 75 also lists the SHORTEST side first.  It does not count 4,3,5 as a different triangle, it has the same ordered set of side lengths.

	The shortest side must be no more than 1/3rd of the triangle (less, actually, but my offhand knowledge can't quickly quantify how much less)
	The hypotenuse can be at MOST almost 1/2 of the perimeter

	Not fast enough.

	"If the lengths of all three sides of a right triangle are integers, the triangle is called a Pythagorean triangle and its side lengths are collectively known as a Pythagorean triple." -- Wikipedia

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

type PythTriple struct {
	a, b, c uint32
}

func (l PythTriple) Eq(r PythTriple) bool {
	return l.a == r.a && l.b == r.b && l.c == r.c
}

func Euler0075(max uint32) uint32 {
	var a, b, c, k, n, m, p, ii uint32
	var count uint32
	var pfound []PythTriple
	pfound = make([]PythTriple, max+1)
	_ = pfound

Euler0075_outer:
	for n = 1; ; n++ {
		for m = n + 1; ; m++ {
			if (1 != (n+m)&1) || (1 != euler.GCDbin(n, m)) {
				continue
			}
			a, b, c = (m*m - n*n), (2 * m * n), (m*m + n*n)
			p = a + b + c
			if p > max && m == n+1 {
				break Euler0075_outer // 2 to outer loop
			}
			if p > max {
				break
			}
			//if 0 == a&1 {
			if a > b {
				a, b = b, a
			}
			for ii = p; ii <= max; ii += p {
				if pfound[ii].Eq(PythTriple{0, 0, 0}) {
					// k = ii / p
					// fmt.Printf("%d\t= ( %d,\t%d,\t%d )\n", p*k, a*k, b*k, c*k)
					pfound[ii] = PythTriple{a, b, c}
					continue
				}
				if !pfound[ii].Eq(PythTriple{1, 1, 1}) {
					k = euler.GCDbin(a, b)
					k = euler.GCDbin(k, c)
					k = euler.GCDbin(k, pfound[a].a)
					//k := uint32(1)
					a, b, c = a/k, b/k, c/k
					if !pfound[ii].Eq(PythTriple{a, b, c}) {
						// fmt.Printf("@%9d\t%9d\t!!! ( %d,\t%d,\t%d ) != %v\n", ii, p, a, b, c, pfound[ii])
						pfound[ii] = PythTriple{1, 1, 1}
					}
				}
			}
		}
	}
	for a = 1; a <= max; a++ {
		if !pfound[a].Eq(PythTriple{0, 0, 0}) && !pfound[a].Eq(PythTriple{1, 1, 1}) {
			count++
		}
	}
	return count
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 75 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 75: Singular Integer Right Triangles:     Count: 161667

real    0m0.180s
user    0m0.220s
sys     0m0.068s
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	r := Euler0075(48)
	if 6 != r {
		panic(fmt.Sprintf("Euler 75: Expected 6 got %d", r))
	}

	//run
	r = Euler0075(1_500_000)
	fmt.Printf("Euler 75: Singular Integer Right Triangles:\tCount: %d\n", r)
	if 161667 != r {
		panic("Did not reach expected value.")
	}
}
