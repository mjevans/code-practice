// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=47
https://projecteuler.net/minimal=47

<p>The first two consecutive numbers to have two distinct prime factors are:</p>
\begin{align}
14 &amp;= 2 \times 7\\
15 &amp;= 3 \times 5.
\end{align}
<p>The first three consecutive numbers to have three distinct prime factors are:</p>
\begin{align}
644 &amp;= 2^2 \times 7 \times 23\\
645 &amp;= 3 \times 5 \times 43\\
646 &amp;= 2 \times 17 \times 19.
\end{align}
<p>Find the first four consecutive integers to have four distinct prime factors each. What is the first of these numbers?</p>




*/
/*


 */

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler047(limit, run uint) uint {
	var ii, r, tt uint32
	euler.Primes.Grow(limit)
	ii = 3
Euler047ii:
	for ; uint(ii) <= limit; ii++ {
		if euler.Primes.KnownPrime(uint(ii)) {
			r = 0
			continue
		}
		r++
		if uint32(run) <= r {
			// shift to zero index
			for tt = ii - uint32(run) + 1; tt <= ii; tt++ {
				f := euler.Primes.Factorize(uint(tt))
				// fmt.Printf("\tTesting: %d (%d)\t%d\t%d\n", ii, r, tt, f.Lenbase)
				if uint(f.Lenbase) != uint(run) {
					// fmt.Printf("\tAborted run at %d\n", tt)
					continue Euler047ii
				}
			}
			ans := uint(ii) - run + 1
			fmt.Printf("Start of run of %d composite numbers with exactly %d base factors (and unlimited powers per): %d\n", run, run, ans)
			return ans
		}
	}
	return 0
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 47 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Start of run of 2 composite numbers with exactly 2 base factors (and unlimited powers per): 14
Euler 47: Test 2: 14
Start of run of 3 composite numbers with exactly 3 base factors (and unlimited powers per): 644
Euler 47: Test 3: 644
Start of run of 4 composite numbers with exactly 4 base factors (and unlimited powers per): 134043
Euler 47: Distinct Primes Factors: 134043

*/
func main() {
	//test
	fmt.Printf("Euler 47: Test 2: %d\n", Euler047(1_000, 2))
	fmt.Printf("Euler 47: Test 3: %d\n", Euler047(1_000, 3))

	//run
	a := Euler047(500_000, 4)
	fmt.Printf("Euler 47: Distinct Primes Factors: %d\n", a)
}
