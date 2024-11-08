// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=76
https://projecteuler.net/minimal=76

<p>It is possible to write five as a sum in exactly six different ways:</p>
\begin{align}
&amp;4 + 1\\
&amp;3 + 2\\
&amp;3 + 1 + 1\\
&amp;2 + 2 + 1\\
&amp;2 + 1 + 1 + 1\\
&amp;1 + 1 + 1 + 1 + 1
\end{align}
<p>How many different ways can one hundred be written as a sum of at least two positive integers?</p>


/
*/
/*
/	My first thoughts were to consider that set of 100 individual ones and how they can link.  However the information about where those single ones exist is extra / ignored.  2 + 1 + 1 + 1 is the unique sorted order set.  So thinking of a shape based on the links in order won't work.

	Bins of objects?  Like 0 bins for 100 ones, 1 bin for 2 + (100-bin1) up to 99 + 1?  Etc.  This might work from a formal math stance but for a program rolling (n/2)+1 counters to enumerate up to 2 x 50 (2s) just doesn't seem right.  Could be improved by a multiplier for the bin size, but then counting becomes a total nightmare and intuitively it still seems too inelegant to be the proper solution.

	For such a deceptively simple problem statement, this is starting to sound more like some variation of Knapsack problem.
	https://en.wikipedia.org/wiki/Knapsack_problem
	https://en.wikipedia.org/wiki/Subset_sum_problem "The analogous counting problem #SSP, which asks to enumerate the number of subsets summing to the target, is #P-complete.[4]"
	https://en.wikipedia.org/wiki/%E2%99%AFP

	Maybe that inelegant idea I had is one of the better methods, it at least allows an algorithm that keeps the sizes sorted greatest to smallest.

	This problem was in the back of my mind while I bought some food and, though it's far more common to use credit/debit cards similar, I was reminded of change.
	https://en.wikipedia.org/wiki/Change-making_problem

	The large number of possible 'coins' make this problem both difficult to spot as a coins problem, and extremely tedious to consider with the typical coin value matrix.

	Also interesting...
	https://en.wikipedia.org/wiki/Coin_problem#McNugget_numbers
	https://en.wikipedia.org/wiki/Integer_partition

	This is an Integer Partition problem
	https://en.wikipedia.org/wiki/Integer_partition#Partition_function
	https://en.wikipedia.org/wiki/Partition_function_(number_theory)  Oh there's an answer on this page, but I don't want to cheat... However I will add the result as the desired test case.

	A good refresher example for Dynamic Programming https://algorithmist.com/wiki/Coin_change  Divide to simplify when 'optimal substructure' is possible and recurs to solve
/
*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func ChangeCombos(coins []uint8, total int32) uint64 {
	if 0 > total || 0 == len(coins) {
		return 0
	}
	if 0 == total {
		return 1
	}
	return ChangeCombos(coins[:len(coins)-1], total) + ChangeCombos(coins, total-int32(coins[len(coins)-1]))
}

func Euler0076_dynamic_no_cache(total int32) uint64 {
	coins := make([]uint8, 0, total)

	// 1 to (total - 1)
	for ii := int32(1); ii < total; ii++ {
		coins = append(coins, uint8(ii))
	}

	return ChangeCombos(coins, total)
}

func Euler0076(total int16) uint64 {
	coins := make([]uint8, 0, total)
	cache := make(map[uint32]uint64)

	var CachableCombos func(maxcoin, total int16) uint64
	CachableCombos = func(maxcoin, total int16) uint64 {
		if 0 > total || 0 == maxcoin {
			return 0
		}
		if 0 == total {
			return 1
		}
		var val uint64
		key := uint32(maxcoin)<<16 | uint32(total)
		if val, exists := cache[key]; exists {
			return val
		}
		val = CachableCombos(maxcoin-1, total) + CachableCombos(maxcoin, total-int16(coins[maxcoin-1]))
		cache[key] = val
		return val
	}

	// 1 to (total - 1)
	for ii := int16(1); ii < total; ii++ {
		coins = append(coins, uint8(ii))
	}

	return CachableCombos(int16(len(coins)), total)
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 76 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

// This is good enough for Euler but it's probably not fast enough for competitive answers.
Euler 76: :     Count: 190569291

real    0m8.244s
user    0m8.290s
sys     0m0.056s

The cache-able (memoized) version is FAR (almost 100x) faster since it dramatically cuts recursions over small values, given the number of recursions placing the coins array in the closure's inherited scope and only passing the maximum index saves at least 2 registers / stack size per call.

Euler 76: :     Count: 190569291

real    0m0.106s
user    0m0.142s
sys     0m0.064s

Euler 76: :     Count: 190569291

real    0m0.095s
user    0m0.157s
sys     0m0.034s

.
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	r := Euler0076(5)
	if 6 != r {
		panic(fmt.Sprintf("Euler 76: Expected 6 got %d", r))
	}

	//run
	r = Euler0076(100)
	fmt.Printf("Euler 76: :\tCount: %d\n", r)
	// Note: this is NOT the accepted answer, because the solution with a single 100 coin is NOT allowed
	// https://en.wikipedia.org/wiki/Partition_function_(number_theory)
	// Some exact values: p(100) 190_569_292, p(1000) 24_061_467_864_032_622_473_692_149_727_991 (E31) , the check digits for p(10000) are listed: 36_167_251_325_..._906_916_435_144 (E106)
	if 190_569_291 != r {
		panic("Did not reach expected value.")
	}
}
