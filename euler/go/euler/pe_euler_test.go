// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package euler_test

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0 license https://creativecommons.org/licenses/by-sa/4.0/

import (
	// . "euler" // https://go.dev/wiki/CodeReviewComments#import-dot
	"container/heap"
	"euler" // I would __really__ like Go to just support path relative imports E.G. "./euler" means, just look in the CWD, it's here.
	"fmt"
	"testing"
)

/*-

	https://github.com/golang/go/issues/25223

	https://stackoverflow.com/questions/19998250/proper-package-naming-for-testing-with-the-go-language/31443271#31443271

	https://pkg.go.dev/cmd/go#hdr-Test_packages

	https://go.dev/src/strings/search_test.go

	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in $(seq 1 68) ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

	for ii in *\/*.go ; do go fmt "$ii" ; done ; go clean -testcache ; go test -v euler/
-*/

// https://pkg.go.dev/testing@go1.23.1

/*	https://pkg.go.dev/testing@go1.23.1#T
	func (c *T) Cleanup(f func())
	func (t *T) Deadline() (deadline time.Time, ok bool)
	func (c *T) Error(args ...any)
	func (c *T) Errorf(format string, args ...any)
	func (c *T) Fail()
	func (c *T) FailNow()
	func (c *T) Failed() bool
	func (c *T) Fatal(args ...any)
	func (c *T) Fatalf(format string, args ...any)
	func (c *T) Helper()
	func (c *T) Log(args ...any)
	func (c *T) Logf(format string, args ...any)
	func (c *T) Name() string
	func (t *T) Parallel()
	func (t *T) Run(name string, f func(t *T)) bool
	func (t *T) Setenv(key, value string)
	func (c *T) Skip(args ...any)
	func (c *T) SkipNow()
	func (c *T) Skipf(format string, args ...any)
	func (c *T) Skipped() bool
	func (c *T) TempDir() string
*/

// var primes = []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173}

// func TestXxx(t *testing.T) { t.  }

func TestDeprPrimesFactor(t *testing.T) {
	// prLim := len(primes)
	var pr *([]int)
	pr = euler.GetPrimes(pr, euler.PrimesSmallU8MxVal)
	for ii := 0; ii <= euler.PrimesSmallU8Mx; ii++ {
		if uint64((*pr)[ii]) != uint64(euler.PrimesSmallU8[ii]) {
			t.Errorf("Expected value %d got %d\t(%d / %d)\n", euler.PrimesSmallU8[ii], (*pr)[ii], ii, euler.PrimesSmallU8Mx)
		}
	}
	if euler.PrimesSmallU8[euler.PrimesSmallU8Mx] != euler.PrimesSmallU8MxVal {
		t.Errorf("euler.PrimesSmallU8MxVal set to %d got %d\n", euler.PrimesSmallU8MxVal, euler.PrimesSmallU8[euler.PrimesSmallU8Mx])
	}

	// Copied from TestFactorizeVsFactorMul
	tests := []struct {
		test uint
		ans  []uint
	}{
		{5309, []uint{5309}},
		{867, []uint{3, 17, 17}},
		{2024, []uint{2, 2, 2, 11, 23}},
		{1885, []uint{5, 13, 29}},
	}
	for _, test := range tests {
		pint := euler.Factor(pr, int(test.test))
		for ii := 0; ii < len(test.ans); ii++ {
			if uint((*pint)[ii]) != test.ans[ii] {
				t.Errorf("Expected factor %d got %d\t%v\n", test.ans[ii], (*pint)[ii], *pint)
			}
		}
	}
}

func TestBVPrimesPrimeAfter(t *testing.T) {
	//p := euler.NewBVPrimes()
	testsPrimeAfter := []struct{ q, r uint64 }{
		{2, 3},
		{3, 5},
		{5, 7},
		{16, 17},
		{18, 19},
		{70, 71},
		{71, 73},
		{502, 503},
		{503, 509},
		{509, 521},
		{522, 523},
		{524, 541},
		{542, 547},
		{548, 557},
		{558, 563},
		{564, 569},
		{570, 571},
		{7908, 7919},
		{7906, 7907},
		// {  , },
	}
	for ii := 0; ii+1 < euler.PrimesSmallU8Mx; ii++ {
		res := euler.Primes.PrimeAfter(uint64(euler.PrimesSmallU8[ii]))
		if uint64(euler.PrimesSmallU8[ii+1]) != res {
			t.Errorf("Bad PrimeAfter(%d) %d != %d (expected)\n", euler.PrimesSmallU8[ii], res, euler.PrimesSmallU8[ii+1])
		} else {
			// t.Logf("PASS PrimeAfter(%d) == %d\n", primes[ii], res)
		}
	}
	for _, tc := range testsPrimeAfter {
		res := euler.Primes.PrimeAfter(tc.q)
		if tc.r != res {
			t.Errorf("Bad PrimeAfter(%d) %d != %d (expected)\n", tc.q, res, tc.r)
		} else {
			// t.Logf("PASS PrimeAfter(%d) == %d\n", tc.q, res)
		}
	}
}

func TestGCDbin(t *testing.T) {
	tests := [][3]uint{
		{0, 18, 18},
		{6, 18, 6},
		{18, 48, 6},
		{12, 48, 12},
		{7, 11, 1},
	}
	for ii := 0; ii < len(tests); ii++ {
		res := euler.GCDbin(tests[ii][0], tests[ii][1])
		if tests[ii][2] != res {
			t.Errorf("Bad result: %d != %d (%d, %d)", res, tests[ii][2], tests[ii][0], tests[ii][1])
		}
	}
}

func TestFactorpairQueue(t *testing.T) {
	fq := &euler.FactorpairQueue{
		euler.Factorpair{Base: 19, Power: 1},
		euler.Factorpair{Base: 17, Power: 1},
		euler.Factorpair{Base: 13, Power: 1},
		euler.Factorpair{Base: 11, Power: 1},
		euler.Factorpair{Base: 7, Power: 1},
		euler.Factorpair{Base: 5, Power: 1},
		euler.Factorpair{Base: 3, Power: 1},
		euler.Factorpair{Base: 2, Power: 1},
	}
	heap.Init(fq)
	heap.Push(fq, euler.Factorpair{Base: 23, Power: 1})
	fqraw := fq.Raw()
	t.Logf("[0] = %d ~~~ [%d] = %d", (*fqraw)[0].Base, fq.Len(), (*fqraw)[fq.Len()-1].Base)
	mark := uint32(0)
	for 0 < fq.Len() {
		base := heap.Pop(fq).(euler.Factorpair).Base
		if mark > base {
			t.Errorf("Bad result, wanted > %d ; got < : %d", mark, base)
		}
		mark = base
	}
}

func TestFactorizeVsFactorMul(t *testing.T) {
	// Test cases
	// 1885 = 5 13 29
	// 2024 = 2^3 11 23
	// I _was_ going to use the meme phone number from the song... 8675309 but that's prime (thanks coreutils factor!)
	// 867 = 3 17^2
	// 5309 == prime
	tests := []struct {
		test uint64
		ans  []uint64
	}{
		{5309, []uint64{5309}},
		{867, []uint64{3, 17, 17}},
		{2024, []uint64{2, 2, 2, 11, 23}},
		{1885, []uint64{5, 13, 29}},
	}
	// p := euler.NewBVPrimes()
	for _, test := range tests {
		left := euler.Primes.Factorize(test.test)
		right := euler.Primes.Factorize(1)
		for _, subfact := range test.ans {
			right.Mul(euler.Primes.Factorize(subfact))
		}
		if false == left.Eq(right) || left.Uint64() != right.Uint64() {
			t.Errorf("Failed Test Case %v\n\t%d != %d\n%v\n%v", test, left.Uint64(), right.Uint64(), left, right)
		}
	}
}

//     pe_euler_test.go:182: Seed 9 for 19321
//     pe_euler_test.go:192: Failed to successfully factor: 32761 -> $ factor => 32761: 181 181
//     pe_euler_test.go:194: Profile complete: successful reseeds: 193     maximum: 9      failures: 1

func TestOverkillSeed(t *testing.T) {
	// const limit = uint( 0x100000 - 1) // 1M
	// const limit = uint( 0x200000 - 1) // 2M
	// const limit = uint( 0x400000 - 1) // 4M
	// const limit = uint(0x1000000 - 1) // 16777215 ~ VERY SLOW ~15-20min on my home NAS/server ... Unless it's REALLY important to quickly know if a number is a prime or not, probably too long.  Only ~8MB of mem though, so easily worth computing once, SAVING, and then loading if doing repeatedly.

	// const limit = 0x100000 + 64
	const limit = 1_250_000
	//const limit = uint(65535)
	// const limit = uint(1025)
	t.Logf("Seed primes upto: %d\n", limit)
	// p := euler.Primes
	euler.Primes.Grow(limit)
}

func TestPrimesChecksum(t *testing.T) {
	const limit = 0x100000
	//const limit = uint(65535)
	// const limit = uint(0x40000)
	// const limit = uint(1025)
	t.Logf("Verify primes upto: %d\n", limit)
	// p := euler.Primes
	euler.Primes.Grow(limit)
	priorProblematicResults := []struct {
		test  uint64
		prime bool
	}{{529409, false}, {532481, false}, {537601, false}, {587777, false}, {654337, false}, {713729, false}, {780289, false}, {801793, false}, {873473, false}, {903169, false}, {976897, false}, {998401, false}, {1039361, false}}
	tests := []struct {
		test uint64
		ans  uint64
	}{
		{72, 20},
		{71, 20},
		{173, 40},
		{281, 60},
		{409, 80},
		{541, 100},
		{660, 120},
		{7920, 1000},
		{0x100000, 82025},
		//{65535, 6542},      // WARNING: These numbers obtained by local calculation, not validated by external reference
		//{0x40000, 104846},  // WARNING:  These numbers obtained by local calculation, not validated by external reference
		//{0x100000, 498062}, // WARNING:  These numbers obtained by local calculation, not validated by external reference
	}
	// WARNING:  These numbers obtained by local calculation, not validated by external reference
	csums := []struct {
		end uint64
		sum uint64
	}{
		{1025, 171}, {2049, 137}, {3073, 130}, {4097, 125}, {5121, 121}, {6145, 116}, {7169, 115}, {8193, 112}, {9217, 114}, {10241, 112}, {11265, 108}, {12289, 108}, {13313, 111}, {14337, 100}, {15361, 114}, {16385, 105},
		{17409, 102}, {18433, 109}, {19457, 95}, {20481, 106}, {21505, 101}, {22529, 104}, {23553, 101}, {24577, 107}, {25601, 94}, {26625, 99}, {27649, 98}, {28673, 108}, {29697, 97}, {30721, 93}, {31745, 100}, {32769, 98},
		{33793, 107}, {34817, 97}, {35841, 92}, {36865, 100}, {37889, 99}, {38913, 91}, {39937, 100}, {40961, 91}, {41985, 102}, {43009, 104}, {44033, 88}, {45057, 95}, {46081, 89}, {47105, 91}, {48129, 98}, {49153, 95},
		{50177, 102}, {51201, 86}, {52225, 100}, {53249, 93}, {54273, 88}, {55297, 96}, {56321, 95}, {57345, 103}, {58369, 95}, {59393, 95}, {60417, 90}, {61441, 86}, {62465, 90}, {63489, 93}, {64513, 91}, {65537, 89},
		{66561, 93}, {67585, 98}, {68609, 86}, {69633, 86}, {70657, 94}, {71681, 95}, {72705, 95}, {73729, 91}, {74753, 90}, {75777, 94}, {76801, 86}, {77825, 98}, {78849, 84}, {79873, 92}, {80897, 91}, {81921, 93},
		{82945, 94}, {83969, 85}, {84993, 89}, {86017, 87}, {87041, 90}, {88065, 94}, {89089, 82}, {90113, 97}, {91137, 85}, {92161, 87}, {93185, 101}, {94209, 87}, {95233, 94}, {96257, 90}, {97281, 86}, {98305, 82},
		{99329, 94}, {100353, 85}, {101377, 91}, {102401, 96}, {103425, 81}, {104449, 85}, {105473, 91}, {106497, 89}, {107521, 85}, {108545, 88}, {109569, 92}, {110593, 83}, {111617, 86}, {112641, 89}, {113665, 87}, {114689, 86},
		{115713, 83}, {116737, 91}, {117761, 89}, {118785, 86}, {119809, 87}, {120833, 89}, {121857, 90}, {122881, 92}, {123905, 87}, {124929, 88}, {125953, 87}, {126977, 84}, {128001, 87}, {129025, 92}, {130049, 84}, {131073, 88},
		{132097, 82}, {133121, 93}, {134145, 81}, {135169, 82}, {136193, 92}, {137217, 96}, {138241, 85}, {139265, 79}, {140289, 91}, {141313, 93}, {142337, 87}, {143361, 77}, {144385, 83}, {145409, 80}, {146433, 94}, {147457, 86},
		{148481, 80}, {149505, 95}, {150529, 87}, {151553, 89}, {152577, 89}, {153601, 87}, {154625, 83}, {155649, 85}, {156673, 80}, {157697, 90}, {158721, 80}, {159745, 85}, {160769, 88}, {161793, 85}, {162817, 80}, {163841, 85},
		{164865, 86}, {165889, 81}, {166913, 80}, {167937, 87}, {168961, 74}, {169985, 87}, {171009, 90}, {172033, 85}, {173057, 88}, {174081, 85}, {175105, 82}, {176129, 79}, {177153, 93}, {178177, 76}, {179201, 91}, {180225, 93},
		{181249, 76}, {182273, 86}, {183297, 82}, {184321, 90}, {185345, 84}, {186369, 94}, {187393, 86}, {188417, 78}, {189441, 86}, {190465, 81}, {191489, 82}, {192513, 87}, {193537, 88}, {194561, 79}, {195585, 87}, {196609, 81},
		{197633, 85}, {198657, 88}, {199681, 80}, {200705, 81}, {201729, 80}, {202753, 88}, {203777, 79}, {204801, 82}, {205825, 83}, {206849, 85}, {207873, 89}, {208897, 86}, {209921, 88}, {210945, 91}, {211969, 86}, {212993, 70},
		{214017, 85}, {215041, 83}, {216065, 77}, {217089, 80}, {218113, 83}, {219137, 88}, {220161, 87}, {221185, 81}, {222209, 88}, {223233, 79}, {224257, 84}, {225281, 83}, {226305, 83}, {227329, 80}, {228353, 83}, {229377, 90},
		{230401, 97}, {231425, 76}, {232449, 80}, {233473, 77}, {234497, 78}, {235521, 83}, {236545, 76}, {237569, 76}, {238593, 87}, {239617, 82}, {240641, 82}, {241665, 86}, {242689, 88}, {243713, 80}, {244737, 86}, {245761, 82},
		{246785, 81}, {247809, 83}, {248833, 91}, {249857, 82}, {250881, 74}, {251905, 90}, {252929, 80}, {253953, 82}, {254977, 78}, {256001, 89}, {257025, 74}, {258049, 79}, {259073, 87}, {260097, 80}, {261121, 80}, {262145, 75},
		{263169, 79}, {264193, 89}, {265217, 79}, {266241, 85}, {267265, 81}, {268289, 96}, {269313, 80}, {270337, 84}, {271361, 82}, {272385, 86}, {273409, 77}, {274433, 74}, {275457, 83}, {276481, 87}, {277505, 78}, {278529, 78},
		{279553, 80}, {280577, 79}, {281601, 87}, {282625, 87}, {283649, 78}, {284673, 83}, {285697, 86}, {286721, 76}, {287745, 71}, {288769, 75}, {289793, 88}, {290817, 85}, {291841, 80}, {292865, 81}, {293889, 70}, {294913, 83},
		{295937, 78}, {296961, 83}, {297985, 76}, {299009, 75}, {300033, 81}, {301057, 87}, {302081, 80}, {303105, 81}, {304129, 86}, {305153, 90}, {306177, 79}, {307201, 86}, {308225, 72}, {309249, 80}, {310273, 80}, {311297, 78},
		{312321, 80}, {313345, 84}, {314369, 87}, {315393, 73}, {316417, 84}, {317441, 86}, {318465, 87}, {319489, 87}, {320513, 85}, {321537, 78}, {322561, 83}, {323585, 78}, {324609, 79}, {325633, 87}, {326657, 81}, {327681, 85},
		{328705, 79}, {329729, 83}, {330753, 81}, {331777, 81}, {332801, 77}, {333825, 85}, {334849, 76}, {335873, 80}, {336897, 83}, {337921, 84}, {338945, 79}, {339969, 70}, {340993, 84}, {342017, 78}, {343041, 79}, {344065, 82},
		{345089, 84}, {346113, 76}, {347137, 84}, {348161, 80}, {349185, 81}, {350209, 81}, {351233, 81}, {352257, 81}, {353281, 82}, {354305, 83}, {355329, 84}, {356353, 79}, {357377, 72}, {358401, 80}, {359425, 89}, {360449, 68},
		{361473, 73}, {362497, 86}, {363521, 74}, {364545, 85}, {365569, 84}, {366593, 79}, {367617, 82}, {368641, 79}, {369665, 66}, {370689, 80}, {371713, 70}, {372737, 78}, {373761, 81}, {374785, 78}, {375809, 93}, {376833, 78},
		{377857, 79}, {378881, 72}, {379905, 85}, {380929, 76}, {381953, 78}, {382977, 75}, {384001, 84}, {385025, 79}, {386049, 82}, {387073, 79}, {388097, 69}, {389121, 82}, {390145, 83}, {391169, 85}, {392193, 80}, {393217, 92},
		{394241, 82}, {395265, 81}, {396289, 71}, {397313, 78}, {398337, 79}, {399361, 81}, {400385, 83}, {401409, 71}, {402433, 75}, {403457, 75}, {404481, 80}, {405505, 71}, {406529, 78}, {407553, 74}, {408577, 85}, {409601, 80},
		{410625, 83}, {411649, 79}, {412673, 85}, {413697, 65}, {414721, 86}, {415745, 83}, {416769, 75}, {417793, 79}, {418817, 82}, {419841, 85}, {420865, 77}, {421889, 73}, {422913, 80}, {423937, 82}, {424961, 84}, {425985, 74},
		{427009, 74}, {428033, 81}, {429057, 70}, {430081, 93}, {431105, 78}, {432129, 82}, {433153, 86}, {434177, 78}, {435201, 88}, {436225, 82}, {437249, 76}, {438273, 70}, {439297, 74}, {440321, 81}, {441345, 83}, {442369, 79},
		{443393, 88}, {444417, 83}, {445441, 78}, {446465, 75}, {447489, 71}, {448513, 70}, {449537, 77}, {450561, 89}, {451585, 80}, {452609, 77}, {453633, 66}, {454657, 77}, {455681, 83}, {456705, 79}, {457729, 80}, {458753, 69},
		{459777, 78}, {460801, 70}, {461825, 77}, {462849, 73}, {463873, 79}, {464897, 83}, {465921, 81}, {466945, 71}, {467969, 85}, {468993, 86}, {470017, 72}, {471041, 89}, {472065, 81}, {473089, 72}, {474113, 83}, {475137, 79},
		{476161, 93}, {477185, 72}, {478209, 73}, {479233, 84}, {480257, 77}, {481281, 82}, {482305, 74}, {483329, 82}, {484353, 79}, {485377, 63}, {486401, 73}, {487425, 79}, {488449, 89}, {489473, 77}, {490497, 79}, {491521, 80},
		{492545, 72}, {493569, 80}, {494593, 80}, {495617, 88}, {496641, 76}, {497665, 78}, {498689, 76}, {499713, 86}, {500737, 77}, {501761, 78}, {502785, 71}, {503809, 72}, {504833, 77}, {505857, 89}, {506881, 75}, {507905, 74},
		{508929, 76}, {509953, 80}, {510977, 74}, {512001, 83}, {513025, 73}, {514049, 79}, {515073, 77}, {516097, 71}, {517121, 88}, {518145, 86}, {519169, 78}, {520193, 77}, {521217, 78}, {522241, 85}, {523265, 73}, {524289, 80},
		{525313, 75}, {526337, 83}, {527361, 83}, {528385, 72}, {529409, 78}, {530433, 79}, {531457, 75}, {532481, 76}, {533505, 79}, {534529, 75}, {535553, 70}, {536577, 84}, {537601, 78}, {538625, 77}, {539649, 75}, {540673, 73},
		{541697, 77}, {542721, 79}, {543745, 80}, {544769, 72}, {545793, 75}, {546817, 73}, {547841, 80}, {548865, 77}, {549889, 81}, {550913, 77}, {551937, 79}, {552961, 72}, {553985, 82}, {555009, 78}, {556033, 64}, {557057, 91},
		{558081, 67}, {559105, 76}, {560129, 83}, {561153, 85}, {562177, 58}, {563201, 92}, {564225, 65}, {565249, 78}, {566273, 81}, {567297, 70}, {568321, 80}, {569345, 75}, {570369, 77}, {571393, 86}, {572417, 76}, {573441, 80},
		{574465, 77}, {575489, 75}, {576513, 80}, {577537, 82}, {578561, 74}, {579585, 79}, {580609, 72}, {581633, 84}, {582657, 74}, {583681, 84}, {584705, 69}, {585729, 80}, {586753, 82}, {587777, 92}, {588801, 77}, {589825, 71},
		{590849, 73}, {591873, 69}, {592897, 79}, {593921, 79}, {594945, 79}, {595969, 81}, {596993, 80}, {598017, 74}, {599041, 75}, {600065, 72}, {601089, 81}, {602113, 76}, {603137, 82}, {604161, 74}, {605185, 79}, {606209, 76},
		{607233, 80}, {608257, 70}, {609281, 84}, {610305, 75}, {611329, 78}, {612353, 81}, {613377, 74}, {614401, 74}, {615425, 75}, {616449, 85}, {617473, 88}, {618497, 72}, {619521, 69}, {620545, 79}, {621569, 75}, {622593, 77},
		{623617, 68}, {624641, 85}, {625665, 77}, {626689, 66}, {627713, 75}, {628737, 69}, {629761, 74}, {630785, 73}, {631809, 78}, {632833, 82}, {633857, 74}, {634881, 83}, {635905, 77}, {636929, 77}, {637953, 80}, {638977, 57},
		{640001, 74}, {641025, 77}, {642049, 74}, {643073, 73}, {644097, 69}, {645121, 91}, {646145, 71}, {647169, 80}, {648193, 82}, {649217, 73}, {650241, 79}, {651265, 72}, {652289, 71}, {653313, 82}, {654337, 80}, {655361, 80},
		{656385, 70}, {657409, 73}, {658433, 77}, {659457, 72}, {660481, 82}, {661505, 82}, {662529, 71}, {663553, 73}, {664577, 82}, {665601, 85}, {666625, 70}, {667649, 77}, {668673, 67}, {669697, 77}, {670721, 78}, {671745, 75},
		{672769, 75}, {673793, 87}, {674817, 70}, {675841, 78}, {676865, 71}, {677889, 72}, {678913, 82}, {679937, 82}, {680961, 78}, {681985, 81}, {683009, 68}, {684033, 79}, {685057, 76}, {686081, 77}, {687105, 76}, {688129, 81},
		{689153, 77}, {690177, 76}, {691201, 67}, {692225, 77}, {693249, 78}, {694273, 76}, {695297, 81}, {696321, 73}, {697345, 79}, {698369, 76}, {699393, 72}, {700417, 72}, {701441, 72}, {702465, 80}, {703489, 79}, {704513, 77},
		{705537, 84}, {706561, 74}, {707585, 73}, {708609, 90}, {709633, 67}, {710657, 84}, {711681, 75}, {712705, 79}, {713729, 74}, {714753, 65}, {715777, 81}, {716801, 78}, {717825, 70}, {718849, 68}, {719873, 74}, {720897, 76},
		{721921, 83}, {722945, 74}, {723969, 87}, {724993, 78}, {726017, 71}, {727041, 78}, {728065, 75}, {729089, 78}, {730113, 76}, {731137, 69}, {732161, 79}, {733185, 75}, {734209, 77}, {735233, 79}, {736257, 78}, {737281, 71},
		{738305, 75}, {739329, 80}, {740353, 77}, {741377, 71}, {742401, 76}, {743425, 75}, {744449, 62}, {745473, 66}, {746497, 76}, {747521, 72}, {748545, 73}, {749569, 82}, {750593, 67}, {751617, 81}, {752641, 81}, {753665, 71},
		{754689, 82}, {755713, 78}, {756737, 71}, {757761, 67}, {758785, 78}, {759809, 81}, {760833, 74}, {761857, 84}, {762881, 69}, {763905, 80}, {764929, 74}, {765953, 86}, {766977, 80}, {768001, 69}, {769025, 80}, {770049, 75},
		{771073, 81}, {772097, 75}, {773121, 73}, {774145, 78}, {775169, 73}, {776193, 69}, {777217, 74}, {778241, 77}, {779265, 78}, {780289, 72}, {781313, 77}, {782337, 84}, {783361, 69}, {784385, 77}, {785409, 79}, {786433, 74},
		{787457, 69}, {788481, 84}, {789505, 71}, {790529, 77}, {791553, 75}, {792577, 79}, {793601, 73}, {794625, 85}, {795649, 69}, {796673, 79}, {797697, 81}, {798721, 75}, {799745, 69}, {800769, 87}, {801793, 70}, {802817, 72},
		{803841, 66}, {804865, 79}, {805889, 85}, {806913, 67}, {807937, 72}, {808961, 71}, {809985, 80}, {811009, 79}, {812033, 74}, {813057, 72}, {814081, 79}, {815105, 71}, {816129, 73}, {817153, 81}, {818177, 71}, {819201, 66},
		{820225, 83}, {821249, 80}, {822273, 74}, {823297, 83}, {824321, 78}, {825345, 80}, {826369, 78}, {827393, 74}, {828417, 74}, {829441, 67}, {830465, 78}, {831489, 73}, {832513, 85}, {833537, 79}, {834561, 65}, {835585, 71},
		{836609, 77}, {837633, 72}, {838657, 74}, {839681, 74}, {840705, 68}, {841729, 78}, {842753, 79}, {843777, 74}, {844801, 81}, {845825, 65}, {846849, 69}, {847873, 82}, {848897, 77}, {849921, 70}, {850945, 74}, {851969, 79},
		{852993, 74}, {854017, 68}, {855041, 71}, {856065, 76}, {857089, 84}, {858113, 74}, {859137, 72}, {860161, 74}, {861185, 80}, {862209, 75}, {863233, 78}, {864257, 73}, {865281, 77}, {866305, 79}, {867329, 73}, {868353, 79},
		{869377, 80}, {870401, 76}, {871425, 69}, {872449, 79}, {873473, 67}, {874497, 69}, {875521, 82}, {876545, 79}, {877569, 71}, {878593, 72}, {879617, 78}, {880641, 80}, {881665, 81}, {882689, 68}, {883713, 78}, {884737, 70},
		{885761, 70}, {886785, 83}, {887809, 74}, {888833, 78}, {889857, 67}, {890881, 72}, {891905, 83}, {892929, 71}, {893953, 75}, {894977, 75}, {896001, 81}, {897025, 67}, {898049, 75}, {899073, 77}, {900097, 64}, {901121, 76},
		{902145, 75}, {903169, 69}, {904193, 69}, {905217, 84}, {906241, 69}, {907265, 83}, {908289, 78}, {909313, 87}, {910337, 79}, {911361, 73}, {912385, 70}, {913409, 69}, {914433, 59}, {915457, 78}, {916481, 78}, {917505, 70},
		{918529, 82}, {919553, 73}, {920577, 74}, {921601, 73}, {922625, 78}, {923649, 77}, {924673, 76}, {925697, 87}, {926721, 75}, {927745, 67}, {928769, 69}, {929793, 82}, {930817, 71}, {931841, 66}, {932865, 79}, {933889, 71},
		{934913, 83}, {935937, 73}, {936961, 72}, {937985, 71}, {939009, 74}, {940033, 71}, {941057, 75}, {942081, 83}, {943105, 74}, {944129, 72}, {945153, 73}, {946177, 73}, {947201, 67}, {948225, 74}, {949249, 68}, {950273, 78},
		{951297, 78}, {952321, 74}, {953345, 73}, {954369, 71}, {955393, 76}, {956417, 72}, {957441, 65}, {958465, 70}, {959489, 78}, {960513, 67}, {961537, 79}, {962561, 80}, {963585, 74}, {964609, 80}, {965633, 76}, {966657, 73},
		{967681, 62}, {968705, 82}, {969729, 72}, {970753, 71}, {971777, 86}, {972801, 78}, {973825, 69}, {974849, 75}, {975873, 78}, {976897, 69}, {977921, 63}, {978945, 80}, {979969, 67}, {980993, 67}, {982017, 76}, {983041, 70},
		{984065, 75}, {985089, 72}, {986113, 75}, {987137, 83}, {988161, 63}, {989185, 68}, {990209, 74}, {991233, 80}, {992257, 73}, {993281, 71}, {994305, 70}, {995329, 73}, {996353, 78}, {997377, 80}, {998401, 73}, {999425, 73},
		{1000449, 75}, {1001473, 78}, {1002497, 76}, {1003521, 80}, {1004545, 70}, {1005569, 75}, {1006593, 79}, {1007617, 70}, {1008641, 86}, {1009665, 81}, {1010689, 62}, {1011713, 76}, {1012737, 86}, {1013761, 65}, {1014785, 75}, {1015809, 73},
		{1016833, 73}, {1017857, 79}, {1018881, 65}, {1019905, 84}, {1020929, 71}, {1021953, 81}, {1022977, 74}, {1024001, 71}, {1025025, 72}, {1026049, 81}, {1027073, 74}, {1028097, 71}, {1029121, 71}, {1030145, 81}, {1031169, 68}, {1032193, 62},
		{1033217, 78}, {1034241, 80}, {1035265, 73}, {1036289, 74}, {1037313, 72}, {1038337, 76}, {1039361, 68}, {1040385, 66}, {1041409, 79}, {1042433, 76}, {1043457, 66}, {1044481, 84}, {1045505, 76}, {1046529, 67}, {1047553, 75}, {1048577, 70},
	}
	abort := false
	for _, test := range tests {
		res := euler.Primes.CountPrimesLE(test.test)
		if res != test.ans {
			t.Errorf("expected count %d got %d", test.ans, res)
			abort = true
		}
	}
	for _, test := range csums {
		res := uint64(len(euler.Primes.PrimesOnPage(test.end)))
		if res != test.sum {
			t.Errorf("expected count %d got %d on cache line ending in %d", test.sum, res, test.end)
			abort = true
		}
	}
	for _, test := range priorProblematicResults {
		if test.prime != euler.Primes.KnownPrime(test.test) {
			t.Errorf("Incorrect result: %d %t should be %t", test.test, euler.Primes.KnownPrime(test.test), test.prime)
			abort = true
		}
	}

	t.Logf("Count at Limit %d", euler.Primes.CountPrimesLE(limit))
	if abort {
		panic("Prime Checksum failed")
	}
	_ = fmt.Sprintf("%s", "")
	/*
		out := ""
		for ii := 0; ii < len(euler.Primes.PV); ii++ {
			out += fmt.Sprintf("%+v", euler.Primes.PV[ii])
		}
		t.Logf("Logic Dump\n\n%s\n\n", out)

		csum := "\n"
		const BVl1 = 64
		var ii, end uint
		for end < limit {
			ii++
			end = (((ii) * BVl1) << 4) + 1
			// t.Logf("\t\t{%d, %d},\n", end, len(euler.Primes.PrimesOnPage(end)))
			ccount := len(euler.Primes.PrimesOnPage(end))
			csum += fmt.Sprintf("\t{%d, %d},", end, ccount)
			if 0 == ii & 0xF {
				csum += "\n"
			}
			if 512 == ccount {
				t.Errorf("Too many primes on a page, unset?")
				break
			}
		}
		t.Log(csum)
	*/
}

func TestOverkillVerifyFactor1980Pollard(t *testing.T) {
	var reseed, maxseed, failseed uint64
	const limit = 65535
	t.Logf("Verify / Profile Factor1980PollardMonteCarlo(p, seed) upto: %d\n", limit)
	// p := euler.Primes
	euler.Primes.Grow(limit)
	reseed, maxseed, failseed = 0, 0, 0
TestOverkillVerifyOuter:
	for ii := uint64(2); ii <= limit; ii++ {
		if 0 == ii&0x3fff {
			t.Logf("... %d", ii)
		}
		for seed := uint64(0); seed <= 20; seed++ {
			factor := euler.Factor1980PollardMonteCarlo(ii, seed)
			if 1 < factor || euler.Primes.KnownPrime(ii) {
				if 0 != seed {
					if testing.Verbose() {
						// t.Logf("Seed %d for %d", seed, ii)
					}
					reseed++
					if maxseed < seed {
						maxseed = seed
					}
				}
				continue TestOverkillVerifyOuter
			}
		}
		failseed++
		seedMax := uint64(1)
		for t := ii >> 1; seedMax < t; t >>= 1 {
			seedMax <<= 1
		}
		if testing.Verbose() {
			// func Factor1980AutoPMC(q uint) uint {...}
			// Resolves the only soft fail 32761 with the >= square root limit test // 17 (289) 79 (6241) 139 (19321) 181 (32761)
			t.Logf("SOFT fail to successfully factor: %d, iterating from seed 21 to %d", ii, seedMax)
		}
		for seed := uint64(21); seed <= seedMax; seed++ {
			factor := euler.Factor1980PollardMonteCarlo(ii, seed)
			if 1 < factor {
				if testing.Verbose() {
					t.Logf("Seed %d for %d", seed, ii)
				}
				reseed++
				if maxseed < seed {
					maxseed = seed
				}
				continue TestOverkillVerifyOuter
			}
		}
		t.Errorf("FAILED to successfully factor: %d", ii)
		// _ = euler.Factor1980AutoPMC(ii)
	}
	t.Logf("Profile complete: successful reseeds: %d\tmaximum: %d\tfailures: %d", reseed, maxseed, failseed)
}

func TestOverkillVerifyFactor1980AutoPMC(t *testing.T) {
	t.Skip("Slow")
	// const limit = 65535
	const limit = 1_250_000
	t.Logf("Verify Factor1980AutoPMC(p) upto: %d\n", limit)
	// p := euler.Primes
	euler.Primes.Grow(limit)
	for ii := uint64(2); ii <= limit; ii += 2 {
		if 0 == ii%3 || euler.Primes.KnownPrime(uint64(ii)) {
			continue
		}
		f := euler.Factor1980AutoPMC(ii, true)
		if false == euler.Primes.KnownPrime(f) {
			t.Errorf("FAILED to successfully factor: %d ~ %d", ii, f)
			t.Fatal("")
		}
	}
}

func TestOverkillVerifyFactorLenstraECW(t *testing.T) {
	t.Skip("Slow")
	// go test -run TestOverkillVerifyFactorLenstraECW -cpuprofile Lenstra.goprof -v euler/
	// go tool pprof euler.test Lenstra.goprof
	// const limit = 65535
	const limit = 1_250_000
	t.Logf("Verify FactorLenstraECW(p, 2) (force the extra checks path too) upto: %d\n", limit)
	// p := euler.Primes
	euler.Primes.Grow(limit)
	for ii := int64(3); ii <= limit; ii += 2 {
		if 0 == ii%3 || euler.Primes.KnownPrime(uint64(ii)) {
			continue
		}
		f := euler.FactorLenstraECW(ii, 2)
		if 2 > f || 0 != ii%f {
			t.Errorf("FactorLenstraECW failed to factor %d, returned %d\n", ii, f)
		}
	}
}

func TestOverkillVerifyPrimeTests(t *testing.T) {
	// t.Skipf("Known Broken - working on E.G. PowIntMod\n")
	const limit = 750_000
	var errors, factor uint64
	euler.Primes.Grow(limit)

	testProbablyPrime := []struct {
		num     uint64
		isPrime bool
	}{
		//{ , },
		{5915587277, true},
		{1500450271, true},
		{3267000013, true},
		{5754853343, true},
		{4093082899, true},
		{9576890767, true},
		{3628273133, true},
		{2860486313, true},
		{5463458053, true},
		{3367900313, true},

		// Test times out... +600 seconds
		// {12764787846358441471, true},
	}
	// https://t5k.org/lists/2small/0bit.html
	testPowPrimes := []struct {
		shift uint8
		subK  []uint16
	}{
		{8, []uint16{5, 15, 17, 23, 27, 29, 33, 45, 57, 59}},
		{9, []uint16{3, 9, 13, 21, 25, 33, 45, 49, 51, 55}},
		{10, []uint16{3, 5, 11, 15, 27, 33, 41, 47, 53, 57}},
		{11, []uint16{9, 19, 21, 31, 37, 45, 49, 51, 55, 61}},
		{12, []uint16{3, 5, 17, 23, 39, 45, 47, 69, 75, 77}},
		{13, []uint16{1, 13, 21, 25, 31, 45, 69, 75, 81, 91}},
		{14, []uint16{3, 15, 21, 23, 35, 45, 51, 65, 83, 111}},
		{15, []uint16{19, 49, 51, 55, 61, 75, 81, 115, 121, 135}},
		{16, []uint16{15, 17, 39, 57, 87, 89, 99, 113, 117, 123}},
		{17, []uint16{1, 9, 13, 31, 49, 61, 63, 85, 91, 99}},
		{18, []uint16{5, 11, 17, 23, 33, 35, 41, 65, 75, 93}},
		{19, []uint16{1, 19, 27, 31, 45, 57, 67, 69, 85, 87}},
		{20, []uint16{3, 5, 17, 27, 59, 69, 129, 143, 153, 185}},
		{21, []uint16{9, 19, 21, 55, 61, 69, 105, 111, 121, 129}},
		{22, []uint16{3, 17, 27, 33, 57, 87, 105, 113, 117, 123}},
		{23, []uint16{15, 21, 27, 37, 61, 69, 135, 147, 157, 159}},
		{24, []uint16{3, 17, 33, 63, 75, 77, 89, 95, 117, 167}},
		{25, []uint16{39, 49, 61, 85, 91, 115, 141, 159, 165, 183}},
		{26, []uint16{5, 27, 45, 87, 101, 107, 111, 117, 125, 135}},
		{27, []uint16{39, 79, 111, 115, 135, 187, 199, 219, 231, 235}},
		{28, []uint16{57, 89, 95, 119, 125, 143, 165, 183, 213, 273}},
		{29, []uint16{3, 33, 43, 63, 73, 75, 93, 99, 121, 133}},
		{30, []uint16{35, 41, 83, 101, 105, 107, 135, 153, 161, 173}},
		{31, []uint16{1, 19, 61, 69, 85, 99, 105, 151, 159, 171}},
		{32, []uint16{5, 17, 65, 99, 107, 135, 153, 185, 209, 267}},
		{33, []uint16{9, 25, 49, 79, 105, 285, 301, 303, 321, 355}},
		{34, []uint16{41, 77, 113, 131, 143, 165, 185, 207, 227, 281}},
		{35, []uint16{31, 49, 61, 69, 79, 121, 141, 247, 309, 325}},
		{36, []uint16{5, 17, 23, 65, 117, 137, 159, 173, 189, 233}},
		{37, []uint16{25, 31, 45, 69, 123, 141, 199, 201, 351, 375}},
		{38, []uint16{45, 87, 107, 131, 153, 185, 191, 227, 231, 257}},
		{39, []uint16{7, 19, 67, 91, 135, 165, 219, 231, 241, 301}},
		{40, []uint16{87, 167, 195, 203, 213, 285, 293, 299, 389, 437}},
		{41, []uint16{21, 31, 55, 63, 73, 75, 91, 111, 133, 139}},
		{42, []uint16{11, 17, 33, 53, 65, 143, 161, 165, 215, 227}},
		{43, []uint16{57, 67, 117, 175, 255, 267, 291, 309, 319, 369}},
		{44, []uint16{17, 117, 119, 129, 143, 149, 287, 327, 359, 377}},
		{45, []uint16{55, 69, 81, 93, 121, 133, 139, 159, 193, 229}},
		{46, []uint16{21, 57, 63, 77, 167, 197, 237, 287, 305, 311}},
		{47, []uint16{115, 127, 147, 279, 297, 339, 435, 541, 619, 649}},
		{48, []uint16{59, 65, 89, 93, 147, 165, 189, 233, 243, 257}},
		{49, []uint16{81, 111, 123, 139, 181, 201, 213, 265, 283, 339}},
		{50, []uint16{27, 35, 51, 71, 113, 117, 131, 161, 195, 233}},
		{51, []uint16{129, 139, 165, 231, 237, 247, 355, 391, 397, 439}},
		{52, []uint16{47, 143, 173, 183, 197, 209, 269, 285, 335, 395}},
		{53, []uint16{111, 145, 231, 265, 315, 339, 343, 369, 379, 421}},
		{54, []uint16{33, 53, 131, 165, 195, 245, 255, 257, 315, 327}},
		{55, []uint16{55, 67, 99, 127, 147, 169, 171, 199, 207, 267}},
		{56, []uint16{5, 27, 47, 57, 89, 93, 147, 177, 189, 195}},
		{57, []uint16{13, 25, 49, 61, 69, 111, 195, 273, 363, 423}},
		{58, []uint16{27, 57, 63, 137, 141, 147, 161, 203, 213, 251}},
		{59, []uint16{55, 99, 225, 427, 517, 607, 649, 687, 861, 871}},
		{60, []uint16{93, 107, 173, 179, 257, 279, 369, 395, 399, 453}},
		{61, []uint16{1, 31, 45, 229, 259, 283, 339, 391, 403, 465}},
		{62, []uint16{57, 87, 117, 143, 153, 167, 171, 195, 203, 273}},
		{63, []uint16{25, 165, 259, 301, 375, 387, 391, 409, 457, 471}},
		{64, []uint16{59, 83, 95, 179, 189, 257, 279, 323, 353, 363}},
	}
	for _, row := range testPowPrimes {
		if 64 == row.shift {
			// 64 bits is the integer size, must use modular math wraparound (mod 2^64)
			for _, val := range row.subK {
				testProbablyPrime = append(testProbablyPrime, struct {
					num     uint64
					isPrime bool
				}{uint64(0) - uint64(val), true})
			}
			continue
		}
		for _, val := range row.subK {
			testProbablyPrime = append(testProbablyPrime, struct {
				num     uint64
				isPrime bool
			}{(uint64(1) << row.shift) - uint64(val), true})
		}
	}
	_ = testProbablyPrime
	t.Logf("\n\nSKIPPING the SOME large PrimeOptiTestMillerRabin tests, still working on the 128bit UU64DivQD for slow-path modulus required for numbers >32 bits\n\n")
	for _, test := range testProbablyPrime {
		if 0x1_0000_0000 < test.num {
			continue
		}
		if test.isPrime != (0 == euler.PrimeOptiTestMillerRabin(test.num)) {
			t.Errorf("PrimeOptiTestMillerRabin returned incorrect result for %d\n", test.num)
		}
	}
	for ii := uint64(3); ii <= limit && 10 > errors; ii++ {
		if 0 == ii&0xFFFF {
			t.Logf("\t@%d", ii)
		}
		factor = euler.PrimeOptiTestMillerRabin(ii)
		if euler.Primes.KnownPrime(ii) != (0 == factor) {
			t.Errorf("PrimeOptiTestMillerRabin returned incorrect result for %d, got factor %d\n", ii, factor)
			errors++
		}
	}
}

func TestOverkillPrimesLists(t *testing.T) {
	const limit = 65536 << 1
	// p := euler.Primes
	euler.Primes.PrimeGlobalList(limit)
	var prime uint64
	var ii, iiLim int
	prime = 1
	for ii = 0; ii <= euler.PrimesSmallU8Mx; ii++ {
		prime = euler.Primes.PrimeAfter(prime)
		if euler.PrimesSmallU8[ii] != uint8(prime) {
			t.Errorf("PrimesLists U8 expected %d returned %d\n", prime, euler.PrimesSmallU8[ii])
		}
	}
	iiLim = len(euler.PrimesMoreU16)
	for ii = 0; ii < iiLim; ii++ {
		prime = euler.Primes.PrimeAfter(prime)
		if euler.PrimesMoreU16[ii] != uint16(prime) {
			t.Errorf("PrimesLists U16 expected %d returned %d\n", prime, euler.PrimesMoreU16[ii])
		}
	}
	iiLim = len(euler.PrimesMoreU32)
	for ii = 0; ii < iiLim; ii++ {
		prime = euler.Primes.PrimeAfter(prime)
		if euler.PrimesMoreU32[ii] != uint32(prime) {
			t.Errorf("PrimesLists U32 expected %d returned %d\n", prime, euler.PrimesMoreU32[ii])
		}
	}
}

func TestFactorizeProperDivisors(t *testing.T) {
	tests := []struct {
		test uint64
		ans  []uint64
	}{
		{5309, []uint64{1}},
		{16, []uint64{1, 2, 4, 8}},
		{867, []uint64{1, 3, 17, 51, 289}},
	}
	// {1885, []uint64{1, 5, 13, 29}},
	// {2024, []uint64{1, 2, 2, 2, 11, 23}},
	// p := euler.Primes
	for _, test := range tests {
		propdiv := *(euler.Primes.Factorize(uint64(test.test)).ProperDivisors())
		if len(test.ans) != len(propdiv) {
			t.Errorf("Lengths do not match:\n%v\n%v\n", test.ans, propdiv)
		}
		for ii := 0; ii < len(propdiv); ii++ {
			if test.ans[ii] != propdiv[ii] {
				t.Errorf("Factor mismatch: [%d] expected %d got %d\n%v\n%v\n", ii, test.ans[ii], propdiv[ii], test.ans, propdiv)
			}
		}
	}
}

/*
func TestMegaTestPrimes0x100000(t *testing.T) {
	t.Fatalf("TODO: Write this test.")
	// I manually tweaked the generator to run several passes and identified why any differences happened.
	// This is now covered by per 'cache line' sized chunk checksums.
}
*/

func TestRotateDecDigits(t *testing.T) {
	tests := []struct {
		test uint64
		ans  []uint64
	}{
		{5309, []uint64{5309, 9530, 953, 3095}},
		{16, []uint64{16, 61}},
		{867, []uint64{867, 786, 678}},
	}
	for _, test := range tests {
		rots := euler.RotateDecDigits(test.test)
		if len(test.ans) != len(rots) {
			t.Errorf("Lengths do not match:\n%v\n%v\n", test.ans, rots)
		}
		for ii := 0; ii < len(rots); ii++ {
			if test.ans[ii] != rots[ii] {
				t.Errorf("Rotation mismatch: [%d] expected %d got %d\n%v\n%v\n", ii, test.ans[ii], rots[ii], test.ans, rots)
			}
		}
	}
}

func TestPalindromeFuncs(t *testing.T) {
	testFlipBin := []struct {
		test, ans uint64
	}{
		{0, 0},
		{0b_1, 0b_1},
		{0b_101, 0b_101},
		{0b_1001, 0b_1001},
		{0b_11011, 0b_11011},
		{0b_10111, 0b_11101},
		{0xF000000000000000, 0x000000000000000F},
	}
	for _, test := range testFlipBin {
		res := euler.PalindromeFlipBinary(test.test)
		if res != test.ans {
			t.Errorf("Flip mismatch: expected %b got %b\n", test.ans, res)
		}
	}
	testMakeDec := []struct {
		test, addZeros, even, odd uint64
	}{
		{1, 0, 11, 1},
		{9, 0, 99, 9},
		{9, 1, 9009, 909},
		{42, 0, 2442, 242},
		{867, 0, 768867, 76867},
		{5309, 0, 90355309, 9035309},
	}
	for _, test := range testMakeDec {
		even := euler.PalindromeMakeDec(test.test, test.addZeros, false)
		odd := euler.PalindromeMakeDec(test.test, test.addZeros, true)
		if even != test.even || odd != test.odd {
			t.Errorf("Even and/or odd failed: Expected %d and %d got %d and %d\n", test.odd, test.even, odd, even)
		}
	}

	testPals := []struct {
		in, base uint64
		pal      bool
	}{
		{11, 10, true},
		{9009, 10, true},
		{909, 10, true},
		{90355309, 10, true},
		{90353309, 10, false},
		{0b_11011, 2, true},
		{0b_111011, 2, false},
	}
	for _, test := range testPals {
		if test.pal != euler.IsPalindrome(test.in, test.base) {
			t.Errorf("expected %t, got %t for %d\n", test.pal, euler.IsPalindrome(test.in, test.base), test.in)
		}
	}

}

func TestPandigital(t *testing.T) {
	testPan := []struct {
		test                 uint64
		reserved, bset, used uint16
		DigitShift           uint64
	}{
		{0, 1, 0, 0b1, 10},
		{1, 1, 1, 0b10, 10},
		{12, 1, 2, 0b110, 100},
		{123, 1, 3, 0b1110, 1000},
		{1234, 1, 4, 0b11110, 10_000},
		{12345, 1, 5, 0b111110, 100_000},
		{123456789, 1, 9, 0b1111_111110, 1_000_000_000},
		{1023456789, 1, 0, 0b1111_111111, 10_000_000_000},
	}
	for _, test := range testPan {
		// func Pandigital(test uint64, used, reserved uint16) (fullPD bool, biton, usedRe uint16, DigitShift uint64) {
		fullPD, bset, used, DigitShift := euler.Pandigital(test.test, 0, test.reserved)
		if used != test.used || bset != test.bset || DigitShift != test.DigitShift || (uint16((uint64(1)<<(bset+1))-2) == used) != fullPD {
			t.Errorf("Pandigital: expected %v got %d, %d, %d\n", test, used, bset, DigitShift)
		}
	}
}

func TestSlicePop(t *testing.T) {
	testSlicePop := []struct {
		deck     []uint8
		requests []int
		result   []uint8
	}{
		{[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{0, 0, 0, 0, 0, 0, 0, 0, 0}, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{1, 1, 1, 1, 1, 1, 1, 1, 0}, []uint8{2, 3, 4, 5, 6, 7, 8, 9, 1}},
		{[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{8, 7, 6, 5, 4, 3, 2, 1, 0}, []uint8{9, 8, 7, 6, 5, 4, 3, 2, 1}},
		{[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{4, 4, 4, 3, 0, 0, 0, 0, 0}, []uint8{5, 6, 7, 4, 1, 2, 3, 8, 9}},
	}
	for tt, test := range testSlicePop {
		clone := make([]uint8, len(test.deck))
		copy(clone, test.deck)
		res := make([]uint8, 0, len(test.result))
		for ii := 0; ii < len(test.requests); ii++ {
			res = append(res, euler.SlicePopUint8(clone, test.requests[ii]))
		}
		if len(res) != len(test.result) {
			t.Errorf("Length mismatch, expected %d, got %d\n", len(test.result), len(res))
		}
		for ii := 0; ii < len(res); ii++ {
			if res[ii] != test.result[ii] {
				t.Errorf("Item mismatch in test %d index %d: expected %d, got %d\n", tt, ii, test.result[ii], res[ii])
			}
		}
	}
}

func TestNgonalNumbers(t *testing.T) {
	var ii, pi, pfi uint64
	for ii = 0; ii < 101; ii++ {
		pi = euler.TriangleNumber(ii)
		pfi = euler.TriangleNumberReverseFloor(pi)
		if ii != pfi {
			t.Errorf("Loop failed: %d == TriangleNumberReverseFloor( TriangleNumber() ~ %d ) got %d\n", ii, pi, pfi)
		}
		pi = euler.SquareNumber(ii)
		pfi = euler.SquareNumberReverseFloor(pi)
		if ii != pfi {
			t.Errorf("Loop failed: %d == SquareNumberReverseFloor( SquareNumber() ~ %d ) got %d\n", ii, pi, pfi)
		}
		pi = euler.PentagonalNumber(ii)
		pfi = euler.PentagonalNumberReverseFloor(pi)
		if ii != pfi {
			t.Errorf("Loop failed: %d == PentagonalNumberReverseFloor( PentagonalNumber() ~ %d ) got %d\n", ii, pi, pfi)
		}
		pi = euler.HexagonalNumber(ii)
		pfi = euler.HexagonalNumberReverseFloor(pi)
		if ii != pfi {
			t.Errorf("Loop failed: %d == HexagonalNumberReverseFloor( HexagonalNumber() ~ %d ) got %d\n", ii, pi, pfi)
		}
		pi = euler.HeptagonalNumber(ii)
		pfi = euler.HeptagonalNumberReverseFloor(pi)
		if ii != pfi {
			t.Errorf("Loop failed: %d == HeptagonalNumberReverseFloor( HeptagonalNumber() ~ %d ) got %d\n", ii, pi, pfi)
		}
		pi = euler.OctagonalNumber(ii)
		pfi = euler.OctagonalNumberReverseFloor(pi)
		if ii != pfi {
			t.Errorf("Loop failed: %d == OctagonalNumberReverseFloor( OctagonalNumber() ~ %d ) got %d\n", ii, pi, pfi)
		}
	}
	test3gon := []struct {
		N, GN uint64
	}{
		{1, 1},
		{2, 3},
		{3, 6},
		{4, 10},
		{5, 15},
	}
	for _, test := range test3gon {
		rev := euler.TriangleNumberReverseFloor(test.GN)
		res := euler.TriangleNumber(test.N)
		if res != test.GN || rev != test.N {
			t.Errorf("6gon failed, expected %d, %d got %d, %d\n", test.N, test.GN, rev, res)
		}
	}
	test4gon := []struct {
		N, GN uint64
	}{
		{1, 1},
		{2, 4},
		{3, 9},
		{4, 16},
		{5, 25},
	}
	for _, test := range test4gon {
		rev := euler.SquareNumberReverseFloor(test.GN)
		res := euler.SquareNumber(test.N)
		if res != test.GN || rev != test.N {
			t.Errorf("6gon failed, expected %d, %d got %d, %d\n", test.N, test.GN, rev, res)
		}
	}
	test5gon := []struct {
		N, GN uint64
	}{
		{1, 1},
		{2, 5},
		{3, 12},
		{4, 22},
		{5, 35},
	}
	for _, test := range test5gon {
		rev := euler.PentagonalNumberReverseFloor(test.GN)
		res := euler.PentagonalNumber(test.N)
		if res != test.GN || rev != test.N {
			t.Errorf("6gon failed, expected %d, %d got %d, %d\n", test.N, test.GN, rev, res)
		}
	}
	test6gon := []struct {
		N, GN uint64
	}{
		{1, 1},
		{2, 6},
		{3, 15},
		{4, 28},
		{5, 45},
	}
	for _, test := range test6gon {
		rev := euler.HexagonalNumberReverseFloor(test.GN)
		res := euler.HexagonalNumber(test.N)
		if res != test.GN || rev != test.N {
			t.Errorf("6gon failed, expected %d, %d got %d, %d\n", test.N, test.GN, rev, res)
		}
	}
	test7gon := []struct {
		N, GN uint64
	}{
		{1, 1},
		{2, 7},
		{3, 18},
		{4, 34},
		{5, 55},
	}
	for _, test := range test7gon {
		rev := euler.HeptagonalNumberReverseFloor(test.GN)
		res := euler.HeptagonalNumber(test.N)
		if res != test.GN || rev != test.N {
			t.Errorf("7gon failed, expected %d, %d got %d, %d\n", test.N, test.GN, rev, res)
		}
	}
	test8gon := []struct {
		N, GN uint64
	}{
		{1, 1},
		{2, 8},
		{3, 21},
		{4, 40},
		{5, 65},
	}
	for _, test := range test8gon {
		rev := euler.OctagonalNumberReverseFloor(test.GN)
		res := euler.OctagonalNumber(test.N)
		if res != test.GN || rev != test.N {
			t.Errorf("8gon failed, expected %d, %d got %d, %d\n", test.N, test.GN, rev, res)
		}
	}
}

//type Card uint8

//func (o []Card) SLUint8() []uint8 { return []uint8(o) }

func TestBaseConversions(t *testing.T) {
	for ii := uint64(1); ii < 0xFFFFFF; ii *= 3 * 5 * 7 {
		for bb := uint64(2); bb <= 16; bb++ {
			tt := euler.Uint8DigitsToUint64(euler.Uint64ToDigitsUint8(ii, bb), bb)
			if ii != tt {
				t.Errorf("Loop Convert: %d base %d : expected %d got %d\n", ii, bb, ii, tt)
			}
		}
	}
	testUint8Sort := []struct {
		src  []uint8
		res  []uint8
		same bool
	}{
		{[]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, true},
		{[]uint8{2, 3, 4, 5, 6, 7, 8, 9, 1}, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, true},
		{[]uint8{9, 8, 7, 6, 5, 4, 3, 2, 1}, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, true},
		{[]uint8{5, 6, 7, 4, 1, 2, 3, 8, 0}, []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8}, true},
		{[]uint8{5, 6, 7, 4, 1, 2, 3, 8, 0}, []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, false},
	}
	for _, test := range testUint8Sort {
		res := euler.Uint8CopyInsertSort(test.src)
		cmp := euler.Uint8Compare(res, test.res)
		if test.same != (0 == cmp) {
			t.Errorf("Expected results: %t %v ~~ got %v %t\n", test.same, test.res, res, 0 == cmp)
		}
	}
	testConcatU64 := []struct {
		x, y, base, res uint64
	}{
		{1, 0, 10, 10},
		{1, 1, 10, 11},
		{10, 1, 10, 101},
		{10, 9, 10, 109},
		{9, 10, 10, 910},
	}
	for _, test := range testConcatU64 {
		if test.res != euler.ConcatDigitsU64(test.x, test.y, test.base) {
			t.Errorf("ConcatDigitsU64 expected %d got %d\n", test.res, euler.ConcatDigitsU64(test.x, test.y, test.base))
		}
	}

	// combomax := euler.FactorialUint64(3)
	// for ii := uint64(0) ; ii < combomax ; ii++ {
	//	t.Logf("PermutationSlUint8 %d : %v\n", ii, euler.PermutationSlUint8(ii, []uint8{0, 1, 2}))
	// }
}

func TestFactorial(t *testing.T) {
	testFactorial := []struct {
		in, div, res uint64
	}{
		{2, 1, 2},
		{9, 1, 362_880},
		{9, 5, 3024},
	}
	var res uint64
	for _, test := range testFactorial {
		res = euler.FactorialDivFactU64toBig(test.in, test.div).Uint64()
		if res != test.res {
			t.Errorf("FactorialDivFactU64toBig: Expected results: %d got %d\n", test.res, res)
		}
		res = uint64(euler.FactorialUint64(test.in) / euler.FactorialUint64(test.div))
		if res != test.res {
			t.Errorf("FactorialUint64: Expected results: %d got %d\n", test.res, res)
		}
		res = uint64(euler.Factorial(int(test.in)) / euler.Factorial(int(test.div)))
		if res != test.res {
			t.Errorf("Factorial: Expected results: %d got %d\n", test.res, res)
		}
	}
}

func TestCardsPoker(t *testing.T) {
	//	Val	Suit	Num	Name	Desc
	//  0x1FF0_0000	same5	inc5	Royal (straight) Flush	AceHigh unbroken run of cards in the same suit
	//  0x1FF0_0000	same5	inc5	Straight Flush	unbroken run of cards in the same suit
	//  0x_F0F_0000	-	same4	4 of a kind (value)	Note the same 0xF (15) for 'card' value in Flush slot
	//  0x_F00_FF00	-	same3+2	Full House	3+2 of a kind (value) -- Uses 0xF (15) for 'card' value in Flush slot
	//  0x_F00_0000	same5	-	Flush	All cards in the same suit, but no other match
	//  0x__F0_0000	-	inc5	Straight All cards in sequence
	//	0x_F000	-	same3	Three of a kind
	//	0x__FF0	-	same2	Two Pair	2+2 of a kind
	//	0x___F0	-	same2	One Pair	2 of a kind
	//	0x____F		highest	High Card
	testCardsPoker := []struct {
		hand  []string
		pub   []string
		score uint
	}{
		{[]string{"2C", "3C", "4C", "5C", "7D"}, []string{}, 0x7},
		{[]string{"2C", "3C", "4C", "5C", "6C"}, []string{}, 0x1660_0000},
		{[]string{"2C", "2D", "4C", "5C", "6C"}, []string{}, 0x0026},
		{[]string{"2C", "2D", "3C", "3H", "6C"}, []string{}, 0x0326},
		{[]string{"2C", "2D", "2H", "5C", "7C"}, []string{}, 0x2007},
		{[]string{"2S", "3C", "4C", "5C", "6C"}, []string{}, 0x0_60_0000},
		{[]string{"2C", "3C", "4C", "5C", "7C"}, []string{}, 0x_700_0000},
		{[]string{"2C", "2D", "3C", "3H", "3S"}, []string{}, 0x_F00_3200},
		{[]string{"2C", "2D", "2H", "2S", "7C"}, []string{}, 0x_F02_0007},
		{[]string{"AC", "KC", "QC", "JC", "TC"}, []string{}, 0x1EE0_0000},
	}

	for _, test := range testCardsPoker {
		cards, pub := make([]uint8, 0), make([]uint8, 0)
		for _, card := range test.hand {
			cards = append(cards, euler.CardParseENG(card))
		}
		for _, card := range test.pub {
			pub = append(pub, euler.CardParseENG(card))
		}
		score := euler.CardPokerScore(cards, pub)
		if test.score != score {
			t.Errorf("Expected score %8x got score %8x: %v %v\n", test.score, score, test.hand, test.pub)
		}
	}
}

func TestSliceFuncs(t *testing.T) {
	testBsearchSl := []struct {
		sl     []uint8
		index  int
		target uint8
	}{
		{[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 4, 4},
		{[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, -1, 255},
		{[]uint8{0, 1, 2, 3, 5, 6, 7, 8, 9}, -1, 4},
		// {[]uint8{}, , },
	}
	for _, test := range testBsearchSl {
		if test.index != euler.BsearchSlice(test.sl, test.target, false) {
			t.Errorf("Bsearch: Expected index %d got %d\n", test.index, euler.BsearchSlice(test.sl, test.target, false))
		}
	}
	testSlCom := []struct {
		a, b, c []uint8
	}{
		{[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []uint8{3, 6, 9}, []uint8{3, 6, 9}},
		{[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []uint8{3, 6, 9, 12}, []uint8{3, 6, 9}},
		{[]uint8{0, 1, 2, 4, 5, 7, 8, 9}, []uint8{3, 6, 9}, []uint8{9}},
		{[]uint8{0, 1, 2, 4, 5, 7, 8, 9}, []uint8{3, 6}, []uint8{}},
		// {[]uint8{}, , },
	}
	for _, test := range testSlCom {
		res := euler.SliceCommon(test.a, test.b)
		if len(res) != len(test.c) {
			t.Errorf("SliceCommon: wanted result length %d, got %d\n", len(test.c), len(res))
			continue
		}
		for ii := 0; ii < len(res); ii++ {
			if test.c[ii] != res[ii] {
				t.Errorf("SliceCommon: differ at index %d: %d != %d\n", ii, test.c[ii], res[ii])
				break
			}
		}
	}
}

func TestEulerTotientPhi(t *testing.T) {
	// https://en.wikipedia.org/wiki/Euler's_totient_function#Proof_of_Euler's_product_formula
	testTotientPhi := []uint64{
		0, 1, 1, 2, 2, 4, 2, 6, 4, 6, 4,
		10, 4, 12, 6, 8, 8, 16, 6, 18, 8,
		12, 10, 22, 8, 20, 12, 18, 12, 28, 8,
		30, 16, 20, 16, 24, 12, 36, 18, 24, 16,
		40, 12, 42, 20, 24, 22, 46, 16, 42, 20,
		32, 24, 52, 18, 40, 24, 36, 28, 58, 16,
		60, 30, 36, 32, 48, 20, 66, 32, 44, 24,
		70, 24, 72, 36, 40, 36, 60, 24, 78, 32,
		54, 40, 82, 24, 64, 42, 56, 40, 88, 24,
		72, 44, 60, 46, 72, 32, 96, 42, 60, 40,
	}
	for ii := 1; ii < len(testTotientPhi); ii++ {
		// for ii := 1; ii < 20; ii++ {
		phi := euler.EulerTotientPhi_old(uint64(ii))
		if phi != testTotientPhi[ii] {
			t.Errorf("Expected results: EulerTotientPhi(%d) => %d got %d\n", ii, testTotientPhi[ii], phi)
		}
	}
	for ii := 1; ii < len(testTotientPhi); ii++ {
		// for ii := 1; ii < 20; ii++ {
		phi := euler.EulerTotientPhi(uint64(ii), 0)
		if phi != testTotientPhi[ii] {
			t.Errorf("Expected results: EulerTotientPhi_exp(%d) => %d got %d\n", ii, testTotientPhi[ii], phi)
		}
	}
}

func TestGeneralMaths(t *testing.T) {
	testLeadingZeros := []struct {
		in uint64
		r  int
	}{
		{1, 63},
		{0, 64},
		{0x8000_0000_0000_0000, 0},
		{0x8000_0000, 32},
		{0x8000, 48},
		{0x80, 56},
		{0x8, 60},
	}
	for _, test := range testLeadingZeros {
		r := euler.BitsLeadingZeros64(test.in)
		if test.r != r {
			t.Errorf("Expected results: BitsLeadingZeros64(%d) = %d got %d\n", test.in, test.r, r)
		}
	}

	testPowInt := []struct {
		n, p, r uint64
	}{
		{64, 0, 1},
		{64, 1, 64},
		{2, 64, 0},
		{2, 63, 0x8000_0000_0000_0000},
		{2, 2, 4},
		{11, 18, 0x4D28CB56C33FA539},
		{17, 10, 0x1D56299ADA1},
		{3, 29, 0x3E6B41437D93},
	}
	for _, test := range testPowInt {
		r := euler.PowInt(test.n, test.p)
		if test.r != r {
			t.Errorf("Expected results: PowInt(%d, %d) = %d got %d\n", test.n, test.p, test.r, r)
		}
	}

	testSqrtI := []struct {
		n, sq uint64
	}{
		{0, 0},
		{1, 1},
		{4, 2},
		{25, 5},
		{81, 9},
		{1_000_000, 1000},
		{2, 1},
		{5, 2},
		{26, 5},
		{82, 9},
		{1_000_001, 1000},
	}
	for _, test := range testSqrtI {
		res := euler.SqrtU64(test.n)
		if test.sq != res {
			t.Errorf("Expected results: SqrtI64(%d) => %d got %d\n", test.n, test.sq, res)
		}
	}
	for _, test := range testSqrtI {
		res := euler.RootU64(test.n, 2)
		if test.sq != res {
			t.Errorf("Expected results: RootU64(%d, 2) => %d got %d\n", test.n, test.sq, res)
		}
	}
	for _, test := range testSqrtI {
		res := uint64(euler.RootF64(float64(test.n), 2, 16))
		if test.sq != res {
			t.Errorf("Expected results: RootF64(%d, 2, 16) => %d got %d\n", test.n, test.sq, res)
		}
	}
	testRootI := []struct {
		n, root, res uint32
	}{

		{8, 3, 2},
		{27, 3, 3},
		{28, 3, 3},
		{256, 8, 2},
		{257, 8, 2},
	}

	for _, test := range testRootI {
		res := uint32(euler.RootU64(uint64(test.n), uint64(test.root)))
		if test.res != res {
			t.Errorf("Expected results: RootU64(%d, %d) => %d got %d\n", test.n, test.root, test.res, res)
		}
	}
	for _, test := range testRootI {
		res := uint32(euler.RootI64(int64(test.n), test.root, 32)) // NOTE: Precision matters greatly for higher roots and for larger numbers, in that order!
		if test.res != res {
			t.Errorf("Expected results: RootI64(%d, %d) => %d got %d\n", test.n, test.root, test.res, res)
		}
	}

	for ii := 0; ii <= 512; ii++ {
		if euler.PSRand.RandU32() == euler.PSRand.RandU32() {
			t.Logf("This should only VERY rarely happen, got the same random value twice in a row.")
		}
	}

	testExtdGCD := []struct {
		a, b, s, t, gcd int64
	}{
		{3, 5, 3, 5, 1},
		{-15, -45, -1, -3, 15},
		{25, -45, 5, -9, 5},
		{-25, 45, -5, 9, 5},
		{-45, 25, -9, 5, 5},
		{45, -25, 9, -5, 5},
		{0, -45, 0, -1, 45},
		{45, 0, 1, 0, 45},
		// { , , , , },
	}
	for _, test := range testExtdGCD {
		s, tN, gcd := euler.ExtendedGCDI64(test.a, test.b)
		if test.s != s || test.t != tN || test.gcd != gcd {
			t.Errorf("Expected results: ExtendedGCDI64(%d, %d) => %d, %d, %d got %d, %d, %d\n", test.a, test.b, test.s, test.t, test.gcd, s, tN, gcd)
		}
	}
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; go clean -testcache ; go test -v euler/
/*/
