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


--- NOPE tests disproved this --- Based on the wording, they don't care that advance by 1 isn't possible, nor are they thinking about epicycles caused by the folds backwards / forwards on Go2J...

2dx	1	2	3	4	5	6
1	2*	3	4	5	6	7
2	3	4*	5	6	7	8
3	4	5	6*	7	8	9
4	5	6	7	8*	9	10
5	6	7	8	9	10*	11
6	7	8	9	10	11	12*

Doubles: N / (N*N) or 1/N chance per roll, so 1/(N*N*N) for three rolls in a row.

	1.0 - (1/(N*N*N))	GoToJAIL
	R / 40			Raw Chance to every square, but apply #30's chance to #10 (jail) since it's Go 2 Jail
	3 cells are CC and re-distribute 1/16th of their value to GO and 1/16th to JAIL
	3 cells are CH and follow the table above...

The question in my mind is now, how to deal with the epicycles?
	* Do they make the total across the board greater than 100% (1.0)?  I think I'm going to tentatively say yes, since the literal magnitude doesn't matter THOUGH, I should decrease later cell values by the leaked epicycles, even as the epicycles are added back to the board.
	* Maybe a recursion limit counter?

	Need to step outside of the code for a moment and workshop the logic 'aloud'...
	This was based on a static probability where everything happened at once and thus I didn't need to consider if E.G. landing on a given square happened.  Every square just started with an equal slice of the pie.
	Now I'm trying to use ONE loop around the board, mod board, to apply that same probability but recurse loop starting on the epicycle when an event sends the cursor to a new location.
	However that doesn't happen _every_ trip around the board, it only happens IF that cell is the one that's landed upon.
	How often is a cell landed on?
	With infinite loops, 1/n grows to an even share of the pie, without any detours considered...
	So when it's passed the probability of that detour needs to be calculated.
	1/n times that it's passed, it will be landed on, and then the probability will strike.

	The check values indicate I'm possibly further from the answer trying to figure out probability with just one circle around the board and epicycles.  There's probably some higher math way of calculating the probability of dice rolls and factoring in how short rolls involve more rolls to progress around the board.  There's also short rolls is more likely to jump to jail.

	I think I might need to go for drastic measures and do this one the full brute force way.  Just play it and record visit counts before calculating.


	// Attempt #3 convolution (discrete signal processing style) of dice rolls (Markov Chain of some sort?) //

	Search : board game probability algorithm
	https://stackoverflow.com/a/73353222
	https://en.wikipedia.org/wiki/Monte_Carlo_method
	https://en.wikipedia.org/wiki/Markov_chain_Monte_Carlo		I've heard of this as a 'pre LLM' chat bot algorithm, such as this result https://rosettacode.org/wiki/Markov_chain_text_generator
	https://en.wikipedia.org/wiki/Markov_chain
	https://in-thread.sonic-pi.net/t/markov-chains-for-beginners/5304

	I wonder if I could apply the statistics of the dice rolls to going around the board, sort of like a fusion of my first really simple pass that didn't account for dice rolls at all because they'd tend to land on every square somewhat evenly over infinite iterations, and like the diversions epicycle loop that isn't walking correctly, but with a statistical shadow of it.

	Initial state: (for 100% and the simplest initial representation) 1.0 on the GO square.  Each round convolve (like signal processing) the entire board by the dice probability matrix, with an added adjustment for the three in a row doubles goes to jail.  Probably OK to stop when the shortest set of steps would form a loop.

	The algorithm I'm using seems to need 3-4 loops to sufficiently stabilize.  I'd though the even/odd settling would likely happen after 2 loops, but it makes sense that each additional power of 2 loops would increase precision.

	My value for JAIL is a tiny bit sweet compared to Euler 84's reference.  I attribute this to calculating dice roll go to jail first and reducing the overall outcome after that.
	After 4 rounds the total probability across the board has decayed to 0.9999999999999966 (for 2d6) which isn't sufficient to explain the discrepancy otherwise.
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

func Euler0084(sides, rlimit int) int {
	// Oversimplified 'Monopoly' board, but with each side as a line: (40 total spots)
	boardNames := []string{
		"GO", "A1", "CC1", "A2", "T1", "R1", "B1", "CH1", "B2", "B3",
		"JAIL", "C1", "U1", "C2", "C3", "R2", "D1", "CC2", "D2", "D3",
		"FP", "E1", "CH2", "E2", "E3", "R3", "F1", "F2", "U2", "F3",
		"G2J", "G1", "G2", "CC3", "G3", "R4", "CH3", "H1", "T2", "H2",
	}
	_ = boardNames
	boardF, prior, diceprob := make([]float64, 40), make([]float64, 40), make([]float64, (sides<<1)+1)
	var check float64
	_ = check

	// Initialize the dice probability pattern, [0] is stay in place (0%) and [1] is roll a 1 (also 0% for 2 dice labeled 1..n)
	for ii := 1; ii <= sides; ii++ {
		p := float64(ii) / float64(sides*sides)
		diceprob[ii+1], diceprob[(sides<<1)+1-ii] = p, p
	}
	// fmt.Println(diceprob)

	var visit func(ii int, p float64)
	visit = func(ii int, p float64) {
		for 40 <= ii {
			ii -= 40
		}
		switch ii {
		case 2, 17, 33:
			// Community Chest, redirect% to Go and Jail - CC (Community Chest) : 1/16th Go to GO, 1/16th Go to Jail
			pcc := p / 16.0
			boardF[0] += pcc  // GO
			boardF[10] += pcc // JAIL
			boardF[ii] += p - (pcc * 2.0)
		case 7, 22, 36:
			// Chance redirect% CH (Chance) 10/16 movement cards
			pch := p / 16.0
			boardF[0] += pch  // Go to GO (#0)
			boardF[5] += pch  // Go to R1 (#5)
			boardF[10] += pch // Go to Jail (#10)
			boardF[11] += pch // Go to C1 (#11)
			boardF[24] += pch // Go to E3 (#34)
			boardF[39] += pch // Go to H2 (#39)
			//	Go to 'Next RR' (Multiple ending in 5)
			// 2x	Go to 'Next RR' (Multiple ending in 5)
			switch ii {
			case 7:
				boardF[15] += pch // RR #1
				boardF[15] += pch // RR #2
				boardF[12] += pch // U
			case 22:
				boardF[25] += pch // RR #1
				boardF[25] += pch // RR #2
				boardF[28] += pch // U
			case 36:
				boardF[5] += pch  // RR #1
				boardF[5] += pch  // RR #2
				boardF[12] += pch // U
			}
			visit(ii-3, pch) // Go 'back 3 squares', this has a chance of hitting another special case (33)
			//
			boardF[ii] += p - (pch * 10.0)
		case 30:
			boardF[10] += p
		default:
			boardF[ii] += p
		}
	}

	var rollvisit func(reduce float64)
	rollvisit = func(reduce float64) {
		prior, boardF = boardF, prior
		for ii := 0; ii < 40; ii++ {
			boardF[ii] = 0.0 // reset the board
		}
		for ii := 0; ii < 40; ii++ {
			if 0.0 == prior[ii] {
				continue
			}
			// 0 and 1 are known to be 0.0, they just exist to make it easier for humans to remember what the roll of the dice was and to project forward in ii correctly
			for jj := 2; jj <= sides<<1; jj++ {
				visit(ii+jj, prior[ii]*diceprob[jj]*reduce)
			}
		}
	}

	// Initial state, pawn at go, will be swapped to the prior array in rollvisit()
	boardF[0] = 1.0

	// first two moves have 0% chance of going to jail on doubles, special case
	rollvisit(1.0)

	// fmt.Println(boardF)
	// check = 0.0
	// for ii := 0; ii < 40; ii++ {
	//	check += boardF[ii]
	// }
	// fmt.Println(check)

	rollvisit(1.0)

	// MINIMUM one trip around the board
	if 1 > rlimit {
		rlimit = 1
	}
	rlimit *= 20 // smallest possible steps always (2) ignoring jail (no doubles) and redirects, fully around the board once to fully spread the probability
	// Doubles: N / (N*N) or 1/N chance per roll, so 1/(N*N*N) for three rolls in a row.
	chance := float64(1.0) / float64(sides*sides*sides)
	reduce := 1.0 - chance
	for ii := 0; ii < rlimit; ii++ {
		rollvisit(reduce)
		boardF[10] += chance // this must be after rollvisit projects the new probability result
	}

	// Why is E3 a popular square in their example from the simplified game?  I expect Jail and Go to be the two most popular.  Are they not ignoring the epicycles after all?

	var ii0, ii1, ii2 int
	var ff0, ff1, ff2 float64

	for ii := 0; ii < 40; ii++ {
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

	// The check informed me that, contrary to the basic implication, someone solving this IS expected to account for the epicycles... "Without any further rules we would expect to visit each square with equal probability: 2.5%. However, landing on G2J (Go To Jail), CC (community chest), and CH (chance) changes this distribution."
	if 6 == sides {
		ok := true
		if !(10 == ii0 && (0.0623 < ff0 && ff0 < 0.0631)) {
			ok = false
			fmt.Printf("Expected First Place JAIL with 6.24%%, got %s (%d) with %f\n", boardNames[ii0], ii0, ff0)
		}
		if !(24 == ii1 && (0.0317 < ff1 && ff1 < 0.0319)) {
			ok = false
			fmt.Printf("Expected Second Place E3 with 3.18%%, got %s (%d) with %f\n", boardNames[ii1], ii1, ff1)
		}
		if !(0 == ii2 && (0.0308 < ff2 && ff2 < 0.0310)) {
			ok = false
			fmt.Printf("Expected Third Place GO with 3.09%%, got %s (%d) with %f\n", boardNames[ii2], ii2, ff2)
		}
		if !ok {
			panic("Check values failed, adjust replication of test case.")
		} else {
			fmt.Printf("First Place JAIL with 6.24%%, got %s (%d) with %f\n", boardNames[ii0], ii0, ff0)
			fmt.Printf("Second Place E3 with 3.18%%, got %s (%d) with %f\n", boardNames[ii1], ii1, ff1)
			fmt.Printf("Third Place GO with 3.09%%, got %s (%d) with %f\n", boardNames[ii2], ii2, ff2)
		}
	}
	// check = 0.0
	// for ii := 0; ii < 40; ii++ {
	//	check += boardF[ii]
	// }
	// fmt.Println(check)

	return ii0*10000 + ii1*100 + ii2
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 84 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

First Place JAIL with 6.24%, got JAIL (10) with 0.063089
Second Place E3 with 3.18%, got E3 (24) with 0.031841
Third Place GO with 3.09%, got GO (0) with 0.030885
Euler 84: Monopoly Odds: 101524

real    0m0.097s
user    0m0.136s
sys     0m0.055s
.
*/
func main() {
	var r int

	//test
	r = Euler0084(6, 4)
	if 102400 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	r = Euler0084(4, 4)
	fmt.Printf("Euler 84: Monopoly Odds: %d\n", r)
	if 101524 != r {
		panic("Did not reach expected value.")
	}
}
