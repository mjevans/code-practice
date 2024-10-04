// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=11
https://projecteuler.net/minimal=11

<p>In the $20 \times 20$ grid below, four numbers along a diagonal line have been marked in red.</p>
<p class="monospace center">
08 02 22 97 38 15 00 40 00 75 04 05 07 78 52 12 50 77 91 08<br>
49 49 99 40 17 81 18 57 60 87 17 40 98 43 69 48 04 56 62 00<br>
81 49 31 73 55 79 14 29 93 71 40 67 53 88 30 03 49 13 36 65<br>
52 70 95 23 04 60 11 42 69 24 68 56 01 32 56 71 37 02 36 91<br>
22 31 16 71 51 67 63 89 41 92 36 54 22 40 40 28 66 33 13 80<br>
24 47 32 60 99 03 45 02 44 75 33 53 78 36 84 20 35 17 12 50<br>
32 98 81 28 64 23 67 10 <span class="red"><b>26</b></span> 38 40 67 59 54 70 66 18 38 64 70<br>
67 26 20 68 02 62 12 20 95 <span class="red"><b>63</b></span> 94 39 63 08 40 91 66 49 94 21<br>
24 55 58 05 66 73 99 26 97 17 <span class="red"><b>78</b></span> 78 96 83 14 88 34 89 63 72<br>
21 36 23 09 75 00 76 44 20 45 35 <span class="red"><b>14</b></span> 00 61 33 97 34 31 33 95<br>
78 17 53 28 22 75 31 67 15 94 03 80 04 62 16 14 09 53 56 92<br>
16 39 05 42 96 35 31 47 55 58 88 24 00 17 54 24 36 29 85 57<br>
86 56 00 48 35 71 89 07 05 44 44 37 44 60 21 58 51 54 17 58<br>
19 80 81 68 05 94 47 69 28 73 92 13 86 52 17 77 04 89 55 40<br>
04 52 08 83 97 35 99 16 07 97 57 32 16 26 26 79 33 27 98 66<br>
88 36 68 87 57 62 20 72 03 46 33 67 46 55 12 32 63 93 53 69<br>
04 42 16 73 38 25 39 11 24 94 72 18 08 46 29 32 40 62 76 36<br>
20 69 36 41 72 30 23 88 34 62 99 69 82 67 59 85 74 04 36 16<br>
20 73 35 29 78 31 90 01 74 31 49 71 48 86 81 16 23 57 05 54<br>
01 70 54 71 83 51 54 69 16 92 33 48 61 43 52 01 89 19 67 48<br></p>
<p>The product of these numbers is $26 \times 63 \times 78 \times 14 = 1788696$.</p>
<p>What is the greatest product of four adjacent numbers in the same direction (up, down, left, right, or diagonally) in the $20 \times 20$ grid?</p>



*/

import (
	"euler"
	"fmt"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler011(lim int, grid [][]int) (int, []int) {
	var runtmp, runbest []int
	runbest = make([]int, 4, 4)
	var max int
	for gg := 0; gg < len(grid)-1-lim; gg++ {
		for ii := 0; ii < len(grid[gg])-1-lim; ii++ {
			prod := euler.ListMul(grid[gg][ii : ii+lim])
			if prod > max {
				max = prod
				// runbest = nil
				copy(runbest, grid[gg][ii:ii+lim])
				fmt.Println("Horiz\t", max, runbest)
			}
			runtmp = []int{grid[gg][ii], grid[gg][ii+1], grid[gg][ii+2], grid[gg][ii+3]}
			prod = euler.ListMul(runtmp)
			if prod > max {
				max = prod
				// runbest = nil
				copy(runbest, runtmp)
				fmt.Println("Vert\t", max, runbest)
			}
			runtmp = []int{grid[gg][ii], grid[gg+1][ii+1], grid[gg+2][ii+2], grid[gg+3][ii+3]}
			prod = euler.ListMul(runtmp)
			// if 26 == runtmp[0] {
			// fmt.Println("Dbg \\\t", prod, runtmp)
			// }
			if prod > max {
				max = prod
				// runbest = nil
				copy(runbest, runtmp)
				fmt.Println("Diag \\\t", max, runbest)
			}
			runtmp = []int{grid[gg+3][ii], grid[gg+2][ii+1], grid[gg+1][ii+2], grid[gg][ii+3]}
			prod = euler.ListMul(runtmp)
			if prod > max {
				max = prod
				// runbest = nil
				copy(runbest, runtmp)
				fmt.Println("Diag /\t", max, runbest)
			}
		}
	}
	return max, runbest
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 11 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Horiz    34144 [8 2 22 97]
Diag \   279496 [8 49 31 23]
Diag /   24468444 [52 49 99 97]
Diag /   34826064 [68 67 98 78]
Diag /   41076896 [92 68 67 98]
Horiz    48477312 [78 78 96 83]
Diag /   70600674 [87 97 94 89]
Euler011:       Max for 4:       70600674

*/

func main() {
	var grid [][]int
	grid = append(grid, []int{8, 2, 22, 97, 38, 15, 0, 40, 0, 75, 4, 5, 7, 78, 52, 12, 50, 77, 91, 8})
	grid = append(grid, []int{49, 49, 99, 40, 17, 81, 18, 57, 60, 87, 17, 40, 98, 43, 69, 48, 4, 56, 62, 0})
	grid = append(grid, []int{81, 49, 31, 73, 55, 79, 14, 29, 93, 71, 40, 67, 53, 88, 30, 3, 49, 13, 36, 65})
	grid = append(grid, []int{52, 70, 95, 23, 4, 60, 11, 42, 69, 24, 68, 56, 1, 32, 56, 71, 37, 2, 36, 91})
	grid = append(grid, []int{22, 31, 16, 71, 51, 67, 63, 89, 41, 92, 36, 54, 22, 40, 40, 28, 66, 33, 13, 80})
	grid = append(grid, []int{24, 47, 32, 60, 99, 3, 45, 2, 44, 75, 33, 53, 78, 36, 84, 20, 35, 17, 12, 50})
	grid = append(grid, []int{32, 98, 81, 28, 64, 23, 67, 10, 26, 38, 40, 67, 59, 54, 70, 66, 18, 38, 64, 70})
	grid = append(grid, []int{67, 26, 20, 68, 2, 62, 12, 20, 95, 63, 94, 39, 63, 8, 40, 91, 66, 49, 94, 21})
	grid = append(grid, []int{24, 55, 58, 5, 66, 73, 99, 26, 97, 17, 78, 78, 96, 83, 14, 88, 34, 89, 63, 72})
	grid = append(grid, []int{21, 36, 23, 9, 75, 0, 76, 44, 20, 45, 35, 14, 0, 61, 33, 97, 34, 31, 33, 95})
	grid = append(grid, []int{78, 17, 53, 28, 22, 75, 31, 67, 15, 94, 3, 80, 4, 62, 16, 14, 9, 53, 56, 92})
	grid = append(grid, []int{16, 39, 5, 42, 96, 35, 31, 47, 55, 58, 88, 24, 0, 17, 54, 24, 36, 29, 85, 57})
	grid = append(grid, []int{86, 56, 0, 48, 35, 71, 89, 7, 5, 44, 44, 37, 44, 60, 21, 58, 51, 54, 17, 58})
	grid = append(grid, []int{19, 80, 81, 68, 5, 94, 47, 69, 28, 73, 92, 13, 86, 52, 17, 77, 4, 89, 55, 40})
	grid = append(grid, []int{4, 52, 8, 83, 97, 35, 99, 16, 7, 97, 57, 32, 16, 26, 26, 79, 33, 27, 98, 66})
	grid = append(grid, []int{88, 36, 68, 87, 57, 62, 20, 72, 3, 46, 33, 67, 46, 55, 12, 32, 63, 93, 53, 69})
	grid = append(grid, []int{4, 42, 16, 73, 38, 25, 39, 11, 24, 94, 72, 18, 8, 46, 29, 32, 40, 62, 76, 36})
	grid = append(grid, []int{20, 69, 36, 41, 72, 30, 23, 88, 34, 62, 99, 69, 82, 67, 59, 85, 74, 4, 36, 16})
	grid = append(grid, []int{20, 73, 35, 29, 78, 31, 90, 1, 74, 31, 49, 71, 48, 86, 81, 16, 23, 57, 5, 54})
	grid = append(grid, []int{1, 70, 54, 71, 83, 51, 54, 69, 16, 92, 33, 48, 61, 43, 52, 1, 89, 19, 67, 48})

	// fmt.Println(grid)
	//test
	prod, _ := Euler011(4, grid)
	fmt.Println("Euler011:\tMax for 4:\t", prod)
	//run
	// fmt.Println("Euler010:\t", Euler010(2000000))
}
