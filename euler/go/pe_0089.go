// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=89
https://projecteuler.net/minimal=89

<p>For a number written in Roman numerals to be considered valid there are basic rules which must be followed. Even though the rules allow some numbers to be expressed in more than one way there is always a "best" way of writing a particular number.</p>
<p>For example, it would appear that there are at least six ways of writing the number sixteen:</p>
<p class="margin_left monospace">IIIIIIIIIIIIIIII<br>
VIIIIIIIIIII<br>
VVIIIIII<br>
XIIIIII<br>
VVVI<br>
XVI</p>
<p>However, according to the rules only <span class="monospace">XIIIIII</span> and <span class="monospace">XVI</span> are valid, and the last example is considered to be the most efficient, as it uses the least number of numerals.</p>
<p>The 11K text file, <a href="resources/documents/0089_roman.txt">roman.txt</a> (right click and 'Save Link/Target As...'), contains one thousand numbers written in valid, but not necessarily minimal, Roman numerals; see <a href="about=roman_numerals">About... Roman Numerals</a> for the definitive rules for this problem.</p>
<p>Find the number of characters saved by writing each of these in their minimal form.</p>
<p class="smaller">Note: You can assume that all the Roman numerals in the file contain no more than four consecutive identical units.</p>



https://projecteuler.net/about=roman_numerals

Traditional Roman numerals are made up of the following denominations:

I = 1
V = 5
X = 10
L = 50
C = 100
D = 500
M = 1000

In order for a number written in Roman numerals to be considered valid there are three basic rules which must be followed.

*	Numerals must be arranged in descending order of size.
*	M, C, and X cannot be equalled or exceeded by smaller denominations.
*	D, L, and V can each only appear once.

For example, the number sixteen could be written as XVI or XIIIIII, with the first being the preferred form as it uses the least number of numerals. We could not write IIIIIIIIIIIIIIII because we are making X (ten) from smaller denominations, nor could we write VVVI because the second and third rule are being broken.

The "descending size" rule was introduced to allow the use of subtractive combinations. For example, four can be written IV because it is one before five. As the rule requires that the numerals be arranged in order of size it should be clear to a reader that the presence of a smaller numeral out of place, so to speak, was unambiguously to be subtracted from the following numeral rather than added.

For example, nineteen could be written XIX = X (ten) + IX (nine). Note also how the rule requires X (ten) be placed before IX (nine), and IXX would not be an acceptable configuration (descending size rule). Similarly, XVIV would be invalid because V can only appear once in a number.

Generally the Romans tried to use as few numerals as possible when displaying numbers. For this reason, XIX would be the preferred form of nineteen over other valid combinations, like XIIIIIIIII or XVIIII.

By mediaeval times it had become standard practice to avoid more than three consecutive identical numerals by taking advantage of the more compact subtractive combinations. That is, IV would be written instead of IIII, IX would be used instead of IIIIIIIII or VIIII, and so on.

In addition to the three rules given above, if subtractive combinations are used then the following four rules must be followed.

*	Only one I, X, and C can be used as the leading numeral in part of a subtractive pair.
*	I can only be placed before V and X.
*	X can only be placed before L and C.
*	C can only be placed before D and M.

Which means that IL would be considered to be an invalid way of writing forty-nine, and whereas XXXXIIIIIIIII, XXXXVIIII, XXXXIX, XLIIIIIIIII, XLVIIII, and XLIX are all quite legitimate, the latter is the preferred (minimal) form. However, minimal form was not a rule and there still remain in Rome many examples where economy of numerals has not been employed. For example, in the famous Colosseum the numerals above the forty-ninth entrance is written XXXXVIIII rather than XLIX.

It is also expected, but not required, that higher denominations should be used whenever possible; for example, V should be used in place of IIIII, L should be used in place of XXXXX, and D should be used in place of CCCCC. However, in the church of Sant'Agnese fuori le Mura (St Agnes' outside the walls), found in Rome, the date, MCCCCCCVI (1606), is written on the gilded and coffered wooden ceiling; I am sure that many would argue that it should have been written MDCVI.

So if we believe the adage, "when in Rome do as the Romans do," and we see how the Romans write numerals, then it clearly gives us much more freedom than many would care to admit.
/
*/
/*
	Possible sub-components: Decode (any, even malformed) Roman Numerals to a number ; Calculate the length of a roman numeral string (Probably nearly as fast to just construct one, then copy the shorter buffer out...)

	There might be some tricks to do this a tiny bit faster in place destructively... That seems inelegant.

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

var RomanNumOuts []struct {
	min, ln uint16
	chr     string
}
var RomanNumIns []uint16

func RomanNumParseB(rn []byte) uint16 {
	var pv1, pc1, pv2, pc2, ret, val uint16
	lim := len(rn)
	for ii := 0; ii < lim; ii++ {
		if 'a' <= rn[ii] && 'x' >= rn[ii] {
			val = RomanNumIns[rn[ii] - 'a']
		} else if 'A' <= rn[ii] && 'X' >= rn[ii] {
			val = RomanNumIns[rn[ii] - 'A']
		} else {
			val = 4
		}
		if 0 < val && 4 != val {
			if val == pv1 {
				pc1++
				continue
			}
			if val < pv1 {
				if 0 == pc2 {
					ret += (pv1 * pc1)
				} else if pv2 < pv1 {
					ret += (pv1 * pc1) - (pv2 * pc2)
				}
				pc2, pc1, pv2, pv1 = 0, 1, pv1, val
			} else {
				// Bigger number to the right?  Push stack and setup for subtraction
				pc2, pc1, pv2, pv1 = pc1, 1, pv1, val
			}
		} else {
			fmt.Printf("ERROR: encountered non-roman numeral [%d] => %d (%c) in string: %s\n", ii, rn[ii], rn[ii], string(rn))
		}
	}
	if 0 == pc2 {
		ret += (pv1 * pc1)
	} else if pv2 < pv1 {
		ret += (pv1 * pc1) - (pv2 * pc2)
	}
	return ret
}

func RomanNumLen(num uint16) uint16 {
	var ret uint16
	lim := len(RomanNumOuts)
	for ii := 0; ii < lim; ii++ {
		for num >= RomanNumOuts[ii].min {
			ret, num = ret+RomanNumOuts[ii].ln, num-RomanNumOuts[ii].min
		}
	}
	return ret
}

func Euler0089(fn string) uint32 {
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler0089 unable to open: " + fn)
	}
	defer fh.Close()
	var pos, saved uint32
	scanner := bufio.NewScanner(fh)
	// split lines is default
	for scanner.Scan() {
		pos++
		line := scanner.Bytes()
		saved += uint32(len(line)) - uint32(RomanNumLen(RomanNumParseB(line)))
	}
	return saved
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 89 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 89: Roman Numerals: 743

real    0m0.099s
user    0m0.157s
sys     0m0.044s
.
*/
func main() {
	RomanNumOuts = []struct {
		min, ln uint16
		chr     string
	}{
		{1000, 1, "M"},
		{900, 2, "CM"},
		{500, 1, "D"},
		{400, 2, "CD"},
		{100, 1, "C"},
		{90, 2, "XC"},
		{50, 1, "L"},
		{40, 2, "XL"},
		{10, 1, "X"},
		{9, 2, "IX"},
		{5, 1, "V"},
		{4, 2, "IV"},
		{1, 1, "I"},
	}
	// x is the highest char value so...
	RomanNumIns = make([]uint16, 'x'-'a'+1)
	RomanNumIns['i'-'a'] = 1
	RomanNumIns['v'-'a'] = 5
	RomanNumIns['x'-'a'] = 10
	RomanNumIns['l'-'a'] = 50
	RomanNumIns['c'-'a'] = 100
	RomanNumIns['d'-'a'] = 500
	RomanNumIns['m'-'a'] = 1000

	var r uint32
	//test
	RomanNumTestCases := []struct {
		s string
		v uint16
	}{
		{"MMMMDCLXXII", 4672},
		{"MMDCCCLXXXIII", 2883},
		{"MMMDLXVIIII", 3569},
		{"MMMMDXCV", 4595},
		{"DCCCLXXII", 872},
		{"MMCCCVI", 2306},
		{"MMMCDLXXXVII", 3487},
		{"MMMMCCXXI", 4221},
		{"MMMCCXX", 3220},
		{"MMMMDCCCLXXIII", 4873},
	}
	for _, test := range RomanNumTestCases {
		r := RomanNumParseB([]byte(test.s))
		if test.v != r {
			fmt.Printf("Euler 89: failed test case %s = %d got %d\n", test.s, test.v, r)
			panic("")
		}
	}

	//run
	r = Euler0089("0089_roman.txt")
	fmt.Printf("Euler 89: Roman Numerals: %d\n", r)
	if 743 != r {
		panic("Did not reach expected value.")
	}
}
