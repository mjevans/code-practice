// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=69
https://projecteuler.net/minimal=69

<p>Euler's totient function, $\phi(n)$ [sometimes called the phi function], is defined as the number of positive integers not exceeding $n$ which are relatively prime to $n$. For example, as $1$, $2$, $4$, $5$, $7$, and $8$, are all less than or equal to nine and relatively prime to nine, $\phi(9)=6$.</p>
<div class="center">
<table class="grid center"><tr><td><b>$n$</b></td>
<td><b>Relatively Prime</b></td>
<td><b>$\phi(n)$</b></td>
<td><b>$n/\phi(n)$</b></td>
</tr><tr><td>2</td>
<td>1</td>
<td>1</td>
<td>2</td>
</tr><tr><td>3</td>
<td>1,2</td>
<td>2</td>
<td>1.5</td>
</tr><tr><td>4</td>
<td>1,3</td>
<td>2</td>
<td>2</td>
</tr><tr><td>5</td>
<td>1,2,3,4</td>
<td>4</td>
<td>1.25</td>
</tr><tr><td>6</td>
<td>1,5</td>
<td>2</td>
<td>3</td>
</tr><tr><td>7</td>
<td>1,2,3,4,5,6</td>
<td>6</td>
<td>1.1666...</td>
</tr><tr><td>8</td>
<td>1,3,5,7</td>
<td>4</td>
<td>2</td>
</tr><tr><td>9</td>
<td>1,2,4,5,7,8</td>
<td>6</td>
<td>1.5</td>
</tr><tr><td>10</td>
<td>1,3,7,9</td>
<td>4</td>
<td>2.5</td>
</tr></table></div>
<p>It can be seen that $n = 6$ produces a maximum $n/\phi(n)$ for $n\leq 10$.</p>
<p>Find the value of $n\leq 1\,000\,000$ for which $n/\phi(n)$ is a maximum.</p>

*/
/*

https://brilliant.org/wiki/eulers-totient-function/
https://mathworld.wolfram.com/TotientFunction.html

https://en.wikipedia.org/wiki/Euler%27s_totient_function
"In number theory, Euler's totient function counts the positive integers up to a given integer n that are relatively prime to n. It is written using the Greek letter phi as φ ( n ) {\displaystyle \varphi (n)} or ϕ ( n ) {\displaystyle \phi (n)}, and may also be called Euler's phi function. In other words, it is the number of integers k in the range 1 ≤ k ≤ n for which the greatest common divisor gcd(n, k) is equal to 1."

https://en.wikipedia.org/wiki/Euler%27s_totient_function#Proof_of_Euler's_product_formula

Both of these require that N is factorized.

I like the pen and paper and human simplicity of:

N * (1 - 1/fact1) * ... * (1 - 1/factF)

It conveys how over the range of N any multiple of that factor is sieved out as a number that shares a 'base' system with N (thus it is not 'relatively (to a base sysetm) prime').

[]FactorizedN -> b0^(p0-1) * (b0-1) * ... * bn^(pn-1) * (bn-1)

The alternative turns that question on it's head and does away with fiddly fractions.

Rather than exclude numbers by division, subtract the repetitions of each multiple from the range covered by each prime base.  I don't quite visualize where the minus one power originates, but it works.

Damn, the range N <= 1_000_000 means the uint16s my Factorized class depends on no longer work.  It might be time for an overhaul to my personal Euler math library.

For the moment, any non math/big related whole / unit number should probably go to either uint64 or int64, and the internal 'small' values like for primes and exponents can go to 32 bit values.  It seems super unlikely that a Euler problem will ask to work with primes greater than 32 bits.



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

func Euler0069(toMax uint64) uint64 {
	var bestI, ii, phi, bestFI, fi uint64
	var bestF, flNphi float64
	euler.Primes.Grow(euler.SqrtU64(toMax))
	for ii = 1; ii <= toMax; ii++ {
		phi = euler.EulerTotientPhi(ii, 0)
		fi = ii / phi
		if fi >= bestFI {
			flNphi = float64(ii) / float64(phi)
			if flNphi > bestF {
				fmt.Printf("Found new best N/phi: %d/%d ~= %f > %f\n", ii, phi, flNphi, bestF)
				bestF, bestFI, bestI = flNphi, fi, ii
			}
		}
	}
	return bestI
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 69 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Found new best N/phi: 1/1 ~= 1.000000 > 0.000000
Found new best N/phi: 2/1 ~= 2.000000 > 1.000000
Found new best N/phi: 6/2 ~= 3.000000 > 2.000000
Euler 69: Totient Maximum: TEST 6 true
Found new best N/phi: 1/1 ~= 1.000000 > 0.000000
Found new best N/phi: 2/1 ~= 2.000000 > 1.000000
Found new best N/phi: 6/2 ~= 3.000000 > 2.000000
Found new best N/phi: 30/8 ~= 3.750000 > 3.000000
Found new best N/phi: 210/48 ~= 4.375000 > 3.750000
Found new best N/phi: 2310/480 ~= 4.812500 > 4.375000
Found new best N/phi: 30030/5760 ~= 5.213542 > 4.812500
Found new best N/phi: 510510/92160 ~= 5.539388 > 5.213542
Euler 69: Totient Maximum: 510510
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	test := Euler0069(10)
	fmt.Printf("Euler 69: Totient Maximum: TEST %d %t\n", test, test == 6)

	//run
	fmt.Printf("Euler 69: Totient Maximum: %d\n", Euler0069(1_000_000))
}
