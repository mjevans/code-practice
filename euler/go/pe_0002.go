// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=2
https://projecteuler.net/minimal=2

c 2001
Even Fibonacci Numbers
[Show HTML problem content]
Problem 2

Each new term in the Fibonacci sequence is generated by adding the previous two terms. By starting with 1 and 2, the first terms will be:
1,2,3,5,8,13,21,34,55,89,...
By considering the terms in the Fibonacci sequence whose values do not exceed four million, find the sum of the even-valued terms.

++
https://en.wikipedia.org/wiki/C_data_types#inttypes.h
32 bit integers are upto ~2B

*/

import (
	"fmt"
	// "strings"
	// "os" // os.Stdout
)

func main() {
	var sum, xx, yy int   // use native word size, 32 or 64 bits, either is enough, default 0.
	sum, xx, yy = 2, 1, 2 // Start with 2 already added and evaluated, xx is 'left' / first round.
	for {                 // forever
		// Instead of swapping, unroll the loop and do a left / right style version.
		// left
		xx += yy
		if xx > 4000000 {
			break
		} else if 0 == xx&0x1 {
			sum += xx
		}
		// right
		yy += xx
		if yy > 4000000 {
			break
		} else if 0 == yy&0x1 {
			sum += yy
		}
	}
	fmt.Printf("%d\n", sum)
}
