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
	FIXME: Refactor as above, include a 'no solution' return path too.
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
)

var bitLUT [32]uint8 // init() to lookup of bit popcount for 5 bits

func SuDoKuSolver(num *[81]uint8) {
	//Golang... I want to hard-ref this array NOT copy it inside the funcion! num := *ptrg //[81]uint8)
	// var num [81]uint8
	// *num := ptrg //[81]uint8)
	note := [81]uint16{}
	nummask := uint16(0b111_111_111_1) // 0 means unsolved
	var drow0, drow1, dcol0, dcol1, dbox0, dbox1 uint16
	_, _, _, _, _, _ = drow0, drow1, dcol0, dcol1, dbox0, dbox1
	rem := uint8(81)

	// Process initial grid state
	for ii := 0; ii < 81; ii++ {
		if 0 == num[ii] {
			note[ii] = nummask
		} else {
			note[ii], rem = 1<<num[ii], rem-1
		}
	}
	drow1, dcol1, dbox1 = 0x1FF, 0x1FF, 0x1FF // All dirty

	iter := 0
	fmt.Printf("Iter %5d\t%d remain\n", iter, rem)
	// Scan, Mask, Solve
	var bigrow, bigcol, brb, bcb, row, rb, col, cell uint8
	var only [9]uint8
	for 0 < rem {
		iter++
		if 0 == iter&0xFFFF || (0 == drow1 && 0 == dcol1 && 0 == dbox1) {
			fmt.Printf("Iter %5d\t%d remain\n", iter, rem)
			for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
				for col = 0; col < 9; col++ {
					cell = rb + col
					if 1 == bitLUT[note[cell]>>5]+bitLUT[note[cell]&0x1F] {
						fmt.Print("            .")
					} else {
						fmt.Printf("   %10b", note[cell])
					}
				}
				fmt.Println()
			}
			for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
				for col = 0; col < 9; col++ {
					cell = rb + col
					if 0 == (*num)[cell] {
						fmt.Print(" .")
					} else {
						fmt.Printf(" %d", (*num)[cell])
					}
				}
				fmt.Println()
			}
			panic("could not solve")
			// return
		}

		// rows
		drow0, drow1 = drow1, 0
		only = [9]uint8{}
		for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
			// if 0 == drow0&(1<<row) {
			// continue
			// }

			// Scan known numbers
			seen := uint16(0)
			for col = 0; col < 9; col++ {
				cell = rb + col
				if 0 == note[cell]&1 {
					seen |= note[cell]
					only[(*num)[cell]-1] = 0xFF
				}
			}
			seen = ^seen // Any number that was NOT seen is allowed...

			// MASK known numbers out of unknown numbers
			for col = 0; col < 9; col++ {
				cell = rb + col
				if 1 == note[cell]&1 {
					possible := note[cell] & seen
					if possible != note[cell] {
						// Test if ONLY one possible number remains
						pscan, testnum := possible>>1, uint8(1) // shift out the zero flag
						if 0 == pscan {
							fmt.Printf("Error: Cell %d has no solution: %b\n", cell, possible)
							return
						}
						justone := uint8(0)
						// scan for if this is the only N, or contested
						for ; 0 < pscan; pscan, testnum = pscan>>1, testnum+1 {
							if 0 != pscan&1 {
								if 0 == justone {
									justone = testnum
								} else {
									justone = 0xFF
								}
								if 0 == only[testnum-1] {
									only[testnum-1] = cell
								} else {
									only[testnum-1] = 0xFF // mark contested
								}
							}
						}
						if 0 < justone && justone < 0xFF {
							fmt.Printf("found ro [%d] = %d\n", cell, justone)
							only[justone-1] = 0
							possible &^= 1    // clear zero flag
							drow1 |= 1 << row // Re-dirty this row
							(*num)[cell], only[justone-1], rem = justone, 0xFF, rem-1
						}
						// fmt.Printf("update (row) [%2d] = %10b\t%d\n", cell, possible, iter)
						note[cell] = possible
						dcol1 |= 1 << col // Dirty flags
						dbox1 |= 1 << (((row / 3) * 3) + (col / 3))
					}
				}
			}
		}
		// Number discovered?
		for rb = 0; rb < 9; rb++ {
			if 0 < only[rb] && only[rb] < 0xFF {
				cell, bcb = only[rb], rb+1
				fmt.Printf("found rd [%d] = %d\n", cell, bcb)
				(*num)[cell], note[cell], rem = bcb, 1<<bcb, rem-1
				drow1 |= 1 << (cell / 9) // recover row from cell co-ord
			}
		}

		// cols
		dcol0, dcol1 = dcol1, 0
		only = [9]uint8{}
		for col = 0; col < 9; col++ {
			// if 0 == dcol0&(1<<row) {
			// continue
			// }

			// Scan known numbers
			seen := uint16(0)
			for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
				cell = rb + col
				if 0 == note[cell]&1 {
					seen |= note[cell]
					only[(*num)[cell]-1] = 0xFF
				}
			}
			seen = ^seen // Any number that was NOT seen is allowed...

			// MASK known numbers out of unknown numbers
			for row, rb = 0, 0; row < 9; row, rb = row+1, rb+9 {
				cell = rb + col
				if 1 == note[cell]&1 {
					possible := note[cell] & seen
					if possible != note[cell] {
						// Test if ONLY one possible number remains
						pscan, testnum := possible>>1, uint8(1) // shift out the zero flag
						if 0 == pscan {
							fmt.Printf("Error: Cell %d has no solution: %b\n", cell, possible)
							return
						}
						justone := uint8(0)
						// scan for if this is the only N, or contested
						for ; 0 < pscan; pscan, testnum = pscan>>1, testnum+1 {
							if 0 != pscan&1 {
								if 0 == justone {
									justone = testnum
								} else {
									justone = 0xFF
								}
								if 0 == only[testnum-1] {
									only[testnum-1] = cell
								} else {
									only[testnum-1] = 0xFF // mark contested
								}
							}
						}
						if 0 < justone && justone < 0xFF {
							fmt.Printf("found co [%d] = %d\n", cell, justone)
							only[justone-1] = 0
							possible &^= 1    // clear zero flag
							dcol1 |= 1 << col // Re-dirty this col
							(*num)[cell], only[justone-1], rem = justone, 0xFF, rem-1
						}
						// fmt.Printf("update (row) [%2d] = %10b\t%d\n", cell, possible, iter)
						note[cell] = possible
						drow1 |= 1 << row // Dirty flags
						dbox1 |= 1 << (((row / 3) * 3) + (col / 3))
					}
				}
			}
		}
		// Number discovered?
		for rb = 0; rb < 9; rb++ {
			if 0 < only[rb] && only[rb] < 0xFF {
				cell, bcb = only[rb], rb+1
				fmt.Printf("found cd [%d] = %d\n", cell, bcb)
				(*num)[cell], note[cell], rem = bcb, 1<<bcb, rem-1
				dcol1 |= 1 << (only[rb] % 9) // recover col from cell co-ord
			}
		}

		// 'boxes' additive wheels
		dbox0, dbox1 = dbox1, 0
		only = [9]uint8{}
		for bigrow, brb = 0, 0; bigrow < 3; bigrow, brb = bigrow+1, brb+27 {
			for bigcol, bcb = 0, 0; bigcol < 3; bigcol, bcb = bigcol+1, bcb+3 {
				// if 0 == dbox0&(1<<(row*3+col)) {
				// continue
				// }

				// Scan known numbers
				seen := uint16(0)
				for row, rb = 0, 0; row < 3; row, rb = row+1, rb+9 {
					for col = 0; col < 3; col++ {
						cell = brb + bcb + rb + col
						if 0 == note[cell]&1 {
							seen |= note[cell]
							only[(*num)[cell]-1] = 0xFF
						}
					}
				}
				seen = ^seen // Any number that was NOT seen is allowed...

				// MASK known numbers out of unknown numbers
				for row, rb = 0, 0; row < 3; row, rb = row+1, rb+9 {
					for col = 0; col < 3; col++ {
						cell = brb + bcb + rb + col
						if 1 == note[cell]&1 {
							possible := note[cell] & seen
							if possible != note[cell] {
								// Test if ONLY one possible number remains
								pscan, testnum := possible>>1, uint8(1) // shift out the zero flag
								if 0 == pscan {
									fmt.Printf("Error: Cell %d has no solution: %b\n", cell, possible)
									return
								}
								justone := uint8(0)
								// scan for if this is the only N, or contested
								for ; 0 < pscan; pscan, testnum = pscan>>1, testnum+1 {
									if 0 != pscan&1 {
										if 0 == justone {
											justone = testnum
										} else {
											justone = 0xFF
										}
										if 0 == only[testnum-1] {
											only[testnum-1] = cell
										} else {
											only[testnum-1] = 0xFF // mark contested
										}
									}
								}
								if 0 < justone && justone < 0xFF {
									fmt.Printf("found bo [%d] = %d\n", cell, justone)
									only[justone-1] = 0
									possible &^= 1 // clear zero flag
									dbox1 |= 1 << (bigrow*3 + bigcol)
									(*num)[cell], only[justone-1], rem = justone, 0xFF, rem-1

								}
								// fmt.Printf("update (row) [%2d] = %10b\t%d\n", cell, possible, iter)
								note[cell] = possible
								drow1 |= 1 << (bigrow + row)
								dcol1 |= 1 << (bigcol + col) // Dirty flags
							}
						}
					}
				}
			}
		}
		// Number discovered?
		for rb = 0; rb < 9; rb++ {
			if 0 < only[rb] && only[rb] < 0xFF {
				cell, bcb = only[rb], rb+1
				fmt.Printf("found bd [%d] = %d\n", cell, bcb)
				(*num)[cell], note[cell], rem = bcb, 1<<bcb, rem-1
				drow1 |= 1 << ((only[rb]/27)*3 + ((only[rb] % 9) / 3)) // recover box from cell co-ord
			}
		}

	}
	fmt.Printf("Iter %5d\t%d remain\n", iter, rem)

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
	grid := [81]uint8{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if 'G' == line[0] {
			pos = 0 // New puzzle
			continue
		}
		for ii, iiLm := 0, len(line); ii < iiLm && pos < 81; ii++ {
			if '0' <= line[ii] && line[ii] <= '9' {
				grid[pos] = line[ii] - '0'
				pos++
			}
		}
		if 81 == pos {
			if 0 == grid[0] || 0 == grid[1] || 0 == grid[2] {
				SuDoKuSolver(&grid) // Solve if the required data wasn't already provided
			}
			ret += uint16(grid[0]) + 10*uint16(grid[1]) + 100*uint16(grid[2])
		}
	}
	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 96 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

[0 0 3 0 2 0 6 0 0 9 0 0 3 0 5 0 0 1 0 0 1 8 0 6 4 0 0 0 0 8 1 0 2 9 0 0 7 0 0 0 0 0 0 0 8 0 0 6 7 0 8 2 0 0 0 0 2 6 0 9 5 0 0 8 0 0 2 0 3 0 0 9 0 0 5 0 1 0 3 0 0]
[0 0 3 0 2 1 6 0 0 9 0 4 3 7 5 8 0 1 0 0 1 8 9 6 4 0 0 0 0 8 1 0 2 9 0 0 7 0 9 5 3 4 1 0 8 0 0 6 7 5 8 2 0 0 0 0 2 6 4 9 5 0 0 8 0 4 2 5 3 7 0 9 0 0 5 0 1 7 3 0 0]
[4 8 3 9 2 1 6 5 7 9 6 7 3 4 5 8 2 1 2 5 1 8 7 6 4 9 3 5 4 8 1 3 2 9 7 6 7 2 9 5 6 4 1 3 8 1 3 6 7 9 8 2 4 5 3 7 2 6 8 9 5 1 4 8 1 4 2 5 3 7 6 9 6 9 5 4 1 7 3 8 2]

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
	SuDoKuSolver(&tgrid)
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

	//run
	sum = Euler0096("p096_sudoku.txt")
	fmt.Printf("Euler 96: Su Doku: %d\n", sum)
	if 10935 != sum {
		panic("Did not reach expected value.")
	}
}
