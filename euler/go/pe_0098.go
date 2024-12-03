// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=98
https://projecteuler.net/minimal=98

<p>By replacing each of the letters in the word CARE with $1$, $2$, $9$, and $6$ respectively, we form a square number: $1296 = 36^2$. What is remarkable is that, by using the same digital substitutions, the anagram, RACE, also forms a square number: $9216 = 96^2$. We shall call CARE (and RACE) a square anagram word pair and specify further that leading zeroes are not permitted, neither may a different letter have the same digital value as another letter.</p>
<p>Using <a href="resources/documents/0098_words.txt">words.txt</a> (right click and 'Save Link/Target As...'), a 16K text file containing nearly two-thousand common English words, find all the square anagram word pairs (a palindromic word is NOT considered to be an anagram of itself).</p>
<p>What is the largest square number formed by any member of such a pair?</p>
<p class="smaller">NOTE: All anagrams formed must be contained in the given text file.</p>



/
*/
/*
	Still a bit sleepy, I already really dislike this problem.  There's no inherent relation between numbers and letters, it's entirely arbitrary.
	Real rules:
	+ _A_ Non-leading letter may be assigned the value 0
	+ Values 0 through 9 may be assigned to 10 unique letters
	- An exactly reversed value is not considered an 'anagram'
	+ The patterns of allowed permutations are based on two permutations which each equal a different squared number.
	tr ',' '\n' < 0098_words.txt | tr -d '"' | while read LINE ; do printf "%s" "$LINE" | wc -c ; done | sort -nu | tail -n 1
	14
	The longest word has 14 letters, so this must handle 64 bit int sized squares up to length 14...
	Doesn't care what permutation number a permutation is, just that it has other permutations with the same signature; though they must be distinct permutations (not the same, not reversed).

	Offhand, any word longer than 15 characters can be classified as a special case.
	15 bins for squares (1..15) (so maybe 16 bins just to make indexing a tiny bit faster and the code easier)
	How many times can each unique entry occur?
	Times	Can Happen
	15	1		7	2	2	<8
	14	1		6	2	2
	13	1		5	3	2
	12	1		4	3	2
	11	1		3	5	3	<4
	10	1		2	7	3
	9	1		1	10*	4
	8	1		1	(technically it can happen up to 15 times, but no decimal number has 15 unique symbols as, base 10.)
	A minimum of 26 bits is required to represent a signature for this problem.
	No bonus for saving bits up to the next power of 2 (32 bits)
	MSB			LSB
	[15-8]	[7-4]	[3,2]	[1,-]
	Highest bits, 1 bit flags for counts, then a byte of 2 bit registers, then the low nibbles, with 0 in the base.
	The lowest index could also be the count of word length.  That information is redundant to the decomposed count of symbol occurrence.

	Initial pass:
		Collect signatures of square numbers that are possible permutation pairs (Later, words which are not numbers might match these...)
		If a square number lacks a pair DELETE it
	Key: ^^^ occurrence count
	[]VALUE:
	Re-map each decimal number to a string:
		From lowest (1s) to highest (Ns), whatever digit occurs uniquely is the next digit (start with 1 since 0 is reserved, for decimal 1..10)
		16 digits can nibble (4bit)-pack into a 64 bit number.
	A similar scheme can be used to occurrence count letters, and then map the pattern to a number format to compare to a list...
	I'd like to sort the list, but it's unlikely to be long so don't bother.


	The input happens to be ALLCAPS; but (In | 0x20) - 'a' would nicely condense english letters...

	+++

	This is very dis-satisfyingly slow.  Mostly the growth rate of squared numbers slows down considerably for larger numbers, contrary to my expectations based on smaller more human-scaled numbers.  Oops.

	The 'anagram' part is comparatively MUCH faster, and currently yields zero results due to a lack of pattern matches that are also anagrams.
	Some choice I made must be incorrect.  I should try again after sleep refreshes the slate.


/
*/

import (
	"bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	"os" // os.Stdout
	// "strconv"
	// "strings"
)

// type SqPatNum struct{ Pat, Sq uint64 }
type SqPatStr struct {
	Pat, Sq uint64
	Str     string
}

func SquarePatterns() [16]map[uint32]map[uint64]uint64 {
	ret := [16]map[uint32]map[uint64]uint64{}
	var qq, sq, place, q, r, pattern uint64 // qq is one ahead in the for loop to make parallel assignment correct
	var stats, found [10]uint8
	var histg uint32
	var uniq, uid, uu uint8
	for uu = 0; uu < 16; uu++ {
		ret[uu] = make(map[uint32]map[uint64]uint64)
	}
	// SquarePatternsOuter:
	for qq, sq = 2, 1; sq < 1_000_000_000_000_000; qq, sq = qq+1, qq*qq {
		stats, found, uniq, pattern = [10]uint8{}, [10]uint8{}, 0, 0
		for place, q = 0, sq; 0 < q; place++ {
			q, r = q/10, q%10
			stats[r]++
			if 1 == stats[r] {
				uniq, uid = uniq+1, uniq
				found[uid] = uint8(r)
			} else {
				for uu = 0; uu < 10; uu++ {
					if uint8(r) == found[uu] {
						uid = uu
						break // 1
					}
				}
			}
			pattern |= uint64(uid) << (place << 2)
		}
		histg, found = 0, [10]uint8{} // Reuse for histogram collection
		for uu = 0; uu < 10; uu++ {
			if 8 <= stats[uu] {
				histg |= 1 << (stats[uu] + 16) // 24 through 31 - can only happen once, no need to track / add
			} else {
				found[stats[uu]]++
			}
		}
		//	MSB			LSB
		//	[15-8]	[7-4]	[3,2]	[1,-]
		// it's easier to just unroll and combine the two loops myself
		histg |= uint32(stats[7])<<22 | uint32(stats[6])<<20 | uint32(stats[5])<<18 | uint32(stats[4])<<16 | uint32(stats[3])<<12 | uint32(stats[2])<<8 | uint32(stats[1])<<4
		if _, exists := ret[place][histg]; !exists {
			ret[place][histg] = make(map[uint64]uint64)
		}
		ret[place][histg][pattern] = sq
	}
	// I think this isn't required.  If it is I'd need to also track multiple numbers with the same pattern but different value
	// for uu = 0; uu < 16; uu++ {
	//	trim := make([]uint32, 0, 8)
	//	for histg, _ := range ret[uu] {
	//		if 2 > len(ret[uu][histg]) {
	//			trim = append(trim, histg)
	//		}
	//	}
	//	for histg = 0; histg < uint32(len(trim)); histg++ {
	//		delete(ret[uu], trim[histg])
	//	}
	//}
	return ret
}

func BytesHistPat(bstr []byte) (uint32, uint64) {
	if 15 < len(bstr) {
		fmt.Printf("Special case / FIXME: %d = len( %s )\n", len(bstr), string(bstr))
		return 0, 0
	}
	var pattern uint64
	var histg uint32
	var found [10]uint8
	var stats [26]uint8
	var uniq, uid, uu, pos, char, place uint8
	for pos = uint8(len(bstr)); 0 < pos; place++ {
		pos--
		char = bstr[pos]
		if ('A' <= char && char < 'Z') || ('a' <= char && char < 'z') {
			char |= 'a' - 'A'
			char -= 'a'
			stats[char]++
			if 1 == stats[char] {
				if 10 <= uniq {
					fmt.Printf("Too many unique characters, skipping: %s\n", string(bstr))
					return 0, 0
				}
				uniq, uid = uniq+1, uniq
				found[uid] = char
			} else {
				for uu = 0; uu < 10; uu++ {
					if char == found[uu] {
						uid = uu
						break // 1
					}
				}
			}
			pattern |= uint64(uid) << (place << 2)
		}
	}
	histg, found = 0, [10]uint8{} // Reuse for histogram collection
	for uu = 0; uu < 26; uu++ {
		if 8 <= stats[uu] {
			histg |= 1 << (stats[uu] + 16) // 24 through 31 - can only happen once, no need to track / add
		} else {
			found[stats[uu]]++
		}
	}
	//	MSB			LSB
	//	[15-8]	[7-4]	[3,2]	[1,-]
	// it's easier to just unroll and combine the two loops myself
	histg |= uint32(stats[7])<<22 | uint32(stats[6])<<20 | uint32(stats[5])<<18 | uint32(stats[4])<<16 | uint32(stats[3])<<12 | uint32(stats[2])<<8 | uint32(stats[1])<<4
	return histg, pattern
}

func InsertSortString(con string) string {
	// This is probably slow.  I should build a library of sorting algorithms soon for sites that don't have "slices" or other standard sorts included.
	lim := len(con)
	ret := make([]byte, 0, lim)
	for ii := 0; ii < lim; ii++ {
		ret = append(ret, byte(con[ii]))
		for jj := ii; 0 < jj && ret[jj] < ret[jj-1]; jj-- {
			ret[jj-1], ret[jj] = ret[jj], ret[jj-1]
		}
	}
	return string(ret)
}

func Euler0098(fn string) uint64 {
	patterns := SquarePatterns()
	// fmt.Println(patterns)
	words := [16]map[uint32][]SqPatStr{}
	for uu := 0; uu < 16; uu++ {
		words[uu] = make(map[uint32][]SqPatStr)
	}

	// Load the words, if they match a square number pattern
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler0098 unable to open: " + fn)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	// split lines is default, use one that chomps all ", (and whitespace) to output 'words'
	scanner.Split(euler.ScannerSplitNLDQ)
	for scanner.Scan() {
		word := scanner.Bytes()
		histg, pat := BytesHistPat(word)
		if _, exists := patterns[len(word)][histg]; exists {
			if sq, exists := patterns[len(word)][histg][pat]; exists {
				// fmt.Printf("DEBUG: %2d %16x add %s\n", len(word), pat, word)
				words[len(word)][histg] = append(words[len(word)][histg], SqPatStr{Pat: pat, Sq: sq, Str: string(word)})
			} else {
				// fmt.Printf("DEBUG: %2d %16x SKIP %s\n", len(word), pat, word)
			}
		}
	}
	fmt.Println(words)

	// Find the biggest square number that's a match
	var ret uint64
	for uu := 15; 0 < uu; uu-- {
		// fmt.Printf("Considering: %d\n%v\n", uu, words[uu])
		// key = histg ; hmVal == []SqPatStr
		for _, hmVal := range words[uu] {
			if 2 > len(hmVal) {
				continue
			}
			anagrams := make(map[string][]SqPatStr)
			for ii, lim := 0, len(hmVal); ii < lim; ii++ {
				anagrams[InsertSortString(hmVal[ii].Str)] = append(anagrams[InsertSortString(hmVal[ii].Str)], hmVal[ii])
			}
			fmt.Printf("Next batch of Anagrams:\n%v\n", anagrams)
			for _, ags := range anagrams {
				if 2 > len(ags) {
					continue
				}
				// Assuming the words are unique, this is a list of at least 2 words which all use the same 10 or less letters...
				possible := ags[0]
				for ii, lim := 1, len(ags); ii < lim; ii++ {
					// Check that assumption
					if possible.Str != ags[ii].Str {
						if possible.Sq < ags[ii].Sq {
							possible = ags[ii]
							ret = ags[ii].Sq
						}
					} else {
						fmt.Printf("DEBUG: skip %d == %d\n", possible.Str, ags[ii].Str)
					}
				}
			}
		}
		if 0 < ret {
			break // return ret
		}
	}

	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 98 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 97: Large Non-Mersenne Prime (2004): 8739992577

real    0m0.099s
user    0m0.136s
sys     0m0.059s
.
*/
func main() {
	var r uint64
	//test
	testISS := []struct{ raw, sorted string }{
		{"abcdefg", "abcdefg"},
		{"cccbbbaaa", "aaabbbccc"},
		{"cabcabcab", "aaabbbccc"},
		{"coco", "ccoo"},
		{"compile", "ceilmop"},
		// {"", ""},
	}
	pass := true
	for _, test := range testISS {
		if test.sorted != InsertSortString(test.raw) {
			fmt.Printf("FAILED: InsertSortString(%s) expected: %s got: %s\n", test.raw, test.sorted, InsertSortString(test.raw))
			pass = false
		}
	}
	if !pass {
		panic("Abort for Debug")
	}

	//run
	r = Euler0098("0098_words.txt")
	fmt.Printf("Euler 98: Anagramic Squares: %d\n", r)
	if 0 != r {
		panic("Did not reach expected value.")
	}
}
