// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=74
https://projecteuler.net/minimal=74

<p>The number $145$ is well known for the property that the sum of the factorial of its digits is equal to $145$:
$$1! + 4! + 5! = 1 + 24 + 120 = 145.$$</p>
<p>Perhaps less well known is $169$, in that it produces the longest chain of numbers that link back to $169$; it turns out that there are only three such loops that exist:</p>
\begin{align}
&amp;169 \to 363601 \to 1454 \to 169\\
&amp;871 \to 45361 \to 871\\
&amp;872 \to 45362 \to 872
\end{align}
<p>It is not difficult to prove that EVERY starting number will eventually get stuck in a loop. For example,</p>
\begin{align}
&amp;69 \to 363600 \to 1454 \to 169 \to 363601 (\to 1454)\\
&amp;78 \to 45360 \to 871 \to 45361 (\to 871)\\
&amp;540 \to 145 (\to 145)
\end{align}
<p>Starting with $69$ produces a chain of five non-repeating terms, but the longest non-repeating chain with a starting number below one million is sixty terms.</p>
<p>How many chains, with a starting number below one million, contain exactly sixty non-repeating terms?</p>


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
	// "os" // os.Stdout
	// "strconv"
	// "strings"
)

func Euler0074(min, max, base uint32) (uint32, []uint32) {
	var chains []uint32
	var facts, lut map[uint32]uint32
	var lenChains, link, cur, next uint32
	var exists bool

	lut = make(map[uint32]uint32, 1_000_210)
	for ii := min; ii <= max; ii++ {
		facts = make(map[uint32]uint32, 60)
		cur = ii
		for link = 1; ; link++ {
			// This cuts the time about in half on my test system, it's still pretty slow thanks to the choice of a hash OR a handful of divisions
			if next, exists = lut[cur]; !exists {
				next = euler.DigitFactorialSum(cur, base)
				lut[cur] = next // Re-Tested with a guard for 1 < link here, the Look Up Table can be as small as 4017 in that case, but it didn't speed things up that much and the runtime became more volatile.  It might be worth for a hash table with a controlled hash method.
			}
			facts[cur] = next
			// fmt.Printf("%d: %d -> %d\n", ii, cur, next)
			if _, exists = facts[next]; exists {
				// loop detected
				if lenChains < link {
					lenChains = link
					chains = append(make([]uint32, 0, 4), ii)
				} else if lenChains == link {
					chains = append(chains, ii)
				}
				break // 1, link loop
			}
			cur = next
		}
	}
	// fmt.Printf("Needed a factorial lookup table of %d\n", len(lut))
	return lenChains, chains
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 74 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Needed a factorial lookup table of 242
Needed a factorial lookup table of 1000209
Euler 74: Digit Factorial Chains:       len: 60 Count: 402
[1479 1497 1749 1794 1947 1974 4079 4097 4179 4197 4709 4719 4790 4791 4907 4917 4970 4971 7049 7094 7149 7194 7409 7419 7490 7491 7904 7914 7940 7941 9047 9074 9147 9174 9407 9417 9470 9471 9704 9714 9740 9741 223479 223497 223749 223794 223947 223974 224379 224397 224739 224793 224937 224973 227349 227394 227439 227493 227934 227943 229347 229374 229437 229473 229734 229743 232479 232497 232749 232794 232947 232974 234279 234297 234729 234792 234927 234972 237249 237294 237429 237492 237924 237942 239247 239274 239427 239472 239724 239742 242379 242397 242739 242793 242937 242973 243279 243297 243729 243792 243927 243972 247239 247293 247329 247392 247923 247932 249237 249273 249327 249372 249723 249732 272349 272394 272439 272493 272934 272943 273249 273294 273429 273492 273924 273942 274239 274293 274329 274392 274923 274932 279234 279243 279324 279342 279423 279432 292347 292374 292437 292473 292734 292743 293247 293274 293427 293472 293724 293742 294237 294273 294327 294372 294723 294732 297234 297243 297324 297342 297423 297432 322479 322497 322749 322794 322947 322974 324279 324297 324729 324792 324927 324972 327249 327294 327429 327492 327924 327942 329247 329274 329427 329472 329724 329742 342279 342297 342729 342792 342927 342972 347229 347292 347922 349227 349272 349722 372249 372294 372429 372492 372924 372942 374229 374292 374922 379224 379242 379422 392247 392274 392427 392472 392724 392742 394227 394272 394722 397224 397242 397422 422379 422397 422739 422793 422937 422973 423279 423297 423729 423792 423927 423972 427239 427293 427329 427392 427923 427932 429237 429273 429327 429372 429723 429732 432279 432297 432729 432792 432927 432972 437229 437292 437922 439227 439272 439722 472239 472293 472329 472392 472923 472932 473229 473292 473922 479223 479232 479322 492237 492273 492327 492372 492723 492732 493227 493272 493722 497223 497232 497322 722349 722394 722439 722493 722934 722943 723249 723294 723429 723492 723924 723942 724239 724293 724329 724392 724923 724932 729234 729243 729324 729342 729423 729432 732249 732294 732429 732492 732924 732942 734229 734292 734922 739224 739242 739422 742239 742293 742329 742392 742923 742932 743229 743292 743922 749223 749232 749322 792234 792243 792324 792342 792423 792432 793224 793242 793422 794223 794232 794322 922347 922374 922437 922473 922734 922743 923247 923274 923427 923472 923724 923742 924237 924273 924327 924372 924723 924732 927234 927243 927324 927342 927423 927432 932247 932274 932427 932472 932724 932742 934227 934272 934722 937224 937242 937422 942237 942273 942327 942372 942723 942732 943227 943272 943722 947223 947232 947322 972234 972243 972324 972342 972423 972432 973224 973242 973422 974223 974232 974322]

real    0m2.537s
user    0m2.599s
sys     0m0.117s
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"
	r, rl := Euler0074(1, 100, 10)
	if 54 != r {
		panic(fmt.Sprintf("Euler 74: Expected 54 got %d (%d) %v", r, len(rl), rl))
	}

	//run
	r, rl = Euler0074(1, 1_000_000, 10)
	fmt.Printf("Euler 74: Digit Factorial Chains:\tlen: %d\tCount: %d\n%v\n", r, len(rl), rl)
	if 402 != len(rl) {
		panic("Did not reach expected value.")
	}
}
