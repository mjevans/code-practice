// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=33

<p>The fraction $49/98$ is a curious fraction, as an inexperienced mathematician in attempting to simplify it may incorrectly believe that $49/98 = 4/8$, which is correct, is obtained by cancelling the $9$s.</p>
<p>We shall consider fractions like, $30/50 = 3/5$, to be trivial examples.</p>
<p>There are exactly four non-trivial examples of this type of fraction, less than one in value, and containing two digits in the numerator and denominator.</p>
<p>If the product of these four fractions is given in its lowest common terms, find the value of the denominator.</p>


Well there are 4, and they want to ignore 'easy' ones like chop off extra 0s...
* So a digit is common between top and bottom, and the remaining digits are a correct reduction (if not fully proper).
* Less than 1 (num < den) in value

49/98 = 4/8 ( = 1/2 )



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

func Euler033() *euler.Rational {
	ret := euler.NewRational(1, 1)
	for num := 11; num <= 98; num++ {
		if 0 == num%10 {
			continue
		}
		for den := num + 1; den <= 99; den++ {
			if 0 == den%10 {
				continue
			}
			n10, n1 := num/10, num%10
			d10, d1 := den/10, den%10
			if n1 == d10 && float64(num)/float64(den) == float64(n10)/float64(d1) {
				fmt.Printf("%d/%d\n", num, den)
				ret = ret.MulRat(euler.NewRational(int64(num), int64(den)))
			}
			if n10 == d1 && float64(num)/float64(den) == float64(n1)/float64(d10) {
				fmt.Printf("%d/%d\n", num, den)
				ret = ret.MulRat(euler.NewRational(int64(num), int64(den)))
			}
			if n1 == d1 && float64(num)/float64(den) == float64(n10)/float64(d10) {
				fmt.Printf("%d/%d\n", num, den)
				ret = ret.MulRat(euler.NewRational(int64(num), int64(den)))
			}
			if n10 == d10 && float64(num)/float64(den) == float64(n1)/float64(d1) {
				fmt.Printf("%d/%d\n", num, den)
				ret = ret.MulRat(euler.NewRational(int64(num), int64(den)))
			}
		}
	}
	return ret
}

//	for ii in */*.go ; do go fmt "$ii" ; done ; for ii in 33 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
/*

16/64
19/95
26/65
49/98
Euler033: 1 / 100 ~~ 100



 */
func main() {
	//test

	//run
	ra := Euler033()
	fmt.Printf("Euler033: %d / %d ~~ %d\n", ra.Num, ra.Den, ra.Den)

}
