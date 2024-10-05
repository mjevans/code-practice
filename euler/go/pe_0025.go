// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=25
https://projecteuler.net/minimal=25

<p>The Fibonacci sequence is defined by the recurrence relation:</p>
<blockquote>$F_n = F_{n - 1} + F_{n - 2}$, where $F_1 = 1$ and $F_2 = 1$.</blockquote>
<p>Hence the first $12$ terms will be:</p>
\begin{align}
F_1 &amp;= 1\\
F_2 &amp;= 1\\
F_3 &amp;= 2\\
F_4 &amp;= 3\\
F_5 &amp;= 5\\
F_6 &amp;= 8\\
F_7 &amp;= 13\\
F_8 &amp;= 21\\
F_9 &amp;= 34\\
F_{10} &amp;= 55\\
F_{11} &amp;= 89\\
F_{12} &amp;= 144
\end{align}
<p>The $12$th term, $F_{12}$, is the first term to contain three digits.</p>
<p>What is the index of the first term in the Fibonacci sequence to contain $1000$ digits?</p>


https://en.wikipedia.org/wiki/Fibonacci_sequence#Matrix_form
https://www.nayuki.io/page/fast-fibonacci-algorithms
"""
Given F(k) and F(k+1)

F(2k) = F(k)[2F(k+1)−F(k)]
F(2k+1) = F(k+1)^2+F(k)^2
"""



*/

import (
	// "bufio"
	// "bitvector"
	"euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler025(digits int64) (fibnum, fibval *big.Int) {
	exp := big.NewInt(digits - 1)
	base := big.NewInt(int64(10))
	base.Exp(base, exp, nil)
	fmt.Printf("%d digits in %s\n", len(base.String()), base)
	one := big.NewInt(int64(1))
	two := big.NewInt(int64(2))
	search := big.NewInt(int64(2))
	// Search Power of 2 Fibs until target found.
	var test1, test2 *big.Int
	for ii := 0; ii < 0x0ffffff; ii++ {
		test1, test2 = euler.BigFib(search)
		c := base.Cmp(test1)
		if -1 == c {
			break
		} else if 0 == c {
			// test1 exactly equals the target
			return search, test1
		} else if 0 >= base.Cmp(test2) {
			// test2 is the first number given 2n, 2n + 1
			return search.Add(search, one), test2
		} else {
			search.Mul(search, two)
		}
		if ii > 0x0f && 0 == ii%0x0f {
			test1.Div(base, test1)
			fmt.Println("Euler025 searching for ", digits, " digits: targetMin / BigFib (iter ", ii, "): ", test1)
		}
	}

	// Between search / 2 (exclusive) and search (inclusive)
	var fmin, fmax *big.Int
	fmin = big.NewInt(int64(0))
	fmax = big.NewInt(int64(0))
	fmin.Div(search, two)
	fmax.Set(search)

	for ii := 0; ii < 0x0ffffff; ii++ {
		search.Add(fmax, fmin)
		search.Div(search, two)

		test1, test2 = euler.BigFib(search)
		c1 := base.Cmp(test1)
		c2 := base.Cmp(test2)
		if 1 == c1 && 0 >= c2 {
			// test2 is the first number given 2n, 2n + 1
			return search.Add(search, one), test2
		} else if 1 == c2 {
			// F(search)+1 too small
			search.Add(search, one)
			fmin.Set(search)
		} else {
			// F(search) equal or larger
			fmax.Set(search)
		}
		if true || 0 == ii%0x0f {
			test1.Div(base, test1)
			fmt.Println("Euler025 binary(ish) search: ", ii, ":\t", fmin, "..", fmax)
		}
	}

	return fmin, fmax
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 25 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Corrected off by one digit power...

What is the index of the first term in the Fibonacci sequence to contain 1000 digits?
1000 digits in 1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
Euler025 binary(ish) search:  0 :        4096 .. 6144
Euler025 binary(ish) search:  1 :        4096 .. 5120
Euler025 binary(ish) search:  2 :        4609 .. 5120
Euler025 binary(ish) search:  3 :        4609 .. 4864
Euler025 binary(ish) search:  4 :        4737 .. 4864
Euler025 binary(ish) search:  5 :        4737 .. 4800
Euler025 binary(ish) search:  6 :        4769 .. 4800
Euler025 binary(ish) search:  7 :        4769 .. 4784
Euler025 binary(ish) search:  8 :        4777 .. 4784
Euler025 binary(ish) search:  9 :        4781 .. 4784
Euler025 binary(ish) search:  10 :       4781 .. 4782
F( 4782 )       =>       1070066266382758936764980584457396885083683896632151665013235203375314520604694040621889147582489792657804694888177591957484336466672569959512996030461262748092482186144069433051234774442750273781753087579391666192149259186759553966422837148943113074699503439547001985432609723067290192870526447243726117715821825548491120525013201478612965931381792235559657452039506137551467837543229119602129934048260706175397706847068202895486902666185435124521900369480641357447470911707619766945691070098024393439617474103736912503231365532164773697023167755051595173518460579954919410967778373229665796581646513903488154256310184224190259846088000110186255550245493937113651657039447629584714548523425950428582425306083544435428212611008992863795048006894330309773217834864543113205765659868456288616808718693835297350643986297640660000723562917905207051164077614812491885830945940566688339109350944456576357666151619317753792891661581327159616877487983821820492520348473874384736771934512787029218636250627816
*/
func main() {
	//test
	last := int64(1)
	this := int64(1)
	passed := true
	for ii := int64(2); ii < int64(64); ii++ {
		test, _ := euler.BigFib(big.NewInt(ii))
		// fmt.Println(ii, this, test.Int64(), this == test.Int64())
		if this != test.Int64() {
			passed = false
		}
		tmp := this
		this += last
		last = tmp
	}
	if false == passed {
		panic("Failed BigFib test")
	}

	testFib := []struct{ test, res uint64 }{
		{1, 1}, {2, 1}, {3, 2}, {4, 3}, {5, 5}, {6, 8}, {7, 13}, {8, 21}, {9, 34}, {10, 55}, {11, 89}, {12, 144},
	}
	for _, test := range testFib {
		res, resPlus := euler.BigFib(big.NewInt(int64(test.test)))
		if uint64(res.Int64()) != test.res {
			fmt.Printf("Test case failed: Ref %d != %s (maybe %s?)\n", test.res, res, resPlus)
		}
	}

	//run
	fmt.Println("What is the index of the first term in the Fibonacci sequence to contain 1000 digits?")
	iter, val := Euler025(int64(1000))
	fmt.Println("F(", iter, ")\t=>\t", val)
}
