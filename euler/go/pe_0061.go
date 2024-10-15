// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=61
https://projecteuler.net/minimal=61

<p>Triangle, square, pentagonal, hexagonal, heptagonal, and octagonal numbers are all figurate (polygonal) numbers and are generated by the following formulae:</p>
<table><tr><td>Triangle</td>
<td> </td>
<td>$P_{3,n}=n(n+1)/2$</td>
<td> </td>
<td>$1, 3, 6, 10, 15, \dots$</td>
</tr><tr><td>Square</td>
<td> </td>
<td>$P_{4,n}=n^2$</td>
<td> </td>
<td>$1, 4, 9, 16, 25, \dots$</td>
</tr><tr><td>Pentagonal</td>
<td> </td>
<td>$P_{5,n}=n(3n-1)/2$</td>
<td> </td>
<td>$1, 5, 12, 22, 35, \dots$</td>
</tr><tr><td>Hexagonal</td>
<td> </td>
<td>$P_{6,n}=n(2n-1)$</td>
<td> </td>
<td>$1, 6, 15, 28, 45, \dots$</td>
</tr><tr><td>Heptagonal</td>
<td> </td>
<td>$P_{7,n}=n(5n-3)/2$</td>
<td> </td>
<td>$1, 7, 18, 34, 55, \dots$</td>
</tr><tr><td>Octagonal</td>
<td> </td>
<td>$P_{8,n}=n(3n-2)$</td>
<td> </td>
<td>$1, 8, 21, 40, 65, \dots$</td>
</tr></table><p>The ordered set of three $4$-digit numbers: $8128$, $2882$, $8281$, has three interesting properties.</p>
<ol><li>The set is cyclic, in that the last two digits of each number is the first two digits of the next number (including the last number with the first).</li>
<li>Each polygonal type: triangle ($P_{3,127}=8128$), square ($P_{4,91}=8281$), and pentagonal ($P_{5,44}=2882$), is represented by a different number in the set.</li>
<li>This is the only set of $4$-digit numbers with this property.</li>
</ol><p>Find the sum of the only ordered set of six cyclic $4$-digit numbers for which each polygonal type: triangle, square, pentagonal, hexagonal, heptagonal, and octagonal, is represented by a different number in the set.</p>

*/
/*

Go 1.21 adds 'clear()' for purging maps / lists but that's every element.
delete(map, key) is in base, and if map is nil or key is absent it returns without error


0061 successfully confused me about the relation of the numbers and thus what the question was:

This is the mathematical relation they wanted to convey:

TriN(127)=8128 PentN(44)=2882 SqN(91)=8281 => 8128 2882 8281

In _any order_ there are three numbers among Ngon numbers (3..5) where the results are all 4 digits, and the lower 2 digits of one result match exactly one pair of 2 digits from another result, such that there is one unbroken chain.

Oh; this is a graph / edge theory problem?

Something like... 5 nodes must be visited, and rules for traversal?

My folly's slightly tightened limit test can also work if modified.


What does a solution look like?

A set of nodes / islands in a given order. ( Have to start iterating on _something_ and that seemed like the best match. )
+
A sequence of numbers...

Third or fourth attempt; I can't discard the numbers.
So:

A set of nodes / islands in a given order. ( Have to start iterating on _something_ and that seemed like the best match. )
+
Each Node's high (inbound) contribution in order.
[][Ngon]uint16
So a list of solutions evaluates this node's matches, but doesn't add that match from the combo list if it matches a prior island's.


I also now realize, offhand, I've never done one of those write your own toy database exercises.  Open Source libre and zero cost solutions abound, be it MySQL/MariaDB/whatever PostgreSQL or SQLite to name a few.  It's surprisingly easy to have a good general idea of how a database works, without connecting that concept to graph theory.


*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

/*
type NgonNum struct {
	Num, Split, N uint16
}

func NewNgonNum(Num, Split, N uint64) NgonNum {
	return NgonNum{Num: uint16(Num), Split: uint16(Split), N: uint16(N)}
}

func (ng NgonNum) GetNum() uint64  { return uint64(ng.Num) }
func (ng NgonNum) GetHigh() uint64 { return uint64(ng.Num / ng.Split) }
func (ng NgonNum) GetLow() uint64  { return uint64(ng.Num % ng.Split) }
func (ng NgonNum) GetN() uint64    { return uint64(ng.N) }

type NgonLink struct {
	Num      uint16
	Src, Dst uint8
}
*/

func Euler0061(min, max, modsplit, base uint64) uint64 {
	var n, ii, jj, added, passes, sum uint64
	modsplitSZ := uint16(modsplit)
	_ = modsplitSZ
	lowMin := modsplit / base
	// Make _slice_ of maps (upper half) to _slice_ of lower half matches
	ng := make([]map[uint16][]uint16, 8-3+1) // NgonMax - NgonMin + 1 for the list size
	for jj = 3; jj <= 8; jj++ {
		ng[jj-3] = make(map[uint16][]uint16)
		fmt.Printf("Searching for Ngon %d : %4d .. %4d", jj, min, max)
		added = 0
		for ii = euler.NgonNumberReverseFloor(min, jj); ; ii++ {
			n = euler.NgonNumber(ii, jj)
			if n < min || n%modsplit < lowMin {
				continue
			}
			if n > max {
				break
			}
			hh, ll := uint16(n/modsplit), uint16(n%modsplit)
			ng[jj-3][hh] = append(ng[jj-3][hh], uint16(ll))
			added++
		}
		fmt.Printf("\tAdded %d\n", added)
	}

	const NgonCount = 6

	// Every ['shape island'] now has ['(high num) inward piers'] which might match A List [] of maybe routes to 'Other Islands (low num)'
	// ? Reduce - Prune all edges that don't lead anywhere
	// var edges [NgonCount][NgonCount][]uint16
	// var edge uint64
	// edges isn't quite correct, but it helped inspire me to a better path
	// 3-3 <= 8-3
	passes = 0
	for del := 1; 0 < del; {
		// edges = [NgonCount][NgonCount][]uint16{}
		// edge = 0
		del = 0 // reset
		for ii = 0; ii < NgonCount; ii++ {
			for iiH, sl := range ng[ii] {
				slmx := len(sl)
				sllink := false
				for slii := 0; !sllink && slii < slmx; slii++ {
					link := false
					lnum := sl[slii]
					if lnum == uint16(modsplit) {
						continue
					}
					for jj = 0; jj < NgonCount; jj++ {
						if ii == jj {
							continue
						}
						if _, ok := ng[jj][lnum]; ok {
							link = true
							sllink = true
							// bridge := iiH*modsplitSZ + lnum
							// edges[ii][jj] = append(edges[ii][jj], bridge)
							// edge++
							break
						}
					}
					if false == link {
						sl[slii] = uint16(modsplit) // invalid value
					}
				}
				if false == sllink {
					delete(ng[ii], iiH)
					del++
				}
			}
		}
		passes++
		fmt.Printf("Reduce pass %d pruned %d\n", passes, del)
	}
	// fmt.Println(edges[5][0])
	//
	// _slice_ :: kn == maps (upper half) to sln == _slice_ of lower half matches

	ValidateNgonSet := func(node []uint8, nums [NgonCount]uint16) uint64 {
		var sum uint64
		ii := 0
		for ; ii < NgonCount-1; ii++ {
			t := uint64(nums[ii]*modsplitSZ + nums[ii+1])
			tr := euler.NgonNumberReverseFloor(t, uint64(node[ii]+3))
			tgon := euler.NgonNumber(tr, uint64(node[ii]+3))
			if tgon != t {
				fmt.Printf("BUG%d: Did not match:\tNgon%d(%d) = %d = %d\n", ii, node[ii]+3, tr, t, tgon)
				return 0
			}
			sum += uint64(tgon)
			fmt.Printf("\tNgon%d(%d) = %d\t", node[ii]+3, tr, t)
		}
		t := uint64(nums[ii]*modsplitSZ + nums[0])
		tr := euler.NgonNumberReverseFloor(t, uint64(node[ii]+3))
		tgon := euler.NgonNumber(tr, uint64(node[ii]+3))
		if tgon != t {
			fmt.Printf("BUG%d: Did not match:\tNgon%d(%d) = %d = %d\n", ii, node[ii]+3, tr, t, tgon)
			return 0
		}
		sum += uint64(tgon)
		fmt.Printf("\tNgon%d(%d) = %d\n", node[ii]+3, tr, t)
		return sum
	}

	// What does a solution look like?

	node := make([]uint8, NgonCount)

	for node[0] = 0; node[0] < NgonCount; node[0]++ {
		for node[1] = 0; node[1] < NgonCount; node[1]++ {
			if -1 != euler.UnsortedSearchSlice(node[0:1], node[1]) {
				continue
			}
			for node[2] = 0; node[2] < NgonCount; node[2]++ {
				if -1 != euler.UnsortedSearchSlice(node[0:2], node[2]) {
					continue
				}
				for node[3] = 0; node[3] < NgonCount; node[3]++ {
					if -1 != euler.UnsortedSearchSlice(node[0:3], node[3]) {
						continue
					}
					for node[4] = 0; node[4] < NgonCount; node[4]++ {
						if -1 != euler.UnsortedSearchSlice(node[0:4], node[4]) {
							continue
						}
					Euler0061NextCombo:
						for node[5] = 0; node[5] < NgonCount; node[5]++ {
							if -1 != euler.UnsortedSearchSlice(node[0:5], node[5]) {
								continue
							}

							// fmt.Printf("Trying path order: %v\n", node)

							// Initial 'row' set.
							rows := make(map[uint16][NgonCount]uint16)
							count := 0
							for iiHi, sl := range ng[node[0]] {
								slmx := len(sl)
								for slii := 0; slii < slmx; slii++ {
									if _, ok := ng[node[1]][sl[slii]]; ok {
										rows[uint16(count)] = [NgonCount]uint16{iiHi, sl[slii]}
										count++
									}
								}
							}
							if 0 == count {
								continue
							}
							idx := uint16(count)

							// rows is like a tuple of high halfs for each ngon, followed by the low half of the next ngon.
							for col := 1; col < NgonCount; col++ {
								count = 0
								colLow := col + 1
								precheck := false
								if NgonCount == colLow {
									colLow = 0
									precheck = true
								}

								added := make([][NgonCount]uint16, 0)
								for k, val := range rows {
									// This will harmlessly be false if asked for 0 since that is not a valid number for this problem
									keepRow := false
									targetLows := ng[node[col]][val[col]]
									tlmx := len(targetLows)
									for tl := 0; tl < tlmx; tl++ {
										if _, ok := ng[node[colLow]][targetLows[tl]]; ok {
											if precheck && val[colLow] != targetLows[tl] {
												continue
											}
											if 0 == tl {
												cval := val
												cval[colLow] = targetLows[tl]
												rows[k] = cval
											} else {
												cval := val
												cval[colLow] = targetLows[tl]
												added = append(added, cval)
												// Might have caused Dupes. deferred insert
											}
											count++
											keepRow = true
										}
									}
									if false == keepRow {
										delete(rows, k)
									}
								}
								for _, v := range added {
									idx++
									rows[idx] = v
								}

								if 0 == count {
									continue Euler0061NextCombo
								}
							}
							for ii, _ := range rows {
								sum += ValidateNgonSet(node, rows[ii])
								fmt.Printf("Solution: %v %v\n", node, rows[ii])
							}
							return sum
						}
					}
				}
			}
		}
	}

	return sum
}

// Oops, I might have been a bit too tired when I initially read the problem and didn't realize that the 'ordered set of numbers' was provided in an arbitrary sort (how their ends match up) rather than how I would have written the statement...
// TriN(127)=8128 PentN(44)=2882 SqN(91)=8281 => 8128 2882 8281
// So this algorithm 'correctly' reaches a state where there isn't a solution, because there isn't a solution.
func Incorrect61_ProblemDefinition(min, max, modsplit, base uint64) uint64 {
	lowMin := modsplit / base

	// R = NGon(ii) => map[ii]R
	// Note: collisions on %modsplit are possible
	m3 := make(map[uint16]uint16)
	m4 := make(map[uint16]uint16)
	m5 := make(map[uint16]uint16)
	m6 := make(map[uint16]uint16)
	m7 := make(map[uint16]uint16)
	m8 := make(map[uint16]uint16)
	// reverse lookup map, note multiple outputs
	m7r := make(map[uint16][]uint16)
	ii := euler.HeptagonalNumberReverseFloor(min)
	for ; ; ii++ {
		n := euler.HeptagonalNumber(ii)
		if n < min || n%modsplit < lowMin {
			//if n >= min {
			//	fmt.Printf("init: Culled low: %d\n", n)
			//}
			continue
		}
		if n > max {
			break
		}
		m7[uint16(ii)] = uint16(n)
		// key := uint16(n / modsplit)
		// m7r[key] = append(m7r[key], uint16(n))
	}
	// sorted not necessary to filter
	// 7 -> 8
	for k, v := range m7 {
		ok := false
		fl := (uint64(v) % modsplit) * modsplit
		cl := fl + modsplit - 1
		fl += lowMin
		ii = euler.OctagonalNumberReverseFloor(fl)
		// fmt.Printf("8: [%4d - %4d] %d (%d)", fl, cl, ii, v)
		for ; ; ii++ {
			n := euler.OctagonalNumber(ii)
			if n < fl {
				//if n >= fl-lowMin {
				//	fmt.Printf("4: Culled low: %d\n", n)
				//}
				continue
			}
			if n > cl {
				break
			}
			// fmt.Printf("\t%d", n)
			m8[uint16(ii)] = uint16(n)
			ok = true
		}
		// fmt.Printf("\t%t\n", ok)
		if !ok {
			delete(m7, k)
		} else {
			key := v / uint16(modsplit)
			m7r[key] = append(m7r[key], uint16(k))
		}
	}
	fmt.Printf("Initial Filter:\n7r: %v\n7: %v\n8: %v\n", m7r, m7, m8)
	// 8 -> 3
	for k, v := range m8 {
		ok := false
		fl := (uint64(v) % modsplit) * modsplit
		cl := fl + modsplit - 1
		fl += lowMin
		ii = euler.TriangleNumberReverseFloor(fl)
		// fmt.Printf("Triangle: [%4d - %4d] %d (%d)", fl, cl, ii, v)
		for ; ; ii++ {
			n := euler.TriangleNumber(ii)
			if n < fl {
				//if n >= fl-lowMin {
				//	fmt.Printf("3: Culled low: %d\n", n)
				//}
				continue
			}
			if n > cl {
				break
			}
			// fmt.Printf("\t%d", n)
			m3[uint16(ii)] = uint16(n)
			ok = true
		}
		// fmt.Printf("\t%t\n", ok)
		if !ok {
			delete(m8, k)
		}
	}
	// 3 -> 4
	for k, v := range m3 {
		ok := false
		fl := (uint64(v) % modsplit) * modsplit
		cl := fl + modsplit - 1
		fl += lowMin
		ii = euler.SquareNumberReverseFloor(fl)
		fmt.Printf("Square: [%4d - %4d] %d (%d)", fl, cl, ii, v)
		for ; ; ii++ {
			n := euler.SquareNumber(ii)
			if n < fl {
				//if n >= fl-lowMin {
				//	fmt.Printf("4: Culled low: %d\n", n)
				//}
				continue
			}
			if n > cl {
				break
			}
			fmt.Printf("\t%d", n)
			m4[uint16(ii)] = uint16(n)
			ok = true
		}
		fmt.Printf("\t%t\n", ok)
		if !ok {
			delete(m3, k)
		}
	}
	// 4 -> 5
	for k, v := range m4 {
		ok := false
		fl := (uint64(v) % modsplit) * modsplit
		cl := fl + modsplit - 1
		fl += lowMin
		ii = euler.PentagonalNumberReverseFloor(fl)
		fmt.Printf("Pentagonal: [%4d - %4d] %d (%d)", fl, cl, ii, v)
		for ; ; ii++ {
			n := euler.PentagonalNumber(ii)
			if n < fl {
				//if n >= fl-lowMin {
				//	fmt.Printf("5: Culled low: %d\n", n)
				//}
				continue
			}
			if n > cl {
				break
			}
			fmt.Printf("\t%d", n)
			m5[uint16(ii)] = uint16(n)
			ok = true
		}
		fmt.Printf("\t%t\n", ok)
		if !ok {
			delete(m4, k)
		}
	}
	// 5 -> 6 + 6->7r match filter
	for k, v := range m5 {
		ok := false
		fl := (uint64(v) % modsplit) * modsplit
		cl := fl + modsplit - 1
		fl += lowMin
		ii = euler.HexagonalNumberReverseFloor(fl)
		fmt.Printf("Hexagonal: [%4d - %4d] %d (%d)\n", fl, cl, ii, v)
		for ; ; ii++ {
			n := euler.HexagonalNumber(ii)
			if n < fl {
				if n >= fl-lowMin {
					fmt.Printf("6: Culled low: %d\n", n)
				}
				continue
			}
			if n > cl {
				break
			}
			key7 := uint16(n % modsplit)
			fmt.Printf("Hexagonal: want m7r %d for %d\n", key7, n)
			if _, ok7r := m7r[key7]; ok7r {
				m6[uint16(ii)] = uint16(n)
				ok = true
			}
		}
		fmt.Printf("\t%t\n", ok)
		if !ok {
			delete(m5, k)
		}
	}
	/*
		// 6 -> 7 + 7->8r match filter
		for k, v := range m6 {
			ok := false
			fl := (uint64(v) % modsplit) * modsplit
			cl := fl + modsplit - 1
			fl += lowMin
			ii = euler.HeptagonalNumberReverseFloor(fl)
			fmt.Printf("Heptagonal: [%4d - %4d] %d (%d)\n", fl, cl, ii, v)
			for ; ; ii++ {
				n := euler.HeptagonalNumber(ii)
				if n < fl {
					if n >= fl-lowMin {
						fmt.Printf("7: Culled low: %d\n", n)
					}
					continue
				}
				if n > cl {
					break
				}
				key8 := uint16(n % modsplit)
				fmt.Printf("Heptagonal: want m8r %d for %d\n", key8, n)
				if _, ok8r := m8r[key8]; ok8r {
					m7[uint16(ii)] = uint16(n)
					ok = true
				}
			}
			// fmt.Printf("\t%t\n", ok)
			if !ok {
				delete(m6, k)
			}
		}
	*/

	var sum uint64
	fmt.Printf("Filtered set total: %d\n3: %v\n4: %v\n5: %v\n6: %v\n7: %v\n8: %v\n8r: %v\n", sum, m3, m4, m5, m6, m7, m8, m7r)
	return sum
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 61 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Searching for Ngon 3 : 1000 .. 9999     Added 88
Searching for Ngon 4 : 1000 .. 9999     Added 53
Searching for Ngon 5 : 1000 .. 9999     Added 47
Searching for Ngon 6 : 1000 .. 9999     Added 44
Searching for Ngon 7 : 1000 .. 9999     Added 40
Searching for Ngon 8 : 1000 .. 9999     Added 30
Reduce pass 1 pruned 11
Reduce pass 2 pruned 4
Reduce pass 3 pruned 0
        Ngon3(128) = 8256               Ngon4(75) = 5625                Ngon7(32) = 2512                Ngon8(21) = 1281                Ngon6(64) = 8128                Ngon5(44) = 2882
Solution: [0 1 4 5 3 2] [82 56 25 12 81 28]
Euler 61: Cyclical Figurate Numbers: 28684


*/
func main() {
	var a uint64
	//test
	// a = Euler0061(1000, 9999, 100, 10)

	//run
	a = Euler0061(1000, 9999, 100, 10)
	fmt.Printf("Euler 61: Cyclical Figurate Numbers: %d\n", a)
}
