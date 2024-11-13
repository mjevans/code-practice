// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=81
https://projecteuler.net/minimal=81

<p>In the $5$ by $5$ matrix below, the minimal path sum from the top left to the bottom right, by <b>only moving to the right and down</b>, is indicated in bold red and is equal to $2427$.</p>
<div class="center">
$$
\begin{pmatrix}
\color{red}{131} &amp; 673 &amp; 234 &amp; 103 &amp; 18\\
\color{red}{201} &amp; \color{red}{96} &amp; \color{red}{342} &amp; 965 &amp; 150\\
630 &amp; 803 &amp; \color{red}{746} &amp; \color{red}{422} &amp; 111\\
537 &amp; 699 &amp; 497 &amp; \color{red}{121} &amp; 956\\
805 &amp; 732 &amp; 524 &amp; \color{red}{37} &amp; \color{red}{331}
\end{pmatrix}
$$
</div>
<p>Find the minimal path sum from the top left to the bottom right by only moving right and down in <a href="resources/documents/0081_matrix.txt">matrix.txt</a> (right click and "Save Link/Target As..."), a 31K text file containing an $80$ by $80$ matrix.</p>


/
*/
/*
	81, 82, and 83 are all variations of the earlier triangle, but instead of an in place algorithm this seems to lean towards a dynamic programming solution.
	81, dynamic programming was sufficient

	81 is similar enough that folding the minimum path values up to the diagonal and then using minimum path back to the start should work.
	However 81 and 82 are subsets of the problem from 83, just with less vectors to evaluate.

	81 Only right/down, only 'forward', implicitly just asks, "which path has the least traversal cost?"
	82 something about go around?
	83 Free Form, still traversal cost, but how to know when to give up on a direction?

	I realized any path, such all along an edge, can be taken as an initial 'best so far'.
	Then modifications can be attempted to shrink it towards a more optimal path.
	If the lowest cell value is known, and a method for calculating the moves to the target is known, then the remaining 'most optimistic' cost could be estimated without traversal.  Without that it can just be assumed to be zero.

	Reconsidering the triangle path algorithm, instead of storing in one row of buffer, either in place (overwright) or a whole results table (at a stretch, a cache) could be used.
	If not in place, it'd be possible to iterate over backwards spread to 'correct' / propagate lower costs into cells.

	That's TraverseEntireMatrix but I might want to update it with the 'best path so far' limit...
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

func Euler0081(fn string) int32 {
	// func LoadMatrix[SL ~[][]INT, INT ~int | ~uint | ~uint32 | ~int32 | ~uint16 | ~int16](fn, split string, limit, base INT) (SL, INT, INT) {
	// func TraverseEntireMatrix[SL ~[][]INT, INT ~int | ~int64 | ~int32 | ~int16](m SL, stR, stC, edR, edC INT, moveRC []INT) INT {
	m, y, x := euler.LoadMatrix(fn, ",", int32(10))
	var z int32
	return euler.TraverseEntireMatrix(m, z, z, y-1, x-1, []int32{1, 0, 0, 1})
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 81 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 81: Path Sum: Two Ways: 427337

real    0m0.192s
user    0m0.221s
sys     0m0.091s
.
*/
func main() {
	var r int32

	//test
	test81 := [][]int16{
		[]int16{131, 673, 234, 103, 18},  // *
		[]int16{201, 96, 342, 965, 150},  // ***
		[]int16{630, 803, 746, 422, 111}, //   **
		[]int16{537, 699, 497, 121, 956}, //    *
		[]int16{805, 732, 524, 37, 331},  //    **
	}
	r = int32(euler.TraverseEntireMatrix(test81, 0, 0, 4, 4, []int16{1, 0, 0, 1}))

	if 2427 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	r = Euler0081("0081_matrix.txt")
	fmt.Printf("Euler 81: Path Sum: Two Ways: %d\n", r)
	if 427337 != r {
		panic("Did not reach expected value.")
	}
}
