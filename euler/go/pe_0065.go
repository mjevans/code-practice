// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=65
https://projecteuler.net/minimal=65

<p>The square root of $2$ can be written as an infinite continued fraction.</p>
<p>$\sqrt{2} = 1 + \dfrac{1}{2 + \dfrac{1}{2 + \dfrac{1}{2 + \dfrac{1}{2 + ...}}}}$</p>
<p>The infinite continued fraction can be written, $\sqrt{2} = [1; (2)]$, $(2)$ indicates that $2$ repeats <i>ad infinitum</i>. In a similar way, $\sqrt{23} = [4; (1, 3, 1, 8)]$.</p>
<p>It turns out that the sequence of partial values of continued fractions for square roots provide the best rational approximations. Let us consider the convergents for $\sqrt{2}$.</p>
<p>$\begin{align}
&amp;1 + \dfrac{1}{2} = \dfrac{3}{2} \\
&amp;1 + \dfrac{1}{2 + \dfrac{1}{2}} = \dfrac{7}{5}\\
&amp;1 + \dfrac{1}{2 + \dfrac{1}{2 + \dfrac{1}{2}}} = \dfrac{17}{12}\\
&amp;1 + \dfrac{1}{2 + \dfrac{1}{2 + \dfrac{1}{2 + \dfrac{1}{2}}}} = \dfrac{41}{29}
\end{align}$</p>
<p>Hence the sequence of the first ten convergents for $\sqrt{2}$ are:</p>
<p>$1, \dfrac{3}{2}, \dfrac{7}{5}, \dfrac{17}{12}, \dfrac{41}{29}, \dfrac{99}{70}, \dfrac{239}{169}, \dfrac{577}{408}, \dfrac{1393}{985}, \dfrac{3363}{2378}, ...$</p>
<p>What is most surprising is that the important mathematical constant,<br>$e = [2; 1, 2, 1, 1, 4, 1, 1, 6, 1, ... , 1, 2k, 1, ...]$.</p>
<p>The first ten terms in the sequence of convergents for $e$ are:</p>
<p>$2, 3, \dfrac{8}{3}, \dfrac{11}{4}, \dfrac{19}{7}, \dfrac{87}{32}, \dfrac{106}{39}, \dfrac{193}{71}, \dfrac{1264}{465}, \dfrac{1457}{536}, ...$</p>
<p>The sum of digits in the numerator of the $10$<sup>th</sup> convergent is $1 + 4 + 5 + 7 = 17$.</p>
<p>Find the sum of digits in the numerator of the $100$<sup>th</sup> convergent of the continued fraction for $e$.</p>



*/
/*

https://en.wikipedia.org/wiki/Continued_fraction#Regular_patterns_in_continued_fractions
https://en.wikipedia.org/wiki/E_(mathematical_constant)#Computing_the_digits

That looks useful but it's too vague... still it's easier to test something that gnarly as code (or a spreadsheet maybe) rather than by hand.

Didn't produce results that fit the problem; though I think I got kind of close

Maybe a table format would help me understand where that went wrong?

I still couldn't figure it out without additional searches, and more obscure results did spoil me slightly for this problem.

That continued fraction sequence matters too.
e = [2; 1, 2, 1, 1, 4, 1, 1, 6, 1, ... , 1, 2k, 1, ...]
That '2k' notation seems to refer to twice the current cycle count?

2k causes a notable jump in the Num, and the Den, though it's harder to tell at first.

A#	N	D	CF
0	2	1	2;
1	3	1	1
2	8	3	2 2k?
3	11	4	1
4	19	7	1
5	87	32	4 2k?
6	106	39	1
7	193	71	1
8	1264	465	6 2k?
9	1457	536	1

c = (a + 1) % 3   ?   2 * ((a + 1 )/3)   :   1
Nn = N(n-2) + c * N(n-1)
Dn = D(n-2) + c * D(n-1)

I can work with that... and math.Big

*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

// Golang allows + - * and << to work with wraparound for unsigned numbers

/*

https://en.wikipedia.org/wiki/E_(mathematical_constant)#Computing_the_digits

Didn't produce results that fit the problem.

func ContinuedFrac_E_ABP(a, b uint64) uint64 {
	if b <= a {
		return 0
	}
	if b == a+1 {
		return 1
	}
	m := (a + b) >> 1
	x, y, z := ContinuedFrac_E_ABP(a, m), ContinuedFrac_E_ABQ(m, b), ContinuedFrac_E_ABP(m, b)
	m = x * y
	if m < x || m < y || m+z < z {
		fmt.Printf("ERROR: ContinuedFrac_E_ABP overflowed and wrapped around for %d * %d = %d OR %d < %d\n", x, y, m, m+z, z)
	}
	return m + z
}

func ContinuedFrac_E_ABQ(a, b uint64) uint64 {
	if b <= a {
		return 1
	}
	if b == a+1 {
		return b
	}
	m := (a + b) >> 1
	x, y := ContinuedFrac_E_ABQ(a, m), ContinuedFrac_E_ABQ(m, b)
	m = x * y
	if m < x || m < y {
		fmt.Printf("ERROR: ContinuedFrac_E_ABQ overflowed and wrapped around for %d * %d = %d\n", x, y, m)
	}
	return m
}

	fmt.Printf("Euler 65: Convergents of /e/ TEST:\n")
	for a = 1; a <= 10; a++ {
		p, q := ContinuedFrac_E_ABP(0, a), ContinuedFrac_E_ABQ(0, a)
		// Small tweak... this thing is 1 + p/q so.... 1* q/q = (q+p)/q for what Euler65 wants?
		p += q
		gcd := euler.GCDbin(p, q)
		p, q = p/gcd, q/gcd
		fmt.Printf("\t%d", p)
	}
	fmt.Printf("\n")
	for a = 1; a <= 10; a++ {
		p, q := ContinuedFrac_E_ABP(0, a), ContinuedFrac_E_ABQ(0, a)
		// Small tweak... this thing is 1 + p/q so.... 1* q/q = (q+p)/q for what Euler65 wants?
		p += q
		gcd := euler.GCDbin(p, q)
		p, q = p/gcd, q/gcd
		fmt.Printf("\t%d", q)
	}
	fmt.Printf("\n")

Euler 65: Convergents of /e/ TEST:
        1       3       5       41      103     1237    433     69281   62353   6235301
        1       2       3       24      60      720     252     40320   36288   3628800
Euler 65: Convergents of /e/ TEST:
Euler 65: Odd Period Square Roots: 11
[aeon@nas go]$ for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 65 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
Euler 65: Convergents of /e/ TEST:
        2       5       8       65      163     1957    685     109601  98641   9864101
        1       2       3       24      60      720     252     40320   36288   3628800

*/

/*

2k causes a notable jump in the Num, and the Den, though it's harder to tell at first.

A#	N	D	CF
0	2	1	2;
1	3	1	1
2	8	3	2 2k?
3	11	4	1
4	19	7	1
5	87	32	4 2k?
6	106	39	1
7	193	71	1
8	1264	465	6 2k?
9	1457	536	1

c = (a + 1) % 3   ?   2 * ((a + 1 )/3)   :   1
Nn = N(n-2) + c * N(n-1)
Dn = D(n-2) + c * D(n-1)

*/

// slight improvement in case I need this in the future
// pn2 is over-written, old pn0 becomes new p1, old p2 becomes p2, old p2 becomes next reuse
func EulerEFracNextRot(pn0, pd, p1 *big.Int, anum int64) (*big.Int, *big.Int, *big.Int) {
	anum++
	if 0 == anum%3 {
		anum = (anum / 3) << 1
	} else {
		anum = 1
	}
	//res := big.NewInt(0)
	pn0.SetInt64(anum)
	pn0.Mul(pn0, p1)
	pn0.Add(pn0, pd)
	return pd, p1, pn0
}

func EulerEFracNext(p2, p1 *big.Int, anum uint64) *big.Int {
	anum++
	if 0 == anum%3 {
		anum = (anum / 3) << 1
	} else {
		anum = 1
	}
	res := big.NewInt(0)
	res.SetUint64(anum)
	res.Mul(res, p1)
	res.Add(res, p2)
	return res
}

func Euler0065(place, base uint64, debug bool) uint64 {
	_, _ = place, debug
	var ii uint64
	if 1 >= place {
		return 2
	}
	if 2 == place {
		return 3
	}
	n2, n1, d2, d1, zero := big.NewInt(2), big.NewInt(3), big.NewInt(1), big.NewInt(1), big.NewInt(0)
	var n0, d0, div, rem *big.Int
	// Euler 65 calls the initial term 1 not 0 so; the above state is configured as completed iteration 2
	for ii = 2; ii < place; ii++ {
		n0 = EulerEFracNext(n2, n1, ii)
		// shuffle the pointers / references back
		n2 = n1
		n1 = n0
		if debug {
			d0 = EulerEFracNext(d2, d1, ii)
			// shuffle the pointers / references back
			d2 = d1
			d1 = d0
			fmt.Printf("%2d:\t\t%s\t/ %s\n", ii+1, n0.Text(10), d0.Text(10))
		}
	}
	fmt.Printf("%2d:\t\t%s\t/ %s\t(not computed if not debug)\n", ii, n0.Text(10), d0.Text(10))
	div = big.NewInt(int64(base))
	rem = big.NewInt(0)
	ii = 0
	for 1 == n0.Cmp(zero) {
		n0.DivMod(n0, div, rem)
		ii += rem.Uint64()
	}
	return ii
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 65 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 65: Convergents of /e/ TEST:

	3:             8       / 3
	4:             11      / 4
	5:             19      / 7
	6:             87      / 32
	7:             106     / 39
	8:             193     / 71
	9:             1264    / 465

10:             1457    / 536
10:             1457    / 536   (not computed if not debug)
Euler 65: Convergents of /e/ TEST:      17 == 17        true
100:            6963524437876961749120273824619538346438023188214475670667      / <nil> (not computed if not debug)
Euler 65: Convergents of /e/: 272
*/
func main() {
	var a uint64
	//test

	fmt.Printf("Euler 65: Convergents of /e/ TEST:\n")
	a = Euler0065(10, 10, true)
	fmt.Printf("Euler 65: Convergents of /e/ TEST:\t17 == %d\t%t\n", a, 17 == a)

	//run
	a = Euler0065(100, 10, false)
	fmt.Printf("Euler 65: Convergents of /e/: %d\n", a)
}
