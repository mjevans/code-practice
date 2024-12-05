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

	+++

	In retrospect 'seen' and 'found' are easily confused... bug fixed.  seen renamed to charOccur
	Found is reused for histogram occurrence later... but naming it occurrence would return the confusion.

Next batch of Anagrams: 90
[{36344967696 923187456 INTRODUCE} {36344967696 923187456 REDUCTION}]
Set new best answer: ...        [9]     [876543210] 83
{36344967696 923187456 REDUCTION}
{36344967696 923187456 INTRODUCE}
Euler 98: Anagramic Squares: 923187456
panic: Did not reach expected value.

	Ha ha, nice hint... it actually is a good hint though.

	What I think is still wrong:
	The numbers have the same 'shape' pattern, but aren't anagrams!  They wanted anagram numbers too?! Yes, I can see how that's a far inferrable in the description, but it's so obtusely described and the entire concept of 'numeric anagrams' is new to me so it was very easy to overlook.  It changes numbers in to strings of characters that happen to represent a decimal value of an integer.

	Not quite correct; length, histogram and pattern are all generic fingerprint elements; they're the _shape_ of the data.
	[16]map[uint32]map[uint64][]uint64{} // [len] [histogram] [pattern(exact)]  [number of that pattern found (tail will be the largest)]
	However the content matters too.  It doesn't matter for adding a word to a list of anagrams, or as I know realize, parallel matched anagrams.
	Oh the other end I'll have to bin the words that match a pattern and thing against the numeric anagram match too.

	After a full review and no obvious faults in any of the other functions, this evening I was sure it had to be in the match function.
	Without a vaccine booster shot headache (like I had last night) I was able to properly focus.
	The mapping process had to be from word, and division result (encode), through another word (decode), as the needle to find if it existed.
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

type SqSortedSq struct{ Sorted, Sq, Nib uint64 }
type SqPatStr struct {
	Pat, Sq uint64
	Str     string
}

func SquarePatterns(lim uint64) [16]map[uint32]map[uint64][]SqSortedSq {
	ret := [16]map[uint32]map[uint64][]SqSortedSq{}       // [len] [histogram] [pattern(exact)]  []SqSortedSq{ Sorted, Sq uint64 }
	var qq, sq, place, q, r, pattern, anagram, nib uint64 // qq is one ahead in the for loop to make parallel assignment correct
	var charOccur, found [10]uint8
	var histg uint32
	var uniq, uid, uu uint8
	for uu = 0; uu < 16; uu++ {
		ret[uu] = make(map[uint32]map[uint64][]SqSortedSq)
	}
	// SquarePatternsOuter:
	fmt.Printf("SquarePatterns up to %d\n", lim)
	for qq, sq = 2, 1; sq < lim; qq, sq = qq+1, qq*qq {
		numsort := make([]byte, 0, 10)
		charOccur, found, uniq, pattern, nib = [10]uint8{}, [10]uint8{}, 0, 0, 0
		for place, q = 0, sq; 0 < q; place++ {
			q, r = q/10, q%10
			charOccur[r]++
			if 1 == charOccur[r] {
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
			nib |= uint64(r) << (place << 2)
			pattern |= uint64(uid) << (place << 2)
			// insert-sort the next found number
			numsort = append(numsort, byte(r))
			for jj := place; 0 < jj && numsort[jj] < numsort[jj-1]; jj-- {
				numsort[jj-1], numsort[jj] = numsort[jj], numsort[jj-1]
			}
		}
		for anagram, r = 0, 0; r < place; r++ {
			anagram |= uint64(numsort[r]) << (r << 2)
		}
		histg, found = 0, [10]uint8{} // Reuse for histogram collection
		for uu = 0; uu < 10; uu++ {
			if 8 <= charOccur[uu] {
				histg |= 1 << (charOccur[uu] + 16) // 24 through 31 - can only happen once, no need to track / add
			} else if 0 != charOccur[uu] {
				found[charOccur[uu]]++
			}
		}
		//	MSB			LSB
		//	[15-8]	[7-4]	[3,2]	[1,-]
		// it's easier to just unroll and combine the loop myself
		histg |= (uint32(found[7]) << 22) | (uint32(found[6]) << 20) | (uint32(found[5]) << 18) | (uint32(found[4]) << 16) | (uint32(found[3]) << 12) | (uint32(found[2]) << 8) | (uint32(found[1]) << 4)
		if _, exists := ret[place][histg]; !exists {
			ret[place][histg] = make(map[uint64][]SqSortedSq)
		}
		ret[place][histg][pattern] = append(ret[place][histg][pattern], SqSortedSq{Sorted: anagram, Sq: sq, Nib: nib})
	}
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
	var charOccur [26]uint8
	var uniq, uid, uu, pos, char, place uint8
	for pos = uint8(len(bstr)); 0 < pos; place++ {
		pos--
		char = bstr[pos]
		if ('A' <= char && char < 'Z') || ('a' <= char && char < 'z') {
			char |= 'a' - 'A'
			char -= 'a'
			charOccur[char]++
			if 1 == charOccur[char] {
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
	// fmt.Println(charOccur)
	histg, found = 0, [10]uint8{} // Reuse for histogram collection
	for uu = 0; uu < 26; uu++ {
		if 8 <= charOccur[uu] {
			histg |= 1 << (charOccur[uu] + 16) // 24 through 31 - can only happen once, no need to track / add
		} else if 0 != charOccur[uu] {
			found[charOccur[uu]]++
		}
	}
	// fmt.Printf("%x\n", found[1]<<4)
	//	MSB			LSB
	//	[15-8]	[7-4]	[3,2]	[1,-]
	// it's easier to just unroll and combine the two loops myself
	histg |= (uint32(found[7]) << 22) | (uint32(found[6]) << 20) | (uint32(found[5]) << 18) | (uint32(found[4]) << 16) | (uint32(found[3]) << 12) | (uint32(found[2]) << 8) | (uint32(found[1]) << 4)
	return histg, pattern
}

// This could be inserted into the loop above like for the numbers... but might be reusable elsewhere too and it's already a function.
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

func E0098HashString(con string) string {
	return InsertSortString(con)
}

// Load a list of ", separated words from a file for Euler 98, binned by size with 0 as 'special case' length words
func E0098LoadWords(fn string) ([16]map[string][]string, int) {
	ret := [16]map[string][]string{}
	for ii := 0; ii < 16; ii++ {
		ret[ii] = make(map[string][]string)
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
		word := scanner.Text()
		wlen, whash := len(word), E0098HashString(word)
		if 15 < wlen {
			wlen = 0
		}
		ret[wlen][whash] = append(ret[wlen][whash], word)
	}
	// Cull any 'word hashes' which lack even a normal anagram partner
	var longest int
	for ii := 0; ii < 16; ii++ {
		dl := make([]string, 0, 16)
		for k, v := range ret[ii] {
			if 2 > len(v) {
				dl = append(dl, k)
			}
		}
		for _, vk := range dl {
			delete(ret[ii], vk)
		}
		if 0 < len(ret[ii]) {
			longest = ii
		}
	}
	return ret, longest
}

func E0098AnagramTransposition(source, dest string, sourceNibNum uint64) uint64 {
	// Euler 98 example
	// CARE =>	RACE
	// 1296	=>	9216
	lm := len(source)
	lutCD := [26]uint8{}
	// load the lookup table between Char and Decimal
	for lim := lm; 0 < lim; {
		lim--
		lutCD[(source[lim]|('a'-'A'))-'a'] = uint8(sourceNibNum & 0xF)
		sourceNibNum >>= 4
	}
	// construct the output number
	var ret uint64
	for ii := 0; ii < lm; ii++ {
		ret = ret*10 + uint64(lutCD[(dest[ii]|('a'-'A'))-'a'])
	}
	return ret
}

func E0098LargestSquareAnagramTransposition(patterns [16]map[uint32]map[uint64][]SqSortedSq, words [16]map[string][]string) uint64 {
	var ret uint64
	for uu := 15; 0 < uu; uu-- {
		// fmt.Printf("Considering: %d\n%v\n", uu, words[uu])
		// key = histg ; hmVal == []SqPatStr
		for _, hmVal := range words[uu] {
			var pat uint64
			var histg uint32 // same for all of the anagrams

			// For every word in an anagram set
			for ss, lim := 0, len(hmVal); ss < lim; ss++ {
				// calc it's histogram and number pattern
				histg, pat = BytesHistPat([]byte(hmVal[ss]))

				// Only if there's a number that matches those patterns
				if sl, exists := patterns[uu][histg][pat]; exists {
					slLm := len(sl)

					// For every word _after_ that word
					for dd := ss + 1; dd < lim; dd++ {

						// Consider each possible input number
						// [16]map[uint32]map[uint64][]SqSortedSq{} // [len] [histogram] [pattern(exact)]  []SqSortedSq{ Sorted, Sq uint64 }
						for aa := 0; aa < slLm; aa++ {

							// Use the word pair and the (cached) division results to calculate the 'needle' to look for.
							needle := E0098AnagramTransposition(hmVal[ss], hmVal[dd], sl[aa].Nib)
							if ret > needle {
								continue
							}
							// Needle must be greater than return, otherwise there's no point, it must also come after this number, since the mappings work from both ways with the start word.
							for bb := aa + 1; bb < slLm; bb++ {
								if needle == sl[bb].Sq {
									ret = needle
									break
								}
							}
						}
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

func Euler0098(fn string) uint64 {
	words, longest := E0098LoadWords(fn)
	patterns := SquarePatterns(euler.PowInt(10, uint64(longest)))
	_, _ = words, patterns
	//	[16]map[uint32]map[uint64][]uint64{}	//	[len] [hash] [pattern(exact)] [number of that pattern found (tail will be the largest)]
	// fmt.Println(patterns)
	fmt.Println(words)
	// fmt.Println(patterns[9][0x90][0x876543210])

	return E0098LargestSquareAnagramTransposition(patterns, words)
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 98 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

SquarePatterns up to 10000
SquarePatterns up to 1000000000
[map[] map[] map[NO:[NO ON]] map[ACT:[ACT CAT] AET:[EAT TEA] DGO:[DOG GOD] HOW:[HOW WHO] IST:[ITS SIT] NOW:[NOW OWN]] map[ACER:[CARE RACE] ADEL:[DEAL LEAD] AEHT:[HATE HEAT] AELM:[MALE MEAL] AEMN:[MEAN NAME] AENR:[EARN NEAR] AERT:[RATE TEAR] AEST:[EAST SEAT] EFIL:[FILE LIFE] EIMT:[ITEM TIME] ENOT:[NOTE TONE] ERSU:[SURE USER] FMOR:[FORM FROM] GINS:[SIGN SING] HSTU:[SHUT THUS] OPST:[POST SPOT STOP]] map[ABDOR:[BOARD BROAD] AEHPS:[PHASE SHAPE] AEHRT:[EARTH HEART] AEIRS:[ARISE RAISE] AELST:[LEAST STEAL] EEHST:[SHEET THESE] EIQTU:[QUIET QUITE] GHINT:[NIGHT THING] HORTW:[THROW WORTH] HOSTU:[SHOUT SOUTH]] map[ADEGNR:[DANGER GARDEN] CDEIRT:[CREDIT DIRECT] CEENRT:[CENTRE RECENT] CEEPTX:[EXCEPT EXPECT] CEORSU:[COURSE SOURCE] EFMORR:[FORMER REFORM] EGINOR:[IGNORE REGION]] map[] map[ACEINORT:[CREATION REACTION]] map[CDEINORTU:[INTRODUCE REDUCTION]] map[] map[] map[] map[] map[] map[]]
Euler 98: Anagramic Squares: 18769

real    0m0.146s
user    0m0.199s
sys     0m0.072s
*/
func main() {
	var r uint64
	pass := true
	//test
	r = E0098AnagramTransposition("CARE", "RACE", 0x1296)
	if 9216 != r {
		fmt.Printf("Expected 9216 from E0098AnagramTransposition got: %d\n", r)
	}
	testISS := []struct{ raw, sorted string }{
		{"abcdefg", "abcdefg"},
		{"cccbbbaaa", "aaabbbccc"},
		{"cabcabcab", "aaabbbccc"},
		{"coco", "ccoo"},
		{"compile", "ceilmop"},
		// {"", ""},
	}
	for _, test := range testISS {
		if test.sorted != InsertSortString(test.raw) {
			fmt.Printf("FAILED: InsertSortString(%s) expected: %s got: %s\n", test.raw, test.sorted, InsertSortString(test.raw))
			pass = false
		}
	}
	words := [16]map[string][]string{}
	for ii := 0; ii < 16; ii++ {
		words[ii] = make(map[string][]string)
	}
	testWords := []string{"CARE", "RACE"}
	for _, word := range testWords {
		wlen, whash := len(word), E0098HashString(word)
		words[wlen][whash] = append(words[wlen][whash], word)
	}
	patterns := SquarePatterns(euler.PowInt(10, uint64(4)))
	// fmt.Println(patterns[4][0x40])
	r = E0098LargestSquareAnagramTransposition(patterns, words)
	if 9216 != r {
		fmt.Printf("Expected 9216 got %d\n", r)
		pass = false
	}

	if !pass {
		panic("Abort for Debug")
	}

	//run
	r = Euler0098("0098_words.txt")
	fmt.Printf("Euler 98: Anagramic Squares: %d\n", r)
	// NOT correct 923187456
	if 18769 != r {
		panic("Did not reach expected value.")
	}
}
