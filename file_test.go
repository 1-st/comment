package comment

import (
	"fmt"
	"testing"
)

/*ssss/*ssss*/

func TestGetCommentStrings(t *testing.T) {
	src := []string{
		"  /*ssss/*ssss*/ssss*/",
		"    /*  ssss*  */ ",
		"/*ssss*/",
		"/*  ssss*  */ ",
		"//ssss",
		"//ssss\n\n",
		"/*ssss*/\n//sss\n",
	}
	results := [][]string{
		{
			"ssss/*ssss",
		},
		{
			"  ssss*  ",
		},
		{
			"ssss",
		},
		{
			"  ssss*  ",
		},
		{
			"ssss",
		},
		{
			"ssss",
		},
		{
			"ssss", "sss",
		},
	}
	var pass = true
	for k, v := range src {
		res := GetCommentStrings([]byte(v))
		should := results[k]
		match := true
		for k, v := range res {
			if v != should[k] {
				match = false
				pass = false

			}
		}
		fmt.Println(res)
		fmt.Println(should)
		fmt.Println(match)
	}
	if !pass {
		t.Error("unmatched")
	}
}
