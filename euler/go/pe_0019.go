// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=19
https://projecteuler.net/minimal=19


https://en.wikipedia.org/wiki/Doomsday_rule#Finding_a_year's_anchor_day


<p>You are given the following information, but you may prefer to do some research for yourself.</p>
<ul><li>1 Jan 1900 was a Monday.</li>
<li>Thirty days has September,<br />
April, June and November.<br />
All the rest have thirty-one,<br />
Saving February alone,<br />
Which has twenty-eight, rain or shine.<br />
And on leap years, twenty-nine.</li>

<li>A leap year occurs on any year evenly divisible by 4, but not on a century unless it is divisible by 400.</li>
</ul><p>How many Sundays fell on the first of the month during the twentieth century (1 Jan 1901 to 31 Dec 2000)?</p>

So 1900 through 2000 (inclusive) two Doomsday_rule periods.



*/

import (
	// "euler"
	"fmt"
	"time"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler019() int {
	ret := 0
	for yy := 1901; yy <= 2000; yy++ {
		for mm := 1; mm <= 12; mm++ {
			if time.Sunday == time.Date(yy, time.Month(mm), 1, 0, 0, 0, 0, time.UTC).Weekday() {
				ret++
			}
		}
	}
	return ret
}

func main() {
	// fmt.Println(grid)
	//test

	//run
	// 173 Sundays fall on the first of the month between 1900 and 2000 (inclusive).
	fmt.Println(Euler019(), "Sundays fall on the first of the month between 1901 and 2000 (inclusive).")
}
