// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=5
https://projecteuler.net/minimal=5

<p>$2520$ is the smallest number that can be divided by each of the numbers from $1$ to $10$ without any remainder.</p>
<p>What is the smallest positive number that is <strong class="tooltip">evenly divisible<span class="tooltiptext">divisible with no remainder</span></strong> by all of the numbers from $1$ to $20$?</p>


*/

import (
	"fmt"
	// "slices" // Doh not in 1.19
	"sort"
	"strings"
	// "os" // os.Stdout
)

func Factor(primes []int, num int) []int {
	// Public school factoring algorithm from memory...

	// With a list of known primes, the largest number that can be factored is Pn * Pn
	for ; nil == primes || num > primes[len(primes)-1]*primes[len(primes)-1]; primes = GetPrimes(primes, 0) {
		// fmt.Println(len(primes), primes[len(primes)-1])
	}

	ret := []int{}
	if num < 2 {
		return ret
	}
	for _, prime := range primes {
		for ; 0 == num%prime; num /= prime {
			ret = append(ret, prime)
		}
		if num < prime*prime {
			break
		} // break if no more prime factors are possible
	}
	if 1 < num {
		ret = append(ret, num)
	}
	// fmt.Println("Factor:\t", num, "\n", ret, primes)
	return ret
}

func GetPrimes(primes []int, primehunt int) []int {
	if nil == primes {
		primes = []int{2, 3, 5, 7, 11, 13, 17, 19}
	}
	// Semi-arbitrary expansion target, find 8 more primes (8, 16, 32, 64 it'll fit within the append growth algo)
	if primehunt < 1 {
		primehunt = 8
	}
PrimeHunt:
	for ; 0 < primehunt; primehunt-- {
		for ii := primes[len(primes)-1] + 1; ; ii++ {
			result := Factor(primes, ii)
			if 1 == len(result) && primes[len(primes)-1] < result[0] {
				//fmt.Println("Found Prime:\t", result[0])
				primes = append(primes, result[0])
				continue PrimeHunt // I could break once, but this documents the intent
			}
		}
	}
	return primes
}

func PrintFactors(factors []int) {
	// Join only takes []string s? fff
	strFact := make([]string, len(factors), len(factors))
	for ii, val := range factors {
		strFact[ii] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strFact, ", "))
}

func ListMul(scale []int) int {
	ret := 1
	for _, val := range scale {
		ret *= val
	}
	return ret
}

func IsPalindrome(num int) bool {
	digits := make([]int, 0, 8)
	for num != 0 {
		digits = append(digits, num%10)
		num /= 10
	}
	// 0 1 2 3 4 5 .. 6
	for ii := 0; ii <= len(digits)/2; ii++ {
		if digits[ii] != digits[len(digits)-1-ii] {
			return false
		}
	}
	return true
}

// CompactInts should behave like slices.Compact(slices.Sort())
func CompactInts(arr []int) []int {
	sort.Ints(arr)
	last := 0
	knext := 0
CompactIntsOuter:
	for k := 0; k < len(arr); k++ {
		// fmt.Println("Arr: ", k, " = ", arr[k], arr)

		// // Happy Path, no dupes, copy mode
		if last < arr[k] {
			last = arr[k]
			continue
		}

		// // Eat Dupes
		// Always pull from ahead
		if knext < k {
			knext = k + 1
		}
		arr[k] = 0 // Zero until replaced
		// fmt.Println("Dup: ", k, " = ", arr[k], "(", arr, ")", knext, knext-k)
		for arr[k] <= last {
			// If the end of the array, calculate the skip and store in knext for the slice
			if knext >= len(arr) {
				knext = knext - k
				break CompactIntsOuter
			}
			// Found Next, pull it back, tested good so advance past.
			if last < arr[knext] {
				arr[k] = arr[knext]
				last = arr[knext]
				k++
			}
			arr[knext] = 0
			knext++
		}
	}
	fmt.Println(knext, arr)
	arr = arr[:len(arr)-knext]
	return arr
}

func PrimeLCD(a, b []int) []int {
	var pa, pb int
	var ret []int
	for {
		if pa < len(a) && pb < len(b) {
			if a[pa] <= b[pb] {
				ret = append(ret, a[pa])
				if a[pa] == b[pb] {
					pb++
				}
				pa++
			} else {
				ret = append(ret, b[pb])
				pb++
			}
		} else { // Take the remaining array
			if pa < len(a) {
				ret = append(ret, a[pa:]...)
			}
			if pb < len(b) {
				ret = append(ret, b[pb:]...)
			}
			break
		}
	}
	// fmt.Println("Prime LCD\n", a, "\n", b, "\n", ret)
	return ret
}

/* Sort Notes
	https://en.wikipedia.org/wiki/Introsort#pdqsort
Pseudocode

If a heapsort implementation and partitioning functions of the type discussed in the quicksort article are available, the introsort can be described succinctly as

procedure sort(A : array):
    maxdepth ← ⌊log2(length(A))⌋ × 2
    introsort(A, maxdepth)

procedure introsort(A, maxdepth):
    n ← length(A)
    if n < 16:
        insertionsort(A)
    else if maxdepth = 0:
        heapsort(A)
    else:
        p ← partition(A)  // assume this function does pivot selection, p is the final position of the pivot
        introsort(A[1:p-1], maxdepth - 1)
        introsort(A[p+1:n], maxdepth - 1)

The factor 2 in the maximum depth is arbitrary; it can be tuned for practical performance. A[i:j] denotes the array slice of items i to j including both A[i] and A[j]. The indices are assumed to start with 1 (the first element of the A array is A[1]).

pdqsort

Pattern-defeating quicksort (pdqsort) is a variant of introsort incorporating the following improvements:[8]

    Median-of-three pivoting,
    "BlockQuicksort" partitioning technique to mitigate branch misprediction penalities,
    Linear time performance for certain input patterns (adaptive sort),
    Use element shuffling on bad cases before trying the slower heapsort.

pdqsort is used by Rust, GAP,[9] and the C++ library Boost.[10]


https://en.wikipedia.org/wiki/Timsort

https://en.wikipedia.org/wiki/Heapsort
	Pattern-defeating quicksort (github.com/orlp)
https://news.ycombinator.com/item?id=14661659

https://news.ycombinator.com/item?id=41066536
	My Favorite Algorithm: Linear Time Median Finding (2018) (rcoh.me)
https://danlark.org/2020/11/11/miniselect-practical-and-generic-selection-algorithms/





*/

/*
func CompactInts(arr []int) []int {
	if 1 >= len(arr) { return arr }
	// Not in place
	arrcap := cap(arr)
	for ; ; {
		var smaller, larger []int
		mid := arr[len(arr)/2]
		for ii := 0 ; ii < len(arr) ; ii++ {

		}
	}
}
*/

func Euler004() [3]int {
	answer := [3]int{0, 0, 0}
	for ii := 999; ii > 99; ii-- {
		if answer[0] > ii*ii {
			// fmt.Println("")
			break
		}
		for kk := ii; kk > 99; kk-- {
			if answer[0] > ii*kk {
				break
			}
			test := ii * kk
			if test > answer[0] && IsPalindrome(test) {
				answer = [3]int{test, ii, kk}
			}
		}
	}
	return answer
}

func Euler005(rmin, rmax int) int {
	var factors, primes []int
	primes = GetPrimes(nil, 0)
	for ii := rmin; ii <= rmax; ii++ {
		iiFact := Factor(primes, ii)
		factors = PrimeLCD(factors, iiFact)
		// fmt.Println("Euler005", ii, iiFact)
	}
	return ListMul(factors)
}

func main() {
	// Tests
	// fmt.Println("CompactInts 1..5 but reversed\t", CompactInts([]int{5, 4, 3, 2, 1}))
	// fmt.Println("CompactInts 5x1\t", CompactInts([]int{1, 1, 1, 1, 1}))
	// fmt.Println("CompactInts 1 2\t", CompactInts([]int{1, 1, 1, 1, 2}))
	// fmt.Println("CompactInts 1 2 3\t", CompactInts([]int{1, 1, 1, 3, 2}))
	fmt.Println(Euler005(1, 10) == 2520, ": 2520 is the smallest number evenly divisible by the range [1..10]")
	// Answer
	fmt.Println(Euler005(1, 20), "\n\tis the smallest number evenly divisible by the range [1..20]")
}
