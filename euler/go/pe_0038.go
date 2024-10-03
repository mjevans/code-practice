// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=38
https://projecteuler.net/minimal=38

<p>Take the number $192$ and multiply it by each of $1$, $2$, and $3$:</p>
\begin{align}
192 \times 1 &amp;= 192\\
192 \times 2 &amp;= 384\\
192 \times 3 &amp;= 576
\end{align}
<p>By concatenating each product we get the $1$ to $9$ pandigital, $192384576$. We will call $192384576$ the concatenated product of $192$ and $(1,2,3)$.</p>
<p>The same can be achieved by starting with $9$ and multiplying by $1$, $2$, $3$, $4$, and $5$, giving the pandigital, $918273645$, which is the concatenated product of $9$ and $(1,2,3,4,5)$.</p>
<p>What is the largest $1$ to $9$ pandigital $9$-digit number that can be formed as the concatenated product of an integer with $(1,2, \dots, n)$ where $n \gt 1$?</p>


*/
/*

The example 9 * ( 1, 2, 3, 4, 5 ) has been provided which yields the 9 digit concatenation: 918273645

So any number that wins HAS to start with 9, otherwise there's no point.

N must be >= 2 so the concatenation must also yield at least two numbers, and presumably the longer is the second, thus 4 digits + 5 digits.

Hence the limit of < 10000 for the hunt number, has to START with 9 and must be greater than 918273645

Do I need to make pandigital a test function? # Euler 32 was the previous use.


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

func Euler038() uint64 {
	max := uint64(918273645)

	// Number * 1 cannot include a zero, nor can any other number...
	// 9 is in front so it can't be used again in the rest of the digits... the first filter case will trap the rest
Euler038Outer:
	for ii := uint64(1); ii <= 876; ii++ {
		var nn, ncat uint64
		if ii >= 100 {
			nn = 9000 + ii
		} else if ii >= 10 {
			nn = 900 + ii
		} else {
			nn = 90 + ii
		}
		ncat = 0
		var used uint16
		used = 0
		// fmt.Printf("\ntest:")
		for jj := uint64(1); 0b0000_0011_1111_1110 != used; jj++ {
			test := nn * jj
			// fmt.Printf("\t%d", test)
			digits := uint64(1)
			for test > 0 {
				bd := uint16(1) << (test % 10)
				test /= 10
				digits *= 10
				if 0 < used&bd || bd == 1 {
					// fmt.Printf("SKIP: %d : dupe or 0 digit : %d\n", nn, test%10)
					continue Euler038Outer
				}
				used |= bd
			}
			ncat = ncat*digits + nn*jj
		}
		// Pandigital, but is it greater?
		if max < ncat && 0b0000_0011_1111_1110 == used {
			fmt.Printf("Found new max: %d < %d (%d)\n", max, ncat, nn)
			max = ncat
		} else {
			fmt.Printf("SKIP: %d > %d (%d) ~= %b \n", max, ncat, nn, used)
		}
	}

	return max
}

func Euler038_harder() uint64 {
	max := uint64(0)

	// Number * 1 cannot include a zero, nor can any other number...
	// 9 is in front so it can't be used again in the rest of the digits... the first filter case will trap the rest
Euler038HOuter:
	for ii := uint64(1); ii <= 10000; ii++ {
		var nn, ncat uint64
		nn = ii
		ncat = 0
		var used uint16
		used = 0
		// fmt.Printf("\ntest:")
		for jj := uint64(1); 0b0000_0011_1111_1110 != used; jj++ {
			test := nn * jj
			// fmt.Printf("\t%d", test)
			digits := uint64(1)
			for test > 0 {
				bd := uint16(1) << (test % 10)
				test /= 10
				digits *= 10
				if 0 < used&bd || bd == 1 {
					// fmt.Printf("SKIP: %d : dupe or 0 digit : %d\n", nn, test%10)
					continue Euler038HOuter
				}
				used |= bd
			}
			// fmt.Printf("%d\t%d : %d\n", ncat, nn*jj, digits)
			ncat = ncat*digits + nn*jj
		}
		// Pandigital, but is it greater?
		if max < ncat && 0b0000_0011_1111_1110 == used {
			fmt.Printf("Found new max: %d < %d (%d)\n", max, ncat, nn)
			max = ncat
		} else {
			fmt.Printf("SKIP: %d > %d (%d) ~= %b \n", max, ncat, nn, used)
		}
	}

	return max
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 38 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Found new max: 918273645 < 926718534 (9267)
Found new max: 926718534 < 927318546 (9273)
Found new max: 927318546 < 932718654 (9327)

Euler038: Pandigit Multiples :  932718654
Found new max: 0 < 123456789 (1)
Found new max: 123456789 < 918273645 (9)
SKIP: 918273645 > 192384576 (192) ~= 1111111110
SKIP: 918273645 > 219438657 (219) ~= 1111111110
SKIP: 918273645 > 273546819 (273) ~= 1111111110
SKIP: 918273645 > 327654981 (327) ~= 1111111110
SKIP: 918273645 > 672913458 (6729) ~= 1111111110
SKIP: 918273645 > 679213584 (6792) ~= 1111111110
SKIP: 918273645 > 692713854 (6927) ~= 1111111110
SKIP: 918273645 > 726914538 (7269) ~= 1111111110
SKIP: 918273645 > 729314586 (7293) ~= 1111111110
SKIP: 918273645 > 732914658 (7329) ~= 1111111110
SKIP: 918273645 > 769215384 (7692) ~= 1111111110
SKIP: 918273645 > 792315846 (7923) ~= 1111111110
SKIP: 918273645 > 793215864 (7932) ~= 1111111110
Found new max: 918273645 < 926718534 (9267)
Found new max: 926718534 < 927318546 (9273)
Found new max: 927318546 < 932718654 (9327)

Euler038: Try harder... :       932718654



*/
func main() {
	//test

	//run
	r := Euler038()
	fmt.Printf("\nEuler038: Pandigit Multiples :\t%d\n", r)

	r = Euler038_harder()
	fmt.Printf("\nEuler038: Try harder... :\t%d\n", r)

}
