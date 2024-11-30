// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=96
https://projecteuler.net/minimal=96

<p>Su Doku (Japanese meaning <i>number place</i>) is the name given to a popular puzzle concept. Its origin is unclear, but credit must be attributed to Leonhard Euler who invented a similar, and much more difficult, puzzle idea called Latin Squares. The objective of Su Doku puzzles, however, is to replace the blanks (or zeros) in a 9 by 9 grid in such that each row, column, and 3 by 3 box contains each of the digits 1 to 9. Below is an example of a typical starting puzzle grid and its solution grid.</p>
<div class="center">
<img src="project/images/p096_1.png" alt="p096_1.png" />     <img src="project/images/p096_2.png" alt="p096_2.png" /></div>
<p>A well constructed Su Doku puzzle has a unique solution and can be solved by logic, although it may be necessary to employ "guess and test" methods in order to eliminate options (there is much contested opinion over this). The complexity of the search determines the difficulty of the puzzle; the example above is considered <i>easy</i> because it can be solved by straight forward direct deduction.</p>
<p>The 6K text file, <a href="project/resources/p096_sudoku.txt">sudoku.txt</a> (right click and 'Save Link/Target As...'), contains fifty different Su Doku puzzles ranging in difficulty, but all with unique solutions (the first puzzle in the file is the example above).</p>
<p>By solving all fifty puzzles find the sum of the 3-digit numbers found in the top left corner of each solution grid; for example, 483 is the 3-digit number found in the top left corner of the solution grid above.</p>


/
*/
/*
	Automated sudoku solver, capture the upper left 3 cells and sum (50x999 max)
	A check solution for the first puzzle has been provided, which should reach 483 there and which has... a bunch of numbers I should record for the test case.
	This is what one puzzle looks like in the file; I want to pass around uint8s but otherwise... yeah.
Grid 01
003020600
900305001
001806400
008102900
700000008
006708200
002609500
800203009
005010300
	There's not a TON of point in compressing the numbers into nibbles, and uint8 is slightly too small to store a bitmask of possible numbers.
	I'm unsure if I want to store more than an array of arrays for the grid, I don't think I care to.
	I will probably just store [81]uint16 as the array and let math figure out the indexes for columns and 3x3s
	It might improve the solver(s) / efficiency if I have some dirty flag bitmaps; probably just one per row, column and 3x3 cell (every logic domain).
	Need to solve generally before turning that on though.

	In addition to reduction that identifies 'only X can go here' by eliminating any options that would conflict...
	Out of the open cells in this set, 'An N can only go in (one) Cell' test is required (for many but not all puzzles).

	I considered that 7 bits * 9 numbers is 63 bits; it could just _barely_ be compressed into a 64 bit number.
	However that would add code complexity for a premature 'optimization', and only saves 5 bytes of storage.

		 0  1  2   3  4  5   6  7  8
		 9 10 11  12 13 14  15 16 17
		18 19 20  21 22 23  24 25 26

		27 28 29  30 31 32  33 34 35
		36 37 38  39 40 41  42 43 44
		45 46 47  48 49 50  51 52 53

		54 55 56  57 58 59  60 61 62
		63 64 65  66 67 68  69 70 71
		72 73 74  75 76 77  78 79 80

	The second puzzle appears to require a guess to solve.
	In golang that pushes me to pass a copy of the puzzle by value (full copy) instead of pointer reference.
	The guess can be made in the cell that has the lowest popcount (including the unsure 0 bitflag), and among values in an increasing order...
	Or should it be made in a cell with the most remaining entropy?  That could be valid too, but I think the least entropy cell is the most likely to return a conflict sooner.

Iter     0      51 remain
found co [22] = 6
Iter     4      50 remain
            .   1010110011   1010100001   1000010001            .   1000010011            .     10100011     11100011
   1000100011            .   1000100001   1000001101            .   1000001111       100011            .            .
    100000011            .    110000001            .            .        10011            .     10000011            .
   1001000001   1010000001   1011001001            .   1000001101            .            .   1010001101            .
   1101100011   1110100011   1111101001   1100011101   1000001101   1100011101   1011100011   1010101111     11100111
            .   1100100011            .            .   1000001001            .   1000100011   1000101011       100011
            .   1100100001            .   1101000101   1000100101            .   1100100001            .       100101
            .            .   1100100001   1100001001            .   1100001001   1100100011            .       100011
   1101100001   1100100001            .   1101000101            .   1100000101   1110100001   1010100101            .
 2 . . . 8 . 3 . .
 . 6 . . 7 . . 8 4
 . 3 . 5 6 . 2 . 9
 . . . 1 . 5 4 . 8
 . . . . . . . . .
 4 . 2 7 . 6 . . .
 3 . 1 . . 7 . 4 .
 7 2 . . 4 . . 6 .
 . . 4 . 1 . . . 3
panic: could not solve

	Guessing isn't quite working and has ballooned an initially barely simple enough core to something that should be broken up.
	Maybe a SuDoKu object that has methods so the shared state can more easily be shared among functions?

	Added some tracking and fixed two bugs caused by a missed element of state tracking.  The complexity is definitely too high.  As much as I'd hoped to avoid a refactor after sleep that does still look required.

	Number three appears to have a structure bare enough to force a full grid entropy consideration, so maybe I _should_ refactor how that's calculated too.

	Refactor effectively complete, though guesses still incorrect.

	The most silly mistake I made, likely due to working on the design over multiple nights, was momentarily forgetting that 0 is a valid cell address, and thus not a good value to use for already seen things.  Got that wrong once on a late bolt-on step; was a nightmare to debug given the rarity of the test case and how easily just one single wrong cell can drive the entire grid off from the value it should have.

	When this started out with just three simple-ish mask and mark sections that was, thanks to three distinct sections, just barely at the level before where a refactor should have happened.
	If I'd refactored at that point instead of trying to add in multiple other complex functionality additions at the same time it's less likely I'd have made a simple mistake like that.

	Overall, this might be the least fun Euler puzzle I've done so far, though one other was nightmarish unfun too.
	However this puzzle at least has redeeming general programming education and practice value.
	It's such a good trap for something that seems deceptively simple, but SHOULD be broken out into easier to digest units of logic sooner.
	The trap is how interconnected all the logic is, which _really_ encourages a programmer to keep things within one function / object / closure if avoiding global variables.

/
*/

import (
	"bufio"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	"os" // os.Stdout
	// "strconv"
	// "strings"
	"container/heap"
)

var bitLUT [32]uint8 // init() to lookup of bit popcount for 5 bits

// https://pkg.go.dev/container/heap@go1.22.6#Pop
// required widths: 4 for cell value, 7 for cell address
// Pack := popcount << 12 | val << 8 | celladdr << 0
type Uint32Queue []uint32

func (uq Uint32Queue) Raw() []uint32 {
	conv := ([]uint32)(uq)
	return conv
}

func (uq Uint32Queue) Len() int { return len(uq) }

func (uq Uint32Queue) Less(queA, queB int) bool {
	// "less" holds items closer to the base of the array
	return uq[queA] < uq[queB]
}

func (uq Uint32Queue) Swap(queA, queB int) {
	uq[queA], uq[queB] = uq[queB], uq[queA]
	// 'Item' lacks priority and lacks index
}

func (uq *Uint32Queue) Push(fp any) {
	*uq = append(*uq, fp.(uint32))
}

func (uq *Uint32Queue) Pop() any {
	n := len(*uq) - 1
	fp := (*uq)[n]
	*uq = (*uq)[0:n]
	return fp
}

type SuDoKu struct {
	DirtyRow, DirtyCol, DirtyBox uint16
	Rem, State                   uint8 // State 0 unset, 1 ArraysValid, ...?
	Note                         [81]uint16
	Num                          [81]uint8
	Only                         [9]uint8
	EntRow, EntCol, EntBox       [9]uint8 // popcount - 9 for every set
}

func (s *SuDoKu) Solve() uint8 {
	s.ResetAllNotes()
	s.SolveInner()
	if 0 != s.Validate(true) {
		s.Print()
		panic("Validate isolated issues")
	}
	return s.Rem
}

func (s *SuDoKu) ResetAllNotes() {
	s.Rem, s.State, s.DirtyRow, s.DirtyCol, s.DirtyBox = 81, 1, 0x1FF, 0x1FF, 0x1FF // All Dirty
	for ii := 0; ii < 81; ii++ {
		if 0 == s.Num[ii] {
			s.Note[ii] = 0b111_111_111_1
		} else {
			s.Note[ii], s.Rem = 1<<s.Num[ii], s.Rem-1
		}
	}
}

func (s *SuDoKu) Validate(info bool) uint8 {
	var bigrow, bigcol, brb, bcb, row, rb, col, cell, ret uint8

	for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
		// Scan known numbers
		seen := uint16(0)
		for col = 0; col < 9; col++ {
			cell = rb + col
			if 0 != s.Num[cell] {
				if 0 != seen&(1<<(s.Num[cell])) {
					if info {
						fmt.Printf("Duplicate %d found in cell %2d\n", s.Num[cell], cell)
					}
					ret = 0xFF
				} else {
					seen |= 1 << (s.Num[cell])
				}
			}
		}
		if info && 0b111_111_111_0 != seen {
			fmt.Printf("Missing numbers in row %d 0b%10b\n", row, seen)
		}
	}

	for col = 0; col < 9; col++ {
		// Scan known numbers
		seen := uint16(0)
		for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
			cell = rb + col
			if 0 != s.Num[cell] {
				if 0 != seen&(1<<(s.Num[cell])) {
					if info {
						fmt.Printf("Duplicate %d found in cell %2d\n", s.Num[cell], cell)
					}
					ret = 0xFF
				} else {
					seen |= 1 << (s.Num[cell])
				}
			}
		}
		if info && 0b111_111_111_0 != seen {
			fmt.Printf("Missing numbers in col %d 0b%10b\n", col, seen)
		}
	}

	for bigrow, brb = 0, 0; bigrow < 3; bigrow, brb = bigrow+1, brb+27 {
		for bigcol, bcb = 0, 0; bigcol < 3; bigcol, bcb = bigcol+1, bcb+3 {
			// Scan known numbers
			seen := uint16(0)
			for row, rb = 0, 0; row < 3; row, rb = row+1, rb+9 {
				for col = 0; col < 3; col++ {
					cell = brb + bcb + rb + col
					if 0 != s.Num[cell] {
						if 0 != seen&(1<<(s.Num[cell])) {
							if info {
								fmt.Printf("Duplicate %d found in cell %2d\n", s.Num[cell], cell)
							}
							ret = 0xFF
						} else {
							seen |= 1 << (s.Num[cell])
						}
					}
				}
			}
			if info && 0b111_111_111_0 != seen {
				fmt.Printf("Missing numbers in box %d 0b%10b_\n", bigrow*3+bigcol, seen)
			}
		}
	}
	return ret
}

// Update the cell's mask, if it's still variable.  Return 0 for no updates, 10 for mask only, Number if set to a number, and 0xFF if error/no matches
func (s *SuDoKu) UpdateCellMask(seen uint16, cell uint8) uint8 {
	var ret uint8
	if 1 == s.Note[cell]&1 {
		possible := s.Note[cell] & seen
		if 1 == possible {
			// fmt.Printf("No possible values for [%d]\n", cell)
			return 0xFF
		}
		if possible != s.Note[cell] { // MUST process to evaluate s.Only numbers
			ret = 10
		}
		pscan, justone, testnum := possible>>1, uint8(0), uint8(1) // shift out the zero flag
		// scan for if this is the only N, or contested
		for ; 0 < pscan; pscan, testnum = pscan>>1, testnum+1 {
			if 0 != pscan&1 {
				if 0 == justone {
					justone = testnum
				} else {
					justone = 0xFF
				}
				if 0xFE == s.Only[testnum-1] {
					s.Only[testnum-1] = cell
				} else {
					s.Only[testnum-1] = 0xFF // mark contested
				}
				// if 2 == testnum { // && 8 == cell%9 {
				// fmt.Printf("[%2d] ~~ %d and Only now [%2d]\n", cell, testnum, s.Only[testnum-1])
				// }
			}
		}
		if 0 < justone && justone < 0xFF {
			possible, ret, s.Num[cell], s.Only[justone-1], s.Rem = 1<<justone, justone, justone, 0xFF, s.Rem-1
			// fmt.Printf("Set JustOne [%2d] = %d\n", cell, justone)
		}
		s.Note[cell] = possible
		// }
	}
	return ret
}

func (s *SuDoKu) numberDiscovered() {
	// fmt.Printf("%v\n", s.Only)
	for ii := uint8(0); ii < 9; ii++ {
		if 0 < s.Only[ii] && s.Only[ii] < 0xFE {
			cell, val := s.Only[ii], ii+1
			// fmt.Printf("found nD [%2d] = %d\n", cell, val)
			s.Num[cell], s.Note[cell], s.Rem = val, 1<<val, s.Rem-1
			s.DirtyRow |= 1 << (cell / 9)                       // recover row from cell co-ord
			s.DirtyCol |= 1 << (cell % 9)                       // recover col from cell co-ord
			s.DirtyBox |= 1 << ((cell/27)*3 + ((cell % 9) / 3)) // recover box from cell co-ord
		}
	}

}

func (s *SuDoKu) ReduceRow() uint8 {
	// skip := s.DirtyRow
	s.DirtyRow = 0
	var row, rb, col, cell, entpopc uint8
	for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
		// if 0 == skip&(1<<row) {
		// continue
		// }
		// fmt.Printf("ReduceRow")

		// Scan known numbers
		seen := uint16(0)
		s.Only = [9]uint8{0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE}

		for col = 0; col < 9; col++ {
			cell = rb + col
			// fmt.Printf("\t%2d", cell)
			if 0 == s.Note[cell]&1 {
				seen |= s.Note[cell]
				s.Only[s.Num[cell]-1] = 0xFF
			}
		}
		seen = (^seen) & 0x3FF // Any number that was NOT seen is allowed...
		// fmt.Printf("\t  %10b\nReduceRow", seen)

		// MASK known numbers out of unknown numbers
		entpopc = 0
		for col = 0; col < 9; col++ {
			cell = rb + col
			// fmt.Printf("\t%2d", cell)
			entpopc += bitLUT[s.Note[cell]>>5] + bitLUT[s.Note[cell]&0x1F]
			cval := s.UpdateCellMask(seen, cell)
			if 0xFF == cval {
				return 0xFF // Failed update / guess
			}
			if 0 < cval && cval < 10 {
				s.DirtyRow |= 1 << row
				seen &^= 1 << cval
			}
			if 0 < cval {
				s.DirtyCol |= 1 << col
				s.DirtyBox |= 1 << (((row / 3) * 3) + (col / 3))
			}
		}
		s.EntRow[row] = entpopc - 9
		// fmt.Printf("\t%d\n", entpopc-9)
		s.numberDiscovered()
	}
	return 0
}

func (s *SuDoKu) ReduceCol() uint8 {
	// skip := s.DirtyCol
	s.DirtyCol = 0
	var row, rb, col, cell, entpopc uint8
	for col = 0; col < 9; col++ {
		// if 0 == skip&(1<<col) {
		// continue
		// }
		// fmt.Printf("ReduceCol")

		// Scan known numbers
		seen := uint16(0)
		s.Only = [9]uint8{0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE}
		for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
			cell = rb + col
			// fmt.Printf("\t%2d", cell)
			if 0 == s.Note[cell]&1 {
				seen |= s.Note[cell]
				s.Only[s.Num[cell]-1] = 0xFF
			}
		}
		seen = (^seen) & 0x3FF // Any number that was NOT seen is allowed...
		// fmt.Printf("\t  %10b\nReduceCol", seen)

		// MASK known numbers out of unknown numbers
		entpopc = 0
		for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
			cell = rb + col
			// fmt.Printf("\t%2d", cell)
			entpopc += bitLUT[s.Note[cell]>>5] + bitLUT[s.Note[cell]&0x1F]
			cval := s.UpdateCellMask(seen, cell)
			if 0xFF == cval {
				return 0xFF // Failed update / guess
			}
			if 0 < cval && cval < 10 {
				s.DirtyCol |= 1 << col
				seen &^= 1 << cval
			}
			if 0 < cval {
				s.DirtyRow |= 1 << row
				s.DirtyBox |= 1 << (((row / 3) * 3) + (col / 3))
			}
		}
		s.EntCol[col] = entpopc - 9
		// fmt.Printf("\t%d\n", entpopc-9)
		s.numberDiscovered()
	}
	return 0
}

func (s *SuDoKu) ReduceBox() uint8 {
	// skip := s.DirtyBox
	s.DirtyBox = 0
	var bigrow, bigcol, brb, bcb, row, rb, col, cell, entpopc uint8
	for bigrow, brb = 0, 0; bigrow < 3; bigrow, brb = bigrow+1, brb+27 {
		for bigcol, bcb = 0, 0; bigcol < 3; bigcol, bcb = bigcol+1, bcb+3 {
			// if 0 == skip&(1<<(bigrow*3+bigcol)) {
			// continue
			// }
			// fmt.Printf("ReduceBox")

			// Scan known numbers
			seen := uint16(0)
			s.Only = [9]uint8{0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE}
			for row, rb = 0, 0; row < 3; row, rb = row+1, rb+9 {
				for col = 0; col < 3; col++ {
					cell = brb + bcb + rb + col
					// fmt.Printf("\t%2d", cell)
					if 0 == s.Note[cell]&1 {
						seen |= s.Note[cell]
						s.Only[s.Num[cell]-1] = 0xFF
					}
				}
			}
			seen = (^seen) & 0x3FF // Any number that was NOT seen is allowed...
			// fmt.Printf("\t  %10b\nReduceBox", seen)

			// MASK known numbers out of unknown numbers
			entpopc = 0
			for row, rb = 0, 0; row < 3; row, rb = row+1, rb+9 {
				for col = 0; col < 3; col++ {
					cell = brb + bcb + rb + col
					// fmt.Printf("\t%2d", cell)
					entpopc += bitLUT[s.Note[cell]>>5] + bitLUT[s.Note[cell]&0x1F]
					cval := s.UpdateCellMask(seen, cell)
					if 0xFF == cval {
						return 0xFF // Failed update / guess
					}
					if 0 < cval && cval < 10 {
						s.DirtyBox |= 1 << (bigrow*3 + bigcol)
						seen &^= 1 << cval
					}
					if 0 < cval {
						s.DirtyCol |= 1 << (bigcol*3 + col)
						s.DirtyCol |= 1 << (bigrow*3 + row)
					}
				}
			}
			s.EntCol[col] = entpopc - 9
			// fmt.Printf("\t%d\n", entpopc-9)
			s.numberDiscovered()
		}
	}
	return 0
}

func (s *SuDoKu) Print() {
	var row, rb, col, cell uint8
	fmt.Println()
	for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
		for col = 0; col < 9; col++ {
			cell = rb + col
			if 1 == bitLUT[s.Note[cell]>>5]+bitLUT[s.Note[cell]&0x1F] && 0 != s.Num[cell] {
				fmt.Print("            .")
			} else {
				fmt.Printf("   %10b", s.Note[cell])
			}
		}
		fmt.Println()
	}
	for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
		for col = 0; col < 9; col++ {
			cell = rb + col
			if 0 == s.Num[cell] {
				fmt.Print(" .")
			} else {
				fmt.Printf(" %d", s.Num[cell])
			}
		}
		fmt.Println()
	}
}

func (s *SuDoKu) GuessRequired() uint8 {
	// fmt.Printf("Guess Required: Remaining: %d\n", s.Rem)

	// Which has the lowest entropy?  Try all three paths until exhausted
	var gl Uint32Queue
	guesses := int(0)
	for ii := 0; ii < 9; ii++ {
		guesses += int(s.EntBox[ii]) // Box pass was last so should be the best estimate
	}
	gl = make([]uint32, 0, guesses)
	// Euler 96 puzzle 3 appears to require more than just the most likely spots in the top location for each iteration style
	var packed uint32
	var cell, row, col, box, ent, bestGscore, bestGcell, bestGnum uint8
	bestGscore = 0xFF
	for cell = 0; cell < 81; cell++ {
		if 0 != s.Note[cell]&1 {
			row, col, box = cell/9, cell%9, (cell/27)*3+((cell%9)/3)                                                  // from cell
			ent = s.EntRow[row] + s.EntCol[col] + s.EntBox[box] + bitLUT[s.Note[cell]>>5] + bitLUT[s.Note[cell]&0x1F] // Worst Case: 253 = 81 * 3 + 10
			packed = uint32(ent)<<16 | uint32(cell)
			for note, val := s.Note[cell]>>1, 1; 0 < note; note, val = note>>1, val+1 {
				if 0 != note&1 {
					gl = append(gl, packed|uint32(val)<<8)
				}
			}
		}
	}
	heap.Init(&gl) // sort the list to a heap
	for 0 < gl.Len() {
		cp := *s                                                    // copy ALL of s by value
		cp.DirtyRow, cp.DirtyCol, cp.DirtyBox = 0x1FF, 0x1FF, 0x1FF // All Dirty
		packed = heap.Pop(&gl).(uint32)
		ent, cell = uint8(packed>>8)&0xF, uint8(packed&0xFF)
		if 0 != s.Num[cell] {
			fmt.Printf("Guess %d: [%d] = %d but was already set to %d\n", s.Rem, cell, ent, s.Num[cell])
			panic("logic error")
		}
		// validate guess correctly decoded // XOR so both the 'zero' and guess bits must be set -- not the problem
		// if 0 != (s.Note[cell]&(1|(1<<ent)))^(1|(1<<ent)) {
		// fmt.Printf("Error decoding guess from %#x to [%d] = %d vs %#b\n", packed, cell, ent, s.Note[cell])
		// panic("logic error")
		// }
		// if 0 == ent {
		//	panic("Somehow, a cell value of 0 made it in.")
		// }
		// fmt.Printf("Trying guess %d: [%d] = %d\n", s.Rem, cell, ent)
		cp.Num[cell], cp.Note[cell], cp.Rem = ent, 1<<ent, cp.Rem-1
		// cp.ResetAllNotes() // This shouldn't be required - Each scan and mask type should have an equivalent effect and must happen anyway
		box = cp.SolveInner()
		// fmt.Printf("RETURN from %d: [%d] = %d\t= %d\n", s.Rem, cell, ent, box)
		if 0 == box {
			*s = cp // Solved, overwrite object and return
			// fmt.Println("Return Answer")
			return 0
		} else if 0xFF == box { // Successfully proved that number CANNOT go in that cell
			s.Note[cell] &^= 1 << ent // At the guess stage, just make updates then full-process the board if required later
			// s.DirtyRow, s.DirtyCol, s.DirtyBox = 0x1FF, 0x1FF, 0x1FF                    // All Dirty -- Not required, since if there isn't a valid guess panic()
			// fmt.Printf("Learned (%d) [%d] != %d\n", s.Rem, cell, ent)
		} else if bestGscore > box {
			// This isn't reached in the current flow of the solver, since a guess exhausts that search space
			// This IS here in case a future version terminates early
			bestGnum, bestGcell, bestGscore = ent, cell, box
			fmt.Printf("Best Guess [%d] = %d (%d)\n", cell, ent, box)
		}
	}

	// This IS here in case a future version terminates early to give guesses weight
	if 0xFF > bestGscore {
		fmt.Printf("Applying Best Guess [%d] = %d (score %d)\n", bestGcell, 1<<bestGnum, bestGscore)
		s.Num[bestGcell], s.Note[bestGcell], s.Rem = bestGnum, 1<<bestGnum, s.Rem-1 // If it wasn't solved, apply the best guess and try again
		s.DirtyRow, s.DirtyCol, s.DirtyBox = 0x1FF, 0x1FF, 0x1FF                    // All Dirty
	} else {
		// iter = 0xFFFF
		return 0xFF
	}
	return s.Rem
}

func (s *SuDoKu) SolveInner() uint8 {
	iter := 0
	// fmt.Printf("Entry\t%d remain\n", s.Rem)
	// Scan, Mask, Solve
	for 0 < s.Rem {
		iter++
		// if 1 == s.Rem {
		// s.Print()
		// }
		if 0 != s.ReduceRow() || 0 != s.ReduceCol() || 0 != s.ReduceBox() || 0 != s.Validate(false) {
			return 0xFF
		}
		if 0 == iter&0xFFFF {
			fmt.Printf("Iter %5d\t%d remain\n", iter, s.Rem)
			s.Print()
			panic("could not solve")
			// return
		}
		// Fully reduced? ... have to make a guess
		if 0 == s.DirtyRow && 0 == s.DirtyCol && 0 == s.DirtyBox && 0xFF == s.GuessRequired() {
			// s.Print()
			return 0xFF // This could be a guess within a guess
			// panic("unsolved and no guess worked")
		}
	}
	// fmt.Printf("Ret %5d\t%d remain\n", iter, s.Rem)
	return s.Rem
}

func Euler0096(fn string) uint16 {
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler0096 unable to open: " + fn)
	}
	defer fh.Close()
	var pos, ret uint16
	scanner := bufio.NewScanner(fh)
	// split lines is default
	sdk := SuDoKu{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if 'G' == line[0] {
			pos = 0 // New puzzle
			// fmt.Printf("\n\n%s\n\n", string(line))
			continue
		}
		for ii, iiLm := 0, len(line); ii < iiLm && pos < 81; ii++ {
			if '0' <= line[ii] && line[ii] <= '9' {
				sdk.Num[pos] = line[ii] - '0'
				pos++
			}
		}
		if 81 == pos {
			if 0 == sdk.Num[0] || 0 == sdk.Num[1] || 0 == sdk.Num[2] {
				if 0 != sdk.Solve() {
					fmt.Println("unsolved")
					sdk.Print()
					// panic("unsolved")
				}
			}
			if 0 == sdk.Num[0] || 0 == sdk.Num[1] || 0 == sdk.Num[2] {
				panic("Should be solved!")
			}
			ret += 100*uint16(sdk.Num[0]) + 10*uint16(sdk.Num[1]) + uint16(sdk.Num[2])
		}
	}
	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 96 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 96: passed puzzle 1 test
Euler 96: Su Doku: 24702

real    0m0.138s
user    0m0.196s
sys     0m0.067s
.
*/
func main() {
	bitLUT = [32]uint8{}
	for ii := 0; ii < 32; ii++ {
		for t := ii; 0 < t; t >>= 1 {
			bitLUT[ii] += uint8(t & 1)
		}
	}

	var sum uint16
	//test
	tgrid := [81]uint8{
		0, 0, 3, 0, 2, 0, 6, 0, 0,
		9, 0, 0, 3, 0, 5, 0, 0, 1,
		0, 0, 1, 8, 0, 6, 4, 0, 0,
		0, 0, 8, 1, 0, 2, 9, 0, 0,
		7, 0, 0, 0, 0, 0, 0, 0, 8,
		0, 0, 6, 7, 0, 8, 2, 0, 0,
		0, 0, 2, 6, 0, 9, 5, 0, 0,
		8, 0, 0, 2, 0, 3, 0, 0, 9,
		0, 0, 5, 0, 1, 0, 3, 0, 0}
	agrid := [81]uint8{
		4, 8, 3, 9, 2, 1, 6, 5, 7,
		9, 6, 7, 3, 4, 5, 8, 2, 1,
		2, 5, 1, 8, 7, 6, 4, 9, 3,
		5, 4, 8, 1, 3, 2, 9, 7, 6,
		7, 2, 9, 5, 6, 4, 1, 3, 8,
		1, 3, 6, 7, 9, 8, 2, 4, 5,
		3, 7, 2, 6, 8, 9, 5, 1, 4,
		8, 1, 4, 2, 5, 3, 7, 6, 9,
		6, 9, 5, 4, 1, 7, 3, 8, 2}
	//	 0  1  2   3  4  5   6  7  8
	//	 9 10 11  12 13 14  15 16 17
	//	18 19 20  21 22 23  24 25 26

	//	27 28 29  30 31 32  33 34 35
	//	36 37 38  39 40 41  42 43 44
	//	45 46 47  48 49 50  51 52 53

	//	54 55 56  57 58 59  60 61 62
	//	63 64 65  66 67 68  69 70 71
	//	72 73 74  75 76 77  78 79 80
	// Grid 29
	ans29 := [81]uint8{
		2, 3, 5, 7, 6, 1, 4, 8, 9,
		4, 1, 9, 3, 2, 8, 5, 7, 6,
		8, 6, 7, 5, 4, 9, 2, 1, 3,
		7, 4, 6, 1, 3, 5, 9, 2, 8,
		5, 2, 1, 8, 9, 6, 7, 3, 4,
		9, 8, 3, 4, 7, 2, 6, 5, 1,
		3, 9, 4, 2, 8, 7, 1, 6, 5,
		6, 5, 2, 9, 1, 3, 8, 4, 7,
		1, 7, 8, 6, 5, 4, 3, 9, 2}
	grid29 := [81]uint8{
		0, 3, 0, 0, 0, 0, 0, 8, 0,
		0, 0, 9, 0, 0, 0, 5, 0, 0,
		0, 0, 7, 5, 0, 9, 2, 0, 0,
		7, 0, 0, 1, 0, 5, 0, 0, 8,
		0, 2, 0, 0, 9, 0, 0, 3, 0,
		9, 0, 0, 4, 0, 2, 0, 0, 1,
		0, 0, 4, 2, 0, 7, 1, 0, 0,
		0, 0, 2, 0, 0, 0, 8, 0, 0,
		0, 7, 0, 0, 0, 0, 0, 9, 0}
	test := SuDoKu{Num: tgrid}
	test.Solve()
	tgrid = test.Num
	ok := true
	for ii := 0; ii < 81; ii++ {
		if tgrid[ii] != agrid[ii] {
			ok = false
			break
		}
	}
	if !ok {
		fmt.Printf("\n%v\n%v\n", tgrid, agrid)
		panic("Euler 96: Test Case failed")
	}
	fmt.Println("Euler 96: passed puzzle 1 test")

	test = SuDoKu{Num: grid29}
	test.ResetAllNotes()

	check29 := func() {
		for ii := 0; ii < 81; ii++ {
			if 0 != test.Num[ii] && ans29[ii] != test.Num[ii] {
				fmt.Printf("Value check 29 failed: [%2d] %d != %d\n", ii, test.Num[ii], ans29[ii])
				panic("Checksum Failed")
			}
		}
	}

	for 0 != test.DirtyRow || 0 != test.DirtyCol || 0 != test.DirtyBox {
		if 0 != test.ReduceRow() || 0 != test.Validate(false) {
			test.Print()
			panic("Failed Row Validation")
		}
		check29()
		if 0 != test.ReduceCol() || 0 != test.Validate(false) {
			test.Print()
			panic("Failed Col Validation")
		}
		check29()
		if 0 != test.ReduceBox() {
			test.Print()
			panic("Failed Box")
		}
		check29()
		if 0 != test.Validate(false) {
			test.Print()
			panic("Failed Box Validation")
		}
	}
	if 0 != test.ReduceRow() || 0 != test.ReduceCol() || 0 != test.ReduceBox() || 0 != test.Validate(true) {
		test.Print()
		panic("Failed Validation")
	}
	// test.Print()
	// return

	//run
	sum = Euler0096("p096_sudoku.txt")
	fmt.Printf("Euler 96: Su Doku: %d\n", sum)
	if 24702 != sum {
		panic("Did not reach expected value.")
	}
}
