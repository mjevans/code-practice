// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=17
https://projecteuler.net/minimal=17

<p>If the numbers $1$ to $5$ are written out in words: one, two, three, four, five, then there are $3 + 3 + 5 + 4 + 4 = 19$ letters used in total.</p>
<p>If all the numbers from $1$ to $1000$ (one thousand) inclusive were written out in words, how many letters would be used? </p>
<br><p class="note"><b>NOTE:</b> Do not count spaces or hyphens. For example, $342$ (three hundred and forty-two) contains $23$ letters and $115$ (one hundred and fifteen) contains $20$ letters. The use of "and" when writing out numbers is in compliance with British usage.</p>




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

func Euler017(start, end int) int {
	// <p>If all the numbers from $1$ to $1000$ (one thousand) inclusive were written out in words, how many letters would be used? </p>
	var printed int
	for ; start <= end; start++ {
		lp, _ := euler.StringBritishCheckNumber(start)
		printed += lp
	}
	return printed
}

/*
1 One
11 Eleven
19 Nineteen
186 One Hundred and Eighty Six
1000 One Thousand
How many printed characters if 1..1000 are written out like in check numbers? :  21034
*/
func SecondIntString(x int, ret string) string { return ret }

func main() {
	// fmt.Println(grid)
	//test
	fmt.Println(1, SecondIntString(euler.StringBritishCheckNumber(1)))
	fmt.Println(11, SecondIntString(euler.StringBritishCheckNumber(11)))
	fmt.Println(19, SecondIntString(euler.StringBritishCheckNumber(19)))
	fmt.Println(186, SecondIntString(euler.StringBritishCheckNumber(186)))
	fmt.Println(1000, SecondIntString(euler.StringBritishCheckNumber(1000)))

	//run
	fmt.Println("How many printed characters if 1..1000 are written out like in check numbers? : ", Euler017(1, 1000))
}
