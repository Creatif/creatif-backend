// from blog of http://angelonotes.blogspot.com/2015/09/golang-utf16-utf8.html
package main

import (
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {
	bs_UTF16LE, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), []byte("測試"))
	bs_UTF16BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder(), []byte("測試"))
	bs_UTF8LE, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), bs_UTF16LE)
	bs_UTF8BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), bs_UTF16BE)

	fmt.Printf("%s\n%s\n%s\n%s\n", bs_UTF16LE, bs_UTF16BE, bs_UTF8LE, bs_UTF8BE)
}
