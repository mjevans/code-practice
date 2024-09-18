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
	// "os" // os.Stdout
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
	ii >>= BVprimeByteBitShift
	pg, pidx := ii/BVpagesize, ii%BVpagesize
	pgmax, pmax := uint(len(p.PV)), ((limit-3)>>BVprimeByteBitShift)%BVpagesize
	if pg <= pgmax {
		for pidx <= pmax {
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
	fmt.Printf("Unable to locate prime after %d under %d\t[%d/%d][%d/%d]\t", input, limit, pg, pgmax, pidx, pmax)
	pg, pidx = ii/BVpagesize, ii%BVpagesize
	fmt.Printf("started near [%d][%d]\n", pg, pidx)
	return 0
}

func (p *BVPrimes) Grow(limit uint) {
	if p.Last >= limit {
		fmt.Printf("Already above requested growth limit, %d, at %d\n", limit, p.Last)
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
		}
	}

	// FIXME: test case, end of BVl1 #0 is 1025 (1026 should also get eaten, but not marked since it's compressed)
	// test case, end of page 1 is 65537 (uint16max + 1)

	base := (p.Last - 3 + 2) // the prime number AFTER ii, E.G. 33 -> 34
	// bidx := (base & BVprimeByteBitMask) >> 1
	base >>= BVprimeByteBitShift
	// pg, pidx := base/BVpagesize, base%BVpagesize

	for line := ((p.Last - 3 + 2) >> BVprimeByteBitShift / uint(BVl1)); line <= cl1z; line++ {
		// pgmax = len(p.PV) // Extended above, cl1z is backed by real array(s)
		pg := (line * BVl1) / BVpagesize
		// inclusive last bit address to set, upper limit
		cl1max := (((line+1)*BVl1-1)*BVbitsPerByte + (BVbitsPerByte - 1)) % (BVpagesize * BVbitsPerByte)
		cl1maxNum := ((pg * BVpagesize) << BVprimeByteBitShift) + cl1max<<1 + 3

		// Emperically it didn't take _too_ long to generate up to 4097, however it became a total slog after that point
		// if 4096 > (line*BVl1<<1)+3 {
		if true {

			//	Last	Next	Pri	Correct	Diff	FlMod(2p)+p
			//	33	34	3	39	6	33
			//	33	34	5	35	2	35
			//	33	34	7	35	2	35
			//	33	34	11	55	22	33
			//	33	34	13	39	6
			//	33	34	17	51	18
			//	33	34	19	57	24
			//	33	34	23	69	36
			//	33	34	29	87	54
			//	33	34	31	93	60

			// bit-bit address to start the process, in this case also the 'base bit' beneath which no marks
			// cl1min := (((p.Last - 3 + 2) >> 1) - (line*uint(BVl1*BVbitsPerByte))%(BVpagesize*BVbitsPerByte))
			// if pg < pgmax {
			prime := uint(3)
			// fmt.Printf("Primes.Grow(%d) INIT [..%d] %d (%d/%d)\n", limit, cl1maxNum, prime, line, cl1z)
			for {
				if prime<<1 >= p.Last {
					newLast := cl1maxNum
					if prime<<1 < newLast {
						// fmt.Printf("cl1max newLast was %d\t", newLast)
						newLast = prime << 1
					}
					// if newLast > 960 {
					if false {
						lineS := ((pg * BVpagesize) << BVprimeByteBitShift) + ((line*BVl1)%(BVpagesize))<<BVprimeByteBitShift + 3
						lineE := cl1maxNum
						fmt.Printf("Adust p.Last %d := %d (%d..%d)\t%d\n", p.Last, newLast, lineS, lineE, prime)
					}
					p.Last = newLast
				}
				// Modern CPUs prefer a branchless path even with a couple possibly redundant operations
				// 33 -> 32 -> 33 -> 35  ~~ 34 -> 33 -> 33 -> 35
				next := uint(((p.Last - 1) | 1) + 2)
				// calculate the next modulus to mark, this will always be an odd multiple of prime of at least 3 * prime...
				flMod3p := (next/(prime<<1))*(prime<<1) + prime
				// cl1bb := (p.Last - 3 + (prime + (prime << 1) - (p.Last % (prime << 1)))) >> 1
				cl1bb := (flMod3p - 3) >> 1
				// if 5 == prime || 7 == prime {
				// fmt.Printf("Start %d ~ %d (at [%d][%d]|%x\n", prime, cl1bb<<1+3, pg, cl1bb>>(BVprimeByteBitShift-1), uint8(1)<<(cl1bb&BVprimeByteBitMaskPost))
				// }
				for cl1bb <= cl1max {
					// fmt.Printf("Mark %d ~ %d (at [%d][%d]|%x\n", prime, cl1bb<<1+3, pg, cl1bb>>(BVprimeByteBitShift-1), uint8(1)<<(cl1bb&BVprimeByteBitMaskPost))
					// if (5 == prime || 7 == prime) && cl1bb < (64-3)>>1 {
					// fmt.Printf("Mark %d ~ %d (at [%d][%d]|%x\n", prime, cl1bb<<1+3, pg, cl1bb>>(BVprimeByteBitShift-1), uint8(1)<<(cl1bb&BVprimeByteBitMaskPost))
					// }
					p.PV[pg][cl1bb>>(BVprimeByteBitShift-1)] |= uint8(1) << (cl1bb & BVprimeByteBitMaskPost)
					cl1bb += prime // compressed (/2) prime bitvector, this advances by the prime * 2
					// if cl1bb > 64 {
					// 	panic("debug")
					// }
				}
				// if prime > 2000 {
				// fmt.Printf("Grow prime %d\n", prime)
				// }
				prime = p.primeAfterUnsafe(prime, prime<<1)
				if 0 == prime {
					panic("primeAfterUnsafe Returned Zero")
				}
				if prime<<1 >= cl1maxNum {
					// fmt.Printf("Primes.Grow(%d) upto %d\tRound %d\n", limit, cl1maxNum, prime)
					break
				}
			}
			//continue
		} else {
			// fmt.Printf("Primes.Grow(%d) reached unsupported quantity, p.Last = %d\n", limit, p.Last)
			// panic(p.Last)
		}
		p.Last = cl1maxNum
		// fmt.Printf("Primes.Grow(%d) upto %d\t(%d/%d)\n", limit, p.Last, line, cl1z)

	}

}

/*
func Factor(primes *[]int, num int) *[]int {
	//
	// Public school factoring algorithm from memory...

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

// func GetPrimes(primes *[]int, primehunt int) *[]int

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

type Factorpair struct {
	base  uint16
	power uint16
}

type Factorized struct {
	// Euler 29 wants a list of unique numbers up to 100**100 (100^100) ...
	// Factorized graduates from a []int type number to a structured number, and also stores the effective lengths ahead of time.
	// I'd like to make a version something like lenbase uint8 ; lenpow uint24 but the latter doesn't exist and the []uint16 (still worth it for data size in cache lines) is about to utilize abus-width int and pointer anyway...
	lenbase uint
	lenpow  uint
	fact    []Factorpair
}

func NewFactorized(primes *[]int, n uint) *Factorized {
	return &Factorized{}
}
