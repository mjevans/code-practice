// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=43
https://projecteuler.net/minimal=43

<p>The number, $1406357289$, is a $0$ to $9$ pandigital number because it is made up of each of the digits $0$ to $9$ in some order, but it also has a rather interesting sub-string divisibility property.</p>
<p>Let $d_1$ be the $1$<sup>st</sup> digit, $d_2$ be the $2$<sup>nd</sup> digit, and so on. In this way, we note the following:</p>
<ul><li>$d_2d_3d_4=406$ is divisible by $2$</li>
<li>$d_3d_4d_5=063$ is divisible by $3$</li>
<li>$d_4d_5d_6=635$ is divisible by $5$</li>
<li>$d_5d_6d_7=357$ is divisible by $7$</li>
<li>$d_6d_7d_8=572$ is divisible by $11$</li>
<li>$d_7d_8d_9=728$ is divisible by $13$</li>
<li>$d_8d_9d_{10}=289$ is divisible by $17$</li>
</ul><p>Find the sum of all $0$ to $9$ pandigital numbers with this property.</p>

*/
/*

This is poorly defined.  What property?

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

func Euler043_test(t []uint8) bool {
	primes := [...]uint8{2, 3, 5, 7, 11, 13, 17}
	winMax := len(t) - 3
	for ii := 0; ii < winMax; ii++ {
		if 0 != (uint(t[ii])*100+uint(t[ii+1])*10+uint(t[ii+2]))%uint(primes[ii]) {
			return false
		}
	}
	return true
}

func Euler043() uint64 {
	// primes, of interest to this problem
	primes := [...]uint8{2, 3, 5, 7, 11, 13, 17}
	deck := []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var res uint64

	// 10! factorial shuffle modulus == 0
	permLimit := euler.FactorialUint64(uint64(len(deck))) - 1
	fmt.Printf("check: %d => %d\n", permLimit, euler.Uint8DigitsToUint64(euler.PermutationSlUint8(permLimit, deck), 10))
	winMax := len(deck) - 3
Euler043perm:
	for ii := uint64(0); ii < permLimit; ii++ {
		shuf := euler.PermutationSlUint8(ii, deck)
		test := shuf[1:]
		for ii := 0; ii < winMax; ii++ {
			if 0 != (uint(test[ii])*100+uint(test[ii+1])*10+uint(test[ii+2]))%uint(primes[ii]) {
				continue Euler043perm
			}
		}
		for ii := 0; ii < winMax; ii++ {
			wn := uint(test[ii])*100 + uint(test[ii+1])*10 + uint(test[ii+2])
			fmt.Printf("%d/%d = %d\t", wn, uint(primes[ii]), (uint(test[ii])*100+uint(test[ii+1])*10+uint(test[ii+2]))%uint(primes[ii]))
		}
		num := euler.Uint8DigitsToUint64(shuf, 10)
		res += num
		fmt.Printf("Found: %d\tsum: %d\n", num, res)
	}
	return res
}

//
/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 43 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 43: test [1 4 0 6 3 5 7 2 8 9]    true == 1406357289
check: 3628799 => 9876543210
406/2 = 0       63/3 = 0        635/5 = 0       357/7 = 0       572/11 = 0      728/13 = 0      289/17 = 0      Found: 1406357289       sum: 1406357289
430/2 = 0       309/3 = 0       95/5 = 0        952/7 = 0       528/11 = 0      286/13 = 0      867/17 = 0      Found: 1430952867       sum: 2837310156
460/2 = 0       603/3 = 0       35/5 = 0        357/7 = 0       572/11 = 0      728/13 = 0      289/17 = 0      Found: 1460357289       sum: 4297667445
106/2 = 0       63/3 = 0        635/5 = 0       357/7 = 0       572/11 = 0      728/13 = 0      289/17 = 0      Found: 4106357289       sum: 8404024734
130/2 = 0       309/3 = 0       95/5 = 0        952/7 = 0       528/11 = 0      286/13 = 0      867/17 = 0      Found: 4130952867       sum: 12534977601
160/2 = 0       603/3 = 0       35/5 = 0        357/7 = 0       572/11 = 0      728/13 = 0      289/17 = 0      Found: 4160357289       sum: 16695334890
Euler 43: Sub-string Divisibility: 16695334890

*/
func main() {
	//test
	test := []uint8{1, 4, 0, 6, 3, 5, 7, 2, 8, 9}
	fmt.Printf("Euler 43: test %v\t%t == %d\n", test, Euler043_test(test[1:]), euler.Uint8DigitsToUint64(test, 10))

	//run
	fmt.Printf("Euler 43: Sub-string Divisibility: %d\n", Euler043())
}
