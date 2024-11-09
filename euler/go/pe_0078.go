// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=78
https://projecteuler.net/minimal=78

<p>Let $p(n)$ represent the number of different ways in which $n$ coins can be separated into piles. For example, five coins can be separated into piles in exactly seven different ways, so $p(5)=7$.</p>
<div class="margin_left">
OOOOO<br>
OOOO   O<br>
OOO   OO<br>
OOO   O   O<br>
OO   OO   O<br>
OO   O   O   O<br>
O   O   O   O   O
</div>
<p>Find the least value of $n$ for which $p(n)$ is divisible by one million.</p>

/
*/
/*
	https://en.wikipedia.org/wiki/Integer_partition#Partition_function
	Hans Rademacher found a function that involves square roots, fractions, Pi, e, and Sinh as well as two summations; It's probably at least faster for large numbers given it's mentioned.
	https://en.wikipedia.org/wiki/Pentagonal_number_theorem
	Pentagonal numbers, PowMod, this might be a little slower, but the terms are integers and can be cached.
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

func Euler0078(div int64) int64 {
	// Pentagonal numbers are correct for values greater than zero, hard code zero (result 1) and negative numbers (all zero)
	partcache := append(make([]int64, 0, 1<<15), 1)

	PentagonalNumberI64 := func(n int64) int64 {
		return (n * (3*n - 1)) >> 1
	}
	var EulerPentagonalPartitionGen func(n int64) int64
	EulerPentagonalPartitionGen = func(n int64) int64 {
		if n < int64(len(partcache)) {
			return partcache[n]
		}
		// https://en.wikipedia.org/wiki/Pentagonal_number_theorem
		var k, g, val int64
		k = 1
		for {
			g = PentagonalNumberI64(k)
			// Pentagonal numbers are correct for values greater than zero, hard code zero (result 1) and negative numbers (all zero)
			if 0 > (n - g) {
				// fmt.Printf("%d:\t%d\t%d\t%d\tbreak\n", n, k, g, val)
				break
			}
			if 1 == k&1 {
				val += partcache[n-g]
			} else {
				val -= partcache[n-g]
			}
			val %= div
			// fmt.Printf("%d:\t%d\t%d\t%d\n", n, k, g, val)
			if k < 0 {
				k = -k + 1
			} else {
				k = -k
			}
		}
		return val
	}

	for pnum := int64(1); pnum <= div; pnum++ {
		res := EulerPentagonalPartitionGen(pnum)
		// fmt.Printf("%d %% %d @ %d\n", res, div, pnum)
		if 0 == res%div {
			return pnum
		}
		partcache = append(partcache, res)
	}
	return 0

}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 78 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 78: Coin Partitions:     Count: 55374

real    0m0.275s
user    0m0.334s
sys     0m0.057s
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	tests := []struct {
		res, a int64
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 5},
		{5, 7},
		{6, 11},
		{7, 15},
		{8, 22},
		{9, 30},
		{10, 42},
		{11, 56},
		{12, 77},
	}
	for _, test := range tests {
		r := Euler0078(test.a)
		if test.res != r {
			panic(fmt.Sprintf("Euler 78: Expected %d got %d", test.res, r))
		}
	}

	//run
	r := Euler0078(1_000_000)
	fmt.Printf("Euler 78: Coin Partitions:\tCount: %d\n", r)
	if 55374 != r {
		panic("Did not reach expected value.")
	}
}
