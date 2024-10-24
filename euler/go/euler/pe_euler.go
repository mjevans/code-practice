// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package euler

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans
// MOST Code in this file is CC BY-SA 4.0 https://creativecommons.org/licenses/by-sa/4.0/
// An algorithm from a public paper over 40 years old has been included in one function, and it's source is cited.
// This algorithm has been adapted from Pascal to golang, but where possible the original variable names have been kept for academic clarity.

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


for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in $(seq 1 42) ; do go fmt $(printf "pe_%04d.go" "$ii") ; printf "NEXT: %s\n" $ii ; go run $(printf "pe_%04d.go" "$ii") || break ; read JUNK ; done



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
	"math"
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

// Go 1.18+ ALMOST has C style 'casting' but these all copy...

func Uint64[T ~uint64](c T) uint64 {
	return uint64(c)
}

func Uint[T ~uint](c T) uint {
	return uint(c)
}

func Uint32[T ~uint32](c T) uint32 {
	return uint32(c)
}

func Uint16[T ~uint16](c T) uint16 {
	return uint16(c)
}

func Uint8[T ~uint8](c T) uint8 {
	return uint8(c)
}

func Int64[T ~int64](c T) int64 {
	return int64(c)
}

func Int[T ~int](c T) int {
	return int(c)
}

func Int32[T ~int32](c T) int32 {
	return int32(c)
}

func Int16[T ~int16](c T) int16 {
	return int16(c)
}

func Int8[T ~int8](c T) int8 {
	return int8(c)
}

// These don't work though...
// Another work around:
// func (o []Card) SLUint8() []uint8 { return []uint8(o) }
/*
func SLUint64[T ~[]uint64](c T) []uint64 {
	return []uint64(c)
}

func SLUint[T ~[]uint](c T) []uint {
	return []uint(c)
}

func SLUint32[T ~[]uint32](c T) []uint32 {
	return []uint32(c)
}

func SLUint16[T ~[]uint16](c T) []uint16 {
	return []uint16(c)
}

func SLUint8[T ~[]uint8](c T) []uint8 {
	return []uint8(c)
}

func SLInt64[T ~[]int64](c T) []int64 {
	return []int64(c)
}

func SLInt[T ~[]int](c T) []int {
	return []int(c)
}

func SLInt32[T ~[]int32](c T) []int32 {
	return []int32(c)
}

func SLInt16[T ~[]int16](c T) []int16 {
	return []int16(c)
}

func SLInt8[T ~[]int8](c T) []int8 {
	return []int8(c)
}
*/

/*	I
 *	I
 *	I
 *	I
 *	I
 *	I
		for ii in *\/*.go ; do go fmt "$ii" ; done ; go test -v euler/
*/

// globals

const (
// Golang... WHY can't this be a constant?  If I had 16 manually defined uint8s it would be, and an array is different than a slice/vector; this is tediously annoying for no obvious reason.
// PrimesSmallU8 = [...]uint8{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53}
)

var (
	Primes        *BVPrimes
	PrimesSmallU8 [16]uint8 // {2, 3, 5, 7, 11,	13, 17, 19, 23, 29,	31, 37, 41}
	PSRand        *PSRandPGC32
)

func init() {
	Primes = NewBVPrimes()
	PrimesSmallU8 = [...]uint8{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53} // 41 required for reasons, 53 nice for 16 total numbers; 16 bytes of memory, 1/4th cache line
	PSRand = NewPSRandPGC32(0x4d595df4d0f33173, 1442695040888963407)                       // Seed could be _anything_, inc must be any odd constant values from https://en.wikipedia.org/wiki/Permuted_congruential_generator#Example_code
}

// Deprecated, DO NOT USE, no replacement planned
func PrintFactors(factors []int) {
	// Join only takes []string s? fff
	strFact := make([]string, len(factors), len(factors))
	for ii, val := range factors {
		strFact[ii] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strFact, ", "))
}

// Deprecated, DO NOT USE, see type Factorized
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

// Deprecated, DO NOT USE
func AlphaSum(str string) int64 {
	var ret, limit int64
	limit = int64(len(str))
	str = strings.ToUpper(str)
	for ii := int64(0); ii < limit; ii++ {
		ret += int64(byte(str[ii]) - 'A' + 1)
	}
	return ret
}

// Deprecated, DO NOT USE
func ListSum(scale []int) int {
	ret := 0
	for _, val := range scale {
		ret += val
	}
	return ret
}

// Deprecated, DO NOT USE
func ListSumUint64(scale []uint64) uint64 {
	ret := uint64(0)
	ll := len(scale)
	for ii := 0; ii < ll; ii++ {
		ret += scale[ii]
	}
	return ret
}

// Deprecated, DO NOT USE
func ListMul(scale []int) int {
	ret := 1
	for _, val := range scale {
		ret *= val
	}
	return ret
}

// Deprecated, DO NOT USE
func Factorial(ii int) int {
	ret := 1
	for ii > 1 {
		ret *= ii
		ii--
	}
	return ret
}

func FactorialUint64(ii uint64) uint64 {
	ret := uint64(1)
	for ii > 1 {
		ret *= ii
		ii--
	}
	return ret
}

func FactorialDivFactU64toBig(ii, div uint64) *big.Int {
	// fmt.Println(ii, div)
	ret := big.NewInt(int64(1))
	one := big.NewInt(int64(1))
	if 1 > div {
		div = 1
	}
	if 1 > ii || ii < div {
		ii = 1
	}
	bifl := big.NewInt(int64(div))
	bi := big.NewInt(int64(ii))
	// limit := 0xFFFF
	// Cmp == -  =>  bifl - bi
	for 0 > bifl.Cmp(bi) {
		ret.Mul(ret, bi)
		bi.Sub(bi, one)
		// limit--
		// if 0 == limit { 			panic("BigFactorial - Iter Limit Reached")		}
	}
	return ret
}

// I don't know any shortcuts, just do it the obvious way
func PowU64(n, pow uint64) uint64 {
	if 0 == pow {
		return 1
	}
	ret := n
	for pow--; pow > 0; pow-- {
		ret *= n
	}
	return ret
}

func SqrtU64(ii uint64) uint64 {
	// https://en.wikipedia.org/wiki/Integer_square_root#Algorithm_using_Newton's_method
	// f(x) { x*x } // dxf(x) { 2x }
	if 1 >= ii {
		return ii
	}
	var x0, x1 uint64
	x0 = ii >> 1 // must be above the answer
	x1 = (x0 + ii/x0) >> 1
	for x1 < x0 {
		x0 = x1
		x1 = (x0 + ii/x0) >> 1
	}
	return x0
}

func RootU64(ii, root uint64) uint64 {
	// https://en.wikipedia.org/wiki/Integer_square_root#Algorithm_using_Newton's_method
	// f(x) { x*x } // dxf(x) { 2x }
	if 0 == root {
		return 0
	}
	if 1 >= ii || 1 == root {
		return ii
	}
	var x0, x1, mid, test uint64
	rmm := root - 1
	x0 = ii / rmm // must be above the answer
	x1 = (x0 + ii/PowU64(x0, rmm)) / root
	// fmt.Printf("RootU64(%d, %d)\tx0: %d\tx1: %d\n", ii, root, x0, x1)
	for x1 < x0 {
		x0 = x1
		// fmt.Printf("x0: %d", x1)
		x1 = (rmm*x0 + ii/PowU64(x0, rmm)) / root
		// fmt.Printf("\tx1: %d\n", x1)
	}

	// Newton's method can be unstable;, this fails some integer tests like 2^3 = 8 ( Root(8,3) )
	if x1 > x0 && ii < PowU64(x1, root) {
		// fmt.Printf("Correcting: %d < %d (%d^%d)\n", ii, PowU64(x1, root), x1, root)
		for x0 != x1 {
			mid = (x0 + x1) >> 1
			test = PowU64(mid, root)
			if test == ii {
				return mid
			}
			if test < ii {
				x0 = mid + 1
			} else {
				x1 = mid
			}
		}
	}
	// fmt.Printf("RootU64(%d, %d)\tx0: %d\tx1: %d\n", ii, root, x0, x1)

	return x0
}

/*

	left, right := 0, len(sl)
	if 0 == right {
		return -1
	}
	right--
	// I looked up the correct algorithm, because close enough but not textbook was painful
	// https://en.wikipedia.org/wiki/Binary_search  Alternative Procedure
	for left != right {
		// Alt uses ceil but bitshift math floors so I've reversed the comparison and side motions
		pos := (left + right) >> 1
		if sl[pos] < val {
			left = pos + 1
		} else {
			right = pos
		}
	}
	if val == sl[right] || always {
		return right
	}
	return -1
}
*/

// NOTE: Precision matters greatly for higher roots and for larger numbers, in that order!
func RootI64(ii int64, root, precision uint32) int64 {
	est := RootF64(float64(ii), root, precision)
	if 0 > ii {
		est -= 0.1
	} else {
		est += 0.1
	}
	return int64(est)
}

// NOTE: Precision matters greatly for higher roots and for larger numbers, in that order!
func RootF64(ii float64, root, precision uint32) float64 {
	// root must be a positive whole integer
	if 0.0 == ii || 0 == root {
		return 0.0
	}
	if 1 == root {
		return ii
	}
	negative := ii < 0.0
	if negative {
		ii = -ii
	}
	// https://en.wikipedia.org/wiki/Nth_root#Computing_principal_roots
	// x(k+1) = (r - 1) / r * x(k) + ( ii / r ) * 1 / ( x(k) ^ (n - 1) )
	var x0, x1, cA, cB, cmpD, cmpF float64
	rf := float64(root)
	// root--
	rpow := float64(root - 1)
	cA, cB = rpow/rf, ii/rf
	x0 = ii / rf
	for ; 0 < precision; precision-- {
		x1 = cB / math.Pow(x0, rpow)
		// root-- above to hoist root-1 out of the loop
		// for ii := root; 0 < ii; ii-- {
		// 	x1 /= x0
		// }
		x1 += cA * x0
		cmpD, cmpF = math.Abs(x1-x0), math.Abs((math.Nextafter(x0, x1)-x0)*rf*2.0)
		if cmpD < cmpF {
			// fmt.Printf("RF64(%f,%d) EXIT x1: %.20f\tx0: %.20f\tFuz: %.40f\n", ii, root, x1, x0, cmpF)
			break
			// } else if 0 == precision&0x3 {
			//	fmt.Printf("RF64(%f,%d) x1: %.20f\tx0: %.20f\tFuz: %.40f\n", ii, root, x1, x0, cmpF)
		}
		x0 = x1
	}
	if negative {
		return -x1
	}
	// fmt.Printf("RF64(%f,%d) PreC x1: %.20f\tx0: %.20f\tFuz: %.40f\n", ii, root, x1, x0, cmpF)
	return x1
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

// Tested by Euler 0015
func PascalTri(row, index uint64) uint64 {
	// Pascal Triangle?
	// Rows and Cols start with 0
	// A row may be calculated in isolation by: (just roll forward)  NOTE: 'forward' might also be backwards since it's symmetrical!
	// F(n Row, k Index) = (F(n, k - 1) * ( n + 1 - k) / k )
	// 0	    1
	// 1	   1 1
	// 2	  1 2 1
	// 3	 1 3 3 1
	// 4	1 4 6 4 1
	// F(4 Row, 2 Index) = (F(4, 2 - 1) * ( 4 + 1 - 2) / 2 ) = (1 * ( 4 + 1 - 1) / 1 ) *  ( 4 + 1 - 2) / 2 )
	// if 0 == index { return 1 }
	if row+1 < index {
		return 0
	}
	if 0 < row && (row+1)>>1 <= index {
		index = row - index
	}
	ret := uint64(1)
	for ii := uint64(1); ii <= index; ii++ {
		ret = (ret * (row + 1 - ii)) / ii
	}
	return ret
}

func TriangleNumber(n uint64) uint64 {
	// Euler 45
	// n*(n+1) / 2
	return (n * (n + 1)) >> 1
}

func TriangleNumberReverseFloor(n uint64) uint64 {
	// (n * (n + 1)) / 2
	// 1/2 * n*n + 1/2 * n == Tn
	// n*n + n - 2Tn == 0
	// a = 1 ; b = 1 ; c = -2Tn
	// n == ( -1 [+/-] sqrt(1 + 8 * Tn) ) / ( 2 )
	// n == ( -1 + sqrt(1 + 8 * Tn) ) / ( 2 )
	return (uint64(math.Sqrt(float64(1)+float64(8)*float64(n))) - 1) >> 1
}

func IsTriangleNumber(n uint64) bool {
	return n == TriangleNumber(TriangleNumberReverseFloor(n))
}

func SquareNumber(n uint64) uint64 {
	return n * n
}

func SquareNumberReverseFloor(n uint64) uint64 {
	return uint64(math.Sqrt(float64(n)))
}

func IsSquareNumber(n uint64) bool {
	return n == SquareNumber(SquareNumberReverseFloor(n))
}

func PentagonalNumber(n uint64) uint64 {
	// Euler 44
	// ( n * ( 3*n - 1 ) ) / 2
	return (n * (3*n - 1)) >> 1
}

func PentagonalNumberReverseFloor(n uint64) uint64 {
	// Euler 44
	// Quadratic Formula (looked up)
	// Given a*x*x + b*x + c == 0
	// Xn == ( -b [+,-] sqrt(b*b - 4 * a * c) ) / ( 2 * a )
	// (n * (3*n - 1)) >> 1 == Pn
	// (3*n*n / 2) - n/2 == Pn
	// 3*n*n - n == 2 * Pn
	// 3*n*n - n - 2 * Pn == 0
	// c = -2*Pn ; b = -1 ; a = 3
	// Only care about positive answers so...
	// n == ( 1 + sqrt(1 + 24 * Pn) ) / ( 6 )
	return (1 + uint64(math.Sqrt(float64(1)+float64(24)*float64(n)))) / uint64(6)
}

func IsPentagonalNumber(n uint64) bool {
	return n == PentagonalNumber(PentagonalNumberReverseFloor(n))
}

func HexagonalNumber(n uint64) uint64 {
	// Euler 45
	// n*(2*n - 1)
	return ((n * n) << 1) - n
}

func HexagonalNumberReverseFloor(n uint64) uint64 {
	// 2 * n*n - n - Hn == 0
	// a = 2 ; b = -1 ; c = - Hn
	// Xn == ( -b [+,-] sqrt(b*b - 4 * a * c) ) / ( 2 * a )
	// n == ( 1 + sqrt(1 + 8 * Hn ) ) / ( 4 )
	return (uint64(math.Sqrt(float64(1)+float64(8)*float64(n))) + 1) >> 2
}

func IsHexagonalNumber(n uint64) bool {
	return n == HexagonalNumber(HexagonalNumberReverseFloor(n))
}

func HeptagonalNumber(n uint64) uint64 {
	// Euler 61
	// n*(5*n-3)/2
	return (n * (5*n - 3)) >> 1
}

func HeptagonalNumberReverseFloor(n uint64) uint64 {
	// Pn = n*(5*n-3)/2
	// 2 * Pn = 5*n*n - 3*n
	// 0 = 5*n*n - 3*n - 2 * Pn
	// a = 5 ; b = -3 ; c = -2*Hn
	// Xn == ( -b [+,-] sqrt(b*b - 4 * a * c) ) / ( 2 * a )
	// n == ( 3 + sqrt(9 + 40 * Hn ) ) / 10
	return (uint64(math.Sqrt(float64(9)+float64(40)*float64(n))) + 3) / 10
}

func IsHeptagonalNumber(n uint64) bool {
	return n == HeptagonalNumber(HeptagonalNumberReverseFloor(n))
}

func OctagonalNumber(n uint64) uint64 {
	// Euler 61
	// n*(3*n-2)
	return n * (3*n - 2)
}

func OctagonalNumberReverseFloor(n uint64) uint64 {
	// On = n*(3*n-2)
	// 0 = 3*n*n -2*n -On
	// a = 3 ; b = -2 ; c = -1*On
	// Xn == ( -b [+,-] sqrt(b*b - 4 * a * c) ) / ( 2 * a )
	// n == ( 2 + sqrt(4 + 12*c) / 6
	return (uint64(math.Sqrt(float64(4)+float64(12)*float64(n))) + 2) / 6
}

func IsOctagonalNumber(n uint64) bool {
	return n == OctagonalNumber(OctagonalNumberReverseFloor(n))
}

func NgonNumber(n, gon uint64) uint64 {
	switch gon {
	case 3:
		return TriangleNumber(n)
	case 4:
		return SquareNumber(n)
	case 5:
		return PentagonalNumber(n)
	case 6:
		return HexagonalNumber(n)
	case 7:
		return HeptagonalNumber(n)
	case 8:
		return OctagonalNumber(n)
	default:
		return 0
	}
}

func NgonNumberReverseFloor(n, gon uint64) uint64 {
	switch gon {
	case 3:
		return TriangleNumberReverseFloor(n)
	case 4:
		return SquareNumberReverseFloor(n)
	case 5:
		return PentagonalNumberReverseFloor(n)
	case 6:
		return HexagonalNumberReverseFloor(n)
	case 7:
		return HeptagonalNumberReverseFloor(n)
	case 8:
		return OctagonalNumberReverseFloor(n)
	default:
		return 0
	}
}

func IsNgonNumber(n, gon uint64) bool {
	switch gon {
	case 3:
		return n == TriangleNumber(TriangleNumberReverseFloor(n))
	case 4:
		return n == SquareNumber(SquareNumberReverseFloor(n))
	case 5:
		return n == PentagonalNumber(PentagonalNumberReverseFloor(n))
	case 6:
		return n == HexagonalNumber(HexagonalNumberReverseFloor(n))
	case 7:
		return n == HeptagonalNumber(HeptagonalNumberReverseFloor(n))
	case 8:
		return n == OctagonalNumber(OctagonalNumberReverseFloor(n))
	default:
		return false
	}
}

/*
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
*/

// CompactInts should behave like slices.Compact(slices.Sort())
// Deprecated, DO NOT USE
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

// Deprecated, DO NOT USE
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

// Euler 0054 - Cards!
// I was going to // type Card uint8 // but a raw uint8 will work with my existing library better and present less issues...
// card & 0xF == 4 bits for NON SUITE value E.G. 2 = 2 9=9 10=10 11=Jack ... ? 12=Queen 13=King 14=Ace OR 1=Ace depending on game rules...
// Unicode went with Spades, Hearts, Diamonds, and Clubs https://en.wikipedia.org/wiki/Playing_cards_in_Unicode

const (
	CardSpade   = 0x10
	CardHeart   = 0x20
	CardDiamond = 0x30
	CardClub    = 0x40
	// These probably won't get used much, more for documentation
	CardAceHigh = 14
	CardKing    = 13
	CardQueen   = 12
	// Tarot Knight? goes here but that would make straights annoying
	CardJack   = 11
	CardAceOne = 1
)

func CardParseENG(s string) uint8 {
	// Works with Euler 0054's syntax, note T = 10 , returns 0 on error
	var suite, num uint8
	lim := len(s)
	for ii := 0; ii < lim; ii++ {
		switch s[ii] {
		case ' ':
			continue
		case 'S':
			suite = CardSpade
		case 'H':
			suite = CardHeart
		case 'D':
			suite = CardDiamond
		case 'C':
			suite = CardClub
		case 'A':
			num = 14
		case 'K':
			num = 13
		case 'Q':
			num = 12
		case 'J':
			num = 11
		case 'T':
			num = 10
		default:
			if '2' <= s[ii] && s[ii] <= '9' {
				num = uint8(s[ii] - '0')
			}
		}
	}
	if 0 < suite && 0 < num {
		return suite | num
	}
	return 0
}

func CardCompareValue(a, b uint8) int {
	return int(a&0xF) - int(b&0xF)
}

func CardPokerScore(hand, pub []uint8) uint {
	// NOTE: This assumes NO DUP cards; this shouldn't matter for 'full' hands, but it'll have edge cases for high cards or '5 of a kind' (which isn't in the score system)
	// NOTE: Poker hands are size 5
	// Euler 0054 - There's art and sometimes music, but in the 80s / 90s people sold legit games that were scantly more than this, some UI and standard library.
	lh, lp := len(hand), len(pub)
	var score uint
	// po := make([]uint8, 0, lh+lp)
	// po = append(po, hand...)
	// po = append(po, pub...)

	//	Val	Suit	Num	Name	Desc
	//  0x1FF0_0000	same5	inc5	Royal (straight) Flush	AceHigh unbroken run of cards in the same suit
	//  0x1FF0_0000	same5	inc5	Straight Flush	unbroken run of cards in the same suit
	//  0x_F0F_0000	-	same4	4 of a kind (value)	Note the same 0xF (15) for 'card' value in Flush slot
	//  0x_F00_FF00	-	same3+2	Full House	3+2 of a kind (value) -- Uses 0xF (15) for 'card' value in Flush slot
	//  0x_F00_0000	same5	-	Flush	All cards in the same suit, but no other match
	//  0x__F0_0000	-	inc5	Straight All cards in sequence
	//	0x_F000	-	same3	Three of a kind
	//	0x__FF0	-	same2	Two Pair	2+2 of a kind
	//	0x___F0	-	same2	One Pair	2 of a kind
	//	0x____F		highest	High Card

	// var bvSpade, bvHeart, bvDiamond, bvClub, bvAny uint16
	// var  nSpade,  nHeart,  nDiamond,  nClub,  nAny uint8
	var bv [5]uint16 // 0 is 'any' 1-4 are suites
	var cc [5]uint8  // card counts
	var cv [16]uint8 // Count of values
	for ii := 0; ii < lh; ii++ {
		suite, num := (hand[ii]&0xF0)>>4, hand[ii]&0x0F
		cc[0]++
		cc[suite]++
		cv[num]++
		bvNum := uint16(1) << num
		bv[0] |= bvNum
		if 0 < bv[suite]&bvNum {
			fmt.Printf("Poker: Duplicate card (hand): %x\n", hand[ii])
		}
		if 0 == suite {
			fmt.Printf("Poker: Null suite? (hand): %x\n", hand[ii])
		}
		bv[suite] |= bvNum
	}
	for ii := 0; ii < lp; ii++ {
		suite, num := (pub[ii]&0xF0)>>4, pub[ii]&0x0F
		cc[0]++
		cc[suite]++
		cv[num]++
		bvNum := uint16(1) << num
		bv[0] |= bvNum
		if 0 < bv[suite]&bvNum {
			fmt.Printf("Poker: Duplicate card (public): %x\n", pub[ii])
		}
		if 0 == suite {
			fmt.Printf("Poker: Null suite? (public): %x\n", pub[ii])
		}
		bv[suite] |= bvNum
	}
	// fmt.Printf("card bitvectors:\n%16b\n%16b\n%16b\n%16b\n%16b\n", bv[0], bv[1], bv[2], bv[3], bv[4])
	// fmt.Printf("card bitvector 4: %16b\n", bv[4])
	var flush, straight uint8
	const stmask = 0b_11111
	// Scan for any flushes
	for ii := 1; ii < 5; ii++ {
		// at least a flush...
		if 5 <= cc[ii] {
			tv := bv[ii]
			var cur uint8
			cur = 4 // cards go from 0..15 a 5 wide match at 0 would have a high value of 4
			for 0 < tv {
				if stmask == tv&stmask {
					stmp := uint(0x1000_0000) | uint(cur)<<24 | uint(cur)<<20
					if score < stmp {
						straight = cur
						flush = cur
						score = stmp
					}
				}
				tv >>= 1
				cur++
			}
			if 0 == flush || flush != straight {
				if flush < cur-5 {
					flush = cur - 5
				}
			}
		}
	}
	if 0 < score {
		return score
	}
	if 5 <= cc[0] {
		tv := bv[0]
		var cur uint8
		cur = 4 // cards go from 0..15 a 5 wide match at 0 would have a high value of 4
		for 0 < tv {
			if stmask == tv&stmask {
				straight = cur
			}
			tv >>= 1
			cur++
		}
	}
	var v2low, v2, v3, v4, hc uint8
	for ii := 0; ii <= 15; ii++ {
		switch {
		case 4 <= cv[ii]:
			v4 = uint8(ii)
		case 3 == cv[ii]:
			v3 = uint8(ii)
		case 2 == cv[ii]:
			v2low = v2
			v2 = uint8(ii)
		case 1 == cv[ii]:
			hc = uint8(ii)
		}
	}
	//  0x_F0F0000	-	same4	4 of a kind (value)	Note the same 0xF (15) for 'card' value in Flush slot
	//  0x_F00FF00	-	same3+2	Full House	3+2 of a kind (value) -- Uses 0xF (15) for 'card' value in Flush slot
	//  0x_F000000	same5	-	Flush	All cards in the same suit, but no other match
	//  0x__F00000	-	inc5	Straight All cards in sequence
	if 0 < v4 {
		return 0xF00_0000 | uint(v4)<<16 | uint(hc)
	}
	if 0 < v3 && 0 < v2 {
		return 0xF00_0000 | uint(v3)<<12 | uint(v2)<<8
	}
	if 0 < flush {
		return uint(flush) << 24
	}
	if 0 < straight {
		return uint(straight) << 20
	}
	//		High card could matter for '4 of a kind' and these last cases which don't use up all the cards
	//	0xF000	-	same3	Three of a kind
	//	0x_FF0	-	same2	Two Pair	2+2 of a kind
	//	0x__F0	-	same2	One Pair	2 of a kind
	//	0x___F		highest	High Card
	score |= uint(hc)
	if 0 < v3 {
		return score | uint(v3)<<12
	}
	if 0 < v2 && 0 < v2low {
		return score | uint(v2)<<8 | uint(v2low)<<4
	}
	if 0 < v2 {
		return score | uint(v2)<<4
	}
	return score
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
			"Fifteen",
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
			"Forty",
			"Fifty",
			"Sixty",
			"Seventy",
			"Eighty",
			"Ninety"}
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

// Deprecated, DO NOT USE (the max path sums use it)
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaximumPathSum(tri [][]int) int {
	lnMx := len(tri)
	if 0 == lnMx {
		return 0
	}
	lnMx--
	dist := make([]int, len(tri[lnMx])+1)
	for line := lnMx; line >= 0; line-- {
		iiLim := len(tri[line])
		for ii := 0; ii < iiLim; ii++ {
			dist[ii] = tri[line][ii] + MaxInt(dist[ii], dist[ii+1])
		}
	}
	return dist[0]
}

func MaximumPathSumAppendShrink(dst, c []int) ([]int, int, int) {
	var maxV, maxI, ld, lc, ii int
	ld, lc = len(dst), len(c)
	// Invalid use, make error more obvious
	if 0 == lc || ld+1 < lc {
		fmt.Printf("MaximumPathSumAppendShrink: incorrect use: lists of length %d and %d provided\n", ld, lc)
		return []int{}, 0, 0
	}
	// dst	2 1 2 1 2 1 2
	// c	 0 1 4 5 3 7
	// dst	 2 3 6 7 5 9
	for ii = 0; ii < lc; ii++ {
		dst[ii] = c[ii] + MaxInt(dst[ii], dst[ii+1])
		if maxV < dst[ii] {
			maxV, maxI = dst[ii], ii
		}
	}
	return dst, maxV, maxI
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
		// fmt.Printf("NQDL EOF + %d >%s< final return '%s'\n", ii, string(data), data[start:ii])
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

/*
// 64 bit / 8 byte points don't make sense for shuffling a deck on a modern CPU
type DLLuint64 struct {
	prev, next *DLLuint64
	data uint64
}
*/
/*
Shuffles / Factorial Permutations

	// PermutationString answered Euler 24 sufficiently as such, but is it optimal?
	for slot := 0 ; slot < end ; slot++ {
		fact := Factorial(end - 1 - slot)
		idx := perm / fact
		perm %= fact
		pull card and shrink the deck by shift copy ;	}

More research yields Wikipedia's article hum, that's similar
https://en.wikipedia.org/wiki/Factorial_number_system#Definition

#24 asked; Permutation(1_000_000 - 1 , "0123456789")

*/

// CPU Timing tested in ex_permutation.go but needs platform specific "syscall"
// Usually ~20% faster, may require a temp buffer if unable to store the index to use in the output array.
//
//	8	7	6	5	4	3	2	1
//	7!	6!	5!	4!	3!	2!	1!	0!	Place Value
//	5040	720	120	24	6	2	1	0	Place Value (in Dec)
//	FactorialUint64(len(con)) == Looped list value, subtract 1
func PermutationSlUint8(perm uint64, con []uint8) []uint8 {
	end := len(con)
	tmp := make([]uint8, end)
	copy(tmp, con)
	res := make([]uint8, end)
	mod := uint64(1)
	seed := perm
	for ii := end; 0 < ii; ii-- {
		res[uint64(end)-mod] = uint8(seed % mod)
		seed /= mod
		mod++
	}
	for ii := 0; ii < end; ii++ {
		idx := res[ii]
		res[ii] = tmp[idx]
		for idx < uint8(end-1-ii) {
			tmp[idx] = tmp[idx+1]
			idx++
		}
	}
	return res
}

func Uint8DigitsToUint64(sl []uint8, base uint64) uint64 {
	var ret uint64
	ii := len(sl)
	for 0 < ii {
		ii--
		ret *= base
		ret += uint64(sl[ii])
	}
	return ret
}

func Uint64ToDigitsUint8(n, base uint64) []uint8 {
	ret := make([]uint8, 0, 8)
	if 0 == n {
		ret = append(ret, 0)
	}
	for n > 0 {
		ret = append(ret, uint8(n%base))
		n /= base
	}
	return ret
}

func Uint8CopyInsertSort(con []uint8) []uint8 {
	l := len(con)
	ret := make([]uint8, 0, l)
	for ii := 0; ii < l; ii++ {
		ret = append(ret, con[ii])
		for jj := ii; jj > 0; jj-- {
			if ret[jj] < ret[jj-1] {
				ret[jj], ret[jj-1] = ret[jj-1], ret[jj]
			} else {
				break // jj
			}
		}
	}
	return ret
}

func Uint8Reverse(mod []uint8) []uint8 {
	l := len(mod)
	lh := l >> 1
	for ii := 0; ii < lh; ii++ {
		mod[ii], mod[l-1-ii] = mod[l-1-ii], mod[ii]
	}
	return mod
}

func Uint8Compare(a, b []uint8) int {
	la, lb := len(a), len(b)
	if la != lb {
		return la - lb
	}
	for ii := 0; ii < la; ii++ {
		if a[ii] != b[ii] {
			return int(a[ii]) - int(b[ii])
		}
	}
	return 0
}

func PermutationString(perm int64, str string) string {
	end := int64(len(str))
	tmp := make([]byte, end)
	copy(tmp, str)
	res := make([]byte, end)
	slot := int64(0)
	for slot < end {
		fact := int64(FactorialUint64(uint64(end - 1 - slot)))
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

func ConcatDigitsU64(x, y, base uint64) uint64 {
	var pow, temp uint64
	pow = base
	for temp = y / base; 0 < temp; temp /= base {
		pow *= base
	}
	return x*pow + y
}

func Pandigital(test uint64, used, reserved uint16) (fullPD bool, biton, usedRe uint16, DigitShift uint64) {
	DigitShift = uint64(1)
	ok := true
	if 0 == test {
		return false, 0, 1, 10
	}
	for test > 0 {
		bd := uint16(1) << (test % 10)
		test /= 10
		DigitShift *= 10
		if 0 < used&bd || 1 == reserved&bd {
			// fmt.Printf("SKIP: %d : dupe or 0 digit : %d\n", nn, test%10)
			ok = false
		}
		used |= bd
	}
	if ok {
		bt := uint16(0)
		// t := used AND (allowed) (XOR reserved)
		for t := used & (^reserved) >> 1; 0 < t; t >>= 1 {
			bt++
		}
		// Turn the bitcount into a binary similar to 0b_xxxx_xxx0 by shifting one more bit left, then subtracting the zero bit, and the bit to make it a 2s compliment negative number.
		return uint16((uint64(1)<<(bt+1))-2) == used, bt, used, DigitShift
	}
	return false, 0, used, DigitShift
}

func PalindromeFlipBinary(x uint64) uint64 {
	var ret uint64
	for 0 < x {
		ret <<= 1
		ret |= x & 1
		x >>= 1
	}
	return ret
}

func IsPalindrome(x, base uint64) bool {
	buf := make([]uint8, 0, 24)
	if 255 < base || 0 == base {
		fmt.Printf("ERROR: IsPalindrome does not support base = %d\n", base)
		return false
	}
	for 0 < x {
		buf = append(buf, uint8(x%base))
		x /= base
	}
	lim := len(buf)
	for ii := 0; ii<<1 < lim; ii++ {
		// 0..3 (4) 4-0-1 == 3~0 4-1-1 == 2~1 // (5) 5-0-1 4:0 5-1-1 3:1 5-2-1 2:2
		if buf[ii] != buf[lim-ii-1] {
			return false
		}
	}
	return true
}

func PalindromeMakeDec(x, addZeros uint64, odd bool) uint64 {
	ret := x
	buf := make([]uint8, 0, 20)
	pow := uint64(1)
	for 0 < x {
		buf = append(buf, uint8(x%10))
		x /= 10
		pow *= 10
	}
	if true == odd && 0 == addZeros {
		buf = buf[0 : len(buf)-1]
	}
	if odd && 0 < addZeros {
		pow /= 10
	}
	for 0 < addZeros {
		pow *= 100
		addZeros--
	}
	for ii := len(buf) - 1; 0 <= ii; ii-- {
		ret += uint64(buf[ii]) * pow
		pow *= 10
	}
	return ret
}

func SlicePopUint8(deck []uint8, index int) uint8 {
	// deck is a slice which MAY contain 'blanked' / empty 'card' slots (value 0), that are skipped in the count.
	// Shuffle would use a random index to obtain the next card
	// __IMPORTANT__ The size of deck does not mutate, only a VALUE within the pointed to array will change, hence this should work.  Thanks Go... your pain here probably improved Rust.
	dLen := len(deck)
	for ii := 0; ii < dLen; ii++ {
		if 0 != deck[ii] {
			if 0 == index {
				ret := deck[ii]
				deck[ii] = 0
				return uint8(ret)
			}
			index--
		}
	}
	return 0
}

func SliceCommon[S ~[]E, E ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string](a S, b S) S {
	al, bl := len(a), len(b)
	t := al
	if t > bl {
		t = bl
	}
	res := make(S, 0, t)
	ai, bi := 0, 0
	for ai < al && bi < bl {
		if a[ai] == b[bi] {
			res = append(res, a[ai])
			ai++
			bi++
			continue
		}
		for ai < al && a[ai] < b[bi] {
			ai++
		}
		for bi < bl && a[ai] > b[bi] {
			bi++
		}
	}
	return res
}

// https://stackoverflow.com/a/70562597 Go 1.21 added cmp.Ordered
func BsearchSlice[S ~[]E, E ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string](sl S, val E, always bool) int {
	left, right := 0, len(sl)
	if 0 == right {
		return -1
	}
	right--
	// I looked up the correct algorithm, because close enough but not textbook was painful
	// https://en.wikipedia.org/wiki/Binary_search  Alternative Procedure
	for left != right {
		// Alt uses ceil but bitshift math floors so I've reversed the comparison and side motions
		pos := (left + right) >> 1
		if sl[pos] < val {
			left = pos + 1
		} else {
			right = pos
		}
	}
	if val == sl[right] || always {
		return right
	}
	return -1
}

func UnsortedSearchSlice[S ~[]E, E ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string](sl S, val E) int {
	mx := len(sl)
	for ii := 0; ii < mx; ii++ {
		if val == sl[ii] {
			return ii
		}
	}
	return -1
}

func BsearchInt(list *[]int, val int) bool {
	idx := BsearchSlice(*list, val, false)
	if 0 > idx {
		return false
	}
	return val == (*list)[idx]
}

/*
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
			pos += (end + 1 - pos) >> 1
		} else { // gt
			end = pos - 1
			pos -= (pos + 1 - left) >> 1
		}
		// fmt.Printf("BsearchInt: next\t%d <= %d <= %d\t%d\n", left, pos, end, (*list)[pos])
	}
	// fmt.Printf("BsearchInt: false : %d\n", val)
	return false
}
*/

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

/*
https://en.wikipedia.org/wiki/Category:Pseudorandom_number_generators
Older, but might have been heard of https://en.wikipedia.org/wiki/Mersenne_Twister

Current state of good enough art for projects like this (if /dev/urandom isn't easy to use..)

https://en.wikipedia.org/wiki/Permuted_congruential_generator

Current (2020+) Linux /dev/random state of the art
https://en.wikipedia.org/wiki/BLAKE_(hash_function)#BLAKE2

*/

func Rotr32(reg uint32, rot uint8) uint32 {
	return reg>>rot | reg<<((^rot)&31)
}

const PSRandPGC32mul = 6364136223846793005

type PSRandPGC32 struct {
	State, Inc uint64 // Inc can be any arbitrary odd constant
}

func NewPSRandPGC32(seed, inc uint64) *PSRandPGC32 {
	r := &PSRandPGC32{State: seed + inc, Inc: inc}
	r.RandU32()
	return r
}

func (r *PSRandPGC32) RandU32() uint32 {
	count := uint8(r.State >> 59) // 64 - 5 (bit size of output index)
	temp := r.State
	r.State = temp*PSRandPGC32mul + r.Inc
	temp ^= temp >> 18                     // (64 - (32 - 5)) / 2
	return Rotr32(uint32(temp>>27), count) // 32 - 5 = 27
}

/*
https://en.wikipedia.org/wiki/Modular_arithmetic

https://en.wikipedia.org/wiki/Modular_arithmetic#Integers_modulo_m



*/

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

	Offhand, from a pragmatic viewpoint, a list of primes between 0 and the largest under 65536 is _probably_ more memory than a practical program should use, though 0..255 is clearly too limited.
	[]uint16 might be a good format for the primes list, if not a bitvector directly.

	2..7919 contains 1000 prime numbers; stored as a compressed (inherently 2 is prime so 3..7919) bitvector, that would take 3958 bits or 495 bytes (rounded up)
	It's entirely practical to throw a 512 or 4096 byte bucket of primes at the issue and simplify life.
	Page ~= 64 bytes
	3..130 = Page 0 The highest prime is 127 which has a square root of ~11.27 (121<>144)
	3..18 == BYTE 0 17,13,11,7,5,3
**/

// BVpagesize >= BVl1 // Both MUST be a power of 2 ( Pow(2, n) )
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
	Last uint64
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
	// WARNING: Set Last according to last tested value, UNTESTED but 2-7 should work with all 0s as the first three bits (3,5,7) are primes (0).
	return &BVPrimes{PV: append(make([]*BVpage, 0, 1), ov), Last: 33}
}

func (p *BVPrimes) PrimeOrDown(ii uint64) uint64 {
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
					return ((pg*BVpagesize + pidx) << BVprimeByteBitShift) + uint64(bidx)<<1 + 3
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
			// fmt.Println("Hit the floor fallback = 3; PrimeOrDown()")
			return 3
		}
		pg--
		pidx = BVpagesize - 1
	}
}

func (p *BVPrimes) PrimeRemove(ii uint64) {
	// testII := ii
	if ii > p.Last {
		return
	}
	ii = (ii - 3)
	bidx := (ii & BVprimeByteBitMask) >> 1
	ii >>= BVprimeByteBitShift
	pg, pidx := ii/BVpagesize, ii%BVpagesize
	// test := ((pg*BVpagesize + pidx) << BVprimeByteBitShift) + uint(bidx)<<1 + 3
	// if 1487 == testII {
	// fmt.Printf("PrimeRemove: %d = %d\n", testII, test)
	// }
	p.PV[pg][pidx] |= (uint8(1) << bidx)
}

func (p *BVPrimes) PrimeAfter(ii uint64) uint64 {
	if 2 > ii {
		return 2
	}
	// Guard an underflow, 2 doesn't really exist to step after
	if 2 == ii {
		return 3
	}
	lastPrime := p.PrimeOrDown(p.Last)
	if ii >= lastPrime {
		newLimit := (((((ii-3)>>1)/uint64(BVl1))+uint64(1))*uint64(BVl1)+uint64(BVprimeByteBitMaskPost))<<1 + 3
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

func (p *BVPrimes) primeAfterUnsafe(input, limit uint64) uint64 {
	ii := (input - 3 + 2) // the prime number AFTER ii, E.G. 6 -> 7
	bidx := (ii & BVprimeByteBitMask) >> 1
	bidx0 := bidx
	ii >>= BVprimeByteBitShift
	pg, pidx := ii/BVpagesize, ii%BVpagesize
	bbMM := ((limit - 3) & BVprimeByteBitMask) >> 1
	ooMM := (limit - 3) >> BVprimeByteBitShift
	pgMM, idxMM := ooMM/BVpagesize, ooMM%BVpagesize
	pgmax, pimax, pbmax := uint64(len(p.PV)), uint64(BVpagesize-1), uint64(BVprimeByteBitMaskPost)
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
					return ((pg*BVpagesize + pidx) << BVprimeByteBitShift) + uint64(bidx)<<1 + 3
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
	var cl, end uint64
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

func (p *BVPrimes) wheelFactCL1Unsafe(start, prime, maxPrime uint64) (uint64, uint64) {
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
		start = uint64(((p.Last - 1) | 1) + 2)
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

func (p *BVPrimes) autoFactorPMCforBVl1(start uint64) {
	if 0 == start {
		// 2 + (IF even (step back to last odd) ELSE odd noop )
		start = uint64(((p.Last - 1) | 1) + 2)
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

func (p *BVPrimes) libraryProbPrimeBVl1(start uint64) {
	if 0 == start {
		// 2 + (IF even (step back to last odd) ELSE odd noop )
		start = uint64(((p.Last - 1) | 1) + 2)
	}
	ABS_bbStart := (start - 3) >> 1
	ABS_ooStart := ABS_bbStart >> (BVprimeByteBitShift - 1)
	pg, pgLine := ABS_ooStart/BVpagesize, (ABS_ooStart%BVpagesize)/BVl1
	bbLimit := ((pgLine*BVl1 + BVl1 - 1) << (BVprimeByteBitShift - 1)) | BVprimeByteBitMaskPost
	bbPos := ABS_bbStart % (BVpagesize * BVbitsPerByte)
	for bbPos <= bbLimit {
		if 0 == p.PV[pg][bbPos>>(BVprimeByteBitShift-1)]&(uint8(1)<<(bbPos&BVprimeByteBitMaskPost)) {
			posPrime := (pg * BVpagesize << BVprimeByteBitShift) + (bbPos << 1) + 3
			if false == big.NewInt(int64(posPrime)).ProbablyPrime(int(8)) {
				// fmt.Printf(" %d", posPrime)
				p.PV[pg][bbPos>>(BVprimeByteBitShift-1)] |= uint8(1) << (bbPos & BVprimeByteBitMaskPost)
			}
		}
		bbPos++
	}
}

func (p *BVPrimes) PrimesOnPage(start uint64) []uint64 {
	if 0 == start {
		// 2 + (IF even (step back to last odd) ELSE odd noop )
		start = uint64(((p.Last - 1) | 1) + 2)
	}
	ret := make([]uint64, 0, 64)
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

func (p *BVPrimes) Grow(limit uint64) {
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
	cl1z := (((limit - 3) >> BVprimeByteBitShift) / uint64(BVl1)) + uint64(1)

	// Ensure the bitvector arrays exist
	pagez := cl1z/(BVpagesize/BVl1) + 1
	lenpv := uint64(len(p.PV))
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
	line := ((next - 3) >> BVprimeByteBitShift) / uint64(BVl1)

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
	// for line <= cl1z && line < 20480 {
	// vs the big.Int.ProbablyPrime test ~2048 lines was the fastest tested cutoff.
	// 1024 ~= 1M/16 ~ 41.36s
	// 2048 ~= 2M/16 ~ 40.91s
	// 3072 ~= 3M/16 ~ 41.88s
	// 4096 ~= 4M/16 ~ 43.65s
	for line <= cl1z && line < 2048 {
		primeStart := uint64(3)
		var end uint64
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

	// if line <= cl1z {
	// fmt.Printf("Primes wheel factorized to %d, switching gears to wheel filter and ProbablyPrime test.\n", p.Last)
	// }
	// AFTER 64K this is somehow so fast that I suspect the results... FIXME: Have I added a total torture test to cover up to 1024*1024 yet?

	for line <= cl1z {
		// p.wheelFactCL1Unsafe(next, 3, 509) // 5 = 498062; 503 = 498062
		// p.autoFactorPMCforBVl1(next)
		p.wheelFactCL1Unsafe(next, 3, 41)
		p.libraryProbPrimeBVl1(next)
		next = 3 + (((line + 1) * BVl1) << BVprimeByteBitShift)
		p.Last = next - 2
		// ccount := len(p.PrimesOnPage(p.Last))
		// fmt.Printf("%d: %d\n", line, p.Last)
		line++
	}

}

func (p *BVPrimes) MaybePrime(q uint64) bool {
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

func (p *BVPrimes) KnownPrime(q uint64) bool {
	// Use base2 storage inherent test for division by 2
	if (0 == q&0x01 && 2 < q) || q > p.Last {
		return false
	}
	pd := p.PrimeOrDown(q)
	return pd == q
}

func (p *BVPrimes) ProbPrime(q uint64) bool {
	// Use base2 storage inherent test for division by 2
	if 0 == q&0x01 && 2 < q {
		return false
	}
	pd := p.PrimeOrDown(q)
	if pd == q {
		return true
	}
	pLim := len(PrimesSmallU8)
	for ii := 0; ii < pLim; ii++ {
		if 0 == q%uint64(PrimesSmallU8[ii]) {
			return false
		}
	}
	return big.NewInt(int64(q)).ProbablyPrime(int(8))
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
	prime := uint64((*primes)[ii-1])
	// fmt.Printf("GetPrimesInt: cap %d\n", lim)
	for ; ii < lim; ii++ {
		prime = p.PrimeAfter(prime)
		*primes = append(*primes, int(prime))
	}
	// fmt.Printf("GetPrimesInt: %v\n", *primes)
	return primes
}

func (p *BVPrimes) CountPrimesLE(ii uint64) uint64 {
	if 2 > ii {
		return 0
	}
	if 2 == ii {
		return 1
	}
	lastPrime := p.PrimeOrDown(p.Last)
	if ii >= lastPrime {
		newLimit := (((((ii-3)>>1)/uint64(BVl1))+uint64(1))*uint64(BVl1)+uint64(BVprimeByteBitMaskPost))<<1 + 3
		// fmt.Printf("Primes.PAfter .Grow triggered:   \t%d\t< %d\t-> %d\n", ii, p.Last, newLimit)
		p.Grow(newLimit)
	}
	return p.countPrimesLEUnsafe(ii)
}
func (p *BVPrimes) countPrimesLEUnsafe(ii uint64) uint64 {
	// 2 isn't in the list, it's implied
	ret := uint64(1)
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

// func GCDbin(a, b uint) uint {
func GCDbin[INT ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](a, b INT) INT {
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
	if 0 > a {
		a = -a
	}
	if 0 > b {
		b = -b
	}

	//fmt.Printf("GCDbin %d, %d\n", a0, b0)
	//if a0 > 0xffffff || b0 > 0xffffff {
	//	panic("overflow") }
	if 0 == b {
		return a
	}
	if 0 == a {
		return b
	}

	var k, ka, kb uint8
	// k == count of common 2 factors
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
func Factor1980PollardMonteCarlo(N, x0 uint64) uint64 {
	// https://en.wikipedia.org/wiki/Pollard%27s_rho_algorithm

	// https://maths-people.anu.edu.au/~brent/pub/pub051.html
	// https://maths-people.anu.edu.au/~brent/pd/rpb051i.pdf
	//if 0 == x0 { x0 = 2 }

	// if x0 > 0xffffff {
	//	fmt.Printf("x0 unlikely large: %d\n", x0)
	//	panic("unexpected value")
	//}

	// print pg182-183 p"2
	fx := func(x, Nin uint64) uint64 {
		return (x*x + 3) % Nin
	}
	umin := func(a, b uint64) uint64 {
		if a < b {
			return a
		}
		return b
	}
	abssub := func(a, b uint64) uint64 {
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
	var y, r, q, m uint64
	y, r, q, m = uint64(x0), 1, 1, 1
	var k, ii, G, x, ys uint64
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

func Factor1980AutoPMC(q uint64, singlePrimeOnly bool) uint64 {
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
	return Factor1980AutoPMC_Pass2(q, singlePrimeOnly)
}

func Factor1980AutoPMC_Pass2(q uint64, singlePrimeOnly bool) uint64 {
	var pollard, pollard_limit, roottest uint64
	pollard = 1

	pollard_limit = SqrtU64(q)
	if q == pollard_limit*pollard_limit {
		return pollard_limit
	}

	for ii := uint64(3); ii < 5; ii++ {
		roottest = RootU64(q, ii)
		if q == PowU64(roottest, ii) {
			return roottest
		}
	}

	/* Old method
	// __approximate__ an integer square root, POW(pollard_limit, 2) _MUST_ be > q == pl*pl means square root factor
	pollard := uint64(1)
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
	*/

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

func (p *BVPrimes) Factorize(q uint64) *Factorized {
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
		heap.Push(facts, Factorpair{Base: 2, Power: uint32(k)})
	}

	// Quickly test some small primes; 2, 3 (~66%), 5 (~73%), 7 (<77%) -- https://en.wikipedia.org/wiki/Wheel_factorization#Description
	// smallPrimes := []uint32{3, 5, 7, 11}
	// for cur := 0; 1 < q && cur < len(smallPrimes); cur++ {
	pMax := len(PrimesSmallU8) - 1
	var base, power uint32
	// Start at 1 : skip already checked 2
	for ii := 1; ii <= pMax; ii++ {
		qd := uint64(PrimesSmallU8[ii])
		for 0 == q%qd {
			q /= qd
			power++
		}
		if 0 < power {
			heap.Push(facts, Factorpair{Base: uint32(qd), Power: uint32(power)})
			power = 0
		}
	}

	// zz := 1000
	// for 1 < q && zz > 0 {
	for 1 < q {
		var unk uint64

		// is the number already known to be a prime?
		if p.KnownPrime(q) {
			heap.Push(facts, Factorpair{Base: uint32(q), Power: 1})
			break
		}

		// FIXME: Factorization

		// try just one iteration, returns 0 on 'no factors found' (but search not exhausted)
		// KEEP one pass (this is faster when it works!)
		unk = Factor1980PollardMonteCarlo(q, 0)

		// Aiming to replace this with code that extracts factors rather than just testing for factors...
		/*
			roottest = SqrtU64(q)
			if q == roottest*roottest {
				return roottest
			}

			// Test higher roots until beneath the highest division tested prime, 41^4 ~= 2.8M ; 41^5 ~= 115.8M
			for ii := uint64(3); roottest > PrimesSmallU8[pMax] ; ii++ {
				roottest = RootU64(q, ii)
				if q == PowU64(roottest, ii) {
					return roottest
				}
			}

			See beneath for the better tests...
		*/

		// aggressive search
		if unk == 0 {
			if big.NewInt(int64(q)).ProbablyPrime(int(8)) {
				unk = q
			} else {
				// Replace this with LECF
				unk = Factor1980AutoPMC_Pass2(q, false)
			}
		}

		// Probably Prime
		if unk == q {
			heap.Push(facts, Factorpair{Base: uint32(q), Power: 1})
			break
		}

		q /= unk
		sf := p.Factorize(unk)
		for ii := uint32(0); ii < sf.Lenbase; ii++ {
			b := uint64(sf.Fact[ii].Base)
			pow := uint32(0)
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

	base, power = 0, 0
	fact := make([]Factorpair, 0, facts.Len())
	for 0 < facts.Len() {
		fp := heap.Pop(facts).(Factorpair)
		base++
		power += uint32(fp.Power)
		fact = append(fact, fp)
	}
	return &Factorized{Lenbase: base, Lenpow: power, Fact: fact}
}

func ProbablyPrimeI64(num, kStrBases int64) bool {
	return big.NewInt(num).ProbablyPrime(int(kStrBases))
}

/*


FIXME: At the very least, clean this section up in a future update...

It was very difficult to find good reference material for the harder 'math' subject elements.  Wikipedia's pages tend to be higher level for most subjects, glossing over key steps that someone in the field likely assumes a reader would know.  Similarly they use maths notation since someone who's not bumping into a subject on an extremely rare basis should know that shorthand.

I'm also not ENTIRELY sure the addition function is correct.  slopeY / slopeX is a rational not necessary an integer...



Integer Factorization is likely very similar to Primality (Probably Prime) tests, with the difference that proof (a factor that cleanly divides) must be found, rather than statistical lack of evidence of a proof.

Step 1 : trial division / wheel factorization / filtering of known small primes.  Just shoot a bunch at it; really only the first 4 primes are required but, __testing higher primes helps later nth root tests!__

Note: due to 2 and 3 in the derivative, primes 2 and 3 MUST be ruled out early for EC tests.

// https://stackoverflow.com/a/2274520

Step ? : (maybe?) ONE round of 1980PollardMonteCarlo

Step ? : nth root tests (down to largest tested prime?) this scales up pretty fast nice 41^5 ~= 115.8M

Step ? : Lenstra Elliptic Curve Factorization (good up to ~10^50 or beyond, and thus well past uint64)

Step +1 : <= 10^100 Quadratic Sieve
Step +2 : >= 10^100 General Number Field Sieve

https://en.wikipedia.org/wiki/Lenstra_elliptic-curve_factorization

Wikipedia summary:
Factor N
1.a :	Pick a random Curve over (mod N) using an equation of the form: y*y = x*x*x - a*x + b (mod N) ((NOTE: FIXME is that modN for every operation, or for both sides?))
1.b :	Pick a non-trivial (?) point on it P(x, y) ; which can be done most easily working _backwards_ from x, y ; b = (F(x0, y0)) ~~ y*y  - x*x*x - a * x (mod N)

2.	'Addition' of two points on the curve can define a Group (see Wikipedia) ... ?? 'division by some V mod N includes calculation of GCD(V, N)
	Math curves group theory I don't understand well enough to confidently summarize.
	When it has a GCD != 1 that indicates a 'non-trivial' factor of N

3.	Compute a sum of curves? ( mod N ) ; k is the product ( * ) of many small primes raised to small powers (**Maybe I can use PrimesSmallU8 rather than Factorial?)
	If no 'non-invertable' elements were found, then the curves' ( mod **Primes ) order wasn't Smooth enough, try again with a different starting Point (and thus A+B)
	If gcd(V,N) != 1, n that's a non-trivial factor of N.

// FIXME: Twisted Edwards curve variation is better? https://en.wikipedia.org/wiki/Lenstra_elliptic-curve_factorization#Twisted_Edwards_curves
*/

type Point3DI64 struct {
	X, Y, Z int64
}

func ExtendedGCDI64(a, b int64) (int64, int64, int64) {
	var quo, r, rP, s, sP, t, tP int64
	// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
	rP, r = a, b
	sP, s = 1, 0
	tP, t = 0, 1

	for 0 != r {
		// Haven't tested it, but golang is probably smart enough to collect the remainder and quotient in one div operation
		rP, quo, r = r, rP/r, rP%r
		sP, s = s, sP-quo*s
		tP, t = t, tP-quo*t
	}
	// Sign correction - didn't test if rP being wrong means both of the other two must be corrected, but this guards more edge cases, and it IS possible for t or s to __independently__ be incorrect while the other is correct.
	if 0 > rP {
		rP = -rP
	}
	if (0 > a && 0 < t) || (0 < a && 0 > t) {
		t = -t
	}
	if (0 > b && 0 < s) || (0 < b && 0 > s) {
		s = -s
	}
	// Bezout CoEff sP, tP
	// GCD rP
	// Quo divided by GCD, (t, s) (respectively)
	return t, s, rP
}

// More refs...
// https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication
// https://stackoverflow.com/questions/30017367/lenstras-elliptic-curve-factorization-problems

// DANGER NOTE: Some of the operations (E.G. 3 * pX * pX + A) could overflow the 63 bit limit before the mod; I'm _mostly_ sure that their construction method should keep them within bounds, but some larger N to factor might make them overflow...
// https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law
// https://sites.math.washington.edu/~morrow/336_16/2016papers/thomas.pdf (pg9-10)
// https://math.uchicago.edu/~may/REU2014/REUPapers/Parker.pdf
func ECaddMod(p, q Point3DI64, A, mod int64) Point3DI64 {
	var slopeY, slopeX, x, y, l1C int64
	if 0 == p.Z {
		return q
	}
	if 0 == q.Z {
		return p
	}
	// P and Q are two Points on the Curve ~ https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law ELSE is Tangent
	// Slope (ycombinator) IF P != Q { (Yp - Yq) / (Xp - Xq) } ELSE (P==Q) { (3x*x + a) / (2y) }
	// If x.=x. but yp = -yq then ??? P+Q = 'Infinity' and P contains a factor
	slopeX = (q.X - p.X) % mod
	// Test first, Lenstra.pdf
	// Why test first?  1) a factor might be in the GCD, done.  2) if slopeX == 0 => GCD == mod (N) ELSE it equals 1
	_, _, gcd := ExtendedGCDI64(slopeX, mod)
	if 1 == gcd {
		slopeY = (q.Y - p.Y) % mod
	} else if gcd < mod {
		return Point3DI64{Z: gcd}
	} else {
		if 0 == (p.Y+q.Y)%mod {
			return Point3DI64{X: 0, Y: 1, Z: 0} // Infinity? identity
		}
		slopeY, slopeX = (3*p.X*p.X+A)%mod, (2*p.Y)%mod
	}
	slopeY, slopeX, _ = ExtendedGCDI64(slopeY, slopeX)

	// DANGER: FIXME: the point calculation is wrong, and it's ~ 3am ; this isn't in a production code segment.

	// These are correct for pen and paper... or Rational Numbers - Are __integer__ forms possible?
	// L1 =: (slopeY / slopeX) * ( x - p.X ) + p.Y
	// L1 R(x3, y3) =:	( (slopeY / slopeX)*(slopeY / slopeX) - p.X - q.X	, (slopeY / slopeX)*(x3 - p.X) + p.Y )
	// L2 R(x3, -y3) =:	( (slopeY / slopeX)*(slopeY / slopeX) - p.X - q.X	, (slopeY / slopeX)*(2*p.X + q.X - (slopeY / slopeX)*(slopeY / slopeX))   -  p.Y )
	// x = ((slopeY/slopeX)*(slopeY/slopeX) - p.X - q.X) % mod
	// x = ((slopeY * slopeY)/(slopeX * slopeX) - p.X - q.X) % mod
	// x = ((sYY)/(sXX) - p.X - q.X) % mod
	// x = ((sYY)/(sXX) - p.X*(sXX)/(sXX) - q.X*(sXX)/(sXX)) % mod
	// x = (sYY - p.X*(sXX) - q.X*(sXX))/(sXX) % mod
	//
	// Start from scratch...
	//
	// y = dY/dX * x + Ip
	// Ip = dY/dX * p.X - p.Y // Does not HAVE to be an integer
	// https://en.wikipedia.org/wiki/Algebraic_curve#Intersection_with_a_line
	// dY*x - dX*y + C = 0
	// C = dX*y - dY*x
	l1C = slopeX*p.Y - slopeY*p.X
	_ = l1C
	// https://en.wikipedia.org/wiki/B%C3%A9zout%27s_theorem#Examples_(plane_curves)
	// dX*y = dY*x - C
	// y = (dY*x - C)/dX
	// ###	y*y = x*x*x + A*x + B (mod N)	###
	// (dY*x - l1C)/dX)*((dY*x - l1C)/dX) = x*x*x + A*x + B
	// (dY*x - l1C)*(dY*x - l1C)/(dX*dX) = x*x*x + A*x + B
	// (dY*x - l1C)*(dY*x - l1C) = (dX*dX)*x*x*x + (dX*dX)*A*x + (dX*dX)*B
	// dY*dY*x*x - l1C*dY*x - l1C*dY*x + l1C*l1C = (dX*dX)*x*x*x + (dX*dX)*A*x + (dX*dX)*B
	// dY*dY*x*x - 2*l1C*dY*x + l1C*l1C = (dX*dX)*x*x*x + (dX*dX)*A*x + (dX*dX)*B
	// 0 = dX*dX*x*x*x - dY*dY*x*x + dX*dX*A*x + 2*l1C*dY*x + dX*dX*B - l1C*l1C
	// 0 = dX*dX * x*x*x  -  dY*dY * x*x  +  (dX*dX*A* + 2*l1C*dY) * x  +  dX*dX*B - l1C*l1C
	// So... https://en.wikipedia.org/wiki/B%C3%A9zout%27s_theorem#A_line_and_a_curve
	// 3 points exist where X becomes zero.  (x + pX)*(x + qX)*(x + rX)
	//

	sYY, sXX := slopeY*slopeY, slopeX*slopeX
	x = sYY - p.X*sXX - q.X*sXX
	if 0 != x%sXX {
		fmt.Printf("ECadd, slope is wrong? %d / %d == 0\n", x, sXX)
		panic("Still wrong")
	}
	x /= sXX

	y = ((slopeY/slopeX)*(2*p.X+q.X-(slopeY/slopeX)*(slopeY/slopeX)) - p.Y) % mod
	return Point3DI64{X: x, Y: y, Z: 1}
}

// https://sites.math.washington.edu/~morrow/336_16/2016papers/thomas.pdf (pg11)
// https://math.uchicago.edu/~may/REU2014/REUPapers/Parker.pdf (pg7-8)
func ECmulMod(k, A, mod int64, pow2 Point3DI64) Point3DI64 {
	ret := Point3DI64{0, 1, 0} // the infinity / O point
	// Use P for Powers of 2, prepare the next P at the end of each cycle
	// 'repeated doubling and k's binary expression
	for {
		// If this power of 2 is part of the result, add it
		if 0 < k&1 {
			ret = ECaddMod(ret, pow2, A, mod)
		}
		// every addition should be a valid point, so every GCD return is desired
		if 1 < ret.Z {
			return ret
		}

		k >>= 1
		if 0 == k {
			break
		}

		// Compute the next power of 2
		pow2 = ECaddMod(pow2, pow2, A, mod)
		// every addition should be a valid point, so every GCD return is desired
		if 1 < pow2.Z {
			return pow2
		}
	}
	return ret
}

func LenstraECFI64(n int64) int64 {
	// NOTE: MUST have already tested for /2 and /3 minimum, this version picks BLCM on the PRECONDITION primes in PrimesSmallU8 have all been tested
	var N, A, B, X0, Y0, BLCM, gcd int64
	spMax := len(PrimesSmallU8) - 1
	if 0 < n {
		N = int64(n)
	} else {
		N = -int64(n)
	}

	// TODO more effective curve forms? https://en.wikipedia.org/wiki/Elliptic_curve#Non-Weierstrass_curves

euler_LenstraECFI64_recurse:
	for {
		// https://en.wikipedia.org/wiki/Lenstra_elliptic-curve_factorization#The_algorithm_with_projective_coordinates
		seed := PSRand.RandU32()
		// 1.
		// Strongly random isn't as necessary as just 'not trivial' and 'not always the same'
		// Mod N the inputs, A can probably be above N by one?
		A, X0, Y0 = int64((seed&0x003F_FC00)>>10)%N, int64((seed&0xFFFF_0000)>>16)%N, int64((seed&0xFFFF))%N
		// A is not allowed to be 0
		if 0 == A || 0 == X0 || 0 == Y0 {
			seed = ^seed
			if 0 == A {
				A = int64((seed&0x003F_FC00)>>10) % N
			}
			if 0 == A {
				A++
			}
			if 0 == X0 {
				X0 = int64((seed&0xFFFF_0000)>>16) % N
			}
			if 0 == X0 {
				X0++
			}
			if 0 == Y0 {
				Y0 = int64((seed & 0xFFFF)) % N
			}
			if 0 == Y0 {
				Y0++
			}
		}
		// Note 'choice of B' https://wstein.org/edu/Fall2001/124/lenstra/lenstra.pdf pg 662

		// 2. y*y = x*x*x + A*x + B (mod N)
		B = (Y0*Y0 - X0*X0*X0 - A*X0) % N

		// math.SE Step 4
		// https://math.stackexchange.com/questions/859116/lenstras-elliptic-curve-algorithm
		// for the Y^2 = X^3... curve (derivative?)
		// check gcd(4 * A*A*A + 27 * B * B , N ) == 1 (OK) || N (Bad, New A, recheck) || between 1 and N => YAY (maybe composite)Factor found

		// Test that the curve is 'square free' in X? https://en.wikipedia.org/wiki/Elliptic_curve
		// Discriminant https://en.wikipedia.org/wiki/Discriminant
		// """ the quantity which appears under the square root in the quadratic formula. If a ≠ 0 , {\displaystyle a\neq 0,} this discriminant is zero if and only if the polynomial has a double root. In the case of real coefficients, it is positive if the polynomial has two distinct real roots, and negative if it has two distinct complex conjugate roots.[1] Similarly, the discriminant of a cubic polynomial is zero if and only if the polynomial has a multiple root. In the case of a cubic with real coefficients, the discriminant is positive if the polynomial has three distinct real roots, and negative if it has one real root and two distinct complex conjugate roots.
		// tests properties of roots of the curve
		Discr := 4*A*A*A + 27*B*B
		if 0 == Discr {
			// Bad roll, try again
			continue euler_LenstraECFI64_recurse
		}
		gcd = GCDbin(Discr, N)
		if gcd == N {
			// Bad roll, try again
			continue euler_LenstraECFI64_recurse
		}
		if 1 < gcd {
			// Lucky roll, found a GCD between 1 and N
			return gcd
		}

		// https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law
		// Since the curve is symmetric about the x-axis, given any point P, we can take −P to be the point opposite it. We then have − O = O {\displaystyle -O=O}, as O {\displaystyle O} lies on the XZ-plane, so that − O {\displaystyle -O} is also the symmetrical of O {\displaystyle O} about the origin, and thus represents the same projective point.
		// FIXME: is P(x0,y0) -> Q(x0, - y0) given this symmetry?
		// Z is 1, multiplicative identity
		// P, Q := Point3DI64{X: x0, Y: y0, Z: 1}, Point3DI64{X: x0, Y: -y0, Z: 1} // is it (-y0) % N to force it to wrap to some positive number?
		P := Point3DI64{X: X0, Y: Y0, Z: 1}

		// 3. + 4. ???
		// math.SE Step 5 (BLCM == k)
		//	Choose B -- All Prime Factors must be Less than or Equal to B ;
		BLCMtarget := N / int64(PrimesSmallU8[spMax])
		BLCM = 30 // 2 * 3 * 5
		ii := 3
		for BLCM < BLCMtarget {
			// https://math.mit.edu/research/highschool/primes/materials/2018/conf/7-2%20Rhee.pdf
			// I think that's slide 15, with Lenstra's Factorization Method and steps 1 .. 11
			// This is Step 5, but 6..9 are built within this loop
			// Compute ? Q using d * P (mod n) 'and set P = Q'
			P = ECmulMod(BLCM, A, N, P)
			if 1 < P.Z && P.Z < N {
				return P.Z
			}
			ii++
			if ii > spMax {
				ii = 0
			}
			BLCM *= int64(PrimesSmallU8[ii])

		}
		// 5.

		// math.SE Step 6
		// Compute k*P = ( (Ak/D^2vk , Bk /d^3vk ) ... Multiply the first X,Y point P?  FIXME: Wikipedia lookup

		// math.SE Step 7
		// Calc D = gcd(Dk, n) Same general rules as Step 5 above, 1<D<N => D is a non-trivial factor, else Raise K or roll again.

		// P and Q are two Points on the Curve ~ https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law ELSE is Tangent
		// Slope (ycombinator) IF P != Q { (Yp - Yq) / (Xp - Xq) } ELSE (P==Q) { (3x*x + a) / (2y) }
		// If x.=x. but yp = -yq then ??? P+Q = 'Infinity' and P contains a factor

		// Curve times the LCM?  lcm(...) * P (mod N)  Valid point, or a Factor found
	}
}

// By Euler 0058 it was clear that a better primality test is REQUIRED
// https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test

// https://en.wikipedia.org/wiki/Primality_test#Miller%E2%80%93Rabin_and_Solovay%E2%80%93Strassen_primality_test
// https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test
// == MR + LL
// https://en.wikipedia.org/wiki/Lucas_pseudoprime

/* Go 1.8+ (at least) has https://pkg.go.dev/math/big@go1.22.6#Int.ProbablyPrime
   It might be useful for uint64 but for trivial code just use the library that already does the right test!
func (p *BVPrimes) PrimalityTestBig(qt ???) ??? {
	// 2 > n
	if 1 == two.Cmp(n) { return n }
	// Returns N if probably prime, >=2 as a likely prime factor if composite, 1 on ??? and 0 on error
	//	https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test

	//	== 1 ==
	//	They 'wheel factorize filter' for some small list of N
	//	https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Testing_against_small_sets_of_bases
	//	MR can use rounds of 'base' prime up to 37 for extremely high confidence in numbers < 2^64
	//	With 41 showing extremely high confidence.  2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41

	smallPrimes := [...]uint8{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41}
	smallPrimesLen := len(smallPrimes)
	for ii := 0 ; ii < smallPrimesLen ; ii++ {
		if 0 == qt % smallPrimes[ii] { return smallPrimes[ii] }
	}

	// == 2 ==
	// A is the base BPSW test only uses MR with Base 2 mode anyway.
	// MR with base 2

	// Miller–Rabin primality test
	// Given an integer n, choose some positive integer a < n. Let 2^(s)*d = n − 1, where d is odd.

	// FIRST: find: 2^(s)*d = n − 1 ; where d is odd.
	// 2 ^ s  * d = (n - 1)
	// D must be an odd Integer (whole number), 2 ^ s clearly <= (n - 1)
	// d = (n - 1) / 2 ^ s

	// If BOTH
	// a^d ≢ ± 1 ( mod n )
	// AND
	// a^(d*2^r) ≢ − 1 ( mod n )
	// {\displaystyle a^{2^{r}d}\not \equiv -1{\pmod {n}}}
	// FOR ALL
	// 0 ≤ r ≤ s − 1
	// THEN (if it's true) n is a witness that N is composite
	// ELSE N might or might not be prime

	// If I need to write my own powmul / pow+mul https://www.khanacademy.org/computing/computer-science/cryptography/modarithmetic/a/fast-modular-exponentiation
	// It uses fast finding A^B mod C where B is a power of 2 by reducing it to A^1 mod C => A^2 mod C = (A^1 mod C) * (A^1 mod C) mod C  Then extends to non Pow2 B using B's binary to make powers of 2.
	x =


	// == 3 ==
	// NOTE: Fails when Num is a perfect square, past a limit perform such a check.
	// https://en.wikipedia.org/wiki/Jacobi_symbol
	// D = Primes: start at 5, every other prime is negative.
	// n == the Number
	// Find (first): -1 = Jacobi(D/num)
	// https://en.wikipedia.org/wiki/Jacobi_symbol#Implementation_in_Lua

Based on this:  P = 1 ; Q = (1 − D) / 4

== 4 == (draw the owl)
Perform a STRONG https://en.wikipedia.org/wiki/Lucas_pseudoprime#Strong_Lucas_pseudoprimes
}
*/

// Deprecated function supported by shim interface to Primes
func GetPrimes(primes *[]int, num int) *[]int {
	return Primes.GetPrimesInt(primes, num)
}

// Deprecated function supported by shim interface to Primes
func Factor(primes *[]int, num int) *[]int {
	fp := Primes.Factorize(uint64(num))
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
	Base  uint32
	Power uint32
}

type Factorized struct {
	// Euler 29 wants a list of unique numbers up to 100**100 (100^100) ...
	// Factorized graduates from a []int type number to a structured number, and also stores the effective lengths ahead of time.
	// I'd like to make a version something like lenbase uint8 ; lenpow uint24 but the latter doesn't exist and the []uint16 (still worth it for data size in cache lines) is about to utilize abus-width int and pointer anyway...
	Lenbase uint32
	Lenpow  uint32
	Fact    []Factorpair
}

//func NewFactorized(primes *[]int, n uint) *Factorized {
//	return &Factorized{}
//}

func (facts *Factorized) Mul(fin *Factorized) *Factorized {
	temp := make([]Factorpair, 0, facts.Lenbase+fin.Lenbase)
	fbuf := (*FactorpairQueue)(&temp)
	heap.Init(fbuf)
	var fr, fi uint32
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
	power := uint32(0)
	for 0 < fbuf.Len() {
		fp := heap.Pop(fbuf).(Factorpair)
		power += uint32(fp.Power)
		leak = append(leak, fp)
	}
	// fmt.Printf("Mul base check: %d\n", leak[0].Base)
	if 0 == leak[0].Base {
		facts.Lenbase = 0
		facts.Lenpow = 0
		power = 1
		leak = append(make([]Factorpair, 0, 1), Factorpair{Base: 0, Power: 1})
	}
	facts.Lenbase, facts.Lenpow, facts.Fact = uint32(len(leak)), power, leak
	return facts
}

func (fl Factorized) Eq(fr *Factorized) bool {
	// I already wrote this and it's good to see that reflect.DeepEqual() is notably slower than a directed codepath. https://stackoverflow.com/a/15312182
	if fl.Lenbase != fr.Lenbase || fl.Lenpow != fr.Lenpow {
		return false
	}
	for ii := uint32(0); ii < fl.Lenbase; ii++ {
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
	for ii := uint32(0); ii < fl.Lenbase; ii++ {
		base := big.NewInt(int64(fl.Fact[ii].Base))
		for ee := uint32(0); ee < fl.Fact[ii].Power; ee++ {
			ret = ret.Mul(ret, base) // math.Pow(x, y float64) float64 {...}
		}
	}
	return ret
}

func (fact *Factorized) Uint64() uint64 {
	ret := uint64(1)
	for ii := uint32(0); ii < fact.Lenbase; ii++ {
		for ee := uint32(0); ee < fact.Fact[ii].Power; ee++ {
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
func (fact *Factorized) ExtractPower() (*Factorized, uint32) {
	// Simplify the code at a small memory cost
	if 0 == fact.Lenbase {
		return fact.Copy(), 0
	}
	powbuf := make([]uint32, 0, fact.Lenbase+1)
	for ii := uint32(0); ii < fact.Lenbase; ii++ {
		powbuf = append(powbuf, uint32(fact.Fact[ii].Power))
	}

	for terms := fact.Lenbase; 1 < terms; terms >>= 1 {
		// fmt.Printf("ExtractPower() Round: %v\n", powbuf)
		var ii uint32
		for ii = 0; ii+1 < terms; ii += 2 {
			powbuf[ii>>1] = GCDbin(powbuf[ii], powbuf[ii+1])
			if 1 == powbuf[ii>>1] {
				return fact.Copy(), 1
			}
		}
		if ii+1 == terms {
			if 1 == powbuf[ii] {
				return fact.Copy(), 1
			}
			powbuf[ii>>1] = powbuf[ii]
		}
		terms++
		terms >>= 1 // flooring binary division
	}
	// fmt.Printf("ExtractPower() Final: %v\n", powbuf)
	if 1 == powbuf[0] {
		return fact.Copy(), 1
	}
	iiLim := len(fact.Fact)
	ret := &Factorized{Lenbase: fact.Lenbase, Lenpow: fact.Lenpow / powbuf[0], Fact: make([]Factorpair, iiLim)}
	for ii := 0; ii < iiLim; ii++ {
		ret.Fact[ii].Base = fact.Fact[ii].Base
		ret.Fact[ii].Power = fact.Fact[ii].Power / uint32(powbuf[0])
	}
	return ret, powbuf[0]
}

func (fact *Factorized) Pow(p uint32) *Factorized {
	// Simplify the code at a small memory cost
	if 0 == p {
		return &Factorized{Lenbase: 1, Lenpow: 1, Fact: append(make([]Factorpair, 0, 1), Factorpair{Base: 2, Power: 0})}
	}
	iiLim := len(fact.Fact)
	ret := &Factorized{Lenbase: fact.Lenbase, Lenpow: fact.Lenpow * p, Fact: make([]Factorpair, iiLim)}
	for ii := 0; ii < iiLim; ii++ {
		ret.Fact[ii].Base = fact.Fact[ii].Base
		ret.Fact[ii].Power = fact.Fact[ii].Power * uint32(p)
	}
	return ret
}

func (fact *Factorized) PowDivMul(num, den uint32) *Factorized {
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
		ret.Fact[ii].Power = (fact.Fact[ii].Power / uint32(den)) * uint32(num)
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
	if uint32(flen) != f.Lenbase {
		fmt.Printf("ERROR ProperDivisors(): Lenbase != len(Fact): %v\n", f)
	}
	sf := make([]uint32, 0, f.Lenpow)
	for ii := uint32(0); ii < f.Lenbase; ii++ {
		for pp := uint32(0); pp < f.Fact[ii].Power; pp++ {
			sf = append(sf, uint32(f.Fact[ii].Base))
		}
	}
	var limit uint64
	if 64 == f.Lenpow {
		limit ^= 1
	} else {
		limit = (uint64(1) << f.Lenpow) - 1
	}

	almost := uint64(1)
	for ff := uint32(1); ff < f.Lenpow; ff++ {
		almost *= uint64(sf[ff])
	}
	bitVec := bitvector.NewBitVector(almost)
	bitVec.Set(1)      // All 0s
	bitVec.Set(almost) // ^1 // ~1
	for ii := uint64(1); ii < limit; ii++ {
		bit := uint64(1)
		ar := uint64(1)
		for ff := uint32(0); ff < f.Lenpow; ff++ {
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

func EulerTotientPhi(n uint64) uint64 {
	if 1 == n {
		return 1
	}
	return Primes.Factorize(n).EulerTotientPhi()
}

func (f *Factorized) EulerTotientPhi() uint64 {
	// []FactorizedN -> b0^(p0-1) * (b0-1) * ... * bn^(pn-1) * (bn-1)
	var ret uint64
	ret = 1
	iiLim := len(f.Fact)
	for ii := 0; ii < iiLim; ii++ {
		// 'subtract' one via < (less than)
		for pow := uint32(1); pow < f.Fact[ii].Power; pow++ {
			ret *= uint64(f.Fact[ii].Base)
		}
		ret *= uint64(f.Fact[ii].Base) - 1
	}
	return ret
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
