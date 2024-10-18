// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=64
https://projecteuler.net/minimal=64

Odd Period Square Roots

<p>All square roots are periodic when written as continued fractions and can be written in the form:</p>

$\displaystyle \quad \quad \sqrt{N}=a_0+\frac 1 {a_1+\frac 1 {a_2+ \frac 1 {a3+ \dots}}}$

<p>For example, let us consider $\sqrt{23}:$</p>
$\quad \quad \sqrt{23}=4+\sqrt{23}-4=4+\frac 1 {\frac 1 {\sqrt{23}-4}}=4+\frac 1  {1+\frac{\sqrt{23}-3}7}$

<p>If we continue we would get the following expansion:</p>

$\displaystyle \quad \quad \sqrt{23}=4+\frac 1 {1+\frac 1 {3+ \frac 1 {1+\frac 1 {8+ \dots}}}}$

<p>The process can be summarised as follows:</p>
<p>
$\quad \quad a_0=4, \frac 1 {\sqrt{23}-4}=\frac {\sqrt{23}+4} 7=1+\frac {\sqrt{23}-3} 7$<br>
$\quad \quad a_1=1, \frac 7 {\sqrt{23}-3}=\frac {7(\sqrt{23}+3)} {14}=3+\frac {\sqrt{23}-3} 2$<br>
$\quad \quad a_2=3, \frac 2 {\sqrt{23}-3}=\frac {2(\sqrt{23}+3)} {14}=1+\frac {\sqrt{23}-4} 7$<br>
$\quad \quad a_3=1, \frac 7 {\sqrt{23}-4}=\frac {7(\sqrt{23}+4)} 7=8+\sqrt{23}-4$<br>
$\quad \quad a_4=8, \frac 1 {\sqrt{23}-4}=\frac {\sqrt{23}+4} 7=1+\frac {\sqrt{23}-3} 7$<br>
$\quad \quad a_5=1, \frac 7 {\sqrt{23}-3}=\frac {7 (\sqrt{23}+3)} {14}=3+\frac {\sqrt{23}-3} 2$<br>

$\quad \quad a_6=3, \frac 2 {\sqrt{23}-3}=\frac {2(\sqrt{23}+3)} {14}=1+\frac {\sqrt{23}-4} 7$<br>
$\quad \quad a_7=1, \frac 7 {\sqrt{23}-4}=\frac {7(\sqrt{23}+4)} {7}=8+\sqrt{23}-4$<br>
</p>

<p>It can be seen that the sequence is repeating. For conciseness, we use the notation $\sqrt{23}=[4;(1,3,1,8)]$, to indicate that the block (1,3,1,8) repeats indefinitely.</p>

<p>The first ten continued fraction representations of (irrational) square roots are:</p>
<p>
$\quad \quad \sqrt{2}=[1;(2)]$, period=$1$<br>
$\quad \quad \sqrt{3}=[1;(1,2)]$, period=$2$<br>
$\quad \quad \sqrt{5}=[2;(4)]$, period=$1$<br>
$\quad \quad \sqrt{6}=[2;(2,4)]$, period=$2$<br>
$\quad \quad \sqrt{7}=[2;(1,1,1,4)]$, period=$4$<br>
$\quad \quad \sqrt{8}=[2;(1,4)]$, period=$2$<br>
$\quad \quad \sqrt{10}=[3;(6)]$, period=$1$<br>
$\quad \quad \sqrt{11}=[3;(3,6)]$, period=$2$<br>
$\quad \quad \sqrt{12}=[3;(2,6)]$, period=$2$<br>
$\quad \quad \sqrt{13}=[3;(1,1,1,1,6)]$, period=$5$
</p>
<p>Exactly four continued fractions, for $N \le 13$, have an odd period.</p>
<p>How many continued fractions for $N \le 10\,000$ have an odd period?</p>



*/
/*

Here 'we' go again...

Maybe better results / explanation for the topic?

https://en.wikipedia.org/wiki/Continued_fraction
https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Continued_fraction_expansion

Quadratic Irrationals ~= ( a + Sqrt(b) ) / ( c )

S == a positive number for which the square root is to be found
a == initial guess
r == 'remainder'
S = a*a + r

Solve for R
r = S - a*a

That one quadratic formula... x^2 + x - x - a^2 ?
(x + a)(x - a)

S - a*a = (Sqrt(S) + a)(Sqrt(S) - a) = r
(Sqrt(S) + a)(Sqrt(S) - a) = r
(Sqrt(S) - a) = r / (Sqrt(S) + a)
Sqrt(S) - a = ( r / (Sqrt(S) + a) )
Sqrt(S) = a + ( r / (Sqrt(S) + a) )
That's it!
Sqrt(S) = a + ( r / ( a + Sqrt(S) ) )

Now how to get the initial estimate of a?

S - a*a = r  // Find the biggest a which leaves a positive r

The example given is Sqrt(23)

A0 = ? ; r = ?
S = a*a + r

For Sqrt(23)...

23 = 4*4 + r => 23-16 = 7 = r

Grab that reference...
Sqrt(S) = a + ( r / ( a + Sqrt(S) ) )
...
S - a*2 = (Sqrt(S) + a)(Sqrt(S) - a) = r
23 - 4*2 = (Sqrt(S) + 4)(Sqrt(S) - 4) = 7

a0:Sqrt(S) = 4 + ( 7 / ( 4 + Sqrt(S) ) )

a1:Sqrt(S) = 4 + ( 7 / ( 4 + 4 + ( 7 / ( 4 + Sqrt(S) ) ) ) )


a0:Sqrt(S) = a0 + ( r / ( a0 + Sqrt(S) ) )

a0 ??? Sqrt(S) - a0 = ( r / ( a0 + Sqrt(S) ) )
a0 ??? (Sqrt(S) - a0)( a0 + Sqrt(S) ) = r
a0 ??? S - a*a = r


a1:Sqrt(S) = a0 + ( r / ( a0 + a1 + ( r / ( a1 + Sqrt(S) ) ) ) )
a1:Sqrt(S) = a0 + ( r / ( a0 + a1 + ( r / ( a1 + Sqrt(S) ) ) ) )

Solve for a?
S - a^2 = r
S - r = a^2
a = (S - r) / a
4 = (23-7) :> 16 / 4 ... but that's not helpful

This explained a key bit more clearly...
https://math.stackexchange.com/questions/265690/continued-fraction-of-a-square-root

a(n) == Floor(Sqrt(S)) == Floor(a0 + ( r / ( a0 + Sqrt(S) ) ))

The initial a0 can be obtained experimentally by trial, or more easily in a modern program...
math.Floor(math.Sqrt(S)) to get the next equal or smaller perfect square.

https://math.stackexchange.com/questions/4487826/properties-of-continued-fractions-of-sqrt-n

The continued fraction thus _begins_ life as the term:

Sqrt(S) ~= a + 1/x  // ONE is chosen as the numerator since it is the div/multiplicative Identity number and easy to extract: it scales x to itself, and signifies that X is less than one...
Sqrt(S) - a ~= 1/x
x * ( Sqrt(S) - a ) ~= 1
x = 1 / ( Sqrt(S) - a )		//  <<< That's the start of the expansion for a0 from the problem
x = ( 1 / ( Sqrt(S) - a ) ) * 1
x = ( 1 / ( Sqrt(S) - a ) ) * ( ( Sqrt(S) + a ) / ( Sqrt(S) + a ))
x = ( ( Sqrt(S) + a ) / ( ( Sqrt(S) + a )) * ( Sqrt(S) - a ) ) )
x = ( ( Sqrt(S) + a ) / ( S - a*a )

is 1/x == r ?

Hum... "This reminds me of a Puzzle" *grep* Euler 0057?

I didn't want to reference that, since my goal with these isn't just to write better Golang code (rather than know the general concept of programming better from other languages); it's also to better understand the problems that are presented and thus expand my knowledge overall.

I'll take the one breadcrumb I needed and continue forward.
Though I'm keeping that initial A and why the **** the first numerator was a ONE.


Grab that reference...
Sqrt(S) = a + ( r / ( a + Sqrt(S) ) )
Sqrt(S) = a + ( r / ( a + n--/d-- ) )
Sqrt(S) = a + ( r / ( (( a * d-- ) / d--) + n--/d-- ) )
Sqrt(S) = a + ( r / ( a * d--  + n--/d-- ) )
Sqrt(S) = a + ( r / ( a * d--  + n--/d-- ) ) * ( d-- / d-- )
Sqrt(S) = a * Y/Y + ( r * d-- ) / ( a * d--  +  n-- )
Sqrt(S) = ( a*a * d-- + a * n-- + r * d-- ) / ( a * d--  +  n-- )
Sqrt(S) = ( a * n--  +  a*a * d--  +  r * d-- ) / ( a * d--  +  n-- )

This, just looks ugly because all the A's (and r) aren't 1 like they were in Sqrt(2)

Sqrt(2) = ( 1 * n--  +  1*1 * d--  +  1 * d-- ) / ( 1 * d--  +  n-- )
Sqrt(2) = ( n--  +  d--  +  d-- ) / ( d--  +  n-- )
Sqrt(2) = ( 2 * d-- + n-- ) / ( d--  +  n-- )

However, that didn't work out and continued deriving larger fractions like in 57, even with the whole number removed.


Sqrt(S) ~= a + 1/x
x = 1 / ( Sqrt(S) - a )
x = ( ( Sqrt(S) + a ) / ( S - a*a )
?
Sqrt(S) ~= a + 1/x
a ~= Sqrt(S) - 1/x

Sqrt(S) ~= a + 1/x
x = 1 / ( Sqrt(S) - a )
x = ( ( Sqrt(S) + a ) / ( S - a*a )
Put X back into the initial equation???
Sqrt(S) ~= a + 1/x


https://projecteuler.net/problem=64
https://en.wikipedia.org/wiki/Continued_fraction
https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Continued_fraction_expansion
https://en.wikipedia.org/wiki/Solving_quadratic_equations_with_continued_fractions
https://en.wikipedia.org/wiki/Quadratic_irrational_number#Square_root_of_non-square_is_irrational
https://math.stackexchange.com/questions/265690/continued-fraction-of-a-square-root
https://math.stackexchange.com/questions/4487826/properties-of-continued-fractions-of-sqrt-n
https://mathworld.wolfram.com/PeriodicContinuedFraction.html

A0 = Floor(Sqrt(S))
Column 1.2 : combines with An values per recursion
x = 1 / ( Sqrt(S) - a )		//  <<< That's the start of the expansion for a0 from the problem
x = n / d
n0 = 1
dn = ( Sqrt(S) - An )

Column 2 : Special term first term. (Probably R = 1?)
x = n/d = ( ( Sqrt(S) + a ) / ( S - a*a )

Column 3 : Next An has been extracted from the fraction
::WANT::
Anext = Floor(x) << Column 2->3
? ~= Anext + (Sqrt(S) - v) / ( w )
w MIGHT be ( S - anext*anext )
v MIGHT be ( a - w )


Concept notes for myself / future readers
* Continued Fractions represent a close 'estimate' (less than the real value)
* r or frequently 1/x are the Remainder, which is always less than one
* r or 1/x uses 1/ as the identity value; x = 1/r - it's just an easier way to express a fractional remainder
* Each An + ( 1/x || r ) pair refines the precision, but prevents the iterative exponentiation of the improper fractions...

* Given how An is extracted; each An must be between 1 and (base - 1) ??? (Maybe not?)
* && if a == 0 the there is no remainder, the fraction has terminated
* SIGNATURE = aPrev, ???

Euler 64 expresses the a0+1/tree as: Sqrt(S) = A0 + Sqrt(S) - A0 => A0 + 1/(1/(Sqrt(S) - 4)) [[Stage 3]] => Stage 4

Sqrt(S) ~= a + r = a + 1/x
r = 1/x	=>	x = 1/r	=>	x = 1 / (Sqrt(S) - aCur)
Sqrt(S) ~= a + r = a + 1/x =>	a + 1/ ( 1/ (Sqrt(S) - a) ) == Sqrt(S) <<<== Euler 64, equation line 2, 'consider Sqrt(23) stage 3

Introspect Further

	(Added: x step 1)
Sqrt(S) ~= a + 1/x

How does Stage 3 (Added: X step 2)
a0 + 1/ ( 1/ (Sqrt(S) - a0) )

Become Stage 4? (Added: X step 3)
a0 + 1/ ( a1 + (Sqrt(S) - v) / ( w )

Work backwards?
a0 + 1/ ( a1*w/w + (Sqrt(S) - v) / ( w )
a0 + 1/ ( ( a1*w + Sqrt(S) - v) / ( w )
a0 + 1/ ( ( Sqrt(S) - v + a1*w ) / ( w )

That early quadratic trick: a*a - b*b <=> (a - b) * (a + b)
w = (Sqrt(S) - ?)*(Sqrt(S) + ?)

Stage 3 (X step 2)
a0 + 1/ ( **1** * 1/ (Sqrt(S) - a0) )
a0 + 1/ ( (Sqrt(S) + a0)/(Sqrt(S) + a0) * 1/ (Sqrt(S) - a0) )
a0 + 1/ ( (Sqrt(S) + a0)/((Sqrt(S) + a0) * (Sqrt(S) - a0) ))
a0 + 1/ ( (Sqrt(S) + a0)/(S - a0 * a0) )  <<-- Derived  X step 2
Sqrt(23) = 4 + 1 / ( (Sqrt(23) + 4)/(23 - 16)) => R ~ 1/X Term Only == 1 / ( (Sqrt(23) + 4)/7 ) == R

A column header list would have REALLY helped me understand Euler 64's table

Term (a#)	| Val	| x (== 1/r ==) ... => ... =>
a0		| 4	| 1 / (Sqrt(23) - 4) => ...

xn = (Sqrt(S) + aPrev)/(S - aPrev * aPrev)
(calculator) 1.256...
an = Floor(xn)

a0 = aPrev = Floor(Sqrt(S))
a#	Floor(xn) xn = (Sqrt(S) + a)/(S - a*a) == num / den
4	1	1.256...
1		... not quite correct, it got flip-flopped (price is right, not logic gates)

xn = (Sqrt(S) + a)/(S - a*a) == num / den
Den is an integer, easy to flatten... Num is the aPrev and Sqrt(S) - lossy to flatten? Might be possible as a float?

Equation
Sqrt(S) ~= a + 1/x
x = 1 / ( Sqrt(S) - a )
x = n / d

INITIAL
ap = aPrev = Floor(Sqrt(S))
dp = 1

1/x flips it over so...
nn = dp
dn = Sqrt(S) - ap


dp / (Sqrt(S) - a)
dp * (Sqrt(S) + a)/(Sqrt(S) + a) * 1/ (Sqrt(S) - a)
dp * (Sqrt(S) + a)/((Sqrt(S) + a) * (Sqrt(S) - a) )
dp * (Sqrt(S) + a)/(S - a * a)  <<-- Derived  X step 2


// Some sort of logic error...
		an? = Floor(dp * (Sqrt(S) + a)/(S - a * a))
F		xn = 1/ ( Sqrt(S) - a ) = an + 1/ ( Sqrt(S) - a ) - an
A		an + 1/ ( Sqrt(S) - a ) - an
I		an + 1/ ( Sqrt(S) - a ) - an * ( Sqrt(S) - a ) / ( Sqrt(S) - a )
L		an + ( 1 - an * ( Sqrt(S) - a )) / ( Sqrt(S) - a )
		an + ( 0 - an * Sqrt(S) - a * an + 1)) / ( Sqrt(S) - a )
F		( 0 - an * Sqrt(S) - a * an + 1)) / ( Sqrt(S) - a )   *   ( Sqrt(S) + a )/( Sqrt(S) + a )
A		( 0 - an * Sqrt(S) - a * an + 1)) / ( Sqrt(S) - a )   *   ( Sqrt(S) + a )/( Sqrt(S) + a )
I		( Sqrt(S) + a ) * ( 0 - an * Sqrt(S) - a * an + 1)) / ( S - a*a )
L		( -an * S  -  a*an * Sqrt(S) + Sqrt(S) - a*an * Sqrt(S) - a*a*an + a ) / ( S - a*a )
		( -an * S  -  2*a*an * Sqrt(S) + Sqrt(S) - a*a*an + a ) / ( S - a*a )
		Eval for example 23 step 0
		( -1 * 23  -  8 * Sqrt(S) + Sqrt(S) - 16 + 4 ) / ( 23 - 16 )
		( -  8 * Sqrt(S) + Sqrt(S) - 16 + 4 - 35 ) / ( 23 - 16 )
// Do this again after some brain rest

I did some other things for an hour or two and am trying again.

This would be so much easier to work on on paper, or if I used an equation syntax that graphs (like Euler has)...

https://projecteuler.net/problem=64

From the top
Sqrt(S) = a0 + 1/x0  => Recurse x0 = a1 + 1/x1 =Loop

0:=	Sqrt(S)
1:=	a0 + Sqrt(S) - a0
=>	a0 + ( Sqrt(S) - a0 ) * ( ( 1 / ( Sqrt(S) - a0 ) ) / ( 1 / ( Sqrt(S) - a0 ) ) )
2:=	a0 + 1/ ( 1/ ( Sqrt(S) - a0 ) )
=>	a0 + 1/ ( ( 1/ ( Sqrt(S) - a0 ) ) * (( Sqrt(S) + a0 ) / ( Sqrt(S) + a0 )) )
=>	a0 + 1/ ( ( Sqrt(S) + a0 ) / ( S - a0*a0 ) )		<<<===	X step 2
x0 = ( Sqrt(S) + a0 ) / ( S - a0*a0 )
x0 = a1 + 1/x1
How do they transform?
=>	( Sqrt(S) + a0 ) / ( S - a0*a0 )
=>	Sqrt(S) / ( S - a0*a0 )  +  a0 / ( S - a0*a0 )

Time to hit the references again and see what the more detailed algorithm is...
https://en.wikipedia.org/wiki/Continued_fraction#Basic_formula
Interesting, 'may be negative', so that's how negative numbers can be represented, and 0 for fractions less than a whole unit which are continued.  Neither of those matter for Euler 64

https://en.wikipedia.org/wiki/Continued_fraction#Calculating_continued_fraction_representations

Integer part (a) and then just subtract it, which explains...

=>	( Sqrt(S) + a0 ) / ( S - a0*a0 )
!!! aNext = Floor( n / d )
=>	aNext + ( Sqrt(S) + a0 ) / ( S - a0*a0 ) - aNext
=>	aNext + ( Sqrt(S) + a0 ) / ( S - a0*a0 ) - aNext* ( ( S - a0*a0 )/( S - a0*a0 ) )
=>	aNext + ( Sqrt(S) + a0 ) / ( S - a0*a0 ) - ( aNext*S - a0*a0*aNext )/( S - a0*a0 )
3:=	aNext + ( Sqrt(S) + a0 - aNext*S + a0*a0*aNext ) / ( S - a0*a0 )

Check23 1 + ( Sqrt(S) + 4 - 23 + 16 ) / ( S - a0*a0 ) => 1 + ( Sqrt(S) - 3 ) / ( S - a0*a0 )

Oh... that happens to match A0 in the first example, but it's really the variable part of the numerator...

SPECIAL CASE: Iteration 0
a0 = Floor(Sqrt(S))
dSave = 1
nSave = a0

!! GCD(dSave, dNew)
aNext = Floor( n / d ) = Floor(   )

NOTE: If using  ( S - nSave*nSave ) without GCD reduction, divide by dSave too to match extracting dSave from the numerator:  ( S - nSave*nSave ) * dSave / dSave

=>	aNext + ( Sqrt(S) + a0 ) / ( d ) - aNext
=>	aNext + ( Sqrt(S) + a0 ) / ( d ) - aNext * d/d
( Sqrt(S) + a0 - d*aNext ) / ( d )



Lesson learned?   The next time I need to do higher Maths ; use some paper, or at least a stylus (pen and paper is presently more cost effective)

*/

import (
	// "bufio"
	"euler"
	"fmt"
	"math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

type ContinuedFracPartV3 struct {
	a, n, d int64 // NOTE: the Sqrt(S) is not stored, it's assumed as part of the context
}

func Euler0064(limit uint64, debug bool) uint64 {
	var cfii, cfmx, iter, countOdd, countEven, countTerm uint64
	var Si, a, gcd, d, nSave, dSave int64
	_ = gcd
	_ = euler.Primes
	var S, Sqr float64
Euler0064NextV3: // There were two old versions that used incorrect math basis
	for Si = 2; Si <= int64(limit); Si++ {
		iter = 0
		cf := make([]ContinuedFracPartV3, 0, 8)
		cfmx = 0
		S = float64(Si)
		Sqr = math.Sqrt(S)
		// a0 = uint64(math.Floor(Sqr)) // but the int cast already floors it since positive number
		dSave = 1
		nSave = int64(Sqr) // implicit floor
		// fmt.Printf("%d\tnSave=%d\t%f\n", Si, nSave, Sqr)
		for {
			iter++
			d = Si - nSave*nSave
			if 0 == d {
				countTerm++
				// Fraction terminated, no repeat
				continue Euler0064NextV3
			}

			a = int64(float64(dSave) / (Sqr - (float64(nSave))))
			// fmt.Printf("%d = %f\tdS=%d\tnS=%d\td=%d\n", a, af, dSave, nSave, d)

			// GCD _after_ creating the next A term
			gcd = euler.GCDbin(dSave, d)
			d, dSave = d/gcd, dSave/gcd
			if 1 != dSave {
				fmt.Printf("debug:\tShouldn't dSave == 1 after GCD?\ta=%d\tn=%d\tdS=%d\td=%d\n", a, nSave, dSave, d)
			}

			// dSave * ( Sqrt(S) + a0 - d*aNext ) / ( d )
			// Stored withOUT the context negation
			nSave = 0 - (nSave - a*d)
			dSave = d
			// fmt.Printf("debug:\ta=%d\tn=%d\td=%d\n", a, nSave, dSave)
			for cfii = 0; cfii < cfmx; cfii++ {
				if a == cf[cfii].a && nSave == cf[cfii].n && dSave == cf[cfii].d {
					// Repeat found
					rlen := cfmx - cfii
					if 0 < rlen&1 {
						countOdd++
						if debug {
							fmt.Printf("Sqrt(%5d) Odd# %d\t\tRepeat Period=%d\t\t(offset %d)\n", Si, countOdd, rlen, cfii)
						}
					} else {
						countEven++
						if debug {
							fmt.Printf("Sqrt(%5d) EVEN# %d\t\tRepeat Period=%d\t\t(offset %d)\n", Si, countEven, rlen, cfii)
						}
					}
					continue Euler0064NextV3
				}
			}
			cf = append(cf, ContinuedFracPartV3{a, nSave, dSave})
			cfmx++
			if 0 == iter&0xFF {
				fmt.Printf("Stall check: %5d: iter %d: a=%d\tnSave=%d\tdSave=%d\n", Si, iter, a, nSave, dSave)
			}
		}
	}
	fmt.Printf("Euler0064 2..%d found %d terminating, %d even, and %d odd repeats.\n", limit, countTerm, countEven, countOdd)
	return countOdd
}

/*
Confusion version 2

type ContinuedRatPartU64 struct {
	a, n, d uint64
}

func Euler0064(limit uint64) uint64 {
	var countOdd, countEven, countTerm, Si, cfii, cfmx, iter uint64
	var a, n, d, r, ap, np, dp, g uint64
Euler0064NextSi:
	for Si = 2; Si <= limit; Si++ {

		iter = 0
		cf := make([]ContinuedRatPartU64, 0, 5)
		cfmx = 0
		a = uint64(math.Sqrt(float64(Si))) // implicit Floor()
		r = Si - a*a
		if 0 == r {
			countTerm++
			// Fraction terminated, no repeat
			continue Euler0064NextSi
		}
		np, dp = 1, 1
		for {
			iter++
			// Sqrt(S) = ( a * n--  +  a*a * d--  +  r * d-- ) / ( a * d--  +  n-- )
			n = a*np + a*a*dp + r*dp
			d = a*dp + np
			a = n / d // implicit Floor()
			fmt.Printf("a=%d\tn=%d\td=%d\n", a, n, d)
			// n/d = (ad + p) / d
			// n = ad + p
			// n - ad = p
			n -= a * d
			if 0 == n {
				countTerm++
				// Fraction terminated, no repeat
				continue Euler0064NextSi
			}
			g = euler.GCDbin(n, d)
			n, d = n/g, d/g
			fmt.Printf("Debug check: %5d: iter %d: a=%d ~ %f\tn=%d\td=%d\tr=%d\n", Si, iter, a, float64(n)/float64(d), n, d, r)

			for cfii = 0; cfii < cfmx; cfii++ {
				if ap == cf[cfii].a && n == cf[cfii].n && d == cf[cfii].d {
					// Repeat found
					rlen := cfmx - cfii
					if 0 < rlen&1 {
						countOdd++
						fmt.Printf("Sqrt(%5d) Odd %d\t\tRepeat Length %d (offset %d)\n", Si, countOdd, rlen, cfii)
					} else {
						countEven++
					}
					continue Euler0064NextSi
				}
			}
			cf = append(cf, ContinuedRatPartU64{ap, n, d})
			ap, np, dp = a, n, d
			cfmx++
			if 0 == iter&0xF {
				fmt.Printf("Stall check: %5d: iter %d: a=%d ~ %f\tn=%d\td=%d\tr=%d\n", Si, iter, a, float64(n)/float64(d), n, d, r)
				break
			}
		}
	}
	fmt.Printf("Euler0064 2..%d found %d terminating, %d even, and %d odd repeats.\n", limit, countTerm, countEven, countOdd)
	return countOdd
}
*/

/*
// Trash ALL of this and use Sqrt(S) = ( a * n--  +  a*a * d--  +  r * d-- ) / ( a * d--  +  n-- )
type ContinuedPartFloat64 struct {
	a, x float64
}


func Euler0064(limit uint64) uint64 {
	var countOdd, countEven, countTerm, Si, cfii, cfmx, iter uint64
	var a, S, Sqr, x float64
Euler0064NextS:
	for Si = 2; Si <= limit; Si++ {
		iter = 0
		cf := make([]ContinuedPartFloat64, 0, 8)
		cfmx = 0
		S = float64(Si)
		Sqr = math.Sqrt(S)
		a = 0.0
		x = S // inject into loop
		for {
			iter++
			a = math.Floor(math.Sqrt(x - a))
			if S == a*a {
				countTerm++
				// Fraction terminated, no repeat
				continue Euler0064NextS
			}
			x = (Sqr + a) / (S - (a * a))
			fmt.Printf("Debug check: %5d: iter %d: a=%.0f\tx=%f\n", Si, iter, a, x)
			for cfii = 0; cfii < cfmx; cfii++ {
				if a == cf[cfii].a && x == cf[cfii].x {
					// Repeat found
					rlen := cfmx - cfii
					if 0 < rlen&1 {
						countOdd++
						fmt.Printf("Sqrt(%5d) Odd %d\t\tRepeat Lngth %d (offset %d)\n", Si, countOdd, rlen, cfii)
					} else {
						countEven++
					}
					continue Euler0064NextS
				}
			}
			cf = append(cf, ContinuedPartFloat64{a, x})
			cfmx++
			if 0 == iter&0xFF {
				fmt.Printf("Stall check: %5d: iter %d: a=%.0f\tx=%f\n", Si, iter, a, x)
			}
		}
	}
	fmt.Printf("Euler0064 2..%d found %d terminating, %d even, and %d odd repeats.\n", limit, countTerm, countEven, countOdd)
	return countOdd
}
*/

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 64 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Sqrt(    2) Odd# 1              Repeat Period=1         (offset 0)
Sqrt(    3) EVEN# 1             Repeat Period=2         (offset 0)
Sqrt(    5) Odd# 2              Repeat Period=1         (offset 0)
Sqrt(    6) EVEN# 2             Repeat Period=2         (offset 0)
Sqrt(    7) EVEN# 3             Repeat Period=4         (offset 0)
Sqrt(    8) EVEN# 4             Repeat Period=2         (offset 0)
Sqrt(   10) Odd# 3              Repeat Period=1         (offset 0)
Sqrt(   11) EVEN# 5             Repeat Period=2         (offset 0)
Sqrt(   12) EVEN# 6             Repeat Period=2         (offset 0)
Sqrt(   13) Odd# 4              Repeat Period=5         (offset 0)
Euler0064 2..13 found 2 terminating, 6 even, and 4 odd repeats.
Euler 64: Odd Period Square Roots: TEST 4 true
Euler0064 2..10000 found 99 terminating, 8578 even, and 1322 odd repeats.
Euler 64: Odd Period Square Roots: 1322

*/
func main() {
	var a uint64
	//test
	a = Euler0064(13, true)
	fmt.Printf("Euler 64: Odd Period Square Roots: TEST %d %t\n", a, a == 4)

	//run
	a = Euler0064(10000, false)
	fmt.Printf("Euler 64: Odd Period Square Roots: %d\n", a)
}
