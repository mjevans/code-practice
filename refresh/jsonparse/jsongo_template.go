// kate: space-indent off; indent-width 8; tab-width 8; mixedindent off; indent-mode tab;
package main

// 2024 - Michael J Evans
// Code in this file is CC BY-SA 4.0 license https://creativecommons.org/licenses/by-sa/4.0/

/*

go fmt *.go ; go run .go

*/

import (
	// "bufio"
	// "bitvector"
	// "euler"
	// "encoding/json"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

func TestJSONs(ajson, kreq []string) (pass, fail, errdecode int) {
	return
}

func main() {
	//test
	key_req := []string{"id", "name"}
	// pass, pass, errdecode, fail, fail, errdecode?
	// 2 3 2
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
	fmt.Println("\njson decode tests:\t", 2 == p && 3 == f && 2 == e, "\t results:\t", p, f, e)
	if 2 == p && 3 == f && 2 == e {
		return
	}
	panic("Failed JSON test.")
}
