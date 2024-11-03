// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=73
https://projecteuler.net/minimal=73

<p>Consider the fraction, $\dfrac n d$, where $n$ and $d$ are positive integers. If $n \lt d$ and $\operatorname{HCF}(n, d)=1$, it is called a reduced proper fraction.</p>
<p>If we list the set of reduced proper fractions for $d \le 8$ in ascending order of size, we get:
$$\frac 1 8, \frac 1 7, \frac 1 6, \frac 1 5, \frac 1 4, \frac 2 7, \frac 1 3, \mathbf{\frac 3 8, \frac 2 5, \frac 3 7}, \frac 1 2, \frac 4 7, \frac 3 5, \frac 5 8, \frac 2 3, \frac 5 7, \frac 3 4, \frac 4 5, \frac 5 6, \frac 6 7, \frac 7 8$$</p>
<p>It can be seen that there are $3$ fractions between $\dfrac 1 3$ and $\dfrac 1 2$.</p>
<p>How many fractions lie between $\dfrac 1 3$ and $\dfrac 1 2$ in the sorted set of reduced proper fractions for $d \le 12\,000$?</p>

*/
/*

Repeatedly doing the thing from Problem 71 seems too slow, Phi / Euler's Totient seems too slow too.

Searching for : totient next reduced fraction

Within the first page or two
https://ibmathsresources.com/2018/05/25/farey-sequences/

(might be easier to read https://en.wikipedia.org/wiki/Farey_sequence )

There's also a paper
https://www.researchgate.net/publication/353555140_Revisited_Carmichael%27s_Reduced_Totient_Function
"The modified Totient function of Carmichael λ(.) is revisited, where important properties have been highlighted. Particularly, an iterative scheme is given for calculating the λ(.) function. A comparison between the Euler ϕ and the reduced totient λ(.) functions aiming to quantify the reduction between is given."


https://en.wikipedia.org/wiki/Farey_sequence#Sequence_length_and_index_of_a_fraction
https://en.wikipedia.org/wiki/Carmichael_function

Search revised: farey  reduced fraction
Might be useful, but not very easy for a visual / kinesthetic learner https://uu.diva-portal.org/smash/get/diva2:1116979/FULLTEXT01.pdf

https://blogs.sas.com/content/iml/2021/03/17/farey-sequence.html

Mediant Property ?
https://en.wikipedia.org/wiki/Farey_sequence#Next_term

https://mathoverflow.net/a/458425
"∑kd=1(⌊n/d⌋−⌈k/d⌉+1)⌊k/d⌋μ(d) lets one compute it efficiently."

An interesting detour https://www.cut-the-knot.org/blue/Fusc.shtml


Just above Farey Neighbors
https://en.wikipedia.org/wiki/Farey_sequence#Sequence_length_and_index_of_a_fraction
"The index of 1 / k {\displaystyle 1/k} where n / ( i + 1 ) < k ≤ n / i {\displaystyle n/(i+1)<k\leq n/i} and n {\displaystyle n} is the least common multiple of the first i {\displaystyle i} numbers, n = l c m ( [ 2 , i ] ) {\displaystyle n={\rm {lcm}}([2,i])}, is given by:[7]

    I n ( 1 / k ) = 1 + n ∑ j = 1 i φ ( j ) j − k Φ ( i ) . {\displaystyle I_{n}(1/k)=1+n\sum _{j=1}^{i}{\frac {\varphi (j)}{j}}-k\Phi (i).}"

Too post-doc math heavy for my present knowledge https://cs.uwaterloo.ca/journals/JIS/VOL25/Tomas/tomas5.pdf

Various useful properties... (wherein Fn is the Farey order length)

Idx (1/2) = (|Fn|-1)/2
Idx (h/k) = |Fn| - 1 - Idx((k-h)/k)

Need the length of a Farey order though...

farey farey order length calculation formula

https://www.nature.com/articles/s41598-021-99545-w

'memoization still too slow for large values of n'




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

func Euler0073(d uint64) uint64 {
	fn := euler.FareyLengthAlgE(d)
	return euler.FareyIndex(fn, d, 1, 2) - euler.FareyIndex(fn, d, 1, 3) - 1 // Remove the midpoint index
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 73 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

.
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	r := Euler0073(8)
	if 3 != r {
		panic(fmt.Sprintf("Euler 73: Expected 21 got %d", r))
	}

	//run
	r = Euler0073(12_000)
	fmt.Printf("Euler 73: Counting Fractions in a Range:\t%d\n", r)
	if 0 != r {
		//panic("Did not reach expected value.")
	}
}
