// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/* https://projecteuler.net/minimal=22
<p>Using <a href="resources/documents/0022_names.txt">names.txt</a> (right click and 'Save Link/Target As...'), a 46K text file containing over five-thousand first names, begin by sorting it into alphabetical order. Then working out the alphabetical value for each name, multiply this value by its alphabetical position in the list to obtain a name score.</p>
<p>For example, when the list is sorted into alphabetical order, COLIN, which is worth $3 + 15 + 12 + 9 + 14 = 53$, is the $938$th name in the list. So, COLIN would obtain a score of $938 \times 53 = 49714$.</p>
<p>What is the total of all the name scores in the file?</p>





*/

import (
	"bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	"sort"
	"strings"
	// "strconv"
	"os" // os.Stdout
)

func Euler022(fn string) int64 {
	var names []string
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler022 unable to open: " + fn)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	scanner.Split(euler.ScannerSplitNLDQ)
	for scanner.Scan() {
		str := strings.ToUpper(scanner.Text())
		if str != "," {
			names = append(names, str)
		}
	}
	fmt.Println("Euler022 read ", len(names), " names of 5163 expected.\t", 5163 == len(names))
	sort.Strings(names)
	limit := len(names)
	ret := int64(0)
	// os.Remove("test.txt")
	// fhw, err := os.Create("test.txt")
	// if nil != err {
	// panic("Euler022 unable to open: test.txt")
	// }
	// defer fhw.Close()
	// out := bufio.NewWriter(fhw)
	for ii := 0; ii < limit; ii++ {
		ret += int64(ii+1) * euler.AlphaSum(names[ii])
		// str := names[ii]
		// out.WriteString(names[ii])
		// out.Write([]byte{'\n'})
	}
	// out.Flush()
	if len(names) < 939 {
		fmt.Println("Warning check value COLIN != names[937] ~ missing: len ", len(names), names)
	} else if names[937] != "COLIN" {
		fmt.Println("Warning check value COLIN != names[937] ~ ", names[937])
		fmt.Println(0, names[0])
		fmt.Println(1, names[1])
		fmt.Println(936, names[936])
		fmt.Println(937, names[937])
		fmt.Println(938, names[938])
	}
	return ret
}

/*
Euler022 read  5163  names of 5163 expected.     true
871198282  score.
*/
func main() {
	// fmt.Println(grid)
	//test
	fmt.Println(int(euler.AlphaSum("Abc")) == 6)
	fmt.Println(int(euler.AlphaSum("Colin"))*938 == 49714)

	//run
	fmt.Println(Euler022("0022_names.txt"), " score.")
}
