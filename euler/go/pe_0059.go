// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0, though Euler's problems are under another NC version of the license https://creativecommons.org/licenses/by-sa/4.0/

/*
https://projecteuler.net/copyright
https://creativecommons.org/licenses/by-nc-sa/4.0/
https://projecteuler.net/problem=59
https://projecteuler.net/minimal=59

<p>Each character on a computer is assigned a unique code and the preferred standard is ASCII (American Standard Code for Information Interchange). For example, uppercase A = 65, asterisk (*) = 42, and lowercase k = 107.</p>
<p>A modern encryption method is to take a text file, convert the bytes to ASCII, then XOR each byte with a given value, taken from a secret key. The advantage with the XOR function is that using the same encryption key on the cipher text, restores the plain text; for example, 65 XOR 42 = 107, then 107 XOR 42 = 65.</p>
<p>For unbreakable encryption, the key is the same length as the plain text message, and the key is made up of random bytes. The user would keep the encrypted message and the encryption key in different locations, and without both "halves", it is impossible to decrypt the message.</p>
<p>Unfortunately, this method is impractical for most users, so the modified method is to use a password as a key. If the password is shorter than the message, which is likely, the key is repeated cyclically throughout the message. The balance for this method is using a sufficiently long password key for security, but short enough to be memorable.</p>
<p>Your task has been made easy, as the encryption key consists of three lower case characters. Using <a href="resources/documents/0059_cipher.txt">0059_cipher.txt</a> (right click and 'Save Link/Target As...'), a file containing the encrypted ASCII codes, and the knowledge that the plain text must contain common English words, decrypt the message and find the sum of the ASCII values in the original text.</p>


*/
/*

3 lower case characters (implying a-z ?)
'encrypted ascii codes' you say...

Sigh... they did it dirty.  This would be somewhat easy in C where E.G. atoi(...) will assume base 10 unless told otherwise.  I dislike the commas more.

'36,22,80,0,0,4,23,25,19,17,88'

XOR might make it hard for humans to read, but...
0x61 = a = 0b_0110?_????
0x7A = z = 0b_0111?_????

Capitol Letters V have that bit cleared, which means they _should_ have 0b_000 as their highest bits, while lowercase would be 0b_001

0x41 = a = 0b_0100?_????
0x5A = z = 0b_0101?_????

Additionally we've been told the 'key' length is 3 so it's possible to divide and perform a statistical comparison of the characters in each bin...

Does this only encipher A-Za-z with the key or is it everything?

0x20 =' '= 0b_0010_0000 which XORed with a lower case letter is 0b_011?_???? and recovers the letter.




*/

import (
	"bufio"
	"euler"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	"os" // os.Stdout
	"strconv"
	"strings"
)

func Euler059(fn string) uint {
	// wc 0059_cipher.txt = 4062 /3 = 1354 but round up a bit
	crypt := make([]byte, 0, 2048)
	fh, err := os.Open(fn)
	if nil != err {
		panic("Euler0059 unable to open: " + fn)
	}
	defer fh.Close()
	var pos uint
	scanner := bufio.NewScanner(fh)
	scanner.Split(euler.ScannerSplitNLDQ)
	for scanner.Scan() {
		dec := strings.ToUpper(scanner.Text())
		if dec != "," {
			tmp, err := strconv.ParseInt(dec, 10, 0)
			b := byte(tmp)
			if nil != err {
				fmt.Printf("%d : %v\n", pos, err)
			}
			if 0b011 == (b >> 5) {
				fmt.Printf("SpaceLeaked: %d = %c\n", pos%3, b)
			}
			crypt = append(crypt, b)
			pos++
		} else {
			fmt.Printf("WARNING: encountered unknown\n\n%s\n\n\n", dec)
		}
	}
	lc := len(crypt)
	fmt.Printf("Read %d base 10 encoded XORed values (/3 = %d) - While a cryptologist probably uses statistical comparison to guess likely codes... brute force it given the small keyspace.\n", lc, lc/3)

	isValidChar := func(c byte) bool {
		// return (' ' <= c && c <= '~') || c == '\r' || c == '\n' || c == '\t' || c == '\v'
		return c == '\r' || c == '\n' || c == '\t' || c == '\v' || c == ' ' || ('0' <= c && c <= '9') || ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
	}

	key := [3]byte{}
	for l := 0; l < 3; l++ {
		var b, bestB, bestErr, ec byte
		bestErr = 0x7f
		for b = 'a'; b <= 'z'; b++ {
			ec = 0
			for ii := l; ii < lc; ii += 3 {
				dc := crypt[ii] ^ b
				if false == isValidChar(dc) {
					ec++
				}
			}
			if ec < bestErr {
				bestB = b
				bestErr = ec
			}
		}
		key[l] = bestB
	}
	fmt.Printf("Recovered likely key: %c%c%c\n", key[0], key[1], key[2])
	// "decrypt the message and find the sum of the ASCII values in the original text"
	pos = 0
	for l := 0; l < 3; l++ {
		b := key[l]
		for ii := l; ii < lc; ii += 3 {
			crypt[ii] ^= b
			pos += uint(crypt[ii])
		}
	}
	fmt.Printf("\n%s\n\n", string(crypt))

	// for ii := 0; ii < lc; ii++ {
	// crypt[ii] ^= xored[ii%3]
	// }
	return pos
}

/*
	for ii in *\/*.go ; do go fmt "$ii" ; done ; for ii in 59 ; do go fmt $(printf "pe_%04d.go" "$ii") ; go run $(printf "pe_%04d.go" "$ii") || break ; done

Read 1455 base 10 encoded XORed values (/3 = 485) - While a cryptologist probably uses statistical comparison to guess likely codes... brute force it given the small keyspace.
Recovered likely key: exp

An extract taken from the introduction of one of Euler's most celebrated papers, "De summis serierum reciprocarum" [On the sums of series of reciprocals]: I have recently found, quite unexpectedly, an elegant expression for the entire sum of this series 1 + 1/4 + 1/9 + 1/16 + etc., which depends on the quadrature of the circle, so that if the true sum of this series is obtained, from it at once the quadrature of the circle follows. Namely, I have found that the sum of this series is a sixth part of the square of the perimeter of the circle whose diameter is 1; or by putting the sum of this series equal to s, it has the ratio sqrt(6) multiplied by s to 1 of the perimeter to the diameter. I will soon show that the sum of this series to be approximately 1.644934066842264364; and from multiplying this number by six, and then taking the square root, the number 3.141592653589793238 is indeed produced, which expresses the perimeter of a circle whose diameter is 1. Following again the same steps by which I had arrived at this sum, I have discovered that the sum of the series 1 + 1/16 + 1/81 + 1/256 + 1/625 + etc. also depends on the quadrature of the circle. Namely, the sum of this multiplied by 90 gives the biquadrate (fourth power) of the circumference of the perimeter of a circle whose diameter is 1. And by similar reasoning I have likewise been able to determine the sums of the subsequent series in which the exponents are even numbers.

Euler 59: XOR Decryption: 129448
*/
func main() {
	//test

	//run
	fmt.Printf("Euler 59: XOR Decryption: %d\n", Euler059("0059_cipher.txt"))
}
