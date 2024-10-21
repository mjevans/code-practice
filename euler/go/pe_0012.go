// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=12
https://projecteuler.net/minimal=12

<p>The sequence of triangle numbers is generated by adding the natural numbers. So the $7$<sup>th</sup> triangle number would be $1 + 2 + 3 + 4 + 5 + 6 + 7 = 28$. The first ten terms would be:
$$1, 3, 6, 10, 15, 21, 28, 36, 45, 55, \dots$$</p>
<p>Let us list the factors of the first seven triangle numbers:</p>
\begin{align}
\mathbf 1 &amp;\colon 1\\
\mathbf 3 &amp;\colon 1,3\\
\mathbf 6 &amp;\colon 1,2,3,6\\
\mathbf{10} &amp;\colon 1,2,5,10\\
\mathbf{15} &amp;\colon 1,3,5,15\\
\mathbf{21} &amp;\colon 1,3,7,21\\
\mathbf{28} &amp;\colon 1,2,4,7,14,28
\end{align}
<p>We can see that $28$ is the first triangle number to have over five divisors.</p>
<p>What is the value of the first triangle number to have over five hundred divisors?</p>

*/

import (
	"euler"
	"fmt"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func Euler012(minDivisors uint32) uint32 {
	var triangle, ii, target uint32
	for t := minDivisors; 0 < t; t >>= 1 {
		target++
	}
	euler.Primes.Grow(1_000_000)
	for {
		triangle += ii
		f := euler.Primes.Factorize(uint64(triangle))
		if target <= f.Lenpow {
			pd := *(f.ProperDivisors())
			divisors := uint32(len(pd)) + 1
			if divisors >= minDivisors {
				fmt.Println(triangle, ii, divisors, pd)
				break
			}
		}
		ii++
	}
	return triangle
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 12 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

The new interface makes pre-testing a count of the number of factors (Lenpow) a great guard against the more costly verification step.

28 7 6 [1 2 4 7 14]
Euler012:       First number with 5 factors:     28 true

76576500 12375 576 [1 2 3 4 5 6 7 9 10 11 12 13 14 15 17 18 20 21 22 25 26 28 30 33 34 35 36 39 42 44 45 50 51 52 55 60 63 65 66 68 70 75 77 78 84 85 90 91 99 100 102 105 110 117 119 125 126 130 132 140 143 150 153 154 156 165 170 175 180 182 187 195 198 204 210 220 221 225 231 234 238 250 252 255 260 273 275 286 300 306 308 315 325 330 340 350 357 364 374 375 385 390 396 420 425 429 442 450 455 462 468 476 495 500 510 525 546 550 561 572 585 595 612 630 650 660 663 693 700 714 715 748 750 765 770 780 819 825 850 858 875 884 900 910 924 935 975 990 1001 1020 1050 1071 1092 1100 1105 1122 1125 1155 1170 1190 1260 1275 1287 1300 1309 1326 1365 1375 1386 1428 1430 1500 1530 1540 1547 1575 1625 1638 1650 1683 1700 1716 1750 1785 1820 1870 1925 1950 1980 1989 2002 2100 2125 2142 2145 2210 2244 2250 2275 2310 2340 2380 2431 2475 2550 2574 2618 2625 2652 2730 2750 2772 2805 2860 2925 2975 3003 3060 3094 3150 3250 3276 3300 3315 3366 3465 3500 3570 3575 3740 3825 3850 3900 3927 3978 4004 4095 4125 4250 4284 4290 4420 4500 4550 4620 4641 4675 4862 4875 4950 5005 5100 5148 5236 5250 5355 5460 5500 5525 5610 5775 5850 5950 6006 6188 6300 6375 6435 6500 6545 6630 6732 6825 6930 7140 7150 7293 7650 7700 7735 7854 7875 7956 8190 8250 8415 8500 8580 8925 9009 9100 9282 9350 9625 9724 9750 9900 9945 10010 10500 10710 10725 11050 11220 11375 11550 11700 11781 11900 12012 12155 12375 12750 12870 13090 13260 13650 13860 13923 14025 14300 14586 14625 14875 15015 15300 15470 15708 15750 16380 16500 16575 16830 17017 17325 17850 17875 18018 18564 18700 19125 19250 19500 19635 19890 20020 20475 21420 21450 21879 22100 22750 23100 23205 23375 23562 24310 24750 25025 25500 25740 26180 26775 27300 27625 27846 28050 28875 29172 29250 29750 30030 30940 31500 32175 32725 33150 33660 34034 34125 34650 35700 35750 36036 36465 38250 38500 38675 39270 39780 40950 42075 42900 43758 44625 45045 45500 46410 46750 47124 48620 49500 49725 50050 51051 53550 53625 55250 55692 56100 57750 58500 58905 59500 60060 60775 64350 65450 66300 68068 68250 69300 69615 70125 71500 72930 75075 76500 77350 78540 81900 82875 84150 85085 86625 87516 89250 90090 92820 93500 98175 99450 100100 102102 102375 107100 107250 109395 110500 115500 116025 117810 121550 125125 128700 130900 133875 136500 139230 140250 145860 150150 153153 154700 160875 163625 165750 168300 170170 173250 178500 180180 182325 193375 196350 198900 204204 204750 210375 214500 218790 225225 232050 235620 243100 248625 250250 255255 267750 278460 280500 294525 300300 303875 306306 321750 327250 331500 340340 346500 348075 364650 375375 386750 392700 409500 420750 425425 437580 450450 464100 490875 497250 500500 510510 535500 546975 580125 589050 607750 612612 643500 654500 696150 729300 750750 765765 773500 841500 850850 900900 911625 981750 994500 1021020 1093950 1126125 1160250 1178100 1215500 1276275 1392300 1472625 1501500 1531530 1701700 1740375 1823250 1963500 2127125 2187900 2252250 2320500 2552550 2734875 2945250 3063060 3480750 3646500 3828825 4254250 4504500 5105100 5469750 5890500 6381375 6961500 7657650 8508500 10939500 12762750 15315300 19144125 25525500 38288250]

Euler012:       First number with 501+ factors:  76576500


*/

func main() {
	// fmt.Println(grid)
	//test
	ans := Euler012(5)
	fmt.Println("Euler012:\tFirst number with 5 factors:\t", ans, ans == 28)
	//run
	ans = Euler012(501)
	fmt.Println("\nEuler012:\tFirst number with 501+ factors:\t", ans)
}
