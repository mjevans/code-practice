// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package euler

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=5
<p>$2520$ is the smallest number that can be divided by each of the numbers from $1$ to $10$ without any remainder.</p>
<p>What is the smallest positive number that is <strong class="tooltip">evenly divisible<span class="tooltiptext">divisible with no remainder</span></strong> by all of the numbers from $1$ to $20$?</p>


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

func FactorsToDivisors_old(factors []int) []int {
	if 12 < len(factors) {
		fmt.Println("FTD: ", ListMul(factors), len(factors), "=~", Factorial(len(factors)))
		return []int{}
	}
	divisors := make([]int, 0, Factorial(len(factors)))
	divisors = append(divisors, 1)
	for ii := 0; ii < len(factors); ii++ {
		mmlim := len(divisors)
		for mm := 0; mm < mmlim; mm++ {
			divisors = append(divisors, divisors[mm]*factors[ii])
		}
	}
	return CompactInts(divisors[:len(divisors)-1])
}

func FactorsToProperDivisors(factors []int) []int {
	fl := len(factors)
	if 0 == fl {
		return factors
	}
	if 2 > fl {
		return []int{1}
	}
	if fl > 63 {
		panic("FtD does not support more than 63 factors.")
	}
	limit := (uint64(1) << fl) - 1
	bitVec := bitvector.NewBitVector(int64(ListMul(factors[1:])))
	bitVec.Set(uint64(1))
	for ii := uint64(0); ii < limit; ii++ {
		div := 1
		bb := uint64(1)
		for ff := 0; ff < fl; ff++ {
			if 0 < ii&bb {
				div *= factors[ff]
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
*/
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
