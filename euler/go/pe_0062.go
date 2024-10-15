// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=62
https://projecteuler.net/minimal=62

<p>The cube, $41063625$ ($345^3$), can be permuted to produce two other cubes: $56623104$ ($384^3$) and $66430125$ ($405^3$). In fact, $41063625$ is the smallest cube which has exactly three permutations of its digits which are also cube.</p>
<p>Find the smallest cube for which exactly five permutations of its digits are cube.</p>

*/
/*



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

func Euler0062(exact_perms, base uint64) uint64 {
	var ii, ii3, combo, res, maxseen, ckey uint64
	cubes := make(map[uint64][]uint64)
	// var count map[uint64]uint64
	ii = 1
	for {
		ii++
		ii3 = ii * ii * ii
		ii3sl := euler.Uint64ToDigitsUint8(ii3, base)
		ckey = euler.Uint8DigitsToUint64(euler.Uint8CopyInsertSort(ii3sl), base)
		cubes[ckey] = append(cubes[ckey], ii)
		seen := uint64(len(cubes[ckey]))
		if maxseen < seen {
			fmt.Printf("New Max seen: @ %d : %d < %d\t%v\n", ii, maxseen, seen, cubes[ckey])
			maxseen = seen
		}
		if 0 == ii&0xFFF {
			fmt.Printf("Iter %d max seen %d\n", ii, maxseen)
		}
		if seen == exact_perms {
			res = cubes[ckey][0]
			res *= res * res
			_ = combo
			break

			// Possible, but a very bad idea
			// Length 12 has max combos: 479001600
			/*
				ii3len := len(ii3sl)
				count = make(map[uint64]uint64)
				count[ii3] = ii
				res = ii3
				combomax := euler.FactorialUint64(uint64(ii3len))
				fmt.Printf("Length %d has max combos: %d\n", ii3len, combomax)
				if combomax < exact_perms {
					continue
				}
				for combo = 1; combo < combomax; combo++ {
					pp3 := euler.Uint8DigitsToUint64(euler.PermutationSlUint8(combo, ii3sl), base)
					// Smaller numbers covered by prior iterations
					if _, there := count[pp3]; there {
						continue
					}
					pp := uint64(math.Pow(float64(pp3), 1.0/3.0))
					if pp3 == pp*pp*pp {
						count[pp3] = pp
						if pp3 < res {
							res = pp3
						}
					}
				}
				if uint64(len(count)) == exact_perms {
					break
				}
				if uint64(len(count)) > exact_perms {
					fmt.Printf("Too many %d %v\n", len(count), count)
				}
			*/
		}
	}
	fmt.Printf("Found %v\n", cubes[ckey])
	return res
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 62 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

New Max seen: @ 2 : 0 < 1       [2]
New Max seen: @ 8 : 1 < 2       [5 8]
New Max seen: @ 405 : 2 < 3     [345 384 405]
Found [345 384 405]
Euler 62: TEST Cubic Permutations: 41063625
New Max seen: @ 2 : 0 < 1       [2]
New Max seen: @ 8 : 1 < 2       [5 8]
New Max seen: @ 405 : 2 < 3     [345 384 405]
New Max seen: @ 2010 : 3 < 4    [1002 1020 2001 2010]
Iter 4096 max seen 4
Iter 8192 max seen 4
New Max seen: @ 8384 : 4 < 5    [5027 7061 7202 8288 8384]
Found [5027 7061 7202 8288 8384]
Euler 62: Cubic Permutations: 127035954683


*/
func main() {
	var a uint64
	//test
	a = Euler0062(3, 10)
	fmt.Printf("Euler 62: TEST Cubic Permutations: %d\n", a)

	//run
	a = Euler0062(5, 10)
	fmt.Printf("Euler 62: Cubic Permutations: %d\n", a)
}
