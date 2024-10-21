// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=10
https://projecteuler.net/minimal=10

<p>The sum of the primes below $10$ is $2 + 3 + 5 + 7 = 17$.</p>
<p>Find the sum of all the primes below two million.</p>



*/
/*

Revisit; I replaced the older naive version with an interface that works OK for a small range, but kind of fails hard with re-allocations on HUGE chunks.


*/

import (
	"euler"
	"fmt"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler010(lim uint64) uint64 {
	if 2 > lim {
		return 0
	}
	if 2 == lim {
		return 2
	}

	if 500000 < lim {
		fmt.Printf("Euler 010... This might take a while, finding primes to a bit over %d\n", lim)
	}
	euler.Primes.Grow(lim)
	if 500000 < lim {
		fmt.Println("Euler 010... Time to scan and sum")
	}

	// 2 is the first prime, and all evens are compressed out
	var ret, pg, pidx, bidx, prime uint64
	ret = 2

	for {
		if euler.BVbitsPerByte <= bidx {
			bidx = 0
			pidx++
		}
		if euler.BVpagesize <= pidx {
			pidx = 0
			pg++
		}
		for ; bidx < euler.BVbitsPerByte; bidx++ {
			if 0 == euler.Primes.PV[pg][pidx]&(uint8(1)<<bidx) {
				prime = ((pg*euler.BVpagesize + pidx) << euler.BVprimeByteBitShift) + uint64(bidx)<<1 + 3
				if lim < prime {
					return ret
				}
				ret += prime
			}
		}
	}
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in $(seq 1 42) ; do go fmt $(printf "pe_%04d.go" "$ii") ; printf "NEXT: %s\n" $ii ; go run $(printf "pe_%04d.go" "$ii") || break ; read JUNK ; done

Euler010:        17 true
Euler010... This might take a while, finding primes to a bit over 2000000
Euler010... Time to scan and sum
Euler010:        142913828922



*/

func main() {
	//test
	fmt.Println("Euler010:\t", Euler010(10), Euler010(10) == 17)
	//run
	fmt.Println("Euler010:\t", Euler010(2000000))
}
