// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0

import (
	"fmt"
	"syscall"
)

// a struct syscall.Rusage contains two syscall.Timeval fields https://pkg.go.dev/syscall@go1.23.2#Rusage https://pkg.go.dev/syscall@go1.23.2#Timeval
type CPUticks struct {
	Utime, Stime int64
}

func CPUtime() CPUticks {
	acc := new(syscall.Rusage)
	syscall.Getrusage(syscall.RUSAGE_SELF, acc)
	return CPUticks{Utime: acc.Utime.Nano(), Stime: acc.Stime.Nano()}
}

func FactorialUint64(ii uint64) uint64 {
	ret := uint64(1)
	for ii > 1 {
		ret *= ii
		ii--
	}
	return ret
}

func PermutationString(perm int, str string) string {
	end := len(str)
	tmp := make([]byte, end)
	copy(tmp, str)
	res := make([]byte, end)
	slot := 0
	for slot < end {
		fact := int(FactorialUint64(uint64(end - 1 - slot)))
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

// Good enough as is on smaller inputs... overflows on slot 99
func PermutationOldSlUint8(perm uint64, con []uint8) []uint8 {
	end := len(con)
	tmp := make([]uint8, end)
	copy(tmp, con)
	res := make([]uint8, end)
	slot := 0
	for slot < end {
		fact := FactorialUint64(uint64(end - 1 - slot))
		if 0 == fact {
			fmt.Printf("Fact returned 0, inputs: %d\t%d\t%d\n", uint64(end-1-slot), end, slot)
		}
		idx := perm / fact
		perm %= fact
		res[slot] = tmp[idx]
		// fmt.Print(slot, idx, "\t", res, "\t", tmp, "\t")
		for idx < uint64(end-1-slot) {
			tmp[idx] = tmp[idx+1]
			idx++
		}
		// fmt.Println(tmp)
		slot++
	}
	return res
}

// Usually ~20% faster, may require a temp buffer if unable to store the index to use in the output array.
func PermutationNewSlUint8(perm uint64, con []uint8) []uint8 {
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

func CompareMethods(perm uint64, arr []uint8) bool {
	end := len(arr)
	tmp := make([]uint8, end)
	//copy(tmp, str)
	for ii := 0; ii < end; ii++ {
		tmp[ii] = uint8(ii)
	}
	res := make([]uint8, end)
	res1 := make([]uint8, end)
	slot := 0
	seed := perm
	for slot < end {
		fact := FactorialUint64(uint64(end - 1 - slot))
		idx := seed / fact
		seed %= fact
		res[slot] = tmp[idx]
		res1[slot] = uint8(idx)
		// fmt.Print(slot, idx, "\t", res, "\t", tmp, "\t")
		for idx < uint64(end-1-slot) {
			tmp[idx] = tmp[idx+1]
			idx++
		}
		// fmt.Println(tmp)
		slot++
	}

	for ii := 0; ii < end; ii++ {
		tmp[ii] = uint8(ii)
	}
	res2 := make([]uint8, end)
	res3 := make([]uint8, end)
	mod := uint64(1)
	seed = perm
	for ii := end; 0 < ii; ii-- {
		res2[uint64(end)-mod] = uint8(seed % mod)
		seed /= mod
		mod++
	}
	for ii := 0; ii < end; ii++ {
		idx := res2[ii]
		res3[ii] = tmp[idx]
		for idx < uint8(end-1-ii) {
			tmp[idx] = tmp[idx+1]
			idx++
		}
	}
	for ii := 0; ii < end; ii++ {
		if res[ii] != res3[ii] {
			fmt.Printf("Failed %d\n%v\n%v\n\n%v\n%v\n\n", perm, res, res3, res1, res2)
			return false
		}
	}
	return true
}

/*
for ii in *\/*.go ; do go fmt "$ii" ; done ; go fmt ex_permutation.go ; go run ex_permutation.go
*/
func main() {
	//test
	fmt.Printf("Comparing 2 million seeds\n")
	start := CPUtime()
	for ii := uint64(0); ii <= uint64(2_000_000); ii++ {
		CompareMethods(ii, make([]uint8, 10))
	}
	end := CPUtime()
	fmt.Printf("Took %f mS\n", float64((end.Utime-start.Utime)+(end.Stime-start.Stime))/float64(1000))

	//run

	deck := make([]uint8, 10)
	for ii := 0; ii < len(deck); ii++ {
		deck[ii] = uint8(ii)
	}
	fmt.Printf("\n\nNew Deck: %v\n", deck)

	fmt.Printf("PermutationOldSlUint8 [10] * 2.5 million seeds\n")
	start = CPUtime()
	for ii := uint64(0); ii <= uint64(2_500_000); ii++ {
		_ = PermutationOldSlUint8(ii, deck)
	}
	end = CPUtime()
	fmt.Printf("Took %f mS\n", float64((end.Utime-start.Utime)+(end.Stime-start.Stime))/float64(1000))

	fmt.Printf("PermutationNewSlUint8 [10] * 2.5 million seeds\n")
	start = CPUtime()
	for ii := uint64(0); ii <= uint64(2_500_000); ii++ {
		_ = PermutationNewSlUint8(ii, deck)
	}
	end = CPUtime()
	fmt.Printf("Took %f mS\n", float64((end.Utime-start.Utime)+(end.Stime-start.Stime))/float64(1000))

	deck = make([]uint8, 52)
	for ii := 0; ii < len(deck); ii++ {
		deck[ii] = uint8(ii)
	}
	fmt.Printf("\n\nNew Deck: %v\n", deck)

	fmt.Printf("PermutationOldSlUint8 [%d] * 2.5 million seeds\n", len(deck))
	start = CPUtime()
	for ii := uint64(0); ii <= uint64(2_500_000); ii++ {
		_ = PermutationOldSlUint8(ii, deck)
	}
	end = CPUtime()
	fmt.Printf("Took %f mS\n", float64((end.Utime-start.Utime)+(end.Stime-start.Stime))/float64(1000))

	fmt.Printf("PermutationNewSlUint8 [%d] * 2.5 million seeds\n", len(deck))
	start = CPUtime()
	for ii := uint64(0); ii <= uint64(2_500_000); ii++ {
		_ = PermutationNewSlUint8(ii, deck)
	}
	end = CPUtime()
	fmt.Printf("Took %f mS\n", float64((end.Utime-start.Utime)+(end.Stime-start.Stime))/float64(1000))
}
