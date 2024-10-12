// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=56
https://projecteuler.net/minimal=56

<p>A googol ($10^{100}$) is a massive number: one followed by one-hundred zeros; $100^{100}$ is almost unimaginably large: one followed by two-hundred zeros. Despite their size, the sum of the digits in each number is only $1$.</p>
<p>Considering natural numbers of the form, $a^b$, where $a, b \lt 100$, what is the maximum digital sum?</p>


*/
/*

"math/big" can do numbers that large

However, when should I stop the backward search?

Review:

1 * 10 = 10	(2 digits)
10 * 10 = 100	(3 digits)
100 * 10 = 1000 (4 digits)
100 * 100 = 10000 (5 digits)

99 * 99 = 4 digits
9 * 99 = 3 digits
9 * 9 = 2 digits

A 2 digit base can _at most_ add 2 digits (or 18 for a digital sum) for every time it's power is evaluated.



*/

import (
	// "bufio"
	// "euler"
	"fmt"
	// "math"
	"math/big"
	// "slices" // Doh not in 1.19
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler056(basemax, powmax, basemod int64) (string, int64) {
	var bestSumDig, sum, possible, digcount int64
	var bestNumStr string
	debug := false
	biBaseMod := big.NewInt(int64(basemod))
	biBase, biPowB, biZero, biResMod := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	digcount = 0
	for t := basemax; 0 < t; t /= basemod {
		digcount++
	}
	possible = digcount * powmax * (basemod - 1)

	basefloor := basemax / basemod
	pow := powmax
	// this will fail before pow reaches 1 let alone 0
	for bestSumDig < possible {
		for base := basemax; basefloor < base; base-- {
			biBase = biBase.SetInt64(base)
			biPowB = biPowB.SetInt64(base)
			for bigPow := pow - 1; 0 < bigPow; bigPow-- {
				biBase.Mul(biBase, biPowB)
			}
			if debug {
				fmt.Printf("%d^%d\t= %s\t+++", base, pow, biBase.Text(int(basemod)))
			}
			biPowB.Set(biBase)
			sum = 0
			// while biBase > 0
			for 1 == biPowB.Cmp(biZero) {
				biPowB.DivMod(biPowB, biBaseMod, biResMod)
				sum += biResMod.Int64()
			}
			if debug {
				fmt.Printf("\t%d\n", sum)
			}
			if bestSumDig < sum {
				bestSumDig = sum
				bestNumStr = biBase.Text(int(basemod))
				fmt.Printf("New Best Sum: %d^%d = (%s) %d\n", base, pow, bestNumStr, bestSumDig)
			}
			if 1 != sum && sum*basemod < bestSumDig {
				if debug {
					fmt.Printf("BREAK:\n%d = sum * basemod\n%d = bestSumDig\n", sum*basemod, bestSumDig)
				}
				break
			} else if debug {
				// fmt.Printf("Compared:\n%d = sum * basemod\n%d = bestSumDig\n", sum*basemod, bestSumDig)
			}
		}
		pow--
		possible -= (basemod - 1)
	}
	return bestNumStr, bestSumDig
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 56 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Just playing around 100, 100
Euler 56: Powerful Digit Sum (3848960788934848611927795802824596789608451156087366034658627953530148126008534258032267383768627487094610968554286692697374726725853195657679460590239636893953692985541958490801973870359499) 972

New Best Sum: 99^99 = (369729637649726772657187905628805440595668764281741102430259972423552570455277523421410650010128232727940978889548326540119429996769494359451621570193644014418071060667659301384999779999159200499899) 936
New Best Sum: 89^99 = (9763615267845346070648947915415150318379739756055909464947086388973389093238094962788773085346058287871726868180109913308128757543020083607530804980833255899439481274385221883982904798463876809) 953
New Best Sum: 79^99 = (73296297596894671194679467281666017876936144887237684747595248746205105058136664662597084059678669544148492467642033016554778566668452647644046940902679823589082567348293472699134109889519) 955
New Best Sum: 94^98 = (23255712658709810541561304330833699959871506998612464798533130670377694999325158896591777210986795684851107725454068188288822567764912694521874079483339544658453938914789983271676836345351766016) 970
New Best Sum: 99^95 = (3848960788934848611927795802824596789608451156087366034658627953530148126008534258032267383768627487094610968554286692697374726725853195657679460590239636893953692985541958490801973870359499) 972
Euler 56: Powerful Digit Sum (3848960788934848611927795802824596789608451156087366034658627953530148126008534258032267383768627487094610968554286692697374726725853195657679460590239636893953692985541958490801973870359499) 972
*/
func main() {
	//test
	// Euler055()
	num, sum := Euler056(99, 99, 10)

	//run
	fmt.Printf("Euler 56: Powerful Digit Sum (%s) %d\n", num, sum)
}
