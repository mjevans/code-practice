// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=100
https://projecteuler.net/minimal=100

<p>If a box contains twenty-one coloured discs, composed of fifteen blue discs and six red discs, and two discs were taken at random, it can be seen that the probability of taking two blue discs, $P(\text{BB}) = (15/21) \times (14/20) = 1/2$.</p>
<p>The next such arrangement, for which there is exactly $50\%$ chance of taking two blue discs at random, is a box containing eighty-five blue discs and thirty-five red discs.</p>
<p>By finding the first arrangement to contain over $10^{12} = 1\,000\,000\,000\,000$ discs in total, determine the number of blue discs that the box would contain.</p>

/
*/
/*
	Preface: Euler  100 will be the last problem that I plan to post completed sections of to the web, per https://projecteuler.net/about
		It's possible I'll continue to work on later problems.
		I might try to factor out not-problem specific utility methods to a library like "euler/pe_euler.go" is within this.
		However my goal will shift from a focus on continued progression in this set of problems.


	Generally, the way this works is... Some probability (over a given balance point)  * some probability (under a given balance point...) (removing one/N of the tokens) yields 50% exactly.
	That surely happens along a line or curve of some sort, but it's only expressed as a whole integer for every number occasionally.

	P		B0		B1		0		1			Delta
	4		3		2		0.75		2/3			0.08(3)
	21		15		14		0.(714285)	0.7			0.0(142857)
	120		85		84		0.708(3)	0.70588...		0.0024509...
	697		493
	4060		2871
	23661		16731
	137904		97513

	0.50 =?= (a)/(a+b) * (a-1)/(a-1+b) == 1/2
	That looks like it's going to get messy fast...
	(a)/(a+b) * (a-1)/(a-1+b) == 1/2
	((a)*(a-1))/((a+b)*(a-1+b)) == 1/2
	((a)*(a-1))<<1 == ((a+b)*(a-1+b))

	While watching this run after a validation run I noticed something...
ii:        1000000000000                707106781182            707106781188      0xfffffca3ebf9dfe8
ii:        1000000000000                707106781186            707106781188           0x45b8e288000
ii:        1000000000000                707106781186            707106781187      0xffffff36775ebff4
ii:        1000000000000                707106781187            707106781187           0x1c902c39ffc
.7071... that seems like a number I recall offhand.
https://en.wikipedia.org/wiki/Square_root_of_2#Multiplicative_inverse
    0.70710678118654752440084436210484903928483593768847...
That sort of does make sense.  Trying to get a '2' as output, even if it's inversed... 1/2 rather than 2.
https://pkg.go.dev/math@go1.22.6#pkg-constants
Sqrt2   = 1.41421356237309504880168872420969807856967187537694807317667974

This is still taking 'forever' but I think anything faster would need to approach finding a pattern in the recurrence or might use some sort of continued fraction representation.
It wouldn't surprise me if a series of numbers like Euler 57 or 65 but a little different was that pattern or if there were a really fast way of doing this that's math (trivia) knowledge heavy.

Yeah... that's almost what I thought, but I didn't realize #94 (done about 2 weeks ago for me with thanksgiving at family between) was related at all.
Also I think the post-spoilers thread post #6 from Fri, 22 Jul 2005, 13:40 combined with this problem has slightly increased how much I understand that sort of number sequence.

My first runtime to completion on this was ~32 min, but I prefer to run them a second time at least with the answer checked in the code as well.

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

const Sqrt2 = 1.41421356237309504880168872420969807856967187537694807317667974

func Euler0100(min uint64) (uint64, uint64) {
	// ii > min  In Go unsigned ints overflow to mod-wrap back to 0 without error, so the test is valid.  However...
	var ii, q, left, right, pos uint64
	for ii = min; ; ii++ {
		q = ii * (ii - 1)
		if q < min {
			return Euler0100UUF(ii)
			// panic("Overflow before solution found.  Unable to calculate ii*(ii-1) within uint64")
		}
		left, right = ii>>1+ii>>3+ii>>4, ii>>1+ii>>2 // 0.75 is the maximum right will ever get, and that's in a range of 4, so whatever; the 0.6875 (1/2+1/8+1/16) is an OK approximation for the low end
		for left != right {
			pos = (left + right) >> 1 // Ceil NORMALLY required, but this rounds up rather than down, so FLOOR and L = pos +1
			if (pos * (pos - 1) << 1) < q {
				left = pos + 1
			} else {
				right = pos
			}
		}
		if ((right * (right - 1)) << 1) == q {
			fmt.Printf("Found: 0.5 == ( %d / %d ) * ( %d / %d ) \n", right, ii, right-1, ii-1)
			return right, ii
		}
	}
}

func Euler0100UU(min uint64) (uint64, uint64) {
	// ii > min  In Go unsigned ints overflow to mod-wrap back to 0 without error, so the test is valid.  However...
	var ii, qh, ql, th, tl, left, right, pos uint64
	fmt.Printf("Entered Euler0100UU(%d)\n", min)
	for ii = min; ; ii++ {
		if ii < min {
			panic("ii Overflowed uint64")
		}
		if 0 == ii&0xFFFFFF {
			fmt.Printf("\t%d", ii)
		}
		qh, ql = euler.UU64Mul(ii, ii-1)
		left, right = ii>>1+ii>>3+ii>>4+ii>>6, ii>>1+ii>>3+ii>>4+ii>>5 // 0.75 is the maximum right will ever get, and that's in a range of 4, so whatever; the 0.6875 (1/2+1/8+1/16) is an OK approximation for the low end
		for left != right {
			_, th, tl = euler.UU64AddUU64(0, left, 0, right) // Ceil NORMALLY required, but this rounds up rather than down, so FLOOR and L = pos +1
			pos = (th << 63) | (tl >> 1)
			th, tl = euler.UU64Mul(pos, pos-1)
			th, tl = (th<<1)|(tl>>63), tl<<1 // prepare (pos*(pos-1)) << 1 for comparison
			// 0 > ( (pos*(pos-1)) << 1 ) - Q
			if -1 == euler.UU64Cmp(th, tl, qh, ql) {
				left = pos + 1
			} else {
				right = pos
			}
			// th, tl = euler.UU64SubUU64(qh, ql, th, tl)
			// fmt.Printf("ii: %20d\t%20d\t%20d\t%#20x\n", ii, left, right, tl)
		}
		th, tl = euler.UU64Mul(right, right-1)
		th, tl = (th<<1)|(tl>>63), tl<<1 // prepare (R*(R-1)) << 1 for comparison
		if 0 == euler.UU64Cmp(th, tl, qh, ql) {
			fmt.Printf("Found: 0.5 == ( %d / %d ) * ( %d / %d ) \n", right, ii, right-1, ii-1)
			return right, ii
		}
		// return 0, 0
	}
}

func Euler0100UUF(min uint64) (uint64, uint64) {
	var ii, qh, ql, th, tl, dd uint64
	var cmp int
	fmt.Printf("Entered Euler0100UUF(%d)\n", min)
	for ii = min; ; ii++ {
		if ii < min {
			panic("ii Overflowed uint64")
		}
		if 0 == ii&0xFFFFFFF {
			fmt.Printf("\t%d", ii)
		}
		qh, ql = euler.UU64Mul(ii, ii-1)
		for dd = uint64(float64(ii) / Sqrt2); ; dd++ {
			th, tl = euler.UU64Mul(dd, dd-1)
			th, tl = (th<<1)|(tl>>63), tl<<1 // prepare (pos*(pos-1)) << 1 for comparison
			// 0 > ( (pos*(pos-1)) << 1 ) - Q
			cmp = euler.UU64Cmp(th, tl, qh, ql)
			if 0 == cmp {
				fmt.Println()
				return dd, ii
			} else if 1 == cmp {
				break
			}
		}
	}
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 100 ; do go fmt $(printf "pe_%04d.go" "$ii") ; time go run $(printf "pe_%04d.go" "$ii") || break ; done

Found: 0.5 == ( 3 / 4 ) * ( 2 / 3 )
Found: 0.5 == ( 15 / 21 ) * ( 14 / 20 )
Found: 0.5 == ( 85 / 120 ) * ( 84 / 119 )
Found: 0.5 == ( 493 / 697 ) * ( 492 / 696 )
Found: 0.5 == ( 2871 / 4060 ) * ( 2870 / 4059 )
Found: 0.5 == ( 16731 / 23661 ) * ( 16730 / 23660 )
Found: 0.5 == ( 97513 / 137904 ) * ( 97512 / 137903 )
Found: 0.5 == ( 568345 / 803761 ) * ( 568344 / 803760 )
Found: 0.5 == ( 3312555 / 4684660 ) * ( 3312554 / 4684659 )
Entered Euler0100UUF(1000008221457)
        1000190509056   1000458944512   1000727379968   1000995815424   1001264250880   1001532686336   1001801121792   1002069557248   1002337992704   1002606428160   1002874863616   1003143299072   1003411734528
1003680169984   1003948605440   1004217040896   1004485476352   1004753911808   1005022347264   1005290782720   1005559218176   1005827653632   1006096089088   1006364524544   1006632960000   1006901395456   1007169830912        1007438266368   1007706701824   1007975137280   1008243572736   1008512008192   1008780443648   1009048879104   1009317314560   1009585750016   1009854185472   1010122620928   1010391056384   1010659491840        1010927927296   1011196362752   1011464798208   1011733233664   1012001669120   1012270104576   1012538540032   1012806975488   1013075410944   1013343846400   1013612281856   1013880717312   1014149152768        1014417588224   1014686023680   1014954459136   1015222894592   1015491330048   1015759765504   1016028200960   1016296636416   1016565071872   1016833507328   1017101942784   1017370378240   1017638813696        1017907249152   1018175684608   1018444120064   1018712555520   1018980990976   1019249426432   1019517861888   1019786297344   1020054732800   1020323168256   1020591603712   1020860039168   1021128474624        1021396910080   1021665345536   1021933780992   1022202216448   1022470651904   1022739087360   1023007522816   1023275958272   1023544393728   1023812829184   1024081264640   1024349700096   1024618135552        1024886571008   1025155006464   1025423441920   1025691877376   1025960312832   1026228748288   1026497183744   1026765619200   1027034054656   1027302490112   1027570925568   1027839361024   1028107796480        1028376231936   1028644667392   1028913102848   1029181538304   1029449973760   1029718409216   1029986844672   1030255280128   1030523715584   1030792151040   1031060586496   1031329021952   1031597457408        1031865892864   1032134328320   1032402763776   1032671199232   1032939634688   1033208070144   1033476505600   1033744941056   1034013376512   1034281811968   1034550247424   1034818682880   1035087118336        1035355553792   1035623989248   1035892424704   1036160860160   1036429295616   1036697731072   1036966166528   1037234601984   1037503037440   1037771472896   1038039908352   1038308343808   1038576779264        1038845214720   1039113650176   1039382085632   1039650521088   1039918956544   1040187392000   1040455827456   1040724262912   1040992698368   1041261133824   1041529569280   1041798004736   1042066440192        1042334875648   1042603311104   1042871746560   1043140182016   1043408617472   1043677052928   1043945488384   1044213923840   1044482359296   1044750794752   1045019230208   1045287665664   1045556101120        1045824536576   1046092972032   1046361407488   1046629842944   1046898278400   1047166713856   1047435149312   1047703584768   1047972020224   1048240455680   1048508891136   1048777326592   1049045762048        1049314197504   1049582632960   1049851068416   1050119503872   1050387939328   1050656374784   1050924810240   1051193245696   1051461681152   1051730116608   1051998552064   1052266987520   1052535422976        1052803858432   1053072293888   1053340729344   1053609164800   1053877600256   1054146035712   1054414471168   1054682906624   1054951342080   1055219777536   1055488212992   1055756648448   1056025083904        1056293519360   1056561954816   1056830390272   1057098825728   1057367261184   1057635696640   1057904132096   1058172567552   1058441003008   1058709438464   1058977873920   1059246309376   1059514744832        1059783180288   1060051615744   1060320051200   1060588486656   1060856922112   1061125357568   1061393793024   1061662228480   1061930663936   1062199099392   1062467534848   1062735970304   1063004405760        1063272841216   1063541276672   1063809712128   1064078147584   1064346583040   1064615018496   1064883453952   1065151889408   1065420324864   1065688760320   1065957195776   1066225631232   1066494066688        1066762502144   1067030937600   1067299373056   1067567808512   1067836243968   1068104679424   1068373114880   1068641550336   1068909985792   1069178421248   1069446856704   1069715292160   1069983727616        1070252163072
Euler 100: Arranged Probability: 756872327473 ( / 1070379110497 )

real    31m40.284s
user    31m39.002s
sys     0m1.902s

.
*/
func main() {
	var r, q uint64
	//test
	for q = 2; q < 1_000_000; q++ {
		r, q = Euler0100(q)
		// fmt.Printf("Euler 100: test: %d\n", r)
	}

	//run
	r, q = Euler0100(1_000_000_000_000)
	fmt.Printf("Euler 100: Arranged Probability: %d ( / %d )\n", r, q)
	// NOT correct 923187456
	if 756_872_327_473 != r {
		panic("Did not reach expected value.")
	}
}
