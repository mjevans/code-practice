// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=41
https://projecteuler.net/minimal=41

Pandigital Prime

<p>We shall say that an $n$-digit number is pandigital if it makes use of all the digits $1$ to $n$ exactly once. For example, $2143$ is a $4$-digit pandigital and is also prime.</p>
<p>What is the largest $n$-digit pandigital prime that exists?</p>




*/
/*

I saw Pandigital in the title and figured it was more validating Pandigital numbers...

However after reading the problem... Oops, it'll be far faster to construct them and then use the faster factorization shot in the dark function to test if suitable candidates are prime.

1-9 without 0, so 987_654_321 is the largest possible number to test...

The squareroot of that number is (rounded up) 31426, so it wouldn't even be _that_ tough to brute force factor numbers in this range.  Though that's still over 10000 modulus operations per candidate.

One other thing... a pandigital number comprised of the digits 1..9 totals 45 which totals 9 which is divisible by 3, and thus, so are all the 9 digit pandigital numbers.

Try #2 8 digit pandigital numbers.

Try #3 7 digit pandigital numbers.

*/

import (
	// "bufio"
	// "bitvector"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler041_deck_fail() uint64 {
	// p := euler.NewBVPrimes()
	// p.Grow(uint(limit))
	// func Factor1980AutoPMC(q uint, singlePrimeOnly bool) uint
	// uses euler.Primes to quickly test if a number is a _known_ prime, but that only matters if fully factoring.

	const decksize = 8

	deck := make([]uint8, decksize, decksize)

	// for nna := 0; decksize > nna; nna++ {
	for nnb := 0; decksize-1 > nnb; nnb++ {
		for nnc := 0; decksize-2 > nnc; nnc++ {
			for nnd := 0; decksize-3 > nnd; nnd++ {
				for nne := 0; decksize-4 > nne; nne++ {
					for nnf := 0; decksize-5 > nnf; nnf++ {
						for nng := 0; decksize-6 > nng; nng++ {
							for ii := 0; decksize > ii; ii++ {
								deck[ii] = uint8(decksize - ii)
							}
							// can := uint64(10_000_000) * uint64(euler.SlicePopUint8(deck, nna))
							can := uint64(1_000_000) * uint64(euler.SlicePopUint8(deck, nnb))
							can += uint64(100_000) * uint64(euler.SlicePopUint8(deck, nnc))
							can += uint64(10_000) * uint64(euler.SlicePopUint8(deck, nnd))
							can += uint64(1_000) * uint64(euler.SlicePopUint8(deck, nne))
							can += uint64(100) * uint64(euler.SlicePopUint8(deck, nnf))
							can += uint64(10) * uint64(euler.SlicePopUint8(deck, nng))
							can += uint64(euler.SlicePopUint8(deck, 0))
							// 2s check already a guard for
							if uint(can) == euler.Factor1980AutoPMC(uint(can), false) {
								fmt.Printf("Found: %d\n", can)
								return can
							}

						}
					}
				}
			}
		}
	}
	// }

	return 0
}

func Euler041() uint64 {
	// func Factor1980AutoPMC(q uint, singlePrimeOnly bool) uint
	// uses euler.Primes to quickly test if a number is a _known_ prime, but that only matters if fully factoring.

	const decksize = 8

	// ii |= 1 // always odd
	for ii := uint64(87_654_321); 1 < ii; ii -= 2 {
		// Pandigital test is probably faster than prime test
		// func Pandigital(test uint64, used uint16) (biton, usedRe uint16, DigitShift uint64) {
		fullPD, _, _, _ := euler.Pandigital(ii, 0)
		if fullPD {
			if uint(ii) == euler.Factor1980AutoPMC(uint(ii), false) {
				fmt.Printf("Found: %d\n", ii)
				return ii
			}
		}
	}
	return 0
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 41 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Found: 7652413
Euler041: Pandigital Prime :    7652413


*/
func main() {
	//test

	//run
	ans := Euler041()
	fmt.Printf("Euler041: Pandigital Prime :\t%d\n", ans)
}
