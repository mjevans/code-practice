// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=48
https://projecteuler.net/minimal=48

<p>The series, $1^1 + 2^2 + 3^3 + \cdots + 10^{10} = 1_04050_71317$.</p>
<p>Find the last ten digits of the series, $1^1 + 2^2 + 3^3 + \cdots + 1000^{1000}$.</p>

*/
/*

1	1
2	4
3	27
4	256
5	3_125
6	46_656
7	823_543
8	16_777_216
9	387_420_489
10	10_000_000_000

If there's some formula that transforms the series, it's outside of my math knowledge.  I don't even have a vague recollection of encountering that trivia before.

*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler048(limit, mod uint64) uint64 {
	var ii, mm uint64
	res := big.NewInt(0)
	bii := big.NewInt(0)
	mul := big.NewInt(0)
	bmd := big.NewInt(int64(mod))
	for ii = 1; ii <= limit; ii++ {
		bii = bii.SetUint64(ii)
		mul = mul.Set(bii)
		// self power
		for mm = 1; mm < ii; mm++ {
			mul = mul.Mul(mul, bii)
		}
		res = res.Add(res, mul)
		res = res.Mod(res, bmd)
	}
	return res.Uint64()
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 48 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 48: test 1..10 = 10405071317 ? true 10405071317
Euler 48: Self Powers: 9110846700

*/
func main() {
	//test
	a := Euler048(10, 2_00000_00000)
	fmt.Printf("Euler 48: test 1..10 = 10405071317 ? %t %d\n", 10405071317 == a, a)

	//run
	b := Euler048(1_000, 1_00000_00000)
	fmt.Printf("Euler 48: Self Powers: %d\n", b)
}
