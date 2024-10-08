// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=46
https://projecteuler.net/minimal=46

<p>It was proposed by Christian Goldbach that every odd composite number can be written as the sum of a prime and twice a square.</p>
\begin{align}
9 = 7 + 2 \times 1^2\\
15 = 7 + 2 \times 2^2\\
21 = 3 + 2 \times 3^2\\
25 = 7 + 2 \times 3^2\\
27 = 19 + 2 \times 2^2\\
33 = 31 + 2 \times 1^2
\end{align}
<p>It turns out that the conjecture was false.</p>
<p>What is the smallest odd composite that cannot be written as the sum of a prime and twice a square?</p>



*/
/*

sqrt ( 500_000 / 2 ) = 500 OK, just go for it

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

func Euler046(limit uint64) uint64 {
	var ii, jj, check uint64
	iiEnd := uint64(math.Sqrt(float64(limit>>1))) + 2
	dsq := make([]uint64, 0, iiEnd)
	for ii = 0; ii < iiEnd; ii++ {
		dsq = append(dsq, (ii*ii)<<1)
	}
	euler.Primes.Grow(uint(limit))
	// 9 is the smallest odd composite number
Euler046composites:
	for ii = 9; ii < limit; ii += 2 {
		if euler.Primes.KnownPrime(uint(ii)) {
			continue
		}
		jj = 0
		for {
			jj++
			if ii <= dsq[jj] {
				fmt.Printf("Unable to find composite solution for %d ( < %d [%d])\n", ii, dsq[jj], jj)
				return ii
			}
			check = ii - dsq[jj]
			if euler.Primes.KnownPrime(uint(check)) {
				continue Euler046composites
			}
		}
	}
	return 0
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 46 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Unable to find composite solution for 5777 ( < 5832 [54])
Euler 46: Goldbach's Other Conjecture: 5777

*/
func main() {
	//test

	//run
	a := Euler046(500_000)
	fmt.Printf("Euler 46: Goldbach's Other Conjecture: %d\n", a)
}
