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
)

func Euler0070(toMax uint64) uint64 {
	var bestI, ii, phi, bestFI, fi uint64
	var bestF, flNphi float64
	ptarget := toMax // Most of the runtime is trying to factor difficult numbers, the edge cases dominate...
	// ptarget := uint64(500_000) // Tried adding better root / square root functions, Factorizing is still too expensive...  Unsure if I need a better Factorizing method, better Phi method, or both.
	fmt.Printf("Finding primes to %d\n", ptarget)
	euler.Primes.Grow(ptarget)
	bestF, bestFI = 3.00, 3 // Not really but I know this value is higher than a real best which is 2 so...
Euler0070ii:
	for ii = 2; ii <= toMax; ii++ {
		phi = euler.EulerTotientPhi(ii)
		fi = ii / phi
		if fi <= bestFI {
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
				fmt.Printf("Found new best N/phi: %d/%d ~= %1.40f > %1.40f\n", ii, phi, flNphi, bestF)
				bestF, bestFI, bestI = flNphi, fi, ii
			}
		}
	}
	fmt.Printf("Found best N/phi: %d ~= %1.40f\n", bestI, bestF)
	return bestI
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 70 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

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
.
*/
func main() {
	//test
	// tested in the golang tests for "euler"

	//run
	fmt.Printf("Euler 70: Totient Permutation: %d\n", Euler0070(10_000_000))
}
