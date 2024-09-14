// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=28

<p>Starting with the number $1$ and moving to the right in a clockwise direction a $5$ by $5$ spiral is formed as follows:</p>
<p class="monospace center"><span class="red"><b>21</b></span> 22 23 24 <span class="red"><b>25</b></span><br>
20  <span class="red"><b>7</b></span>  8  <span class="red"><b>9</b></span> 10<br>
19  6  <span class="red"><b>1</b></span>  2 11<br>
18  <span class="red"><b>5</b></span>  4  <span class="red"><b>3</b></span> 12<br><span class="red"><b>17</b></span> 16 15 14 <span class="red"><b>13</b></span></p>
<p>It can be verified that the sum of the numbers on the diagonals is $101$.</p>
<p>What is the sum of the numbers on the diagonals in a $1001$ by $1001$ spiral formed in the same way?</p>

... So within the 5x5 in the center... 101 == 21 + 7 + 1 + 3 + 13 + 17 + 5 + 9 + 25

	73	74	75	76	77	78	79	80	81
	72	43	44	45	46	47	48	49	50
	71	42	21	22	23	24	25	26	51
	70	41	20	7	8	9	10	27	52
	69	40	19	6	1	2	11	28	53
	68	39	18	5	4	3	12	29	54
	67	38	17	16	15	14	13	30	55
	66	37	36	35	34	33	32	31	56
	65	64	63	62	61	60	59	58	57

 25 == 7 + 1 + 3 + 5 + 9
101 == 21 + 7 + 1 + 3 + 13 + 17 + 5 + 9 + 25
261 == 43 + 21 + 7 + 1 + 3 + 13 + 31 + 37 + 17 + 5 + 9 + 25 + 49
537 == 73 + 43 + 21 + 7 + 1 + 3 + 13 + 31 + 57 + 65 + 37 + 17 + 5 + 9 + 25 + 49 + 81

Euler 26 reminded me, if I don't understand what's going on, make a logic / states table and see what patterns exist.

step	Grid	Sum	RingSZ	Area	Corners	lr	ll	ul	ur
0	1	1	1	1	1
1	3	25	8	9	24	3	5	7	9
2	5	101	16	25	79	13	17	21	25
3	7	261	24	49	160	31	37	43	49
4	9	537	32	81	276	57	65	73	81

Observations...

Upper Right == 'grid' * 'grid'
'Corners' == 4 ('grid' * 'grid') - 6 * ('grid' - 1)  // Where's the 6 come from?  That's the 1, 2 and 3 'back' the other number of elements on each side

I'm not a Maths major, but if someone with that background were on the team, or I needed to work it out on my own, I'd investigate series equation transformations.  There might be some way of making that faster...

As this is, with ~500 steps of N operations (bitshift by 2 ; two muls, two subtracts) I think this sufficiently solvable.


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

func Euler028(outer int) int {
	ret := 1
	if 1 == outer {
		return ret
	}
	if 1 > outer || 1 != outer&0x1 {
		return 0
	}
	for ii := 3; ii <= outer; ii += 2 {
		ret += ii*ii<<2 - 6*(ii-1)
	}
	return ret
}

// for ii in */*.go ; do go fmt "$ii" ; done ; for ii in 28 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
/*

Euler028: test:  true  Grid  0  x  0  ==  0
Euler028: test:  true  Grid  1  x  1  ==  1
Euler028: test:  true  Grid  3  x  3  ==  25
Euler028: test:  true  Grid  5  x  5  ==  101
Euler028: test:  true  Grid  7  x  7  ==  261
Euler028: test:  true  Grid  9  x  9  ==  537
Euler028: Result:               Grid  1001  x  1001  ==  669171001

*/
func main() {
	//test
	test := []int{0, 0, 1, 1, 3, 25, 5, 101, 7, 261, 9, 537}
	for ii := 0; ii < len(test); ii += 2 {
		sum := Euler028(test[ii])
		fmt.Println("Euler028: test:\t", sum == test[ii+1], " Grid ", test[ii], " x ", test[ii], " == ", sum)
	}

	//run
	grid := 1001
	sum := Euler028(grid)
	fmt.Println("Euler028: Result:\t\tGrid ", grid, " x ", grid, " == ", sum)

}
