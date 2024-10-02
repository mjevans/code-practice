// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=20
https://projecteuler.net/minimal=20


<p>$n!$ means $n \times (n - 1) \times \cdots \times 3 \times 2 \times 1$.</p>
<p>For example, $10! = 10 \times 9 \times \cdots \times 3 \times 2 \times 1 = 3628800$,<br>and the sum of the digits in the number $10!$ is $3 + 6 + 2 + 8 + 8 + 0 + 0 = 27$.</p>
<p>Find the sum of the digits in the number $100!$.</p>



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

func Euler020(fact int) int64 {
	return euler.AddBigIntDecDigits(euler.BigFactorial(int64(fact)))
}

/*
10! Euler20sum = 27? 27 true
100! Euler20sum = ??? 648
*/
func main() {
	// fmt.Println(grid)
	//test
	fmt.Println("10! Euler20sum = 27?", Euler020(10), Euler020(10) == int64(27))

	//run
	fmt.Println("100! Euler20sum = ???", Euler020(100))
}
