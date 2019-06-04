# nocomment
nocomment is a tool to remove comments of code, it support /**/, //, #, \<!-- -->

## Usage
```go
package main

import (
	"fmt"

	"github.com/Greyh4t/nocomment"
)

func main() {
	s := &nocomment.Stripper{}
	// remove /* */
	s.RemoveBlockComment = true
	// remove //
	s.RemoveLineComment = true
	// remove #
	s.RemoveShellComment = true
	// remove <!-- -->
	s.RemoveHtmlComment = true

	fmt.Println(string(s.Clean([]byte(""))))
	fmt.Println(string(s.Clean([]byte("hello world"))))
	fmt.Println(string(s.Clean([]byte("//this is a comment\nHello World\n"))))
	fmt.Println(string(s.Clean([]byte("//this is a comment\r\nHello World\r\n"))))
	fmt.Println(string(s.Clean([]byte("//this is a comment\nHello World// another comment\n"))))
	fmt.Println(string(s.Clean([]byte("//this is a comment\r\nHello World// another comment\r\n"))))
	fmt.Println(string(s.Clean([]byte("#this is a comment\nHello World\n"))))
	fmt.Println(string(s.Clean([]byte("#this is a comment\r\nHello World\r\n"))))
	fmt.Println(string(s.Clean([]byte("#this is a comment\nHello World# another comment\r\n"))))
	fmt.Println(string(s.Clean([]byte("#this is a comment\r\nHello World# another comment\r\n"))))
	fmt.Println(string(s.Clean([]byte("//this is a comment\nHello World# another comment\n"))))
	fmt.Println(string(s.Clean([]byte("//this is a comment\r\nHello World# another comment\r\n"))))
	fmt.Println(string(s.Clean([]byte("/*this is a comment*/\nHello World\n"))))
	fmt.Println(string(s.Clean([]byte("/*this is a comment*/\r\nHello World\r\n"))))
	fmt.Println(string(s.Clean([]byte("/*this is a comment\n*/Hello World/* another comment*/\n"))))
	fmt.Println(string(s.Clean([]byte("/*this is a comment\r\n*/Hello World/* another comment*/\r\n"))))
	fmt.Println(string(s.Clean([]byte("/*this\n is a\n comment\n*/Hello World\n"))))
	fmt.Println(string(s.Clean([]byte("/*this\r\n is a\r\n comment\r\n*/Hello World\r\n"))))
	fmt.Println(string(s.Clean([]byte(`This is some text. "#This is not a comment // neither is this /* or this */" sooo, no comments!`))))
	fmt.Println(string(s.Clean([]byte("/* this is a broken block comment"))))
	fmt.Println(string(s.Clean([]byte("\" this is an unlcosed quote"))))
}

```