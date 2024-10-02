// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=16
https://projecteuler.net/minimal=16

<p>$2^{15} = 32768$ and the sum of its digits is $3 + 2 + 7 + 6 + 8 = 26$.</p>
<p>What is the sum of the digits of the number $2^{1000}$?</p>





*/

import (
	// "euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler016(exp int64) int64 {

	// check := int64(0)
	// if exp < 64 {
		// check = int64(math.Pow(2, float64(exp)))
	// }

	// 1000 is many more bits than are supported by normal machine types...
	zero := big.NewInt(int64(0))
	two := big.NewInt(int64(2))
	b := big.NewInt(int64(1))
	for ; exp > 0; exp-- {
		b.Mul(b, two)
	}
	// if check > 0 {
		// fmt.Println(check == b.Int64(), check, " == ", b.Int64())
	// }

	ret := int64(0)
	ten := big.NewInt(int64(10))
	rem := big.NewInt(int64(0))
	limit := 0x7FFFFFFF
	for 0 < b.Cmp(zero) && limit > 0 {
		limit--
		b.DivMod(b, ten, rem)
		ret += rem.Int64()
	}
	return ret
}
/*
true 32768  ==  32768
15 26
1000 1366
*/
func main() {
	// fmt.Println(grid)
	//test
	fmt.Println(15, Euler016(15))
	//run
	fmt.Println(1000, Euler016(1000))
}
