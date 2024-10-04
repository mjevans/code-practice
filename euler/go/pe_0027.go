// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=27
https://projecteuler.net/minimal=27


<p>Euler discovered the remarkable quadratic formula:</p>
<p class="center">$n^2 + n + 41$</p>
<p>It turns out that the formula will produce $40$ primes for the consecutive integer values $0 \le n \le 39$. However, when $n = 40, 40^2 + 40 + 41 = 40(40 + 1) + 41$ is divisible by $41$, and certainly when $n = 41, 41^2 + 41 + 41$ is clearly divisible by $41$.</p>
<p>The incredible formula $n^2 - 79n + 1601$ was discovered, which produces $80$ primes for the consecutive values $0 \le n \le 79$. The product of the coefficients, $-79$ and $1601$, is $-126479$.</p>
<p>Considering quadratics of the form:</p>
<blockquote>
$n^2 + an + b$, where $|a| &lt; 1000$ and $|b| \le 1000$<br><br><div>where $|n|$ is the modulus/absolute value of $n$<br>e.g. $|11| = 11$ and $|-4| = 4$</div>
</blockquote>
<p>Find the product of the coefficients, $a$ and $b$, for the quadratic expression that produces the maximum number of primes for consecutive values of $n$, starting with $n = 0$.</p>

*/
/*

Writing this in my terms...

// Valid primes 0..n for biggest n
n*n + a*n + b

-1000 <= a <= 1000
-1000 <= b <= 1000
0 <= n <= (until condition violated)

At first I considered a cache on things like n*n, but modern CPUs probably multiply faster than the cache miss on lookups; and many numbers are going to fail early anyway.
1000*1000 + 1000*1000 + 1000 is 2 million, however I seriously doubt a valid prime will need to be _that_ large...

80 was considered incredible so 200*200 + 1000*200 + 1000 ~= 40000 + 200000 + 1000 = 241000, eh prisoner 24601_9

... Revisit updates

This spends so much time asking

*/

import (
	// "bufio"
	// "bitvector"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler027(maxa, maxb, maxn, expect int) (int, int, int) {
	// lim := maxn*maxn + maxn*maxa + maxb

	if maxn <= 0 {
		return 0, 0, 0
	}

	fmt.Printf("Euler 027 Primes.Grow(%d)\n", uint(maxn*maxn+maxn*maxa+maxb))
	euler.Primes.Grow(uint(maxn*maxn + maxn*maxa + maxb))

	var besta, bestb, bestn int
	bestn = expect

	for a := -maxa; a <= maxa; a++ {
		for b := -maxb; b <= maxb; b++ {
			for n := 0; n <= maxn; n++ {
				test := n*n + n*a + b
				if false == euler.Primes.KnownPrime(uint(test)) {
					n--
					if n > bestn {
						besta, bestb, bestn = a, b, n
						fmt.Printf("\nEuler027: NEW\t%d\tn*n + %d*n + %d\n", n, a, b)
					} else if n == bestn {
						// fmt.Printf("\nEuler027: tie\t%d\tn*n + %d*n + %d\n", n, a, b)
					}
					break // n
				}
			}
		}
	}
	return besta, bestb, bestn
}

/*
Euler 027 Primes.Grow(1722)

Euler027: NEW   10      n*n + -1*n + 11

Euler027: NEW   16      n*n + -1*n + 17

Euler027: NEW   39      n*n + 1*n + 41
Euler027: test:  true   0.. 39  FOR n*n +  1 *n +  41
Euler 027 Primes.Grow(111000)

Euler027: NEW   10      n*n + -105*n + 967

Euler027: NEW   70      n*n + -61*n + 971
Euler027: run:          0.. 70  FOR n*n +  -61 *n +  971        Thus     -59231
*/
func main() {
	//test
	a, b, n := Euler027(2, 42, 40, 9)
	fmt.Println("Euler027: test:\t", 1 == a && 41 == b && 39 == n, "\t0..", n, " FOR n*n + ", a, "*n + ", b)

	//run
	a, b, n = Euler027(1000, 1000, 100, 9)
	fmt.Println("Euler027: run:\t", "\t0..", n, " FOR n*n + ", a, "*n + ", b, "\tThus\t", a * b)

}
