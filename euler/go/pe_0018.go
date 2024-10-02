// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=18
https://projecteuler.net/minimal=18

<p>By starting at the top of the triangle below and moving to adjacent numbers on the row below, the maximum total from top to bottom is $23$.</p>
<p class="monospace center"><span class="red"><b>3</b></span><br><span class="red"><b>7</b></span> 4<br>
2 <span class="red"><b>4</b></span> 6<br>
8 5 <span class="red"><b>9</b></span> 3</p>
<p>That is, $3 + 7 + 4 + 9 = 23$.</p>
<p>Find the maximum total from top to bottom of the triangle below:</p>
<p class="monospace center">75<br>
95 64<br>
17 47 82<br>
18 35 87 10<br>
20 04 82 47 65<br>
19 01 23 75 03 34<br>
88 02 77 73 07 63 67<br>
99 65 04 28 06 16 70 92<br>
41 41 26 56 83 40 80 70 33<br>
41 48 72 33 47 32 37 16 94 29<br>
53 71 44 65 25 43 91 52 97 51 14<br>
70 11 33 28 77 73 17 78 39 68 17 57<br>
91 71 52 38 17 14 91 43 58 50 27 29 48<br>
63 66 04 68 89 53 67 30 73 16 69 87 40 31<br>
04 62 98 27 23 09 70 98 73 93 38 53 60 04 23</p>
<p class="note"><b>NOTE:</b> As there are only $16384$ routes, it is possible to solve this problem by trying every route. However, <a href="problem=67">Problem 67</a>, is the same challenge with a triangle containing one-hundred rows; it cannot be solved by brute force, and requires a clever method! ;o)</p>




*/

import (
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler018() int {
	var triTest, tri [][]int
	triTest = append(triTest, []int{3})
	triTest = append(triTest, []int{7, 4})
	triTest = append(triTest, []int{2, 4, 6})
	triTest = append(triTest, []int{8, 5, 9, 3})
	test := euler.MaximumPathSum(triTest)
	if 23 != test {
		panic(fmt.Sprint("Failed test, expected 23, got: ", test))
	}
	tri = append(tri, []int{75})
	tri = append(tri, []int{95, 64})
	tri = append(tri, []int{17, 47, 82})
	tri = append(tri, []int{18, 35, 87, 10})
	tri = append(tri, []int{20, 4, 82, 47, 65})
	tri = append(tri, []int{19, 1, 23, 75, 3, 34})
	tri = append(tri, []int{88, 2, 77, 73, 7, 63, 67})
	tri = append(tri, []int{99, 65, 4, 28, 6, 16, 70, 92})
	tri = append(tri, []int{41, 41, 26, 56, 83, 40, 80, 70, 33})
	tri = append(tri, []int{41, 48, 72, 33, 47, 32, 37, 16, 94, 29})
	tri = append(tri, []int{53, 71, 44, 65, 25, 43, 91, 52, 97, 51, 14})
	tri = append(tri, []int{70, 11, 33, 28, 77, 73, 17, 78, 39, 68, 17, 57})
	tri = append(tri, []int{91, 71, 52, 38, 17, 14, 91, 43, 58, 50, 27, 29, 48})
	tri = append(tri, []int{63, 66, 4, 68, 89, 53, 67, 30, 73, 16, 69, 87, 40, 31})
	tri = append(tri, []int{4, 62, 98, 27, 23, 9, 70, 98, 73, 93, 38, 53, 60, 4, 23})
	ans := euler.MaximumPathSum(tri)
	fmt.Println("Euler018 maxumum path:\t", ans)
	return ans
	// Euler018 maxumum path:   1074
}

func main() {
	// fmt.Println(grid)
	//test

	//run
	Euler018()
}
