// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=70
https://projecteuler.net/minimal=70

<p>Euler's totient function, $\phi(n)$ [sometimes called the phi function], is used to determine the number of positive numbers less than or equal to $n$ which are relatively prime to $n$. For example, as $1, 2, 4, 5, 7$, and $8$, are all less than nine and relatively prime to nine, $\phi(9)=6$.<br>The number $1$ is considered to be relatively prime to every positive number, so $\phi(1)=1$. </p>
<p>Interestingly, $\phi(87109)=79180$, and it can be seen that $87109$ is a permutation of $79180$.</p>
<p>Find the value of $n$, $1 \lt n \lt 10^7$, for which $\phi(n)$ is a permutation of $n$ and the ratio $n/\phi(n)$ produces a minimum.</p>


*/
/*

The permutation checking signature involves a at least a divide per digit, as well as character shuffles, so that'll get tested last...

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
	// "runtime/pprof"
)

func Euler0070(toMax uint64) uint64 {
	var bestI, ii, phi, phiMin, bestFI, fi uint64
	var bestF, flNphi float64
	_, _ = bestFI, fi
	// ptarget := euler.SqrtU64(toMax) // Most of the runtime is trying to factor difficult numbers, the edge cases dominate...
	// ptarget := uint64(500_000) // Tried adding better root / square root functions, Factorizing is still too expensive...  Unsure if I need a better Factorizing method, better Phi method, or both.
	// fmt.Printf("Finding primes to %d\n", ptarget)
	// euler.Primes.Grow(ptarget)
	// euler.Primes.PrimeGlobalList(ptarget)
	bestF, bestFI = 3.00, 3 // Not really but I know this value is higher than a real best which is 2 so...
Euler0070ii:
	for ii = toMax; 2 <= ii; ii-- {
		phiMin = uint64(float64(ii) / bestF)
		phi = euler.EulerTotientPhi(ii, phiMin)
		// fi = ii / phi
		// if fi <= bestFI {
		if phi > phiMin {
			flNphi = float64(ii) / float64(phi)
			if flNphi < bestF {
				slii := euler.Uint8CopyInsertSort(euler.Uint64ToDigitsUint8(ii, 10))
				slphi := euler.Uint8CopyInsertSort(euler.Uint64ToDigitsUint8(phi, 10))
				limII, limPHI := len(slii), len(slphi)
				if limII != limPHI {
					continue
				}
				for jj := 0; jj < limII; jj++ {
					if slii[jj] != slphi[jj] {
						continue Euler0070ii
					}
				}
				fmt.Printf("Found new best N/phi: %d/%d ~= %1.40f < %1.40f\n", ii, phi, flNphi, bestF)
				bestF, bestFI, bestI = flNphi, fi, ii
			}
		}
	}
	fmt.Printf("Found best N/phi: %d ~= %1.40f\n", bestI, bestF)
	return bestI
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 70 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Found new best N/phi: 5050429/5045920 ~= 1.0008935932396867407589979848125949501991 > 1.0010763255282499883946911722887307405472
Found new best N/phi: 5380657/5375860 ~= 1.0008923223447039330125107881030999124050 > 1.0008935932396867407589979848125949501991
Found new best N/phi: 5886817/5881876 ~= 1.0008400381102899867613587048253975808620 > 1.0008923223447039330125107881030999124050
Found new best N/phi: 6018163/6013168 ~= 1.0008306769410069136938545852899551391602 > 1.0008400381102899867613587048253975808620
Found new best N/phi: 6636841/6631684 ~= 1.0007776305384876724957621263456530869007 > 1.0008306769410069136938545852899551391602
Found new best N/phi: 7026037/7020736 ~= 1.0007550490432912670968335078214295208454 > 1.0007776305384876724957621263456530869007
Found new best N/phi: 7357291/7351792 ~= 1.0007479809004389270654655774706043303013 > 1.0007550490432912670968335078214295208454
Found new best N/phi: 7507321/7501732 ~= 1.0007450279482124066987580590648576617241 > 1.0007479809004389270654655774706043303013
Found new best N/phi: 8316907/8310976 ~= 1.0007136345959848355846588674467056989670 > 1.0007450279482124066987580590648576617241
Found new best N/phi: 8319823/8313928 ~= 1.0007090511248113440245788297033868730068 > 1.0007136345959848355846588674467056989670
Found best N/phi: 8319823 ~= 1.0007090511248113440245788297033868730068
Euler 70: Totient Permutation: 8319823

After a change in algorithm to allow for early aborts it now runs in less than 1 min of runtime.

for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 70 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go build -o $(printf "pe_%04d" "$ii") $(printf "pe_%04d.go" "$ii") ; time ./$(printf "pe_%04d" "$ii") || break ; done
Finding primes to 3162
Found new best N/phi: 9983167/9973816 ~= 1.0009375548937338162858168288948945701122 < 3.0000000000000000000000000000000000000000
Found new best N/phi: 9848203/9840328 ~= 1.0008002782021088172825784567976370453835 < 1.0009375548937338162858168288948945701122
Found new best N/phi: 8357821/8351872 ~= 1.0007122953991631764125713743851520121098 < 1.0008002782021088172825784567976370453835
Found new best N/phi: 8319823/8313928 ~= 1.0007090511248113440245788297033868730068 < 1.0007122953991631764125713743851520121098
Found best N/phi: 8319823 ~= 1.0007090511248113440245788297033868730068
Euler 70: Totient Permutation: 8319823

go tool pprof pe_0070 latest.goprof
File: pe_0070
Type: cpu
Time: Oct 25, 2024 at 10:56pm (UTC)
Duration: 43.86s, Total samples = 45.77s (104.34%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 33.05s, 72.21% of 45.77s total
Dropped 213 nodes (cum <= 0.23s)
Showing top 10 nodes out of 64

	  flat  flat%   sum%        cum   cum%
	14.66s 32.03% 32.03%     14.71s 32.14%  math/rand.seedrand (inline) << math.big BigInt.ProbablyPrime()  I do need a non math.big probably prime test
	 4.89s 10.68% 42.71%     19.60s 42.82%  math/rand.(*rngSource).Seed
	 2.77s  6.05% 48.77%      4.60s 10.05%  math/big.divWVW
	 2.17s  4.74% 53.51%     12.72s 27.79%  math/big.nat.expNN
	 1.86s  4.06% 57.57%     43.44s 94.91%  euler.EulerTotientPhi
	 1.56s  3.41% 60.98%      8.82s 19.27%  math/big.nat.div
	 1.43s  3.12% 64.10%      1.43s  3.12%  math/big.nat.norm (inline)
	 1.38s  3.02% 67.12%      1.38s  3.02%  euler.GCDbin[go.shape.uint64]
	 1.20s  2.62% 69.74%      2.74s  5.99%  runtime.mallocgc
	 1.13s  2.47% 72.21%      1.14s  2.49%  math/big.reciprocalWord (inline)

math/big ProbablyPrime(0) is sufficient for 64 bit numbers, it runs the MR(2) + a Lucas Strong = B-PSW test internally, that's where the math/rand sink crept in.

While writing a 64 bit (plus some 128 bit internal ops) version of the tests Lenstra was also improved and as a whole with both a Probably Prime filter AND Pollard + Lenstra factorization, it's possible to find Totient values without exhaustively proving all primes have been tested other than a final prime.

Untuned:

~ 4.98 S user time average across three runs.

.
*/
func main() {
	//test
	// tested in the golang tests for "euler"

	// f, err := os.Create("latest.goprof")
	// if nil != err {
	//	panic("Unable to create latest.goprof file\n")
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	//run
	fmt.Printf("Euler 70: Totient Permutation: %d\n", Euler0070(10_000_000))
}
