// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=60
https://projecteuler.net/minimal=60

<p>The primes $3$, $7$, $109$, and $673$, are quite remarkable. By taking any two primes and concatenating them in any order the result will always be prime. For example, taking $7$ and $109$, both $7109$ and $1097$ are prime. The sum of these four primes, $792$, represents the lowest sum for a set of four primes with this property.</p>
<p>Find the lowest sum for a set of five primes for which any two primes concatenate to produce another prime.</p>



*/
/*

Counter clock / hill climb for 5 slots to cover all combinations.

The highest slot can start at 673 since that was the lowest set that worked for 4 digits...

Evaluating all of the combinations is even slower than expected; I'd hoped between ConcatDigitsU64 using small inputs and thus few mul/div and all the primes of concern already known this would be trivial.

Can't eliminate whole sets of numbers since individual components might be reused in a solution...

Testing combinations of individual pairs is the inner most loop and thus most likely candidate for optimization.


Use an array of pages of answers, to reduce copying as more terms are added.

map[BIGPRIME]([]lesserprimeSorted) ??

Workshop the data-flow first to refine what the cache shape.

Offhand PROBABLY fits in uint16 space (primes <65535)
Just store things for both the low and the high end to aid matching
map[uint16]([]uint16)

2 can't be in the list because it will always fail the low side (Even number) ; could do the (n -3) >> 1 trick



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

type CoCatPrimeSet struct {
	Tested   map[uint16][]uint16
	Idx, Sub uint16
}

func NewCoCatPrimeSet() *CoCatPrimeSet {
	return &CoCatPrimeSet{Tested: make(map[uint16][]uint16), Idx: 2, Sub: 2}
}

/*
func (co *CoCatPrimeSet) NextCoCatPrime() (uint16, uint16) {
	for {
		co.Sub = uint16(euler.Primes.PrimeAfter(uint(co.Sub)))
	if co.Idx <= co.Sub {
		co.Sub = 2 // Idx is rechecked so becomes 3 in the new cycle
		tmp := euler.Primes.PrimeAfter(uint(co.Idx))
		if 0 < tmp&0xFFFF {
			panic("Overflowed uint16 design limit")
		}
		co.Idx = uint16(tmp)
		continue
	}
		is := uint(euler.ConcatDigitsU64(uint64(co.Idx), uint64(co.Sub), 10))
		si := uint(euler.ConcatDigitsU64(uint64(co.Sub), uint64(co.Idx), 10))
		if euler.Primes.ProbPrime(is) && euler.Primes.ProbPrime(si) {
			co.Tested[co.Sub] = append(co.Tested[co.Sub], co.Idx)
			co.Tested[co.Idx] = append(co.Tested[co.Idx], co.Sub)
			return co.Idx, co.Sub
		}
	}
	//return 0, 0
}
*/

func (co *CoCatPrimeSet) NextCoCatSlice() (uint16, []uint16) {
	var sub uint16
	for {
		tmp := euler.Primes.PrimeAfter(uint(co.Idx))
		if (tmp & 0xFFFF) < uint(co.Idx) {
			fmt.Printf("Got prime %d\n", tmp)
			panic("Overflowed uint16 design limit")
		}
		co.Idx = uint16(tmp)
		co.Tested[co.Idx] = make([]uint16, 0, 16)
		sub = 3
		for sub < co.Idx {
			is := uint(euler.ConcatDigitsU64(uint64(co.Idx), uint64(sub), 10))
			si := uint(euler.ConcatDigitsU64(uint64(sub), uint64(co.Idx), 10))
			if euler.Primes.ProbPrime(is) && euler.Primes.ProbPrime(si) {
				co.Tested[sub] = append(co.Tested[sub], co.Idx)
				co.Tested[co.Idx] = append(co.Tested[co.Idx], sub)
			}
			sub = uint16(euler.Primes.PrimeAfter(uint(sub)))
		}
		if 0 != len(co.Tested[co.Idx]) {
			break
		}
	}
	return co.Idx, co.Tested[co.Idx]
}

func Euler060() uint {
	euler.Primes.Grow(4096)
	var combos uint
	var pl []uint16
	p := [5]uint16{}
	co := NewCoCatPrimeSet()
	// Largest Prime Loop
	for {
		p[0], pl = co.NextCoCatSlice()
		pllen := len(pl)
		for ii3 := 3; ii3 < pllen; ii3++ {
			p[4] = pl[ii3]
			for ii2 := 2; ii2 < ii3; ii2++ {
				p[3] = pl[ii2]
				for ii1 := 1; ii1 < ii2; ii1++ {
					p[2] = pl[ii1]
				Euler060failed:
					for ii0 := 0; ii0 < ii1; ii0++ {
						p[1] = pl[ii0]
						combos++
						for ii := 1; ii < 5; ii++ {
							for jj := ii + 1; jj < 5; jj++ {
								if -1 == euler.BsearchSlice(co.Tested[p[ii]], p[jj]) {
									continue Euler060failed
								}
							}
						}
						for ii := 1; ii < 5; ii++ {
							for jj := ii + 1; jj < 5; jj++ {
								fmt.Printf("Validate: %d %d = %d => %t\n", p[ii], p[jj], euler.ConcatDigitsU64(uint64(p[ii]), uint64(p[jj]), 10), euler.Primes.ProbPrime(uint(euler.ConcatDigitsU64(uint64(p[ii]), uint64(p[jj]), 10))))
								fmt.Printf("Validate: %d %d = %d => %t\n", p[jj], p[ii], euler.ConcatDigitsU64(uint64(p[jj]), uint64(p[ii]), 10), euler.Primes.ProbPrime(uint(euler.ConcatDigitsU64(uint64(p[jj]), uint64(p[ii]), 10))))
							}
						}
						sum := uint(p[0]) + uint(p[1]) + uint(p[2]) + uint(p[3]) + uint(p[4])
						fmt.Printf("Found: %d = %v\n", sum, p)
						return sum
					}
				}
			}
		}
		// fmt.Printf("Tested %d failed combos, used up %d\n", combos, p[0])
	}
}

// This was WAY too slow
func Euler060_slowmode() uint {
	euler.Primes.Grow(4096)
	var combos uint
	p := [5]uint{}
	for p[0] = 673; ; p[0] = euler.Primes.PrimeAfter(p[0]) {
		for p[1] = 7; p[1] < p[0]; p[1] = euler.Primes.PrimeAfter(p[1]) {
			for p[2] = 5; p[2] < p[1]; p[2] = euler.Primes.PrimeAfter(p[2]) {
				for p[3] = 3; p[3] < p[2]; p[3] = euler.Primes.PrimeAfter(p[3]) {
				Euler060_slowmodefailed:
					for p[4] = 2; p[4] < p[3]; p[4] = euler.Primes.PrimeAfter(p[4]) {
						combos++
						for ii := 1; ii < 5; ii++ {
							for jj := ii + 1; jj < 5; jj++ {
								if false == euler.Primes.ProbPrime(uint(euler.ConcatDigitsU64(uint64(p[ii]), uint64(p[jj]), 10))) || false == euler.Primes.ProbPrime(uint(euler.ConcatDigitsU64(uint64(p[jj]), uint64(p[ii]), 10))) {
									continue Euler060_slowmodefailed
								}
							}
						}
						for ii := 1; ii < 5; ii++ {
							for jj := ii + 1; jj < 5; jj++ {
								fmt.Printf("Validate: %d %d = %d => %t\n", p[ii], p[jj], euler.ConcatDigitsU64(uint64(p[ii]), uint64(p[jj]), 10), euler.Primes.ProbPrime(uint(euler.ConcatDigitsU64(uint64(p[ii]), uint64(p[jj]), 10))))
								fmt.Printf("Validate: %d %d = %d => %t\n", p[jj], p[ii], euler.ConcatDigitsU64(uint64(p[jj]), uint64(p[ii]), 10), euler.Primes.ProbPrime(uint(euler.ConcatDigitsU64(uint64(p[jj]), uint64(p[ii]), 10))))
							}
						}
						sum := p[0] + p[1] + p[2] + p[3] + p[4]
						fmt.Printf("Found: %d = %v\n", sum, p)
						return sum
					}
				}
			}
		}
		fmt.Printf("Tested %d failed combos, used up %d\n", combos, p[0])
	}
	return 0
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 60 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Limited list (only coprimes with the current lead prime) is still slow, but reasonable?
Validate: 13 5197 = 135197 => true
Validate: 5197 13 = 519713 => true
Validate: 13 5701 = 135701 => true
Validate: 5701 13 = 570113 => true
Validate: 13 6733 = 136733 => true
Validate: 6733 13 = 673313 => true
Validate: 5197 5701 = 51975701 => true
Validate: 5701 5197 = 57015197 => true
Validate: 5197 6733 = 51976733 => true
Validate: 6733 5197 = 67335197 => true
Validate: 5701 6733 = 57016733 => true
Validate: 6733 5701 = 67335701 => true
Found: 26033 = [8389 13 5197 5701 6733]
Euler 60: Prime Pair Sets: 26033

real    0m37.952s
user    0m38.797s
sys     0m0.386s

Use the cache?  that is ~11x faster
Validate: 13 5197 = 135197 => true
Validate: 5197 13 = 519713 => true
Validate: 13 5701 = 135701 => true
Validate: 5701 13 = 570113 => true
Validate: 13 6733 = 136733 => true
Validate: 6733 13 = 673313 => true
Validate: 5197 5701 = 51975701 => true
Validate: 5701 5197 = 57015197 => true
Validate: 5197 6733 = 51976733 => true
Validate: 6733 5197 = 67335197 => true
Validate: 5701 6733 = 57016733 => true
Validate: 6733 5701 = 67335701 => true
Found: 26033 = [8389 13 5197 5701 6733]
Euler 60: Prime Pair Sets: 26033

real    0m3.265s
user    0m3.331s
sys     0m0.124s


*/
func main() {
	//test

	//run
	fmt.Printf("Euler 60: Prime Pair Sets: %v\n", Euler060())
}
