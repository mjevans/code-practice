// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=77
https://projecteuler.net/minimal=77

<p>It is possible to write ten as the sum of primes in exactly five different ways:</p>
\begin{align}
&amp;7 + 3\\
&amp;5 + 5\\
&amp;5 + 3 + 2\\
&amp;3 + 3 + 2 + 2\\
&amp;2 + 2 + 2 + 2 + 2
\end{align}
<p>What is the first value which can be written as the sum of primes in over five thousand different ways?</p>



/
*/
/*
/	Reuse the partition system from 76
/
*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

var PrimesSmallU8 []uint8

func Euler0077(req uint64) uint64 {
	//coins := make([]uint8, 0, total)
	coins := PrimesSmallU8
	cache := make(map[uint32]uint64)

	var CachableCombos func(maxcoin, total int16) uint64
	CachableCombos = func(maxcoin, total int16) uint64 {
		if 0 > total || 0 == maxcoin {
			return 0
		}
		if 0 == total {
			return 1
		}
		var val uint64
		key := uint32(maxcoin)<<16 | uint32(total)
		if val, exists := cache[key]; exists {
			return val
		}
		val = CachableCombos(maxcoin-1, total) + CachableCombos(maxcoin, total-int16(coins[maxcoin-1]))
		cache[key] = val
		return val
	}

	//coincount := int16(len(coins))
	coincount := int16(20)
	for total := int16(0); total < (1<<15)-1; total++ {
		res := CachableCombos(coincount, total)
		if res >= req {
			fmt.Printf("%d >= %d @ %d\n", res, req, total)
			return uint64(total)
		}
	}
	return 0
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 77 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

5 >= 5 @ 10
5007 >= 5000 @ 71
Euler 77: :     Count: 71

real    0m0.098s
user    0m0.159s
sys     0m0.031s
.
*/
func main() {
	PrimesSmallU8 = []uint8{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251} // 41 required for reasons, 53 nice for 16 total numbers; 16 bytes of memory, 1/4th cache line

	//test
	// tested in the golang tests for "euler"
	r := Euler0077(5)
	if 10 != r {
		panic(fmt.Sprintf("Euler 77: Expected 10 got %d", r))
	}

	//run
	r = Euler0077(5000)
	fmt.Printf("Euler 77: :\tCount: %d\n", r)
	if 71 != r {
		panic("Did not reach expected value.")
	}
}
