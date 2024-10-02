// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=15
https://projecteuler.net/minimal=15

<p>Starting in the top left corner of a $2 \times 2$ grid, and only being able to move to the right and down, there are exactly $6$ routes to the bottom right corner.</p>
<div class="center">
<img src="resources/images/0015.png?1678992052" class="dark_img" alt=""></div>
<p>How many such routes are there through a $20 \times 20$ grid?</p>






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
