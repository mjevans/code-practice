// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=31
https://projecteuler.net/minimal=31


<p>In the United Kingdom the currency is made up of pound (£) and pence (p). There are eight coins in general circulation:</p>
<blockquote>1p, 2p, 5p, 10p, 20p, 50p, £1 (100p), and £2 (200p).</blockquote>
<p>It is possible to make £2 in the following way:</p>
<blockquote>1×£1 + 1×50p + 2×20p + 1×5p + 1×2p + 3×1p</blockquote>
<p>How many different ways can £2 be made using any number of coins?</p>




*/

import (
	// "bufio"
	// "bitvector"
	// "euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler031(change int) int {
	var combos int
	// coins, lim := [8]int{1,2,5,10,20,50,100,200}, [8]int{}
	for c200 := 0; c200*200 <= change; c200++ {
		for c100 := 0; c200*200+c100*100 <= change; c100++ {
			for c050 := 0; c200*200+c100*100+c050*50 <= change; c050++ {
				for c020 := 0; c200*200+c100*100+c050*50+c020*20 <= change; c020++ {
					for c010 := 0; c200*200+c100*100+c050*50+c020*20+c010*10 <= change; c010++ {
						for c005 := 0; c200*200+c100*100+c050*50+c020*20+c010*10+c005*5 <= change; c005++ {
							for c002 := 0; c200*200+c100*100+c050*50+c020*20+c010*10+c005*5+c002*2 <= change; c002++ {
								// for c001 := 0; c200*200+c100*100+c050*50+c020*20+c010*10+c005*5+c002*2+c001 <= change; c001++ {
								c001 := change - (c200*200 + c100*100 + c050*50 + c020*20 + c010*10 + c005*5 + c002*2)
								if change != c200*200+c100*100+c050*50+c020*20+c010*10+c005*5+c002*2+c001 {
									fmt.Printf("Logic check failed, why? %d != %d\n", change, c200*200+c100*100+c050*50+c020*20+c010*10+c005*5+c002*2+c001)
								}
								combos++
								if 11 > change {
									fmt.Printf("%d:\t", combos)
									/*
										if 0 < c200 {
											fmt.Printf("200p%d\t", c200)
										} else {
											fmt.Printf("\t")
										}
										if 0 < c100 {
											fmt.Printf("100p%d\t", c100)
										} else {
											fmt.Printf("\t")
										}
										if 0 < c050 {
											fmt.Printf("050p%d\t", c050)
										} else {
											fmt.Printf("\t")
										}
										if 0 < c020 {
											fmt.Printf("020p%d\t", c020)
										} else {
											fmt.Printf("\t")
										}*/
									if 0 < c010 {
										fmt.Printf("010p%d\t", c010)
									} else {
										fmt.Printf("\t")
									}
									if 0 < c005 {
										fmt.Printf("005p%d\t", c005)
									} else {
										fmt.Printf("\t")
									}
									if 0 < c002 {
										fmt.Printf("002p%d\t", c002)
									} else {
										fmt.Printf("\t")
									}
									if 0 < c001 {
										fmt.Printf("001p%d\n", c001)
									} else {
										fmt.Printf("\n")
									}
								}
								// }
							}
						}
					}
				}
			}
		}
	}
	return combos
}

func Euler031_Alt(change uint16) uint64 {
	var combos uint64
	const ccount = 8
	coins, vals := make([]uint16, 8, 8), []uint16{1, 2, 5, 10, 20, 50, 100, 200} // {1,2,5,10,20,50,100,200} // {200,100,50,20,10,5,2,1}
	slotLim := func(change uint16, coins, vals []uint16) uint16 {
		up := uint16(0)
		for ii := 1; ii < len(vals); ii++ {
			up += coins[ii] * vals[ii]
		}
		return (change - up) / vals[0]
	}

	for {
		cur := uint16(0)
		for ii := 1; ii < ccount; ii++ {
			cur += coins[ii] * vals[ii]
		}
		// The rest are single unit pieces for this currency
		if cur <= change {
			combos++
			if change < 11 {
				coins[0] = change - cur
				fmt.Printf("p%d\t%dc\t\t%v\n", change, combos, coins)
				coins[0] = 0
			}
		}
		// increase coins...
		for ii := 1; ii < ccount; ii++ {
			lim := slotLim(change, coins[ii:], vals[ii:])
			if coins[ii] < lim {
				coins[ii]++
				break
			} else {
				if ccount == ii+1 {
					// fmt.Println("exit: coin iter")
					return combos
				}
				coins[ii] = 0
			}
		}
	}
	return combos
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 31 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

An obvious mistake when I looked at the code again, 'any number of pennies less than the correct number' is not proper change.

1:                              001p10
2:                      002p1   001p8
3:                      002p2   001p6
4:                      002p3   001p4
5:                      002p4   001p2
6:                      002p5
7:              005p1           001p5
8:              005p1   002p1   001p3
9:              005p1   002p2   001p1
10:             005p2
11:     010p1
Test 10p:        11
Euler031:        73682
*/
func main() {
	//test
	fmt.Println("Test p5:\t", Euler031(5))
	fmt.Println("Test p8:\t", Euler031(8))
	fmt.Println("Test p10:\t", Euler031(10))
	fmt.Println("Test p20:\t", Euler031(20))
	fmt.Println("Test p50:\t", Euler031(50))
	fmt.Println("Test p100:\t", Euler031(100))
	fmt.Println("Test p200:\t", Euler031(200))

	//run
	permu := Euler031(200)
	fmt.Println("Euler031:\t", permu)
	fmt.Println("Also a more generic variation...")

	fmt.Println("Test p5:\t", Euler031_Alt(5))
	fmt.Println("Test p8:\t", Euler031_Alt(8))
	fmt.Println("Test p10:\t", Euler031_Alt(10))
	fmt.Println("Test p20:\t", Euler031_Alt(20))
	fmt.Println("Test p50:\t", Euler031_Alt(50))
	fmt.Println("Test p100:\t", Euler031_Alt(100))
	fmt.Println("Test p200:\t", Euler031_Alt(200))

	//run
	p2 := Euler031_Alt(200)
	fmt.Println("Euler031_Alt:\t", p2)

}
