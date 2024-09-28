// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=31

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
								for c001 := 0; c200*200+c100*100+c050*50+c020*20+c010*10+c005*5+c002*2+c001 <= change; c001++ {
									combos++
									if 21 > change {
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
								}
							}
						}
					}
				}
			}
		}
	}
	return combos
}

//	for ii in */*.go ; do go fmt "$ii" ; done ; for ii in 31 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done
/*

37:             005p1
38:             005p1           001p1
39:             005p1           001p2
40:             005p1           001p3
41:             005p1           001p4
42:             005p1           001p5
43:             005p1   002p1
44:             005p1   002p1   001p1
45:             005p1   002p1   001p2
46:             005p1   002p1   001p3
47:             005p1   002p2
48:             005p1   002p2   001p1
49:             005p2
50:     010p1
Test 10p:        50
Euler031:        2886726


 */
func main() {
	//test
	fmt.Println("Test 10p:\t", Euler031(10))

	//run
	permu := Euler031(200)
	fmt.Println("Euler031:\t", permu)

}
