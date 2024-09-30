// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=35

<p>The number, $197$, is called a circular prime because all rotations of the digits: $197$, $971$, and $719$, are themselves prime.</p>
<p>There are thirteen such primes below $100$: $2, 3, 5, 7, 11,		13, 17, 31, 37, 71, 	73, 79$, and $97$.</p>
<p>How many circular primes are there below one million?</p>







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

func Euler035(limit uint) uint {
	if 0 == limit {
		limit = 1_000_000
	}
	p := euler.NewBVPrimes() // Local vector since it'll be mutated
	p.Grow(limit)
	var found, prime uint
	if prime < 2 {
		found++
		prime = 2
	}
	for prime < limit {
		prime = p.PrimeAfter(prime)
		rots := euler.RotateDecDigits(uint64(prime))
		ok := true
		seen := make(map[uint64]uint64)
		iiLim := len(rots)
		for ii := 0; ii < iiLim; ii++ {
			if _, old := seen[rots[ii]]; false == old {
				seen[rots[ii]] = rots[ii]
				ok = ok && p.KnownPrime(uint(rots[ii]))
				if ok {
					bbAbs := (uint(rots[ii]) - 3) >> 1
					ooAbs := (uint(rots[ii]) - 3) >> (euler.BVprimeByteBitShift)
					pg, ooPos, bbBit := ooAbs/(euler.BVpagesize), ooAbs%(euler.BVpagesize), bbAbs&euler.BVprimeByteBitMaskPost
					// fmt.Printf("[%d][%d]:%d <= %d\n", pg, ooPos, bbBit, rots[ii])
					p.PV[pg][ooPos] |= uint8(1) << bbBit
				}
			}
		}
		if ok {
			found += uint(len(seen))
			fmt.Printf("%d\t+++\t%v\n", found, rots)
		}
		seen = nil

	}
	return found
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 35 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done



2       +++     [3]
3       +++     [5]
4       +++     [7]
5       +++     [11 11]
7       +++     [13 31]
9       +++     [17 71]
11      +++     [37 73]
13      +++     [79 97]
Euler035: TEST: 100 == 13?      13
2       +++     [3]
3       +++     [5]
4       +++     [7]
5       +++     [11 11]
7       +++     [13 31]
9       +++     [17 71]
11      +++     [37 73]
13      +++     [79 97]
16      +++     [113 311 131]
19      +++     [197 719 971]
22      +++     [199 919 991]
25      +++     [337 733 373]
29      +++     [1193 3119 9311 1931]
33      +++     [3779 9377 7937 7793]
38      +++     [11939 91193 39119 93911 19391]
43      +++     [19937 71993 37199 93719 99371]
49      +++     [193939 919393 391939 939193 393919 939391]
55      +++     [199933 319993 331999 933199 993319 999331]
Euler035: RUN : 1_000_000 ==    55

*/
func main() {
	//test
	count := Euler035(100)
	fmt.Printf("Euler035: TEST: 100 == 13?\t%d\n", count)

	//run
	count = Euler035(1_000_000)
	fmt.Printf("Euler035: RUN : 1_000_000 ==\t%d\n", count)

}
