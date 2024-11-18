// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=87
https://projecteuler.net/minimal=87

<p>The smallest number expressible as the sum of a prime square, prime cube, and prime fourth power is $28$. In fact, there are exactly four numbers below fifty that can be expressed in such a way:</p>
\begin{align}
28 &amp;= 2^2 + 2^3 + 2^4\\
33 &amp;= 3^2 + 2^3 + 2^4\\
49 &amp;= 5^2 + 2^3 + 2^4\\
47 &amp;= 2^2 + 3^3 + 2^4
\end{align}
<p>How many numbers below fifty million can be expressed as the sum of a prime square, prime cube, and prime fourth power?</p>


/
*/
/*
	At a glance this seems fairly obvious, and the listed example even biases towards that solution.
	Three wheels, prime^4 + prime^3 + prime^2 and cut things off when the entry to any wheel's loop is greater than the target.
	sqrt(50M) = 7071+ so technically it'd be possible to cache all the multiplications; that isn't the expensive part on a modern CPU, it might have been in the late 1990s / early 2000s
/
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

func Euler0087(target uint32) uint32 {
	t := euler.SqrtU64(uint64(target))
	fmt.Printf("Euler 87: Primes.Grow(%d + 8)\n", t)
	euler.Primes.Grow(t + 8)
	var found, q, c, s, vq, vc, vs uint32
	check := make(map[uint32]uint8, 1136692)

	q = 1
	for {
		q = uint32(euler.Primes.PrimeAfter(uint64(q)))
		vq = q * q * q * q
		if target < vq {
			break
		}
		c = 1
		for {
			c = uint32(euler.Primes.PrimeAfter(uint64(c)))
			vc = vq + c*c*c
			if target < vc {
				break
			}
			s = 1
			for {
				s = uint32(euler.Primes.PrimeAfter(uint64(s)))
				vs = vc + s*s
				if target < vs {
					break
				}
				if _, exists := check[vs]; !exists {
					found++
				}
				check[vs]++
			}
		}
		fmt.Printf("Euler 87: exhausted %4d^4 with final set (%4d, %4d, %4d) = %12d ; found (so far) %12d\n", q, q, c, s, vs, found)
	}
	return found
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 87 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Euler 87: Primes.Grow(7 + 8)
Euler 87: exhausted    2^4 with final set (   2,    5,    3) =           52 ; found (so far)            4
Euler 87: Primes.Grow(7071 + 8)
Euler 87: exhausted    2^4 with final set (   2,  373,  757) =     50003928 ; found (so far)        58013
Euler 87: exhausted    3^4 with final set (   3,  373,  757) =     50003993 ; found (so far)       116024
Euler 87: exhausted    5^4 with final set (   5,  373,  757) =     50004537 ; found (so far)       173867
Euler 87: exhausted    7^4 with final set (   7,  373,  757) =     50006313 ; found (so far)       231554
Euler 87: exhausted   11^4 with final set (  11,  373,  751) =     50009505 ; found (so far)       288622
Euler 87: exhausted   13^4 with final set (  13,  373,  739) =     50005545 ; found (so far)       345957
Euler 87: exhausted   17^4 with final set (  17,  373,  701) =     50005785 ; found (so far)       402489
Euler 87: exhausted   19^4 with final set (  19,  373,  673) =     50014113 ; found (so far)       458708
Euler 87: exhausted   23^4 with final set (  23,  373,  541) =     50003385 ; found (so far)       514435
Euler 87: exhausted   29^4 with final set (  29,  367, 1741) =     50006641 ; found (so far)       570169
Euler 87: exhausted   31^4 with final set (  31,  367, 1693) =     50058049 ; found (so far)       625074
Euler 87: exhausted   37^4 with final set (  37,  367, 1367) =     50011129 ; found (so far)       678538
Euler 87: exhausted   41^4 with final set (  41,  367,  953) =     50002249 ; found (so far)       731649
Euler 87: exhausted   43^4 with final set (  43,  367,  563) =     50004049 ; found (so far)       784156
Euler 87: exhausted   47^4 with final set (  47,  359, 1069) =     50009419 ; found (so far)       835346
Euler 87: exhausted   53^4 with final set (  53,  349,  577) =     50005333 ; found (so far)       883429
Euler 87: exhausted   59^4 with final set (  59,  337, 1277) =     50012781 ; found (so far)       927719
Euler 87: exhausted   61^4 with final set (  61,  331, 2081) =     50031415 ; found (so far)       971155
Euler 87: exhausted   67^4 with final set (  67,  311,  967) =     50020653 ; found (so far)      1008970
Euler 87: exhausted   71^4 with final set (  71,  293, 1399) =     50034069 ; found (so far)      1042157
Euler 87: exhausted   73^4 with final set (  73,  281,  593) =     50003823 ; found (so far)      1072259
Euler 87: exhausted   79^4 with final set (  79,  223, 1289) =     50005533 ; found (so far)      1090501
Euler 87: exhausted   83^4 with final set (  83,  137,  547) =     50005621 ; found (so far)      1097343
Euler 87: Prime Power Triples: 1097343

real    0m0.262s
user    0m0.268s
sys     0m0.097s
.
*/
func main() {
	var r uint32
	//test
	r = Euler0087(50 - 1)
	if 4 != r {
		panic(fmt.Sprintf("Did not reach expected test value. Got: %d", r))
	}

	//run
	r = Euler0087(50_000_000 - 1)
	fmt.Printf("Euler 87: Prime Power Triples: %d\n", r)
	if 1097343 != r {
		panic("Did not reach expected value.")
	}
}
