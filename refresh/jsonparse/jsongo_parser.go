// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// golang 1.19 is current Debian stable
// 2024 - Michael J Evans ***REMOVED***

/*

go fmt *.go ; go run .go

*/

import (
	// "bufio"
	// "bitvector"
	// "euler"
	// "encoding/json"
	"errors"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	"strings"
	// "strconv"
	// "os" // os.Stdout
)

// go 1.18+ generic array of type T
type Stack[T any] []T // []interface{T}

func (s *Stack[T]) Push(v T) *Stack[T] {
	*s = append(*s, v)
	return s
}

func (s *Stack[T]) Peek() T {
	var ret T
	end := len(*s) - 1
	if 0 > end {
		return ret // default value
	}
	ret = (*s)[end]
	return ret
}

func (s *Stack[T]) Pop() T {
	var ret T
	end := len(*s) - 1
	if 0 > end {
		return ret // default value
	}
	ret = (*s)[end]
	*s = (*s)[:end]
	return ret
}

// map... is syntax sugar for (*map)... and indirects to functions in a core library.
// Slice is view into an array, already a pointer, but the view needs manual assignment.

type JSONel struct {
	Dict map[string]*JSONel // String indexed map/dict
	Arr  []*JSONel          // numeric indexed array
	Raw  string             // raw value of element (populated if a value)
}

func (j *JSONel) AddChildKV(k string, v *JSONel) {
	if nil == j.Dict {
		j.Dict = make(map[string]*JSONel, 0)
	}
	j.Dict[k] = v
}

func (j *JSONel) AddChildV(v *JSONel) {
	if nil == j.Arr {
		j.Arr = make([]*JSONel, 0)
	}
	j.Arr = append(j.Arr, v)
}

func (j JSONel) String() string {
	var ret strings.Builder
	if 0 < len(j.Dict) {
		ret.WriteString("{")
		limII := len(j.Dict)
		ii := 0 // track if we're near the end yet...
		for k, v := range j.Dict {
			ret.WriteString(k)
			ret.WriteString(":")
			ret.WriteString(v.String())
			if ii+1 < limII {
				ret.WriteString(",")
			}
			ii++
		}
		ret.WriteString("}")
		return ret.String()
	}
	if 0 < len(j.Arr) {
		ret.WriteString("[")
		limII := len(j.Arr)
		for ii := 0; ii < limII; ii++ {
			ret.WriteString(j.Arr[ii].String())
			if ii+1 < limII {
				ret.WriteString(",")
			}
		}
		ret.WriteString("]")
		return ret.String()
	}
	return j.Raw
}

// (lowercase only): null false true (floatnum) "string" { [
const (
	TypeDict = iota
	TypeArray
	TypeString
	TypeFloat
	TypeInt
	TypeBool
	TypeNil
)

func (j JSONel) GuessType() int {
	if 0 < len(j.Dict) {
		return TypeDict
	}
	if 0 < len(j.Arr) {
		return TypeArray
	}
	// Non-strings can always be used as strings...
	if "" == j.Raw {
		return TypeString
	}
	if "null" == j.Raw {
		return TypeNil
	}
	if "flase" == j.Raw || "true" == j.Raw {
		return TypeBool
	}
	// (guarded len at least 1), Might be INT , or FLOAT , or String
	var ii int
	limII := len(j.Raw)
	if '-' == j.Raw[0] {
		ii++
	}
	if 1 < limII-ii && '0' == j.Raw[ii] && '.' == j.Raw[ii+1] {
		ii += 2
	} else if 1 == limII-ii && '0' == j.Raw[ii] {
		return TypeInt
	} else if 1 <= limII-ii && 1 <= j.Raw[ii]-'0' && 9 >= j.Raw[ii]-'0' {
		ii++ // Must start with a number that is not zero, unless it is only zero or a float between -1 and 1 (exclusive)
	} else {
		return TypeString
	}
	// is int?
	for ii < limII && 0 <= j.Raw[ii]-'0' && 9 >= j.Raw[ii]-'0' {
		ii++
	}
	if ii == limII {
		return TypeInt
	}
	if 2 <= limII-ii && '.' == j.Raw[ii] && 0 <= j.Raw[ii+1]-'0' && 9 >= j.Raw[ii+1]-'0' {
		ii += 2
		for ii < limII && 0 <= j.Raw[ii]-'0' && 9 >= j.Raw[ii]-'0' {
			ii++
		}
		if ii == limII {
			return TypeFloat
		}
	}
	if 'e' == j.Raw[ii] || 'E' == j.Raw[ii] {
		ii++
		if 1 <= limII-ii && ('-' == j.Raw[ii] || '+' == j.Raw[ii]) {
			ii++
		}
		if 1 <= limII-ii && 1 <= j.Raw[ii]-'0' && 9 >= j.Raw[ii]-'0' {
			for ii < limII && 0 <= j.Raw[ii]-'0' && 9 >= j.Raw[ii]-'0' {
				ii++
			}
		}
		if 1 <= limII-ii && '0' == j.Raw[ii] {
			ii++
		}
		if ii == limII {
			return TypeFloat
		}
	}
	return TypeString
}

func isSpace(s byte) bool {
	return ' ' == s || '\t' == s || '\r' == s || '\n' == s || '\v' == s
}

func isHexByte(b byte) bool {
	return ('0' <= b && '9' >= b) || ('A' <= b && 'F' >= b) || ('a' <= b && 'f' >= b)
}

func byteADecodeHex(b []byte) (uint64, error) {
	var ret, nib uint64
	limII := len(b)
	for ii := 0; ii < limII; ii++ {
		if '0' <= b[ii] && '9' >= b[ii] {
			nib = uint64(b[ii] - '0')
		} else if 'A' <= b[ii] && 'F' >= b[ii] {
			nib = uint64(b[ii] - 'A')
		} else if 'a' <= b[ii] && 'f' >= b[ii] {
			nib = uint64(b[ii] - 'a')
		} else {
			return uint64(0), InvalidUTF16Error
		}
		ret = (ret << 4) | nib
	}
	return ret, nil
}

// type WriterCloser interface {
// 	Write(p []byte) (n int, err error)
// 	Close() error
// }
//func JSONNewDecoder(r io.Reader) *Decoder{
//	Internal function interface should be like the
//}

// https://pkg.go.dev/errors@go1.23.0#New

// All of the exampls I found elsewhere focus on
// interface { Unwrap() error }
// However, "errors" also supports
// interface { Unwrap() []error }

type ErrArray struct {
	err []error
}

func NewErrArray(err ...error) *ErrArray {
	return &ErrArray{err}
}

func (e *ErrArray) Unwrap() []error {
	return (*e).err
}

func (e *ErrArray) Error() string {
	var tmp string
	limII := len((*e).err)
	if 0 >= limII {
		return ""
	}
	// inefficient, but this should be both rare and not pathalogical...
	tmp = ((*e).err)[0].Error()
	for ii := 1; ii < limII; ii++ {
		tmp += ": " + ((*e).err)[ii].Error()
	}
	return tmp
}

// var _ interface{...} = (*TypeStruct)(nil) // *TypeStruct implements the namedInterface


type JSONDecoder struct {
	stB   Stack[byte]
	stJ   Stack[*JSONel]
	read  int
	flags int
	key   string
}

const (
	AllowSingleValue = 1 << iota
	IgnoreRepeatedComma
	AllowUnquotedString
)

func jsonNewDecoderFull(sz, flags int) *JSONDecoder {
	ret := new(JSONDecoder)
	ret.read = 0
	ret.key = "" // These should be the defaults, more to document
	ret.stB = append(make(Stack[byte], 0, sz), byte(0), byte(0))
	tmp := &JSONel{}
	ret.stJ = append(make(Stack[*JSONel], 0, sz), tmp, tmp)
	ret.flags = flags
	return ret
}

func JSONNewDecoder(flags int) *JSONDecoder {
	return jsonNewDecoderFull(32, flags)
	// return jsonNewDecoderFull(int(^(uint(0)>>1)), 32, flags)
}

var (
	JSONIncomplete               = errors.New("JSON Incomplete")
	JSONError                    = errors.New("JSON Error")
	SingleValueError             = errors.New("Expected Initial Object { or Array [")
	ExpectedKeyError             = errors.New("Expected Quoted Key")
	ExpectedSemicolonError       = errors.New("Expected Semicolon")
	ExpectedValueError           = errors.New("Expected Value")
	ExpectedCloseContinueError   = errors.New("Expected a , to continue or } or ] to close")
	CommaError                   = errors.New("Multiple commas without value between")
	InvalidUTF16Error            = errors.New("Invalid UTF-16, valid escape values 0x0000 : 0xD7FF , 0x9000 : 0xFFFF , or a surrogate pair with the high bits first")
	UnsupportedEscapeSequence    = errors.New("Unsupported JSON escape sequence")
	UnsupportedCharacterSequence = errors.New("Unsupported JSON character sequence")
)

func (j *JSONDecoder) NewJSONErrS(e error, ii int) *ErrArray {
	return NewErrArray(JSONError, e, fmt.Errorf("%d bytes processed ~ %d bytes in latest call", j.read+ii, ii))
}

func JSONDecodeString(raw []byte) (string, error) {
	// Uggh Thanks JavaScript / ECMAScript UTF-16 / UCS-2 further immortalized forever...
	// https://www.unicode.org/versions/Unicode15.0.0/ch03.pdf
	// Table 3-5 (file page 55) Conformance 124 3.9 Unicode Encoding Forms
	// UCS-2 past 0xFFFF extended to cover 0x00010000 through 0x0010FFFF by reserving 0xD800 through 0xDFFF , and worse using all of it to smuggle only up to 20 bits + the base 0x10000 ... at the cost of every OTHER unicode encoding excluding the 0x800 possible glyphs.
	// It even infects UTF-8 with the madness; as table 3-6 and 3-7 show on the following pages.
	// 0x40000 is where ease could return to the UTF-8 style process (if uncapped through 111111 10B / 0xFD)
	// Quick aside: If 'exceptions' are going to be added, it would have been WAY cleaner to for Chr >= 0x10000 to just use 0xF as a prefix, and claim the 0x0F nibble of that first byte as 3 + nibble 6 bit byte blocks.
	// 6 bits * 6 octets = 36 bits (at a storage cost of 7 bytes) - which would easily cover an entire 32 bit 'full' Unicode character, and likely satisfy as many alien and fictional languages as we care to invent before civilization ends or we cease caring about the cost of storage.
	// Though even more practically speaking, it's obvious why the chart stops at 4 bytes (though aside from supporting UTF-16's uglyness)
	// Though not not why UTF-8 can't use up to the 0x1FFFFF it could so justly store past UTF-16's madness to cram in support for slightly less than half those glyphs.
	// Upper, Lower = ((u32 >> 16) - 1) | ((u32 >> 10) & 0x3f) | 0xD800  ,  (u32 & 0x3ff) | 0xDC00
	limII := len(raw)
	sz := limII
JSONDecodeString_outer:
	for ii := 0; ii < limII; ii++ {
		// https://datatracker.ietf.org/doc/html/rfc8259#section-7
		// valid unescaped %x20-21 / %x23-5B / %x5D-10FFFF
		if '\x20' == raw[ii] || '\x21' == raw[ii] || ('\x23' <= raw[ii] && '\x5B' >= raw[ii]) || ('\x5D' <= raw[ii]) { // golang's too serious to get the joke: && '\U0010FFFF' >= raw[ii]
			continue JSONDecodeString_outer
		}
		if ii+2 > limII || '\\' != raw[ii] {
			return "", UnsupportedCharacterSequence
		}
		ii++ // Advance over \
		switch raw[ii] {
		case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			sz--
			continue JSONDecodeString_outer
		case 'u':
			ii++ // Advance over u
			if ii+4 > limII {
				return "", UnsupportedEscapeSequence
			}
			var err error
			var su, sl int
			sur, err := byteADecodeHex(raw[ii : ii+4])
			su = int(sur)
			if nil != err {
				return "", err
			}
			ii += 3 // Advance to final Hex of normal UTF-16 escape
			// Validate UTF-16...
			if su < 0x80 {
				sz -= 6 - 1
				continue JSONDecodeString_outer
			}
			if su < 0x800 {
				sz -= 6 - 2
				continue JSONDecodeString_outer
			}
			// if su < 0x10000 { sz -= 6 - 3 ; continue }
			// if su < 0xD800 || ( su >= 0x9000 && su < 0x10000) {
			if su < 0xD800 || su >= 0x9000 { // 0xFFFF is the max, since input 4
				sz -= 6 - 3 // 6 bytes becomes 4 bytes
				continue JSONDecodeString_outer
			}
			// https://www.unicode.org/versions/Unicode15.0.0/ch03.pdf
			// Surrogate Pair check, leading pair?
			if 0xD800 != (su&0xFC00) || ii+7 > limII || '\\' != raw[ii+1] || 'u' != raw[ii+2] {
				return "", InvalidUTF16Error
			}
			ii += 3 // Advance past \u
			sur, err = byteADecodeHex(raw[ii : ii+4])
			sl = int(sur)
			if nil != err {
				return "", err
			}
			if 0xDC00 != (sl & 0xFC00) {
				return "", InvalidUTF16Error
			}
			ii += 3 // Advance to final Hex of Surrogate Pair
			sz -= 12 - 4
			continue JSONDecodeString_outer
		}

	}
	// sz now contains the size to process...
	buf := make([]byte, sz)
	bb := 0
JSONDecodeString_copy:
	for ii := 0; ii < limII; ii++ {
		// Copy over all valid characters...
		if '\x20' == raw[ii] || '\x21' == raw[ii] || ('\x23' <= raw[ii] && '\x5B' >= raw[ii]) || ('\x5D' <= raw[ii]) {
			buf[bb] = raw[ii]
			bb++
			continue JSONDecodeString_copy
		}
		// All Escapes were validated during result length calculation
		ii++ // Advance over \
		switch raw[ii] {
		// , '\\', '/', 'b', 'f', 'n', 'l', 't':
		case '"':
			buf[bb] = '"'
		case '\\':
			buf[bb] = '\\'
		case '/':
			buf[bb] = '/'
		case 'b':
			buf[bb] = '\b'
		case 'f':
			buf[bb] = '\f'
		case 'n':
			buf[bb] = '\n'
		case 'r':
			buf[bb] = '\r'
		case 't':
			buf[bb] = '\t'
		case 'u':
			ii++ // Advance over u
			var su, sl int
			sur, _ := byteADecodeHex(raw[ii : ii+4])
			su = int(sur)
			ii += 3 // Advance to final Hex of normal UTF-16 escape
			// Validate UTF-16...
			if su < 0x80 {
				buf[bb] = byte(su)
				bb += 1
				continue JSONDecodeString_copy
			}
			if su < 0x800 {
				buf[bb+0] = byte(0xC0 | ((su >> 6) & 0x1F))
				buf[bb+1] = byte(0x80 | (su & 0x3F))
				bb += 2
				continue JSONDecodeString_copy
			}
			// if su < 0x10000 { sz -= 6 - 3 ; continue }
			// if su < 0xD800 || ( su >= 0x9000 && su < 0x10000) {
			if su < 0xD800 || su >= 0x9000 { // 0xFFFF is the max, since input 4
				buf[bb+0] = byte(0xE0 | ((su >> 12) & 0x0F))
				buf[bb+1] = byte(0x80 | ((su >> 6) & 0x3F))
				buf[bb+2] = byte(0x80 | (su & 0x3F))
				bb += 3
				continue JSONDecodeString_copy
			}
			// https://www.unicode.org/versions/Unicode15.0.0/ch03.pdf
			ii += 3 // Advance past \u
			sur, _ = byteADecodeHex(raw[ii : ii+4])
			sl = int(sur)
			su = ((su & 0x3FF) << 10) | (sl & 0x3FF)
			buf[bb+0] = byte(0xF0 | ((su >> 18) & 0x07))
			buf[bb+1] = byte(0x80 | ((su >> 12) & 0x3F))
			buf[bb+2] = byte(0x80 | ((su >> 6) & 0x3F))
			buf[bb+3] = byte(0x80 | (su & 0x3F))
			bb += 4
			ii += 3 // Advance to final Hex of Surrogate Pair
			continue JSONDecodeString_copy
		}
		bb++
		continue JSONDecodeString_copy
	}
	return string(buf), nil
}

func (j *JSONDecoder) Write(raw []byte) (int, error) {
	parent := j.stJ.Peek()
	state := j.stB.Peek()
	kvStart := -1
	lastGood := -1
	// singleValue := false
	var err error
	unquotedString := false
	limII := len(raw)
	var ii int
JSONDecoder_Write_ii:
	for ii = 0; ii < limII; ii++ {

		// skip 'free' space
		for ii < limII && isSpace(raw[ii]) {
			ii++
		}
		if ii >= limII {
			break JSONDecoder_Write_ii
		}
		if ',' == raw[ii] {
			if 0 < j.flags&IgnoreRepeatedComma {
				continue JSONDecoder_Write_ii
			}
			return 0, j.NewJSONErrS(CommaError, ii)
		}

		lastGood = ii
		switch state {
		case '\x00': // Expect / Accept only object or array (strict json)

			j.stB.Pop() // Replace state

			switch raw[ii] {
			case '{', '[':
				state = raw[ii]
				j.stB.Push(raw[ii])
				continue JSONDecoder_Write_ii
			default:
				if 0 < j.flags&AllowSingleValue {
					state = '['
					j.stB.Push('[')
					// singleValue = true
					continue JSONDecoder_Write_ii
				}
				// j.read += ii - 1
				return 0, j.NewJSONErrS(SingleValueError, ii)
			}
		case '{':
			// Require Key then value
			if "" == j.key {
				if '"' != raw[ii] {
					if 0 < j.flags&AllowUnquotedString {
						unquotedString = true
					} else {
						// j.read += ii - 1
						return 0, j.NewJSONErrS(ExpectedKeyError, ii)
					}
				} else {
					ii += 1
				}
				kvStart = ii
				if false == unquotedString {
					for ii < limII && ('"' != raw[ii] || '\\' == raw[ii-1]) {
						ii++
					}
					if ii >= limII {
						ii = lastGood
						break JSONDecoder_Write_ii
					}
					j.key, err = JSONDecodeString(raw[kvStart:ii])
					if nil != err {
						return 0, err
					}
					ii++ // "

					// skip 'free' space
					for ii < limII && isSpace(raw[ii]) {
						ii++
					}

					// MUST have : to promote to key
					if ':' != raw[ii] || ii+1 >= limII {
						fmt.Printf("%c", raw[ii])
						return 0, j.NewJSONErrS(ExpectedSemicolonError, ii)
					}
					ii += 1

					// skip 'free' space
					for ii < limII && isSpace(raw[ii]) {
						ii++
					}
					if ii >= limII {
						ii = lastGood
						break JSONDecoder_Write_ii
					}
				} else {
					kvEnd := kvStart
					for ii < limII && ':' != raw[ii] {
						if false == isSpace(raw[ii]) {
							kvEnd = ii
						}
						ii++
					}
					if ii >= limII {
						ii = lastGood
						break JSONDecoder_Write_ii
					}
					j.key, err = JSONDecodeString(raw[kvStart:kvEnd])
					if nil != err {
						return 0, err
					}
				}
				if ii+1 >= limII {
					ii = limII
					break JSONDecoder_Write_ii
				}
				kvStart = -1
				lastGood = ii
			}
			// Read Value
			fallthrough
		case '[':
			child := &JSONel{} // NewJSONel()
			// Expect: https://datatracker.ietf.org/doc/html/rfc8259#section-3
			// (lowercase only): null false true (floatnum) "string" { [
			switch raw[ii] {
			case '"': // Quoted String
				ii += 1
				kvStart = ii
				for ii < limII && ('"' != raw[ii] || '\\' == raw[ii-1]) {
					ii++
				}
				if ii >= limII {
					ii = lastGood
					break JSONDecoder_Write_ii
				}
				child.Raw, err = JSONDecodeString(raw[kvStart:ii])
				if nil != err {
					return 0, err
				}
				ii++ // "
			default:
				// find the next control character: , } ]
				for kvStart = ii; ii < limII && ',' != raw[ii] && '}' != raw[ii] && ']' != raw[ii]; ii++ {
				}
				if ii >= limII {
					ii = lastGood
					break JSONDecoder_Write_ii
				}
				for kvStart < ii && isSpace(raw[ii]) {
					ii--
				}
				if "null" == string(raw[kvStart:ii]) || "false" == string(raw[kvStart:ii]) || "true" == string(raw[kvStart:ii]) {
					child.Raw, err = JSONDecodeString(raw[kvStart:ii])
					if nil != err {
						return 0, err
					}
				} else { // It must be a float, OR an out of spec unquoted string
					// FIXME: Optional, support unquoted srings?
					// FIXME: Validate case insensitive, base10, no leading zeros (unless 1s place)... E.G. -0.231E[+-]?12
					// https://datatracker.ietf.org/doc/html/rfc8259#section-6
					child.Raw, err = JSONDecodeString(raw[kvStart:ii])
					if nil != err {
						return 0, err
					}
				}
			case '{', '[': // Object OR Array
				state = raw[ii]
				j.stB.Push(raw[ii])
				if "" != j.key {
					// parent.Dict[j.key] = child
					parent.AddChildKV(j.key, child)
				} else {
					// parent.Arr = append(parent.Arr, child)
					parent.AddChildV(child)
				}
				j.stJ.Push(parent)
				parent = child
				continue JSONDecoder_Write_ii
			}
			kvStart = -1
			if "" != j.key {
				// parent.Dict[j.key] = child
				parent.AddChildKV(j.key, child)
				j.key = ""
			} else {
				// parent.Arr = append(parent.Arr, child)
				parent.AddChildV(child)
			}
		}

		// skip 'free' space
		for ii < limII && isSpace(raw[ii]) {
			ii++
		}
		if ii >= limII {
			ii = limII
			break JSONDecoder_Write_ii
		}

		// NEXT control code , ] } // Just read a key:value or value
		switch raw[ii] {
		case ',':
			// Continue as is
		case '}', ']':
			// close object or list
			state = j.stB.Pop()
			parent = j.stJ.Pop()
		default:
			// NOT a control character
			return 0, j.NewJSONErrS(ExpectedCloseContinueError, ii)
		}
	}
	j.read += ii
	return ii, nil
}

func (j *JSONDecoder) Close() error {
	lb, lj := len(j.stB), len(j.stJ)
	if lb != lj || 0 == j.read {
		return JSONError
	}
	// fully closed [ or {
	if lb == 0 {
		return nil
	}
	// 'single value' [ mode
	if lb == 1 && '\000' == j.stB[0] {
		return nil
	}
	// Otherwise error
	return JSONError
}

func (j *JSONDecoder) Root() (*JSONel, error) {
	if 1 == len(j.stJ) || (2 == len(j.stJ) && 2 == len(j.stB) && '\000' == j.stB[0]) {
		return j.stJ[0], nil
	}
	return nil, JSONError
}

func JSONUnmarshal(raw string) (*JSONel, error) {
	j := JSONNewDecoder(0)
	if nil == j {
		return nil, JSONError
	}

	_, err := j.Write([]byte(raw))
	if nil != err {
		return nil, err
	}

	err = j.Close()
	if nil != err {
		return nil, err
	}

	return j.Root()

}

func TestJSONs(ajson, kreq []string) (pass, fail, errdecode int) {
	// Write a parser, use in place since a string array is handed to the func already... for real return deep copies

	// JSON OUTER: { object, [ array
	// JSON Inner: quoted string, boolean (any case), 'float(unquoted text)', more objects, 'null' as a null object / value.
	// Whitespace " \t\r\n\v"
	// No trailing commas

	// Yeah, I was doomed the moment I tried to write a JSON parser in a timed test.

	// State	stack	Expect
	// (nil)	(nil)	{ [
	// [		[ ,	('value') :: null, true / false, float, "string", { , [
	// {		{ ,	"key"\s*:
	// :		{ ,	('value') :: null, true / false, float, "string", { , [

	// FIXME :

	limJS := len(ajson)
	for kk := 0; kk < limJS; kk++ {
		je, err := JSONUnmarshal(ajson[kk])
		fmt.Printf("%s\n%v\n%v\n", ajson[kk], je, err)
		if nil != err {
			errdecode++
			continue
		}
		var okn, oki bool
		_, okn = je.Dict["name"]
		_, oki = je.Dict["id"] //
		if true != okn || true != oki {
			fail++
			continue
		}
		pass++
		_ = je
		_ = err
	}
	return
}

func main() {
	//test
	key_req := []string{"id", "name"}
	// pass, pass, errdecode, fail, fail, errdecode?
	// 2 3 2 -- pass pass err p? f? f f err
	test_json := []string{
		"{\"id\": \"200\", \"name\": \"Test User\"}",
		"\r\n  \t{\"id\": \"200\", \"name\": \"Test User\", \"none\": \"Test User\"}\t\r\n   \n",
		"{id: 200, name: \"Test User\"}",
		"{\"id\": 200, \"name\": \"Test User\"}",
		"{\"id\": \"200\", \"none\": \"Test User\"}",
		"{\"did\": \"200\", \"name\": \"Test User\"}",
		"\"id\": \"200\", \"name\": \"Test User\"",
	}

	p, f, e := TestJSONs(test_json, key_req)
	fmt.Println("\njson decode tests, alternate outcome from encoding/json :\t", 3 == p && 2 == f && 2 == e, "\t results:\t", p, f, e)
	fmt.Println("\nFIXME: More comprehensive unit tests: UTF-16, invalid UTF-16, and tests for every relax flag.\nMaybe also JSON5 / other flavor constructors...")
	if 3 == p && 2 == f && 2 == e {
		return
	}
	panic("Failed JSON test.")
}
