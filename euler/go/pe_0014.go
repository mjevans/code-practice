// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=14
https://projecteuler.net/minimal=14

<p>The following iterative sequence is defined for the set of positive integers:</p>
<ul style="list-style-type:none;">
<li>$n \to n/2$ ($n$ is even)</li>
<li>$n \to 3n + 1$ ($n$ is odd)</li></ul>
<p>Using the rule above and starting with $13$, we generate the following sequence:
$$13 \to 40 \to 20 \to 10 \to 5 \to 16 \to 8 \to 4 \to 2 \to 1.$$</p>
<p>It can be seen that this sequence (starting at $13$ and finishing at $1$) contains $10$ terms. Although it has not been proved yet (Collatz Problem), it is thought that all starting numbers finish at $1$.</p>
<p>Which starting number, under one million, produces the longest chain?</p>
<p class="note"><b>NOTE:</b> Once the chain starts the terms are allowed to go above one million.</p>






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

func Euler014(start, end int64) int64 {
	// This could use a cache... or just run
	var maxloops, maxii int64
	for ii := start; ii <= end; ii++ {
		loops := int64(0)
		for val := ii; loops < 0x7FFFFFFFFFFFFFFF && val > 1; loops++ {
			if 0 == val&0x1 {
				val /= 2
			} else {
				val = 3*val + 1
			}
		}
		loops++
		if maxloops < loops {
			maxloops = loops
			maxii = ii
			// fmt.Println("Euler014: New Longest Chain:\t", maxii, "\tloops ", maxloops)
		}
	}
	return maxii
}

func main() {
	// fmt.Println(grid)
	//test
	_ = Euler014(13, 13)
	// fmt.Println(euler.BCDadd([]string{"5", "5", "10"}))
	//run
	// Euler014: New Longest Chain:     837799         loops  525
	at := Euler014(1, 1000000)
	fmt.Println("Euler014:\tLongest run happend on: ", at)
}
