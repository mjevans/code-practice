// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=93
https://projecteuler.net/minimal=93

<p>By using each of the digits from the set, $\{1, 2, 3, 4\}$, exactly once, and making use of the four arithmetic operations ($+, -, \times, /$) and brackets/parentheses, it is possible to form different positive integer targets.</p>
<p>For example,</p>
\begin{align}
8 &amp;= (4 \times (1 + 3)) / 2\\
14 &amp;= 4 \times (3 + 1 / 2)\\
19 &amp;= 4 \times (2 + 3) - 1\\
36 &amp;= 3 \times 4 \times (2 + 1)
\end{align}
<p>Note that concatenations of the digits, like $12 + 34$, are not allowed.</p>
<p>Using the set, $\{1, 2, 3, 4\}$, it is possible to obtain thirty-one different target numbers of which $36$ is the maximum, and each of the numbers $1$ to $28$ can be obtained before encountering the first non-expressible number.</p>
<p>Find the set of four distinct digits, $a \lt b \lt c \lt d$, for which the longest set of consecutive positive integers, $1$ to $n$, can be obtained, giving your answer as a string: <i>abcd</i>.</p>


/
*/
/*
	The parentheses make this rather annoying.  Not only do the operations happen to the terms in any possible order, it's possible to evaluate sets of terms in multiple different ways.
	(((a+b)+c)+d) OR (a+b)*(c+d) OR (a+(b+(c+d)))  << The right side is a mirror of the left...
	For + or * the order of operations looks might decieve a qucik reader.
	(((a/b)-c)*d) OR (a/b)-(c*d) OR (a/(b-(c*d)))
	Now all three are distinct, but the left and right will happen across the set if the numbers and operations all permuate.  However that middle version has an entirely different type of flow.

	Example 2: 14 = 4 x (3 + 1/2)
	Rationals... great.

/
*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0092TestCombo(dig []uint8) (int32, int32, int32) {
	var mn, mx, tmp int32
	var r1, r2 euler.Rat2
	var x, y, z uint8
	if 4 != len(dig) {
		return 0, 0, 0 // Incorrect number of terms
	}
	opPrint := "+-*/"
	_ = opPrint
	mn, mx = 0x7FFFFFFF, -0x7FFFFFFF
	bv := [1 + 3024>>3]uint8{1} // this has to be initialized every run, though 0 is free // 9*8*7*6=3024
	// bv[0] |= 1

	for ii := 0; ii < 24; ii++ {
		shf := euler.PermutationSlUint8(uint64(ii), dig)
		// euler.NewRat2(int64(), 1)
		a, b, c, d := euler.NewRat2(int64(shf[0]), 1), euler.NewRat2(int64(shf[1]), 1), euler.NewRat2(int64(shf[2]), 1), euler.NewRat2(int64(shf[3]), 1) // These will be referenced a lot
		for x = 0; x < 4; x++ {
			for y = 0; y < 4; y++ {
				for z = 0; z < 4; z++ {
					// (((a~b)~c)~d)
					r1 = a
					switch x {
					case 0:
						r1 = r1.AddRat(b)
					case 1:
						r1 = r1.SubRat(b)
					case 2:
						r1 = r1.MulRat(b)
					case 3:
						// Special case, /0 will be forced to 0 since that has the least negative impact on the resuls
						if 0 == b.Num {
							r1.Num = 0
						} else {
							r1 = r1.DivRat(b)
						}
					}
					switch y {
					case 0:
						r1 = r1.AddRat(c)
					case 1:
						r1 = r1.SubRat(c)
					case 2:
						r1 = r1.MulRat(c)
					case 3:
						// Special case, /0 will be forced to 0 since that has the least negative impact on the resuls
						if 0 == c.Num {
							r1.Num = 0
						} else {
							r1 = r1.DivRat(c)
						}
					}
					switch z {
					case 0:
						r1 = r1.AddRat(d)
					case 1:
						r1 = r1.SubRat(d)
					case 2:
						r1 = r1.MulRat(d)
					case 3:
						// Special case, /0 will be forced to 0 since that has the least negative impact on the resuls
						if 0 == d.Num {
							r1.Num = 0
						} else {
							r1 = r1.DivRat(d)
						}
					}
					r1 = euler.ReduceRat2(r1)
					if 1 == r1.Den {
						tmp = int32(r1.Num)
						if 0 < tmp && tmp < 3025 {
							bv[tmp>>3] |= 1 << (tmp & 0b111)
						}
						if tmp < mn {
							mn = tmp
						}
						if tmp > mx {
							mx = tmp
						}
					}

					// (a~b)~(c~d)
					r1 = a
					switch x {
					case 0:
						r1 = r1.AddRat(b)
					case 1:
						r1 = r1.SubRat(b)
					case 2:
						r1 = r1.MulRat(b)
					case 3:
						// Special case, /0 will be forced to 0 since that has the least negative impact on the resuls
						if 0 == b.Num {
							r1.Num = 0
						} else {
							r1 = r1.DivRat(b)
						}
					}
					r2 = c
					switch z {
					case 0:
						r2 = r2.AddRat(d)
					case 1:
						r2 = r2.SubRat(d)
					case 2:
						r2 = r2.MulRat(d)
					case 3:
						// Special case, /0 will be forced to 0 since that has the least negative impact on the resuls
						if 0 == c.Num {
							r2.Num = 0
						} else {
							r2 = r2.DivRat(d)
						}
					}
					switch y {
					case 0:
						r1 = r1.AddRat(r2)
					case 1:
						r1 = r1.SubRat(r2)
					case 2:
						r1 = r1.MulRat(r2)
					case 3:
						// Special case, /0 will be forced to 0 since that has the least negative impact on the resuls
						if 0 == r2.Num {
							r1.Num = 0
						} else {
							r1 = r1.DivRat(r2)
						}
					}
					r1 = euler.ReduceRat2(r1)
					if 1 == r1.Den {
						tmp = int32(r1.Num)
						if 0 < tmp && tmp < 3025 {
							bv[tmp>>3] |= 1 << (tmp & 0b111)
						}
						if tmp < mn {
							mn = tmp
						}
						if tmp > mx {
							mx = tmp
						}
					}
				}
			}
		}
	}
	for tmp = 0; 0xFF == bv[tmp]; tmp++ {
	}
	for x = 0; 0 != bv[tmp]&(1<<x); x++ {
	}
	tmp = ((tmp) << 3) + int32(x) - 1
	return tmp, mn, mx
}

func Euler0093(min, max uint8) ([]uint8, int32, int32, int32) {
	// 9*9*9*9=6561 so even an int16 would work
	// 9*8*7*6=3024
	dig, mxDig := make([]uint8, 4), make([]uint8, 4)
	var mr, mn, mx int32
	for dig[0] = min; dig[0] <= max-3; dig[0]++ {
		for dig[1] = dig[0] + 1; dig[1] <= max-2; dig[1]++ {
			for dig[2] = dig[1] + 1; dig[2] <= max-1; dig[2]++ {
				for dig[3] = dig[2] + 1; dig[3] <= max-0; dig[3]++ {
					maxRun, minVal, maxVal := Euler0092TestCombo(dig)
					if maxRun > mr {
						mr, mn, mx = maxRun, minVal, maxVal
						copy(mxDig, dig)
					}
				}
			}
		}
	}
	return mxDig, mr, mn, mx
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 93 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 93: test: [1 2 3 4] = 28, -18, 36
Euler 93: Arithmetic Expressions: [1 2 5 8] = 51, -70, 120 = 1258

real    0m0.133s
user    0m0.170s
sys     0m0.071s
.
*/
func main() {
	var mr, mn, mx int32
	var digits []uint8
	//test
	digits, mr, mn, mx = Euler0093(1, 4)
	fmt.Printf("Euler 93: test: %v = %d, %d, %d\n", digits, mr, mn, mx)
	if 36 != mx || 28 != mr {
		panic("Did not reach expected value.")
	}

	//run
	digits, mr, mn, mx = Euler0093(1, 9)
	fmt.Printf("Euler 93: Arithmetic Expressions: %v = %d, %d, %d = %d%d%d%d\n", digits, mr, mn, mx, digits[0], digits[1], digits[2], digits[3])
	if 51 != mr || 120 != mx || -70 != mn {
		panic("Did not reach expected value.")
	}
}
