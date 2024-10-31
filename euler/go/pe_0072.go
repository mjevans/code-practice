// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=72
https://projecteuler.net/minimal=72

<p>Consider the fraction, $\dfrac n d$, where $n$ and $d$ are positive integers. If $n \lt d$ and $\operatorname{HCF}(n,d)=1$, it is called a reduced proper fraction.</p>
<p>If we list the set of reduced proper fractions for $d \le 8$ in ascending order of size, we get:
$$\frac 1 8, \frac 1 7, \frac 1 6, \frac 1 5, \frac 1 4, \frac 2 7, \frac 1 3, \frac 3 8, \frac 2 5, \frac 3 7, \frac 1 2, \frac 4 7, \frac 3 5, \frac 5 8, \frac 2 3, \frac 5 7, \frac 3 4, \frac 4 5, \frac 5 6, \frac 6 7, \frac 7 8$$</p>
<p>It can be seen that there are $21$ elements in this set.</p>
<p>How many elements would be contained in the set of reduced proper fractions for $d \le 1\,000\,000$?</p>


*/
/*

New	ii	factors
0	1	sp
1	2	 p	1/2
2	3	 p	1/3 2/3
2	4	2	1/4 3/4
4	5	 p	1/5 2/5 3/5 4/5
2	6	2 3	1/6 5/6
6	7	 p	1/7 2/7 3/7 4/7 5/7 6/7
4	8	2	1/8 3/8 5/8 7/8

6	9	3	1/9 2/9 4/9 5/9 7/9 8/9
4	10	2 5	1/ 3/ 7/ 9/
10	11	 p	1/ 2/ 3/ 4/ 5/ 6/ 7/ 8/ 9/ 10/

21 in SUM(1,8) ... It is Totient again isn't it?

*/

import (
	// "bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0072(limit uint64) uint64 {
	var ret, ii uint64
	// sqrt = euler.SqrtU64(limit)
	// fmt.Printf("Euler 72\tGrow primes list to %d\n", sqrt)
	// euler.Primes.PrimeGlobalList(sqrt)
	fmt.Printf("Euler 72\n")
	for ii = 2; ii <= limit; ii++ {
		ret += euler.EulerTotientPhi(ii, 0)
		if 0 == ii&0xFFFF {
			fmt.Printf("Euler 72 ...\t%12d\t%15d\n", ii, ret)
		}
	}
	return ret
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 72 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

user: 5.5  5.7  5.75  6.6  6.8

With a working (fast) Lenstra factorization option, the runtime is cut by almost /4
Euler 72
Euler 72
Euler 72 ...           65536         1305514925
Euler 72 ...          131072         5222105723
Euler 72 ...          196608        11749660201
Euler 72 ...          262144        20888264643
Euler 72 ...          327680        32637919475
Euler 72 ...          393216        46998531267
Euler 72 ...          458752        63970257971
Euler 72 ...          524288        83553060809
Euler 72 ...          589824       105746731133
Euler 72 ...          655360       130551492759
Euler 72 ...          720896       157967275095
Euler 72 ...          786432       187994104135
Euler 72 ...          851968       220631912747
Euler 72 ...          917504       255880910551
Euler 72 ...          983040       293740710159
Euler 72: Counting Fractions:  303963552391

real    0m1.934s
user    0m1.980s
sys     0m0.056s
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	r := Euler0072(8)
	if 21 != r {
		panic(fmt.Sprintf("Euler 72: Expected 21 got %d", r))
	}

	//run
	r = Euler0072(1_000_000)
	fmt.Printf("Euler 72: Counting Fractions:  %d\n", r)
	if 303963552391 != r {
		panic("Did not reach expected value.")
	}
}
