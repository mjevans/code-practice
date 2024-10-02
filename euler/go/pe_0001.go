// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=1
https://projecteuler.net/minimal=1

c 2001
Multiples of 3 or 5
[Show HTML problem content]
Problem 1

If we list all the natural numbers below 10 that are multiples of 3 or 5, we get and 3, 5, 6, 9. The sum of these multiples is 23.

Find the sum of all the multiples of 3 or 5 below 1000.

A: 233168
*/

import (
	"fmt"
	// "strings"
	// "os" // os.Stdout
)

func main() {
	// About: for Range https://stackoverflow.com/questions/21950244/is-there-a-way-to-iterate-over-a-range-of-integers
	var sum int // use native word size, 32 or 64 bits, either is enough, default 0.
	for ii := 1; ii < 1000; ii++ {
		if (0 == ii%5) || (0 == ii%3) {
			sum += ii
		}
	}
	fmt.Printf("%d\n", sum)
}
