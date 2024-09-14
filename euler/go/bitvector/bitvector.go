package bitvector

// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

import (
// "math/big" // Do I want a BigIntBitVector?
)

type BitVector struct {
	Minset uint64
	Maxset uint64
	vec    []uint64 // 2^6 = 64 numbers per unit, 0x ffff ffff ffff ffff
}

func NewBitVector(max uint64) *BitVector {
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
	return 0 < bv.vec[num>>6]&uint64(1)<<shift
}

func (bv *BitVector) TestAndSet(num uint64) bool {
	shift := num & 0x3f
	ret := 0 < bv.vec[num>>6]&uint64(1)<<shift
	if bv.Minset > num {
		bv.Minset = num
	}
	if bv.Maxset < num {
		bv.Maxset = num
	}
	bv.vec[num>>6] |= uint64(1) << shift
	return ret
}

func (bv *BitVector) GetInts() *[]int {
	ret := make([]int, 0, 4)
	limit := bv.Maxset>>6 + 1
	for ii := bv.Minset >> 6; ii < limit; ii++ {
		bb := uint64(1)
		for ff := 0; ff < 64; ff++ {
			if 0 < bv.vec[ii]&bb {
				ret = append(ret, int(ii<<6)+ff)
			}
			bb <<= 1
		}
	}
	return &ret
}

func (bv *BitVector) GetUInt64s() *[]uint64 {
	ret := make([]uint64, 0, 4)
	limit := bv.Maxset>>6 + 1
	for ii := bv.Minset >> 6; ii < limit; ii++ {
		bb := uint64(1)
		for ff := 0; ff < 64; ff++ {
			if 0 < bv.vec[ii]&bb {
				ret = append(ret, ii<<6+uint64(ff))
			}
			bb <<= 1
		}
	}
	return &ret
}

type OffsetBitVector struct {
	Minset int64
	Maxset int64
	Offset int64
	vec    []uint64 // 2^6 = 64 numbers per unit, 0x ffff ffff ffff ffff
}

func NewOffsetBitVector(min, max int64) *OffsetBitVector {
	vecsz := ((max - min) >> 6) + 1
	return &OffsetBitVector{0, 0, min, make([]uint64, vecsz, vecsz)}
}

func (bv *OffsetBitVector) Set(num int64) {
	if bv.Minset > num {
		bv.Minset = num
	}
	if bv.Maxset < num {
		bv.Maxset = num
	}
	offnum := uint64(num - bv.Offset)
	shift := offnum & 0x3f
	bv.vec[offnum>>6] |= uint64(1) << shift
}

func (bv *OffsetBitVector) Clear(num int64) {
	offnum := uint64(num - bv.Offset)
	shift := offnum & 0x3f
	bv.vec[offnum>>6] &^= uint64(1) << shift
}

func (bv *OffsetBitVector) Test(num int64) bool {
	offnum := uint64(num - bv.Offset)
	shift := offnum & 0x3f
	return 0 < bv.vec[offnum>>6]&uint64(1)<<shift
}

func (bv *OffsetBitVector) TestAndSet(num int64) bool {
	offnum := uint64(num - bv.Offset)
	shift := offnum & 0x3f
	ret := 0 < bv.vec[offnum>>6]&uint64(1)<<shift
	if bv.Minset > num {
		bv.Minset = num
	}
	if bv.Maxset < num {
		bv.Maxset = num
	}
	bv.vec[offnum>>6] |= uint64(1) << shift
	return ret
}

func (bv *OffsetBitVector) GetInts() []int {
	ret := make([]int, 0, 4)
	limit := uint64(bv.Maxset-bv.Offset)>>6 + 1
	for ii := uint64(bv.Minset-bv.Offset) >> 6; ii < limit; ii++ {
		bb := uint64(1)
		for ff := uint64(0); ff < uint64(64); ff++ {
			if 0 < bv.vec[ii]&bb {
				// FIXME incomplete
				ret = append(ret, int(int64(ii<<6+ff)+bv.Offset))
			}
			bb <<= 1
		}
	}
	return ret
}

func (bv *OffsetBitVector) GetInt64s() []int64 {
	ret := make([]int64, 0, 4)
	limit := uint64(bv.Maxset-bv.Offset)>>6 + 1
	for ii := uint64(bv.Minset-bv.Offset) >> 6; ii < limit; ii++ {
		bb := uint64(1)
		for ff := int64(0); ff < int64(64); ff++ {
			if 0 < bv.vec[ii]&bb {
				// FIXME incomplete
				ret = append(ret, int64(ii<<6)+ff+bv.Offset)
			}
			bb <<= 1
		}
	}
	return ret
}
