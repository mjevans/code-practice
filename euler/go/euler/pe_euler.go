// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package euler

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/*

module main

require euler v1.0.0
replace euler v1.0.0 => ./euler

require bitvector v1.0.0
replace bitvector v1.0.0 => ./bitvector

go 1.19


https://go.dev/blog/package-names
https://google.github.io/styleguide/go/decisions.html
https://go.dev/ref/spec
https://pkg.go.dev/std

https://en.wikipedia.org/wiki/C_data_types#inttypes.h

https://projecteuler.net/
https://projecteuler.net/archives
https://projecteuler.net/minimal=NUM

export NUM=25 ; export FN="$(printf "pe_%04d.go" $NUM)" ; go fmt "$FN" ; go fmt euler/*.go bitvector/*.go ; go build euler/pe_euler.go ; go run "$FN"


FIXME: REMINDER -- https://go101.org/article/value-part.html
https://github.com/go101/go101/wiki/About-the-terminology-%22reference-type%22-in-Go

For greater clarity of anyone who thinks in C terms:

Value is Full Copy (single direct value part)
== single allocation, fully copied
boolean
numeric (all ints, floats etc)
pointer
unsafe.Pointer
array
struct

Value is dangerous shallow copy
== multiple allocations, management information not copied
string (len, *bytes) // However the compiler treats strings as immutable, so in practice []byte access allocates it's own copy!
slice (len, cap, *T) // DANGER: append() can allocate a larger *T (copy of one depth) also len and cap only update on the direct handle.

Value is (all pointers) Shallow Copy (indirect / reference / pointer to something within)
== Base type is pointer, or struct of only pointers.
map
channel
function
interface (specification)


https://go.dev/wiki/SliceTricks

*/

import (
	"bitvector"
	"bufio"
	"fmt"
	// "slices" // Doh not in 1.19
	"math/big"
	"sort"
	"strings"
	"sync"
	// "os" // os.Stdout
	"container/heap"
)

// 1.18+ has generics and a lot of places aren't at 1.21 yet

func maxT[T int](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func minT[T int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

/*	I
 *	I
 *	I
 *	I
 *	I
 *	I
		for ii in *\/*.go ; do go fmt "$ii" ; done ; go test -v euler/
*/

// globals
var (
	Primes *BVPrimes
)

func init() {
	Primes = NewBVPrimes()
}

/*
// Deprecated function supported by shim interface to Primes
func Factor(primes *[]int, num int) *[]int {
	// Public school factoring algorithm from memory...
	// Trial Division - https://en.wikipedia.org/wiki/Integer_factorization#Factoring_algorithms

	// With a list of known primes, the largest number that can be factored is Pn * Pn
	for ; nil == primes || num > (*primes)[len(*primes)-1]*(*primes)[len(*primes)-1]; primes = GetPrimes(primes, 0) {
		// fmt.Println(len(primes), primes[len(primes)-1])
	}

	ret := &[]int{}
	if num < 2 {
		return ret
	}
	for _, prime := range *primes {
		for ; 0 == num%prime; num /= prime {
			*ret = append(*ret, prime)
		}
		if num < prime*prime {
			break
		} // break if no more prime factors are possible
	}
	if 1 < num {
		*ret = append(*ret, num)
	}
	// fmt.Println("Factor:\t", num, "\n", ret, primes)
	return ret
}

// Deprecated function supported by shim interface to Primes
func GetPrimes(primes *[]int, primehunt int) *[]int {
	if nil == primes {
		primes = &[]int{2, 3, 5, 7, 11, 13, 17, 19}
	}
	// Semi-arbitrary expansion target, find 8 more primes (8, 16, 32, 64 it'll fit within the append growth algo)
	if primehunt < 1 {
		primehunt = 8
	}
PrimeHunt:
	for ; 0 < primehunt; primehunt-- {
		for ii := (*primes)[len(*primes)-1] + 1; ; ii++ {
			result := Factor(primes, ii)
			if 1 == len(*result) && (*primes)[len(*primes)-1] < (*result)[0] {
				//fmt.Println("Found Prime:\t", result[0])
				*primes = append(*primes, (*result)[0])
				continue PrimeHunt // I could break once, but this documents the intent
			}
		}
	}
	return primes
}
*/

func PrintFactors(factors []int) {
	// Join only takes []string s? fff
	strFact := make([]string, len(factors), len(factors))
	for ii, val := range factors {
		strFact[ii] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strFact, ", "))
}

/*
func FactorsToDivisors_old(factors *[]int) *[]int {
	fact_len := len(*factors)
	if 12 < fact_len {
		fmt.Println("FTD: ", ListMul(*factors), fact_len, "=~", Factorial(fact_len))
		return []int{}
	}
	divisors := make([]int, 0, Factorial(fact_len ))
	divisors = append(divisors, 1)
	for ii := 0; ii < fact_len; ii++ {
		mmlim := fact_len
		for mm := 0; mm < mmlim; mm++ {
			divisors = append(divisors, divisors[mm]*factors[ii])
		}
	}
	return CompactInts(divisors[:len(*divisors)-1])
}
*/

func FactorsToProperDivisors(factors *[]int) *[]int {
	fl := len(*factors)
	if 0 == fl {
		return factors
	}
	if 2 > fl {
		return &[]int{1}
	}
	if fl > 63 {
		panic("FtD does not support more than 63 factors.")
	}
	limit := (uint64(1) << fl) - 1
	bitVec := bitvector.NewBitVector(uint64(ListMul((*factors)[1:])))
	bitVec.Set(uint64(1))
	for ii := uint64(0); ii < limit; ii++ {
		div := 1
		bb := uint64(1)
		for ff := 0; ff < fl; ff++ {
			if 0 < ii&bb {
				div *= (*factors)[ff]
			}
			bb <<= 1
		}
		bitVec.Set(uint64(div))
	}
	return bitVec.GetInts()
}

func AlphaSum(str string) int64 {
	var ret, limit int64
	limit = int64(len(str))
	str = strings.ToUpper(str)
	for ii := int64(0); ii < limit; ii++ {
		ret += int64(byte(str[ii]) - 'A' + 1)
	}
	return ret
}

func ListSum(scale []int) int {
	ret := 0
	for _, val := range scale {
		ret += val
	}
	return ret
}

func ListMul(scale []int) int {
	ret := 1
	for _, val := range scale {
		ret *= val
	}
	return ret
}

func Factorial(ii int) int {
	ret := 1
	for ii > 1 {
		ret *= ii
		ii--
	}
	return ret
}

func AddInt64DecDigits(ii int64) int {
	ret := int64(0)
	for 0 < ii {
		ret += ii % 10
		ii /= 10
	}
	return int(ret)
}

/*
https://en.wikipedia.org/wiki/Fibonacci_sequence#Matrix_form
https://www.nayuki.io/page/fast-fibonacci-algorithms
"""
Given F(k) and F(k+1)

F(2k) = F(k)[2F(k+1)−F(k)]
F(2k+1) = F(k+1)^2+F(k)^2

Isolate Terms
F(k) == H
F(k+1) == J
F(k)
F(k+1)
F(k)

F(2k) = h ( 2j-h )
F(2k+1) = j^2 + h^2


*/

func BigFib(n *big.Int) (*big.Int, *big.Int) {
	zero := big.NewInt(int64(0))
	two := big.NewInt(int64(2))
	if 0 == n.Cmp(zero) {
		return big.NewInt(int64(0)), big.NewInt(int64(1))
	}
	recurse := big.NewInt(int64(0))
	recurse.Div(n, two)
	h, j := BigFib(recurse)
	// fmt.Print("BigFib rec\t", n, recurse, "\t", h, j)

	// BigFib is fed 2k : recurse with k

	// Differnt K, used to avoid X and other common variables
	k := big.NewInt(int64(0))
	// F(2k) = h ( 2j-h )
	k.Mul(j, two)
	k.Sub(k, h)
	k.Mul(k, h)
	// F(2k+1) = j^2 + h^2
	h.Mul(h, h)
	j.Mul(j, j)
	j.Add(j, h)
	// Clone N : Reuse H for modulus by two
	h.Set(n)
	h.Mod(h, two)
	// fmt.Println("\tresults: ", k, j)
	// If N was even, F(n) and F(n+1) were the returned terms.
	if 0 == h.Cmp(zero) {
		return k, j
	} else { // Calculated desired term n, but n-1...
		return j, k.Add(k, j)
	}
}

func BigFactorial(ii int64) *big.Int {
	ret := big.NewInt(int64(1))
	one := big.NewInt(int64(1))
	bi := big.NewInt(ii)
	limit := 0xFFFF
	for 0 < bi.Cmp(one) {
		ret.Mul(ret, bi)
		bi.Sub(bi, one)
		limit--
		if 0 == limit {
			panic("BigFactorial - Iter Limit Reached")
		}
	}
	return ret
}

func AddBigIntDecDigits(bi *big.Int) int64 {
	ret := int64(0)
	b := big.NewInt(ret)
	b.Set(bi)
	zero := big.NewInt(int64(0))
	ten := big.NewInt(int64(10))
	rem := big.NewInt(int64(0))
	// limit := 0x7FFF ; && limit > 0 ; limit--
	limit := 0xFFFF
	for 0 < b.Cmp(zero) {
		b.DivMod(b, ten, rem)
		ret += rem.Int64()
		limit--
		if 0 == limit {
			panic("AddBigIntDecDigits - Iter Limit Reached")
		}
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

		// // Happy Path, no dupes, scan mode
		if last < arr[k] {
			last = arr[k]
			continue
		}

		// // Eat Dupes
		// If ANY dupes, the zeros will trigger the knext / break beneath at the end when they're hit.
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
			// Zero the new gap be it from next Dup or found Next
			arr[knext] = 0
			knext++
		}
	}
	// fmt.Println(knext, arr)
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
    "BlockQuicksort" partitioning technique to mitigate branch misprediction penalties,
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

func BCDadd(in []string) string {
	accum := []int{0}
	for _, line := range in {
		line = strings.TrimSpace(line)
		carry := 0
		a := make([]int, 0, len(accum))
		for ii := 0; ii < len(accum) || ii < len(line); ii++ {
			da := 0
			if ii < len(accum) {
				da = accum[ii]
			}
			dline := 0
			if ii < len(line) {
				dline = int(line[len(line)-1-ii]) - int('0')
			}
			dsum := da + dline + carry
			carry = dsum / 10
			a = append(a, dsum%10)
		}
		if carry > 0 {
			a = append(a, 1)
		}
		accum = a
	}
	buf := make([]byte, len(accum))
	for ii := 0; ii < len(accum); ii++ {
		buf[len(buf)-1-ii] = byte(int('0') + accum[ii])
	}
	return string(buf)
}

var WrittenNumbersLow, WrittenNumbersTens []string

func InitWrittenNumbers() {
	if nil == WrittenNumbersLow {
		WrittenNumbersLow = []string{"",
			"One",
			"Two",
			"Three",
			"Four",
			"Five",
			"Six",
			"Seven",
			"Eight",
			"Nine",
			"Ten",
			"Eleven",
			"Twelve",
			"Thirteen",
			"Fourteen",
			"Fiveteen",
			"Sixteen",
			"Seventeen",
			"Eighteen",
			"Nineteen"}
	}

	if nil == WrittenNumbersTens {
		WrittenNumbersTens = []string{"",
			"",
			"Twenty",
			"Thirty",
			"Fourty",
			"Fifty",
			"Sixty",
			"Sevent",
			"Eighty",
			"Ninty"}
	}
}

func StringBritishCheckNumber(num int) (int, string) {
	InitWrittenNumbers()
	// FIXME: support more than thousands later...
	var typed int
	var ret string
	if num >= 1000 {
		ths := num / 1000
		if ths > 19 {
			panic("StringBritishCheckNumber: Fixme, number greater than 19999.")
		}
		ret += " " + WrittenNumbersLow[ths] + " Thousand"
		typed += len(WrittenNumbersLow[ths]) + len("Thousand")
		num %= 1000
	}
	if num >= 100 {
		hun := num / 100
		ret += " " + WrittenNumbersLow[hun] + " Hundred"
		typed += len(WrittenNumbersLow[hun]) + len("Hundred")
		num %= 100
		if num > 0 {
			ret += " and"
			typed += 3
		}
	}
	if num > 19 {
		tens := num / 10
		ret += " " + WrittenNumbersTens[tens]
		typed += len(WrittenNumbersTens[tens])
		num %= 10
	}
	ret += " " + WrittenNumbersLow[num]
	typed += len(WrittenNumbersLow[num])
	return typed, strings.TrimSpace(ret)
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaximumPathSum(tri [][]int) int {
	dist := make([]int, len(tri[len(tri)-1])+1)
	for line := int(len(tri)) - 1; line >= 0; line-- {
		for ii := 0; ii < len(tri[line]); ii++ {
			dist[ii] = MaxInt(tri[line][ii]+dist[ii], tri[line][ii]+dist[ii+1])
		}
	}
	return dist[0]
}

func ScannerSplitNLDQ(data []byte, atEOF bool) (advance int, token []byte, err error) {
	isJunk := func(b byte) bool {
		return '\r' == b || '\n' == b || '"' == b || ',' == b
	}
	limit := len(data)
	if 0 == limit {
		// fmt.Println("NQDL 0 limit, more data")
		return 0, nil, nil
	}
	var ii int
	for ii < limit && isJunk(data[ii]) {
		ii++
	}
	start := ii
	for ii < limit {
		if isJunk(data[ii]) {
			// fmt.Println("NQDL + ", ii, " >", string(data[0:ii]), "<")
			return ii, data[start:ii], nil
		}
		ii++
	}
	if atEOF && ii > start {
		fmt.Println("NQDL EOF + ", ii, " >", string(data), "<")
		return ii, data[start:ii], bufio.ErrFinalToken
	} else {
		// fmt.Println("NQDL no token, request more data than ", ii, " >", string(data), "<")
		return 0, nil, nil
	}
}

/*
	DoomsDayRule https://en.wikipedia.org/wiki/Doomsday_rule#Finding_a_year's_anchor_day

For the Gregorian calendar:

	YearAnchor := make(map[int]int,0,7)
	YearAnchor[1600] = 2
	YearAnchor[1700] = 0
	YearAnchor[1800] = 5
	YearAnchor[1900] = 3
	YearAnchor[2000] = 2
	YearAnchor[2100] = 0
	YearAnchor[2200] = 5

Julian dates only

	Mathematical formula
	5 × (c mod 4) mod 7 + Tuesday = anchor.
	Algorithmic
	Let r = c mod 4
	if r = 0 then anchor = Tuesday
	if r = 1 then anchor = Sunday
	if r = 2 then anchor = Friday
	if r = 3 then anchor = Wednesday

Next, find the year's anchor day. To accomplish that according to Conway:[11]

1    Divide the year's last two digits (call this y) by 12 and let a be the floor of the quotient.
2    Let b be the remainder of the same quotient.
3    Divide that remainder by 4 and let c be the floor of the quotient.
4    Let d be the sum of the three numbers (d = a + b + c). (It is again possible here to divide by seven and take the remainder. This number is equivalent, as it must be, to y plus the floor of y divided by four.)
5    Count forward the specified number of days (d or the remainder of ⁠d/7⁠) from the anchor day to get the year's one.

	( ⌊ y 12 ⌋ + y mod 1 2 + ⌊ y mod 1 2 4 ⌋ ) mod 7 + a n c h o r = D o o m s d a y {\displaystyle {\begin{matrix}\left({\left\lfloor {\frac {y}{12}}\right\rfloor +y{\bmod {1}}2+\left\lfloor {\frac {y{\bmod {1}}2}{4}}\right\rfloor }\right){\bmod {7}}+{\rm {{anchor}={\rm {Doomsday}}}}\end{matrix}}}

For the twentieth-century year 1966, for example:

	( ⌊ 66 12 ⌋ + 66 mod 1 2 + ⌊ 66 mod 1 2 4 ⌋ ) mod 7 + W e d n e s d a y = ( 5 + 6 + 1 ) mod 7 + W e d n e s d a y   = M o n d a y {\displaystyle {\begin{matrix}\left({\left\lfloor {\frac {66}{12}}\right\rfloor +66{\bmod {1}}2+\left\lfloor {\frac {66{\bmod {1}}2}{4}}\right\rfloor }\right){\bmod {7}}+{\rm {Wednesday}}&=&\left(5+6+1\right){\bmod {7}}+{\rm {Wednesday}}\\\ &=&{\rm {Monday}}\end{matrix}}}

As described in bullet 4, above, this is equivalent to:

	( 66 + ⌊ 66 4 ⌋ ) mod 7 + W e d n e s d a y = ( 66 + 16 ) mod 7 + W e d n e s d a y   = M o n d a y {\displaystyle {\begin{matrix}\left({66+\left\lfloor {\frac {66}{4}}\right\rfloor }\right){\bmod {7}}+{\rm {Wednesday}}&=&\left(66+16\right){\bmod {7}}+{\rm {Wednesday}}\\\ &=&{\rm {Monday}}\end{matrix}}}

So doomsday in 1966 fell on Monday.

Similarly, doomsday in 2005 is on a Monday:

	( ⌊ 5 12 ⌋ + 5 mod 1 2 + ⌊ 5 mod 1 2 4 ⌋ ) mod 7 + T u e s d a y = M o n d a y {\displaystyle \left({\left\lfloor {\frac {5}{12}}\right\rfloor +5{\bmod {1}}2+\left\lfloor {\frac {5{\bmod {1}}2}{4}}\right\rfloor }\right){\bmod {7}}+{\rm {{Tuesday}={\rm {Monday}}}}}

	func DoomsDayRule(year int) {
	cent := (year / 100) * 100 // lossy division
	centanchor := (5*(cent%4) + 2) % 7

	y := year % 100
	a, b := y/12, y%12
	c := b / 4
	d := a + b + c

	_ = centanchor
	_ = d
	// FIXME : This isn't worth the payoff.
}

*/

func PermutationString(perm int, str string) string {
	end := len(str)
	tmp := make([]byte, end)
	copy(tmp, str)
	res := make([]byte, end)
	slot := 0
	for slot < end {
		fact := Factorial(end - 1 - slot)
		idx := perm / fact
		perm %= fact
		res[slot] = tmp[idx]
		// fmt.Print(slot, idx, "\t", res, "\t", tmp, "\t")
		for idx < end-1-slot {
			tmp[idx] = tmp[idx+1]
			idx++
		}
		// fmt.Println(tmp)
		slot++
	}
	return string(res)
}

func RotateDecDigits(x uint64) []uint64 {
	y := x
	temp := make([]uint8, 0, 40) // a unit64 needs at most 20 digits but
	for y > 0 {
		temp = append(temp, uint8(y%10))
		y /= 10
	}
	rots := len(temp)
	temp = append(temp, temp...)
	ret := append(make([]uint64, 0, rots), x)
	// 0 1 2 3 0 1 2 3 // rots == 4
	//   1 2 3 4
	//     1 2 3 4
	//       1 2 3 4
	//         1 2 3 4
	for ii := 0; ii+1 < rots; ii++ {
		r := uint64(0)
		for d := rots; 0 < d; d-- {
			r *= 10
			r += uint64(temp[ii+d])
		}
		ret = append(ret, r)
	}
	return ret
}

func BsearchInt(list *[]int, val int) bool {
	end := len(*list)
	if nil == list || 1 > end {
		return false
	}
	left := 0
	pos := end >> 1
	end--
	for left <= pos && pos <= end {
		if (*list)[pos] == val {
			// fmt.Printf("BsearchInt: TRUE : %d\n", val)
			return true
		}
		// fmt.Printf("BsearchInt: NOW\t%d <= %d <= %d\t%d <= %d <= %d\n", left, pos, end, (*list)[left], (*list)[pos], (*list)[end])
		if (*list)[pos] < val {
			left = pos + 1
			pos += (end - pos + 1) >> 1
		} else { // gt
			end = pos - 1
			pos -= (pos - left + 1) >> 1
		}
		// fmt.Printf("BsearchInt: next\t%d <= %d <= %d\t%d\n", left, pos, end, (*list)[pos])
	}
	// fmt.Printf("BsearchInt: false : %d\n", val)
	return false
}

type Rational struct {
	Num int64
	Den int64
	Res int
	Ree int
	Imp []int8
	Quo []int8
}

func NewRational(num, den int64) *Rational {
	return &Rational{num, den, 0, 0, []int8{}, []int8{}}
}

func (ra *Rational) MulRat(rr *Rational) *Rational {
	rNum := ra.Num * rr.Num
	rDen := ra.Den * rr.Den
	rGCD := int64(GCDbin(uint(rDen), uint(rNum)))
	if 1 < rGCD {
		rNum /= rGCD
		rDen /= rGCD
	}
	return &Rational{rNum, rDen, 0, 0, []int8{}, []int8{}}
}

func (ra *Rational) Divide() {
	ra.Imp = []int8{}
	ra.Quo = []int8{}
	ra.Res = 0
	ra.Ree = 0
	if 0 == ra.Den {
		return
	}
	n := ra.Num
	d := ra.Den
	neg := false
	if n < 0 {
		n = -n
		neg = !neg
	}
	if d < 0 {
		d = -d
		neg = !neg
	}
	q := n / d
	r := n % d
	for q > 0 {
		ra.Imp = append(ra.Imp, int8(q%10))
		q /= 10
	}

	//	r*10	d	q	r	seen
	//	-	7	0.	1	-
	//	10	7	1	3	0
	//	30	7	4	2	1
	//	20	7	2	6	2
	//	60	7	8	4	3
	//	40	7	5	5	4
	//	50	7	7	1	5
	//	!! already seen

	// last remainder pos cache -- FIXME ?? Remainder can't ever be larger than ra.Den, arr possible, but... wasteful for many cases and mem clear performance / human readability.
	remCache := make(map[int64]int)
	idx := 0
	for {
		if 0 == r {
			return
		}
		if start, seen := remCache[r]; seen {
			ra.Res = start
			ra.Ree = idx
			// if 7 == d {
			// fmt.Println(remCache)
			// }
			return
		}
		ra.Quo = append(ra.Quo, int8((r*10)/d))
		remCache[r] = idx
		r = (r * 10) % d
		idx++
		if idx > 200000 {
			panic("Limit reached while in .Divide() : " + fmt.Sprint(*ra))
		}
	}

}

/**
	https://en.wikipedia.org/wiki/Integer_factorization#Factoring_algorithms
	Trial Division seems a bit easier and has the benefit of pre-sorting the result array.
	The other algorithms seem to benefit larger numbers, since I've effectively made an infinite wheel algorithm with the prime list, at the cost of memory.
	I like how simple Euler's method looks in pseudo-code, however that's a lot of squareroot operations!
	A https://stackoverflow.com/questions/2267146/what-is-the-fastest-integer-factorization-algorithm
	B https://stackoverflow.com/questions/1877255/problems-with-prime-numbers
	<(short)	Lookup Table
	<2^70		Richard Brent's mod of Pollard's Rho algo http://wwwmaths.anu.edu.au/~brent/pub/pub051.html
	<10^50		Lenstra Elliptic curve http://en.wikipedia.org/wiki/Lenstra_elliptic_curve_factorization
	<10^100		Quadratic Sieve http://en.wikipedia.org/wiki/Quadratic_sieve
	>10^100		GNFS http://en.wikipedia.org/wiki/General_number_field_sieve

	Offhand, from a pragmatic viewpoint, a list of primes betten 0 and the largest under 65536 is _probably_ more memory than a practical program should use, though 0..255 is clearly too limited.
	[]uint16 might be a good format for the primes list, if not a bitvector directly.

	2..7919 contains 1000 prime numbers; stored as a compressed (inherently 2 is prime so 3..7919) bitvector, that would take 3958 bits or 495 bytes (rounded up)
	It's entirely practical to throw a 512 or 4096 byte bucket of primes at the issue and simplify life.
	Page ~= 64 bytes
	3..130 = Page 0 The highest prime is 127 which has a square root of ~11.27 (121<>144)
	3..18 == BYTE 0 17,13,11,7,5,3
**/

// BVpagesize >= BVl1 // Both MUST be a power of 2 ( Pow(2, n) )
// WARNING: Populate more known primes if increasing BVl1 size
const BVl1 = 64
const BVpagesize = 4096
const BVbitsPerByte = 8
const BVprimeByteBitMaskPost = BVbitsPerByte - 1
const BVprimeByteBitMask = BVprimeByteBitMaskPost << 1 // 0b_1110 // The 2^0 = 1s bit is discarded in the compression shift
// const BVprimeByteBitMask = 0b_1110 // The 2^0 = 1s bit is discarded in the compression shift
// const BVprimeByteBitMaskPost = BVprimeByteBitMask >> 1
const BVprimeByteBitShift = 3 + 1 // 3 bits for 8 bit index, plus 1 bit for discard all even numbers

type BVpage [BVpagesize]uint8

type BVPrimes struct {
	Last uint
	Mu   sync.Mutex
	PV   []*BVpage // starting from bit 0 (set) == 3 (prime), record all odd primes with SET bits
	// MAYBE primes are any unset bits > Last, unset bits < Last == composite
}

func NewBVPrimes() *BVPrimes {
	ov := new(BVpage)
	// no even 7_31 _753
	ov[0] = 0b_0100_1000
	//          19   3 9 // 33 is not prime, but it is the last tested number 33*33 = 1089 the first l1 cacheline is safe to factor in place.
	ov[1] = 0b_1001_1010
	// WARNING: Populate more known primes if increasing BVl1 size
	return &BVPrimes{PV: append(make([]*BVpage, 0, 1), ov), Last: 33}
}

func (p *BVPrimes) PrimeOrDown(ii uint) uint {
	if 2 > ii {
		return 0
	}
	if 2 == ii {
		return 2
	}
	// NOTE: Strictly the 'or' part is in use here, find the last number _known_ to be a prime
	if ii > p.Last {
		ii = p.Last
	}
	// in := ii
	ii = (ii - 3)
	bidx := (ii & BVprimeByteBitMask) >> 1
	ii >>= BVprimeByteBitShift
	pg, pidx := ii/BVpagesize, ii%BVpagesize
	// fmt.Printf("PrimeOrDown from [%d][%d]&%x == %d\n", pg, pidx, (uint8(1) << bidx), in)
	// pg
	for {
		// pidx
		for {
			// bidx
			for {
				if 0 == p.PV[pg][pidx]&(uint8(1)<<bidx) {
					return ((pg*BVpagesize + pidx) << BVprimeByteBitShift) + uint(bidx)<<1 + 3
				}
				if 0 == bidx {
					break
				}
				bidx--
			}
			bidx = BVbitsPerByte - 1 // reset after scanning the initial index bits
			if 0 == pidx {
				break
			}
			pidx--
		}
		if 0 == pg && 0 == pidx {
			fmt.Println("Hit the floor fallback = 3; PrimeOrDown()")
			return 3
		}
		pg--
		pidx = BVpagesize - 1
	}
}

func (p *BVPrimes) PrimeAfter(ii uint) uint {
	if 2 > ii {
		return 2
	}
	// Guard an underflow, 2 doesn't really exist to step after
	if 2 == ii {
		return 3
	}
	lastPrime := p.PrimeOrDown(p.Last)
	if ii >= lastPrime {
		newLimit := (((((ii-3)>>1)/uint(BVl1))+uint(1))*uint(BVl1)+uint(BVprimeByteBitMaskPost))<<1 + 3
		// fmt.Printf("Primes.PAfter .Grow triggered:   \t%d\t< %d\t-> %d\n", ii, p.Last, newLimit)
		p.Grow(newLimit)
	}
	// } else {
	// fmt.Printf("Prime.PAfter last prime:\t%d\t< %d\n", lastPrime, p.Last)
	// if 7600 < ii {
	// fmt.Printf("Prime.PAfter last prime:\t%d\t< %d\n", lastPrime, p.Last)
	// }
	// }
	return p.primeAfterUnsafe(ii, p.Last)
}

func (p *BVPrimes) primeAfterUnsafe(input, limit uint) uint {
	ii := (input - 3 + 2) // the prime number AFTER ii, E.G. 6 -> 7
	bidx := (ii & BVprimeByteBitMask) >> 1
	bidx0 := bidx
	ii >>= BVprimeByteBitShift
	pg, pidx := ii/BVpagesize, ii%BVpagesize
	bbMM := ((limit - 3) & BVprimeByteBitMask) >> 1
	ooMM := (limit - 3) >> BVprimeByteBitShift
	pgMM, idxMM := ooMM/BVpagesize, ooMM%BVpagesize
	pgmax, pimax, pbmax := uint(len(p.PV)), uint(BVpagesize-1), uint(BVprimeByteBitMaskPost)
	if pgmax < pgMM {
		pgMM, idxMM, bbMM = pgmax, pimax, pbmax
	}
	_ = bbMM // scan the whole byte to simplify the logic check

	// pg
	for pg <= pgMM {
		// pidx
		for (pg < pgMM && pidx <= pimax) || (pg == pgMM && pidx <= idxMM) {
			// bidx
			for ; bidx < BVbitsPerByte; bidx++ {
				if 0 == p.PV[pg][pidx]&(uint8(1)<<bidx) {
					return ((pg*BVpagesize + pidx) << BVprimeByteBitShift) + uint(bidx)<<1 + 3
				}
			}
			bidx = 0 // reset after scanning the initial index bits
			pidx++
		}
		if pidx >= BVpagesize {
			pg++
			pidx = 0
		}
	}
	fmt.Printf("Unable to locate prime after %d under %d\t[%d/%d][%d/%d]\t", input, limit, pg, pgMM, pidx, idxMM)
	pg, pidx = ii/BVpagesize, ii%BVpagesize
	fmt.Printf("started near [%d][%d]:%d (%d)\n", pg, pidx, bidx0, input)
	var cl, end uint
	cl = ii / BVl1
	for end < limit {
		cl++
		end = (((cl) * BVl1) << 4) + 1
		ccount := len(p.PrimesOnPage(end))
		fmt.Printf("\t{%d, %d},", end, ccount)
		if 512 == ccount {
			fmt.Printf("<<<ERROR\t")
			break
		}
		if 0 == cl&7 {
			fmt.Println()
		}
	}
	fmt.Println()
	return 0
}

func (p *BVPrimes) wheelFactCL1Unsafe(start, prime, maxPrime uint) (uint, uint) {
	// https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes
	// https://en.wikipedia.org/wiki/Sieve_of_Sundaram
	// https://en.wikipedia.org/wiki/List_of_prime_numbers

	//	Last	Next	Pri	Correct	Diff	FlMod(2p)+p
	//	33	35	3	39	4	33
	//	33	35	5	35	0	35
	//	33	35	7	35	0	35
	//	33	35	11	55	20	33
	//	33	35	13	39	4
	//	33	35	17	51	16
	//	33	35	19	57	22
	//	33	35	23	69	34
	//	33	35	29	87	52
	//	33	35	31	93	58

	if 0 == start {
		// Modern CPUs prefer a branchless path even with a couple possibly redundant operations
		// 2 + (IF even (step back to last odd) ELSE odd noop )
		start = uint(((p.Last - 1) | 1) + 2)
	}
	start |= 1 // Evens inherently compressed out; always odd
	if 3 > prime {
		prime = 3
	}
	// bits,	 octets (bytes),	 and page / pgLine (line on page)
	ABS_bbStart := (start - 3) >> 1
	ABS_ooStart := ABS_bbStart >> (BVprimeByteBitShift - 1)
	pg, pgLine := ABS_ooStart/BVpagesize, (ABS_ooStart%BVpagesize)/BVl1
	bbLimit := ((pgLine*BVl1 + BVl1 - 1) << (BVprimeByteBitShift - 1)) | BVprimeByteBitMaskPost
	ABS_maxPrimeReq := 3 + (bbLimit << 1) + (pg * BVpagesize * BVbitsPerByte << 1)
	if 0 == maxPrime {
		// There's... probably a prime on this page?  For page 0 this returns 3 which at least terminates early if uselessly, for any other page it should be sufficient.
		maxPrime = 3 + (ABS_ooStart << BVprimeByteBitShift)
	}
	bbStartPg := ABS_bbStart % (BVpagesize * BVbitsPerByte)
	//if 0 < pg {
	//	fmt.Printf("TRACE: [%d][%d]\tprime = %d\tstart=%d\t%d\n", pg, bbStartPg, ABS_maxPrimeReq, start, bbLimit)
	//}
	for 0 != prime && prime <= maxPrime && prime<<1 < ABS_maxPrimeReq {
		// calculate the next modulus to mark, this will always be an odd multiple of prime of at least 3 * prime...
		startModPr := (((start-prime-1)/(prime<<1))+1)*(prime<<1) + prime
		if startModPr <= ABS_maxPrimeReq {
			// pgC, bbPos := ((startModPr-3)>>1)/(BVpagesize*BVbitsPerByte), ((startModPr-3)>>1)%(BVpagesize*BVbitsPerByte)
			bbPos := ((startModPr - 3) >> 1) % (BVpagesize * BVbitsPerByte)
			if bbPos < bbStartPg {
				fmt.Printf("TRACE: %d\tstart=%d\tprime = %d\tsModPr=%d\t[%d][%d]\n", bbPos, start, prime, startModPr, pg, pgLine)
				fmt.Printf("Logic Error, %d < %d\n", bbPos, bbStartPg)
				panic("debug")
			}
			for bbPos <= bbLimit {
				// Spot check early problems : 457 509
				// if 0 == pg && (227 == bbPos || 253 == bbPos) {
				// 	fmt.Printf("TRACE: %d\tprime = %d\tstart=%d\n", bbPos, prime, start)
				// }
				p.PV[pg][bbPos>>(BVprimeByteBitShift-1)] |= uint8(1) << (bbPos & BVprimeByteBitMaskPost)
				bbPos += prime // compressed (/2) prime bitvector, this advances by the prime * 2
			}
		}
		// else { // Minimum iteration outside this window }
		prime = p.primeAfterUnsafe(prime, prime<<1)
		if 0 == prime {
			fmt.Printf("TRACE: %d .. %d\t%d\t[%d~%d][...]\n", start, ABS_maxPrimeReq, startModPr, pg, pgLine)
			panic("primeAfterUnsafe Returned Zero")
		}
	}
	return prime, ABS_maxPrimeReq
}

func (p *BVPrimes) autoFactorPMCforBVl1(start uint) {
	if 0 == start {
		// 2 + (IF even (step back to last odd) ELSE odd noop )
		start = uint(((p.Last - 1) | 1) + 2)
	}
	ABS_bbStart := (start - 3) >> 1
	ABS_ooStart := ABS_bbStart >> (BVprimeByteBitShift - 1)
	pg, pgLine := ABS_ooStart/BVpagesize, (ABS_ooStart%BVpagesize)/BVl1
	bbLimit := ((pgLine*BVl1 + BVl1 - 1) << (BVprimeByteBitShift - 1)) | BVprimeByteBitMaskPost
	bbPos := ABS_bbStart % (BVpagesize * BVbitsPerByte)
	for bbPos <= bbLimit {
		if 0 == p.PV[pg][bbPos>>(BVprimeByteBitShift-1)]&(uint8(1)<<(bbPos&BVprimeByteBitMaskPost)) {
			posPrime := (pg * BVpagesize << BVprimeByteBitShift) + (bbPos << 1) + 3
			posRes := Factor1980AutoPMC(posPrime, false)
			if posPrime != posRes {
				// fmt.Printf(" %d", posPrime)
				p.PV[pg][bbPos>>(BVprimeByteBitShift-1)] |= uint8(1) << (bbPos & BVprimeByteBitMaskPost)
			}
		}
		bbPos++
	}
}

func (p *BVPrimes) PrimesOnPage(start uint) []uint {
	if 0 == start {
		// 2 + (IF even (step back to last odd) ELSE odd noop )
		start = uint(((p.Last - 1) | 1) + 2)
	}
	ret := make([]uint, 0, 64)
	bbStart := (start - 3) >> 1
	ooStart := bbStart >> (BVprimeByteBitShift - 1)
	pg, pgLine := ooStart/BVpagesize, (ooStart%BVpagesize)/BVl1
	bbLimit := ((pgLine*BVl1 + BVl1 - 1) << (BVprimeByteBitShift - 1)) | BVprimeByteBitMaskPost

	bbPos := (pgLine * BVl1) << (BVprimeByteBitShift - 1)
	for bbPos <= bbLimit {
		if 0 == p.PV[pg][bbPos>>(BVprimeByteBitShift-1)]&(uint8(1)<<(bbPos&BVprimeByteBitMaskPost)) {
			ret = append(ret, (pg*BVpagesize<<BVprimeByteBitShift)+(bbPos<<1)+3)
		}
		bbPos++
	}
	return ret
}

func (p *BVPrimes) Grow(limit uint) {
	if 0x10000f00 < limit {
		fmt.Printf("Emperically refusing to grow past ~2sec runtime (~2015 era Xeon 1 CPU core) %d < %d", 0x100000, limit)
		panic("Likely overflow")
		// https://en.wikipedia.org/wiki/Primality_test#Number-theoretic_methods
	}

	if p.Last >= limit {
		// fmt.Printf("Already above requested growth limit, %d, at %d\n", limit, p.Last)
		return
	}

	// Attempt to lock, if another goroutine (thread) is updating p.Last has probably changed
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if p.Last >= limit {
		return
	}

	// last l1 cache line
	cl1z := (((limit - 3) >> BVprimeByteBitShift) / uint(BVl1)) + uint(1)

	// Ensure the bitvector arrays exist
	pagez := cl1z/(BVpagesize/BVl1) + 1
	lenpv := uint(len(p.PV))
	if pagez > lenpv {
		// Extend Capacity https://go.dev/wiki/SliceTricks
		p.PV = append(make([]*BVpage, 0, pagez), p.PV...)
		for lenpv <= pagez {
			p.PV = append(p.PV, new(BVpage))
			lenpv++
		}
	}
	// for ii := uint(0); ii < pagez; ii++ {
	//	fmt.Printf("Pointer check Primes Page %4d = %p\n", ii, p.PV[ii])
	//}

	next := ((p.Last - 1) | 1) + 2
	line := ((next - 3) >> BVprimeByteBitShift) / uint(BVl1)

	// ??? FIXME ???
	// This might be seen as a refined and optimized version of 'first gear'; extending the concepts of trial division, wheel, and sieves.
	// As it marks prime repeats (the odd multiples), the repeat of the 'next' prime must extend past the current page of the array.
	// (duh) More primes must be tested the deeper numbers progress.
	// FIXME: It's well past the scope of this library (educational / toy work) to quantify the computational growth pattern, or even an approximation, of how quickly the cost grows...
	// Though it does seem clear that it's well under the number 1,000,000
	// At some point it must make sense to switch to a different gear
	//

	// Gear 1 : wheelFactCL1Unsafe Faster than autoFactorPMC through about 64K
	// 0x100000 ~= 0.74s
	// 0x200000 ~= 2.82s
	// 0x400000 ~= 10.75s // Tried 2048 for the cutoff with this and autoFactorPMCforBVl1 was _still_ slower on a decade old Xeon CPU, it's correct, but the 'need to be 100%sure' spin cycle makes it too slow.  Perceptibly >> 201 seconds (I killed the run) vs 10.78s with the extra logic test.
	for line <= cl1z && line < 20480 {
		primeStart := uint(3)
		var end uint
		for {
			maxPrimeCall := p.PrimeOrDown(p.Last)
			primeStart, end = p.wheelFactCL1Unsafe(next, primeStart, maxPrimeCall)

			if primeStart<<1 > end {
				break
			}
			next = (primeStart << 1) | 1
			p.Last = next
			// fmt.Printf("%d:\t%d\t@%d\tmaybe %d\n", line, primeStart, next, p.countPrimesLEUnsafe(end))
		}
		// Last = 3 + (((line*BVl1 + BVl1 - 1) << (BVprimeByteBitShift)) | BVprimeByteBitMask)
		next = 3 + (((line + 1) * BVl1) << BVprimeByteBitShift)
		p.Last = next - 2
		ccount := len(p.PrimesOnPage(p.Last))
		if 512 == ccount {
			fmt.Printf("%d: %d\t\t%d\tmaybe %d == %d\t\t%d primes on page\n", line, p.Last, primeStart, p.countPrimesLEUnsafe(p.Last), end, ccount)
			panic("too many primes")
			// This _reliably_ fails on page !0, why?
		}
		if 256 < ccount {
			fmt.Printf("SUS: about 80%% of the numbers should be filtered as a minimum...\n%d: %d\t\t%d\tmaybe %d == %d\t\t%d primes on page\n", line, p.Last, primeStart, p.countPrimesLEUnsafe(p.Last), end, ccount)
		}
		line++
	}

	// Gear 2: _partial_ filter pages for primes up to a reasonable cost... Likely less than one cache line's worth of bitfield; probably way less if the growth rate is any indication.
	//         Then test each not known-composite number on the current page with another algo, E.G. Factor1980PollardMonteCarlo
	//
	// Gear 2 might be worth it if there's a major need to know all the primes under a given value, rather than just factoring. This seive cost VS GCDbin VS Pollard?

	// AFTER 64K this is somehow so fast that I suspect the results... FIXME: Have I added a total torture test to cover up to 1024*1024 yet?

	for line <= cl1z {
		p.wheelFactCL1Unsafe(next, 3, 509) // 5 = 498062; 503 = 498062
		p.autoFactorPMCforBVl1(next)
		next = 3 + (((line + 1) * BVl1) << BVprimeByteBitShift)
		p.Last = next - 2
		// ccount := len(p.PrimesOnPage(p.Last))
		fmt.Printf("%d: %d\n", line, p.Last)
		line++
	}

}

func (p *BVPrimes) MaybePrime(q uint) bool {
	// Use base2 storage inherent test for division by 2
	if 0 == q&0x01 && 2 < q {
		return false
	}
	if q > p.Last {
		return true
	}
	pd := p.PrimeOrDown(q)
	return pd == q
}

func (p *BVPrimes) KnownPrime(q uint) bool {
	// Use base2 storage inherent test for division by 2
	if (0 == q&0x01 && 2 < q) || q > p.Last {
		return false
	}
	pd := p.PrimeOrDown(q)
	return pd == q
}

func (p *BVPrimes) GetPrimesInt(primes *[]int, num int) *[]int {
	if nil == primes {
		primes = &[]int{}
		//*primes = make([]int, 0, 8+num)
		*primes = append(make([]int, 0, 8+num), 2)
	} else {
		*primes = append(make([]int, 0, 8+num), (*primes)...)
	}

	ii := len(*primes)
	lim := cap(*primes)
	prime := uint((*primes)[ii-1])
	// fmt.Printf("GetPrimesInt: cap %d\n", lim)
	for ; ii < lim; ii++ {
		prime = p.PrimeAfter(prime)
		*primes = append(*primes, int(prime))
	}
	// fmt.Printf("GetPrimesInt: %v\n", *primes)
	return primes
}

func (p *BVPrimes) CountPrimesLE(ii uint) uint {
	if 2 > ii {
		return 0
	}
	if 2 == ii {
		return 1
	}
	lastPrime := p.PrimeOrDown(p.Last)
	if ii >= lastPrime {
		newLimit := (((((ii-3)>>1)/uint(BVl1))+uint(1))*uint(BVl1)+uint(BVprimeByteBitMaskPost))<<1 + 3
		// fmt.Printf("Primes.PAfter .Grow triggered:   \t%d\t< %d\t-> %d\n", ii, p.Last, newLimit)
		p.Grow(newLimit)
	}
	return p.countPrimesLEUnsafe(ii)
}
func (p *BVPrimes) countPrimesLEUnsafe(ii uint) uint {
	// 2 isn't in the list, it's implied
	ret := uint(1)
	// out := 0
	// fmt.Printf("\nfactor ")

	ii = (ii - 3)
	bidx := (ii & BVprimeByteBitMask) >> 1
	ii >>= BVprimeByteBitShift
	pg, pidx := ii/BVpagesize, ii%BVpagesize
	// fmt.Printf("PrimeOrDown from [%d][%d]&%x == %d\n", pg, pidx, (uint8(1) << bidx), in)
	// pg
	for {
		// pidx
		for {
			// bidx
			for {
				if 0 == p.PV[pg][pidx]&(uint8(1)<<bidx) {
					ret++
					// fmt.Printf(" %d", ((pg*BVpagesize+pidx)<<BVprimeByteBitShift)+uint(bidx)<<1+3)
					// out++
					// if 20 < out && 0x1000 < ii { panic("debug") }
				}
				if 0 == bidx {
					break
				}
				bidx--
			}
			bidx = BVbitsPerByte - 1 // reset after scanning the initial index bits
			if 0 == pidx {
				break
			}
			pidx--
		}
		if 0 == pg && 0 == pidx {
			return ret
		}
		pg--
		pidx = BVpagesize - 1
	}
}

func GCDbin(a, b uint) uint {
	// https://en.wikipedia.org/wiki/Binary_GCD_algorithm
	// https://pkg.go.dev/math/bits#TrailingZeros
	// https://cs.opensource.google/go/go/+/go1.23.1:src/math/bits/bits.go;l=59
	// if x == 0 { return 64 } ; return int(deBruijn64tab[(x&-x)*deBruijn64>>(64-6)])
	// See http://supertech.csail.mit.edu/papers/debruijn.pdf
	// const deBruijn32 = 0x077CB531
	// var deBruijn32tab = [32]byte{		0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,		31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,	}
	// const deBruijn64 = 0x03f79d71b4ca8b09
	// var deBruijn64tab = [64]byte{		0, 1, 56, 2, 57, 49, 28, 3, 61, 58, 42, 50, 38, 29, 17, 4,		62, 47, 59, 36, 45, 43, 51, 22, 53, 39, 33, 30, 24, 18, 12, 5,		63, 55, 48, 27, 60, 41, 37, 16, 46, 35, 44, 21, 52, 32, 23, 11,		54, 26, 40, 15, 34, 20, 31, 10, 25, 14, 19, 9, 13, 8, 7, 6,	}

	// There is a better way to count trailing zeros, but it uses annoying magic numbers or imports "math/bits"

	a0, b0 := a, b
	//fmt.Printf("GCDbin %d, %d\n", a0, b0)
	//if a0 > 0xffffff || b0 > 0xffffff {
	//	panic("overflow") }
	if 0 == b {
		return a
	}
	if 0 == a {
		return b
	}

	var ka, kb int
	k := 0 // k == count of common 2 factors
	for 0 == a&1 {
		a >>= 1
		ka++
	}
	for 0 == b&1 {
		b >>= 1
		kb++
	}
	if ka > kb {
		k = kb
	} else {
		k = ka
	}

	for {
		if a < b {
			a, b = b, a
		}
		a -= b // a is now even
		if a == 0 {
			return b << k
		}
		// b is odd, therefore no more common 2 factors, discard
		g := 0
		for 0 == a&1 {
			a >>= 1
			g++
			if g > 63 {
				fmt.Printf("GCDbin failed to converge with %d, %d\n", a0, b0)
				panic("overflow")
			}
		}
	}
}

/*
			fmt.Printf("
Pollard: %d ??\t (%d >= %d) || (1 < %d)\t%d\t%d\n",

N,			k, r, 		G, 	x, 	y	)

Pollard: 5309 ??         (1 >= 1) || (1 < 1)    0       147
Pollard: 5309 ??         (1 >= 2) || (1 < 1)    147     4833
Pollard: 5309 ??         (2 >= 2) || (1 < 1)    147     2626
Pollard: 5309 ??         (1 >= 4) || (1 < 1)    2626    5021
Pollard: 5309 ??         (2 >= 4) || (1 < 1)    2626    953
Pollard: 5309 ??         (3 >= 4) || (1 < 1)    2626    1098
Pollard: 5309 ??         (4 >= 4) || (1 < 1)    2626    2939
Pollard: 5309 ??         (1 >= 8) || (1 < 1)    2939    4130
Pollard: 5309 ??         (2 >= 8) || (1 < 1)    2939    1886
Pollard: 5309 ??         (3 >= 8) || (1 < 1)    2939    964
Pollard: 5309 ??         (4 >= 8) || (1 < 1)    2939    2398
Pollard: 5309 ??         (5 >= 8) || (1 < 1)    2939    4231
Pollard: 5309 ??         (6 >= 8) || (1 < 1)    2939    1283
Pollard: 5309 ??         (7 >= 8) || (1 < 1)    2939    954
Pollard: 5309 ??         (8 >= 8) || (1 < 1)    2939    892
Pollard: 5309 ??         (1 >= 16) || (1 < 1)   892     60
Pollard: 5309 ??         (2 >= 16) || (1 < 1)   892     1107
Pollard: 5309 ??         (3 >= 16) || (1 < 1)   892     4583
Pollard: 5309 ??         (4 >= 16) || (1 < 1)   892     294
Pollard: 5309 ??         (5 >= 16) || (1 < 1)   892     5248
Pollard: 5309 ??         (6 >= 16) || (1 < 1)   892     1071
Pollard: 5309 ??         (7 >= 16) || (1 < 1)   892     5059
Pollard: 5309 ??         (8 >= 16) || (1 < 1)   892     2671
Pollard: 5309 ??         (9 >= 16) || (1 < 1)   892     2435
Pollard: 5309 ??         (10 >= 16) || (1 < 1)  892     879
Pollard: 5309 ??         (11 >= 16) || (1 < 1)  892     862
Pollard: 5309 ??         (12 >= 16) || (1 < 1)  892     2900
Pollard: 5309 ??         (13 >= 16) || (1 < 1)  892     1908
Pollard: 5309 ??         (14 >= 16) || (1 < 1)  892     4109
Pollard: 5309 ??         (15 >= 16) || (1 < 1)  892     4999
Pollard: 5309 ??         (16 >= 16) || (1 < 1)  892     689
Pollard: 5309 ??         (1 >= 32) || (1 < 1)   689     1908
Pollard: 5309 ??         (2 >= 32) || (1 < 1)   689     4109
Pollard: 5309 ??         (3 >= 32) || (1 < 1)   689     4999
Pollard: 5309 ??         (4 >= 32) || (1 < 5309)        689     689

*/

// Returns _a_ factor OR 0 on failure (means MAYBE prime, 289 fails) (0 or 1 are never returned on success) 'took x0 := 0 ; m = 1'
func Factor1980PollardMonteCarlo(N, x0 uint) uint {
	// https://en.wikipedia.org/wiki/Pollard%27s_rho_algorithm

	// https://maths-people.anu.edu.au/~brent/pub/pub051.html
	// https://maths-people.anu.edu.au/~brent/pd/rpb051i.pdf
	//if 0 == x0 { x0 = 2 }

	// if x0 > 0xffffff {
	//	fmt.Printf("x0 unlikely large: %d\n", x0)
	//	panic("unexpected value")
	//}

	// print pg182-183 p"2
	fx := func(x, Nin uint) uint {
		return (x*x + 3) % Nin
	}
	umin := func(a, b uint) uint {
		if a < b {
			return a
		}
		return b
	}
	abssub := func(a, b uint) uint {
		if a < b {
			return b - a
		}
		return a - b
	}

	// they 'took x0 := 0 ; m = 1' -- pg183 8.
	//								// Pascal ? https://www.freepascal.org/docs-html/ref/refsu60.html
	//								y := x0 ; r := 1 ; q := 1 ;
	//								repeat x := y ;
	//									for i := 1 to r do y := f(y) ; k := 0 ;
	//									repeat ys := y ;
	//										for i := 1 to min(m, r-k) do
	//											begin y := f(y) ; q := q * abs(x-y) mod N
	//											end;
	//										G := GCD(q,N) ; k := k + m
	//									until (k >= r) or (G > 1) ; r := 2 * r
	//								until G > 1 ;
	//								if G == N then
	//									repeat ys := f(ys) ; G := GCD(abs(x - ys), N)
	//									until G > 1 ;
	//								if G == N then {failure} else {success}

	// x := r Round? / Rotate / Roll? (bitshift / pass) (fast hare)
	// y := k seems to be the slow tortoise
	//
	//								y := x0 ; r := 1 ; q := 1 ;
	y, r, q, m := uint(x0), uint(1), uint(1), uint(1)
	var k, ii, G, x, ys uint
	//								repeat x := y ;
	for {
		x = y
		//								for i := 1 to r do y := f(y) ; k := 0 ;
		for ii = 1; ii <= r; ii++ {
			y = fx(y, N) // tortoise
		}
		k = 0
		//								repeat ys := y ;
		for {
			ys = y // save Y for GCD extraction?
			//								for i := 1 to min(m, r-k) do
			iz := umin(m, r-k)
			for ii = 0; ii <= iz; ii++ {
				//								begin y := f(y) ; q := q * abs(x-y) mod N
				y = fx(y, N)
				q = (q * abssub(x, y)) % N
			}
			//									end;
			//								G := GCD(q,N) ; k := k + m
			G = GCDbin(q, N)
			k += m
			//							until (k >= r) or (G > 1) ; r := 2 * r
			// fmt.Printf("Pollard: %d ??\t (%d >= %d) || (1 < %d)\t%d\t%d\n", N, k, r, G, x, y)
			if (k >= r) || (1 < G) {
				break
			}
		}
		r <<= 1
		//							until G > 1 ;
		if 1 < G {
			break
		}
	}
	//								if G == N then
	if G == N {
		//								repeat ys := f(ys) ; G := GCD(abs(x - ys), N)
		for {
			ys = fx(ys, N)
			G = GCDbin(abssub(x, ys), N)
			//							until G > 1 ;
			if 1 < G {
				break
			}
		}
	}
	//								if G == N then {failure} else {success}
	if G == N {
		// fmt.Printf("Pollard: FAILED %d / %d\n", N, G)
		return 0
	} // 0 == Failed
	// fmt.Printf("Pollard: %d / %d\n", N, G)
	return G
}

func Factor1980AutoPMC(q uint, singlePrimeOnly bool) uint {
	if 0 == q&1 {
		return 2
	}
	// This test appears to yield a ~200x improvement in factor speed 20.66s vs 0.01s for 2..65535
	// Primes is a global instance which hopefully knows if small Q are prime or not...
	if Primes.PrimeOrDown(q) == q {
		return q
	}

	unk := Factor1980PollardMonteCarlo(q, 0)
	if 0 != unk {
		if singlePrimeOnly {
			return Factor1980AutoPMC(unk, singlePrimeOnly)
		}
		return unk
	}
	// Usually works in one pass, but if not...

	// __approximate__ an integer square root, POW(pollard_limit, 2) _MUST_ be > q == pl*pl means square root factor
	pollard := uint(1)
	pollard_limit := pollard
	for t := q >> 1; pollard_limit < t; t >>= 1 {
		pollard_limit <<= 1
	}
	for dist := pollard_limit >> 1; 0 < dist; dist >>= 1 {
		t := pollard_limit + dist
		if t*t < q {
			pollard_limit += dist
		}
	}
	// Only found 2 which deserves a guard at the top
	// if q == pollard_limit*pollard_limit {
	//	fmt.Printf("Found squared pre at %d (%d)\n", pollard_limit, q)
	//	return pollard_limit }
	pollard_limit++
	// 17 (289) 79 (6241) 139 (19321) 181 (32761)
	// Test for Sqaure Root factor
	if q == pollard_limit*pollard_limit {
		// fmt.Printf("Found squared post at %d (%d)\n", pollard_limit, q)
		if singlePrimeOnly {
			return Factor1980AutoPMC(pollard_limit, singlePrimeOnly)
		}
		return pollard_limit
	}
	// fmt.Printf("pollard_limit(%d) => %d\n", q, pollard_limit)

	for pollard <= pollard_limit {
		unk := Factor1980PollardMonteCarlo(q, pollard)
		if 0 != unk {
			if singlePrimeOnly {
				return Factor1980AutoPMC(unk, singlePrimeOnly)
			}
			return unk
		}
		pollard++
	}
	return q
}

func (p *BVPrimes) Factorize(q uint) *Factorized {
	qin := q
	_ = qin
	// Low hanging fruit first
	if 2 > q {
		if 0 == q {
			return &Factorized{Lenbase: 1, Lenpow: 1, Fact: []Factorpair{Factorpair{Base: 0, Power: 1}}}
		}
		return &Factorized{Lenbase: 1, Lenpow: 1, Fact: []Factorpair{Factorpair{Base: 1, Power: 1}}}
	}

	facts := &FactorpairQueue{}
	heap.Init(facts)

	// Special test & extract: base2, /2
	k := 0
	for 0 == q&1 {
		q >>= 1
		k++
	}
	if 0 < k {
		heap.Push(facts, Factorpair{Base: 2, Power: uint16(k)})
	}

	// pLim := uint(7)
	// p.Grow(pLim)
	// for cur := uint(3); 1 < q && cur <= pLim; cur = p.primeAfterUnsafe(cur, pLim) {

	// Quickly test some small primes; 2, 3 (~66%), 5 (~73%), 7 (<77%) -- https://en.wikipedia.org/wiki/Wheel_factorization#Description
	smallPrimes := []uint16{3, 5, 7}
	for cur := 0; 1 < q && cur < len(smallPrimes); cur++ {
		fac := Factorpair{Base: smallPrimes[cur], Power: uint16(0)}
		qd := uint(smallPrimes[cur])
		for 0 == q%qd {
			q /= qd
			fac.Power++
		}
		if 0 < fac.Power {
			heap.Push(facts, fac)
		}
	}
	// zz := 1000
	// for 1 < q && zz > 0 {
	for 1 < q {
		unk := Factor1980AutoPMC(q, false)
		// Probably Prime
		if unk == q {
			heap.Push(facts, Factorpair{Base: uint16(q), Power: 1})
			break
		}
		q /= unk
		sf := p.Factorize(unk)
		for ii := uint(0); ii < sf.Lenbase; ii++ {
			b := uint(sf.Fact[ii].Base)
			pow := uint16(0)
			for 0 == q%b {
				q /= b
				pow++
			}
			if 0 < pow {
				sf.Fact[ii].Power += pow
			}
			heap.Push(facts, sf.Fact[ii])
		}
		// fmt.Printf("Factorize: %d @ %d\n", qin, q)
		// zz--
	}

	var base, power uint
	fact := make([]Factorpair, 0, facts.Len())
	for 0 < facts.Len() {
		fp := heap.Pop(facts).(Factorpair)
		base++
		power += uint(fp.Power)
		fact = append(fact, fp)
	}
	return &Factorized{Lenbase: base, Lenpow: power, Fact: fact}
}

// Deprecated function supported by shim interface to Primes
func GetPrimes(primes *[]int, num int) *[]int {
	return Primes.GetPrimesInt(primes, num)
}

// Deprecated function supported by shim interface to Primes
func Factor(primes *[]int, num int) *[]int {
	fp := Primes.Factorize(uint(num))
	ret := make([]int, 0, fp.Lenpow)
	iiLim := len(fp.Fact)
	if int(fp.Lenbase) != iiLim {
		fmt.Printf("WARNING: malformed factor pair returned by Primes.Factorize(), %d != %d", fp.Lenbase, iiLim)
	}
	for ii := 0; ii < iiLim; ii++ {
		for kk := 0; kk < int(fp.Fact[ii].Power); kk++ {
			ret = append(ret, int(fp.Fact[ii].Base))
		}
	}
	return &ret
}

type Factorpair struct {
	Base  uint16
	Power uint16
}

type Factorized struct {
	// Euler 29 wants a list of unique numbers up to 100**100 (100^100) ...
	// Factorized graduates from a []int type number to a structured number, and also stores the effective lengths ahead of time.
	// I'd like to make a version something like lenbase uint8 ; lenpow uint24 but the latter doesn't exist and the []uint16 (still worth it for data size in cache lines) is about to utilize abus-width int and pointer anyway...
	Lenbase uint
	Lenpow  uint
	Fact    []Factorpair
}

//func NewFactorized(primes *[]int, n uint) *Factorized {
//	return &Factorized{}
//}

func (facts *Factorized) Mul(fin *Factorized) *Factorized {
	temp := make([]Factorpair, 0, facts.Lenbase+fin.Lenbase)
	fbuf := (*FactorpairQueue)(&temp)
	heap.Init(fbuf)
	var fr, fi uint
	for fr < facts.Lenbase && fi < fin.Lenbase {
		// If BOTH have 1 as their base, add it as they're probably both 1...
		if facts.Fact[fr].Base == fin.Fact[fi].Base {
			heap.Push(fbuf, Factorpair{Base: facts.Fact[fr].Base, Power: facts.Fact[fr].Power + fin.Fact[fi].Power})
			fr++
			fi++
			continue
		}
		if 1 == facts.Fact[fr].Base {
			fr++
			continue
		}
		if 1 == fin.Fact[fi].Base {
			fi++
			continue
		}
		if facts.Fact[fr].Base >= fin.Fact[fi].Base {
			heap.Push(fbuf, fin.Fact[fi])
			fi++
		} else {
			heap.Push(fbuf, facts.Fact[fr])
			fr++
		}
	}
	for fr < facts.Lenbase {
		heap.Push(fbuf, facts.Fact[fr])
		fr++
	}
	for fi < fin.Lenbase {
		heap.Push(fbuf, fin.Fact[fi])
		fi++
	}
	leak := make([]Factorpair, 0, fbuf.Len())
	power := uint(0)
	for 0 < fbuf.Len() {
		fp := heap.Pop(fbuf).(Factorpair)
		power += uint(fp.Power)
		leak = append(leak, fp)
	}
	// fmt.Printf("Mul base check: %d\n", leak[0].Base)
	if 0 == leak[0].Base {
		facts.Lenbase = 0
		facts.Lenpow = 0
		power = 1
		leak = append(make([]Factorpair, 0, 1), Factorpair{Base: 0, Power: 1})
	}
	facts.Lenbase, facts.Lenpow, facts.Fact = uint(len(leak)), power, leak
	return facts
}

func (fl Factorized) Eq(fr *Factorized) bool {
	// I already wrote this and it's good to see that reflect.DeepEqual() is notably slower than a directed codepath. https://stackoverflow.com/a/15312182
	if fl.Lenbase != fr.Lenbase || fl.Lenpow != fr.Lenpow {
		return false
	}
	for ii := uint(0); ii < fl.Lenbase; ii++ {
		if fl.Fact[ii].Base != fr.Fact[ii].Base || fl.Fact[ii].Power != fr.Fact[ii].Power {
			return false
		}
	}
	return true
}

func (fl Factorized) Compare(fr *Factorized) int {
	// (unverified) Creating and comparing two BigInts is _PROBABLY_ expensive...
	if fl.Eq(fr) {
		return 0
	}
	left, right := fl.BigInt(), fr.BigInt()
	return left.Cmp(right)
}

func (fl Factorized) Cmp(fr *Factorized) int { return fl.Compare(fr) }

func (fl Factorized) BigInt() *big.Int {
	ret := big.NewInt(int64(1))
	for ii := uint(0); ii < fl.Lenbase; ii++ {
		base := big.NewInt(int64(fl.Fact[ii].Base))
		for ee := uint16(0); ee < fl.Fact[ii].Power; ee++ {
			ret = ret.Mul(ret, base) // math.Pow(x, y float64) float64 {...}
		}
	}
	return ret
}

func (fact *Factorized) Uint64() uint64 {
	ret := uint64(1)
	for ii := uint(0); ii < fact.Lenbase; ii++ {
		for ee := uint16(0); ee < fact.Fact[ii].Power; ee++ {
			ret *= uint64(fact.Fact[ii].Base)
		}
	}
	return ret
}

func (fact *Factorized) Copy() *Factorized {
	// len(Fact) SHOULD == Lenbase ... but this copies even not-normalized versions (without validating)
	ret := &Factorized{Lenbase: fact.Lenbase, Lenpow: fact.Lenpow, Fact: make([]Factorpair, len(fact.Fact))}
	copy(ret.Fact, fact.Fact)
	return ret
}

// Extract(transitive?)Power E.G.  4[^2] == (2^1, 2)[^2]
func (fact *Factorized) ExtractPower() (*Factorized, uint) {
	// Simplify the code at a small memory cost
	if 0 == fact.Lenbase {
		return fact.Copy(), 0
	}
	buf := make([]uint, 0, fact.Lenbase+1)
	for ii := uint(0); ii < fact.Lenbase; ii++ {
		buf = append(buf, uint(fact.Fact[ii].Power))
	}

	for terms := fact.Lenbase; 1 < terms; terms >>= 1 {
		// fmt.Printf("ExtractPower() Round: %v\n", buf)
		var ii uint
		for ii = 0; ii+1 < terms; ii += 2 {
			buf[ii>>1] = GCDbin(buf[ii], buf[ii+1])
			if 1 == buf[ii>>1] {
				return fact.Copy(), 1
			}
		}
		if ii+1 == terms {
			if 1 == buf[ii] {
				return fact.Copy(), 1
			}
			buf[ii>>1] = buf[ii]
		}
		terms++
		terms >>= 1 // flooring binary division
	}
	// fmt.Printf("ExtractPower() Final: %v\n", buf)
	if 1 == buf[0] {
		return fact.Copy(), 1
	}
	iiLim := len(fact.Fact)
	ret := &Factorized{Lenbase: fact.Lenbase, Lenpow: fact.Lenpow / buf[0], Fact: make([]Factorpair, iiLim)}
	for ii := 0; ii < iiLim; ii++ {
		ret.Fact[ii].Base = fact.Fact[ii].Base
		ret.Fact[ii].Power = fact.Fact[ii].Power / uint16(buf[0])
	}
	return ret, buf[0]
}

func (fact *Factorized) Pow(p uint) *Factorized {
	// Simplify the code at a small memory cost
	if 0 == p {
		return &Factorized{Lenbase: 1, Lenpow: 1, Fact: append(make([]Factorpair, 0, 1), Factorpair{Base: 2, Power: 0})}
	}
	iiLim := len(fact.Fact)
	ret := &Factorized{Lenbase: fact.Lenbase, Lenpow: fact.Lenpow * p, Fact: make([]Factorpair, iiLim)}
	for ii := 0; ii < iiLim; ii++ {
		ret.Fact[ii].Base = fact.Fact[ii].Base
		ret.Fact[ii].Power = fact.Fact[ii].Power * uint16(p)
	}
	return ret
}

func (fact *Factorized) PowDivMul(num, den uint) *Factorized {
	// Simplify the code at a small memory cost
	// Divide by zero is not legal, this is the closest I've got to NaN at the moment.
	if 0 == den {
		return &Factorized{}
	}
	if 0 == num {
		return &Factorized{Lenbase: 1, Lenpow: 1, Fact: append(make([]Factorpair, 0, 1), Factorpair{Base: 2, Power: 0})}
	}
	iiLim := len(fact.Fact)
	ret := &Factorized{Lenbase: fact.Lenbase, Lenpow: (fact.Lenpow / den) * num, Fact: make([]Factorpair, iiLim)}
	for ii := 0; ii < iiLim; ii++ {
		ret.Fact[ii].Base = fact.Fact[ii].Base
		ret.Fact[ii].Power = (fact.Fact[ii].Power / uint16(den)) * uint16(num)
	}
	return ret
}

func (f *Factorized) ProperDivisors() *[]uint64 {
	flen := len(f.Fact)
	if 0 == flen {
		return &[]uint64{1}
	}
	//if 1 == flen {
	//	return append(make([]uint, 0, 1), uint(f.Fact[0]))
	//}
	if flen > 64 {
		panic("Factorized.ProperDivisors() does not support more than 64 factors")
	}
	if uint(flen) != f.Lenbase {
		fmt.Printf("ERROR ProperDivisors(): Lenbase != len(Fact): %v\n", f)
	}
	sf := make([]uint, 0, f.Lenpow)
	for ii := uint(0); ii < f.Lenbase; ii++ {
		for pp := uint16(0); pp < f.Fact[ii].Power; pp++ {
			sf = append(sf, uint(f.Fact[ii].Base))
		}
	}
	var limit uint64
	if 64 == f.Lenpow {
		limit ^= 1
	} else {
		limit = (uint64(1) << f.Lenpow) - 1
	}

	almost := uint64(1)
	for ff := uint(1); ff < f.Lenpow; ff++ {
		almost *= uint64(sf[ff])
	}
	bitVec := bitvector.NewBitVector(almost)
	bitVec.Set(1)      // All 0s
	bitVec.Set(almost) // ^1 // ~1
	for ii := uint64(1); ii < limit; ii++ {
		bit := uint64(1)
		ar := uint64(1)
		for ff := uint(0); ff < f.Lenpow; ff++ {
			if 0 < ii&bit {
				ar *= uint64(sf[ff])
			}
			bit <<= 1
		}
		bitVec.Set(ar)
	}
	res := bitVec.GetUInt64s()
	// fmt.Printf("ProperDivisors() %d : 1 .. %d ??\t%v\n", f.Lenpow, almost, res)
	return res
}

/*
func FactorsToProperDivisors(factors *[]int) *[]int {
	fl := len(*factors)
	if 0 == fl {
		return factors
	}
	if 2 > fl {
		return &[]int{1}
	}
	if fl > 63 {
		panic("FtD does not support more than 63 factors.")
	}
	limit := (uint64(1) << fl) - 1
	bitVec := bitvector.NewBitVector(uint64(ListMul((*factors)[1:])))
	bitVec.Set(uint64(1))
	for ii := uint64(0); ii < limit; ii++ {
		div := 1
		bb := uint64(1)
		for ff := 0; ff < fl; ff++ {
			if 0 < ii&bb {
				div *= (*factors)[ff]
			}
			bb <<= 1
		}
		bitVec.Set(uint64(div))
	}
	return bitVec.GetInts()
}
*/

// Priority Queue heap https://pkg.go.dev/container/heap@go1.22.6

// Factorpair has no dynamic/reference based storage, simple copy is fine
type FactorpairQueue []Factorpair

func (pq FactorpairQueue) Raw() *[]Factorpair {
	conv := ([]Factorpair)(pq)
	return &conv
}

// { return &(pq.([]Factorpair)) }

// func (pq FactorpairQueue) Len() int { return len(([]Factorpair)(pq)) }
func (pq FactorpairQueue) Len() int { return len(pq) }

func (pq FactorpairQueue) Less(quea, queb int) bool {
	// "less" holds items closer to the base of the array
	return pq[quea].Base < pq[queb].Base
}

func (pq FactorpairQueue) Swap(quea, queb int) {
	pq[quea], pq[queb] = pq[queb], pq[quea]
	// 'Item' lacks priority and lacks index
}

func (pq *FactorpairQueue) Push(fp any) {
	*pq = append(*pq, fp.(Factorpair))
}

func (pq *FactorpairQueue) Pop() any {
	n := len(*pq) - 1
	fp := (*pq)[n]
	*pq = (*pq)[0:n]
	return fp
}

// Factorpair has no dynamic/reference based storage, simple copy is fine
type UintQueue []uint

func (uq UintQueue) Raw() []uint {
	conv := ([]uint)(uq)
	return conv
}

func (uq UintQueue) Len() int { return len(uq) }

func (uq UintQueue) Less(queA, queB int) bool {
	// "less" holds items closer to the base of the array
	return uq[queA] < uq[queB]
}

func (uq UintQueue) Swap(queA, queB int) {
	uq[queA], uq[queB] = uq[queB], uq[queA]
	// 'Item' lacks priority and lacks index
}

func (uq *UintQueue) Push(fp any) {
	*uq = append(*uq, fp.(uint))
}

func (uq *UintQueue) Pop() any {
	n := len(*uq) - 1
	fp := (*uq)[n]
	*uq = (*uq)[0:n]
	return fp
}
