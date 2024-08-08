// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package bitvector

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

import (
	// "bufio"
	// "euler"
	// "fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

type BitVector struct {
	Minset uint64
	Maxset uint64
	vec []uint64 // 2^6 = 64 numbers per unit, 0x ffff ffff ffff ffff
}

func NewBitVector(max int64) *BitVector {
	vecsz := (max >> 6) + 1
	return &BitVector{0, 0, make([]uint64, vecsz, vecsz)}
}

func (bv *BitVector) Set(num uint64) {
	if bv.Minset > num {
		bv.Minset = num
	}
	if bv.Maxset < num {
		bv.Maxset = num
	}
	shift := num & 0x3f
	bv.vec[num>>6] |= uint64(1) << shift
}

func (bv *BitVector) Clear(num uint64) {
	shift := num & 0x3f
	bv.vec[num>>6] &^= uint64(1) << shift
}

func (bv *BitVector) Test(num uint64) bool {
	shift := num & 0x3f
	return 0 < bv.vec[num>>6] & uint64(1) << shift
}

func (bv *BitVector) GetInts() []int {
	ret := make([]int, 0, 4)
	limit := bv.Maxset >> 6 + 1
	for ii := bv.Minset >> 6 ; ii < limit ; ii++ {
		bb := uint64(1)
		for ff := 0; ff < 64; ff++ {
			if 0 < bv.vec[ii] & bb {
				ret = append(ret, int(ii << 6) + ff )
			}
			bb <<= 1
		}
	}
	return ret
}
