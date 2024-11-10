// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=79
https://projecteuler.net/minimal=79

<p>A common security method used for online banking is to ask the user for three random characters from a passcode. For example, if the passcode was 531278, they may ask for the 2nd, 3rd, and 5th characters; the expected reply would be: 317.</p>
<p>The text file, <a href="resources/documents/0079_keylog.txt">keylog.txt</a>, contains fifty successful login attempts.</p>
<p>Given that the three characters are always asked for in order, analyse the file so as to determine the shortest possible secret passcode of unknown length.</p>


/
*/
/*
	Post-answer preface: Now I understand this unclear phrase 'determine the shortest possible secret passcode of unknown length'.  It wasn't asking for a length, but to find a possible passcode: criteria, the one with the shortest possible length.


	Worst case would be just 0-9 repeated 3 times, or the same backwards.  Without knowing more about the characters requested it's possible for even some reference number like Pi or a crazy large prime to be the answer.
	I considered a variation of the problem while away from my PC, but thought it was 500 list entries rather than 50.  It might have been worth iterating and filtering out candidates that failed a sort of hash of which digits had been seen in each slot (E.G. bits 0-9 for the last digit seen, 10-19 for the middle, and 20-29 for the first digit, with two spare bits).

	What data can be gleaned from the list:
		A set of nodes
			value
			order

	A version of a solution might look like a cursor for each row, that advances to the next slot when a number matches; that seems slow.  It fails to merge later dupes.

	Each digit might have some lists... maybe a 'must also be seen before' set of rules, and a 'must also be seen after'.
	For single numbers / characters a single rule that's part of a longer chain can be discarded as redundant.
	A map could be used to support longer string segments.  The file read time's probably the main time sink.

	A glance at the list of problems under 150, I don't see any obvious from the title reuse so it'd be faster to convert everything to integers

	If thought of as the question: where is the lowest insert cost? (most numbers matched) that doesn't help with the next question: where to insert the not matched portion

	I added a hash-map to make it easy to de-duplicate the entries.  Seeking more inspiration from the data, I realized the file was about 50% dupes.

	Maybe I can recover a passphrase rather than just compute the length; that seems about the same difficulty.
	I'll try left to right first since that's how fmt.Printf can dump stuff.

	Given the passcode seems to use each number once, it would probably be faster to use a uint16 bitvector of prior numbers, and a uint16 (for alignment) count of list size (to speed selection)... however in trying to gain a better understanding of the question asked the answer was already found.
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

func Euler0079(fn string) (uint, string) {
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler 79 unable to open: " + fn)
	}
	defer fh.Close()

	keylog := make(map[string]string)

	var pos, output uint
	scanner := bufio.NewScanner(fh)
	// split lines is default
	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := keylog[line[:3]]; !exists {
			c := (line[:3] + "\000")[:3] // concise (code) 'copy' from the file buffer ; strings.Clone() is 1.20+
			keylog[c] = c
		}
		pos++
	}

	fmt.Printf("Euler 79 input...\n")

	hist := make(map[byte]map[byte]byte) // history / histogram of what came before
	for _, val := range keylog {
		fmt.Printf("%s\n", val)
		if _, exists := hist[val[0]]; !exists {
			hist[val[0]] = make(map[byte]byte)
		}
		// [0] No history, but make the map to mark that it exists
		if _, exists := hist[val[1]]; !exists {
			hist[val[1]] = make(map[byte]byte)
		}
		// Before the middle
		if _, exists := hist[val[1]][val[0]]; !exists {
			hist[val[1]][val[0]] = val[0]
		}
		if _, exists := hist[val[2]]; !exists {
			hist[val[2]] = make(map[byte]byte)
		}
		// Before the end
		if _, exists := hist[val[2]][val[1]]; !exists {
			hist[val[2]][val[1]] = val[1]
		}
	}

	// If nothing came before anything, it can be printed in any order, otherwise print the smallest discrepancy
	var shortest, length byte
	output = 0
	ret := ""
	fmt.Printf("Euler 79 guess:\t")
	for 0 < len(hist) {
		shortest, length = 0, 100 // C-strings use \0 terminators and no text file or number pad has that as an input so this is the traditional in band choice for unassigned value
		for k, v := range hist {
			if int(length) > len(v) {
				shortest, length = k, byte(len(v))
			}
		}
		// about to print 'shortest' purge it from the before / history values
		for k, _ := range hist {
			delete(hist[k], shortest)
		}
		delete(hist, shortest)
		output++
		ret += string(shortest)
		fmt.Printf("%c", shortest)
	}
	fmt.Printf("\n")

	return output, ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 79 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 79 guess: 73162890
Euler 79: Passcode Derivation: shortest possible 8:     73162890

real    0m0.131s
user    0m0.182s
sys     0m0.056s
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"

	//run
	r, str := Euler0079("0079_keylog.txt")
	fmt.Printf("Euler 79: Passcode Derivation: shortest possible %d:\t%s\n", r, str)
	if 8 != r {
		panic("Did not reach expected value.")
	}
}
