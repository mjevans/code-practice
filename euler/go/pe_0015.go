// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=16
<p>$2^{15} = 32768$ and the sum of its digits is $3 + 2 + 7 + 6 + 8 = 26$.</p>
<p>What is the sum of the digits of the number $2^{1000}$?</p>





*/

import (
	// "euler"
	"fmt"
	"math"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler015(sqsz int64) int64 {
	return (int64(math.Pow(2, float64(sqsz))) - 1) * 2
}

/*
Euler15 test: 2 6 true
Euler15 test: 1 2 true
Euler15 test: 3 14 false
Euler15 test: 4 30 false
Euler15 20:  2097150
*/

func main() {
	// fmt.Println(grid)
	//test
	fmt.Println("Euler15 test: 2", Euler015(2), Euler015(2) == 6)
	fmt.Println("Euler15 test: 1", Euler015(1), Euler015(1) == 2)
	fmt.Println("Euler15 test: 3", Euler015(3), Euler015(3) == 0)
	fmt.Println("Euler15 test: 4", Euler015(4), Euler015(4) == 0)

	//run
	fmt.Println("Euler15 20: ", Euler015(20))
}
