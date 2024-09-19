// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package euler_test

import (
	// . "euler" // https://go.dev/wiki/CodeReviewComments#import-dot
	"container/heap"
	"euler" // I would __really__ like Go to just support path relative imports E.G. "./euler" means, just look in the CWD, it's here.
	"testing"
)

/*-

	https://github.com/golang/go/issues/25223

	https://stackoverflow.com/questions/19998250/proper-package-naming-for-testing-with-the-go-language/31443271#31443271

	https://pkg.go.dev/cmd/go#hdr-Test_packages

	https://go.dev/src/strings/search_test.go

	for ii in *\/*.go ; do go fmt "$ii" ; done ; go test -v euler/
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

// func TestXxx(t *testing.T) { t.  }

func TestBVPrimesPrimeAfter(t *testing.T) {
	p := euler.NewBVPrimes()
	testsPrimeAfter := []struct{ q, r uint }{
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
	testsPrimes := []uint{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131}
	for ii := 0; ii+1 < len(testsPrimes); ii++ {
		res := p.PrimeAfter(testsPrimes[ii])
		if testsPrimes[ii+1] != res {
			t.Errorf("Bad PrimeAfter(%d) %d != %d (expected)\n", testsPrimes[ii], res, testsPrimes[ii+1])
		} else {
			// t.Logf("PASS PrimeAfter(%d) == %d\n", testsPrimes[ii], res)
		}
	}
	for _, tc := range testsPrimeAfter {
		res := p.PrimeAfter(tc.q)
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
	mark := uint16(0)
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
		test uint
		ans  []uint
	}{
		{5309, []uint{5309}},
		{867, []uint{3, 17, 17}},
		{2024, []uint{2, 2, 2, 11, 23}},
		{1885, []uint{5, 13, 29}},
	}
	p := euler.NewBVPrimes()
	for _, test := range tests {
		left := p.Factorize(test.test)
		right := p.Factorize(1)
		for _, subfact := range test.ans {
			right.Mul(p.Factorize(subfact))
		}
		if false == left.Eq(right) || left.Uint64() != right.Uint64() {
			t.Errorf("Failed Test Case %v\n\t%d != %d\n%v\n%v", t, left.Uint64(), right.Uint64(), left, right)
		}
	}
}
