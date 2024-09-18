// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package euler_test

import (
	// . "euler" // https://go.dev/wiki/CodeReviewComments#import-dot
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
