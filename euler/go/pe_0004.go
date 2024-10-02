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

func AskAnEasierQuestion() [3]int {
	answer := [3]int{0, 0, 0}
	for ii := 999; ii > 99; ii-- {
		if answer[0] > ii*ii {
			// fmt.Println("")
			break
		}
		for kk := ii; kk > 99; kk-- {
			if answer[0] > ii*kk {
				break
			}
			test := ii * kk
			if test > answer[0] && IsPalindrome(test) {
				answer = [3]int{test, ii, kk}
			}
		}
	}
	return answer
}

func main() {
	// Test
	// fmt.Println("Is Palindrome?\t", 9009, "\t", IsPalindrome(9009))
	// fmt.Println("Is Palindrome?\t", 9090, "\t", IsPalindrome(9090))
	// fmt.Println("Is Palindrome?\t", 90609, "\t", IsPalindrome(90609))
	// fmt.Println("Is Palindrome?\t", 905434509, "\t", IsPalindrome(905434509))
	fmt.Println(AskAnEasierQuestion())
}
