// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=4
https://projecteuler.net/minimal=4

<p>A palindromic number reads the same both ways. The largest palindrome made from the product of two $2$-digit numbers is $9009 = 91 \times 99$.</p>
<p>Find the largest palindrome made from the product of two $3$-digit numbers.</p>

*/

import (
	"fmt"
	"strings"
	// "os" // os.Stdout
)

func Factor(primes []int, num int) []int {
	// Public school factoring algorithm from memory...

	// With a list of known primes, the largest number that can be factored is Pn * Pn
	for ; nil == primes || num > primes[len(primes)-1]*primes[len(primes)-1]; primes = GetPrimes(primes, 0) {
		// fmt.Println(len(primes), primes[len(primes)-1])
	}

	ret := []int{}
	for _, prime := range primes {
		for ; 0 == num%prime; num /= prime {
			ret = append(ret, prime)
		}
		if num < prime*prime {
			break
		} // break if no more prime factors are possible
	}
	if 0 < num {
		ret = append(ret, num)
	}
	return ret
}

func GetPrimes(primes []int, primehunt int) []int {
	if nil == primes {
		primes = []int{2, 3, 5, 7, 11, 13, 17, 19}
	}
	// Semi-arbitrary expansion target, find 8 more primes (8, 16, 32, 64 it'll fit within the append growth algo)
	if primehunt < 1 {
		primehunt = 8
	}
PrimeHunt:
	for ; 0 < primehunt; primehunt-- {
		for ii := primes[len(primes)-1] + 1; ; ii++ {
			result := Factor(primes, ii)
			if 1 == len(result) && primes[len(primes)-1] < result[0] {
				//fmt.Println("Found Prime:\t", result[0])
				primes = append(primes, result[0])
				continue PrimeHunt // I could break once, but this documents the intent
			}
		}
	}
	return primes
}

func PrintFactors(factors []int) {
	// Join only takes []string s? fff
	strFact := make([]string, len(factors), len(factors))
	for ii, val := range factors {
		strFact[ii] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(strFact, ", "))
}

func ListMul(scale []int) int {
	ret := 1
	for _, val := range scale {
		ret *= val
	}
	return ret
}

func Find3DigitNumbers(factors []int) (bool, []int) {
	// Winging it... Zipper back to front, since the list is sorted already uneven splits go to the smaller largest prime factor
	a := make([]int, 0, len(factors))
	b := make([]int, 0, len(factors))
	ii := len(factors) - 1
	for ; ii > 2; ii -= 2 {
		a = append(a, factors[ii])
		b = append(b, factors[ii])
	}
	if 1 == ii {
		b = append(b, factors[ii])
	}
	// Two lists of factors...
	ret = make([]int, 2)
	for {
		ret[0] = ListMul(a)
		ret[1] = ListMul(b)
		if 100 <= ret[0] && ret[0] <= 999 && 100 <= ret[1] && ret[1] <= 999 {
			return true, ret
		}
		if (ret[0] < 100 && ret[1] < 100) || (ret[0] > 999 && ret[1] > 999) {
			return false, nil
		}
		// Plan an Adjustment, division is usually expensive so minimize
		if ret[0] < ret[1] { // If B is bigger, swap the lists
			tmp := ret[0]
			tmpS := a
			ret[0] = ret[1]
			a = b
			ret[1] = tmp
			b = tmpS
		}
		larger := ret[0] / 100
		smaller := 1000 / ret[1] // 500 in worst case
		if smaller < 2 || larger * ret[1] < 1000 {
			return false, nil // Either unpackable (smaller too large) or too small to satisfy
		}
		// Search for the answer FIXME:

	}
}

func AskASillyQuestion() {
	var primes []int
	for primes = GetPrimes(nil, 0); primes[len(primes)-1] < 1000; primes = GetPrimes(primes, 0) {
	}
	// Construct a 6 digit number that MIGHT be a palindrome, biased from largest to smallest.
	for i1 := 9; i1 >= 0; i1-- {
		for i2 := 9; i2 >= 0; i2-- {
			for i3 := 9; i3 >= 0; i3-- {
				maybePal = i1*100000 + i2*10000 + i3*1000 + i3*100 + i2*10 + i1
				factors := Factor(primes, maybePal)
				if ok, nums := Find3DigitNumbers(factors); ok && len(nums) == 2 {
					strNums := make([]string, len(nums), len(nums))
					for ii, val := range nums {
						strNums[ii] = fmt.Sprint(val)
					}
					fmt.Print(maybePal, ": ", strings.Join(strNums, " "))
				}
			}
		}
	}
}

func main() {
	// Test
	// PrintFactors(Factor(nil, 13195))
	// matches factor
	PrintFactors(Factor(nil, 600851475143))
}
