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
	"encoding/json"
	"fmt"
	// "math"
	// "math/big"
	// "slices" // Doh not in 1.19
	// "sort"
	// "strings"
	// "strconv"
	// "os" // os.Stdout
)

// string == a.(type)  // I would REALLY like to directly test the expected type instead
func anyIsString(a any) bool {
	switch a.(type) {
	case string:
		return true
	}
	return false
}

func TestJSONs(ajson, kreq []string) (pass, fail, errdecode int) {
	var jmap map[string]interface{} // map[string]any
	limit_cc := len(kreq)
	limit := len(ajson)
TestJSONs_outer:
	for ii := 0; ii < limit; ii++ {
		jmap = nil
		if err := json.Unmarshal([]byte(ajson[ii]), &jmap); err != nil {
			errdecode += 1
			fmt.Println("TestJSONs errdecode: ", ajson[ii])
			continue
		}
		for cc := 0; cc < limit_cc; cc++ {
			// var v interface{}
			// var ok bool
			v, ok := jmap[kreq[cc]]
			if ok != true || false == anyIsString(v) {
				fail += 1
				fmt.Println("TestJSONs fail: ", ajson[ii], "\tmissing: ", kreq[cc])
				continue TestJSONs_outer
			} else {
				fmt.Print(kreq[cc], ":", v.(string), "\t")
			}
			// Test for kreq here
			_ = v
		}
		pass += 1
		fmt.Println("TestJSONs pass: ", ajson[ii])
	}
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
