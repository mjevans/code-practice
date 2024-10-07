// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=42
https://projecteuler.net/minimal=42

<p>The $n$<sup>th</sup> term of the sequence of triangle numbers is given by, $t_n = \frac12n(n+1)$; so the first ten triangle numbers are:
$$1, 3, 6, 10, 15, 21, 28, 36, 45, 55, \dots$$</p>
<p>By converting each letter in a word to a number corresponding to its alphabetical position and adding these values we form a word value. For example, the word value for SKY is $19 + 11 + 25 = 55 = t_{10}$. If the word value is a triangle number then we shall call the word a triangle word.</p>
<p>Using <a href="resources/documents/0042_words.txt">words.txt</a> (right click and 'Save Link/Target As...'), a 16K text file containing nearly two-thousand common English words, how many are triangle words?</p>



*/
/*
sed -e 's/","/\n/g' < 0042_words.txt | wc -l
1785 words, in the same cursed format as Euler 22 ( No Newlines, only  '"WORD","WORD"' like a one row, columns only, quote everything, CSV file )

*/

import (
	"bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	"strings"
	// "strconv"
	"os" // os.Stdout
)

func Euler042(fn string) int {
	trin := make(map[uint64]uint64)
	found := 0
	seen := 0
	n := uint64(0)
	for ii := 1; ii < 128; ii++ {
		// either ii or ii+1 will be even, and thus at least one factor of 2 is present
		n = uint64(ii*(ii+1)) >> 1
		trin[n] = n
	}
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler042 unable to open: " + fn)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	scanner.Split(euler.ScannerSplitNLDQ)
	for scanner.Scan() {
		str := strings.ToUpper(scanner.Text())
		if str != "," {
			calc := uint64(euler.AlphaSum(str))
			if n < calc {
				fmt.Printf("ERROR: encountered score %d for %s\n", calc, str)
			}
			if _, ex := trin[calc]; true == ex {
				found++
			}
			seen++
		} else {
			fmt.Printf("WARNING: encountered unknown\n\n%s\n\n\n", str)
		}
	}
	if 1786 != seen && "0042_words.txt" == fn {
		fmt.Printf("WARNING: checksum did not match expected value of %d evaluated %d words\n", 1785, seen)
	}
	return found
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 42 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 42: Coded Triangle Numbers: 162
*/
func main() {
	// fmt.Println(grid)
	//test
	// fmt.Println(int(euler.AlphaSum("Abc")) == 6)
	// fmt.Println(int(euler.AlphaSum("Colin"))*938 == 49714)

	//run
	fmt.Printf("Euler 42: Coded Triangle Numbers: %d\n", Euler042("0042_words.txt"))
}
