// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=92
https://projecteuler.net/minimal=92

<p>A number chain is created by continuously adding the square of the digits in a number to form a new number until it has been seen before.</p>
<p>For example,
\begin{align}
&amp;44 \to 32 \to 13 \to 10 \to \mathbf 1 \to \mathbf 1\\
&amp;85 \to \mathbf{89} \to 145 \to 42 \to 20 \to 4 \to 16 \to 37 \to 58 \to \mathbf{89}
\end{align}
</p><p>Therefore any chain that arrives at $1$ or $89$ will become stuck in an endless loop. What is most amazing is that EVERY starting number will eventually arrive at $1$ or $89$.</p>
<p>How many starting numbers below ten million will arrive at $89$?</p>


/
*/
/*


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

func Euler0092_10(max uint32) uint32 {
	var ret, ii, cur, next, bv89valid uint32
	caSq := [10]uint8{}
	for ii = 1; ii < 10; ii++ {
		caSq[ii] = uint8(ii * ii)
	}

	bv89mx := 9*9*9 + 4*4                // a 32 bit number can't quite hold 10 decimal digits, so this is more than enough.  If there are 4_999_999_999
	bv89 := make([]uint8, 1+(bv89mx>>3)) // BitVector: 0 (beneath current test) means ends in 1, while 1 means ends in 89

	for ii = 1; ii <= max; ii++ {
		iter := 0
		for cur = ii; 1 != cur && 89 != cur; cur = next {
			iter++
			for next = 0; 0 < cur; {
				next, cur = next+uint32(caSq[cur%10]), cur/10
			}
			if next <= bv89valid {
				if 0 < ((bv89[next>>3]) & (1 << (next & 0b111))) {
					next = 89
				} else {
					next = 1
				}
			}
			if 20 < iter {
				return 0
			}
		}
		if bv89valid < uint32(bv89mx) {
			if 89 == cur {
				bv89[ii>>3] |= (1 << (ii & 0b111))
			}
			bv89valid = ii
		}
		// fmt.Printf("\t%d\t%d\n", ii, cur)
		if 89 == cur {
			ret++
		}
	}

	return ret
}

func Euler0092(max uint32) uint32 {
	var ret, ii, cur, next, bv89valid uint32
	caSq := [100]uint8{} // two cache lines and 81+81 fits in u8
	for cur = 0; cur < 10; cur++ {
		caSq[cur*10] = uint8(cur * cur)
		for ii = 1; ii < 10; ii++ {
			caSq[cur*10+ii] = uint8(ii*ii + cur*cur)
		}
	}

	bv89mx := 9*9*9 + 4*4                // a 32 bit number can't quite hold 10 decimal digits, so this is more than enough.  If there are 4_999_999_999
	bv89 := make([]uint8, 1+(bv89mx>>3)) // BitVector: 0 (beneath current test) means ends in 1, while 1 means ends in 89

	for ii = 1; ii <= max; ii++ {
		for cur = ii; 1 != cur && 89 != cur; cur = next {
			for next = 0; 0 < cur; {
				next, cur = next+uint32(caSq[cur%100]), cur/100
			}
			if next <= bv89valid {
				if 0 < ((bv89[next>>3]) & (1 << (next & 0b111))) {
					next = 89
				} else {
					next = 1
				}
			}
		}
		if bv89valid < uint32(bv89mx) {
			if 89 == cur {
				bv89[ii>>3] |= (1 << (ii & 0b111))
			}
			bv89valid = ii
		}
		// fmt.Printf("\t%d\t%d\n", ii, cur)
		if 89 == cur {
			ret++
		}
	}

	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 92 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Easy, but slow
Euler 92: Square Digit Chains: 8581146

real    0m0.235s
user    0m0.285s
sys     0m0.053s
.
*/
func main() {
	var r uint32
	//test
	// r = Euler0092(9)
	// fmt.Printf("Euler 92: debug: %d\n", r)
	// return

	//run
	r = Euler0092(10_000_000 - 1)
	fmt.Printf("Euler 92: Square Digit Chains: %d\n", r)
	if 8581146 != r {
		panic("Did not reach expected value.")
	}
}
