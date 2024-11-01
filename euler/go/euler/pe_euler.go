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
	PrimesSmallU8Mx             = 54 - 1 // Not quite a whole cache line (64 bytes) -- If more are needed consider a not-compressed bitvector rather than the primes bitvector (this was added as that seemed too costly for a hot path on initial factorization)
	PrimesSmallU8MxVal          = 251
	PrimesSmallU8MxValCompAfter = 256
	PrimesSmallU8MxValPow2After = 0x1_0000 // 256 * 256
	Debug128bitInts             = false
)

var (
	Primes        *BVPrimes
	PrimesSmallU8 [PrimesSmallU8Mx + 1]uint8 // {2, 3, 5, 7, ...}
	PrimesMoreU16 []uint16
	PrimesMoreU32 []uint32 // I do see where 65535 could be insufficient
	// PrimesMoreU64 []uint64 // I don't foresee any algorithms that need a LIST of continuous known primes greater than 32billion; that's clearly the range of other tools...
	PSRand *PSRandPGC32
)

func init() {
	Primes = NewBVPrimes()
	PrimesSmallU8 = [PrimesSmallU8Mx + 1]uint8{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251} // 41 required for reasons, 53 nice for 16 total numbers; 16 bytes of memory, 1/4th cache line
	PrimesMoreU16 = []uint16{}
	PrimesMoreU32 = []uint32{}
	PSRand = NewPSRandPGC32(0x4d595df4d0f33173, 1442695040888963407) // Seed could be _anything_, inc must be any odd constant values from https://en.wikipedia.org/wiki/Permuted_congruential_generator#Example_code
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

/*
	A challenge to help me better understand a method.
	It's obvious if considering E.G. u8 vs u16 how 256*256 is 65536 or rather, how 255*255 can't overflow 65535.
	If the number is small take the fast path and do things with native ints that fit.
	Otherwise, 'make it fit'.

	E.G. Consider the toy example of nibbles, 4 bit values.  Their values can't overflow an octet / byte, but they are good toy examples.
	5 * 4 = 101 * 100  Oh no, it overflows the nibble as 10100

	An idea I came across in search results; split the value up so it doesn't overflow.  Though I didn't look at enough detail to know exactly how.
	5 * 4	=> 101 * 100
	01 01 * 01 00
	h*h		l*l
	01 (0000)	00	|= 010000  That got the high bits correct, but it's missing 100  The patterns don't make it clear enough where that should come from though...

	01 11 * 10 00 (7*8 = 56 = 111000)
	10 (0000)	00	|= 100000  H0*H1 + L0*L1 missing H0*L1 and H1*L0 ?
	10 (0000) + 0110 (00) + 0000 (00) + 00	|=	111000

	Also, if the two numbers are known to be the same number, the middle term is doubled (<<1); though that seems unlikely to matter for long when a modulus is involved.

	That makes this much easier to follow.

	However, there's better:
	https://en.wikipedia.org/wiki/Karatsuba_algorithm
	https://en.wikipedia.org/wiki/Toom%E2%80%93Cook_multiplication (Karatsuba aka Toom-2)
	2019~2021 Annals of Mathematics https://en.wikipedia.org/wiki/Multiplication_algorithm#Further_improvements

	Modulus though...
	https://en.wikipedia.org/wiki/Division_algorithm#Integer_division_(unsigned)_with_remainder
	https://en.wikipedia.org/wiki/Binary_number#Division
	Maybe good for education of those new to math https://en.wikipedia.org/wiki/Chunking_(division) https://en.wikipedia.org/wiki/Long_division
	Considering these
	https://en.wikipedia.org/wiki/Fourier_division
	** https://en.wikipedia.org/wiki/Division_algorithm#Newton%E2%80%93Raphson_division
	https://en.wikipedia.org/wiki/Short_division#Modulo_division << Appealing at a glance as it's direct BUT, this would require operation across H+L as one word.
	Interesting
	https://en.wikipedia.org/wiki/Division_algorithm#Goldschmidt_division
	https://en.wikipedia.org/wiki/Division_algorithm#SRT  F00F




*/

func Mod1toN[INT ~uint | ~int | ~uint64 | ~int64](a, mod INT) INT {
	if 0 == a {
		return mod
	}
	return a % mod
}

func AddOF64Mod[INT ~uint | ~int | ~uint64 | ~int64](a, b, mod INT) INT {
	// func UU64AddUU64(ah, al, bh, bl uint64) (uint64, uint64, uint64) {
	var minusA, minusB bool
	var Al, Bl, Rh, Rl uint64
	minusA, minusB = false, false
	if 0 > a {
		Al = uint64(-a)
		minusA = true
	} else {
		Al = uint64(a)
	}
	if 0 > b {
		Bl = uint64(-b)
		minusB = true
	} else {
		Bl = uint64(b)
	}
	if minusA == minusB {
		_, Rh, Rl = UU64AddUU64(0, Al, 0, Bl)
	} else if minusB {
		Rh, Rl = UU64SubUU64(0, Al, 0, Bl)
	} else {
		Rh, Rl = UU64SubUU64(0, Bl, 0, Al) // minusA
	}
	minusA = (0 < Rh&(uint64(1)<<63))
	if 0 == Rh || 0 == ^Rh {
		if minusA {
			return -INT(^(Rl % uint64(mod)))
		}
		return INT(Rl % uint64(mod))
	}
	if minusA {
		return -INT(UU64Mod(Al, Bl, uint64(mod)))
	}
	return INT(UU64Mod(Al, Bl, uint64(mod)))
}

func SubOF64Mod[INT ~uint | ~int | ~uint64 | ~int64](a, b, mod INT) INT {
	// func UU64SubUU64(ah, al, bh, bl uint64) (uint64, uint64) {
	var minusA, minusB bool
	var Al, Bl, Rh, Rl uint64
	minusA, minusB = false, true
	if 0 > a {
		Al = uint64(-a)
		minusA = true
	} else {
		Al = uint64(a)
	}
	if 0 > b {
		Bl = uint64(-b)
		minusB = false
	} else {
		Bl = uint64(b)
	}
	if minusA == minusB {
		_, Rh, Rl = UU64AddUU64(0, Al, 0, Bl)
	} else if minusB {
		Rh, Rl = UU64SubUU64(0, Al, 0, Bl)
	} else {
		Rh, Rl = UU64SubUU64(0, Bl, 0, Al) // minusA
	}
	minusA = (0 < Rh&(uint64(1)<<63))
	if 0 == Rh || 0 == ^Rh {
		if minusA {
			return -INT(^(Rl % uint64(mod)))
		}
		return INT(Rl % uint64(mod))
	}
	if minusA {
		return -INT(UU64Mod(Al, Bl, uint64(mod)))
	}
	return INT(UU64Mod(Al, Bl, uint64(mod)))
}

func MulOF64Mod[INT ~uint | ~int | ~uint64 | ~int64](a, b, mod INT) INT {
	var minus bool
	var Al, Bl uint64
	minus = false
	if 0 > a {
		Al = uint64(-a)
		minus = !minus
	} else {
		Al = uint64(a)
	}
	if 0 > b {
		Bl = uint64(-b)
		minus = !minus
	} else {
		Bl = uint64(b)
	}
	// This is _notably_ slower
	// pow2a := BitsPowerOfTwo(Al)
	// pow2b := BitsPowerOfTwo(Bl)
	bits := BitsLeadingZeros64(Al) + BitsLeadingZeros64(Bl)
	if bits >= 64 {
		// a and b fit the fast path this time!!!
		// This misses _some_ cases that fit, but a full check should cost about as much as taking the slow path
		// E.G. 0x7F * 0x02 1 bit + 6 bits == 7 bit fits, but 0x7F * 0x03 does not fit 8 bits.  F*F fits (E1) as would F*0x11 , but not 0x12
		// if 0 <= pow2a {
		//	Al = Bl << pow2a
		// } else if 0 <= pow2b {
		//	Al = Al << pow2b
		// } else {
		// }
		Al *= Bl
		if minus {
			return -INT(Al % uint64(mod))
		}
		return INT(Al % uint64(mod))
	}

	// if 0 <= pow2a {
	//	Al, Bl = Bl>>(64-pow2a), Bl<<pow2a
	//} else if 0 <= pow2b {
	//	Al, Bl = Al>>(64-pow2b), Al<<pow2b
	//} else {
	//
	//}
	Al, Bl = UU64Mul(Al, Bl)

	if minus {
		return -INT(UU64Mod(Al, Bl, uint64(mod)))
	}
	return INT(UU64Mod(Al, Bl, uint64(mod)))
}

func UU64MulWrap(a, b uint64) (uint64, uint64) {
	bits := BitsLeadingZeros64(a) + BitsLeadingZeros64(b)
	if bits >= 64 {
		return 0, a * b
	}
	return UU64Mul(a, b)
}

func UU64Mul(a, b uint64) (uint64, uint64) {
	// https://en.wikipedia.org/wiki/Karatsuba_algorithm
	// https://en.wikipedia.org/wiki/Toom%E2%80%93Cook_multiplication (Karatsuba AKA Toom-2)
	var z1minus bool
	var Ah, Al, Bh, Bl, Z2, Z0 uint64
	Al = uint64(a)
	Bl = uint64(b)
	// This is _notably_ slower
	// pow2a := BitsPowerOfTwo(Al)
	// pow2b := BitsPowerOfTwo(Bl)
	// if BitsLeadingZeros64(Al)+BitsLeadingZeros64(Bl) >= 64 {
	//	 if 0 <= pow2a {
	//		Al = Bl << pow2a
	//	} else if 0 <= pow2b {
	//		Al = Al << pow2b
	//	} else {
	//		Al *= Bl
	//	}
	//	return 0, Al
	//}
	// if 0 <= pow2a {
	//	Al, Bl = Bl>>(64-pow2a), Bl<<pow2a
	//	return Al, Bl
	//}
	//if 0 <= pow2b {
	//	Al, Bl = Al>>(64-pow2b), Al<<pow2b
	//	return Al, Bl
	//}

	// split @ fixed 32 bits
	Ah, Al, Bh, Bl = Al>>32, Al&0xFFFF_FFFF, Bl>>32, Bl&0xFFFF_FFFF
	Z2, Z0 = Ah*Bh, Al*Bl
	z1minus = false
	if Al < Ah {
		Al = Ah - Al
		z1minus = !z1minus
	} else {
		Al = Al - Ah
	}
	if Bh < Bl {
		Bl = Bl - Bh
		z1minus = !z1minus
	} else {
		Bl = Bh - Bl
	}
	Ah = Al * Bl
	// if 0x7ffffffffffffe79 == b {
	//	fmt.Printf("UU64Mul: Z1: %#x =? %x * %x (z1-%t)\n", Ah, Al, Bl, z1minus)
	// 0x9900A28DD20B2B8
	// 0x9900a28dd20b2b8
	//}
	// Every step of Z1 might overflow, in assembly the carry flags could help, but higher level it's just less errant to use the wider functions
	_, Bh, Bl = UU64AddUU64(0, Z2, 0, Z0) // Z1
	if z1minus {
		Bh, Bl = UU64SubUU64(Bh, Bl, 0, Ah)
	} else {
		_, Bh, Bl = UU64AddUU64(Bh, Bl, 0, Ah)
	}
	Ah, Al = Bh<<32|Bl>>32, (Bl&0xFFFF_FFFF)<<32
	// if 0xe4bc2b74f0572bbb == b {
	//	fmt.Printf("UU64Mul: debug: %16x + %x .. %16x + %x (z1-%t)\n", Z2, Ah, Z0, Al, z1minus)
	// }
	Z0 += Al
	if Z0 < Al {
		Ah++ // addition overflowed, add it to the high half
	}
	Z2 += Ah

	//
	// Mul complete... 128 bit split word Z2, Z0
	//
	if Debug128bitInts {
		b1, b0, bd := big.NewInt(1), big.NewInt(0), big.NewInt(0)
		b0.SetUint64(0xFFFF_FFFF_FFFF_FFFF)
		b1.Add(b1, b0)
		b0.SetUint64(a)
		bd.SetUint64(b)
		b0.Mul(b0, bd)
		b0.DivMod(b0, b1, b1)
		if Z2 != b0.Uint64() || Z0 != b1.Uint64() {
			fmt.Printf("UU64Mul: Output does not match (%#x, %#x) should be %#x %x got %#x %x\n", a, b, b0.Uint64(), b1.Uint64(), Z2, Z0)
			panic("UU64Mul")
		}
	}

	return Z2, Z0
}

/*
func UU64MulUUUU(ah, al, bh, bl uint64) (uint64, uint64, uint64, uint64) {
	// https://en.wikipedia.org/wiki/Toom%E2%80%93Cook_multiplication (Toom-4)
	// natural math https://gmplib.org/manual/Toom-4_002dWay-Multiplication
	// X is split across ah,al Y is split across bh,bl
	//          * I'm worried about three lines, all of which bit-shift terms to the left.  Those could overflow my selected split locations.
	//		While I could farm those back out to Karatsuba, that sort of defeats the point of trying to support a 128bit * 128bit op; GMP's documentation suggests Toom-6+1/2 is the next worthwhile step higher.
	// . t=0	x0 * y0 == w[0]
	// 0 t=1/2  *	(x3 + 2*x2 + 4*x1 + 8*x0) * (y3 + 2*y2 + 4*y1 + 8*y0)
	// 1 t=-1/2 *	( -x3 + 2*x2 - 4*x1 + 8*x0) * (-y3 + 2*y2 - 4*y1 + 8*y0)
	// 2 t=1	(x3 + x2 + x1 + x0) * (y3 + y2 + y1 + y0)
	// 3 t=-1	( -x3 + x2 - x1 + x0) * ( -y3 + y2 - y1 + y0)
	// 4 t=2    *	(8*x3 + 4*x2 + 2*x1 + x0) * (8*y3 + 4*y2 + 2*y1 + y0)
	// . t=inf	x3 * y3 == w[inf]
	x, y := [4]uint64{al&0xFFFF_FFFF, al>>32, ah&0xFFFF_FFFF, ah>>32},[4]uint64{bl&0xFFFF_FFFF, bl>>32, bh&0xFFFF_FFFF, bh>>32}
	var w [7]uint64
	var t [5]uint64
	var ti1minus, ti3minus bool
	w[0], w[6] = x[0] * y[0], x[3] * y[3]  // t0 & tInf
	t[4] = (x[3]<<3 + x[2]<<2 + x[1]<<1 + x[0]) * (y[3]<<3 + y[2]<<2 + y[1]<<1 + y[0])
	// uint math is mod (pow2 - 1) ; -1 wraps to all 1s and adding one overflows back around to 0.  It's safe to just add and subtract as directed then correct the result (0-negativevalue)
	a, b := x[0] + x[2] - x[1] - x[3], y[0] + y[2] - y[1] - y[3]
	if 0 < a&(1<<63) {
		ti3minus = !ti3minus
		a = 0 - a
	}
	if 0 < b&(1<<63) {
		ti3minus = !ti3minus
		b = 0 - b
	}
	t[3], t[2] = a * b, (x[3] + x[2] + x[1] + x[0]) * (y[3] + y[2] + y[1] + y[0])
	a, b = x[0]<<3 + x[2]<<1 - x[1]<<2 - x[3], y[0]<<3 + y[2]<<1 - y[1]<<2 - y[3]
	if 0 < a&(1<<63) {
		ti1minus = !ti1minus
		a = 0 - a
	}
	if 0 < b&(1<<63) {
		ti1minus = !ti1minus
		b = 0 - b
	}
	t[1], t[0] = a * b,  x[0]<<3 + x[2]<<1 + x[1]<<2 + x[3], y[0]<<3 + y[2]<<1 + y[1]<<2 + y[3]
	//
	// Abandoned - the bitshifted terms would overflow ; What I wanted to use this for also seems unlikely to yield sufficient benefit when it's already this complex.
	//
}
*/

func UU64AddUU64(ah, al, bh, bl uint64) (uint64, uint64, uint64) {
	var carry uint64
	bl = al + bl
	// modulo overflow
	if bl < al {
		carry++
	}
	bh = ah + bh + carry
	carry = 0
	// modulo overflow
	if bh < ah {
		carry++
	}
	return carry, bh, bl
}

// Note, data is passed around as uint64s but underflow can result in negative numbers that are binary equal to int64s
func UU64SubUU64(ah, al, bh, bl uint64) (uint64, uint64) {
	bl = al - bl
	// modulo wrap occurred
	if bl > al {
		ah--
	}
	return ah - bh, bl
}

func UU64Cmp(ah, al, bh, bl uint64) int {
	if ah > bh {
		return 1
	}
	if ah < bh {
		return -1
	}
	if al > bl {
		return 1
	}
	if al < bl {
		return -1
	}
	return 0 // ==
}

/*
.	FIXME: After a good night's rest I happened to realize the current version of UU64DivQD throws away a lot of information.

	The fit estimate vaguely models just subtracting a matched up bitshift of D and then corrects the estimate for multiplication successively...
	However I don't really _need_ that multiply, bit-shifts (across two machine words, so masks too) could generate successive values to subtract / add.
	Mod (2^128)-1 I could also just subtract whichever is the least 2s comp number and adjust the Quotient accordingly.

	It's sort of like NR, and SRT which uses a lookup table and processes a chain of estimates and corrections to reach the result seems like it'd be even slower...
	https://en.wikipedia.org/wiki/Division_algorithm#SRT_division

	https://stackoverflow.com/a/27808065

	Multiplication to support q[n+1]=q[n]*(2^(k+1)-q[n]*B) >> k looks _way_ more expensive than simple bitshifts and subtractions.  I started to embark on Toom-4 in an attempt to support it but realized the bitshifts meant that terms would overflow, just like they did in the addition section of Karatsuba (that solved by 128bit +/- to support possibly multiple carry events).
.
*/

func UU64DivQD(h, l, d uint64) (uint64, uint64, uint64) {
	// https://en.wikipedia.org/wiki/Division_algorithm#Large-integer_methods
	// ~~ Vaguely inspired by https://en.wikipedia.org/wiki/Division_algorithm#Newton%E2%80%93Raphson_division
	var qh, ql, sh, sl, r, th, tl uint64
	// if 0xffffffffffffffc5 == d {
	//	fmt.Printf("UU64DivQD: Invoked with 0x%x, 0x%x, 0x%x\n", h, l, d)
	//}
	if 0 == d {
		fmt.Printf("UU64DivQD: Division by Zero, no error path, returning zero. %x_%16x / %x", h, l, d)
		return 0, 0, 0
	}
	if 1 == d {
		return h, l, 0
	}

	if 0 == h {
		// hardware word fast path, there's no high part
		l, r = l/d, l%d
		return 0, l, r
	}

	// Measured: No, at least not with BitsPowOfTwo as is... NOT worth trying to detect an exact power of 2 divisor

	// btN is the current MSB (1<<btN) of the Numerator 40e6_278a_a8ed_c6c9
	// btQ = btN - btD (the current MSB (1<<btD) of the Denominator)
	btN := 64 - BitsLeadingZeros64(h) + 64
	btQ := btN - 64 + BitsLeadingZeros64(d)
	th, tl = h, l

	// Establish an initial estimate with quick bit muls, at least 7 bits precision will result
	for 0 <= btQ {
		if 64 <= btQ {
			r = d << (btQ - 64)
		} else {
			r = d >> (64 - btQ) // d>>(64-btQ) How would that work if btQ > 64?
		}
		sh, sl = UU64SubUU64(th, tl, r, d<<btQ)
		if 0 == sh&(1<<63) {
			if 64 <= btQ {
				qh |= 1 << (btQ - 64)
			} else {
				ql |= 1 << btQ
			}
			th, tl = sh, sl
		}
		btN, btQ = btN-1, btQ-1
	}
	return qh, ql, tl
}

func UU64DivQD_old(h, l, d uint64) (uint64, uint64, uint64) {
	// https://en.wikipedia.org/wiki/Division_algorithm#Large-integer_methods
	// ~~ Vaguely inspired by https://en.wikipedia.org/wiki/Division_algorithm#Newton%E2%80%93Raphson_division
	var X, low, high, rh, rl, r, th, tl uint64
	// if 0xffffffffffffffc5 == d {
	//	fmt.Printf("UU64DivQD: Invoked with 0x%x, 0x%x, 0x%x\n", h, l, d)
	//}
	if 0 == d {
		fmt.Printf("UU64DivQD: Division by Zero, no error path, returning zero. %x_%16x / %x", h, l, d)
		return 0, 0, 0
	}
	if 1 == d {
		return h, l, 0
	}

	if 0 == h {
		// Fast path, there's no high part
		l, r = l/d, l%d
		return 0, l, r
	}

	// Measured: No, at least not with the BitsPowOfTwo as is... Ponder: is it worth trying to detect an exact power of 2 divisor?
	// pow2b := BitsPowerOfTwo(d)
	//if 0 <= pow2b {
	//	rh, rl, r = h>>pow2b, (h<<(64-pow2b))|(l>>pow2b), l&((1<<pow2b)-1)
	//	return rh, rl, r
	//}

	// h > 0

	// This is _slightly_ faster, 0.05~0.1 seconds in the unit tests (out of 6 seconds) but still slightly.
	bitsN, bitsD := uint64(64-BitsLeadingZeros64(h))+64, uint64(64-BitsLeadingZeros64(d))
	// low, high = 1<<(bitsN-bitsD-1), bitsN-bitsD+1
	// low must be LESS so assume bitsN and bitsD+1 as the (less than) power of 2 value of the target Q
	// high must be MORE so assume bitsN+1 and bitsD as the (more than) power of 2 value of the target Q
	low, high = bitsN-bitsD-1, bitsN-bitsD+1
	if 64 <= high {
		// Need more data and though thought to establish a better relation than this...
		high = ^uint64(0)
	} else {
		high = (1 << high) / d
	}
	low = (1 << low) / d

	// fmt.Printf("debug trace: bN: %d\tbD: %d\tl: %d\th: %d\n", bitsN, bitsD, low, high)

	low, high = 0, ^uint64(0)
	// https://en.wikipedia.org/wiki/Binary_search -- sort of
	for low != high {
		// X= (low + high) / 2 // without overflow
		X = (low >> 1) + (high >> 1) + (low&1)&(high&1)
		th, tl = UU64Mul(X, d)
		rh, rl = UU64SubUU64(h, l, th, tl)
		cmp := UU64Cmp(h, l, th, tl)
		// if 0xffffffffffffffc5 == d && 0xcc5fb7b7273cc415 == h && 0x4bfc032560925a99 == l {
		// fmt.Printf("debug trace: %x : l %16x h %16x\tX: %16x\tsub %16x %16x\n", d, low, high, X, rh, rl)
		//}
		if 0 == rh {
			// fmt.Printf("Debug return rh 0\n")
			if 0 == rl {
				// Winner
				return 0, X, 0
			}
			// Runner Up
			rh, r = rl/d, rl%d
			rh += X
			if rh < X {
				return 1, rh, r
			}
			return 0, rh, r
		} else if 1 == rh && rl < (d-1) {
			// fmt.Printf("debug rh1: %x : x %16x\th: %16x l: %16x\tth %16x\ttl: %16x\tsub %16x %16x\n", d, X, h, l, th, tl, rh, rl)
			// Edge cases like: 0x2 , 0xfffffffffffff8d5 , 0x3fffffffffffff67
			rl += 1 + ((^uint64(0)) % d)
			if rl >= d {
				X++
				rl %= d
			}
			rh = X + ((^uint64(0)) / d) + (((^uint64(0))%d)+1+rl)/d
			if rh < X {
				return 1, rh, rl
			}
			return 0, rh, rl
		} else if -1 == cmp {
			// X*d > h_l
			high = X - 1
		} else {
			// X*d < h_l -- If rh == 1 the remainder may be at this location
			low = X + 1
		}
	}
	th, tl = UU64Mul(low, d)
	rh, rl = UU64SubUU64(h, l, th, tl)
	if 0 == rh && rl < d {
		return 0, low, rl
	}

	fmt.Printf("debug: UU64DivQD: %x : x %16x\th: %16x l: %16x\tth %16x\ttl: %16x\tsub %16x %16x\n", d, low, h, l, th, tl, rh, rl)
	panic("UU64DivQD: 128bit div: This should not be reached.  Expected return path is via binary search resolution.")
	/*
		th, tl = UU64Mul(low, d)
		rh, rl = UU64SubUU64(h, l, th, tl)
		fmt.Printf("debug outer: %x : l %16x h %16x\tX: %16x\tsub %16x %16x\n", d, low, high, low, rh, rl)
		if 0 == rh && rl < d {
			fmt.Printf("Debug return Outer\n")
			return 0, low, rl
		}
		fmt.Printf("debug outer: %x : x %16x\th: %16x l: %16x\tth %16x\ttl: %16x\tsub %16x %16x\n", d, low, h, l, th, tl, rh, rl)
		b1, b0, bd := big.NewInt(1), big.NewInt(0), big.NewInt(0)
		b0.SetUint64(0xFFFF_FFFF_FFFF_FFFF)
		b1.Add(b1, b0)
		b0.SetUint64(h)
		b0.Mul(b0, b1)
		bd.SetUint64(l)
		b0.Add(b0, bd)
		bd.SetUint64(d)
		b0.DivMod(b0, bd, b1)
		fmt.Printf("Big Div: q: 0x%s\tr: 0x%s\n", b0.Text(16), b1.Text(16))
		// 0xcc5fb7b7273cc444
		//   cc5fb7b7cc5fb7b6
		panic("exit")
		// d is _way_ smaller than h_l
		// low has become the low of the quotient
		th, tl, r = UU64DivQD(rh, rl, d)
		high, th, tl = UU64AddUU64(th, tl, 0, low)
		// fmt.Printf("debug unstk: %x : l %16x r %16x\th: %16x\tadd %16x %16x\n", d, low, r, high, rh, rl)
		if 0 != high {
			fmt.Printf("UU64DivQD: New test case, please add and resolve UU64DivQD(0x%x, 0x%x, 0x%x) -- It should not have overflowed, this should be unreachable.\n", h, l, d)
			panic("too big")
		} else {
			fmt.Printf("Debug return Outer 2\n")
			return th, tl, r
		}
		fmt.Printf("debug back failed: %x : x %16x\th: %16x l: %16x\tth %16x\ttl: %16x\n", d, low, h, l, th, tl)
		panic("what now?")

		fmt.Printf("UU64DivQD: New test case, please add and resolve UU64DivQD(0x%x, 0x%x, 0x%x) -- It should have returned during the binary search rather than exited. This should be unreachable.\n", h, l, d)
		// panic("UU64DivQD")
	*/
	return 0, 0, 0
}

func UU64Mod(h, l, mod uint64) uint64 {
	if 0 == h || 0 == ^h {
		return l % mod
	}
	qh, ql, r := UU64DivQD(h, l, mod)
	if Debug128bitInts {
		b1, b0, bd := big.NewInt(1), big.NewInt(0), big.NewInt(0)
		b0.SetUint64(0xFFFF_FFFF_FFFF_FFFF)
		b1.Add(b1, b0)
		b0.SetUint64(h)
		b0.Mul(b0, b1)
		bd.SetUint64(l)
		b0.Add(b0, bd)
		bd.SetUint64(mod)
		b0.DivMod(b0, bd, b1)
		if b1.Uint64() != r {
			panic(fmt.Sprintf("Big Div: %x %x / %x =>\tq: 0x%s\tr: 0x%s\t got: %x %x r %x\n", h, l, mod, b0.Text(16), b1.Text(16), qh, ql, r))
		}
	}
	return r
}

func BitsPowerOfTwo[INT ~uint | ~uint64](n INT) int {
	shift := 63 - BitsLeadingZeros64(n)
	if 0 > shift || 0 != n&^(1<<shift) {
		return -1
	}
	return shift
}

func BitsLeadingZeros64[INT ~uint | ~int | ~uint64 | ~int64](n INT) int {
	var test uint64
	var zeros, shift int
	shift = 64 // bits
	// var negative bool // if it's negative, or even just a signed int, the leading sign bit should be assumed by the collar
	// math/bits doesn't have Len functions for signed integers
	if 0 > n {
		test = uint64(-(n + 1))
	} else {
		test = uint64(n)
	}
	if 0 == n {
		return shift
	}
	// math/bits uses a look-up table for the last byte, but I'll trust in the branch predictor, cache, and tighter code; or if I'm wrong the compiler to be smart enough to unroll this 6 times
	for 1 < shift {
		shift >>= 1
		if (uint64(1) << shift) <= test {
			test >>= shift
		} else {
			zeros += shift
		}
		// fmt.Printf("BLZ: %d\ts: %d\tz: %d\tt: %d\n", n, shift, zeros, test)
	}
	return zeros
}

// Compatibility interface
func PowU64(n, pow uint64) uint64 {
	return PowInt(n, pow)
}

func PowInt[INT ~uint | ~int | ~uint64 | ~int64](n, pow INT) INT {
	if 0 >= pow {
		return 1
	}
	if 1 == pow {
		return n
	}
	//
	// PowInt does not guard against overflow, the answer won't fit anyway in those cases; and implicitly *mod* (U64Limit) this is the only possible answer.
	//
	if 2 == n {
		return INT(1) << pow
	}

	//  (re?)learned of better ways https://en.wikipedia.org/wiki/Modular_exponentiation#Right-to-left_binary_method
	// MSB (set) to LSB ; each step will be either np*np OR np*np + n
	// 100 == n*n, n*n * n*n ; 101 == n*n, n*n n*n * n
	var np INT
	np = n // Always one is the first N ;; -1 to move from leading zeros to first one, then another -1 to move over the first 1 (MSB) into the first variable point.
	for mask := INT(1) << (64 - BitsLeadingZeros64(pow) - 2); 0 < mask; mask >>= 1 {
		// fmt.Printf("%05b\n%05b\n\n", pow, mask)
		np *= np // double
		if 0 < pow&mask {
			np *= n // times n
			//} else {
		}
	}
	return np
}

func PowIntMod[INT ~uint | ~uint64 | ~int | ~int64](n, pow, mod INT) INT {
	if 1 >= mod {
		if 0 == mod {
			fmt.Printf("PowIntMod(%d, %d, %d) incorrect arguments, modulus cannot be <1", n, pow, mod)
			panic("Modulus cannot be < 1")
		}
		return 0 // Zero should be an incorrect mod...
	}
	if 0 >= pow {
		return 1
	}

	n %= mod
	if 1 == pow {
		return n
	}
	if 2 == n && 64 > pow {
		return (INT(1) << pow) % mod
	}
	var np, n64, pow64, mod64, mask uint64
	// NP == Always one is the first N ;; -1 to move from leading zeros to first one, then another -1 to move over the first 1 (MSB) into the first variable point.
	np, n64, pow64, mod64 = uint64(n), uint64(n), uint64(pow), uint64(mod)

	// Test for overflow of (mod-1) * (mod-1) and divert to 'overflow' path.
	if BitsLeadingZeros64((mod64 - 1)) >= 32 {
		// No overflow, fast path

		// MSB (set) to LSB ; each step will be either np*np OR np*np + n
		// 100 == n*n, n*n * n*n ; 101 == n*n, n*n n*n * n
		for mask = 1 << (64 - BitsLeadingZeros64(pow64) - 2); 0 < mask; mask >>= 1 {
			// fmt.Printf("%05b\n%05b\n\n", pow, mask)
			np = (np * np) % mod64 // double
			if 0 < pow64&mask {
				np = (np * n64) % mod64 // times n
				//} else {
			}
		}
	} else {
		// Overflows possible slow path
		for mask = 1 << (64 - BitsLeadingZeros64(pow64) - 2); 0 < mask; mask >>= 1 {
			// fmt.Printf("%05b\n%05b\n\n", pow, mask)
			np = MulOF64Mod(np, np, mod64) // double
			if 0 < pow64&mask {
				np = MulOF64Mod(np, n64, mod64) // times n
				//} else {
			}
		}
	}
	return INT(np)
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
	// This isn't worth the payoff.
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

func (r *PSRandPGC32) RandIntU64() uint64 {
	return uint64(r.RandU32()<<32) | uint64(r.RandU32())
}

func (r *PSRandPGC32) RandInt(max uint64) uint64 {
	if max < 0x1_0000_0000 {
		return uint64(r.RandU32()) % max
	}
	return (uint64(r.RandU32()<<32) | uint64(r.RandU32())) % max
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

func (p *BVPrimes) PrimeGlobalList(pmax uint64) {
	var pg, pgMx, pidx, pidxMx, bb uint32
	var prime uint64
	p.Grow(pmax)
	PrimesMoreU16 = make([]uint16, 0, BVl1/2)
	PrimesMoreU32 = make([]uint32, 0, BVl1/4)
	prime = PrimesSmallU8MxVal
	if pmax <= prime {
		return
	}
	pidx = uint32(PrimesSmallU8MxVal) + 2 - 3 // +2 for the prime _after_ the current prime...
	bb = (pidx & BVprimeByteBitMask) >> 1
	pidx >>= BVprimeByteBitShift
	pg, pidx = pidx/BVpagesize, pidx%BVpagesize
	pidxMx = (uint32(p.Last) - 3) >> BVprimeByteBitShift
	pgMx, pidxMx = pidxMx/BVpagesize, pidxMx%BVpagesize
BVPrimes_PrimeGlobalList_U16primes:
	for pg <= pgMx {
		// pidx
		for (pg < pgMx && pidx < BVpagesize) || (pg == pgMx && pidx <= pidxMx) {
			// bb
			for ; bb < BVbitsPerByte; bb++ {
				if 0 == p.PV[pg][pidx]&(uint8(1)<<bb) {
					prime = (uint64(pg*BVpagesize|pidx) << BVprimeByteBitShift) | uint64(bb)<<1 + 3
					if prime <= pmax && prime <= 65535 {
						PrimesMoreU16 = append(PrimesMoreU16, uint16(prime))
					} else if prime <= pmax {
						PrimesMoreU32 = append(PrimesMoreU32, uint32(prime))
						bb++
						break BVPrimes_PrimeGlobalList_U16primes // break 3
					} else {
						return
					}
				}
			}
			bb = 0 // reset after scanning the initial index bits
			pidx++
		}
		if pidx >= BVpagesize {
			pg++
			pidx = 0
		}
	}
	for pg <= pgMx {
		// pidx
		for (pg < pgMx && pidx < BVpagesize) || (pg == pgMx && pidx <= pidxMx) {
			// bb
			for ; bb < BVbitsPerByte; bb++ {
				if 0 == p.PV[pg][pidx]&(uint8(1)<<bb) {
					prime = (uint64(pg*BVpagesize|pidx) << BVprimeByteBitShift) | uint64(bb)<<1 + 3
					if prime <= pmax {
						PrimesMoreU32 = append(PrimesMoreU32, uint32(prime))
					} else {
						return
					}
				}
			}
			bb = 0 // reset after scanning the initial index bits
			pidx++
		}
		if pidx >= BVpagesize {
			pg++
			pidx = 0
		}
	}
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
			if false == big.NewInt(int64(posPrime)).ProbablyPrime(int(0)) {
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
	return big.NewInt(int64(q)).ProbablyPrime(int(0))
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
	var base, power uint32
	// Start at 1 : skip already checked 2
	for ii := 1; 1 < q && ii <= PrimesSmallU8Mx; ii++ {
		qd := uint64(PrimesSmallU8[ii])
		if q < qd*qd {
			// q must be prime
			heap.Push(facts, Factorpair{Base: uint32(q), Power: uint32(1)})
			q = 1
			break
		}
		for 0 == q%qd {
			q /= qd
			power++
		}
		if 0 < power {
			heap.Push(facts, Factorpair{Base: uint32(qd), Power: uint32(power)})
			power = 0
		}
	}
	if 1 < q && q < PrimesSmallU8MxValPow2After {
		// q must be prime
		heap.Push(facts, Factorpair{Base: uint32(q), Power: uint32(1)})
		q = 1
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
		if q < 500_000_000 {
			unk = Factor1980PollardMonteCarlo(q, 0)
		} else {
			unk = 0
		}

		// Aiming to replace this with code that extracts factors rather than just testing for factors...
		/*
			roottest = SqrtU64(q)
			if q == roottest*roottest {
				return roottest
			}

			// Test higher roots until beneath the highest division tested prime, 41^4 ~= 2.8M ; 41^5 ~= 115.8M
			for ii := uint64(3); roottest > PrimesSmallU8[PrimesSmallU8Mx] ; ii++ {
				roottest = RootU64(q, ii)
				if q == PowU64(roottest, ii) {
					return roottest
				}
			}

			See beneath for the better tests...
		*/

		// aggressive search
		if 0 == unk {
			//if big.NewInt(int64(q)).ProbablyPrime(int(0)) {
			unk = PrimeProbBailliePSWppo(q)
			if 0 == unk {
				unk = q
			} else {
				unk = FactorStep1RootsFilter(q, PrimesSmallU8MxVal)
				if unk == q {
					unk = uint64(FactorLenstraECW(uint64(q), uint64(PrimesSmallU8MxVal)))
				}
				// Replace this with LECF
				// unk = Factor1980AutoPMC_Pass2(q, false)
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
	return big.NewInt(num).ProbablyPrime(int(0)) // 0 is sufficient for 64 bit integers
}

// returns the first factor found, q if unable to find a factor (Maybe Prime, might just be a big factor) -- anything less than 256*256 or PrimesSmallU8MxValPow2After must be prime.
func FactorStep0TD(q uint64) uint64 {
	if 2 > q {
		return q
	}
	// Special test & extract: base2, /2
	if 0 == q&1 {
		return 2
	}

	// Start at 1 : skip already checked 2
	for ii := 1; ii <= PrimesSmallU8Mx; ii++ {
		qd := uint64(PrimesSmallU8[ii])
		if q < qd*qd {
			return q
		}
		if 0 == q%qd {
			return qd
		}
	}
	return q
	// Other functions can utilize this check
	// if q < PrimesSmallU8MxValPow2After {
	//	// q must be prime
	//	heap.Push(facts, Factorpair{Base: uint32(qd), Power: uint32(1)})
	//	q = 1
	//}
}

func FactorReduceStep0TDlimit(q, limit uint64) uint64 {
	var qd uint64
	if 2 > q {
		return q
	}
	// Special test & extract: base2, /2
	for 0 == q&1 {
		return 2
	}

	// Start at 1 : skip already checked 2
	for ii := 1; ii <= PrimesSmallU8Mx && qd < limit; ii++ {
		qd = uint64(PrimesSmallU8[ii])
		if q < qd*qd {
			return q
		}
		if 0 == q%qd {
			return qd
		}
	}
	return 0
}

func FactorStep1RootsFilter(q, maxTestedPrime uint64) uint64 {
	roottest := SqrtU64(q)
	if roottest < maxTestedPrime {
		return q
	}
	if q == roottest*roottest {
		return roottest
	}

	// Test higher roots until beneath the highest division tested prime, 41^4 ~= 2.8M ; 41^5 ~= 115.8M
	for ii := uint64(3); roottest > maxTestedPrime; ii++ {
		roottest = RootU64(q, ii)
		if q == PowInt(roottest, ii) {
			return roottest
		}
	}
	return q
}

func FactorStep2ProbablyPrimeOrFactor(q uint64) uint64 {
	return PrimeProbBailliePSWppo(q)
}

func FactorStep2ProbablyNotPrime(q uint64) bool {
	//return false == big.NewInt(int64(q)).ProbablyPrime(int(0))
	return 0 != PrimeProbBailliePSWppo(q)
}

func FactorStep3BigPrimeTests(q, maxTestedPrime uint64) uint64 {
	// Reportedly, good up to ~70bits https://stackoverflow.com/questions/2267146/what-is-the-fastest-integer-factorization-algorithm
	// if q < 500_000 {
	// I've found this can miss squares, and for other numbers resists convergence
	r := Factor1980PollardMonteCarlo(q, 0)
	if 0 != r && r != q {
		return r
	}
	// }
	// Lenstra Elliptic Curve Factorization (good up to ~10^50 or beyond, and thus well past uint64)
	return uint64(FactorLenstraECW(uint64(q), uint64(maxTestedPrime)))
}

/*


FIXME: At the very least, clean this section up in a future update...

It was very difficult to find good reference material for the harder 'math' subject elements.  Wikipedia's pages tend to be higher level for most subjects, glossing over key steps that someone in the field likely assumes a reader would know.  Similarly they use maths notation since someone who's not bumping into a subject on an extremely rare basis should know that shorthand.



Integer Factorization is likely very similar to Primality (Probably Prime) tests, with the difference that proof (a factor that cleanly divides) must be found, rather than statistical lack of evidence of a proof.

Step 0: trial division / wheel factorization / filtering of known small primes.  Just shoot a bunch at it; really only the first 4 primes are required but, __testing higher primes helps later nth root tests!__

?
Step 1 : nth root tests (down to largest tested prime?) this scales up pretty fast nice 41^5 ~= 115.8M

?
Step 2 : Primality Test (check that the number is NOT probably prime...)


Note: due to 2 and 3 in the derivative, primes 2 and 3 MUST be ruled out early for EC tests.

// https://stackoverflow.com/a/2274520

Step ? : (maybe?) ONE round of 1980PollardMonteCarlo

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

// type Point3DI64 struct {
//	X, Y, Z int64
//}

func GCDeuc[INT ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](a, b INT) INT {
	// https://en.wikipedia.org/wiki/Euclidean_algorithm
	for b != 0 {
		a, b = b, a%b
	}
	return a
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

func ModMulInv[INT ~int64 | uint64](n, mod INT) INT {
	// Returns EITHER 1 or a possible factor of 'mod' - I might have written this more as documentation than with an intent to use.
	// https://en.wikipedia.org/wiki/Modular_multiplicative_inverse
	// """	If d is the greatest common divisor of a and m then the linear congruence ax ≡ b (mod m) has solutions if and only if d divides b. If d divides b, then there are exactly d solutions.[7]
	// 	A modular multiplicative inverse of an integer a with respect to the modulus m is a solution of the linear congruence
	//	ax === 1 (mod m)
	//	The previous result says that a solution exists if and only if gcd(a, m) = 1, that is, a and m must be relatively prime (i.e. coprime). Furthermore, when this condition holds, there is exactly one solution	"""
	//	https://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Computation
	//	ax + my = gcd(a,m) = 1	=>	ax - 1 = (-y)m	=>	ax = 1 (mod m)
	n %= mod
	return GCDeuc(n, mod)
}

/*

NOTE: Lenstra_elliptic-curve_factorization

IMPORTANT: The most likely to be confused / overlooked as an already known concept is the __modular inverse__ of a number as the Wikipedia article on it helpfully notes 'notation abuse' is likely to confuse a reader that (term)^(-1) is NOT to the negative power or a notation for division by, but rather that it is a special term MEANING the modular inverse therein.

As a preface:
***
* https://en.wikipedia.org/wiki/Modular_multiplicative_inverse
***
* https://en.wikipedia.org/wiki/Lenstra_elliptic-curve_factorization
* https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law
* https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
* https://en.wikipedia.org/wiki/Elliptic_curve
* https://en.wikipedia.org/wiki/Algebraic_curve#Intersection_with_a_line
* https://en.wikipedia.org/wiki/B%C3%A9zout%27s_theorem#A_line_and_a_curve  (but also skim the rest)
* https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication

Educational papers / slides
* https://math.stackexchange.com/questions/859116/lenstras-elliptic-curve-algorithm
* https://wstein.org/edu/124/lenstra/lenstra.pdf
* https://math.uchicago.edu/~may/REU2014/REUPapers/Parker.pdf
* https://sites.math.washington.edu/~morrow/336_16/2016papers/thomas.pdf
* https://web.ma.utexas.edu/users/sl55444/CompsciTalk.pdf
* https://math.mit.edu/research/highschool/primes/materials/2018/conf/7-2%20Rhee.pdf
* https://en.wikipedia.org/wiki/Addition_chain (only LOOSELY related...)
* http://www.additionchains.com/ << Find: Scholz-Brauer Conjecture look at the number chain... Then think about the doubling, repeated multiplication by 2, how that relates to the storage form of a number...
* Also related to the addition chain tangent: https://cr.yp.to/papers/pippenger-20020118-retypeset20220327.pdf  Pippenger's Exponentiation Algorithm (a brief paper by Daniel J. Bernstein)
* https://en.wikipedia.org/wiki/Modular_exponentiation#Right-to-left_binary_method

*** THIS is when I FINALLY found something that mentioned Modular_multiplicative_inverse , with an example so I could understand how it differed from what I thought the words meant!
* https://stackoverflow.com/a/75765566

Additionally; many examples assume (or worse cargo cult) that the programmer understands they're specific to the shape of numeric system described in the initial Lenstra_elliptic-curve_factorization W example.

I think the full (and gory) proof of how the third point R is calculated would probably also help me understand the Discriminant. https://en.wikipedia.org/wiki/Discriminant

2^2*A^3 + 3^2*B^2 ??? 4*A*A*A + 27*B*B  -  I think it relates to the 'secret' of math being solved by the work; the answer to the multi-dimensional version of (x - pX)*(x - pQ)*(x - pR) = 0 .

Having stated that, it's also time for me to throw out the nearly working initial draft and start from the top, better knowing what I didn't when I started.  Just like most initial versions.



*/

func FactorLenstraECW(q, maxTestedPrime uint64) uint64 {

	// TODO more effective curve forms? https://en.wikipedia.org/wiki/Elliptic_curve#Non-Weierstrass_curves

	var N, A, B, X0, Y0, K, Kmx, Kmulmx, gcd, loop uint64
	var ii, iiMx, iiPow, iiPowMx uint8
	// iiMx = PrimesSmallU8Mx
	// Tune // With a different Rand source these numbers might differ; for a reasonably fast PseudoRNG, even LQ random output, this seemed the fastest over a couple runs.
	iiMx, iiPowMx, _ = 3, 2, Kmulmx
	if maxTestedPrime < 3 {
		// MUST have already tested for /2 and /3 minimum - so take reasonable steps...
		A = uint64(FactorStep0TD(uint64(q)))
		if A != q || q < PrimesSmallU8MxValPow2After {
			return A
		}
		A = uint64(FactorStep1RootsFilter(uint64(q), PrimesSmallU8MxVal))
		if A != q {
			return A
		}
		A = FactorStep2ProbablyPrimeOrFactor(uint64(q))
		if 0 == A {
			return 0
		}
		maxTestedPrime = PrimesSmallU8MxVal
	}

	N = q
	//if 0 > q {
	//	N = -q
	//} else {
	//	N = q
	//}

	__LenECWaddMod := func(x0, y0, z0, x1, y1, z1 uint64) (uint64, uint64, uint64) {
		// returns rX, rY, gcd
		// https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
		var modMulInv, dY uint64
		if x0 != x1 || y0 != y1 {
			dY, modMulInv = SubOF64Mod(y1, y0, N), GCDeuc(N, (x1-x0)%N)
		} else {
			dY, modMulInv = AddOF64Mod(MulOF64Mod(3, MulOF64Mod(x0, x0, N), N), A, N), GCDeuc(N, MulOF64Mod(2, y0, N))
		}
		if 1 < modMulInv && modMulInv < N {
			return 0, 1, modMulInv
		}
		// L1 =: (dY * modMulInv) * ( x - p.X ) + p.Y
		// L1 R(x3, y3) =:	( (dY * modMulInv)*(dY * modMulInv) - x0 - x1	, (dY * modMulInv)*(x3 - x0) + y0 )
		// L2 R(x3, -y3) =:	( (dY * modMulInv)*(dY * modMulInv) - x0 - x1	, (dY * modMulInv)*(2*x0 + x1 - (dY * modMulInv)*(dY * modMulInv))   -  y0 )

		// REUSE: dY is now 'slope'^2 and modMullInv is now single slope
		modMulInv = MulOF64Mod(modMulInv, dY, N)
		dY = MulOF64Mod(modMulInv, modMulInv, N)
		// L2 R(x3, -y3) =:	( (dY * modMulInv)*(dY * modMulInv) - x0 - x1	, (dY * modMulInv)*(2*x0 + x1 - (dY * modMulInv)*(dY * modMulInv))   -  y0 )
		//	ALL (mod N)	dY * N - x0 - x1				modMulInv	*			(2*x0		+ x1	- dY)	- y0
		return SubOF64Mod(SubOF64Mod(dY, x0, N), x1, N), SubOF64Mod(MulOF64Mod(modMulInv, SubOF64Mod(AddOF64Mod(MulOF64Mod(2, x0, N), x1, N), dY, N), N), y0, N), 1
		// I *think* it's still correct to modulus every block of operations?  Unsure if it's faster to keep the word size small (my 128bit division binary search is SLOW, which is why I've biased towards the logic that modern CPU hardware has _way_ more human hours put into making it fast.)
	}

	__LenECWmulMod := func(k, x, y uint64) (uint64, uint64, uint64) {
		var rx, ry, rz, z uint64
		rx, ry, rz = 0, 1, 0 // the infinity / O point
		z = 1                // constant input
		// Use P for doubles, prepare the next P at the end of each cycle
		// 'repeated doubling and k's binary expression
		for {
			// If this power of 2 is part of the result, add it
			if 0 < k&1 {
				rx, ry, rz = __LenECWaddMod(rx, ry, rz, x, y, z)
			}
			// every addition should be a valid point, so every GCD return is desired
			if 1 < rz {
				return rx, ry, rz
			}

			k >>= 1
			if 0 == k {
				break
			}

			// Compute the next doubling
			x, y, z = __LenECWaddMod(x, y, z, x, y, z)
			// every addition should be a valid point, so every GCD return is desired
			if 1 < z {
				return x, y, z
			}
		}
		return rx, ry, rz
	}

	__LenECWcurve := func() uint64 {
		// Uses function context variables
		// https://en.wikipedia.org/wiki/Lenstra_elliptic-curve_factorization#The_algorithm_with_projective_coordinates
		// 1.
		// Strongly random isn't as necessary as just 'not trivial' and 'not always the same'
		A, X0, Y0 = 0, 0, 0
		for 0 == A {
			A = uint64(PSRand.RandInt(uint64(N))) // >>4 + 1
		}
		for 0 == X0 {
			X0 = uint64(PSRand.RandInt(uint64(N)))
		}
		for 0 == Y0 {
			Y0 = uint64(PSRand.RandInt(uint64(N)))
		}
		// Note 'choice of B' https://wstein.org/edu/Fall2001/124/lenstra/lenstra.pdf pg 662

		// 2. y*y = x*x*x + A*x + B (mod N)
		// B = (Y0*Y0 - X0*X0*X0 - A*X0) % N
		B = SubOF64Mod(SubOF64Mod(MulOF64Mod(Y0, Y0, N), MulOF64Mod(X0, MulOF64Mod(X0, X0, N), N), N), MulOF64Mod(A, X0, N), N)

		// math.SE Step 4
		// https://math.stackexchange.com/questions/859116/lenstras-elliptic-curve-algorithm
		// for the Y^2 = X^3... curve (derivative?)
		// check gcd(4 * A*A*A + 27 * B * B , N ) == 1 (OK) || N (Bad, New A, recheck) || between 1 and N => YAY (maybe composite)Factor found

		// Test that the curve is 'square free' in X? https://en.wikipedia.org/wiki/Elliptic_curve
		// Discriminant https://en.wikipedia.org/wiki/Discriminant
		// """ the quantity which appears under the square root in the quadratic formula. If a ≠ 0 , {\displaystyle a\neq 0,} this discriminant is zero if and only if the polynomial has a double root. In the case of real coefficients, it is positive if the polynomial has two distinct real roots, and negative if it has two distinct complex conjugate roots.[1] Similarly, the discriminant of a cubic polynomial is zero if and only if the polynomial has a multiple root. In the case of a cubic with real coefficients, the discriminant is positive if the polynomial has three distinct real roots, and negative if it has one real root and two distinct complex conjugate roots.
		// tests properties of roots of the curve
		// return 4*A*A*A + 27*B*B
		return AddOF64Mod(MulOF64Mod(MulOF64Mod(4, A, N), MulOF64Mod(A, A, N), N), MulOF64Mod(27, MulOF64Mod(B, B, N), N), N)
	}

FactorLenstraECW_reroll:
	for {
		Discr := __LenECWcurve()
		if 0 == Discr {
			// Bad roll, try again
			continue FactorLenstraECW_reroll
		}
		gcd = GCDeuc(Discr, N)
		if gcd == N {
			// Bad roll, try again
			continue FactorLenstraECW_reroll
		}
		if 1 < gcd {
			// Lucky roll, found a GCD between 1 and N
			return gcd
		}

		loop++
		if 0 == loop&0xFFF {
			fmt.Printf("FactorLenstraECW(%d, %d) loop %d:\tA: %d\tB:%d\t\t%x\t%x\t%x\t%x\n", q, maxTestedPrime, loop, A, B, PSRand.RandInt(uint64(N)), PSRand.RandInt(uint64(N)), PSRand.RandInt(uint64(N)), PSRand.RandInt(uint64(N)))
			if loop > 0xD000 {
				panic("Too slow")
			}
		}

		// https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law
		// 3. + 4. ???
		// math.SE Step 5 (BLCM == k)
		//	Choose B -- All Prime Factors must be Less than or Equal to B ;
		//Kmx = N / uint64(maxTestedPrime)
		_ = Kmx
		K = 2 // 2 * 3 * 5
		iiPow, ii = 1, 0
		//for K < Kmx {
		for {
			// 5.
			// https://math.mit.edu/research/highschool/primes/materials/2018/conf/7-2%20Rhee.pdf
			// I think that's slide 15, with Lenstra's Factorization Method and steps 1 .. 11
			// This is Step 5, but 6..9 are built within this loop
			// Compute ? Q using d * P (mod n) 'and set P = Q'

			// math.SE Step 6
			X0, Y0, gcd = __LenECWmulMod(K, X0, Y0)
			// math.SE Step 7
			// Calc D = gcd(Dk, n) Same general rules as Step 5 above, 1<D<N => D is a non-trivial factor, else Raise K or roll again.
			if 1 < gcd && gcd < N {
				return gcd
			}
			ii++ //|| Kmulmx < int64(PrimesSmallU8[ii]) {
			if ii > iiMx {
				ii = 0
				iiPow++
			}
			if iiPow > iiPowMx {
				continue FactorLenstraECW_reroll
			}
			K *= uint64(PrimesSmallU8[ii])

		}
	}
	// unreached
	return N
}

// More refs...
// https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication
// https://stackoverflow.com/questions/30017367/lenstras-elliptic-curve-factorization-problems

// DANGER NOTE: Some of the operations (E.G. 3 * pX * pX + A) could overflow the 63 bit limit before the mod; I'm _mostly_ sure that their construction method should keep them within bounds, but some larger N to factor might make them overflow...
// https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law
// https://sites.math.washington.edu/~morrow/336_16/2016papers/thomas.pdf (pg9-10)
// https://math.uchicago.edu/~may/REU2014/REUPapers/Parker.pdf
/*
//
		// I still like the idea of this version of seeding the curve, even if it might not have sufficient range.
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
//
*/

// By Euler 0058 it was clear that a better primality test is REQUIRED
// https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test

// https://en.wikipedia.org/wiki/Primality_test#Miller%E2%80%93Rabin_and_Solovay%E2%80%93Strassen_primality_test
// https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test
// == MR + LL
// https://en.wikipedia.org/wiki/Lucas_pseudoprime

func PrimeProbMillerRabinInner(n, blim uint64, prob uint8) uint64 {
	// Returns 0 if 'might be prime', 1 if it's composite but unfactored, and (multiple of factor(s)) otherwise
	// blim allows targeted tests by an external function https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Testing_against_small_sets_of_bases
	var a, x, s, d, y uint64
	if 0 == n&1 {
		return 2
	}
	if 0 == prob {
		prob = 1
	}
	if 4 > n {
		return 0
	}
	// find s and d
	d = n - 1
	for 0 == d&0xFF {
		s += 8
		d >>= 8
	}
	if 0 == d&0xF {
		s += 4
		d >>= 4
	}
	if 0 == d&0x3 {
		s += 2
		d >>= 2
	}
	if 0 == d&0x1 {
		s++
		d >>= 1
	}
	if n-1 != d<<s || 1 != d&1 || 0 == s {
		fmt.Printf("n-1 ? 2^d => %d == %d\ts: %d\td: %d\n", n-1, d<<s, s, d)
		panic("Programmer Logic Error!")
	}
	// for d * (1 << s) != a {
	//	s++
	//	d = a / (1 << s)
	// }
	for ; 0 < prob; prob-- {
		if 1 < blim {
			a, prob = blim, 1
		} else {
			a = PSRand.RandInt(n-4) + 2
		}
		x = PowIntMod(a, d, n)
		for ; 0 < s; s-- {
			y = PowIntMod(x, 2, n)
			if 1 == y && 1 != x && x != n-1 {
				return GCDbin(x-1, n)
			}
			x = y
		}
		if 1 != y {
			// fmt.Printf("\tx: %d\ty: %d", x, y)
			return 1
		}
	}
	return 0
}

func PrimeOptiTestMillerRabin(q uint64) uint64 {
	var a []uint8
	// https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Testing_against_small_sets_of_bases
	// Pomerance, Selfridge, Wagstaff[4] and Jaeschke
	switch {
	case 2047 > q:
		a = []uint8{2}
	case 1_373_653 > q:
		a = []uint8{2, 3}
	case 9_080_191 > q:
		a = []uint8{31, 73}
	case 25_326_001 > q:
		a = []uint8{2, 3, 5}
	case 3_215_031_751 > q:
		a = []uint8{2, 3, 5, 7}
	case 4_759_123_141 > q:
		a = []uint8{2, 7, 61}
	case 1_122_004_669_633 > q:
		a = []uint8{2, 13, 23} // ,1_662_803  THIS is the only outlier for uint16, let alone uint8 which every other number falls into
		res := PrimeProbMillerRabinInner(q, 1_662_803, 1)
		if 0 != res {
			return res // These don't have to be run in any particular order?  Probably fine...
		}
	case 2_152_302_898_747 > q:
		a = []uint8{2, 3, 5, 7, 11}
	case 3_474_749_660_383 > q:
		a = []uint8{2, 3, 5, 7, 11, 13}
	case 341_550_071_728_321 > q:
		a = []uint8{2, 3, 5, 7, 11, 13, 17}
	case 3_825_123_056_546_413_051 > q:
		a = []uint8{2, 3, 5, 7, 11, 13, 17, 19, 23}
	default:
		// < 2^64 Feitsma and Galway
		a = []uint8{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
		// Sorenson and Webster[13]
		// < 318,665,857,834,031,151,167,461 :: a = 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37
		// < 3,317,044,064,679,887,385,961,981 :: a = 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41
	}
	iiLim := len(a)
	for ii := 0; ii < iiLim; ii++ {
		res := PrimeProbMillerRabinInner(q, uint64(a[ii]), 1)
		if 0 != res {
			return res
		}
	}
	return 0
}

/*
func JacobiSym(n, d int64) int8 {
	var ret int8
	if 1 > d || 0 == d&1 {
		return -127
	}
	// https://en.wikipedia.org/wiki/Jacobi_symbol#Calculating_the_Jacobi_symbol
	ret = 1
	// #1 rule 2 ; #3 == Rule 3 ; re-do step 1# rule 2
	for n %= d; 0 != n; n %= d {
		// #2 rule 9
		for 0 == n&1 {
			n >>= 1
			// rule 9 == 3 || 5
			if 3 == d&(0b111) || 5 == d&(0b111) {
				ret = -ret
			}
		}
		// #4 flip rule 6 (flip the return if both are 3 in mod 4
		n, d = d, n
		if 3 == n&3 && 3 == d&3 {
			ret = -ret
		}
	}
	// Rule 3 use ret (gcd) if == 1, otherwise 0
	if 1 == d {
		return ret
	}
	return 0
}
*/

func JacobiSym(sn int64, d uint64) int8 {
	var ret int8
	if 1 > d || 0 == d&1 {
		return -127
	}
	// https://en.wikipedia.org/wiki/Jacobi_symbol#Calculating_the_Jacobi_symbol
	ret = 1
	n := uint64(sn) % d
	// #1 rule 2 ; #3 == Rule 3 ; re-do step 1# rule 2
	for ; 0 != n; n %= d {
		// #2 rule 9
		for 0 == n&1 {
			n >>= 1
			// rule 9 == 3 || 5
			if 3 == d&(0b111) || 5 == d&(0b111) {
				ret = -ret
			}
		}
		// #4 flip rule 6 (flip the return if both are 3 in mod 4
		n, d = d, n
		if 3 == n&3 && 3 == d&3 {
			ret = -ret
		}
	}
	// Rule 3 use ret (gcd) if == 1, otherwise 0
	if 1 == d {
		return ret
	}
	return 0
}

/*
	I took a couple hours Monday and couldn't quite get the Lucas Probable Prime test to work based on what was on Wikipedia.
	Golang's math/big library has a working implementation and even at a quick skim it mentions some things I'd been unclear about (E.G. the Jacobi Symbol prefs for a stronger test).
	https://cs.opensource.google/go/go/+/refs/tags/go1.23.1:src/math/big/prime.go;l=26

	I'll be trying to annotate where my current attempts and the working example differ.

	* A different set of P D Q values, 'method c' (important for the stronger test, but not for working at all)
	* VERY different V equation(s) that do not utilize the Uk or Qk terms at all.
	* Clear test values rather than a paragraph in an article

	It also references:
	//
	Baillie and Wagstaff, "Lucas Pseudoprimes", Mathematics of Computation 35(152),
	October 1980, pp. 1391-1417, especially page 1401.
	https://www.ams.org/journals/mcom/1980-35-152/S0025-5718-1980-0583518-6/S0025-5718-1980-0583518-6.pdf
	//
	Grantham, "Frobenius Pseudoprimes", Mathematics of Computation 70(234),
	March 2000, pp. 873-891.
	https://www.ams.org/journals/mcom/2001-70-234/S0025-5718-00-01197-2/S0025-5718-00-01197-2.pdf
	//
	Baillie, "Extra strong Lucas pseudoprimes", OEIS A217719, https://oeis.org/A217719.
	//
	Jacobsen, "Pseudoprime Statistics, Tables, and Data", http://ntheory.org/pseudoprimes.html.
	//
	Nicely, "The Baillie-PSW Primality Test", https://web.archive.org/web/20191121062007/http://www.trnicely.net/misc/bpsw.html.
	(Note that Nicely's definition of the "extra strong" test gives the wrong Jacobi condition,
	as pointed out by Jacobsen.)
	//
	Crandall and Pomerance, Prime Numbers: A Computational Perspective, 2nd ed.
	Springer, 2005.


*/

func PrimeProbLucasStrong(n uint64) uint64 {
	// https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test#The_test
	// This assumes that no factors were found using
	// 0# Trial Division (of small primes) E.G. FactorStep0TD(q) // FactorReduceStep0TDlimit(q, 7)
	// 1# PrimeProbMillerRabinInner(q, n, 1)

	var js int8
	var D, P, qd uint64

	if 2 >= n {
		if 2 == n {
			return 0
		}
		return 1 // Neither 1 nor 0 are 'prime' though they also can't be factored...
	}
	if 0 == n&1 {
		return 2 // even number > 2, not prime
	}

	// Baillie-OEIS 'method C' for P D Q ; Initial P=3 , test: P*P - 4 = D => Jacobi(D,n) == -1
	for P = 3; ; P++ {
		// Someone probably already burned several CPU-years of time on the 64 bit integer space but just in case
		if P > 10_000 {
			panic(fmt.Sprintf("math/big said this was thought to be impossible and should be reported to researchers: could not find (D/n) = -1 for %#x", n))
		}
		D = P*P - 4
		js = JacobiSym(int64(D), n)
		if -1 == js {
			break
		} else if 0 == js {
			// 0 means D and n(input) share a prime factor, P*P - 4 == (P-2)(P+2) ; since the loop is increasing and starts wtih p-2=1, the new possible factor is the p+2 side
			if P+2 == n {
				return 0 // the new 'factor' is the prime being tested
			}
			return P + 2
		}
		// math/big uses 40, but is intended for large numbers; I had an effective guess of 10 checks in the first version of my code.  This isn't tuned, it's just a split of the difference.
		if 20 == P {
			qd = uint64(SqrtU64(uint64(n)))
			if qd*qd == n {
				return qd // if n(in) is a square JS will always be 1, test for that if a few iterations don't find an answer
			}
		}
	}

	// math/big mentions """
	// Grantham definition of "extra strong Lucas pseudoprime", after Thm 2.3 on p. 876 (D, P, Q above have become Δ, b, 1):
	// ... 'GCD(n, D~~Δ) == 1 or 0 would have been found above, and GCD(n, 2) == 1 because n is odd.'
	// s = (n - Jacobi(D~~Δ, n)) / 2^r = (n+1) / 2^r.  """
	var r, s, ns2, Vk, Vk1, Pmodsub, Tz, Th, Tl uint64
	for s = n + 1; 0 == s&1; s >>= 1 {
		r++
	}
	ns2 = n - 2
	Vk, Vk1, Pmodsub = 2, P, (n-P)%n
	if P > n {
		fmt.Printf("Unexpected: %#x <P %#x\n", n, P)
	}

	// math/big mentions the 'almost extra strong' test lacks a modular inversion (I assume that's of D with the notation abuse encountered earlier)...
	// " possible to recover U_n using Crandall and Pomerance equation 3.13: U_n = D^-1 (2V_{n+1} - PV_n) allowing us to run the full extra-strong test "

	// Out of the three Lucas components, math/big only calculates V from Vs(b, 1) (since Q was 1 and P == b in some literature)
	// Their V equations: (Remember that Q1 == 1, so that term muls out)
	// V(2k) = V(k)*V(k) - 2
	// V(2k+1) = V(k) * V(k+1) - P <<< THIS EQUATION isn't even on Wikipedia's pages, or at least I didn't see it in a quick visual search

	// math/big calculates 2k+1 and 2k on each pass to avoid the U chain

	// Start with K=0, build up in log2(s) using the MSB first method, uint64s
	for ii := 64 - 1 - BitsLeadingZeros64(s); 0 <= ii; ii-- {
		// fmt.Printf("%#x iter\t%#10x\t%#10x\tif %t\n", n, s, s_verify, 0 < s&(1<<ii))
		if 0 < s&(1<<ii) {
			// V(step) = V(2k+1) = V(k) * V(k+1) - P
			Th, Tl = UU64MulWrap(Vk, Vk1)
			Tz, Th, Tl = UU64AddUU64(Th, Tl, 0, Pmodsub)
			if 0 != Tz {
				fmt.Printf("Overflow: %#x OF2k+1 %#x\n", n, Tz)
			}
			Vk = UU64Mod(Th, Tl, n)

			// V(step + 1) = V(2k+2) = V(old+1) * V(old+1) - 2
			Th, Tl = UU64MulWrap(Vk1, Vk1)
			Tz, Th, Tl = UU64AddUU64(Th, Tl, 0, ns2)
			if 0 != Tz {
				fmt.Printf("Overflow: %#x OF2k+1 %#x\n", n, Tz)
			}
			Vk1 = UU64Mod(Th, Tl, n)
		} else {
			// V(step + 1) = V(2k) = V(k) * V(k+1) - p
			Th, Tl = UU64MulWrap(Vk, Vk1)
			Tz, Th, Tl = UU64AddUU64(Th, Tl, 0, Pmodsub)
			if 0 != Tz {
				fmt.Printf("Overflow: %#x Of2k+1 %#x\n", n, Tz)
			}
			Vk1 = UU64Mod(Th, Tl, n)

			// V(step) = V(2k) = V(old) * V(old) - 2
			Th, Tl = UU64MulWrap(Vk, Vk)
			Tz, Th, Tl = UU64AddUU64(Th, Tl, 0, ns2)
			if 0 != Tz {
				fmt.Printf("Overflow: %#x Of2k+1 %#x\n", n, Tz)
			}
			Vk = UU64Mod(Th, Tl, n)
		}
	}

	// For this form the two roots (p - 2)(p + 2) must be checked
	if 2 == Vk || ns2 == Vk {
		// "Jacobsen, apply Crandall and Pomerance equation 3.13:" check U(s~~k) == 0
		// U(k) = ModMulInv(D, n) * (2*V(k+1) - P * V(k)
		// math/go "U(k) == 0" 'can be checked via ' " 2 * V(k+1) == P * V(k)   mod n "
		qd, Tz = UU64MulWrap(2, Vk1)
		Th, Tl = UU64MulWrap(P, Vk)
		if 0 < UU64Cmp(qd, Tz, Th, Tl) {
			Th, Tl = UU64SubUU64(qd, Tz, Th, Tl)
		} else {
			Th, Tl = UU64SubUU64(Th, Tl, qd, Tz)
		}
		qd = UU64Mod(Th, Tl, n)
		if 0 == qd {
			// fmt.Printf("%d\tprime exit 1\n", n)
			return 0 // Confirmed prime
		}
	}

	// Test V(2^t * s) == 0 for the range 0 <= t <= r - 1 (earlier R was the power of 2 extracted from n+1)
	// R isn't used anywhere after this, so...
	for ; 0 < r; r-- {
		if 0 == Vk {
			// fmt.Printf("%d\tprime exit 2\n", n)
			return 0 // Confirmed prime
		}
		if 2 == Vk {
			// fmt.Printf("%d\t2 loop\n", n)
			return 1 // Vk == 2 'is a fixed point for V(step) = V(k)*V(k) - 2 ' That makes sense.
		}
		// V(step) = Vk * Vk - 2
		Th, Tl = UU64MulWrap(Vk, Vk)
		Tz, Th, Tl = UU64AddUU64(Th, Tl, 0, ns2)
		if 0 != Tz {
			fmt.Printf("Overflow: %#x Of2^r %#x\n", n, Tz)
		}
		Vk = UU64Mod(Th, Tl, n)
	}
	// fmt.Printf("%d\texit composite\n", n)
	return 1 // composite
}

func PrimeProbBailliePSW(N uint64) uint64 {
	var ret uint64
	ret = FactorReduceStep0TDlimit(N, 7)
	if 0 != ret {
		return ret
	}
	return PrimeProbBailliePSWppo(N)
}

func PrimeProbBailliePSWppo(N uint64) uint64 {
	var ret uint64
	// return PrimeOptiTestMillerRabin(N) // FIXME - Covering broken BPSW, remove when ready
	ret = PrimeProbMillerRabinInner(N, 2, 1)
	if 0 != ret {
		return ret
	}
	return PrimeProbLucasStrong(N)
}

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

	factor := FactorReduceStep0TDlimit(qt, 41)
	if 0 != factor { return factor }

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

func EulerTotientPhi_old(n uint64) uint64 {
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

func EulerTotientPhi(q, rmin uint64) uint64 {
	var ret, qd, qdqd uint64
	var ii int
	// var ii, iiLim int
	qdqd, ret = 1, q

	// https://en.wikipedia.org/wiki/Euler%27s_totient_function#Computing_Euler's_totient_function
	// Why is this so slow?  Integer version has SO MANY div / multiplies while the fraction version can ignore that and also give an immediate remaining residual
	// Of course the later wasn't considered as part of a generic interface... so that optimization option is (maybe?) Euler 70 specific

	__ETP_TestQDFactor := func() {
		if 0 == q%qd {
			for 0 == q%qd {
				q /= qd
			}
			ret -= ret / qd
			if qdqd < qd*qd {
				qdqd = qd * qd
			}
		}
	}

	// 2
	if 0 == q&1 {
		ret -= ret >> 1 // *(1 - 1/qd)
		for 0 == q&1 {
			q >>= 1
		}
		qdqd = 4
	}

	// Start at 1 : skip already checked 2 -- PrimesSmallU8Mx // Tested a sample of values, all uint8 primes yielded the fastest result (at the time of the test)
	for ii = 1; 1 < q && q > qdqd && rmin < ret && ii <= PrimesSmallU8Mx; ii++ {
		qd = uint64(PrimesSmallU8[ii])
		__ETP_TestQDFactor()
		if 1 == q || rmin >= ret {
			return ret
		}
		if q <= qdqd {
			return ret - ret/q
		}
	}

	for 1 < q && qdqd < q && rmin < ret {
		qd = PrimeProbBailliePSWppo(q)
		if 0 == qd {
			break
		}

		if 1 < qd {
			if 0 != q%qd {
				panic(fmt.Sprintf("Bad factor returned, this should not happen: %#x not-factor %#x", q, qd))
			}
			__ETP_TestQDFactor()
			continue
		}

		qd = Factor1980PollardMonteCarlo(q, 0)
		if 0 != qd {
			__ETP_TestQDFactor()
			continue
		}
		qd = FactorLenstraECW(q, uint64(PrimesSmallU8MxVal))
		if 0 != qd {
			__ETP_TestQDFactor()
			continue
		}
	}

	if 1 == q {
		return ret
	}

	/*
		ii, iiLim = 0, len(PrimesMoreU16)
		for 1 < q && rmin < ret && ii < iiLim {
			// ProbPrime _expensive_ but WORTH it... My own version rather than math.big would be good though...
			if q < qdqd || Primes.ProbPrime(q) {
				// q must be a single prime number
				return ret - ret/q
			}

			// pprof-ed slow, but much slower without *shrug*
			qd = Factor1980PollardMonteCarlo(q, 0)
			if 0 != qd && 0 == q%qd {
				for 0 == q%qd {
					q /= qd
				}
				ret -= ret / qd // *(1 - 1/qd)
				continue
			}

			// fmt.Printf("Debug q16: q: %d\tqdqd: %d\tii: %d\tiiLim %d\n", q, qdqd, ii, iiLim)
			for q > qdqd && ii < iiLim {
				qd = uint64(PrimesMoreU16[ii]) // Primes.PrimeAfter(qd)
				ii++
				qdqd = qd * qd
				if 0 == q%qd {
					for 0 == q%qd {
						q /= qd
					}
					ret -= ret / qd // *(1 - 1/qd)
					break           // 1
				}
			}
		}
		ii, iiLim = 0, len(PrimesMoreU32)
		for 1 < q && rmin < ret && ii < iiLim {
			// ProbPrime _expensive_ but WORTH it... My own version rather than math.big would be good though...
			if q < qdqd || Primes.ProbPrime(q) {
				// q must be a single prime number
				return ret - ret/q
			}

			// pprof-ed slow, but much slower without *shrug*
			qd = Factor1980PollardMonteCarlo(q, 0)
			if 0 != qd && 0 == q%qd {
				for 0 == q%qd {
					q /= qd
				}
				ret -= ret / qd // *(1 - 1/qd)
				continue
			}

			for q > qdqd && ii < iiLim {
				qd = uint64(PrimesMoreU32[ii]) // Primes.PrimeAfter(qd)
				ii++
				qdqd = qd * qd
				if 0 == q%qd {
					for 0 == q%qd {
						q /= qd
					}
					ret -= ret / qd // *(1 - 1/qd)
					break           // 1
				}
			}
		}
	*/
	// qd := FactorStep1RootsFilter(q, uint64(PrimesSmallU8MxVal))
	// Lenstra Elliptic Curve Factorization (good up to ~10^50 or beyond, and thus well past uint64)
	//return uint64(FactorLenstraECW(int64(q), int64(PrimesSmallU8MxVal)))
	return ret - ret/q
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
