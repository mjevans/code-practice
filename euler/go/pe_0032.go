// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=32
https://projecteuler.net/minimal=32


<p>We shall say that an $n$-digit number is pandigital if it makes use of all the digits $1$ to $n$ exactly once; for example, the $5$-digit number, $15234$, is $1$ through $5$ pandigital.</p>

<p>The product $7254$ is unusual, as the identity, $39 \times 186 = 7254$, containing multiplicand, multiplier, and product is $1$ through $9$ pandigital.</p>

<p>Find the sum of all products whose multiplicand/multiplier/product identity can be written as a $1$ through $9$ pandigital.</p>

<div class="note">HINT: Some products can be obtained in more than one way so be sure to only include it once in your sum.</div>

*/
/*


Zero (0) can't be used.

Balances:

NOPE aaa * bbb = ccc ?  // Nope; 100*100 is 10000
aa * bbb = cccc // At least one example given
NOPE aa * bb = ccccc ?? 100 * 100 is 10000, but that's over the limits and the min result
aa * bbbb = ccc , clearly no

+ Revisit +

I missed an option...

a * bbbb = cccc // This can work too!

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

func Euler032() (uint, uint) {
	var combos, sum, a, b, r uint
	var used uint16
	seen := make(map[uint]bool)
Euler032a:
	for a = 123; a <= 987; a++ {
		// func Pandigital(test uint64, used uint16) (highest, usedRe uint16, DigitShift uint64) {
		/*
			used = 0
			test := a
			for test > 0 {
				bd := uint(1) << (test % 10)
				test = test / 10
				if 0 < used&bd || bd == 1 {
					// fmt.Printf("Skip dup digit: %10b %d\n", used, a)
					continue Euler032a
				}
				used |= bd
			}
			used_a := used
		*/
		bts, used_a, _ := euler.Pandigital(uint64(a), 0)
		if 0 == bts {
			continue Euler032a
		}
	Euler032b:
		for b = 12; b <= 98; b++ {
			/*
				used = used_a
				test = b
				for test > 0 {
					bd := uint(1) << (test % 10)
					test = test / 10
					if 0 < used&bd || bd == 1 {
						// fmt.Printf("Skip dup digit: %10b %d x %d\n", used, a, b)
						continue Euler032b
					}
					used |= bd
				}
			*/
			bts, used, _ = euler.Pandigital(uint64(b), used_a)
			if 0 == bts {
				continue Euler032b
			}
			r = a * b
			if r >= 100000 {
				fmt.Printf("Skip too big: %d\n", r)
				continue Euler032a
			}
			/*
				test = r
				for test > 0 {
					bd := uint(1) << (test % 10)
					test = test / 10
					if 0 < used&bd || bd == 1 {
						// fmt.Printf("Skip dup digit: %10b %d x %d = %d\n", used, a, b, r)
						continue Euler032b
					}
					used |= bd
				}
			*/
			bts, used, _ = euler.Pandigital(uint64(r), used)
			if 9 != bts {
				continue Euler032b
			}
			if _, ex := seen[r]; false == ex && 0b0000_0011_1111_1110 == used {
				combos++
				sum += r
				fmt.Printf("Euler032: seen\t%d x %d = ** %d **\n", a, b, r)
				seen[r] = true
			} else {
				fmt.Printf("Euler032: SKIPPED\t%d x %d = %d\n", a, b, r)
			}
		}
	}
	for a = 1234; a <= 9876; a++ {
		// func Pandigital(test uint64, used uint16) (highest, usedRe uint16, DigitShift uint64) {
		bts, used_a, _ := euler.Pandigital(uint64(a), 0)
		if 0 == bts {
			continue
		}
		for b = 1; b <= 9; b++ {
			bts, used, _ = euler.Pandigital(uint64(b), used_a)
			if 0 == bts {
				continue
			}
			r = a * b
			if r >= 100000 {
				fmt.Printf("Skip too big: %d\n", r)
				continue
			}
			bts, used, _ = euler.Pandigital(uint64(r), used)
			if 9 != bts {
				continue
			}
			if _, ex := seen[r]; false == ex && 0b0000_0011_1111_1110 == used {
				combos++
				sum += r
				fmt.Printf("Euler032: seen\t%d x %d = ** %d **\n", a, b, r)
				seen[r] = true
			} else {
				fmt.Printf("Euler032: SKIPPED\t%d x %d = %d\n", a, b, r)
			}
		}
	}

	return combos, sum
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 32 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler032: seen  138 x 42 = ** 5796 **
Euler032: seen  157 x 28 = ** 4396 **
Euler032: seen  159 x 48 = ** 7632 **
Euler032: seen  186 x 39 = ** 7254 **
Euler032: seen  198 x 27 = ** 5346 **
Euler032: SKIPPED       297 x 18 = 5346
Euler032: SKIPPED       483 x 12 = 5796
Euler032: seen  1738 x 4 = ** 6952 **
Euler032: seen  1963 x 4 = ** 7852 **
Euler032:        7       45228
*/
func main() {
	//test

	//run
	combos, sum := Euler032()
	fmt.Println("Euler032:\t", combos, "\t", sum)
}
