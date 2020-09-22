package test

import (
	"bytes"
	"comment"
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	src = []string{
		"unicode/8",
		`
		// Copyright 2009 The Go Authors. All rights reserved.
		// Use of this source code is governed by a BSD-style
		// license that can be found in the LICENSE file.
		`,
		"/*\n *ssss\n *ssss\n */",
		"/*ssss*/\n//sss\n",
		"//ssss\n\n",
		"//ssss",
		"/******/",
		"/*ssss",
		"  /*ssss/*ssss*/ssss*/",
		"    /*  ssss*  */ ",
		"/*ssss*/",
		"/*  ssss*  */ ",
	}
	results = [][]string{
		{},
		{
			" Copyright 2009 The Go Authors. All rights reserved.",
			" Use of this source code is governed by a BSD-style",
			" license that can be found in the LICENSE file.",
		},
		{
			"\n *ssss\n *ssss\n ",
		},
		{
			"ssss", "sss",
		},
		{
			"ssss",
		},
		{
			"ssss",
		},
		{"****"},
		{},
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
	}
)

func TestGetTokenFromBytes(t *testing.T) {
	var pass = true
	for k, v := range src {
		res := comment.GetTokenFromBytes([]byte(v))
		should := results[k]
		match := true
		for k, v := range res {
			if v.Content != should[k] {
				match = false
				pass = false
			}
		}
		for _,v:= range res{
			fmt.Println(comment.PurifyToken(&v,""))
		}
		fmt.Println(match)
	}
	if !pass {
		t.Error("unmatched")
	}
}

func TestGetTokenFromBytesFile(t *testing.T) {
	var LoadFile = func(input, output string) {
		var b, err = ioutil.ReadFile(input)
		fmt.Println("read from ", input)
		if err != nil {
			t.Error("err:", err)
			return
		}
		res := comment.GetTokenFromBytes(b)
		var buf bytes.Buffer
		for _, v := range res {
			fmt.Println(comment.PurifyToken(&v,input))
			buf.WriteString(fmt.Sprintln(comment.PurifyToken(&v,input)))
		}
		if err := ioutil.WriteFile(output, buf.Bytes(), 0x777); err != nil {
			t.Error("err:", err)
		}
	}

	LoadFile("./example.go.txt", "./example_go_out")
	LoadFile("./example.java.txt", "./example_java_out")
	LoadFile("./example.cpp.txt", "./example_cpp_out")
}