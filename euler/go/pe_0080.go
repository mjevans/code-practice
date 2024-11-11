// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=80
https://projecteuler.net/minimal=80

<p>It is well known that if the square root of a natural number is not an integer, then it is irrational. The decimal expansion of such square roots is infinite without any repeating pattern at all.</p>
<p>The square root of two is $1.41421356237309504880\cdots$, and the digital sum of the first one hundred decimal digits is $475$.</p>
<p>For the first one hundred natural numbers, find the total of the digital sums of the first one hundred decimal digits for all the irrational square roots.</p>

/
*/
/*
	The 'digital' (10s digits) summation of the first 100 (after the decimal) digits for the square roots of the first 100 positive numbers...
	My initial idea was to reuse the (round down) integer square root representation, multiply the number inside by 100 after each time, and subtract the square of the outside number, and start again; but I was concerned about how to propagate the digit shift inside the square root and wanted to check my gut instinct.

	https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Decimal_(base_10)
	Makes sense, to extract (mod)10 (a digit) the square of 10 is 100

	Total of the digits would easily fit standard integers (I correctly guessed it'd PROBABLY fit in uint16 but safely uint32).
	However COMPUTING those digits, past the largest float's precision, requires VERY large integers.  In Wikipedia's variable terms, p is the problem.

	Golang is nice, math/big is part of the basic language standard library.

	In other languages I might reach for GMP or some other well known library for handling BIG numbers.

	If this comes up on a site that runs code for me, and doesn't include that standard library with Golang; I'll pick a language that does support large numbers just for that issue.  If a bunch of problems come up that require 'big' number support maybe I'll make my own limited/big package.
/
*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0080(min, max, digits uint32) uint32 {
	var ii, jj, sum, L, R, P uint32

	// Using the terms from the Wikipedia article for clarity:
	// c : Current Value (under the square root) / Remainder
	// p : part of the root so far  ,  y : new estimate adjustment  ,  x : estimate of the next digit
	c, p, y, x, tmpA := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	ten, hun, twenty, zero := big.NewInt(10), big.NewInt(100), big.NewInt(20), big.NewInt(0)
	// Some constants for operations later, and the some 'constants'

	for ii = min; ii <= max; ii++ {
		// Step 2 (initial)
		p.Sqrt(c.SetInt64(int64(ii))) // The left side of the decimal
		// modification to the full squareroot method; the integer part of the result is ignored, and there's no fraction to bring down... just *100
		// Step 3 + Step 1 (AlgMod: no new digits, just zeros)
		c.Mul(c.Sub(c, x.Mul(p, p)), hun) // c = (c - x*x) * 100 // where initially X is found as the integer part of the square root

		if 0 == zero.Cmp(c) {
			fmt.Printf("Skip %d, a rational square of %s\n", ii, p.String())
			continue
		}

		jj = 0 // count the 'precision' of the leading digits
		for P = uint32(p.Uint64()); 0 < P; P /= 10 {
			jj++
			sum += P % 10
		}

		// Step 4 Continue IF further remainder is not zero AND less than desired precision
		for ; 0 != zero.Cmp(c) && jj < digits; jj++ {
			// Step 2
			// -- The algorithm says to 'guess' an initial x, but offers an exact thing to guess.
			// x.Div(x.Div(c, twenty), p) // goal -- This estimate didn't work out

			// X has to be between [0 and 9] binary search time
			L, R = 0, 9
			for L != R {
				P = (L + R + 1) >> 1 // ceil
				x.SetInt64(int64(P))
				y.Mul(x, tmpA.Add(tmpA.Mul(twenty, p), x)) // y = x*(20*p +x)
				cmp := zero.Cmp(tmpA.Sub(c, y))
				if 1 == cmp {
					R = P - 1
				} else {
					L = P
				}
			}
			// y may have been calculated off an R that was one too large, force refresh outside the loop
			x.SetInt64(int64(R))
			y.Mul(x, tmpA.Add(tmpA.Mul(twenty, p), x)) // y = x*(20*p +x)
			// fmt.Printf("debug: %d %d %d\t%d\t%s\t%s\t%s\n", L, P, R, 0-0, c.String(), y.String(), p.String())

			// sqrt(2) = 1. 41421 35623 73095 04880
			// fmt.Printf("\tc: %s\t\ty: %s\n", c.String(), y.String())
			p.Add(p.Mul(p, ten), x)
			sum += R

			// Step 3
			c.Mul(c.Sub(c, y), hun) // c = (c-y)*100
			// fmt.Print(x.String())
		}
	}
	fmt.Println()
	return sum
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 80 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Skip 0, a rational square of 0
Skip 1, a rational square of 1
Skip 4, a rational square of 2
Skip 9, a rational square of 3
Skip 16, a rational square of 4
Skip 25, a rational square of 5
Skip 36, a rational square of 6
Skip 49, a rational square of 7
Skip 64, a rational square of 8
Skip 81, a rational square of 9
Skip 100, a rational square of 10

Euler 80: Square Root Digital Expansion: 40886

real    0m0.124s
user    0m0.180s
sys     0m0.048s
.
*/
func main() {
	var r uint32
	//test
	r = Euler0080(2, 2, 100)
	if 475 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	r = Euler0080(0, 100, 100)
	fmt.Printf("Euler 80: Square Root Digital Expansion: %d\n", r)
	if 40886 != r {
		panic("Did not reach expected value.")
	}
}
