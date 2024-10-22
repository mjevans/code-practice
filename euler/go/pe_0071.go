// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=71
https://projecteuler.net/minimal=71

<p>Consider the fraction, $\dfrac n d$, where $n$ and $d$ are positive integers. If $n \lt d$ and $\operatorname{HCF}(n,d)=1$, it is called a reduced proper fraction.</p>
<p>If we list the set of reduced proper fractions for $d \le 8$ in ascending order of size, we get:
$$\frac 1 8, \frac 1 7, \frac 1 6, \frac 1 5, \frac 1 4, \frac 2 7, \frac 1 3, \frac 3 8, \mathbf{\frac 2 5}, \frac 3 7, \frac 1 2, \frac 4 7, \frac 3 5, \frac 5 8, \frac 2 3, \frac 5 7, \frac 3 4, \frac 4 5, \frac 5 6, \frac 6 7, \frac 7 8$$</p>
<p>It can be seen that $\dfrac 2 5$ is the fraction immediately to the left of $\dfrac 3 7$.</p>
<p>By listing the set of reduced proper fractions for $d \le 1\,000\,000$ in ascending order of size, find the numerator of the fraction immediately to the left of $\dfrac 3 7$.</p>



*/
/*

Any valid fraction with a denominator 1..1_000_000 ; find the one closest in value to 3/7th without going equal or over

It's _probably_ going to be the largest Totient number under 1 million... or one that has a bunch of prime factors otherwise.

X/Y == N/D
XD/YD == YN/YD
XD == YN
X == YN/D


*/

import (
	// "bufio"
	"euler"
	"fmt"
	"math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0071(limit uint32) (uint32, uint32) {
	var iiD, bestN, bestD uint32
	var bestF, targetF, NF, testF float64
	targetF = float64(3) / float64(7)
	for iiD = limit; 0 < iiD; iiD-- {
		NF = math.Floor(float64(iiD) * targetF)
		testF = NF / float64(iiD)
		if testF > bestF && testF != targetF {
			bestN, bestD = uint32(NF), iiD
			// fmt.Printf("New Best found: %d/%d = %f > %f\n", bestN, bestD, testF, bestF)
			bestF = testF
		}
	}
	gcd := euler.GCDbin(bestN, bestD)
	fmt.Printf("Best found: %d/%d\n\t%1.40f gcd: %d\n\t%1.40f\n", bestN, bestD, bestF, gcd, targetF)
	return bestN / gcd, bestD / gcd
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 71 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Best found: 428570/999997

	0.4285712857138571640902569015452172607183 gcd: 1
	0.4285714285714285476380780437466455623507

Euler 71: Ordered Fractions: Left of 3/7th 428570 / 999997
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"

	//run
	n, d := Euler0071(1_000_000)
	fmt.Printf("Euler 71: Ordered Fractions: Left of 3/7th %d / %d\n", n, d)
}
