// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=34
https://projecteuler.net/minimal=34


<p>$145$ is a curious number, as $1! + 4! + 5! = 1 + 24 + 120 = 145$.</p>
<p>Find the sum of all numbers which are equal to the sum of the factorial of their digits.</p>
<p class="smaller">Note: As $1! = 1$ and $2! = 2$ are not sums they are not included.</p>

9! == 362880 soo things can get big fast.
9!*6 == 2177280
9!*7 == 2540160




*/

import (
	// "bufio"
	// "bitvector"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler034() uint {
	fa := [10]uint{0, 1, 2}
	fac := uint(2)
	for f := uint(3); f < 10; f++ {
		fac *= f
		fa[f] = fac
	}
	fmt.Println(fa)
	var sum uint
	for i7 := uint(0); i7 < 3; i7++ {
		for i6 := uint(0); i6 < 10; i6++ {
			for i5 := uint(0); i5 < 10; i5++ {
				for i4 := uint(0); i4 < 10; i4++ {
					for i3 := uint(0); i3 < 10; i3++ {
						for i2 := uint(0); i2 < 10; i2++ {
							for i1 := uint(0); i1 < 10; i1++ {
								faSum := fa[i7] + fa[i6] + fa[i5] + fa[i4] + fa[i3] + fa[i2] + fa[i1]
								if 2 < faSum && i7*1000000+i6*100000+i5*10000+i4*1000+i3*100+i2*10+i1 == faSum {
									sum += faSum
									fmt.Printf("Found: %d\n", faSum)
								}
							}
						}
					}
				}
			}
			// fmt.Printf("... %d%d....\n", i7, i6)
		}
	}
	return sum
}

//	for ii in */*.go ; do go fmt "$ii" ; done ; for ii in 34 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
/*

Euler034: Sum of Digit Factorials: 145


*/
func main() {
	//test

	//run
	sum := Euler034()
	fmt.Printf("Euler034: Sum of Digit Factorials: %d\n", sum)

	fmt.Printf("\t\t\t*** TODO correct answer before new Euler problems ***\n")
}
