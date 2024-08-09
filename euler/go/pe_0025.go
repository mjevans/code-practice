// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=25
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
	exp := big.NewInt(digits)
	base := big.NewInt(int64(10))
	base.Exp(base, exp, nil)
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

57 365435296162 365435296162 true
58 591286729879 591286729879 true
59 956722026041 956722026041 true
60 1548008755920 1548008755920 true
61 2504730781961 2504730781961 true
62 4052739537881 4052739537881 true
63 6557470319842 6557470319842 true
What is the index of the first term in the Fibonacci sequence to contain 1000 digits?
Euler025 binary(ish) search:  0 :        4096 .. 6144
Euler025 binary(ish) search:  1 :        4096 .. 5120
Euler025 binary(ish) search:  2 :        4609 .. 5120
Euler025 binary(ish) search:  3 :        4609 .. 4864
Euler025 binary(ish) search:  4 :        4737 .. 4864
Euler025 binary(ish) search:  5 :        4737 .. 4800
Euler025 binary(ish) search:  6 :        4769 .. 4800
Euler025 binary(ish) search:  7 :        4785 .. 4800
Euler025 binary(ish) search:  8 :        4785 .. 4792
Euler025 binary(ish) search:  9 :        4785 .. 4788
F( 4787 )       =>       11867216745258291596767088485966669273798582100095758927648586619975930687764095025968215177396570693265703962438125699711941059562545194266075961811883693134762216371218311196004424123489176045121333888565534924242378605373120526670329845322631737678903926970677861161240351447136066048164999599442542656514905088616976279305745609791746515632977790194938965236778055329967326038544356209745856855159058933476416258769264398373862584107011986781891656652294354303384242672408623790331963965457196174228574314820977014549061641307451101774166736940218594168337251710513138183086237827524393177246011800953414994670315197696419455768988692973700193372678236023166645886460311356376355559165284374295661676047742503016358708348137445254264644759334748027290043966390891843744407845769620260120918661264249498568399416752809338209739872047617689422485537053988895817801983866648336679027270843804302586168051835624516823216354234081479331553304809262608491851078404280454207286577699580222132259241827433

*/
func main() {
	//test
	last := int64(1)
	this := int64(1)
	passed := true
	for ii := int64(2); ii < int64(64); ii++ {
		test, _ := euler.BigFib(big.NewInt(ii))
		fmt.Println(ii, this, test.Int64(), this == test.Int64())
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

	//run
	fmt.Println("What is the index of the first term in the Fibonacci sequence to contain 1000 digits?")
	iter, val := Euler025(int64(1000))
	fmt.Println("F(", iter, ")\t=>\t", val)

}
