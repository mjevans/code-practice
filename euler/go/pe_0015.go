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
	"euler"
	"fmt"
	// "math"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler015(sqsz uint64) uint64 {
	return euler.PascalTri(sqsz<<1, sqsz)
}

/*

Wrong:	return (int64(math.Pow(2, float64(sqsz))) - 1) * 2
Wrong:	Euler15 test: 2 6 true
Wrong:	Euler15 test: 1 2 true
Wrong:	Euler15 test: 3 14 false
Wrong:	Euler15 test: 4 30 false
Wrong:	Euler15 20:  2097150

The hint I needed came from LeetCode's rectangle based variation of the same problem, which forced me to think from another angle... looking it up.

I don't often see Pascal's Triangle drawn out as a grid, and the progression of it's sequence is pretty but not so obviously connected to an 8th circle rotation of a rectangle / square, after all it's usually depicted as an equal-lateral or possibly even distorted to render the numbers, triangle; not a right angle at the top triangle.

- Unique Paths - https://leetcode.com/problems/unique-paths/
Pascal Triangle?
Rows and Cols start with 0
A row may be calculated in isolation by: (just roll forward)  NOTE: 'forward' might also be backwards since it's symmetrical!
F(n Row, k Index) = (F(n, k - 1) * ( n + 1 - k) / k )

	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 15 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler015 test: 2 6 true
Euler015 test: 1 2 true
Euler015 test: 3 20 true
Euler015 test: 4 70 true
Euler015 20:  137846528820

*/

func main() {
	// fmt.Println(grid)
	//test
	fmt.Println("Euler015 test: 2", Euler015(2), Euler015(2) == 6)
	fmt.Println("Euler015 test: 1", Euler015(1), Euler015(1) == 2)
	fmt.Println("Euler015 test: 3", Euler015(3), Euler015(3) == 20)
	fmt.Println("Euler015 test: 4", Euler015(4), Euler015(4) == 70)

	//run
	fmt.Println("Euler015 20: ", Euler015(20))
}
