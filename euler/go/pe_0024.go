// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=24
https://projecteuler.net/minimal=24

<p>A permutation is an ordered arrangement of objects. For example, 3124 is one possible permutation of the digits 1, 2, 3 and 4. If all of the permutations are listed numerically or alphabetically, we call it lexicographic order. The lexicographic permutations of 0, 1 and 2 are:</p>
<p class="center">012   021   102   120   201   210</p>
<p>What is the millionth lexicographic permutation of the digits 0, 1, 2, 3, 4, 5, 6, 7, 8 and 9?</p>






*/

import (
	// "bufio"
	// "bitvector"
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

func Euler024() string {
	return ""
}

/*
euler.PermutationString(int64())

321

012
021
102
120
201
210

0 012
1 021
2 102
3 120
4 201
5 210
Euler024 test :  true
What is the millionth lexicographic permutation of the digits 0, 1, 2, 3, 4, 5, 6, 7, 8 and 9?

	1000000 2783915604
	1000000 2783915640
*/
func main() {
	// fmt.Println(grid)
	//test
	for ii := 0; ii <= 5; ii++ {
		fmt.Println(ii, euler.PermutationString(ii, "012"))
	}
	test := euler.PermutationString(4, "012") == "201" &&
		euler.PermutationString(1, "012") == "021" &&
		euler.PermutationString(5, "012") == "210" &&
		euler.PermutationString(0, "012") == "012"
	fmt.Println("Euler024 test : ", test)

	//run
	fmt.Println("What is the millionth lexicographic permutation of the digits 0, 1, 2, 3, 4, 5, 6, 7, 8 and 9?\n",
		999999, euler.PermutationString(999999, "0123456789"), "\n",
		1000000, euler.PermutationString(1000000, "0123456789"), "\n",
		1000001, euler.PermutationString(1000001, "0123456789"))

}
