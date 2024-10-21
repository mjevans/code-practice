// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=66
https://projecteuler.net/minimal=66

<p>Consider quadratic Diophantine equations of the form:
$$x^2 - Dy^2 = 1$$</p>
<p>For example, when $D=13$, the minimal solution in $x$ is $649^2 - 13 \times 180^2 = 1$.</p>
<p>It can be assumed that there are no solutions in positive integers when $D$ is square.</p>
<p>By finding minimal solutions in $x$ for $D = \{2, 3, 5, 6, 7\}$, we obtain the following:</p>
\begin{align}
3^2 - 2 \times 2^2 &amp;= 1\\
2^2 - 3 \times 1^2 &amp;= 1\\
{\color{red}{\mathbf 9}}^2 - 5 \times 4^2 &amp;= 1\\
5^2 - 6 \times 2^2 &amp;= 1\\
8^2 - 7 \times 3^2 &amp;= 1
\end{align}
<p>Hence, by considering minimal solutions in $x$ for $D \le 7$, the largest $x$ is obtained when $D=5$.</p>
<p>Find the value of $D \le 1000$ in minimal solutions of $x$ for which the largest value of $x$ is obtained.</p>




*/
/*

Diophantine Equation

x*x - D*y*y = 1

https://en.wikipedia.org/wiki/Diophantine_equation

One half of Pell's equation?

{\displaystyle x^{2}-ny^{2}=\pm 1} 	This is Pell's equation, which is named after the English mathematician John Pell. It was studied by Brahmagupta in the 7th century, as well as by Fermat in the 17th century.

https://en.wikipedia.org/wiki/Pell%27s_equation

At least some test cases...
https://en.wikipedia.org/wiki/Pell%27s_equation#List_of_fundamental_solutions_of_Pell's_equations

"When D is an integer square, there is no solution except for the trivial solution (1, 0)."

Filter 1: Minimal solutions of X ( over D[1..1000] )
Filter 2: Of those solutions, the largest value of X

x*x = 1 + D*y*y
x = Sqrt(D*y*y + 1)

It's probably time to find the Wikipedia or another math site for derivatives.  I haven't seen the inside math book cover sheet of identities in years.

Intuitively though, that literally says 'x is going to be the square root of a number (x*x) which is 1 above D * the square of another number; which was exactly what the problem statement was.  It's just more obvious to me when written in this second form.

How would I do this without any further prep?

Create a slice of x*x and another for D solutions; for every new X test remaining empty D solutions until all 1000 D are found.

Maybe make that a map[D]x since I don't care about iteration order and I DO care about reducing the tests as the x*x array answers grow.

Could do it with less storage by using math.Sqrt()

https://www.khanacademy.org/math/multivariable-calculus/applications-of-multivariable-derivatives/optimizing-multivariable-functions/a/maximums-minimums-and-saddle-points
https://en.wikipedia.org/wiki/Saddle_point

https://en.wikipedia.org/wiki/Maximum_and_minimum
https://en.wikipedia.org/wiki/Second_partial_derivative_test

https://en.wikipedia.org/wiki/Partial_derivative

https://www.khanacademy.org/math/multivariable-calculus/multivariable-derivatives/partial-derivative-and-gradient-articles/a/introduction-to-partial-derivatives
https://tutorial.math.lamar.edu/classes/calciii/partialderivatives.aspx

Wait, Partial's treat the other term as a constant?  That's just going to have each side see the other's key point as zero, which is clearly incorrect.

https://en.wikipedia.org/wiki/Related_rates#Relative_kinematics_of_two_vehicles


I might have the logic backwards:

x*x = 1 + D*y*y
C = (( x*x - 1 ) / ( y*y )) = C
^^ find the largest D in the range of [1..1000] for which x is still a (local?) minimum

dc/dt = (d/dt) ( (( x*x - 1 ) / ( y*y )) )

... I don't remember how to process that and haven't found a good answer using a search engine.  This is the sort of thing that I'd crack a reference book open for (if that's part of my job, hopefully they've got a good library).

I did think about the two earlier approaches more and memorizing all of the squared values _probably_ costs more on a modern CPU than just a square root.

Some of the answers under 100 are pretty large though, so math/big to the rescue.

x*x = 1 + D*y*y
y*y = (x*x -1) / D
https://en.wikipedia.org/wiki/Derivative#Rules_of_computation

This is slow enough that I'm rather unhappy about it, but...  It's just "I don't know enough and am taking the slow way" bad.

https://en.wikipedia.org/wiki/Pell%27s_equation#Example  Mentiones my new nemesis, Continued Fractions

Didn't help me...
https://crypto.stanford.edu/pbc/notes/contfrac/pell.html


The UMB paper, Chapter 4, started to help; it at least took the detour to example 'what if D is negative'; which an instructive math text or class might offer as an exercise.  Yielding the then clear case of any values other than -1, 0, 1 for X OR Y would 'blow up' the side and make it unequal to the right.

https://www.cs.umb.edu/~eb/458/final/PeterPresentation.pdf

"""Theorem 4.4 If N = 1, then equation (13) is always solvable and the
solution are (Ak, Bk) where Ak/Bk is a convergent of √d. If N = −1,
then (13) is solvable if and only if the length of the period of the continued
expansion of √d is odd"""

However, it continues with continued fractions like they're my entire day to day job.

An interesting claim in https://www.isres.org/books/chapters/CSBET2021_10_03-01-2022.pdf calls Pell's Equation an erroneous attribution by Euler, as Fermat challenged others to solve a similar equation and Brouncker gave a solution (they reference Silverman, 2013 as the credit of that citation)

Maybe related, but lots of math major stuff  https://math.stackexchange.com/questions/2749487/algebra-direct-connect-pell-eqn-soln-p-nk-q-nk-with-p-n-q-n-sqrtd/2752161#2752161



I let the first attempt run while I did some other things, as expected it got much slower as D grew larger.

However, reading some references again (RTFM) another time, I notice how some of the examples E.G.

Example 4.3 https://www.cs.umb.edu/~eb/458/final/PeterPresentation.pdf
Sqrt(7) when D=7 for x*x - D*y*y = 2 ... it mentions that the constant number has to be less than Sqrt(D).

This also uses a square root and fractions
https://crypto.stanford.edu/pbc/notes/contfrac/pell.html

Euler 64 and 65's solutions involved computing Continued Fractions


I wonder if
dc/dt C = dc/dt ( (( x*x - 1 ) / ( y*y )) )

should really have been with the constant 1 removed

dc/dt C + 0 = dc/dt ( (( x*x ) / ( y*y )) )
dc/dt Sqrt(C) + 0 = dc/dt ( (( x ) / ( y )) )
?

I'd still like a better guide about why this becomes a square root.


After considering the previous two problems for a _while_ I realized what I was doing wrong, using the 65's fixed series of anum instead of calculating a new one each time.


*/

import (
	// "bufio"
	// "euler"
	"fmt"
	"math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler0066_slow(limit, base int64, debug bool) uint64 {
	_, _ = limit, debug
	var biggestXD, ii, evals int64
	biggestXval, X, D, Y, one := big.NewFloat(0), big.NewFloat(0), big.NewFloat(0), big.NewFloat(0), big.NewFloat(1)
	if limit < 2 {
		return 0
	}

	// x*x - D * y*y = 1
	// x*x = D * y*y + 1
	for ii = 2; ii <= limit; ii++ {
		// Euler 66 says 'it can be assumed there are no solutions in positive integers when D is square'
		iiSqr := math.Sqrt(float64(ii))
		iiS := int64(iiSqr)
		if ii == iiS*iiS {
			if debug {
				fmt.Printf("@%d:\tSKIP square (%f)\n", ii, iiSqr)
			}
			continue
		}
		D.SetInt64(ii)
		Y.SetInt64(1)
		for {
			evals++
			// // x*x = D * y*y + 1
			X.Mul(D, Y).Mul(X, Y).Add(X, one)
			X.Sqrt(X)
			// Y - D ?
			if debug || 0 == evals&0xFFFFFF {
				fmt.Printf("@%d:\tTried %d square/root solutions, best so far: X: %s\tD: %d\n", ii, evals, biggestXval.String(), biggestXD)
			}
			if X.IsInt() {
				if -1 == biggestXval.Cmp(X) {
					biggestXval.Set(X)
					biggestXD = ii
				}
				break
			}
			Y.Add(Y, one)
		}
	}
	return uint64(biggestXD)
}

/*

@2:     Tried 1 square/root solutions, best so far: X: 0        D: 0
@2:     Tried 2 square/root solutions, best so far: X: 0        D: 0
@3:     Tried 3 square/root solutions, best so far: X: 3        D: 2
@4:     SKIP square (2.000000)
@5:     Tried 4 square/root solutions, best so far: X: 3        D: 2
@5:     Tried 5 square/root solutions, best so far: X: 3        D: 2
@5:     Tried 6 square/root solutions, best so far: X: 3        D: 2
@5:     Tried 7 square/root solutions, best so far: X: 3        D: 2
@6:     Tried 8 square/root solutions, best so far: X: 9        D: 5
@6:     Tried 9 square/root solutions, best so far: X: 9        D: 5
@7:     Tried 10 square/root solutions, best so far: X: 9       D: 5
@7:     Tried 11 square/root solutions, best so far: X: 9       D: 5
@7:     Tried 12 square/root solutions, best so far: X: 9       D: 5
Euler 66: Diophantine Equation: TEST 5 true
...
@998:   Tried 4815060992 WRONG X: 2178548422      D: 778 WRONG
I kept this value in as a reminder and also as a brute force evaluation that might be correct if the approach used in the faster test happens to not pick the lowest possible value.  That chance is small, it's more likely the approach in Euler0066_slow is incorrect or misses some solution criteria.
Euler 66: WRONG: 778

real    90m54.444s
user    93m18.238s
sys     1m4.601s

*/

/*

A#	N	D	anum?
-1	1	0	.
0	2	1	.
1	3	1	?
2	8	3	?
3	11	4	?
4	19	7	?
5	87	32	?
6	106	39	?
7	193	71	?
8	1264	465	?
9	1457	536	?

anum == current A digit
Nn = N(n-2) + c * N(n-1)
Dn = D(n-2) + c * D(n-1)

*/

func EulerContinuedFracNextRot(pn0, pd, p1 *big.Int, anum int64) (*big.Int, *big.Int, *big.Int) {
	pn0.SetInt64(anum)
	pn0.Mul(pn0, p1)
	pn0.Add(pn0, pd)
	return pd, p1, pn0
}

func Euler0066(limit, base int64, debug bool) uint64 {
	_, _ = limit, debug
	if limit < 2 {
		return 0
	}
	var biggestXD, iiD, anum, evals, dSave, nSave int64
	// biggestXval, X, D, Y, one := big.NewFloat(0), big.NewFloat(0), big.NewFloat(0), big.NewFloat(0), big.NewFloat(1)
	biggestXval, X, D, Y, one := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1)
	n2, n1, n0, d2, d1, d0 := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)

	// x*x - D * y*y = 1
	// x*x = D * y*y + 1
	for iiD = 2; iiD <= limit; iiD++ {
		// Euler 66 says 'it can be assumed there are no solutions in positive integers when D is square'
		iiSqr := math.Sqrt(float64(iiD))
		iiS := int64(iiSqr)
		if iiD == iiS*iiS {
			if debug {
				fmt.Printf("@%d:\tSKIP square (%f)\n", iiD, iiSqr)
			}
			continue
		}
		D.SetInt64(iiD)
		// initial continued fraction as state 0
		// A#	N	  D	CF
		// -1	1	  0	(2)
		// 0	Fl(Sqrt)  1	(1)
		d1.SetInt64(0)
		d0.SetInt64(1)
		n1.SetInt64(1)
		n0.SetInt64(iiS)
		anum, nSave, dSave = iiS, 0, 1
		for {
			evals++

			// Calculate the next Anum term
			nSave = anum*dSave - nSave
			dSave = (iiD - nSave*nSave) / dSave
			if 0 == dSave {
				break
			}
			// Form 2
			// anum = int64((iiSqr + (float64(nSave))) / float64(dSave))
			anum = (iiS + nSave) / dSave

			n2, n1, n0 = EulerContinuedFracNextRot(n2, n1, n0, anum)
			d2, d1, d0 = EulerContinuedFracNextRot(d2, d1, d0, anum)

			if debug {
				fmt.Printf("%d: a= %d\tnSave = %d / %d = dSave,\tOUT n= %s\td= %s\n", iiD, anum, nSave, dSave, n0.String(), d0.String())
			}

			// // x*x = D * y*y + 1
			X.Mul(n0, n0)
			Y.Mul(d0, d0).Mul(Y, D).Add(Y, one)
			// Y - D ?
			if debug || 0 == evals&0xFFFFFF {
				fmt.Printf("@%d:\tTried %d square/root solutions, best so far: X: %s\tD: %d\n", iiD, evals, biggestXval.String(), biggestXD)
			}
			// X= x*x ; Y= y*y*D+1
			if 0 == X.Cmp(Y) {
				if -1 == biggestXval.Cmp(n0) {
					biggestXval.Set(n0)
					biggestXD = iiD
					fmt.Printf("@%d:\tTried %d FOUND best so far: X: %s\tD: %d\n", iiD, evals, biggestXval.String(), biggestXD)
				}
				break
			}
		}
	}
	return uint64(biggestXD)
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 66 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

# Diophantine Equation

2: a= 2 nSave = 1 / 1 = dSave,  OUT n= 3        d= 2
@2:     Tried 1 square/root solutions, best so far: X: 0        D: 0
@2:     Tried 1 FOUND best so far: X: 3 D: 2
3: a= 1 nSave = 1 / 2 = dSave,  OUT n= 2        d= 1
@3:     Tried 2 square/root solutions, best so far: X: 3        D: 2
@4:     SKIP square (2.000000)
5: a= 4 nSave = 2 / 1 = dSave,  OUT n= 9        d= 4
@5:     Tried 3 square/root solutions, best so far: X: 3        D: 2
@5:     Tried 3 FOUND best so far: X: 9 D: 5
6: a= 2 nSave = 2 / 2 = dSave,  OUT n= 5        d= 2
@6:     Tried 4 square/root solutions, best so far: X: 9        D: 5
7: a= 1 nSave = 2 / 3 = dSave,  OUT n= 3        d= 1
@7:     Tried 5 square/root solutions, best so far: X: 9        D: 5
7: a= 1 nSave = 1 / 2 = dSave,  OUT n= 5        d= 2
@7:     Tried 6 square/root solutions, best so far: X: 9        D: 5
7: a= 1 nSave = 1 / 3 = dSave,  OUT n= 8        d= 3
@7:     Tried 7 square/root solutions, best so far: X: 9        D: 5
Euler 66: Diophantine Equation: TEST 5 true
@2:     Tried 1 FOUND best so far: X: 3 D: 2
@5:     Tried 3 FOUND best so far: X: 9 D: 5
@10:    Tried 9 FOUND best so far: X: 19        D: 10
@13:    Tried 20 FOUND best so far: X: 649      D: 13
@29:    Tried 60 FOUND best so far: X: 9801     D: 29
@46:    Tried 120 FOUND best so far: X: 24335   D: 46
@53:    Tried 140 FOUND best so far: X: 66249   D: 53
@61:    Tried 196 FOUND best so far: X: 1766319049      D: 61
@109:   Tried 459 FOUND best so far: X: 158070671986249 D: 109
@181:   Tried 970 FOUND best so far: X: 2469645423824185801     D: 181
@277:   Tried 1763 FOUND best so far: X: 159150073798980475849  D: 277
@397:   Tried 2932 FOUND best so far: X: 838721786045180184649  D: 397
@409:   Tried 3009 FOUND best so far: X: 25052977273092427986049        D: 409
@421:   Tried 3187 FOUND best so far: X: 3879474045914926879468217167061449     D: 421
@541:   Tried 4610 FOUND best so far: X: 3707453360023867028800645599667005001  D: 541
@661:   Tried 6326 FOUND best so far: X: 16421658242965910275055840472270471049 D: 661
Euler 66: Diophantine Equation: 661

real    0m0.126s
user    0m0.191s
sys     0m0.040s
*/
func main() {
	var a uint64
	//test
	a = Euler0066(7, 10, true)
	fmt.Printf("Euler 66: Diophantine Equation: TEST %d %t\n", a, a == 5)

	//run
	a = Euler0066(1000, 10, false)
	fmt.Printf("Euler 66: Diophantine Equation: %d\n", a)
}
