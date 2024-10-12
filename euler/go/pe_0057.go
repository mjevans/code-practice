// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=57
https://projecteuler.net/minimal=57

<p>It is possible to show that the square root of two can be expressed as an infinite continued fraction.</p>
<p class="center">$\sqrt 2 =1+ \frac 1 {2+ \frac 1 {2 +\frac 1 {2+ \dots}}}$</p>
<p>By expanding this for the first four iterations, we get:</p>
<p>$1 + \frac 1 2 = \frac  32 = 1.5$<br>
$1 + \frac 1 {2 + \frac 1 2} = \frac 7 5 = 1.4$<br>
$1 + \frac 1 {2 + \frac 1 {2+\frac 1 2}} = \frac {17}{12} = 1.41666 \dots$<br>
$1 + \frac 1 {2 + \frac 1 {2+\frac 1 {2+\frac 1 2}}} = \frac {41}{29} = 1.41379 \dots$<br></p>
<p>The next three expansions are $\frac {99}{70}$, $\frac {239}{169}$, and $\frac {577}{408}$, but the eighth expansion, $\frac {1393}{985}$, is the first example where the number of digits in the numerator exceeds the number of digits in the denominator.</p>
<p>In the first one-thousand expansions, how many fractions contain a numerator with more digits than the denominator?</p>



*/
/*
Square Root Convergents

Offhand, I don't recall ever processing 'sqrt(2)' directly in math, though in Calculus arbitrary power functions happened, with ^(1/2) sometimes one of them.

However, I don't remember ever converting those to division.

Time to hit Wikipedia again for some semi-obscure math detail the average school books probably at best briefly cover in an easily forgotten way.

https://en.wikipedia.org/wiki/Square_root_of_2
https://en.wikipedia.org/wiki/Methods_of_computing_square_roots
https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Heron's_method

Heron's Method / Babylonian Method is close

X[n+1] = 1/2 * ( X[n] + S / X[n])
X[n+1] = Xn/2 + S / (2 * X[n])

Euler 057
sqrt(2) ~= 1 + 1/2 not quite...

f(0) = 1?  Is this S?  No, the shape isn't quite the same...
f(1) = 1 + 1 / 2
f(2) = 1 + 1 / ( 2 + 1/2 )
f(3) = 1 + 1 / ( 2 + 1/ (2 + 1/2) )

Search: square root 2 repeating 1/2
First __useful__ result
https://math.stackexchange.com/questions/14617/proving-the-continued-fraction-representation-of-sqrt2
Spivak's Calculus in the 2nd Edition, it's Chapter 21, Problem 7 'concept' "continued fraction"

Represents (I'll re-index them to 0 and -1 rather than 1 and n+1) F(0) = 1 , F(n) = 1 + 1 / (1 + F(n-1))

Maybe followup material https://en.wikipedia.org/wiki/Banach_fixed-point_theorem

F(0) = 1
F(n) = 1 + 1 / (1 + F(n-1))

F(1) = 1 + 1 / (1 + 1) = 1 + 1/2
F(2) = 1 + 1 / (1 + 1 + 1 / (1 + 1) = 1 + 1/2)
F(2) = 1 + 1 / ( 2 + 1/2 )
F(3) = 1 + 1 / (1 + 1 + 1 / ( 2 + 1/2 ))
F(3) = 1 + 1 / (2 + 1 / ( 2 + 1/2 ))

Yes, that the correct function expression

F(0) = 1
F(n) = 1 + 1 / (1 + F(n-1))

F(n) = (1 + F(n-1) + 1) / (1 + F(n-1))
F(n) = (2 + F(n-1)) / (1 + F(n-1))

That'll be a nightmare to express; have to break up that denominator first

F(x) = 1 + 1 / (1 + F(x-1))
F(x) = 1 + 1 / (1 + n(x-1)/d(x-1))
F(x) = 1 + 1 / ( (d(x-1) + n(x-1))/d(x-1) )
F(x) = 1 + d(x-1) / ( d(x-1) + n(x-1) )
F(x) = ( d(x-1) + n(x-1) + d(x-1) ) / ( d(x-1) + n(x-1) )
F(x) = n(x) / n(d) =  ( 2*d(x-1) + n(x-1) ) / ( d(x-1) + n(x-1) )

Also, offhand, I expect an answer in the ballpark of 1/2 for any given N, since 2d+n / d+n seems likely to flip flop... but that's just a gut feeling guess.

"<p>In the first one-thousand expansions, how many fractions contain a numerator with more digits than the denominator?</p>"

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

func Euler057(terms, base int) uint64 {

	var digitN, digitD, countNgD uint64

	// Start with F(0) n and d in nPrev and dPrev 1 / 1
	biZero, n, d, gcd := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	biOne, nPrev, dPrev := big.NewInt(1), big.NewInt(1), big.NewInt(1)
	biTwo, biBase := big.NewInt(2), big.NewInt(int64(base))

	for ii := 1; ii <= terms; ii++ {
		// F(x) = n(x) / n(d) =  ( 2*d(x-1) + n(x-1) ) / ( d(x-1) + n(x-1) )
		// Calculate N
		//	2*d(x-1) + n(x-1)
		n.Set(dPrev)
		n.Mul(n, biTwo)
		n.Add(n, nPrev)
		// Calculate D
		//	d(x-1) + n(x-1)
		d.Set(dPrev)
		d.Add(d, nPrev)
		// GCD Reduce
		gcd.GCD(nil, nil, n, d)
		if -1 == biOne.Cmp(gcd) {
			n.Div(n, gcd)
			d.Div(d, gcd)
			fmt.Println("divided")
		}
		// Preserve for next round
		nPrev.Set(n)
		dPrev.Set(d)
		// fmt.Printf("Debug 0057: GCD %4d:\t%s\t/\t%s (cmp = %d) %s\n", ii, nPrev.Text(base), dPrev.Text(base), biOne.Cmp(gcd), gcd.Text(base))

		// Count and compare digits
		// "<p>In the first one-thousand expansions, how many fractions contain a numerator with more digits than the denominator?</p>"
		digitN = 0
		for 0 != biZero.Cmp(n) {
			n.Div(n, biBase)
			digitN++
		}
		digitD = 0
		for 0 != biZero.Cmp(d) {
			d.Div(d, biBase)
			digitD++
		}
		if digitN > digitD {
			countNgD++
		}
		// fmt.Printf("Debug 0057: Itr %4d:\t%d / %d\t\t%s\t/\t%s\n", ii, digitN, digitD, nPrev.Text(base), dPrev.Text(base))
	}
	return countNgD
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 57 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 57: Square Root Convergents: TEST(8, 10) = 1 == 1? true
Euler 57: Square Root Convergents: 153

/
*/
func main() {
	//test
	t := Euler057(8, 10)
	fmt.Printf("Euler 57: Square Root Convergents: TEST(8, 10) = %d == 1? %t\n", t, 1 == t)

	//run
	fmt.Printf("Euler 57: Square Root Convergents: %d\n", Euler057(1000, 10))
}
