// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=99
https://projecteuler.net/minimal=99

<p>Comparing two numbers written in index form like $2^{11}$ and $3^7$ is not difficult, as any calculator would confirm that $2^{11} = 2048 \lt 3^7 = 2187$.</p>
<p>However, confirming that $632382^{518061} \gt 519432^{525806}$ would be much more difficult, as both numbers contain over three million digits.</p>
<p>Using <a href="resources/documents/0099_base_exp.txt">base_exp.txt</a> (right click and 'Save Link/Target As...'), a 22K text file containing one thousand lines with a base/exponent pair on each line, determine which line number has the greatest numerical value.</p>
<p class="smaller">NOTE: The first two lines in the file represent the numbers in the example given above.</p>


/
*/
/*
	In the real world, when asked this problem, I'm going to reach for math/big or GMP or a similar library for 'really big numbers'
	Oh math/big lacks a Pow function? ... This is fine.

	I also need to remark; look at that narrow bitlength progression in the best numbers.  Rounded Floats do not seem likely to be sufficient.
	Wonder how that datastream was generated... maybe an RNG initial number, then binary search exponent size with a target that's between two results, and if it doesn't fit discard?
/
*/

import (
	"bufio"
	"euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	"os" // os.Stdout
	"strconv"
	"strings"
)

func PowBigInt(base *big.Int, pow uint64) *big.Int {
	if 1 == pow {
		return base
	}
	// MSB (set) to LSB ; each step will be either np*np OR np*np + n
	// 100 == n*n, n*n * n*n ; 101 == n*n, n*n n*n * n
	np := big.NewInt(0).Set(base)
	for mask := uint64(1) << (64 - euler.BitsLeadingZeros64(pow) - 2); 0 < mask; mask >>= 1 {
		// fmt.Printf("%05b\n%05b\n\n", pow, mask)
		np.Mul(np, np) // double
		if 0 < pow&mask {
			np.Mul(np, base) // times n
		}
	}
	return np
}

func nilOrPanic(err error) {
	if nil != err {
		panic(err.Error())
	}
}

func Euler0099(fn string) uint64 {
	var ret, line uint64
	best, test := big.NewInt(0), big.NewInt(0)

	// Load the words, if they match a square number pattern
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler0099 unable to open: " + fn)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	// split lines is default, use one that chomps all ", (and whitespace) to output 'words'
	// scanner.Split(euler.ScannerSplitNLDQ)
	fmt.Println("Euler 99 takes several min to run...")
	for line = 1; scanner.Scan(); line++ {
		strnum := strings.Split(scanner.Text(), ",")
		pi, err := strconv.ParseInt(strnum[0], 10, 64)
		nilOrPanic(err)
		test.SetInt64(pi)
		pi, err = strconv.ParseInt(strnum[1], 10, 64)
		nilOrPanic(err)
		test = PowBigInt(test, uint64(pi))
		if 0 == line&0x3F {
			fmt.Printf("Line %#5x status (bitlengths): best = %d\tcurrent = %d\n", line, best.BitLen(), test.BitLen())
		}
		// bigint's Cmp is clamped subtract
		if -1 == best.Cmp(test) {
			ret = line
			best.Set(test)
		}
	}
	return ret
}

/*
for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 99 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done


Euler 99 takes several min to run...
Line  0x40 status (bitlengths): best = 9983435  current = 9983335
Line  0x80 status (bitlengths): best = 9983438  current = 9983426
Line  0xc0 status (bitlengths): best = 9983442  current = 9983349
Line 0x100 status (bitlengths): best = 9983442  current = 9983293
Line 0x140 status (bitlengths): best = 9983442  current = 9983427
Line 0x180 status (bitlengths): best = 9983442  current = 9983379
Line 0x1c0 status (bitlengths): best = 9983442  current = 9983367
Line 0x200 status (bitlengths): best = 9983442  current = 9983348
Line 0x240 status (bitlengths): best = 9983442  current = 9983340
Line 0x280 status (bitlengths): best = 9983442  current = 9983359
Line 0x2c0 status (bitlengths): best = 9983442  current = 9983373
Line 0x300 status (bitlengths): best = 9983444  current = 9983405
Line 0x340 status (bitlengths): best = 9983444  current = 9983273
Line 0x380 status (bitlengths): best = 9983444  current = 9983338
Line 0x3c0 status (bitlengths): best = 9983444  current = 9983407
Euler 99: Largest Exponential: 709

real    6m35.654s
user    6m35.286s
sys     0m0.673s
.
*/
func main() {
	var r uint64
	//test

	//run
	r = Euler0099("0099_base_exp.txt")
	fmt.Printf("Euler 99: Largest Exponential: %d\n", r)
	if 709 != r {
		panic("Did not reach expected value.")
	}
}
