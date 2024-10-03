// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=40
https://projecteuler.net/minimal=40

Champernowne's Constant

<p>An irrational decimal fraction is created by concatenating the positive integers:
$$0.12345678910{\color{red}\mathbf 1}112131415161718192021\cdots$$</p>
<p>It can be seen that the $12$<sup>th</sup> digit of the fractional part is $1$.</p>
<p>If $d_n$ represents the $n$<sup>th</sup> digit of the fractional part, find the value of the following expression.
$$d_1 \times d_{10} \times d_{100} \times d_{1000} \times d_{10000} \times d_{100000} \times d_{1000000}$$</p>



*/
/*
                           * 10s place in 11 == digit 12
 0 . 1 2 3 4 5 6 7 8 9  10 11 12 13 14 15 16 17 ...
[ 10 - 1 = 1 * 9      ][ 100 - 10 = 2 * 90          ][ 1000 - 100 = 3 * 900          ]

Digit 12 ?

12 > (span 10^0=1* :: 9)
12 - 9 = 3
3 < (span 10^1=10* :: 90)
3 digits into span 90, 10^n ~ n + 1 digits == 2 per number

digit(1) * ... ~ {1, 10, 100, 1000, 10000, 100000, 1000000 }

Also, never having heard the name https://en.wikipedia.org/wiki/Champernowne_constant

What is this supposed to be zero because one of these digits is zero?

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

func ChampernowneDigit(ii uint64) uint64 {
	var dig, placeNext, place uint64
	if 1 == ii {
		return 1
	}

	// Correct for 0th offset.
	ii--
	dig = 1
	place = 1
	placeNext = 10
	for ii >= placeNext-place {
		ii -= (placeNext - place)
		place = placeNext
		placeNext *= 10
		dig += 1
	}
	// digit power is backwards since Arabic numbers are written largest to smallest, but string index is smallest to largest.
	numPlus, numPow := place+ii/dig, dig-(ii%dig)-1
	for 0 < numPow {
		numPlus /= 10
		numPow--
	}
	return numPlus % 10
}

func Euler040() uint64 {
	seq := [...]uint64{1, 10, 100, 1000, 10000, 100000, 1000000}
	ret := uint64(1)
	iiLen := len(seq)
	for ii := 0; ii < iiLen; ii++ {
		ret *= ChampernowneDigit(seq[ii])
		fmt.Printf("\t%d:\t%d\t%d", ii, seq[ii], ret)
	}
	return ret
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 40 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

        0:      1       1       1:      10      1       2:      100     1       3:      1000    1       4:      10000   1       5:      100000  1       6:      1000000 1
Euler040: Champernowne's Constant (result) :    1




*/
func main() {
	//test
	test := []struct {
		test, res uint64
	}{
		{1, 1},
		{9, 9},
		{10, 1},
		{11, 0},
		{12, 1},
		{13, 1},
		{14, 1},
		{15, 2},
	}
	for _, test := range test {
		res := ChampernowneDigit(test.test)
		if test.res != res {
			fmt.Printf("%d Expected %d got %d\n", test.test, test.res, res)
		}
	}

	//run
	r := Euler040()
	fmt.Printf("\nEuler040: Champernowne's Constant (result) :\t%d\n", r)

}
