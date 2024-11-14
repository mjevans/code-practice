// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=84
https://projecteuler.net/minimal=84

<p>In the game, <strong>Monopoly</strong>, the standard board is set up in the following way:</p>
<div class="center">
<img src="resources/images/0084_monopoly_board.png?1678992052" alt="0084_monopoly_board.png">
</div>
<p>A player starts on the GO square and adds the scores on two 6-sided dice to determine the number of squares they advance in a clockwise direction. Without any further rules we would expect to visit each square with equal probability: 2.5%. However, landing on G2J (Go To Jail), CC (community chest), and CH (chance) changes this distribution.</p>
<p>In addition to G2J, and one card from each of CC and CH, that orders the player to go directly to jail, if a player rolls three consecutive doubles, they do not advance the result of their 3rd roll. Instead they proceed directly to jail.</p>
<p>At the beginning of the game, the CC and CH cards are shuffled. When a player lands on CC or CH they take a card from the top of the respective pile and, after following the instructions, it is returned to the bottom of the pile. There are sixteen cards in each pile, but for the purpose of this problem we are only concerned with cards that order a movement; any instruction not concerned with movement will be ignored and the player will remain on the CC/CH square.</p>
<ul><li>Community Chest (2/16 cards):
<ol><li>Advance to GO</li>
<li>Go to JAIL</li>
</ol></li>
<li>Chance (10/16 cards):
<ol><li>Advance to GO</li>
<li>Go to JAIL</li>
<li>Go to C1</li>
<li>Go to E3</li>
<li>Go to H2</li>
<li>Go to R1</li>
<li>Go to next R (railway company)</li>
<li>Go to next R</li>
<li>Go to next U (utility company)</li>
<li>Go back 3 squares.</li>
</ol></li>
</ul><p>The heart of this problem concerns the likelihood of visiting a particular square. That is, the probability of finishing at that square after a roll. For this reason it should be clear that, with the exception of G2J for which the probability of finishing on it is zero, the CH squares will have the lowest probabilities, as 5/8 request a movement to another square, and it is the final square that the player finishes at on each roll that we are interested in. We shall make no distinction between "Just Visiting" and being sent to JAIL, and we shall also ignore the rule about requiring a double to "get out of jail", assuming that they pay to get out on their next turn.</p>
<p>By starting at GO and numbering the squares sequentially from 00 to 39 we can concatenate these two-digit numbers to produce strings that correspond with sets of squares.</p>
<p>Statistically it can be shown that the three most popular squares, in order, are JAIL (6.24%) = Square 10, E3 (3.18%) = Square 24, and GO (3.09%) = Square 00. So these three most popular squares can be listed with the six-digit modal string: 102400.</p>
<p>If, instead of using two 6-sided dice, two 4-sided dice are used, find the six-digit modal string.</p>

/
*/
/*
	Monopoly board, but with each side as a line: (40 total spots)
	GO	A1	CC1	A2	T1	R1	B1	CH1	B2	B3
	JAIL	C1	U1	C2	C3	R2	D1	CC2	D2	D3
	FP	E1	CH2	E2	E3	R3	F1	F2	U2	F3
	Go2J	G1	G2	CC3	G3	R4	CH3	H1	T2	H2

	Specials:
		3 'doubles' in successive turns
		Go2Jail (100% go to jail)
		CC (Community Chest) : 1/16th Go to GO, 1/16th Go to Jail
		CH (Chance) 10/16 movement cards:
			Go to GO (#0)
			Go to R1 (#5)
			Go to Jail (#10)
			Go to C1 (#11)
			Go to E3 (#34)
			Go to H2 (#39)
			Go to 'Next RR' (Multiple ending in 5)
			Go to 'Next RR' (Multiple ending in 5)
			Go to Next U (utility, 2 spots #12, #28)
			Go 'back 3 squares'
		NOTE: Jail == 'Just Visiting' for this problem

Test case data:
	"Statistically it can be shown that the three most popular squares, in order, are JAIL (6.24%) = Square 10, E3 (3.18%) = Square 24, and GO (3.09%) = Square 00. So these three most popular squares can be listed with the six-digit modal string: 102400."

Real data: Figure out the most popular 3 squares if rolling with 2d4 instead of 2d6


Based on the wording, they don't care that advance by 1 isn't possible, nor are they thinking about epicycles caused by the folds backwards / forwards on Go2J...

2d4	1	2	3	4
1	2*	3	4	5
2	3	4*	5	6
3	4	5	6*	7
4	5	6	7	8*

Doubles: N / (N*N) or 1/N chance per roll, so 1/(N*N*N) for three rolls in a row.

	1.0 - (1/(N*N*N))	GoToJAIL
	R / 40			Raw Chance to every square, but apply #30's chance to #10 (jail) since it's Go 2 Jail
	3 cells are CC and re-distribute 1/16th of their value to GO and 1/16th to JAIL
	3 cells are CH and follow the table above...


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

func Euler0084(sides int) int {
	// Oversimplified 'Monopoly' board, but with each side as a line: (40 total spots)
	boardNames := []string{
	"GO",	"A1",	"CC1",	"A2",	"T1",	"R1",	"B1",	"CH1",	"B2",	"B3",
	"JAIL",	"C1",	"U1",	"C2",	"C3",	"R2",	"D1",	"CC2",	"D2",	"D3",
	"FP",	"E1",	"CH2",	"E2",	"E3",	"R3",	"F1",	"F2",	"U2",	"F3",
	"Go2J",	"G1",	"G2",	"CC3",	"G3",	"R4",	"CH3",	"H1",	"T2",	"H2",
	}
	_ = boardNames
	boardF := make([]float64, 40)

	// Doubles: N / (N*N) or 1/N chance per roll, so 1/(N*N*N) for three rolls in a row.
	r := float64(1.0) / float64(sides * sides * sides)
	boardF[10], r = r, (1.0 - r)/40.0

	// Add base chance to every square...
	for ii:=0 ; ii < 40 ; ii++ {
		boardF[ii] += r
	}

	// NOTE: Jail == 'Just Visiting' for this problem

	// Go to Jail square - (100% go to jail)
	boardF[10] += boardF[30]
	boardF[30] = 0.0

	// Community Chest, redirect% to Go and Jail - 		CC (Community Chest) : 1/16th Go to GO, 1/16th Go to Jail
	rcc := r / 16.0
	boardF[0] += rcc + rcc + rcc
	boardF[10] += rcc + rcc + rcc
	boardF[2] -= rcc + rcc
	boardF[17] -= rcc + rcc
	boardF[33] -= rcc + rcc

	// Chance redirect% CH (Chance) 10/16 movement cards:
	boardF[7] = r * 6.0 / 16.0 // 6 out of 16 times of landing on the square it's the end square
	boardF[22] = boardF[7]
	boardF[36] = boardF[7]
	// Go to GO (#0)
	boardF[0] += rcc + rcc + rcc
	// Go to R1 (#5) + CH3 RR x2
	boardF[5] += rcc + rcc + rcc + rcc + rcc
	// Go to Jail (#10)
	boardF[10] += rcc + rcc + rcc
	// Go to C1 (#11)
	boardF[11] += rcc + rcc + rcc
	// Go to E3 (#34)
	boardF[34] += rcc + rcc + rcc
	// Go to H2 (#39)
	boardF[39] += rcc + rcc + rcc
	//	Go to 'Next RR' (Multiple ending in 5)
	// 2x	Go to 'Next RR' (Multiple ending in 5)
	// CH3 R1 applied above
	boardF[15] += rcc + rcc
	boardF[25] += rcc + rcc
	// Go to Next U (utility, 2 spots #12, #28)
	boardF[12] += rcc + rcc
	boardF[28] += rcc
	// Go 'back 3 squares'
	boardF[4] += rcc
	boardF[29] += rcc
	boardF[33] += rcc * 14.0 / 16.0 // Wait, this is Community Chest!
	rcc /= 16.0
	boardF[0] += rcc // rollback go to go
	boardF[10] += rcc // rollback go to jail

	// Why is E3 a popular square in their example from the simplified game?  I expect Jail and Go to be the two most popular.  Are they not ignoring the epicycles after all?

	var ii0, ii1, ii2 int
	var ff0, ff1, ff2 float64

	for ii := 0 ; ii < 40 ; ii++ {
		if ff2 < boardF[ii] {
			ff2, ii2 = boardF[ii], ii
			if ff1 < ff2 {
				ff2, ff1, ii2, ii1 = ff1, ff2, ii1, ii2
				if ff0 < ff1 {
					ff1, ff0, ii1, ii0 = ff0, ff1, ii0, ii1
				}
			}
		}
	}

	// FIXME DANGER : The check informed me that, contrary to the basic implication, someone solving this IS expected to account for the epicycles... "Without any further rules we would expect to visit each square with equal probability: 2.5%. However, landing on G2J (Go To Jail), CC (community chest), and CH (chance) changes this distribution."
	if 6 == sides {
		ok := true
		if ! ( 10 == ii0 && (0.0623 < ff0 && ff0 < 0.0625)) {
			ok = false
			fmt.Printf("Expected First Place JAIL with 6.24%%, got %s (%d) with %f\n", boardNames[ii0], ii0, ff0)
		}
		if ! ( 24 == ii1 && (0.0317 < ff1 && ff1 < 0.0319)) {
			ok = false
			fmt.Printf("Expected Second Place E3 with 3.18%%, got %s (%d) with %f\n", boardNames[ii1], ii1, ff1)
		}
		if ! ( 0 == ii2 && (0.0308 < ff2 && ff2 < 0.0310)) {
			ok = false
			fmt.Printf("Expected Third Place GO with 3.09%%, got %s (%d) with %f\n", boardNames[ii2], ii2, ff2)
		}
		if !ok {
			panic("Check values failed, adjust replication of test case.")
		}
	}

	return 0
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 84 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

.
*/
func main() {
	var r int

	//test
	r = Euler0084(6)
	if 102400 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	r = Euler0084(4)
	fmt.Printf("Euler 84: Monopoly Odds: %d\n", r)
	if 427337 != r {
		panic("Did not reach expected value.")
	}
}
