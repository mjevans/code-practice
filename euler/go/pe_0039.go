// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=39
https://projecteuler.net/minimal=39

<p>If $p$ is the perimeter of a right angle triangle with integral length sides, $\{a, b, c\}$, there are exactly three solutions for $p = 120$.</p>
<p>$\{20,48,52\}$, $\{24,45,51\}$, $\{30,40,50\}$</p>
<p>For which value of $p \le 1000$, is the number of solutions maximised?</p>



*/
/*

https://en.wikipedia.org/wiki/Right_triangle#Characterizations

It's easy to remember a^2 + b^2 = c^2 ,
a bit harder to recall a <= b < c

I never remember seeing (though a math book probably covered it and never touched it again)

s == 1/2 * (a + b + c)
(s-a)(s-b)=s(s-c)

Wow that looks so much easier to write a program for than caching the squares of numbers and probing.

*/

import (
	// "bufio"
	// "bitvector"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler039(min, limit uint) uint {
	var max, maxCombos uint
	if 0 < min&1 {
		min += 1
	}
	if 0 < limit&1 {
		limit ^= 1
	}
	for ii := min; ii <= limit; ii += 2 {
		s := ii >> 1
		combos := uint(0)
		aS := ii / 3
		for a := uint(1); a < aS; a++ {
			for b := a; b < s; b++ {
				// (s-a)(s-b)=s(s-c) => ((s-a)(s-b))/s=s-c => (((s-a)(s-b))/s)-s=-c
				// c = s-(((s-a)(s-b))/s)
				c := s - (((s - a) * (s - b)) / s)

				// Verify after possible integer truncation
				if (s-a)*(s-b) == s*(s-c) && a*a+b*b == c*c {
					combos++
					if 840 == ii {
						fmt.Printf("%d triangle? %d, %d, %d\t(%d : %d)\n", ii, a, b, c, s, aS)
					}
				}
			}
		}
		if maxCombos < combos {
			max = ii
			maxCombos = combos
			fmt.Printf("%d new max combos %d\n", ii, combos)
		}
	}

	return max
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 39 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

120 new max combos 3

Euler039: Integer Right Triangles TEST 120 :    120
12 new max combos 1
60 new max combos 2
120 new max combos 3
240 new max combos 4
420 new max combos 5
720 new max combos 6
840 triangle? 40, 399, 401      (420 : 280)
840 triangle? 56, 390, 394      (420 : 280)
840 triangle? 105, 360, 375     (420 : 280)
840 triangle? 120, 350, 370     (420 : 280)
840 triangle? 140, 336, 364     (420 : 280)
840 triangle? 168, 315, 357     (420 : 280)
840 triangle? 210, 280, 350     (420 : 280)
840 triangle? 240, 252, 348     (420 : 280)
840 new max combos 8

Euler039: Integer Right Triangles :     840



*/
func main() {
	//test
	r := Euler039(120, 120)
	fmt.Printf("\nEuler039: Integer Right Triangles TEST 120 :\t%d\n", r)

	//run
	r = Euler039(2, 1000)
	fmt.Printf("\nEuler039: Integer Right Triangles :\t%d\n", r)

}
