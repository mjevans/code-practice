// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=67
https://projecteuler.net/minimal=67

<p>By starting at the top of the triangle below and moving to adjacent numbers on the row below, the maximum total from top to bottom is 23.</p>
<p class="monospace center"><span class="red"><b>3</b></span><br><span class="red"><b>7</b></span> 4<br>
2 <span class="red"><b>4</b></span> 6<br>
8 5 <span class="red"><b>9</b></span> 3</p>
<p>That is, 3 + 7 + 4 + 9 = 23.</p>
<p>Find the maximum total from top to bottom in <a href="resources/documents/0067_triangle.txt">triangle.txt</a> (right click and 'Save Link/Target As...'), a 15K text file containing a triangle with one-hundred rows.</p>
<p class="smaller"><b>NOTE:</b> This is a much more difficult version of <a href="problem=18">Problem 18</a>. It is not possible to try every route to solve this problem, as there are 2<sup>99</sup> altogether! If you could check one trillion (10<sup>12</sup>) routes every second it would take over twenty billion years to check them all. There is an efficient algorithm to solve it. ;o)</p>



*/
/*

func MaximumPathSum(tri [][]int) int {}

The same, except line by line (to also offer the maximum value and address in the combined step)
func MaximumPathSumAppendShrink(dst, c []int) ([]int, int, int) {}


*/

import (
	"bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	"os" // os.Stdout
	"strconv"
	"strings"
)

func Euler0067(fn string) int {
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler0067 unable to open: " + fn)
	}
	defer fh.Close()
	var pos, maxnum, ii int
	var tri [][]int
	scanner := bufio.NewScanner(fh)
	// split lines is default
	for scanner.Scan() {
		pos++
		line := scanner.Text()
		snum := strings.SplitN(line, " ", -1)
		slen := len(snum)
		pline := make([]int, 0, slen)
		if maxnum < slen {
			maxnum = slen
		}
		for ii = 0; ii < slen; ii++ {
			tmp, err := strconv.ParseInt(snum[ii], 10, 0)
			if nil != err {
				fmt.Printf("Line %d : %d:\t%v\n", pos, ii, err)
			}
			pline = append(pline, int(tmp))
		}
		tri = append(tri, pline)
	}
	maxnum++
	buf := make([]int, maxnum)

	oneshot := euler.MaximumPathSum(tri)

	ii = len(tri)
	for 0 < ii {
		ii--
		buf, maxnum, _ = euler.MaximumPathSumAppendShrink(buf, tri[ii])
	}

	fmt.Printf("Euler 67 got:\t %d\t %d\n", oneshot, maxnum)

	return oneshot
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 67 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 67 got:    7273    7273
Euler 67: Maximum Path Sum II: 7273
*/
func main() {
	//test

	//run
	fmt.Printf("Euler 67: Maximum Path Sum II: %d\n", Euler0067("0067_triangle.txt"))
}
