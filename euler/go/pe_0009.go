// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=9
https://projecteuler.net/minimal=9

<p>A Pythagorean triplet is a set of three natural numbers, $a \lt b \lt c$, for which,
$$a^2 + b^2 = c^2.$$</p>
<p>For example, $3^2 + 4^2 = 9 + 16 = 25 = 5^2$.</p>
<p>There exists exactly one Pythagorean triplet for which $a + b + c = 1000$.<br>Find the product $abc$.</p>



*/

import (
	// "euler"
	"fmt"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler009() []int {
	for a := 1; a < 1000-2; a++ {
		for b := 1; b < 1000-1-a; b++ {
			c := 1000 - b - a
			// A B C satisfy: a + b + c = 1000
			if a*a+b*b == c*c {
				return []int{a, b, c}
			}
		}
	}
	return []int{0, 0, 0}
}

func main() {
	fmt.Println("Euler009:\t", Euler009())
}
