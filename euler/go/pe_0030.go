// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=30
https://projecteuler.net/minimal=30


<p>Surprisingly there are only three numbers that can be written as the sum of fourth powers of their digits:
\begin{align}
1634 &amp;= 1^4 + 6^4 + 3^4 + 4^4\\
8208 &amp;= 8^4 + 2^4 + 0^4 + 8^4\\
9474 &amp;= 9^4 + 4^4 + 7^4 + 4^4
\end{align}
</p><p class="smaller">As $1 = 1^4$ is not a sum it is not included.</p>
<p>The sum of these numbers is $1634 + 8208 + 9474 = 19316$.</p>
<p>Find the sum of all the numbers that can be written as the sum of fifth powers of their digits.</p>


Is there a property by which I can determine the maximum possible register size of a valid solution
Empirical observation (9^10 = 3486784401 ) suggests the 9^n might match at most a register size of N digits.

100000f+10000e+1000d+100c+10b+1a = f^5+e^5+d^5+c^5+b^5+a^5



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

func Euler030() (uint, []uint) {
	power := uint(5)
	ret, pow, hund := make([]uint, 0, 8), [10]uint{0, 1}, [100]uint{}
	var sum uint
	for ii := uint(2); ii < 10; ii++ {
		t := ii
		for pp := uint(1); pp < power; pp++ {
			t *= ii
		}
		pow[ii] = t
	}
	for a := uint(0); a < 10; a++ {
		a10 := a * 10
		for b := uint(0); b < 10; b++ {
			hund[b+a10] = pow[b] + pow[a]
		}
	}
	var floor uint // Maybe a large place int can roll the lower digits all the way around, so this _was_ incorrect.  It shadowed results with 9 in the 10000s digit.
	for a := uint(0); a < 100; a++ {
		// floor := (hund[a] / 100) % 100
		for b := floor; b < 100; b++ {
			// floor = ((hund[a] + hund[b]) / 10000) % 100
			for c := floor; c < 100; c++ {
				powSum := hund[c] + hund[b] + hund[a]
				if 1 < powSum && 10000*c+100*b+1*a == powSum {
					fmt.Printf("Found: %d %d %d ~~ %d %d %d == %d\n", c, b, a, hund[c], hund[b], hund[a], powSum)
					ret = append(ret, powSum)
					sum += powSum
				}

			}
		}
	}
	return sum, ret
}

//	for ii in */*.go ; do go fmt "$ii" ; done ; for ii in 30 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
/*

Found: 5 47 48 ~~ 3125 17831 33792 == 54748
Found: 0 41 50 ~~ 0 1025 3125 == 4150
Found: 0 41 51 ~~ 0 1025 3126 == 4151
Euler030: Result  63049 [54748 4150 4151]

Found: 9 27 27 ~~ 59049 16839 16839 == 92727
Found: 5 47 48 ~~ 3125 17831 33792 == 54748
Found: 0 41 50 ~~ 0 1025 3125 == 4150
Found: 0 41 51 ~~ 0 1025 3126 == 4151
Found: 19 49 79 ~~ 59050 60073 75856 == 194979
Found: 9 30 84 ~~ 59049 243 33792 == 93084
Euler030: Result  443839 [92727 54748 4150 4151 194979 93084]



*/
func main() {
	//test

	//run
	sum, r := Euler030()
	fmt.Println("Euler030: Result ", sum, r)

}
