// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=90
https://projecteuler.net/minimal=90

<p>Each of the six faces on a cube has a different digit ($0$ to $9$) written on it; the same is done to a second cube. By placing the two cubes side-by-side in different positions we can form a variety of $2$-digit numbers.</p>

<p>For example, the square number $64$ could be formed:</p>

<div class="center">
<img src="resources/images/0090.png?1678992052" class="dark_img" alt=""><br></div>

<p>In fact, by carefully choosing the digits on both cubes it is possible to display all of the square numbers below one-hundred: $01$, $04$, $09$, $16$, $25$, $36$, $49$, $64$, and $81$.</p>

<p>For example, one way this can be achieved is by placing $\{0, 5, 6, 7, 8, 9\}$ on one cube and $\{1, 2, 3, 4, 8, 9\}$ on the other cube.</p>

<p>However, for this problem we shall allow the $6$ or $9$ to be turned upside-down so that an arrangement like $\{0, 5, 6, 7, 8, 9\}$ and $\{1, 2, 3, 4, 6, 7\}$ allows for all nine square numbers to be displayed; otherwise it would be impossible to obtain $09$.</p>

<p>In determining a distinct arrangement we are interested in the digits on each cube, not the order.</p>

<ul style="list-style-type:none;"><li>$\{1, 2, 3, 4, 5, 6\}$ is equivalent to $\{3, 6, 4, 1, 2, 5\}$</li>
<li>$\{1, 2, 3, 4, 5, 6\}$ is distinct from $\{1, 2, 3, 4, 5, 9\}$</li></ul>

<p>But because we are allowing $6$ and $9$ to be reversed, the two distinct sets in the last example both represent the extended set $\{1, 2, 3, 4, 5, 6, 9\}$ for the purpose of forming $2$-digit numbers.</p>

<p>How many distinct arrangements of the two cubes allow for all of the square numbers to be displayed?</p>


/
*/
/*
	They already mentioned 9 and 6 are interchangable.
	Though I think there might be an inherent trap.
	As I started to ponder automating the problem, I realized that the dice could be swapped too.

	0	1	2	3	4	5	6	7	8	9
01	0	1									OK
04	0				4						OK
09	0									9	OK
16		1					6				OK
25			2			5
36				3			6
49					4					9	OK
64					4		6				OK
81		1							8

	[0	2	0	0	2	1	4	0	0	2]
	[3	1	1	1	1	0	1	0	1	0]
	[3	3	1	1	3	1	+5	0	1	-]
	OK	OK	2	3	OK	5	OK		8	OK

Dice MUST have... given zero...
{0,} {1,4,9~6,}
It's a bit annoying, but 1 and 4, the second most frequently occuring numbers, also happen next to a 9 or a 6.
{0,6~9,} {1,4,6~9,}
The above is where I should start if brute forcing stuff... Using 5 slots out of 12 slots... it would be possible to just brute force to exhaust any possibly optimal set.

There are 4 other numbers, which must be placed 2 and 5 must be on different dice, 8 across from a one (already on one/both), and 3 against a 6 (also alrady on both)
There's room to put 1 and 4 on the other dice, and 0 on the right dice...
{0,6~9,2,3,1,4} {1,4,6~9,5,8,0}

Combo powers:
x3,3,3,5,	1,1		x3,3,3,5,	1,1
{0,1,4,6~9,2,3}	{0,1,4,6~9,5,8}

Though I now read again, this problem IS entirely brute force.  They aren't asking for the set of cubes that could display the most combinations, it's the set of all potential permutations which yield all cube matches...

Sigh

Recycling that groundwork:
----a 6 or a 9 MUST be on each dice, more can work---- Not strictly true for ANY valid dice...
0 must be across from 1, 4, 9~6  and 9~6 must be across from 4,3,1,0
AND a 0 MUST be across from 1, 4 --6 or 9--
AND a 2 MUST be across from a 5 ++ NOTE: The problem doesn't mention 2 and 5 looking similar upside down, which some fonts (E.G. 7 seg display) might offer
AND an 8 MUST be across from a 1

Reading even MORE carefully:

The order of the dice digits doesn't matter, it's LITERALLY a sorted bitfield of which either 6 or 7 (if 6 or 9 are ticked) bits must be on and the rest off...

I was about to break out dynamic programming to spin up dice that might even have 6 on 6 sides if they rolled out that way, but nope, it's going to just be two counters and counting bits...


I spent far too many hours wondering why the filters (I wrote two) weren't reaching a correct value.
Not only are the individual dice ordered, the pairing of dice are ALSO ordered / unique.  The phrase "distinct arrangements of the two cubes" is doing a LOT of heavy lifting.



	var tens, ones [10]uint8
	var ii, ret, q, r uint16
	for ii = 1; ii < 10; ii++ {
		q = ii * ii
		q, r = q/10, q%10
		ones[r]++
		tens[q]++
	}
	ones[6] += ones[9]
	tens[6] += tens[9]
	// fmt.Println(ones)	//	[0 2 0 0 2 1 4 0 0 2]
	// fmt.Println(tens)	//	[3 1 1 1 1 0 1 0 1 0]	// Yeah, 6 NON 0 numbers...
	return ret


/
*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0090() uint64 {
	var ret uint64

	// Lookup table to speed up bit counts, it's supposed to easily fit in cache, and helps count 5 bits at a time precisely for this problem
	bitLUT := [32]uint8{}
	for ii := 0; ii < 32; ii++ {
		for t := ii; 0 < t; t >>= 1 {
			bitLUT[ii] += uint8(t & 1)
		}
	}

	sqp := [10]uint8{}
	for ii := 1; ii < 10; ii++ {
		q := ii * ii
		q, r := q/10, q%10
		sqp[ii] = uint8(q)<<4 | uint8(r)
	}

	bruteTest := func(A, B uint16, prt bool) bool {
		var mten, mone uint16
		a, b := A, B
		if 0 < a&0b10010_00000 {
			a |= 0b10010_00000
		}
		if 0 < b&0b10010_00000 {
			b |= 0b10010_00000
		}
		ret := true
		for ii := 1; ii < 10; ii++ {
			ten, one := sqp[ii]>>4, sqp[ii]&0xF
			mten = 1 << ten
			mone = 1 << one
			// Filter failed for 09 with dice: 11101 00101   101 11110
			if !((0 < a&mten && 0 < b&mone) || (0 < b&mten && 0 < a&mone)) {
				if prt {
					fmt.Printf("Filter failed for %d%d with dice: %10b %10b\n", ten, one, A, B)
				}
				ret = false
			}
		}
		return ret
	}
	// _ = bruteTest

	bruteTest_old := func(a, b uint16, prt bool) bool {
		var mten, mone uint16
		ret := true
		for ii := 0; ii < 9; ii++ {
			ten, one := sqp[ii]>>4, sqp[ii]&0xF
			if 6 == ten || 9 == ten {
				mten = 0b10010_00000
			} else {
				mten = 1 << ten
			}
			if 6 == one || 9 == one {
				mone = 0b10010_00000
			} else {
				mone = 1 << one
			}
			// Filter failed for 09 with dice: 11101 00101   101 11110
			if !((0 < a&mten && 0 < b&mone) || (0 < b&mten && 0 < a&mone)) {
				if prt {
					fmt.Printf("Filter failed for %d%d with dice: %10b %10b\n", ten, one, a, b)
				}
				ret = false
			}
		}
		return ret
	}
	_ = bruteTest_old

	isValid := func(a, b uint16, prt bool) bool {
		// if !((0 < a&0b01) && (0 < b&0b01)) {
		// if prt {
		// fmt.Printf("Filter reject: 0 + 0 dice: %10b %10b\n", a, b)
		// }
		// return false
		// }
		// 9~6 must be across from 4,3,1,0 (in the 0 set)
		if !((0 < a&0b10000 && 0 < b&0b10010_00000) || (0 < b&0b10000 && 0 < a&0b10010_00000)) {
			if prt {
				fmt.Printf("Filter reject: 6|9 + 4 dice: %10b %10b\n", a, b)
			}
			return false
		}
		if !((0 < a&0b1000 && 0 < b&0b10010_00000) || (0 < b&0b1000 && 0 < a&0b10010_00000)) {
			if prt {
				fmt.Printf("Filter reject: 6|9 + 3 dice: %10b %10b\n", a, b)
			}
			return false
		}
		if !((0 < a&0b10 && 0 < b&0b10010_00000) || (0 < b&0b10 && 0 < a&0b10010_00000)) {
			if prt {
				fmt.Printf("Filter reject: 6|9 + 1 dice: %10b %10b\n", a, b)
			}
			return false
		}
		// AND a 0 MUST be across from 1, 4 --6 or 9--
		if !((0 < a&0b01 && 0 < b&0b10010_00000) || (0 < b&0b01 && 0 < a&0b10010_00000)) {
			if prt {
				fmt.Printf("Filter reject: 0 + 9||6 dice: %10b %10b\n", a, b)
			}
			return false
		}
		if !((0 < a&0b01 && 0 < b&0b10000) || (0 < b&0b01 && 0 < a&0b10000)) {
			if prt {
				fmt.Printf("Filter reject: 0 + 4 dice: %10b %10b\n", a, b)
			}
			return false
		}
		if !((0 < a&0b01 && 0 < b&0b10) || (0 < b&0b01 && 0 < a&0b10)) {
			if prt {
				fmt.Printf("Filter reject: 0 + 1 dice: %10b %10b\n", a, b)
			}
			return false
		}
		// AND a 2 MUST be across from a 5 ++ NOTE: The problem doesn't mention 2 and 5 looking similar upside down, which some fonts (E.G. 7 seg display) might offer
		if !((0 < a&0b100 && 0 < b&0b1_00000) || (0 < b&0b100 && 0 < a&0b1_00000)) {
			if prt {
				fmt.Printf("Filter reject: 2 + 5 dice: %10b %10b\n", a, b)
			}
			return false
		}
		// AND an 8 MUST be across from a 1
		if !((0 < a&0b10 && 0 < b&0b1000_00000) || (0 < b&0b10 && 0 < a&0b1000_00000)) {
			if prt {
				fmt.Printf("Filter reject: 1 + 8 dice: %10b %10b\n", a, b)
			}
			return false
		}
		ret++
		return true
	}

	var considered uint64
	var a, b uint16
	for a = 0b1_11111; a <= 0b11111_10000; a++ {
		if 6 != bitLUT[a>>5]+bitLUT[a&0x1F] {
			continue
		}
		for b = a + 1; b <= 0b11111_10000; b++ {
			if 6 != bitLUT[b>>5]+bitLUT[b&0x1F] {
				continue
			}
			// Test and increase ret if it's valid
			considered++
			if isValid(a, b, false) != bruteTest(a, b, false) {
				fmt.Printf("Filter failed for dice: %10b %10b\n", a, b)
				bruteTest(a, b, true)
				isValid(a, b, true)
			}
		}
	}

	fmt.Printf("Cosidered: %d\n", considered)
	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 90 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Cosidered: 21945
Euler 90: Cube Digit Pairs: 1217

real    0m0.115s
user    0m0.166s
sys     0m0.064s
.
*/
func main() {
	var r uint64
	//test

	//run
	r = Euler0090()
	fmt.Printf("Euler 90: Cube Digit Pairs: %d\n", r)
	if 1217 != r {
		panic("Did not reach expected value.")
	}
}
