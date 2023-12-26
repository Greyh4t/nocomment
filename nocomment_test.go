package nocomment

import (
	"testing"
)

func TestClean(t *testing.T) {
	codeList := [][]string{
		{"", ""},
		{"hello world", "hello world"},
		{"//this is a comment\nHello World\n", "\nHello World\n"},
		{"//this is a comment\r\nHello World\r\n", "\r\nHello World\r\n"},
		{"//this is a comment\nHello World// another comment\n", "\nHello World\n"},
		{"//this is a comment\r\nHello World// another comment\r\n", "\r\nHello World\r\n"},
		{"#this is a comment\nHello World\n", "\nHello World\n"},
		{"#this is a comment\r\nHello World\r\n", "\r\nHello World\r\n"},
		{"#this is a comment\nHello World# another comment\n", "\nHello World\n"},
		{"#this is a comment\r\nHello World# another comment\r\n", "\r\nHello World\r\n"},
		{"--this is a comment\nHello World\n", "\nHello World\n"},
		{"--this is a comment\r\nHello World\r\n", "\r\nHello World\r\n"},
		{"--this is a comment\nHello World-- another comment\n", "\nHello World\n"},
		{"--this is a comment\r\nHello World-- another comment\r\n", "\r\nHello World\r\n"},
		{"//this is a comment\nHello World# another comment\n", "\nHello World\n"},
		{"//this is a comment\r\nHello World# another comment\r\n", "\r\nHello World\r\n"},
		{"/*this is a comment*/Hello World", "Hello World"},
		{"/*this is a comment*/\nHello World\n", "\nHello World\n"},
		{"/*this is a comment*/\r\nHello World\r\n", "\r\nHello World\r\n"},
		{"/*this\n is a\n comment\n*/Hello World/* another comment*/\n", "Hello World\n"},
		{"/*this\r\n is a\r\n comment\r\n*/Hello World/* another comment*/\r\n", "Hello World\r\n"},
		{"<!--this is a comment-->Hello World", "Hello World"},
		{"<!--this is a comment-->\nHello World\n", "\nHello World\n"},
		{"<!--this is a comment-->\r\nHello World\r\n", "\r\nHello World\r\n"},
		{"<!--this\n is a\n comment\n-->Hello World<!-- another comment-->\n", "Hello World\n"},
		{"<!--this\r\n is a\r\n comment\r\n-->Hello World<!-- another comment-->\r\n", "Hello World\r\n"},
		{`This is some text. "#This is not a comment // neither is this /* or this */" sooo, no comments!`, `This is some text. "#This is not a comment // neither is this /* or this */" sooo, no comments!`},
		{"/* this is a broken block comment", ""},
		{"\" this is an unlcosed quote", "\" this is an unlcosed quote"},
		{"select * /* this is\n a comment \n*/\nfrom table\nwhere name=\"aa--\"\nand value like '%/*this is not comment*/%'\nand id in ('14') -- this a comment\n-- this a comment", "select * \nfrom table\nwhere name=\"aa--\"\nand value like '%/*this is not comment*/%'\nand id in ('14') \n"},
	}
	s := &Stripper{
		RemoveBlockComment: true,
		RemoveLineComment:  true,
		RemoveShellComment: true,
		RemoveHtmlComment:  true,
		RemoveSQLComment:   true,
	}

	for i, v := range codeList {
		code := v[0]
		expectCode := v[1]
		r := string(s.Clean([]byte(code)))
		if expectCode != r {
			t.Errorf("id: %d, expect: %#v, actual: %#v", i, expectCode, r)
		}
	}
}
